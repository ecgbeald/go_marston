package server

import (
	"net"
	"net/http"

	marston "github.com/ecgbeald/marston/services/genproto/nodes"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Master struct {
	api         *gin.Engine
	listener    net.Listener
	server      *grpc.Server
	nodeService *NodeServiceGprcServer
}

func (m *Master) Init() (err error) {
	m.listener, err = net.Listen("tcp", ":8888")
	if err != nil {
		return err
	}
	m.server = grpc.NewServer()
	m.nodeService = GetNodeServiceGrpcServer()
	marston.RegisterNodeServiceServer(m.server, m.nodeService)

	m.api = gin.Default()
	// get tasks form tasks endpoint
	m.api.POST("/tasks", func(ctx *gin.Context) {
		var payload struct {
			Cmd string `json:"cmd"`
		}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// the most important line! whack the payload into the channel!
		m.nodeService.channel <- payload.Cmd
		ctx.AbortWithStatus(http.StatusOK)
	})
	return nil
}

func (m *Master) Start() {
	go m.server.Serve(m.listener)
	_ = m.api.Run(":9092")
	m.server.Stop()
}

var master *Master

func GetMasterNode() *Master {
	if master == nil {
		master = &Master{}

		if err := master.Init(); err != nil {
			panic(err)
		}
	}
	return master
}
