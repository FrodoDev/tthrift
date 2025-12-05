package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctxSignal, sigStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer sigStop()

	Registe()
	svr, err := NewServer(ctxSignal)
	if err != nil {
		log.Fatal("NewServer fail", err)
	}

	go func() {
		if err := svr.Start(); err != nil {
			// 服务器停止通常是正常的（因为 Stop() 关闭了 transport）
			log.Printf("Server stopped: %v", err)
		} else {
			log.Println("Server stopped normally")
		}
	}()

	log.Println("Server started, waiting for signal...")
	<-ctxSignal.Done()
	log.Printf("Server shutting down due to signal: %v\n", context.Cause(ctxSignal))

	svr.Stop()

	time.Sleep(200 * time.Millisecond)
	log.Println("Server shutdown complete")

}
