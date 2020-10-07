package model

import (
	"context"
	"fmt"
	blogpb "github.com/JetBrainer/gRPCServer/internal/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type BlogItem struct {
	ID			primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID	string			   `bson:"author_id"`
	Content		string			   `bson:"content"`
	Title		string			   `bson:"title"`
}


type BlogServiceServer struct {
	db	*mongo.Collection
}

func (b *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()
	// Конертируем в БСОН
	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title: 	  blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	// Добавляем наш результат в базу
	result, err := b.db.InsertOne(ctx, data)
	if err != nil{
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}
	// Добавляем наш сгенеренный айди в блог
	oid := result.InsertedID.(primitive.ObjectID)
	// Конвертирование в Обджект Айди
	blog.Id = oid.Hex()
	// Возвращаем наш блог
	return &blogpb.CreateBlogRes{
		Blog: blog,
	}, nil
}

func (b *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	// Конвертирование данных от протофайл до uuid Mongo
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil{
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Could not format to ObejctID: %v", err))
	}
	result := b.db.FindOne(ctx,bson.M{"_id":oid})
	// Записываем полученные данные в структуру
	data := BlogItem{}
	if err := result.Decode(&data); err != nil{
		return nil, status.Errorf(codes.NotFound,
			fmt.Sprintf("Could not found blog with object ID %s:%v",req.GetId(),err))
	}
	// Отдаем записанный респонс
	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: data.AuthorID,
			Content: data.Content,
			Title: data.Content,
		},
	}
	return response, nil
}

func (b *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	// Получаем блог с запроса
	blog := req.GetBlog()
	// Конвертируем айди
	oid , err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil{
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Could not format to ObejctID: %v", err))
	}
	// Конвертируем данные в БСОН документ
	update := bson.M{
		"author_id":blog.GetAuthorId(),
		"title":	blog.GetTitle(),
		"content":	blog.GetContent(),
	}
	filter := bson.M{"_id":oid}
	//Обновляем данные в БСОН-е
	result := b.db.FindOneAndUpdate(ctx,filter,bson.M{"$set":update},options.FindOneAndUpdate().SetReturnDocument(1))

	// Декодируем для результата
	decoded := BlogItem{}
	if err = result.Decode(&decoded); err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find blog with supplied ID: %v", err),
		)
	}
	return &blogpb.UpdateBlogRes{Blog: &blogpb.Blog{
		Id: oid.Hex(),
		AuthorId: decoded.AuthorID,
		Content: decoded.Content,
		Title: decoded.Title,
	}},nil
}

func (b *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	// Берем Обдж Айди
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil{
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Could not convert to obj id %v", err))
	}
	// Удаляем
	_, err = b.db.DeleteOne(ctx, bson.M{"_id":oid})
	if err != nil{
		return nil, status.Errorf(codes.NotFound,"Could not find/delete blog with id %s:%v",req.GetId(),err)
	}
	// В случае успеха отправляем значение true
	return &blogpb.DeleteBlogRes{Success: true}, nil
}

func (b *BlogServiceServer) ListBlogs(request *blogpb.ListBlogRequest, server blogpb.BlogService_ListBlogsServer) error {
	// Для декода создаем структуру
	data := &BlogItem{}
	// Используем Find для возврата курсора
	cursor, err := b.db.Find(context.Background(),bson.M{})
	if err != nil{
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// Закрываем курсор в конце
	defer func() {
		if err := cursor.Close(context.Background()); err != nil{
			log.Println("Cursor close error")
		}
	}()
	// Пробегаемся по курсору
	for cursor.Next(context.Background()){
		// Декодим результаты в нашу созданную структуру
		if err := cursor.Decode(data); err != nil{
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		// Если ошибки не найдены отправляем наши блоги в поток
		if err = server.Send(&blogpb.ListBlogResponse{Blog:
			&blogpb.Blog{
			Id: data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content: data.Content,
			Title: data.Title,
			}}); err != nil{
			log.Println("Sending error", err)
		}
	}
	// проверяем курсор на ошибки
	if err := cursor.Err(); err != nil{
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

