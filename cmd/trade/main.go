package main

import (
	"encoding/json"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/infra/kafka"
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/market/dto"
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/market/entity"
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/market/transformer"
	"sync"
)

// use flag --add-host=host.docker.internal:host-gateway
func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})
	go kafka.Consume(kafkaMsgChan) //T2

	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() //T3

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(outputJson))
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
