package utils

import "reflect"

func TypeName(object any) string {
	t := reflect.TypeOf(object)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
