package golaze

import "github.com/rs/zerolog/log"

type TaskQueue struct {
	tasks   []Task
	enqueue chan Task
	dequeue chan Task
}

func (tq *TaskQueue) Start() {
	go func() {
		for {
			select {
			case task := <-tq.enqueue:
				log.Debug().Msgf("enqueuing task %s", task.Name)
				tq.tasks = append(tq.tasks, task)
			case tq.dequeue <- tq.Dequeue():
			}
		}
	}()
}

func (tq *TaskQueue) Dequeue() Task {
	if len(tq.tasks) == 0 {
		return Task{}
	}

	task := tq.tasks[0]
	log.Debug().Msgf("dequeuing task %s", task.Name)
	tq.tasks = tq.tasks[1:]
	return task
}
