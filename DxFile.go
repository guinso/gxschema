package gxschema

import (
	"fmt"
	"reflect"
)

//DxFile file data type
//file contents of:
//		1. file path
//		2. file name
type DxFile struct {
	Name       string
	IsOptional bool
	IsArray    bool
}

//GetName get name
func (item DxFile) GetName() string { return item.Name }

//IsValueOptional is field value optional
func (item DxFile) IsValueOptional() bool { return item.IsOptional }

//IsValueArray is field value allow to store multiple values
func (item DxFile) IsValueArray() bool { return item.IsArray }

//XML generate definitino into XML format
func (item DxFile) XML(indentLevel int) string {
	var result string
	for i := 0; i < indentLevel; i++ {
		result += "\t"
	}
	result += "<dxfile name=\"" + item.Name + "\""

	if item.IsArray {
		result += " isArray=\"true\""
	}

	if item.IsOptional {
		result += " isOptional=\"true\""
	}

	return result + "></dxfile>"
}

//ValidateData validate input data
func (item DxFile) ValidateData(input map[string]interface{}, name string) error {
	rawValue, keyOK := input[name]

	if !keyOK {
		if !item.IsOptional {
			return fmt.Errorf("map entry '%s' is not exists", name)
		}

		return nil
	}

	if item.IsArray {
		arrMap, arrOK := rawValue.([]map[string]interface{})
		if arrOK {
			for index, tmpMap := range arrMap {
				if err := item.validateNode(tmpMap, fmt.Sprintf("%s[%d]", name, index)); err != nil {
					return err
				}
			}

			return nil
		}

		arrMapStr, arrMapStrOK := rawValue.([]map[string]string)
		if arrMapStrOK {
			for index, tmpMap := range arrMapStr {
				if err := item.validateNodeV2(tmpMap, fmt.Sprintf("%s[%d]", name, index)); err != nil {
					return err
				}
			}

			return nil
		}

		return fmt.Errorf("%s is not array map but %s", name, reflect.TypeOf(rawValue))
	}

	interfaceMap, interfaceMapOK := rawValue.(map[string]interface{})
	if interfaceMapOK {
		return item.validateNode(interfaceMap, name)
	}

	strMap, strMapOK := rawValue.(map[string]string)
	if strMapOK {
		return item.validateNodeV2(strMap, name)
	}

	return fmt.Errorf("%s is not array map but %s", name, reflect.TypeOf(rawValue))
}

func (item DxFile) validateNode(tmpMap map[string]interface{}, name string) error {
	//filename node
	filenameRaw, filenameOK := tmpMap["filename"]
	if !filenameOK {
		return fmt.Errorf("%s has no 'filename' node", name)
	}

	_, OK := filenameRaw.(string)
	if !OK {
		return fmt.Errorf("%s['filename'] value is not string: %s",
			name, reflect.TypeOf(filenameRaw))
	}

	//filepath node
	filepathRaw, filepathOK := tmpMap["filepath"]
	if !filepathOK {
		return fmt.Errorf("%s has no 'filepath' node", name)
	}

	_, OK = filepathRaw.(string)
	if !OK {
		return fmt.Errorf("%s['filepath'] value is not string: %s",
			name, reflect.TypeOf(filepathRaw))
	}

	return nil
}

func (item DxFile) validateNodeV2(tmpMap map[string]string, name string) error {
	//filename node
	_, filenameOK := tmpMap["filename"]
	if !filenameOK {
		return fmt.Errorf("%s has no 'filename' node", name)
	}

	//filepath node
	_, filepathOK := tmpMap["filepath"]
	if !filepathOK {
		return fmt.Errorf("%s has no 'filepath' node", name)
	}

	return nil
}
