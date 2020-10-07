package internal

import "github.com/JetBrainer/gRPCServer/internal/model"

type Store interface {
	User()	model.BlogServiceServer
}
