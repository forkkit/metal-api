package service

type PathParameter struct {
	name        string
	description string
}

func NewPathParameter(name, description string) *PathParameter {
	return &PathParameter{
		name:        name,
		description: description,
	}
}
