package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

type PubSubClientConfig struct {
	ProjectID      string
	TopicID        string
	SubscriptionID string
}

type PubSubClient struct {
	*PubSubClientConfig
	client *pubsub.Client
}

func NewPubSubClient(config *PubSubClientConfig) (*PubSubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	return &PubSubClient{
		config,
		client,
	}, nil
}

func (client *PubSubClient) Publish(message string) error {
	return nil
}

func (client *PubSubClient) Subscribe(ctx context.Context) (chan *pubsub.Message, error) {
	sub := client.client.Subscription(client.SubscriptionID)

	// Create the subscription if it doesn't exist.
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := client.client.CreateSubscription(ctx, client.SubscriptionID, pubsub.SubscriptionConfig{
			Topic: client.client.Topic(client.TopicID),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription: %v", err)
		}
	}

	// Create a channel to pass messages received from Pub/Sub.
	msgCh := make(chan *pubsub.Message)

	// Start a goroutine to receive messages.
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			// Forward message data to channel and then acknowledge the message.
			select {
			case <-ctx.Done():
				return
			case msgCh <- msg:
				log.Info().Msgf("received message: %v", string(msg.Data))
			}
		})
		if err != nil {
			// Handling error, for example logging it. Not closing channel here to avoid panic in case of send on closed channel.
			fmt.Printf("error receiving messages: %v\n", err)
		}
	}()

	return msgCh, nil
}
