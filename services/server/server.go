package server

import (
	"context"

	marston "github.com/ecgbeald/marston/services/genproto/nodes"
)

type NodeServiceGprcServer struct {
	marston.UnimplementedNodeServiceServer
	channel chan string
}

func (c NodeServiceGprcServer) ReturnStatus(ctx context.Context, in *marston.Request) (*marston.Response, error) {
	return &marston.Response{Data: "ok"}, nil
}

func (c NodeServiceGprcServer) AssignTask(req *marston.Request, server marston.NodeService_AssignTaskServer) error {
	for {
		select {
		case cmd := <-c.channel: // when there is stuff in the channel
			if err := server.Send(&marston.Response{Data: cmd}); err != nil {
				return err
			}
		}
	}
}

var server *NodeServiceGprcServer

func GetNodeServiceGrpcServer() *NodeServiceGprcServer {
	if server == nil {
		server = &NodeServiceGprcServer{
			channel: make(chan string),
		}
	}
	return server
}
