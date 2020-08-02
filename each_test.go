package every

import (
	"fmt"
	"reflect"
	"testing"
)

type (
	TestString string
	TestStruct struct {
		Name string
	}
	TestNestedStruct struct {
		Name TestString
		In   *TestNestedStruct
		More []TestNestedStruct
	}
)

var (
	testString       = TestString("")
	testStruct       = TestStruct{}
	testNestedStruct = TestNestedStruct{}

	testStringType       = reflect.TypeOf(testString)
	testStructType       = reflect.TypeOf(testStruct)
	testNestedStructType = reflect.TypeOf(testNestedStruct)

	ptrTestStringType       = reflect.TypeOf(&testString)
	ptrTestStructType       = reflect.TypeOf(&testStruct)
	ptrTestNestedStructType = reflect.TypeOf(&testNestedStruct)
)

func TestRenamedString(t *testing.T) {
	strings := make([]interface{}, 100)
	for i := 0; i < len(strings); i++ {
		if i%2 == 0 {
			strings[i] = TestString(fmt.Sprintf("test%02d", i))
		} else {
			strings[i] = fmt.Sprintf("test%02d", i)
		}
	}

	var testStrings []TestString
	For(testStringType).In(strings).Do(func(typeObj interface{}) {
		testStrings = append(testStrings, typeObj.(TestString))
	})

	for i, v := range testStrings {
		if string(v) != fmt.Sprintf("test%02d", i*2) {
			t.Errorf("got unexpected value: %s", v)
		}
	}
}

func TestAssignStructs(t *testing.T) {
	sliceTestStructs := make([]*TestStruct, 10)
	// assigning values to non-ptr structs or nil values will not work ofc.
	for i := range sliceTestStructs {
		sliceTestStructs[i] = &TestStruct{}
	}

	For(ptrTestStructType).In(sliceTestStructs).Do(func(typeObj interface{}) {
		obj := typeObj.(*TestStruct)
		obj.Name = "test"
	})

	for _, v := range sliceTestStructs {
		if v.Name != "test" {
			t.Errorf("got unexpected value: %s", v.Name)
		}
	}

	For(testStructType).In(sliceTestStructs).Do(func(typeObj interface{}) {
		obj := typeObj.(TestStruct)
		if obj.Name != "test" {
			t.Errorf("got unexpected value: %s", obj.Name)
		}
	})
}

func TestMaps(t *testing.T) {
	mapTestStrings := make(map[int]*TestString)
	For(reflect.TypeOf(mapTestStrings)).In(mapTestStrings).Do(func(typeObj interface{}) {
		obj := typeObj.(map[int]*TestString)
		for i := 0; i < 100; i++ {
			str := TestString("test")
			obj[i] = &str
		}
	})

	for _, v := range mapTestStrings {
		if string(*v) != "test" {
			t.Errorf("got unexpected value: %s", *v)
		}
	}

	var i int
	For(ptrTestStringType).In(mapTestStrings).Do(func(typeObj interface{}) {
		obj := typeObj.(*TestString)
		*obj = TestString(fmt.Sprintf("test%02d", i))
		i++
	})

	// we can not assume the same order of map keys and values.
	for _, v := range mapTestStrings {
		if len(*v) != 6 {
			t.Errorf("got unexpected value: %s", *v)
		}
	}

	For(testStringType).In(mapTestStrings).Do(func(typeObj interface{}) {
		obj := typeObj.(TestString)
		if len(obj) != 6 {
			t.Errorf("got unexpected value: %s", obj)
		}
	})
}

func TestNestedStructures(t *testing.T) {
	someNestedStruct := TestNestedStruct{
		Name: "test",
		In: &TestNestedStruct{
			Name: "ptrTest",
		},
	}
	someNestedStruct.More = []TestNestedStruct{someNestedStruct, someNestedStruct}

	var names []string
	For(testNestedStructType).In(someNestedStruct).Do(func(typeObj interface{}) {
		obj := typeObj.(TestNestedStruct)
		names = append(names, string(obj.Name))
	})

	if len(names) != 6 {
		t.Errorf("got unexpected values: %s", names)
	}

	for i, v := range names {
		if i%2 == 0 {
			if v != "test" {
				t.Errorf("got unexpected value: %s", v)
			}
		} else {
			if v != "ptrTest" {
				t.Errorf("got unexpected value: %s", v)
			}
		}
	}
}
