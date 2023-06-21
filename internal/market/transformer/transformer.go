package transformer

import (
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/market/dto"
	"github.com/soaresenzo/home-broker/bolsa-microservice/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)
	if input.CurrentShares > 0 {
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}
	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderID:    order.ID,
		AssetID:    order.Asset.ID,
		InvestorID: order.Investor.ID,
		Shares:     order.Shares,
		OrderType:  order.OrderType,
		Status:     order.Status,
		Partial:    order.PendingShares,
	}

	var transactionsOutput []*dto.TransactionOutput
	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: t.ID,
			BuyerID:       t.BuyingOrder.ID,
			SellerID:      t.SellingOrder.ID,
			AssetID:       t.SellingOrder.Asset.ID,
			Price:         t.Price,
			Shares:        t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	output.TransactionOutput = transactionsOutput
	return output
}
