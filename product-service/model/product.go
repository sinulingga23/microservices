package model

type Product struct {
	Id          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Qtty        int      `bson:"qtyy"`
	Price       int      `bson:"int"`
	Description string   `bson:"description"`
	CategoryIds []string `bson:"categoryIds"`
}
