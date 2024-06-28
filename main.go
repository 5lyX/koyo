package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"koyo/manager"
	"koyo/task"
	"koyo/worker"
	"os"
	"strconv"
)

func main() {
	whost := os.Getenv("KOYO_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("KOYO_WORKER_PORT"))

	mhost := os.Getenv("KOYO_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("KOYO_MANAGER_PORT"))

	fmt.Println("Starting koyo worker")

	w1 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi1 := worker.Api{Address: whost, Port: wport, Worker: &w1}

	w2 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi2 := worker.Api{Address: whost, Port: wport + 1, Worker: &w2}

	w3 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi3 := worker.Api{Address: whost, Port: wport + 2, Worker: &w3}

	go w1.RunTasks()
	// go w.CollectStats()
	go w1.UpdateTasks()
	go wapi1.Start()

	go w2.RunTasks()
	// go w.CollectStats()
	go w2.UpdateTasks()
	go wapi2.Start()

	go w3.RunTasks()
	// go w.CollectStats()
	go w3.UpdateTasks()
	go wapi3.Start()

	fmt.Println("Starting Koyo manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport), fmt.Sprintf("%s:%d", whost, wport+1), fmt.Sprintf("%s:%d", whost, wport+2)}
	m := manager.New(workers, "epvm")
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()
	go m.DoHealthChecks()

	mapi.Start()
}
