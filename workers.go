package golaze

type WorkerConfig struct{}

type Worker struct {
	*WorkerConfig
}

func NewWorker(config *WorkerConfig) *Worker {
	return &Worker{
		config,
	}
}
