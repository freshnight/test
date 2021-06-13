package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = ":9527"
const addr2 = ":9528"

func main() {

	exit := make(chan os.Signal)
	//监听 信号
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	var g errgroup.Group
	mux := http.NewServeMux()
	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "just another http server  1...")
	})

	var srv = http.Server{
		Addr:    addr,
		Handler: mux,
	}

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "just another http server  2...")
	})
	var srv2 = http.Server{
		Addr:    addr2,
		Handler: mux2,
	}
	g.Go(func() error {

		fmt.Println("listening 1 at " + addr)
		err := srv.ListenAndServe()

		fmt.Println("waiting for 1 the remaining connections to finish...")

		if err != nil && err != http.ErrServerClosed {
			return err
		}
		fmt.Println("gracefully 1 shutdown the http server...")
		return nil
	})

	g.Go(func() error {

		fmt.Println("listening 2 at " + addr2)
		err := srv2.ListenAndServe()

		fmt.Println("waiting for 2 the remaining connections to finish...")

		if err != nil && err != http.ErrServerClosed {
			return err
		}
		fmt.Println("gracefully 2  shutdown the http server...")
		return nil
	})
	g.Go(func() error {
		<-exit
		//使用context控制srv.Shutdown的超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		//关闭1
		err := srv.Shutdown(ctx)

		if err != nil {

			return err
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel2()
		//关闭2
		err = srv2.Shutdown(ctx2)
		if err != nil {

			return err
		}
		return nil
	})
	err := g.Wait()
	if err != nil {

		fmt.Println("error ", err)
	} else {

		fmt.Println("Successfully ")
	}

}
