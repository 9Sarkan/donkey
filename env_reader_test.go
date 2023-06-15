package donkey

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSimpleRead(t *testing.T) {
	type testConfig struct {
		Hello string `mapstructure:"hello,required"`
		Age   int    `mapstructure:"age,required"`
		Flag  bool   `mapstructure:"flag"`
	}

	os.Setenv("HELLO", "world")
	os.Setenv("AGE", "10")
	os.Setenv("FLAG", "1")

	wantConfig := &testConfig{
		Hello: "world",
		Age:   10,
		Flag:  true,
	}
	gotConfig := &testConfig{}
	configReader, err := NewEnvReader(gotConfig)
	if err != nil {
		t.Error(err)
	}
	configReader.SetReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if err := configReader.Read(); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("want: %+v got: %+v", wantConfig, gotConfig)
	}
}

func TestNestedRead(t *testing.T) {
	type nestedConfig struct {
		Hi string `mapstructure:"hi"`
	}
	type testConfig struct {
		Hello        string       `mapstructure:"hello,required"`
		Age          int          `mapstructure:"age,required"`
		Flag         bool         `mapstructure:"flag"`
		NestedConfig nestedConfig `mapstructure:"nested"`
	}

	os.Setenv("HELLO", "world")
	os.Setenv("AGE", "10")
	os.Setenv("FLAG", "1")
	os.Setenv("NESTED_HI", "hi there")

	wantConfig := &testConfig{
		Hello: "world",
		Age:   10,
		Flag:  true,
		NestedConfig: nestedConfig{
			Hi: "hi there",
		},
	}
	gotConfig := &testConfig{}
	configReader, err := NewEnvReader(gotConfig)
	if err != nil {
		t.Error(err)
	}
	configReader.SetReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if err := configReader.Read(); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("want: %+v got: %+v", wantConfig, gotConfig)
	}
}

func TestNestedPointerRead(t *testing.T) {
	type nestedConfig struct {
		Hi string `mapstructure:"hi"`
	}
	type testConfig struct {
		Hello        string        `mapstructure:"hello,required"`
		Age          int           `mapstructure:"age,required"`
		Flag         bool          `mapstructure:"flag"`
		NestedConfig *nestedConfig `mapstructure:"nested"`
	}

	os.Setenv("HELLO", "world")
	os.Setenv("AGE", "10")
	os.Setenv("FLAG", "1")
	os.Setenv("NESTED_HI", "hi there")

	wantConfig := &testConfig{
		Hello:        "world",
		Age:          10,
		Flag:         true,
		NestedConfig: nil,
	}
	gotConfig := &testConfig{}
	configReader, err := NewEnvReader(gotConfig)
	if err != nil {
		t.Error(err)
	}
	configReader.SetReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if err := configReader.Read(); err != nil {
		if err.Error() != "unsupported type" {
			t.Error(err)
		}
		return
	}

	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Errorf("want: %+v got: %+v", wantConfig, gotConfig)
	}
}
