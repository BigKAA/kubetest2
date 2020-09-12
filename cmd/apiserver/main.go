package main

import (
	"log"

	"github.com/BigKAA/kubetest2/internal/app/apiserver"
)

func main() {
	config := apiserver.NewConfig()
	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
