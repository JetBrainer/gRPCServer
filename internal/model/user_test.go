package model_test

import (
	"context"
	"github.com/JetBrainer/gRPCServer/internal/model"
	blogpb "github.com/JetBrainer/gRPCServer/internal/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)
var databaseURL = "mongodb://127.0.0.1:27017"

func TestBlogServiceServer_CreateBlog(t *testing.T) {
	db, teardown := model.TestDB(t, databaseURL)
	defer teardown("blogtest")

	s := model.New(db)

	u := model.TestBlog(t)
	_, err := s.CreateBlog(context.Background(),&blogpb.CreateBlogReq{Blog: &blogpb.Blog{
		AuthorId: u.AuthorID,
		Content: u.Content,
		Title: u.Title,
	}})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
