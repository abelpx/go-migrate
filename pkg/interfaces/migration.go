package interfaces

type Migration interface {
	Up() error
	Down() error
}
