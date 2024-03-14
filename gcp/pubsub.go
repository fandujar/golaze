package gcp

type PubSubClientConfig struct {
	ProjectID string
	TopicID   string
}

type PubSubClient struct {
	*PubSubClientConfig
}

func NewPubSubClient(config *PubSubClientConfig) *PubSubClient {
	return &PubSubClient{
		config,
	}
}

func (client *PubSubClient) Publish(message string) error {
	return nil
}

func (client *PubSubClient) Subscribe() (chan string, error) {
	return nil, nil
}
