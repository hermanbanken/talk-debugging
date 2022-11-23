package job

import (
	"encoding/gob"
	"fmt"
	"strconv"
)

type Job struct {
	Op Op
	A  Expr
	B  Expr
}

func (j Job) String() string {
	var op string
	switch j.Op {
	case Add:
		op = "+"
	case Remove:
		op = "-"
	case Multiply:
		op = "*"
	case Divide:
		op = "/"
	}
	return fmt.Sprintf("%s %s %s", j.A, op, j.B)
}

type Op int

const (
	Add      Op = 0
	Remove   Op = 1
	Multiply Op = 2
	Divide   Op = 3
)

type Expr interface {
	expr()
	String() string
}

type Value int

func (v Value) String() string { return strconv.Itoa(int(v)) }

func (Value) expr() {}

func (Job) expr() {}

func init() {
	gob.Register(Job{})
	gob.Register(Value(0))
}
