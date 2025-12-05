package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github/FrodoDev/tthrift/gen-go/trpc"

	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/protobuf/proto"
)

// var gClientPool = new(clientPool)

type PoolConfig struct {
	connNum           int16         // 每个地址的连接数
	addrs             []string      // 连接地址，可能有多个
	idleDuration      time.Duration // 超过空闲时长后关闭链接
	timeout           time.Duration // connect timeout
	reconnectDuration time.Duration // 重连间隔
}

// todo NewDevelopCfg, NewProductionCfg

type ThriftClient struct {
	transport     thrift.TTransport
	client        *trpc.TServiceClient
	needReconnect bool
	addr          string
}

type clientPool struct {
	ctx           context.Context
	cfg           *PoolConfig
	clients       chan *ThriftClient
	reconnectChan chan string // todo struct{tradeNo, addr, times}
}

func NewClientPool(ctx context.Context, cfg *PoolConfig) *clientPool {
	clientPool := &clientPool{ctx: ctx}
	clientPool.clients = make(chan *ThriftClient, len(cfg.addrs)*int(cfg.connNum))
	clientPool.cfg = cfg
	clientPool.reconnectChan = make(chan string, len(cfg.addrs)*int(cfg.connNum))
	if clientPool.cfg.reconnectDuration < 50*time.Millisecond {
		clientPool.cfg.reconnectDuration = 50 * time.Millisecond
	}
	go clientPool.reconnect()

	for _, addr := range cfg.addrs {
		for i := int16(0); i < cfg.connNum; i++ {
			client, err := clientPool.newClient(addr)
			if err != nil {
				log.Println("connect fail", addr, err)
				clientPool.reconnectChan <- addr
				continue
			}

			clientPool.putClient(client)
		}
	}

	return clientPool
}

func (c *clientPool) ShutDown() {
	close(c.clients)
	for client := range c.clients {
		if client != nil && client.transport != nil {
			client.needReconnect = false
			client.transport.Close()
		}
	}
	log.Println("clientPool ShutDown done")
}

func (c *clientPool) newClient(addr string) (tclient *ThriftClient, err error) {
	// trpc.NewTServiceClient()
	sock := thrift.NewTSocketConf(addr,
		&thrift.TConfiguration{ConnectTimeout: 500 * time.Millisecond, SocketTimeout: 500 * time.Millisecond})

	transport := thrift.NewTFramedTransport(sock)
	if err := transport.Open(); err != nil {
		log.Println("Error opening transport:", err)
	}
	// defer transport.Close()

	protocol := thrift.NewTBinaryProtocolFactoryDefault().GetProtocol(transport)
	client := trpc.NewTServiceClient(thrift.NewTStandardClient(protocol, protocol))
	log.Println("new client", addr)
	return &ThriftClient{transport: transport, client: client, addr: addr}, nil
}

func (c *clientPool) reconnect() {
	ticker := time.NewTicker(c.cfg.reconnectDuration)
	for {
		select {
		case <-c.ctx.Done():
			log.Printf("clientPool reconnect shutting down due to ctx: %v\n", context.Cause(c.ctx))
			return
		case <-ticker.C:
			addr := <-c.reconnectChan
			newClient, err := c.newClient(addr)
			if err != nil {
				log.Println("reconnect fail", addr, err)
				c.reconnectChan <- addr
				continue
			}
			c.putClient(newClient)
		}
	}
}

func (c *clientPool) getClient() *ThriftClient {
	select {
	case conn := <-c.clients:
		if conn.needReconnect {
			c.reconnectChan <- conn.addr
			return nil
		}
		return conn
	case <-time.After(time.Millisecond * 500):
		log.Println("getClient timeout")
		return nil
	}
}

func (c *clientPool) putClient(client *ThriftClient) {
	if client.needReconnect {
		c.reconnectChan <- client.addr
		return
	}
	c.clients <- client
}

func (c *clientPool) TwoWay(ctx context.Context, rqType int32, rqContent []byte) (res *trpc.TResponse, err error) {
	client := c.getClient()
	if client == nil {
		err = fmt.Errorf("TwoWay getClient timeout")
		return
	}

	defer func() {
		if needReconnectErr(err) {
			client.needReconnect = true
		}
		c.putClient(client)
	}()

	shId, ok := ctx.Value(rqShardingIdKey).(int64)
	if !ok {
		err = fmt.Errorf("invalid shardingId:%v", ctx.Value(rqShardingIdKey))
		return
	}

	req := &trpc.TRequest{
		Type:       thrift.Int32Ptr(rqType),
		Content:    rqContent,
		ShardingId: thrift.Int64Ptr(shId),
	}
	log.Println("send twoWay msg")
	return client.client.TwoWay(ctx, req)
}

func needReconnectErr(err error) bool {
	log.Println("needReconnectErr ", err)
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "write: broken pipe") {
		return true
	}

	if strings.Contains(err.Error(), "i/o timeout") {
		return true
	}

	if strings.Contains(err.Error(), "Connection not open") {
		return true
	}
	return false
}

func (c *clientPool) oneWay(ctx context.Context, rqType int32, rqContent []byte) (err error) {
	client := c.getClient()
	if client == nil {
		err = fmt.Errorf("oneWay getClient timeout")
		return
	}

	defer func() {
		if needReconnectErr(err) {
			client.needReconnect = true
		}
		c.putClient(client)
	}()

	shId, ok := ctx.Value(rqShardingIdKey).(int64)
	if !ok {
		err = fmt.Errorf("invalid shardingId:%v", ctx.Value("shardingId"))
		return
	}

	req := &trpc.TRequest{
		Type:       thrift.Int32Ptr(rqType),
		Content:    rqContent,
		ShardingId: thrift.Int64Ptr(shId),
	}
	log.Println("send oneWay msg")
	return client.client.OneWay(ctx, req)
}

func Marshal(rq any) ([]byte, error) {
	if rq == nil {
		return nil, nil
	}

	if msg, ok := rq.(proto.Message); ok {
		return proto.Marshal(msg)
	}

	v := reflect.ValueOf(rq)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		return json.Marshal(rq)
	}

	return nil, fmt.Errorf("invalid type:%v", reflect.TypeOf(rq).Name())
}
