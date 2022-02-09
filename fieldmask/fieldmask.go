package fieldmask

import (
	"fmt"

	"github.com/Neakxs/neatapi/internal"
	"github.com/Neakxs/neatapi/resource"
)

type Maskabler interface {
	resource.Resource
	Maskable(f resource.Field) bool
}

type FieldMask interface {
	Validate(r resource.Resource) error
	BuildMaskEntries(r resource.Resource) ([]*FieldMaskEntry, error)
}

type fieldMask struct {
	mapping map[string]*internal.MapEntry
	paths   []string
}

type FieldMaskEntry internal.MapEntry

func NewFieldMask(paths ...string) FieldMask {
	return &fieldMask{paths: paths}
}

func (m *fieldMask) populateResourceMap(r resource.Resource) error {
	if m.mapping == nil {
		m.mapping = make(map[string]*internal.MapEntry)
		var err error
		if iface, ok := r.(Maskabler); ok {
			err = internal.PopulateResourceMap(r, m.mapping, iface.Maskable)
		} else {
			err = internal.PopulateResourceMap(r, m.mapping, func(f resource.Field) bool { return true })
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *fieldMask) Validate(r resource.Resource) error {
	if err := m.populateResourceMap(r); err != nil {
		return err
	}
	fmt.Println(m.mapping)
	if iface, ok := r.(Maskabler); ok {
		for k, v := range m.mapping {
			if !iface.Maskable(v) {
				return fmt.Errorf("%v: not allowed", k)
			}
		}
	}
	if iface, ok := r.(resource.FieldValidater); ok {
		for _, v := range m.mapping {
			if err := iface.ValidateField(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *fieldMask) BuildMaskEntries(r resource.Resource) ([]*FieldMaskEntry, error) {
	if err := m.populateResourceMap(r); err != nil {
		return nil, err
	}
	res := make([]*FieldMaskEntry, 0)
	for i := 0; i < len(m.paths); i++ {
		if fld, ok := m.mapping[m.paths[i]]; ok {
			res = append(res, &FieldMaskEntry{m.mapping[m.paths[i]].StructField, fld.Name, fld.Value})
		}
	}
	return res, nil
}
