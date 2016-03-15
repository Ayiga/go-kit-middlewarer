package main

type Variable struct {
	name  string
	typ   *Type // nillable
	value string
}

func createVariable(name, value string, typ *Type) Variable {
	return Variable{
		name:  name,
		typ:   typ,
		value: value,
	}
}
