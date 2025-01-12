package worker

type Worker[T any] struct {
	taskChan chan T
}

func NewWorker[T any](bufferSize int) *Worker[T] {
	return &Worker[T]{
		taskChan: make(chan T, bufferSize),
	}
}

func (w *Worker[T]) Handle(tf func(t T)) {
	defer close(w.taskChan)
	for task := range w.taskChan {
		tf(task)
	}
}

func (w *Worker[T]) AddTask(task T) {
	w.taskChan <- task
}
