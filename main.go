package main

import (
	"fmt"
	"log"

	"github.com/elos/mongo"
	"github.com/elos/space"
	"github.com/elos/stack"
	"github.com/robertkrimen/otto"

	"github.com/GeertJohan/go.linenoise"
)

var Otto = otto.New()

func main() {
	fmt.Println("Elos Script\n")
	mongo.DefaultLogger = mongo.NullLogger // hear nothing from mongo

	go mongo.Runner.Start()
	defer mongo.Runner.Stop()

	store := stack.SetupStore("localhost")

	/*
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("ID: ")
		scanner.Scan()
		id := scanner.Text()
		fmt.Println(id)
		fmt.Print("Key: ")
		scanner.Scan()
		key := scanner.Text()
		fmt.Println(key)
	*/

	id := "54dec8f084a588e46d000001"
	key := "cgy1LfpQ7UufLfh2eDBsC2IKLdei02o6_TGm4iC98xhBn_nBeXcqbLYbBrTrE1UH"

	c := &space.Credentials{
		ID:  id,
		Key: key,
	}

	s, err := space.NewSpace(c, store)
	if err != nil {
		log.Fatal(err)
	}

	s.Expose(Otto)

	loop()
}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	value, err := Otto.Run(entry)
	if err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("%v", value)
	}
}

func loop() {
	for {
		entry, err := linenoise.Line("> ")
		if err != nil {
			fmt.Println(err)
			break
		}
		if entry == "exit" {
			fmt.Println("Goodbye")
			break
		}

		output := dispatch(entry)
		if len(output) > 0 {
			fmt.Println(output)
		}
	}
}
