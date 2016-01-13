package interactive

import (
	"fmt"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/robertkrimen/otto"
)

const (
	ottoObjectString = "[object Object]"
	jsonParseFormat  = "JSON.parse(%s)"
)

type Credentials struct {
	ID  string
	Key string
}

type Env struct {
	db   data.DB
	otto *otto.Otto
	user *models.User
}

func NewEnv(db data.DB, u *models.User) *Env {
	e := new(Env)

	e.db = db
	e.otto = otto.New()
	e.otto.Set("me", u)
	e.otto.Set("db", db)
	e.user = u

	return e
}

func (e *Env) Interpret(entry string) string {
	if len(entry) == 0 {
		return entry
	}

	value, err := e.otto.Run(entry)

	if err != nil {
		return err.Error()
	}

	/* results in SyntaxError: invalid character 'o' looking for beginning of value
	if value.IsObject() {
		return e.Interpret(fmt.Sprintf(jsonParseFormat, entry))
	}
	*/

	return fmt.Sprintf("%v", value)
}

func (e *Env) Set(variableName string, value interface{}) {
	e.otto.Set(variableName, value)
}
