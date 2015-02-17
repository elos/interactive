package main

import (
	"fmt"
	"log"

	"github.com/elos/autonomous"
	"github.com/elos/mongo"
	"github.com/elos/space"
	"github.com/elos/stack"
	"github.com/robertkrimen/otto"

	"github.com/GeertJohan/go.linenoise"
)

var Otto = otto.New()

func main() {
	fmt.Println("Elos Script\n")
	mongo.Runner.Logger = mongo.NullLogger // hear nothing from mongo

	h := autonomous.NewHub()
	go h.Start()

	h.StartAgent(mongo.Runner)
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

	r := NewREPL()
	h.StartAgent(r)

	go autonomous.HandleIntercept(h.Stop)

	h.WaitStop()
}

type REPL struct {
	autonomous.Managed
	autonomous.Life
	autonomous.Stopper

	entries chan string
}

func NewREPL() *REPL {
	r := new(REPL)
	r.Life = autonomous.NewLife()
	r.Stopper = make(autonomous.Stopper)
	r.entries = make(chan string)
	return r
}

func (r *REPL) Start() {
	r.Life.Begin()

	go r.read()

Run:
	for {
		select {
		case entry := <-r.entries:
			r.dispatch(entry)
		case <-r.Stopper:
			break Run
		}
	}

	r.Life.End()
}

func (r *REPL) read() {
	entry, err := linenoise.Line("> ")
	if err != nil {
		fmt.Println(err)
		go r.Manager().Stop()
	}
	if entry == "exit" {
		fmt.Println("Goodbye")
		go r.Manager().Stop()
	}

	r.entries <- entry
}

func (repl *REPL) dispatch(entry string) {
	if len(entry) == 0 {
		return
	}
	value, err := Otto.Run(entry)

	var s string
	if err != nil {
		s = err.Error()
	} else {
		s = fmt.Sprintf("%v", value)
	}

	if len(s) > 0 {
		fmt.Println(s)
	}

	go repl.read()
}
