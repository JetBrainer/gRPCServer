package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/JetBrainer/gRPCServer/internal/apiserver"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Будем парсить флаги в это значение
var configPath string

// Парсинг
func init() {
	flag.StringVar(&configPath,"config-path","configs/apiserver.toml","path to config file")
}

func main(){
	// Распарсим флаги
	flag.Parse()

	// Конфиг по умолчанию
	config := apiserver.NewConfig()

	// Декодим Томл Файл
	_, err := toml.DecodeFile(configPath,&config)
	if err != nil{
		log.Fatal("File Parsing error", err)
	}

	// Контекст для остановки
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Старт Сервера
	listener, serv, db := apiserver.Start(config)
	defer func() {
		if err := serv.Stop; err != nil{
			log.Println("Server stop ERROR ", err)
		}
	}()
	defer func() {
		if err := db.Disconnect(ctx); err != nil{
			log.Println("Mongo Disconnection ERROR ", err)
		}
	}()
	defer func() {
		if err := listener.Close(); err != nil{
			log.Println("Listener Closing ERROR ", err)
		}
	}()


	// Канал для сигнала для остановки
	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("\n Grace Stopping Server...")
}
