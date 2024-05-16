package reader

// Reader is an interface used to read a string stream
type Reader interface {
	Open(source string) error
	Read() (string, error)
	Close() error
}
