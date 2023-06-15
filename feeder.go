package donkey

import (
	"fmt"
	"reflect"
	"strconv"
)

// fieldFeeder this function responsible to put a string value to requested field
func fieldFeeder(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("can not set value for field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s, %w", value, err)
		}
		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s, %w", value, err)
		}
		field.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported type: %s", field.Type().String())
	}

	return nil
}
