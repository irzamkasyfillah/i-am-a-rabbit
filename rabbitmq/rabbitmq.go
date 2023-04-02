package rabbitmq

import (
	"context"
	"irzam/rabbitmq/helper"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ch *amqp.Channel
var conn *amqp.Connection
var err error

func Connect() {
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	helper.FailOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ")
}

func OpenChannel() {
	ch, err = conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
}

func DeclareQueue(queue string) {
	_, err = ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	helper.FailOnError(err, "Failed to declare a queue")
}

func Publish(ctx context.Context, exchange string, route string, content string) {
	err := ch.PublishWithContext(context.Background(),
		exchange, // exchange
		route,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	helper.FailOnError(err, "Failed to publish a message")
}

// HandlePub is a function to send message to queue with params: exchange, route, content
func HandlePub(ctx context.Context, exchange string, route string, content string) {
	Connect()
	OpenChannel()
	DeclareQueue(route)
	Publish(ctx, exchange, route, content)
	defer ch.Close()
	defer conn.Close()
}

// HandleWorker is a function to create worker and consume all defined queue.
// Add more queue to consume using: go Consume("queue_name")
func HandleWorker() {
	Connect()
	OpenChannel()

	// consumer
	go Consume("hello")
	go Consume("hell")

	defer ch.Close()
	defer conn.Close()
	// create channel to block main goroutine from exiting
	forever := make(chan struct{})
	<-forever
}

func Consume(queue string) {
	DeclareQueue(queue)
	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	helper.FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		switch queue {
		case "hell":
			log.Printf("Received a message from hell: %s", d.Body)
		case "hello":
			// var body map[string]interface{}
			// json.Unmarshal(d.Body, &body)
			log.Printf("Received a message from hello: %s", d.Body)
		}
	}
}
