package donkey

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type envReader struct {
	configStruct any
	replacer     *strings.Replacer
	prefix       string
}

// NewEnvReader create a new environment variable reader
func NewEnvReader(configStruct any) (*envReader, error) {
	if reflect.TypeOf(configStruct).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("destination must be a struct")
	}
	return &envReader{
		configStruct: configStruct,
	}, nil
}

// SetPrefix for each env
func (r *envReader) SetPrefix(prefix string) {
	r.prefix = prefix
}

// SetReplacer for fields name
func (r *envReader) SetReplacer(replacer *strings.Replacer) {
	r.replacer = replacer
}

// Read config and save it in config object
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

		if !(len(tagValues) >= 1 && tagValues[0] != "-") {
			continue
		}
		name = tagValues[0]

		// generate env key
		if r.replacer != nil {
			name = r.replacer.Replace(name)
		}

		envKey := strings.ToUpper(prefix + name)

		// if it is a struct go to find all of its fields
		fieldKind := field.Type.Kind()
		if fieldKind == reflect.Struct {
			if err := r.readStructFields(envKey+"_", reflect.ValueOf(st).Elem().Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		// pointer field not supported yet
		if fieldKind == reflect.Pointer {
			return fmt.Errorf("unsupported type")
		}

		envValue := os.Getenv(envKey)
		if envValue == "" {
			continue
		}

		if err := fieldFeeder(reflect.ValueOf(st).Elem().Field(i), envValue); err != nil {
			return err
		}
	}
	return nil
}
