package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/clodomar/prova/configs"
	"github.com/clodomar/prova/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Client = configs.ConnectDB()
var clienteCollection *mongo.Collection = configs.GetCollection(DB, "clientes")

func main() {
	configs.ConnectDB()
	fmt.Println("RabbitMQ in Golang")

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		"queue_cliente",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			client := models.Cliente{}
			err := json.Unmarshal(msg.Body, &client)
			if err != nil {
				log.Printf("Failed to decode message body: %v", err)
				continue
			}
			if err := CreateClient(&client); err != nil {
				log.Printf("Failed to save client: %v", err)
				continue
			}
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
}

func CreateClient(client *models.Cliente) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := clienteCollection.InsertOne(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
