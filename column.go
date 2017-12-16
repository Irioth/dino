package dino

type ColumnFactory struct {
}

func (f ColumnFactory) NewColumnData() columnData {
	return &Vector{}
}

type Vector struct {
	data []interface{}
}

func (v *Vector) AddValue(index int, value interface{}) {
	for len(v.data) < index {
		v.data = append(v.data, nil)
	}
	v.data = append(v.data, value)
}
