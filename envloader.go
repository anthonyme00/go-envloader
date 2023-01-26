package envloader

import (
	"strings"
)

type EnvLoader struct {
	processorStack fieldProcessorStack
	config         Config
}

type Config struct {
	Definition         string
	SliceSeparator     string
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

func (e *EnvLoader) Load(i interface{}) (errs []error) {
	return e.processorStack.Load(i)
}

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
