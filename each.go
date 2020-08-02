package every

import (
	"reflect"
)

type F struct {
	typ reflect.Type
}

func For(typ reflect.Type) *F {
	return &F{typ: typ}
}

func (f *F) In(obj interface{}) *I {
	return &I{
		F:   f,
		obj: obj,
	}
}

type I struct {
	*F
	obj interface{}
}

func (o *I) Do(f func(typeObj interface{})) {
	v := reflect.ValueOf(o.obj)
	if v.Type() == o.typ {
		f(v.Interface())
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			o.In(v.Field(i).Interface()).Do(f)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			o.In(v.Index(i).Interface()).Do(f)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			o.In(v.MapIndex(key).Interface()).Do(f)
		}
	case reflect.Ptr:
		element := v.Elem()
		if !element.IsValid() {
			return
		}
		o.In(element.Interface()).Do(f)
	}
}
