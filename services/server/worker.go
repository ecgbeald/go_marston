package server

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	marston "github.com/ecgbeald/marston/services/genproto/nodes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Worker struct {
	conn *grpc.ClientConn
	c    marston.NodeServiceClient
}

func (w *Worker) Init() (err error) {
	fmt.Println("Connecting to port 8888")
	w.conn, err = grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	w.c = marston.NewNodeServiceClient(w.conn)
	return nil
}

func (w *Worker) Start() {
	fmt.Println("worker started")
	// return status of worker to master (just ok)
	_, _ = w.c.ReturnStatus(context.Background(), &marston.Request{})

	stream, _ := w.c.AssignTask(context.Background(), &marston.Request{})
	// main loop
	for {
		res, err := stream.Recv()
		if err != nil {
			return
		}
		fmt.Println("received command:", res.Data)
		parts := strings.Split(res.Data, " ")
		if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
			fmt.Println(err)
		}
	}
}

var worker *Worker

func GetWorker() *Worker {
	if worker == nil {
		worker = &Worker{}
		if err := worker.Init(); err != nil {
			panic(err)
		}
	}
	return worker
}
