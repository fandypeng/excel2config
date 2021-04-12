package rpcclient

import (
	"github.com/fandypeng/e2cdatabus/auth"
	pb "github.com/fandypeng/e2cdatabus/proto"
	"google.golang.org/grpc"
	"log"
)

type Conf struct {
	ServerAddr string
	AppKey     string
	AppSecret  string
}

func NewRpcClient(conf Conf) (client pb.DatabusClient, err error) {
	c, err := grpc.Dial(conf.ServerAddr, grpc.WithInsecure(), grpc.WithPerRPCCredentials(auth.New(conf.AppKey, conf.AppSecret)))
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	client = pb.NewDatabusClient(c)
	return
}
