package donkey

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type envReader struct {
	configStruct any
	replacer     *strings.Replacer
	prefix       string
}

func NewEnvReader(configStruct any) (*envReader, error) {
	if reflect.TypeOf(configStruct).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("destination must be a struct")
	}
	return &envReader{
		configStruct: configStruct,
	}, nil
}

func (r *envReader) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *envReader) SetReplacer(replacer *strings.Replacer) {
	r.replacer = replacer
}

func (r *envReader) Read() error {
	return r.readStructFields(r.prefix, r.configStruct)
}

func (r *envReader) readStructFields(prefix string, st any) error {
	rt := reflect.TypeOf(st)

	for i := 0; i < rt.Elem().NumField(); i++ {
		field := rt.Elem().Field(i)
		value, ok := field.Tag.Lookup(Tag)
		if !ok {
			continue
		}
		tagValues := strings.Split(value, ",")

		var name string
		required := false

		if len(tagValues) >= 1 {
			name = tagValues[0]
		} else {
			continue
		}

		if len(tagValues) >= 2 {
			required = tagValues[1] == "required"
		}

		// check field type
		envKey := strings.ToUpper(prefix + r.replacer.Replace(name))

		// if it is a struct go to find all of its fields
		if field.Type.Kind() == reflect.Struct {
			if err := r.readStructFields(envKey+"_", reflect.ValueOf(st).Elem().Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			return fmt.Errorf("unsupported type")
		}

		envValue := os.Getenv(envKey)
		if envValue == "" {
			if required {
				return fmt.Errorf("%s value is required", name)
			}
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			reflect.ValueOf(st).Elem().Field(i).SetString(envValue)
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			intValue, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid value for %s, %w", envKey, err)
			}
			reflect.ValueOf(st).Elem().Field(i).SetInt(intValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("invalid value for %s, %w", envKey, err)
			}
			reflect.ValueOf(st).Elem().Field(i).SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported type: %s", field.Type.String())
		}
	}
	return nil
}
