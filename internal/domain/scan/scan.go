package scan

type Kit interface {
	Scan() (interface{}, error)
}

type Scanner struct {
}
