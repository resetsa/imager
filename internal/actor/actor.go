package actor

type Actor interface {
	ActMany(paths []string) error
	ActOnce(path string) error
}
