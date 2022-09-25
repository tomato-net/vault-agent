package credentials

type Fake struct {
	Password string
}

func (f *Fake) Read() (PasswordCredential, error) {
	return PasswordCredential{Data: f.Password}, nil
}

func (f *Fake) Has() (bool, error) {
	return f.Password != "", nil
}
