package er

import "fmt"

// Wrap ...
func Wrap(msg string, err error) error {

	return fmt.Errorf("%s: %w", msg, err)
}

// Wrap ...
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
