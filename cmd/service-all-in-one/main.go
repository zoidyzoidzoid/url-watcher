package main

import (
	_ "net/http/pprof"
	"sync"

	service "github.com/zoidbergwill/url-watcher/pkg/service"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.RunGRPCServer()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.RunGRPCGatewayServer()
	}()
	wg.Wait()
}
