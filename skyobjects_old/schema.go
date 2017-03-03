package skyobject

import (

	// "bytes"

	"reflect"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Schema represents the type of Hashed item.
type Schema struct {
	Name   string                `json:"name"`
	Fields []encoder.StructField `json:"fields"`
}

// ReadSchema reads the schema from an object.
func ReadSchema(data interface{}) (sch Schema) {
	st := reflect.TypeOf(data)
	sv := reflect.ValueOf(data)
	sch.Name = st.Name()
	for i := 0; i < st.NumField(); i++ {
		if st.Field(i).Tag.Get("enc") != "-" {
			sch.Fields = append(
				sch.Fields,
				getField(st.Field(i), sv.Field(i)),
			)
		}
	}
	return
}

func (s *Schema) String() string {
	var b = make([]byte, 0, 96) // scratch
	b = append(b, "struct "...)
	b = append(b, s.Name...)
	b = append(b, '\n')
	for i := 0; i < len(s.Fields); i++ {
		b = append(b, s.Fields[i].String()...)
	}
	return string(b)
}

func getField(field reflect.StructField,
	fieldValue reflect.Value) encoder.StructField {

	fieldType := strings.ToLower(fieldValue.Type().Name())
	return encoder.StructField{
		Name: field.Name,
		Type: fieldType,
		Tag:  string(field.Tag),
		Kind: uint32(fieldValue.Type().Kind()),
	}

}