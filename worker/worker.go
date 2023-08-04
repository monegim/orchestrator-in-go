package worker

import (
	"cube/task"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	Stats     *Stats
	TaskCount int
}

func (w *Worker) GetTasks() []*task.Task {
	var tasks []*task.Task
	for _, t := range w.Db {
		tasks = append(tasks, t)
	}
	return tasks
}
func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		fmt.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	taskQueued := t.(task.Task)
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}
	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("We should not get here")
		}
	} else {
		err := fmt.Errorf("Invalid transition from %v to %v", taskPersisted, taskQueued)
		result.Error = err
	}
	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	// Update the StartTime field on the task t
	t.StartTime = time.Now().UTC()
	// Create an instance of the Docker struct to talk to the Docker daemon
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	// Call the Run() method on the Docker struct
	result := d.Run()
	// Check the error
	if result.Error != nil {
		log.Printf("Err running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	//Add the id to the task
	result.ContainerId = t.ContainerID
	// Save the update task to the worker DB
	w.Db[t.ID] = &t
	return result
}
func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v", t.ContainerID, result.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v", t.ContainerID, t.ID)
	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = GetStats()
		w.TaskCount = w.Stats.TaskCount
		time.Sleep(time.Second * 15)
	}
}
