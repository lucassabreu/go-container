package example

// Doer do
type Doer interface {
	Do()
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

type JustDo struct {
	That string
}

func (d JustDo) Do() {
	println(d.That)
}

func NewJustDo(that string) JustDo {
	return JustDo{That: that}
}

type SomethingDo struct {
	Something Doer
}

func NewSomethingDo(do Doer) SomethingDo {
	return SomethingDo{Something: do}
}

func (s SomethingDo) Do() {
	s.Something.Do()
}

type doALot struct {
	things []Doer
}

func NewDoALot(doers []Doer) doALot {
	return doALot{things: doers}
}

func (a doALot) Do() {
	for _, doer := range a.things {
		doer.Do()
	}
}
