package main

import (
	"github.com/raoqu/gojs"
)

func main() {
	config := gojs.LoadConfig("gojs.yaml")
	instance := gojs.CreateInstance(config)
	instance.Run()
	select {}
}
