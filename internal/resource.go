package internal

import (
	"reflect"

	"github.com/Neakxs/neatapi/resource"
)

type StructField struct {
	field reflect.StructField
}

func (f *StructField) GetTag(s string) string {
	value, _ := f.field.Tag.Lookup(s)
	return value
}

func (f *StructField) GetName() string {
	return f.field.Name
}

type MapEntry struct {
	*StructField
	Name  string
	Value interface{}
}

func PopulateResourceMap(r resource.Resource, m map[string]*MapEntry) error {
	return populateResourceMap(r, m, "", "")
}

func populateResourceMap(r resource.Resource, mapping map[string]*MapEntry, publicNs, privateNs string) error {
	rootValue := reflect.ValueOf(r)
	for rootValue.Kind() == reflect.Ptr {
		rootValue = rootValue.Elem()
	}
	for fieldNo := 0; fieldNo < rootValue.NumField(); fieldNo++ {
		indexField := &StructField{rootValue.Type().Field(fieldNo)}
		publicNames := r.GetPublicNames(indexField)
		privateName := r.GetPrivateName(indexField)
		fieldValue := rootValue.Field(fieldNo)
		if rr, ok := fieldValue.Interface().(resource.Resource); ok {
			if len(publicNames) == 1 && publicNames[0] == "" {
				populateResourceMap(rr, mapping, publicNs, privateNs)
			} else {
				for i := 0; i < len(publicNames); i++ {
					populateResourceMap(rr, mapping, publicNs+publicNames[i]+".", privateNs+privateName+".")
				}
			}
		} else {
			for _, publicName := range publicNames {
				mapping[publicNs+publicName] = &MapEntry{
					StructField: indexField,
					Name:        privateNs + privateName,
					Value:       fieldValue.Interface(),
				}
			}
		}
	}
	return nil
}
