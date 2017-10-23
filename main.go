package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	"shortURL/core"
)

func main() {

	if err := core.LoadConfg(); err != nil {
		panic(err)
	}

	if runtime.NumCPU() < core.Conf.Procs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(core.Conf.Procs)
	}

	core.LinkDB()

	http.HandleFunc("/", core.GetOriginalURL)
	http.HandleFunc("/_genShortURL", core.GenShortURL)

	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		err := http.ListenAndServe(":1234", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		waitGroup.Done()
	}()

	log.Println("hortURL 启动完成")
	waitGroup.Wait()
	log.Println("ShortURL 程序退出")

}
