package main

import (
	"context"
	"fmt"
	"github/FrodoDev/tthrift/gen-go/trpc"
	"log"
	"sync"
)

// 实现 TService 接口
type THandler struct{}

func (h *THandler) TwoWay(ctx context.Context, req *trpc.TRequest) (*trpc.TResponse, error) {
	log.Printf("Received TwoWay request: Type=%d, ShardingId=%d\n", req.GetType(), req.GetShardingId())

	v, ok := typeMachineTwoWay.Load(req.GetType())
	if !ok {
		return nil, fmt.Errorf("not found processor")
	}
	method, ok1 := v.(twoWayProcessor)
	if !ok1 {
		return nil, fmt.Errorf("invalid registe processor")
	}

	return method(ctx, req)
}

func (h *THandler) OneWay(ctx context.Context, req *trpc.TRequest) error {
	log.Printf("Received OneWay request: Type=%d\n", req.GetType())

	v, ok := typeMachineOneWay.Load(req.GetType())
	if !ok {
		return fmt.Errorf("not found processor")
	}
	method, ok1 := v.(oneWayProcessor)
	if !ok1 {
		return fmt.Errorf("invalid registe processor")
	}

	return method(ctx, req)
}

type twoWayProcessor func(ctx context.Context, req *trpc.TRequest) (*trpc.TResponse, error)
type oneWayProcessor func(ctx context.Context, req *trpc.TRequest) error

// var typeMachine map[int32]processor
var (
	typeMachineTwoWay sync.Map
	typeMachineOneWay sync.Map
)

func RegisteTypeMachineTwoWay(rqType int32, method twoWayProcessor) {
	_, loaded := typeMachineTwoWay.LoadOrStore(rqType, method)
	if loaded {
		log.Println("repeated RegisteTypeMachineTwoWay", rqType)
	}
}

func RegisteTypeMachineOneWay(rqType int32, method oneWayProcessor) {
	_, loaded := typeMachineOneWay.LoadOrStore(rqType, method)
	if loaded {
		log.Println("repeated RegisteTypeMachineOneWay", rqType)
	}
}
