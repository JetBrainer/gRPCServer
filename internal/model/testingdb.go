package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

// Тестовый Сторедж

func TestDB(t *testing.T,databaseURL string) (*mongo.Client, func(...string)){
	t.Helper()
	db, err := mongo.Connect(context.Background(),options.Client().ApplyURI(databaseURL))
	if err !=nil{
		t.Fatal(err)
	}
	if err = db.Ping(context.Background(),nil); err != nil{
		t.Fatal(err)
	}
	//testdbse := db.Database("mydb").Collection("blogtest")
	return db, func(tables ...string) {
		if len(tables) > 0{
			if  err := db.Database("mydb").Collection("blogtest").Drop(context.Background()); err != nil{
				t.Fatal(err)
			}
			if err = db.Disconnect(context.Background()); err != nil{
				t.Fatal(err)
			}
		}
	}
}