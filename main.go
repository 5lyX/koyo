package main

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"koyo/task"
	"koyo/worker"
	"time"
)

//
//import (
//	"fmt"
//	"github.com/docker/docker/client"
//	"github.com/golang-collections/collections/queue"
//	"github.com/google/uuid"
//	"koyo/manager"
//	"koyo/node"
//	"koyo/task"
//	"koyo/worker"
//	"os"
//	"time"
//)
//
//func main() {
//	t := task.Task{
//		ID:     uuid.New(),
//		Name:   "Task-1",
//		State:  task.Pending,
//		Image:  "Image-1",
//		Memory: 1024,
//		Disk:   1,
//	}
//
//	te := task.TaskEvent{
//		ID:        uuid.New(),
//		State:     task.Pending,
//		Timestamp: time.Now(),
//		Task:      t,
//	}
//
//	fmt.Printf("task: %v\n", t)
//	fmt.Printf("task event: %v\n", te)
//
//	w := worker.Worker{
//		Name:  "worker-1",
//		Queue: *queue.New(),
//		Db:    make(map[uuid.UUID]*task.Task),
//	}
//	fmt.Printf("worker: %v\n", w)
//	w.CollectStats()
//	w.RunTask()
//	w.StartTask()
//	w.StopTask()
//
//	m := manager.Manager{
//		Pending: *queue.New(),
//		TaskDb:  make(map[string][]*task.Task),
//		EventDb: make(map[string][]*task.TaskEvent),
//		Workers: []string{w.Name}, //defines a slice of strings and initializes it
//	}
//
//	fmt.Printf("manager: %v\n", m)
//	m.SelectWorker()
//	m.UpdateTasks()
//	m.SendWork()
//
//	n := node.Node{
//		Name:   "Node-1",
//		Ip:     "192.168.1.1",
//		Cores:  4,
//		Memory: 1024,
//		Disk:   25,
//		Role:   "worker",
//	}
//
//	fmt.Printf("node: %v\n", n)
//
//	fmt.Printf("create a test container\n")
//	dockerTask, createResult := createContainer()
//	if createResult.Error != nil {
//		fmt.Printf("%v\n", createResult.Error)
//		os.Exit(1)
//	}
//
//	time.Sleep(time.Second * 5)
//	fmt.Printf("stopping container %s\n", createResult.ContainerId)
//	_ = stopContainer(dockerTask, createResult.ContainerId)
//}
//
//func createContainer() (*task.Docker, *task.DockerResult) {
//	c := task.Config{
//		Name:  "test-container-1",
//		Image: "postgres:13",
//		Env: []string{
//			"POSTGRES_USER=koyo",
//			"POSTGRES_PASSWORD=secret",
//		},
//	}
//
//	dc, _ := client.NewClientWithOpts(client.FromEnv)
//	d := task.Docker{
//		Client: dc,
//		Config: c,
//	}
//
//	result := d.Run()
//	if result.Error != nil {
//		fmt.Printf("%v\n", result.Error)
//		return nil, nil
//	}
//
//	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, c)
//	return &d, &result
//}
//
//func stopContainer(d *task.Docker, id string) *task.DockerResult {
//	result := d.Stop(id)
//	if result.Error != nil {
//		fmt.Printf("%v\n", result.Error)
//		return nil
//	}
//
//	fmt.Printf("Container %s has been stopped and removed\n", result.ContainerId)
//	return &result
//}

func main() {
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "task-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	// first time the worker will see the task
	fmt.Println("starting task")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerId
	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(30 * time.Second)

	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
