package model

type Category struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}
