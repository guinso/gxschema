package gxschema

import (
	"fmt"
	"math"
	"reflect"
)

//DxDecimal floating point data type
type DxDecimal struct {
	Name       string
	IsOptional bool
	IsArray    bool
	Precision  int //decimal precision
}

//GetName get name
func (item DxDecimal) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxDecimal) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxDecimal) IsValueArray() bool { return item.IsArray }

//XML generate XML
func (item DxDecimal) XML(indentLevel int) string {
	var result string
	for i := 0; i < indentLevel; i++ {
		result += "\t"
	}
	result += "<dxdecimal name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	return result + fmt.Sprintf(" precision=\"%d\"></dxdecimal>", item.Precision)
}

//ValidateData validate input data
func (item DxDecimal) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	}

	if item.IsArray {
		arrFloat, arrOK := rawValue.([]float64)
		if arrOK {
			for index, value := range arrFloat {
				if err := item.validateFloat(value, fmt.Sprintf("%s[%d]", name, index)); err != nil {
					return err
				}
			}

			return nil
		}

		arrObj, arrObjOK := rawValue.([]interface{})
		if arrObjOK {
			for index, rawValue := range arrObj {
				value, floatOK := rawValue.(float64)
				if floatOK {
					if err := item.validateFloat(value, fmt.Sprintf("%s[%d]", name, index)); err != nil {
						return err
					}
				}
			}

			return nil
		}

		return fmt.Errorf("%s is not array boolean but %s", name, reflect.TypeOf(rawValue))
	}

	value, floatOK := rawValue.(float64)
	if floatOK {
		return item.validateFloat(value, name)
	}

	return fmt.Errorf("%s is not floating number but %s", name, reflect.TypeOf(rawValue))
}

func (item DxDecimal) validateFloat(value float64, name string) error {
	newValue := value

	if item.Precision > 0 {
		newValue *= math.Pow10(item.Precision)
	}

	if _, residue := math.Modf(newValue); residue != 0 {
		return fmt.Errorf("%s has invalid precision, expected %d: %f",
			name, item.Precision, value)
	}

	return nil
}
