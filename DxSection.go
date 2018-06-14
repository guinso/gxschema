package gxschema

import (
	"fmt"
	"strings"
)

//DxSection group DxItem(s) data type
type DxSection struct {
	Name       string
	IsOptional bool
	IsArray    bool
	Items      []DxItem
}

//GetName get name
func (item DxSection) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxSection) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxSection) IsValueArray() bool { return item.IsArray }

//XML generate XML
func (item DxSection) XML(indentLevel int) string {
	var indent string
	for i := 0; i < indentLevel; i++ {
		indent = "\t"
	}
	result := indent + "<dxsection name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	result += ">"

	for _, item := range item.Items {
		result += "\n" + item.XML(indentLevel+1)
	}

	return result + "\n" + indent + "</dxsection>"
}

//ValidateData validate input data
func (item DxSection) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	}

	if item.IsArray {
		subArr, subOK := rawValue.([]map[string]interface{})
		if !subOK {
			return fmt.Errorf("%s is not map array", name)
		}

		//iterate each array item and validate its value
		for index, tmp := range subArr {
			if err := item.validateItem(tmp, fmt.Sprintf("%s[%d]", name, index)); err != nil {
				return err
			}
		}

		return nil
	}

	return item.validateItem(rawValue, name)
}

func (item DxSection) validateItem(rawValue interface{}, name string) error {
	subItem, subOK := rawValue.(map[string]interface{})
	if !subOK {
		return fmt.Errorf("%s is not map", name)
	}

	//build checkmark
	checkMark := make([]int, len(item.Items))

	for key := range subItem {
		defIndex, def := item.findItem(key)
		if def == nil {
			continue
		}

		validErr := def.ValidateData(subItem, key)
		if validErr != nil {
			return validErr
		}

		checkMark[defIndex]++
	}

	for i := 0; i < len(item.Items); i++ {
		if checkMark[i] == 0 && !item.Items[i].IsValueOptional() {
			return fmt.Errorf("%s has no key '%s'", name, item.Items[i].GetName())
		}
	}

	return nil
}

func (item DxSection) findItem(name string) (int, DxItem) {
	for index, tmp := range item.Items {
		if strings.Compare(name, tmp.GetName()) == 0 {
			return index, tmp
		}
	}

	return -1, nil
}
