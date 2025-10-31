package internal

import "errors"

func Error(err ...string) error {
	msg := ""
	for _, e := range err {
		msg += e
	}
	return errors.New(msg)
}
