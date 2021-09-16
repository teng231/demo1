package main

import (
	"net"
	"time"

	"github.com/teng231/demo1/db"
	"github.com/teng231/demo1/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type IDatabase interface {
	ListUsers(rq *pb.UserRequest) ([]*pb.User, error)
	CountUsers(rq *pb.UserRequest) (int64, error)
	FindUser(rq *pb.UserRequest) (*pb.User, error)
}

type Demo struct {
	db IDatabase
}

func initService() (*Demo, error) {
	d := &db.DB{}
	if err := d.ConnectDb("localhost:5432", "demo", "true"); err != nil {
		return nil, err
	}
	return &Demo{db: d}, nil
}

func GRPCServe(port string, handler *Demo) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 30 * time.Second,
		}),
	}
	serve := grpc.NewServer(opts...)
	pb.RegisterDemoServiceServer(serve, handler)
	reflection.Register(serve)
	return serve.Serve(listen)
}
