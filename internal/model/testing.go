package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

// Тестовый Блог
func TestBlog(t *testing.T) *BlogItem{
	oid, _ := primitive.ObjectIDFromHex("5d2399ef96fb765873a24bae")
	return &BlogItem{
		ID: oid,
		AuthorID: "Sula Go",
		Content: "Some principles of 12 factor",
		Title: "How to control your environment",
	}
}
