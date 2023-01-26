package envloader

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	ARRAY_SEPARATOR = ","
)

type fieldProcessor struct {
	Kind reflect.Kind
	// A function that gets the tag for a field and return a value
	// to assign it to that field
	Processor func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (reflect.Value, error)
}

type fieldProcessorStack []fieldProcessor

func (p fieldProcessorStack) Load(i interface{}) (errs []error) {
	iterable, err := iterableType(i)
	if err != nil {
		errs = append(errs, err)
		return
	}

	for iterable.Next() {
		value, structField := iterable.Get()

		conf, err := CreateConfig(value, structField)
		if err != nil {
			errs = append(errs, fmt.Errorf("ENVLOADER: Invalid tag, for field %s error: %s", structField.Name, err.Error()))
			continue
		}

		if conf.Key == "" && value.Kind() != reflect.Struct {
			continue
		}

		fieldValue, err := p.ProcessField(value, structField, conf)
		if err != nil {
			errs = append(errs, fmt.Errorf("ENVLOADER: Unable to set value for field %s, error: %s", structField.Name, err.Error()))
			continue
		}
		value.Set(fieldValue)
	}

	return
}

func (p fieldProcessorStack) ProcessField(val reflect.Value, field reflect.StructField, conf FieldConfig) (v reflect.Value, err error) {
	for _, processor := range p {
		if val.Kind() == processor.Kind {
			v, err = processor.Processor(val, field, conf, p)
		}
	}
	return
}

func defaultStack() fieldProcessorStack {
	return fieldProcessorStack{
		intProcessor,
		int8Processor,
		int16Processor,
		int32Processor,
		int64Processor,
		float32Processor,
		float64Processor,
		stringProcessor,
		sliceProcessor,
		structProcessor,
	}
}

var intProcessor = fieldProcessor{
	Kind: reflect.Int,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseInt := func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		}

		value, from := conf.GetValue()

		intVal, err := parseInt(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse integer with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(int(intVal))
		return
	},
}

var int8Processor = fieldProcessor{
	Kind: reflect.Int8,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseInt := func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 8)
		}

		value, from := conf.GetValue()

		intVal, err := parseInt(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse integer with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(int8(intVal))
		return
	},
}

var int16Processor = fieldProcessor{
	Kind: reflect.Int16,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseInt := func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 16)
		}

		value, from := conf.GetValue()

		intVal, err := parseInt(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse integer with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(int16(intVal))
		return
	},
}

var int32Processor = fieldProcessor{
	Kind: reflect.Int32,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseInt := func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 32)
		}

		value, from := conf.GetValue()

		intVal, err := parseInt(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse integer with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(int32(intVal))
		return
	},
}

var int64Processor = fieldProcessor{
	Kind: reflect.Int64,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseInt := func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		}

		value, from := conf.GetValue()

		intVal, err := parseInt(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse integer with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(int64(intVal))
		return
	},
}

var float32Processor = fieldProcessor{
	Kind: reflect.Float32,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseFloat := func(s string) (float64, error) {
			return strconv.ParseFloat(s, 32)
		}

		value, from := conf.GetValue()

		floatVar, err := parseFloat(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse float with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(float32(floatVar))
		return
	},
}

var float64Processor = fieldProcessor{
	Kind: reflect.Float64,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		parseFloat := func(s string) (float64, error) {
			return strconv.ParseFloat(s, 64)
		}

		value, from := conf.GetValue()

		floatVar, err := parseFloat(value)
		if err != nil {
			err = fmt.Errorf("Unable to parse float with value : %s | source : %s", value, from)
		}

		v = reflect.ValueOf(float64(floatVar))
		return
	},
}

var stringProcessor = fieldProcessor{
	Kind: reflect.String,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		value, _ := conf.GetValue()

		v = reflect.ValueOf(value)
		return
	},
}

var sliceProcessor = fieldProcessor{
	Kind: reflect.Slice,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		rawValues, from := conf.GetValue()

		rawValueArr := strings.Split(rawValues, ARRAY_SEPARATOR)

		elemType := field.Type.Elem()
		elemVal := reflect.Zero(elemType)

		v = reflect.MakeSlice(field.Type, 0, len(rawValueArr))

		for _, rawValue := range rawValueArr {
			var vEl reflect.Value

			conf.OverrideValue = &rawValue
			vEl, err = stack.ProcessField(elemVal, field, conf)
			if err != nil {
				err = fmt.Errorf("Unable to parse slice of type %s with value : %s | source : %s", elemType.Kind().String(), rawValues, from)
			}

			v = reflect.Append(v, vEl)
		}

		return
	},
}

var structProcessor = fieldProcessor{
	Kind: reflect.Struct,
	Processor: func(val reflect.Value, field reflect.StructField, conf FieldConfig, stack fieldProcessorStack) (v reflect.Value, err error) {
		vAddr := val.Addr()
		stack.Load(vAddr.Interface())
		v = vAddr.Elem()

		return
	},
}
