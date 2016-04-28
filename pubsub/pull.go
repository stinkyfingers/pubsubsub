package pubsub

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"

	"io/ioutil"
)

var (
	ctx context.Context

	ProjectID    = "cp100-john"
	Subscription = "test"
	CallbackURL  = "" //To receive notifications of new messages in the queue, specify an endpoint callback URL. If endpoint is an empty string the backend will not notify the client of new messages.
)

// Pull gets a single message from the PUBSUB queue
func Pull(topic string) error {
	var err error
	if ctx == nil {
		err = Context()
		if err != nil {
			return err
		}
	}
	pubsub.CreateSub(ctx, Subscription, topic, 0, CallbackURL)

	msgs, err := pubsub.Pull(ctx, Subscription, 1)
	if err != nil {
		return err
	}
	if len(msgs) < 1 {
		return nil
	}
	err = pubsub.Ack(ctx, Subscription, msgs[0].AckID)
	if err != nil {
		return err
	}
	return handleMessage(msgs[0])
}

func PullWait(topic string) error {
	var err error
	if ctx == nil {
		err = Context()
		if err != nil {
			return err
		}
	}
	pubsub.CreateSub(ctx, Subscription, topic, 0, CallbackURL)

	for {
		msgs, err := pubsub.PullWait(ctx, Subscription, 1)
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			err = pubsub.Ack(ctx, Subscription, msg.AckID)
			if err != nil {
				return err
			}
			go func(message *pubsub.Message) error {
				err = handleMessage(message)
				if err != nil {
					return err
				}
				return nil
			}(msg)
			if err != nil {
				break
			}
		}
	}
	return err
}

func Context() error {
	jsonKey, err := ioutil.ReadFile("keys/cp100-f39fd3c5c9f5.json")
	if err != nil {
		return err
	}
	conf, err := google.JWTConfigFromJSON(jsonKey, pubsub.ScopeCloudPlatform, pubsub.ScopePubSub)
	if err != nil {
		return err
	}
	ctx = cloud.NewContext(ProjectID, conf.Client(oauth2.NoContext))
	return nil
}
