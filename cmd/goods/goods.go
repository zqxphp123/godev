package main

import (
	"math/rand"
	"mydev/app/goods/srv"
	"os"
	"runtime"
	"time"
)

func main() {
	randSrc := rand.NewSource(time.Now().UnixNano())
	rand.New(randSrc)
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	srv.NewApp("goods-server").Run()
}
