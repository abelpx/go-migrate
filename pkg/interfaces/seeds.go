package interfaces

type Seeds interface {
	error
	Seed(data ...map[string]interface{}) error
}
