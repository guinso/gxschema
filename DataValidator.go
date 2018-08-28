package gxschema

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

//ValidateDataFromXML validate data based on XML string
func ValidateDataFromXML(dataXML string, docSchema *DxDoc) error {
	var n XMLNode

	marshallErr := xml.Unmarshal([]byte(dataXML), &n)
	if marshallErr != nil {
		return fmt.Errorf("Failed to parse XML: %s", marshallErr.Error())
	}

	mapValue := parseMapInterfaceFromXMLNode(&n)

	if tmpMap, mapOK := mapValue[docSchema.Name].(map[string]interface{}); mapOK {
		return docSchema.ValidateData(tmpMap)
	}

	return fmt.Errorf("Invalid XML root element, expect %s: %s", docSchema.Name, dataXML)
}

//parseMapInterfaceFromXMLNode parse map[string]interface{} from XMLNode
func parseMapInterfaceFromXMLNode(node *XMLNode) map[string]interface{} {
	inputData := make(map[string]interface{})

	if len(node.Nodes) > 0 {
		subData := make(map[string][]interface{})

		for _, subNode := range node.Nodes {
			tmpData := parseMapInterfaceFromXMLNode(&subNode)

			if tmpSlice, ok := subData[subNode.XMLName.Local]; ok {
				subData[subNode.XMLName.Local] = append(tmpSlice, tmpData[subNode.XMLName.Local])
			} else {
				subData[subNode.XMLName.Local] = append(make([]interface{}, 0), tmpData[subNode.XMLName.Local])
			}
		}

		tmpSubData := make(map[string]interface{})
		for sliceKey, tmpSlice := range subData {
			if len(tmpSlice) > 1 {
				tmpSubData[sliceKey] = tmpSlice
			} else {
				tmpSubData[sliceKey] = tmpSlice[0]
			}
		}

		if len(tmpSubData) > 0 {
			inputData[node.XMLName.Local] = tmpSubData
		} else {
			inputData[node.XMLName.Local] = nil
		}

	} else if integerValue, parseErr := strconv.Atoi(node.Data); parseErr == nil {
		inputData[node.XMLName.Local] = integerValue
	} else if booleanValue, parseErr := strconv.ParseBool(node.Data); parseErr == nil {
		inputData[node.XMLName.Local] = booleanValue
	} else if decimalValue, parseErr := decimal.NewFromString(node.Data); parseErr == nil {
		inputData[node.XMLName.Local] = decimalValue
	} else if len(node.Data) > 0 {
		inputData[node.XMLName.Local] = node.Data
	} else {
		inputData[node.XMLName.Local] = nil
	}

	return inputData
}

//ValidateDataFromJSON validate data based on JSON string
func ValidateDataFromJSON(dataJSON string, docSchema *DxDoc) error {
	rawMap := make(map[string]interface{})

	parseErr := json.Unmarshal([]byte(dataJSON), &rawMap)
	if parseErr != nil {
		return fmt.Errorf("Failed to parse JSON string: %s", parseErr.Error())
	}

	return docSchema.ValidateData(rawMap)
}
