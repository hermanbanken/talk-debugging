package job

import "encoding/gob"

type Job struct {
	Op Op
	A  Expr
	B  Expr
}

type Op int

const (
	Add      Op = 0
	Remove   Op = 1
	Multiply Op = 2
	Divide   Op = 3
)

type Expr interface{ expr() }

type Value int

func (Value) expr() {}

func (Job) expr() {}

func init() {
	gob.Register(Job{})
	gob.Register(Value(0))
}
