package golaze

type WorkerConfig struct {
	EventBus *EventBus
}

type Worker struct {
	*WorkerConfig
}

func NewWorker(config *WorkerConfig) *Worker {
	if config.EventBus == nil {
		config.EventBus = NewEventBus(&EventBusConfig{})
	}

	return &Worker{
		config,
	}
}
