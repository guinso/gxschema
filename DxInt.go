package gxschema

import (
	"fmt"
	"math"
	"reflect"
)

//DxInt integer data item
type DxInt struct {
	Name       string
	IsOptional bool
	IsArray    bool
}

//GetName get name
func (item DxInt) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxInt) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxInt) IsValueArray() bool { return item.IsArray }

//XML generate XML
func (item DxInt) XML(indentLevel int) string {
	var result string
	for i := 0; i < indentLevel; i++ {
		result += "\t"
	}
	result += "<dxint name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	return result + "></dxint>"
}

//ValidateData validate input data
func (item DxInt) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	}

	if item.IsArray {
		_, arrOK := rawValue.([]int)
		if arrOK {
			return nil
		}

		arrFloat, arrFloatOK := rawValue.([]float64)
		if arrFloatOK {
			for index, tmp := range arrFloat {
				if _, x2 := math.Modf(tmp); x2 != 0 {
					return fmt.Errorf("%s[%d] is not int value: %f", name, index, tmp)
				}
			}
			return nil
		}

		arrObj, arrObjOK := rawValue.([]interface{})
		if arrObjOK {
			for index, tmp := range arrObj {
				if _, ok := tmp.(int); ok {
					continue
				}

				if tmpf, ok := tmp.(float64); ok {
					if _, x2 := math.Modf(tmpf); x2 == 0 {
						continue
					}

					return fmt.Errorf("%s[%d] is not int value: %f", name, index, tmpf)
				}

				return fmt.Errorf("%s[%d] is not int value: %s", name, index, tmp)
			}

			return nil
		}

		return fmt.Errorf("input value is not array int but %s", reflect.TypeOf(rawValue))
	}

	_, intOK := rawValue.(int)
	if intOK {
		return nil
	}

	tmpFloat, floatOK := rawValue.(float64)
	if floatOK {
		_, x2 := math.Modf(tmpFloat)
		if x2 == 0 {
			return nil
		}

		return fmt.Errorf("%s is not int value: %f", name, x2)
	}

	return fmt.Errorf("input value is not int but %s", reflect.TypeOf(rawValue))
}
