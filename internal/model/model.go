package model

type Order struct {
	OrderID    string `bson,json:"orderid"`
	ProductID  string `bson,json:"productid"`
	ConsumerID string `bson,json:"consumerid"`
	VendorID   string `bson,json:"vendorid"`
	Amount     uint   `bson,json:"amount"`
	OrderPrice uint   `bson,json:"orderprice"`
}
