package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/elos/autonomous"
	"github.com/elos/space"
	"github.com/elos/stack"
	"github.com/robertkrimen/otto"

	"github.com/GeertJohan/go.linenoise"
)

var Otto = otto.New()

func main() {
	fmt.Fprintln(os.Stdout, "Elos Script\n")
	// mongo.Runner.Logger = mongo.NullLogger // hear nothing from mongo

	h := autonomous.NewHub()
	go h.Start()

	// h.StartAgent(mongo.Runner)
	store := stack.SetupStore("localhost")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, "ID: ")
	scanner.Scan()
	id := scanner.Text()
	fmt.Fprint(os.Stdout, "Key: ")
	scanner.Scan()
	key := scanner.Text()

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
		return
	}
	if entry == "exit" {
		fmt.Println("Goodbye")
		go r.Manager().Stop()
		return
	}

	r.entries <- entry
	linenoise.AddHistory(entry)
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
