package credentials

type PasswordCredential struct {
	Data string
}

type Reader interface {
	Read() (PasswordCredential, error)
	Has() (bool, error)
}

type Writer interface {
	Write(PasswordCredential) error
	Delete() error
}

type ReadWriter interface {
	Reader
	Writer
}
