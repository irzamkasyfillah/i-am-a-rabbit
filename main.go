package main

import (
	"context"
	"encoding/json"
	"irzam/rabbitmq/rabbitmq"
)

func main() {
	body := map[string]interface{}{
		"message": "Hello World",
	}
	bodyJson, _ := json.Marshal(body)
	rabbitmq.HandlePub(context.Background(), "", "hello", string(bodyJson))
	// rabbitmq.HandlePub(context.Background(), "", "hell", "Hell dian")
	rabbitmq.HandleWorker()
}
