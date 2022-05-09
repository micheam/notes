package testing

import (
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"
)

type entityMatcher struct {
	x any
}

func (e entityMatcher) Matches(x any) bool {
	if e.x == nil || x == nil {
		return reflect.DeepEqual(e.x, x)
	}

	x1Val := reflect.ValueOf(e.x).Elem()
	x2Val := reflect.ValueOf(e.x).Elem()

	id1 := x1Val.FieldByName("ID")
	id2 := x2Val.FieldByName("ID")
	return reflect.DeepEqual(id1.Interface(), id2.Interface())
}

func (e entityMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%[1]T)", e.x)
}

func EqEntity(x any) gomock.Matcher { return entityMatcher{x} }
