package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type rqCtxKey string

const rqShardingIdKey rqCtxKey = "shardingId"

func main() {
	cfg := &PoolConfig{
		connNum:      1,
		addrs:        []string{":9090"},
		idleDuration: 5 * time.Minute,
		timeout:      5000 * time.Millisecond,
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx, stop := context.WithCancel(context.Background())
	ctxSignal, sigStop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	myClientPool := NewClientPool(ctx, cfg)

	ticker := time.NewTicker(time.Second * 5)
	var No int64
	for {
		select {
		case <-ctxSignal.Done():
			log.Printf("Server shutting down due to signal: %v\n", context.Cause(ctxSignal))
			myClientPool.ShutDown()
			stop()
			sigStop()
			log.Println("client exit...")
			return
		case <-ticker.C:
			rqCtx := context.WithValue(ctx, rqShardingIdKey, No)
			rq := &request{Msg: "hello thrift", No: No}
			data, err := Marshal(rq)
			if err != nil {
				log.Fatal("marshal request fail", err)
			}

			rs, err := myClientPool.TwoWay(rqCtx, 1001, data)
			if err != nil {
				log.Println("work TwoWay fail", No, err)
			} else {
				log.Println("work TwoWay success", No, "rs", string(rs.Content), "err", err)
			}

			err = myClientPool.oneWay(rqCtx, 1002, data)
			log.Println("work oneWay", No, "err", err)

			No++
			n := rand.Intn(11) + 5
			ticker.Reset(time.Duration(n) * time.Second)
		}
	}
}

// todo send work put in goroutine

type request struct {
	Msg string
	No  int64
}
