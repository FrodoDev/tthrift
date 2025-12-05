package main

import (
	"context"
	"github/FrodoDev/tthrift/gen-go/trpc"
	"log"
	"math/rand/v2"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)

func Registe() {
	RegisteTypeMachineTwoWay(1001, MyTwoWay)
	RegisteTypeMachineOneWay(1002, MyOneWay)
}

func MyTwoWay(ctx context.Context, req *trpc.TRequest) (*trpc.TResponse, error) {
	n := rand.IntN(500) + 100
	time.Sleep(time.Millisecond * time.Duration(n))
	log.Printf("enter MyTwoWay request: Type=%d msg:%s sleep:%d Millisecond\n", req.GetType(), req.GetContent(), n)

	return &trpc.TResponse{
		ErrCode: thrift.Int64Ptr(0),
		Type:    req.Type,
		Content: []byte("echo: " + string(req.GetContent())),
		Trace:   req.Trace,
	}, nil
}

func MyOneWay(ctx context.Context, req *trpc.TRequest) error {
	log.Printf("enter MyOneWay request: Type=%d msg:%s \n", req.GetType(), req.GetContent())
	return nil
}

// todo
func Unmarshal() {

}
