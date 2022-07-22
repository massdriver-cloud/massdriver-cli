package application

import (
	"errors"
)

var (
	ErrCloneFail = errors.New("failed while cloning the application templates repository")
	ErrCopyFail  = errors.New("failed copying the application templates")
)
