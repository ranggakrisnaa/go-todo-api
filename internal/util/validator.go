package util

import (
	"github.com/go-playground/validator/v10"
)

func IsRequestValid(u interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return false, err
	}
	return true, nil
}
