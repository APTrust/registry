package forms

type Form struct {
	Action string
	Fields map[string]*Field
}

func NewForm() Form {
	return Form{
		Fields: make(map[string]*Field),
	}
}

type Field struct {
	Name         string
	Label        string
	Placeholder  string
	Value        interface{}
	ErrMsg       string
	DisplayError bool
	Options      []ListOption
	Attrs        map[string]string
}

type ListOption struct {
	Value string
	Text  string
}
