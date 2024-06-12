package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))

	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	fmt.Println("Starting Cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	fmt.Println("Starting Cube manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)

	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()

	mapi.Start()

	for i := 0; i < 3; i++ {
		uuidString := uuid.New()
		t := task.Task{
			ID:    uuidString,
			Name:  fmt.Sprintf(("test-container-%d-%s"), i, uuidString.String()[:3]),
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID:    uuid.New(),
			State: task.Running,
			Task:  t,
		}
		m.AddTask(te)

		// Original code was just the m.SendWork() call. But it itermittently failed to send the work to the worker
		// So I added the for loop to keep trying until the work is sent and the sleep to give the worker time to process the task
		for m.Pending.Len() != 0 {
			m.SendWork()
			time.Sleep(1 * time.Second)
		}
	}

}
