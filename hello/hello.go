package main

import (
	"example.com/greetings"
	"fmt"
	"log"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	//message, err := greetings.Hello("Alex")
	names := []string{"Alex", "Felix", "Rejoice"}
	messages, err := greetings.HelloV2(names)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(messages)
}
