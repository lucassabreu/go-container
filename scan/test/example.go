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
