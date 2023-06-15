package donkey

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// fieldFeeder this function responsible to put a string value to requested field
func fieldFeeder(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("can not set value for field")
	}

	if reflect.DeepEqual(field.Kind(), reflect.ValueOf(time.Duration(0)).Kind()) {
		// field is a time duration
		// check if value all in number make convert it to second
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			// parse string value to integer -> 10s = 10000000000
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("failed to parse %s to time.Duration, err: %w", value, err)
			}
			intValue = duration.Nanoseconds()
		}

		field.SetInt(intValue)
		return nil
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
	case reflect.Float64, reflect.Float32:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %s to float, err: %w", value, err)
		}
		field.SetFloat(floatValue)
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
