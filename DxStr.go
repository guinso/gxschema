package gxschema

import (
	"fmt"
	"reflect"
)

//DxStr string data type
type DxStr struct {
	Name           string
	IsOptional     bool
	IsArray        bool
	EnableLenLimit bool
	LenLimit       int
}

//GetName get name
func (item DxStr) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxStr) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxStr) IsValueArray() bool { return item.IsArray }

//XML generate XML
func (item DxStr) XML(indentLevel int) string {
	var result string
	for i := 0; i < indentLevel; i++ {
		result += "\t"
	}
	result += "<dxstr name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	if item.EnableLenLimit {
		result += fmt.Sprintf(" lenLimit=\"%d\"", item.LenLimit)
	}

	return result + "></dxstr>"
}

//ValidateData validate input data
func (item DxStr) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	}

	if item.IsArray {
		strArr, arrOK := rawValue.([]string)
		if arrOK {
			if item.EnableLenLimit {
				for index, tmp := range strArr {
					if len(tmp) != item.LenLimit {
						return fmt.Errorf("%s[%d] length is not %d: %s",
							name, index, item.LenLimit, tmp)
					}
				}
			}

			return nil
		}

		arrStr, arrStrOK := rawValue.([]interface{})
		if arrStrOK {
			for index, tmp := range arrStr {
				tmpStr, OK := tmp.(string)
				if !OK {
					return fmt.Errorf("%s[%d] is not string but %s", name, index, reflect.TypeOf(tmp))
				}

				if item.EnableLenLimit && len(tmpStr) != item.LenLimit {
					return fmt.Errorf("%s[%d] length is not %d:%s", name, index, item.LenLimit, tmpStr)
				}
			}

			return nil
		}

		return fmt.Errorf("%s is not array string but %s", name, reflect.TypeOf(rawValue))
	}

	str, intOK := rawValue.(string)
	if !intOK {
		return fmt.Errorf("%s is not string but %s", name, reflect.TypeOf(rawValue))
	}

	if item.EnableLenLimit {
		if len(str) != item.LenLimit {
			return fmt.Errorf("%s length is not %d: %s",
				name, item.LenLimit, str)
		}
	}

	return nil
}
