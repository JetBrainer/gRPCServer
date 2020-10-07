package model

import "go.mongodb.org/mongo-driver/mongo"

func New(db *mongo.Client) *BlogServiceServer{
	return &BlogServiceServer{
		db: db.Database("mydb").Collection("blogtest"),
	}
}
