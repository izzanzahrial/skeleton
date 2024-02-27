package main

import (
	"context"
	"log"

	"github.com/izzanzahrial/skeleton/config"
	"github.com/izzanzahrial/skeleton/internal/domain/post/broker"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("failed to load environment variables")
	}

	cfg, err := config.NewConsumer()
	if err != nil {
		log.Panicf("failed to get kafka consumer config: %v", err)
	}

	consumer, err := broker.NewConsumer()
	if err != nil {
		log.Panicf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	handler := broker.Handler{}
	for {
		if err := consumer.Consume(context.Background(), cfg.Topics, handler); err != nil {
			log.Fatalf("error from consumer: %v", err)
		}
	}
}
