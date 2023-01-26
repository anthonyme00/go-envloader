package envloader

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

func generateRandomValue() reflect.Value {
	store := []interface{}{
		int64(rand.Int() % 256),
		int32(rand.Int() % 256),
		int16(rand.Int() % 256),
		int8(rand.Int() % 256),
		int(rand.Int() % 256),
		float64(rand.Float64() * 100.0),
		float32(rand.Float32() * 100.0),
		string(generateRandomString(10)),
	}

	i := rand.Int() % len(store)

	return reflect.ValueOf(store[i])
}

func generateRandomValueArr(size int) (reflect.Value, string) {
	typeCount := 8

	randInt64Arr := func() (data []int64) {
		for i := 0; i < size; i++ {
			data = append(data, int64(rand.Int()%256))
		}
		return
	}

	randInt32Arr := func() (data []int32) {
		for i := 0; i < size; i++ {
			data = append(data, int32(rand.Int()%256))
		}
		return
	}

	randInt16Arr := func() (data []int16) {
		for i := 0; i < size; i++ {
			data = append(data, int16(rand.Int()%256))
		}
		return
	}

	randInt8Arr := func() (data []int8) {
		for i := 0; i < size; i++ {
			data = append(data, int8(rand.Int()%256))
		}
		return
	}

	randIntArr := func() (data []int) {
		for i := 0; i < size; i++ {
			data = append(data, rand.Int()%256)
		}
		return
	}

	randFloat64Arr := func() (data []float64) {
		for i := 0; i < size; i++ {
			data = append(data, rand.Float64()*100.0)
		}
		return
	}

	randFloat32Arr := func() (data []float32) {
		for i := 0; i < size; i++ {
			data = append(data, rand.Float32()*100.0)
		}
		return
	}

	randStringArr := func() (data []string) {
		for i := 0; i < size; i++ {
			data = append(data, string(generateRandomString(10)))
		}
		return
	}

	i := rand.Int() % typeCount

	val := reflect.Value{}
	returnStrElems := []string{}

	switch i {
	case 0:
		data := randInt64Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 1:
		data := randInt32Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 2:
		data := randInt16Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 3:
		data := randInt8Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 4:
		data := randIntArr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 5:
		data := randFloat64Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 6:
		data := randFloat32Arr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	case 7:
		data := randStringArr()
		for _, d := range data {
			returnStrElems = append(returnStrElems, fmt.Sprint(d))
		}
		val = reflect.ValueOf(data)
	}

	return val, strings.Join(returnStrElems, ",")
}

type dummyStruct struct {
	Value        reflect.Value
	Struct       reflect.Type
	Expected     map[string]reflect.Value
	RawStr       map[string]string
	ChildStructs map[string]dummyStruct
}

func (t *dummyStruct) rebuildReference() {
	iter, _ := iterableType(t.Value.Interface())

	for iter.Next() {
		val, structEl := iter.Get()
		if child, ok := t.ChildStructs[structEl.Name]; ok {
			newChild := child
			newChild.Value = val.Addr()
			newChild.Struct = structEl.Type
			t.ChildStructs[structEl.Name] = newChild
		}
	}
}

func (t *dummyStruct) compare() bool {
	t.rebuildReference()

	if len(t.ChildStructs) > 0 {
		for _, child := range t.ChildStructs {
			if !child.compare() {
				return false
			}
		}
	}

	iter, err := iterableType(t.Value.Interface())
	if err != nil {
		return false
	}

	for iter.Next() {
		value, structField := iter.Get()

		if _, ok := t.Expected[structField.Name]; !ok {
			continue
		}

		expectedVal := t.Expected[structField.Name]

		if fmt.Sprint(value.Interface()) != fmt.Sprint(expectedVal.Interface()) {
			return false
		}
	}

	return true
}

func (t *dummyStruct) GetEnvStrings() (data []string) {
	for _, child := range t.ChildStructs {
		data = append(data, child.GetEnvStrings()...)
	}

	for k, v := range t.RawStr {
		data = append(data, fmt.Sprintf("%s = %s", k, v))
	}

	return
}

func generateRandomString(length int) string {
	strBank := "ABCDEFGHIJKLMNOPQRSTUVWXYZ 1234567890"

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = strBank[rand.Int()%len(strBank)]
	}

	return string(result)
}

func generateRandomStruct(alwaysUseDefault bool, size int, depth int, maxChild int, idx *int) (t dummyStruct) {
	t.Expected = make(map[string]reflect.Value)
	t.ChildStructs = make(map[string]dummyStruct)
	t.RawStr = make(map[string]string)

	keyGen := func() string {
		return fmt.Sprintf("KEY_%d", *idx)
	}
	keyChildGen := func() string {
		return fmt.Sprintf("STRUCT_%d", *idx)
	}
	keyInvalidGen := func() string {
		return fmt.Sprintf("INVALID_%d", *idx)
	}

	getTag := func(key, expected interface{}) string {
		if rand.Int()%2 == 0 || alwaysUseDefault {
			return fmt.Sprintf(`env:"key:%s;default:%v"`, key, expected)
		}

		os.Setenv(fmt.Sprint(key), fmt.Sprint(expected))

		return fmt.Sprintf(`env:"key:%s"`, key)
	}

	fields := make([]reflect.StructField, 0)
	for i := 0; i < size; i++ {
		*idx = *idx + 1
		if i%2 == 0 {
			key, expected := keyGen(), generateRandomValue()

			tag := getTag(key, expected)
			field := reflect.StructField{
				Name: key,
				Type: expected.Type(),
				Tag:  reflect.StructTag(tag),
			}

			fields = append(fields, field)

			t.Expected[key] = expected
			t.RawStr[key] = fmt.Sprint(expected)
		} else {
			key := keyGen()
			expected, str := generateRandomValueArr(5)
			tag := getTag(key, str)

			field := reflect.StructField{
				Name: key,
				Type: expected.Type(),
				Tag:  reflect.StructTag(tag),
			}

			fields = append(fields, field)

			t.Expected[key] = expected
			t.RawStr[key] = str
		}
	}

	if depth > 0 && maxChild > 0 {
		for i := 0; i < maxChild; i++ {
			*idx = *idx + 1
			child := generateRandomStruct(alwaysUseDefault, size, depth-1, maxChild, idx)
			fields = append(fields, reflect.StructField{
				Name: keyChildGen(),
				Type: child.Struct,
			})

			t.ChildStructs[keyChildGen()] = child
		}
	}

	randInvalidCount := rand.Int() % 5
	for i := 0; i < randInvalidCount; i++ {
		*idx = *idx + 1
		fields = append(fields, reflect.StructField{
			Name: keyInvalidGen(),
			Type: reflect.TypeOf(int(0)),
		})
	}

	t.Struct = reflect.StructOf(fields)

	t.Value = reflect.New(t.Struct)

	return
}

func TestDefaultValues(t *testing.T) {
	for i := 0; i < 20; i++ {
		idx := 0
		testData := generateRandomStruct(false, 10, 1, 2, &idx)

		processor := New(nil)

		errs := processor.Load(testData.Value.Interface())

		if len(errs) != 0 {
			fmt.Println(errs)
			t.FailNow()
		}

		if !testData.compare() {
			t.FailNow()
		}

		os.Clearenv()
	}
}

func TestGeneratedEnv(t *testing.T) {
	matchSlice := func(a []string, b []string, excludeEmpty bool) bool {
		approvedA := 0
		approvedB := 0
		mapA := map[string]struct{}{}
		for _, kA := range a {
			if excludeEmpty && kA == "" {
				continue
			}
			mapA[kA] = struct{}{}
			approvedA++
		}

		for _, kB := range b {
			if excludeEmpty && kB == "" {
				continue
			}
			if _, ok := mapA[kB]; !ok {
				return false
			}
			approvedB++
		}

		if approvedA != approvedB {
			return false
		}

		return true
	}

	for i := 0; i < 20; i++ {
		idx := 0
		testData := generateRandomStruct(true, 10, 1, 2, &idx)

		processor := New(nil)

		str, err := processor.Stringify(testData.Value.Interface())
		strArr := strings.Split(str, "\n")

		expected := testData.GetEnvStrings()

		if err != nil {
			t.FailNow()
		}

		if !matchSlice(expected, strArr, true) {
			t.FailNow()
		}

		os.Clearenv()
	}
}

func TestErrorNonPointer(t *testing.T) {
	processor := New(nil)

	type TestStruct struct {
		Test  string `env:"key:TEST"`
		Test1 struct {
			Test1 string `env:""`
		}
	}

	test := TestStruct{}

	errs := processor.Load(test)

	if len(errs) == 0 {
		t.FailNow()
	}
}

func TestErrorNonPointerStringer(t *testing.T) {
	processor := New(nil)

	type TestStruct struct {
		Test string `env:"key:TEST"`
	}

	test := TestStruct{}

	_, err := processor.Stringify(test)

	if err == nil {
		t.FailNow()
	}
}

func TestErrorInvalidTag(t *testing.T) {
	processor := New(nil)

	type TestStruct struct {
		Test string `env:"key:TEST:234"`
	}

	test := TestStruct{}

	_, err := processor.Stringify(&test)

	if err == nil {
		t.FailNow()
	}
}

func TestErrorEmptyTag(t *testing.T) {
	processor := New(nil)

	type TestStruct struct {
		Test string
	}

	test := TestStruct{}

	_, err := processor.Stringify(&test)

	if err != nil {
		t.FailNow()
	}
}
