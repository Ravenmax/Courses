package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(e.errs)) + " errors occured:\n")

	for _, err := range e.errs {
		sb.WriteString("\t* ")
		sb.WriteString(err.Error())
	}
	sb.WriteString("\n")
	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	var multi *MultiError
	var ok bool
	if multi, ok = err.(*MultiError); !ok {
		multi = &MultiError{}
	}

	// Добавляем остальные errors
	for _, e := range errs {
		if e != nil {
			multi.errs = append(multi.errs, e)
		}
	}

	// Если нет ни одной ошибки, возвращаем nil
	if len(multi.errs) == 0 {
		return nil
	}

	return multi
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
