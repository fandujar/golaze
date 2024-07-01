package golaze

type TaskQueue struct {
	tasks   []Task
	enqueue chan Task
	dequeue chan Task
}
