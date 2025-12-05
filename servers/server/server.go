// server.go
package main

import (
	"context"
	"log"

	"github/FrodoDev/tthrift/gen-go/trpc" // 假设 trpc 在当前目录下；也可用 module 路径

	"github.com/apache/thrift/lib/go/thrift"
)

type TServer struct {
	ctx    context.Context
	server *thrift.TSimpleServer
}

func NewServer(ctx context.Context) (*TServer, error) {
	server := &TServer{}
	server.ctx = ctx
	newSvr, err := newThriftServer()
	if err != nil {
		return nil, err
	}
	server.server = newSvr
	return server, nil
}

func newThriftServer() (*thrift.TSimpleServer, error) {
	handler := &THandler{}
	processor := trpc.NewTServiceProcessor(handler)
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	addr := ":9090"
	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Fatal("Unable to create server socket:", err)
		return nil, err
	}

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	log.Println("Starting Thrift server on", addr)
	return server, nil
}

func (t *TServer) Start() error {
	return t.server.Serve()
}

func (t *TServer) Stop() {
	if t.server == nil {
		return
	}

	transport := t.server.ServerTransport()
	if transport == nil {
		return
	}

	t.server.ServerTransport().Close()
}
