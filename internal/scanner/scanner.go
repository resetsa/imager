package scanner

type Scanner interface {
	Scan(out chan string) error
	RootDir() string
}
