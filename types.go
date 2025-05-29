package tabby

import (
	"errors"
	"strconv"
	"time"
)

var (
	Int = NewTransfer("int", func(s string) (any, error) {
		if i, err := strconv.ParseInt(s, 10, 0); err != nil {
			return 0, err
		} else {
			return int(i), nil
		}
	})

	String = NewTransfer("string", func(s string) (any, error) {
		return s, nil
	})

	Bool = NewTransfer("bool", func(s string) (any, error) {
		if len(s) == 0 {
			return true, nil
		} else {
			return false, errors.New("boolean requires no value")
		}
	})

	Duration = NewTransfer("duration", func(s string) (any, error) {
		if i, err := time.ParseDuration(s); err != nil {
			return time.Duration(0), err
		} else {
			return time.Duration(i), nil
		}
	})
)
