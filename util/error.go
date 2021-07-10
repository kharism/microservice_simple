package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	//ErrUtilJWTMapping error message for JWT mapping
	ErrUtilJWTMapping = "invalid mapping JWT claim with key '%s'"
)

// Error customized error
type Error struct {
	Tag     string
	Message interface{}
	Errors  []error
	Data    map[string]interface{}
}

func (b Error) Error() string {
	var message string
	if IsZeroOfUnderlyingType(b.Message) {
		message = "error occured"
	} else if v, ok := b.Message.(string); ok {
		message = v
	} else if v, ok := b.Message.(error); ok {
		message = v.Error()
	} else {
		message = "error occured"
	}

	if len(b.Errors) > 0 {
		return MergeError(message, b.Errors...).Error()
	}

	return message
}

// MergeError merge error message
func MergeError(message string, err ...error) error {
	errMessages := []string{}
	if len(err) > 0 {
		for _, e := range err {
			errMessages = append(errMessages, e.Error())
		}

		return fmt.Errorf("%s\n\t-> %s", message, strings.Join(errMessages, "\n\t-> "))
	}

	return errors.New(message)
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
