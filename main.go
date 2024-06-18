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

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	fmt.Println("Starting Koyo manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()

	mapi.Start()
}
