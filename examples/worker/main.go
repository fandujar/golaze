package main

import (
	"fmt"

	"github.com/fandujar/golaze"
)

type CustomEventHandler struct{}

func (h *CustomEventHandler) Handle(event *golaze.Event) error {
	// log := golaze.NewLogger()
	// log.Info().Msgf("Event: %v", event)
	fmt.Printf("Event: %v\n", event)
	return nil
}

func main() {
	eventHandler := &CustomEventHandler{}

	worker := golaze.NewWorker(&golaze.WorkerConfig{
		EventBus: golaze.NewEventBus(&golaze.EventBusConfig{}),
	})

	worker.EventBus.Subscribe("custom", eventHandler)
	worker.EventBus.Publish("custom", &golaze.Event{Data: "Hello, World!"})
	worker.EventBus.Unsubscribe("custom", eventHandler)
	worker.EventBus.Publish("custom", &golaze.Event{Data: "Bye, World!"})
}
