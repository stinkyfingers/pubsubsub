package pubsub

import (
	"google.golang.org/cloud/pubsub"

	"log"
)

func handleMessage(msg *pubsub.Message) error {
	log.Print(string(msg.Data))
	return nil
}
