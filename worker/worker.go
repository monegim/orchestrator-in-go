package worker

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"cube/task"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]task.Task
	TaskCount int
}
