package er

import (
	"errors"
	"fmt"
)

var ErrUserExists = errors.New("вы уже зарегистрированы")

// Wrap ...
func Wrap(msg string, err error) error {

	return fmt.Errorf("%s: %w", msg, err)
}
