package main

const (
	minPassphraseLength = 8
	maxPassphraseLength = 64
)

type passphrase []byte

func (pp passphrase) validate() ([]byte, error) {
	if len(pp) < minPassphraseLength {
		return make([]byte, 0), ErrPassTooShort
	} else if len(pp) > maxPassphraseLength {
		return make([]byte, 0), ErrPassTooLong
	}

	return pp, nil
}

func catchPassphrase(passphraseIn []byte, err error) ([]byte, error) {
	if err != nil {
		return make([]byte, 0), ErrPassNotFound
	}

	return passphrase(passphraseIn).validate()
}

