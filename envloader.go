package envloader

import (
	"strings"
)

type EnvLoader struct {
	processorStack fieldProcessorStack
	config         Config
}

// Customize special characters
type Config struct {
	// This string is used to separate attribute key
	// and values, by default it is colon (:)
	// Example:
	// 1. 	Definition = ":"
	// 		Tag = `env:"key:KEY_1"`
	// 2. 	Definition = "="
	// 		Tag = `env:"key=KEY_1"`
	Definition string
	// This string is used to separate
	// slices elements, by default it is comma (,)
	SliceSeparator string
	// This string is used to separate between
	// attributes, by default it is semicolon (;)
	AttributeSeparator string
}

// Create new instance of EnvLoader
//
// c - Configuration struct for configuring behaviour
func New(c *Config) EnvLoader {
	if c == nil {
		c = &Config{
			Definition:         ":",
			SliceSeparator:     ",",
			AttributeSeparator: ";",
		}
	}
	return EnvLoader{
		processorStack: defaultStack(*c),
		config:         *c,
	}
}

// Used to load environment data to struct
//
// **Must use pointer**
func (e *EnvLoader) Load(i interface{}) (errs []error) {
	return e.processorStack.Load(i)
}

// Can be used to create a default .ENV configuration file
//
// **Must use pointer**
func (e *EnvLoader) Stringify(i interface{}) (s string, err error) {
	sb := strings.Builder{}

	iterable, err := iterableType(i)
	if err != nil {
		return
	}

	for iterable.Next() {
		value, structField := iterable.Get()

		c, errConf := CreateConfig(value, structField, e.config)
		if errConf != nil {
			err = errConf
			return
		}

		c.WriteENVString(&sb)
	}

	s = sb.String()

	return
}
