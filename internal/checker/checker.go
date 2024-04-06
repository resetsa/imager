package checker

type Checker interface {
	// return true if need process
	CheckImage(string) (bool, error)
}
