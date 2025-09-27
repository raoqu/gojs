package main

import (
	"log"

	"github.com/raoqu/gojs"
)

func main() {
	config := gojs.LoadConfig("gojs.yaml")
	instance := gojs.CreateInstance(config)
	instance.Run()

	instance.Update("test",
		`
		let data = ""
		console.log(text)
		console.log(user)
		data = "" + user + ": " + text
		data // return value
	`)
	params := map[string]interface{}{
		"text": "text",
		"user": "user1",
	}
	result, err := instance.Execute("test", params)
	if err == nil {
		log.Printf("Result: %v", result)
	} else {
		log.Printf("Error: %v", err)
	}
}
