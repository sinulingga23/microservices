package model

type Product struct {
	Id          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Qtty        int      `bson:"qtty"`
	Price       int      `bson:"price"`
	Description string   `bson:"description"`
	CategoryIds []string `bson:"categoryIds"`
}

type ProductQttyData struct {
	Id   string `bson:"_id"`
	Qtty int    `bson:"qtty"`
}

type OrderProduct struct {
	OrderId   string `bson:"orderId"`
	ProductId string `bson:"productId"`
	Qtty      int    `bson:"qtty"`
}
