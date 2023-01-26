package envloader

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	TAG_KEY = "env"

	SOURCE_ENV      = "ENV"
	SOURCE_DEFAULT  = "DEFAULT VALUE"
	SOURCE_OVERRIDE = "OVERRIDE"
)

type FieldConfig struct {
	Key           string
	DefaultValue  *string
	OverrideValue *string
	Config        Config
	value         reflect.Value
	field         reflect.StructField
}

func (f *FieldConfig) GetValue() (v, source string) {
	if f.OverrideValue != nil {
		return *f.OverrideValue, SOURCE_OVERRIDE
	}

	v = os.Getenv(f.Key)
	source = SOURCE_ENV
	if v == "" && f.DefaultValue != nil {
		v = *f.DefaultValue
		source = SOURCE_DEFAULT
	}

	return
}

func (f *FieldConfig) WriteENVString(sb *strings.Builder) {
	if f.value.Kind() == reflect.Struct {
		sb.WriteRune('\n')
		iter, err := iterableType(f.value.Addr().Interface())
		if err == nil {
			for iter.Next() {
				value, structField := iter.Get()

				conf, err := CreateConfig(value, structField, f.Config)
				if err != nil {
					continue
				}

				conf.WriteENVString(sb)
			}
		}
	}

	if f.Key == "" {
		return
	}

	val := ""
	if f.DefaultValue != nil {
		val = *f.DefaultValue
	}
	sb.WriteString(fmt.Sprintf("%s = %s\n", f.Key, val))
}

func CreateConfig(v reflect.Value, s reflect.StructField, conf Config) (c FieldConfig, err error) {
	c.value = v
	c.field = s
	c.Config = conf

	if s.Type.Kind() == reflect.Struct {
		return
	}

	tags := s.Tag.Get(TAG_KEY)

	if tags == "" {
		return
	}

	for _, tag := range strings.Split(tags, conf.AttributeSeparator) {
		splitTag := strings.SplitN(tag, conf.Definition, 2)
		if len(splitTag) <= 1 {
			err = fmt.Errorf("Invalid tag %s", tag)
			continue
		}

		k, v := splitTag[0], splitTag[1]

		switch k {
		case "key":
			c.Key = v
		case "default":
			c.DefaultValue = &v
		}
	}

	if c.Key == "" {
		err = fmt.Errorf("Field %s does not have key", c.field.Name)
	}

	return
}
