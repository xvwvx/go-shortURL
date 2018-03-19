package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	"shortURL/core"
	"os"
)

func main() {

	if err := core.LoadConfg(); err != nil {
		panic(err)
	}

	if runtime.NumCPU() < core.Conf.MaxProc {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(core.Conf.MaxProc)
	}

	writer, err := os.OpenFile("./shortURL.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer writer.Close()
	log.SetOutput(writer)

	core.LinkDB()

	http.HandleFunc("/", core.Root)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		err := http.ListenAndServe(core.Conf.HTTPAddr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		waitGroup.Done()
	}()

	log.Println("ShortURL 启动完成")
	waitGroup.Wait()
	log.Println("ShortURL 程序退出")

}
