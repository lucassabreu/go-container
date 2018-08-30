package test

import "fmt"

// Doer do
type Doer interface {
	Do()
}

// Config needed for DoByConfig
type Config struct {
	What string
	When string
}

// DoByConfig uses a config
type DoByConfig struct {
	config Config
}

// NewDoByConfig creates a DoByConfig
func NewDoByConfig(c Config) DoByConfig {
	return DoByConfig{config: c}
}

// Do things
func (d DoByConfig) Do() {
	fmt.Printf("%s at %s", d.config.What, d.config.When)
}

type iDo struct{}

func (id iDo) Do() {
	println("i did")
}

// TheyDo things
type TheyDo struct {
	what string
	ToDo func(string)
}

// Do to do
func (td TheyDo) Do() {
	td.ToDo(td.what)
}

// NewIDo return new iDo
func NewIDo() Doer {
	return iDo{}
}

// NewTheyDo return new TheyDo
func NewTheyDo(toDo func(string)) TheyDo {
	return TheyDo{
		what: "something",
		ToDo: toDo,
	}
}

// JustDo is a example
type JustDo struct {
	That string
}

// Do does
func (d JustDo) Do() {
	println(d.That)
}

// NewJustDo receives a scalar
func NewJustDo(that string) JustDo {
	return JustDo{That: that}
}

// SomethingDo calls other Doer
type SomethingDo struct {
	Something Doer
}

// NewSomethingDo dependence
func NewSomethingDo(do Doer) SomethingDo {
	return SomethingDo{Something: do}
}

// Do does
func (s SomethingDo) Do() {
	s.Something.Do()
}

type doALot struct {
	things []Doer
}

// NewDoALot returns a interface pointer
func NewDoALot(doers []Doer) *Doer {
	var doer Doer
	doer = doALot{things: doers}
	return &doer
}

// NewDoALotByPointer receives a pointer slice
func NewDoALotByPointer(doers *[]Doer) Doer {
	doer := NewDoALot(*doers)
	return *doer
}

func (a doALot) Do() {
	for _, doer := range a.things {
		doer.Do()
	}
}

// ToDo has a map
type ToDo struct {
	toDo map[string]Doer
}

// NewToDo recieves a map
func NewToDo(toDos map[string]Doer) *ToDo {
	return &ToDo{
		toDo: toDos,
	}
}

// Do does
func (t *ToDo) Do() {
	for name, doer := range t.toDo {
		fmt.Printf("%s:\n", name)
		doer.Do()
	}
}
