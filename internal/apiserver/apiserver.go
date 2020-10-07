package apiserver

import (
	"context"
	"github.com/JetBrainer/gRPCServer/internal/model"
	blogpb "github.com/JetBrainer/gRPCServer/internal/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start(config *Config) (net.Listener,*grpc.Server,*mongo.Client){
	db := NewDB(config.DatabaseURL)

	log.Println("Start Listen our api")
	listener, err := net.Listen("tcp",":"+config.BindAddr)
	if err != nil{
		log.Fatalf("Unable to listen api port: %v",err)
	}

	// Опции в случае необходимости поддержки TLS
	opts := []grpc.ServerOption{}

	// Создаем gRPC сервер
	s := grpc.NewServer(opts...)

	// Тип Блога
	srv := &model.BlogServiceServer{}

	// Регистрируем Сервер с сервисом
	blogpb.RegisterBlogServiceServer(s, srv)

	go func() {
		if err := s.Serve(listener); err != nil{
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	log.Println("Server successfully started")

	return listener, s, db
}

func NewDB(databaseURL string) *mongo.Client{
	// Подключение к базе данных
	ctx := context.Background()
	db, err := mongo.Connect(ctx,options.Client().ApplyURI(databaseURL))
	if err != nil{
		log.Fatal("Connection ERROR")
	}
	// Пингуем для проверки успешного подключения
	err = db.Ping(ctx, nil)
	if err != nil{
		log.Fatal("Ping ERROR DB")
	}
	log.Println("Connected to MongoDB")

	return db
}
