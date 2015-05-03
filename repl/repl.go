package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/elos/autonomous"
	"github.com/elos/interactive"
	"github.com/elos/models"

	"github.com/GeertJohan/go.linenoise"
)

func main() {
	fmt.Fprintln(os.Stdout, "Welcome to the Elos REPL")
	// mongo.Runner.Logger = mongo.NullLogger // hear nothing from mongo

	h := autonomous.NewHub()
	go h.Start()

	store, err := models.MongoDB("localhost")
	if err != nil {
		log.Fatal("Failed to connect to db: %s", err)
	}

	//_ := RetrieveCredentials(bufio.NewScanner(os.Stdin))

	e := interactive.NewEnv(store, models.NewUser())
	r := NewREPL(e)
	h.StartAgent(r)
	go autonomous.HandleIntercept(h.Stop)
	h.WaitStop()
}

func RetrieveCredentials(scanner *bufio.Scanner) *interactive.Credentials {
	fmt.Fprint(os.Stdout, "ID: ")
	scanner.Scan()
	id := scanner.Text()
	fmt.Fprint(os.Stdout, "Key: ")
	scanner.Scan()
	key := scanner.Text()

	c := &interactive.Credentials{
		ID:  id,
		Key: key,
	}

	return c
}

type REPL struct {
	autonomous.Managed
	autonomous.Life
	autonomous.Stopper

	env *interactive.Env
}

func NewREPL(e *interactive.Env) *REPL {
	r := new(REPL)
	r.Life = autonomous.NewLife()
	r.Stopper = make(autonomous.Stopper)
	r.env = e
	return r
}

func (r *REPL) Start() {
	r.Life.Begin()
	go r.read()
	<-r.Stopper
	r.Life.End()
}

func (r *REPL) exit() {
	fmt.Println("Goodbye")
	go r.Manager().Stop()
}

func (r *REPL) read() {
	entry, err := linenoise.Line("> ")
	if err != nil {
		log.Fatal(err)
	}

	if entry == "exit" {
		r.exit()
		return
	}

	go r.dispatch(entry)
	linenoise.AddHistory(entry)
}

func (r *REPL) dispatch(entry string) {
	fmt.Println(r.env.Interpret(entry))
	go r.read()
}
