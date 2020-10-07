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

func TestBlogServiceServer_ReadBlog(t *testing.T) {
	db, teardown := model.TestDB(t, databaseURL)
	defer teardown("blogtest")

	s := model.New(db)

	u := model.TestBlog(t)
	some, err := s.CreateBlog(context.Background(),&blogpb.CreateBlogReq{Blog: &blogpb.Blog{
		AuthorId: u.AuthorID,
		Content: u.Content,
		Title: u.Title,
	}})
	assert.NoError(t, err)
	assert.NotNil(t, u)

	_, err = s.ReadBlog(context.Background(),&blogpb.ReadBlogReq{Id: some.GetBlog().Id})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestBlogServiceServer_DeleteBlog(t *testing.T) {
	db, teardown := model.TestDB(t, databaseURL)
	defer teardown("blogtest")

	s := model.New(db)

	u := model.TestBlog(t)
	some, err := s.CreateBlog(context.Background(),&blogpb.CreateBlogReq{Blog: &blogpb.Blog{
		AuthorId: u.AuthorID,
		Content: u.Content,
		Title: u.Title,
	}})
	assert.NoError(t, err)
	assert.NotNil(t, u)

	s1, err := s.DeleteBlog(context.Background(),&blogpb.DeleteBlogReq{Id: some.GetBlog().Id})

	assert.Equal(t, true, s1.Success)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestBlogServiceServer_UpdateBlog(t *testing.T) {
	db, teardown := model.TestDB(t, databaseURL)
	defer teardown("blogtest")

	s := model.New(db)

	u := model.TestBlog(t)
	some, err := s.CreateBlog(context.Background(),&blogpb.CreateBlogReq{Blog: &blogpb.Blog{
		AuthorId: u.AuthorID,
		Content: u.Content,
		Title: u.Title,
	}})
	assert.NoError(t, err)
	assert.NotNil(t, u)

	_, err = s.UpdateBlog(context.Background(),&blogpb.UpdateBlogReq{Blog: &blogpb.Blog{
		Id: some.GetBlog().Id,
		AuthorId: u.AuthorID,
		Content: u.Content,
		Title: u.Title,
	}})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestBlogServiceServer_ListBlogs(t *testing.T) {
	db, teardown := model.TestDB(t, databaseURL)
	defer teardown("blogtest")

	s := model.New(db)


	var server interface{
		blogpb.BlogService_ListBlogsServer
	}
	err := s.ListBlogs(&blogpb.ListBlogRequest{},server)

	assert.NoError(t, err)
}