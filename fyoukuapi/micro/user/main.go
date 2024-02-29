package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"time"
)

func main() {

	//服务端代码
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"http://127.0.0.1:2379"}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.srv.fyoukuApi.user"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	service.Init()
	//proto.RegisterUserServiceHandler(service.Server(), controllers.UserRpcController)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
