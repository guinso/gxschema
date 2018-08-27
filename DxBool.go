package gxschema

import (
	"fmt"
	"reflect"
)

//DxBool boolean data type
type DxBool struct {
	Name       string
	IsOptional bool
	IsArray    bool
}

//GetName get name
func (item DxBool) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxBool) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxBool) IsValueArray() bool { return item.IsArray }

//XML generate XML
func (item DxBool) XML(indentLevel int) string {
	var result string
	for i := 0; i < indentLevel; i++ {
		result += "\t"
	}
	result += "<dxbool name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	return result + "></dxbool>"
}

//ValidateData validate input data
func (item DxBool) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	} else if rawValue == nil && item.IsOptional {
		return nil
	}

	if item.IsArray {
		_, arrOK := rawValue.([]bool)
		if arrOK {

			return nil
		}

		arrBool, arrBoolOK := rawValue.([]interface{})
		if arrBoolOK {
			for index, tmp := range arrBool {
				_, OK := tmp.(bool)
				if !OK {
					return fmt.Errorf("%s[%d] is not boolean but %s", name, index, reflect.TypeOf(tmp))
				}
			}

			return nil
		} else if _, intOK := rawValue.(bool); intOK {
			return nil
		}

		return fmt.Errorf("%s is not array boolean but %s", name, reflect.TypeOf(rawValue))
	}

	_, intOK := rawValue.(bool)
	if !intOK {
		return fmt.Errorf("%s is not boolean but %s", name, reflect.TypeOf(rawValue))
	}

	return nil
}
