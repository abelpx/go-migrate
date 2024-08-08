package interfaces

type ForeignFoundation interface {
	Reference(name string) ForeignFoundation
	On(table string) ForeignFoundation
	OnUpdate(action string) ForeignFoundation
	OnDelete(action string) ForeignFoundation
}
