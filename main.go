package main

import (
	"log"

	"one-cd/conf"
	"one-cd/http"
	"one-cd/service"
)

func main() {
	conf.Init()
	s := service.New()
	if err := s.Init(); err != nil {
		log.Println("Init error:", err)
	}
	http.Init(s)
	http.Start(conf.Conf.Listen)
}
