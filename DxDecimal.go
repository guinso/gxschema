package gxschema

import (
	"fmt"
	"math"
	"reflect"

	"github.com/shopspring/decimal"
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
	} else if rawValue == nil && item.IsOptional {
		return nil
	}

	if item.IsArray {
		if arrFloat, arrOK := rawValue.([]float64); arrOK {
			for index, value := range arrFloat {
				if err := item.validateFloat(value,
					fmt.Sprintf("%s[%d]", name, index)); err != nil {
					return err
				}
			}

			return nil
		} else if arrDecimal, arrOK := rawValue.([]decimal.Decimal); arrOK {
			for index, value := range arrDecimal {
				if err := item.validateShopSpringDecimal(value,
					fmt.Sprintf("%s[%d]", name, index)); err != nil {
					return err
				}
			}

			return nil
		} else if arrObj, arrObjOK := rawValue.([]interface{}); arrObjOK {

			for index, rawValue := range arrObj {
				if value, floatOK := rawValue.(float64); floatOK {
					if err := item.validateFloat(value, fmt.Sprintf("%s[%d]", name, index)); err != nil {
						return err
					}
				} else if value, shopspringDecimalOK := rawValue.(decimal.Decimal); shopspringDecimalOK {
					if err := item.validateShopSpringDecimal(value, fmt.Sprintf("%s[%d]", name, index)); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("%s[%d] is not decimal but %s", name, index, reflect.TypeOf(rawValue))
				}
			}

			return nil
		} else if value, floatOK := rawValue.(float64); floatOK {
			return item.validateFloat(value, name)
		} else if value, shopspringDecimalOK := rawValue.(decimal.Decimal); shopspringDecimalOK {
			return item.validateShopSpringDecimal(value, name)
		}

		return fmt.Errorf("%s is not array decimal but %s", name, reflect.TypeOf(rawValue))
	}

	if value, floatOK := rawValue.(float64); floatOK {
		return item.validateFloat(value, name)
	} else if value, shopspringDecimalOK := rawValue.(decimal.Decimal); shopspringDecimalOK {
		return item.validateShopSpringDecimal(value, name)
	}

	return fmt.Errorf("%s is not decimal but %s", name, reflect.TypeOf(rawValue))
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

func (item DxDecimal) validateShopSpringDecimal(value decimal.Decimal, name string) error {
	precision := value.Exponent() * -1

	if precision > int32(item.Precision) {
		return fmt.Errorf("%s has invalid precision, expected %d: %s", name, item.Precision, value.String())
	}

	return nil
}
