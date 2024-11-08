package checker

type Checker interface {
	// return true if need process
	Check(string) (bool, error)
}
