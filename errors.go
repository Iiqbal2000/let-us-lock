package main

import "errors"

var (
	ErrCmd          = errors.New("you have to include 'encrypt' or 'decrypt' command")
	ErrPassWrong    = errors.New("the password/passphrase is invalid")
	ErrPassTooShort = errors.New("the password/passphrase too short (MIN 8 characters)")
	ErrPassTooLong  = errors.New("the password/passphrase too long (MAX 64 characters)")
	ErrPassNotFound = errors.New("the password/passphrase is required")
	ErrFileNotFound = errors.New("the file is not found")
	ErrSaltNotFound = errors.New("failure when reading a salt file")
)
