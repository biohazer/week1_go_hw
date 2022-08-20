package main

import (
	"awesomeProject1/service"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	s1 := service.NewServer("business", "localhost:8080")
	s1.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello 8080 business serv"))
	}))

	s2 := service.NewServer("admin", "localhost:8081")
	s2.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("8082 admin serv"))
	}))
	// 这里魔改了一下，实在不知道这个type Option要怎么往后传
	app := service.NewApp([]*service.Server{s1, s2}, service.WithShutdownCallbacks(StoreCacheToDBCallback))
	app.StartAndServe()
}

func StoreCacheToDBCallback(ctx context.Context) {
	done := make(chan struct{}, 1)
	go func() {
		// 你的业务逻辑，比如说这里我们模拟的是将本地缓存刷新到数据库里面
		// 这里我们简单的睡一段时间来模拟
		// 需求：Redis key过期才持久化到MySQL
		log.Printf("刷新缓存中……")
		done <- struct{}{}
		time.Sleep(2 * time.Second)
	}()
	select {
	case <-ctx.Done():
		log.Printf("刷新缓存超时")
	case <-done:
		log.Printf("缓存被刷新到了 DB")
	}
}
