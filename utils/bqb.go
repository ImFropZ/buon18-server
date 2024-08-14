package utils

import (
	"reflect"
	"server/models"

	"github.com/nullism/bqb"
)

func PrepareUpdateBqbQuery[T models.CommonUpdateRequest](bqbQuery *bqb.Query, i *T) {
	val := reflect.ValueOf(i).Elem()
	typ := reflect.TypeOf(i).Elem()

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for index := range val.NumField() {
		if val.Field(index).IsNil() || val.Field(index).Kind() == reflect.Array || val.Field(index).Kind() == reflect.Slice {
			continue
		}

		err := (*i).MapUpdateFields(bqbQuery, typ.Field(index).Name, val.Field(index).Interface())
		if err != nil {
			continue
		}
	}
}
