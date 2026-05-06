package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			return err
		}

		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, formatValidationError(e))
		}

		return fmt.Errorf("%s", strings.Join(msgs, "; "))
	}
	return nil
}

func formatValidationError(e validator.FieldError) string {
	field := toSnakeCase(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "email":
		return fmt.Sprintf("%s harus berupa email yang valid", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, e.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, e.Param())
	case "len":
		return fmt.Sprintf("%s harus tepat %s karakter", field, e.Param())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", field)
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh berisi huruf dan angka", field)
	case "e164":
		return fmt.Sprintf("%s harus berformat nomor telepon internasional (contoh: +628xxxxxxxxx)", field)
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih dari %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s harus lebih dari atau sama dengan %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s harus kurang dari %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s harus kurang dari atau sama dengan %s", field, e.Param())
	case "url":
		return fmt.Sprintf("%s harus berupa URL yang valid", field)
	case "eqfield":
		return fmt.Sprintf("%s harus sama dengan %s", field, toSnakeCase(e.Param()))
	case "nefield":
		return fmt.Sprintf("%s tidak boleh sama dengan %s", field, toSnakeCase(e.Param()))
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}

// toSnakeCase converts PascalCase/camelCase field names to snake_case for friendlier error messages.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}