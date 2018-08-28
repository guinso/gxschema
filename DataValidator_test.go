package gxschema

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_parseMapInterfaceFromXMLNode(t *testing.T) {
	node := &XMLNode{
		XMLName: xml.Name{Space: "", Local: "book"},
		Nodes: []XMLNode{
			XMLNode{
				XMLName: xml.Name{Space: "", Local: "year"},
				Data:    "1996",
			},
			XMLNode{
				XMLName: xml.Name{Space: "", Local: "name"},
				Data:    "James Foo",
			},
			XMLNode{
				XMLName: xml.Name{Space: "", Local: "price"},
				Data:    "12.45",
			},
		}}

	mapValue := parseMapInterfaceFromXMLNode(node)

	bookValue, bookOk := mapValue["book"]
	if !bookOk {
		t.Error("Book key not found!")
		return
	}
	bookMap, mapOk := bookValue.(map[string]interface{})
	if !mapOk {
		t.Error("Book value is not map[string]interface{} type")
		return
	}

	nameValue, nameErr := getMapValueString(bookMap, "name")
	if nameErr != nil {
		t.Error(nameErr)
		return
	}
	if strings.Compare(nameValue, "James Foo") != 0 {
		t.Errorf("Expect name value is 'James Foo' but get '%s'", nameValue)
		return
	}

	yearValue, yearErr := getMapValueInt(bookMap, "year")
	if yearErr != nil {
		t.Error(yearErr)
		return
	}
	if yearValue != 1996 {
		t.Errorf("Expect year value is '1996' but get '%d'", yearValue)
	}

	priceValue, priceErr := getMapValueShopSpringDecimal(bookMap, "price")
	if priceErr != nil {
		t.Error(priceErr)
		return
	}
	if priceValue.String() != "12.45" {
		t.Errorf("Expect price value is '12.45' but get '%s'", priceValue.String())
	}
}

func getMapValueInt(mapValue map[string]interface{}, key string) (int, error) {
	value, ok := mapValue[key]
	if !ok {
		return 0, fmt.Errorf("map key [%s] not found", key)
	}
	typeValue, intOk := value.(int)
	if !intOk {
		return 0, fmt.Errorf("%s is not int type: %s", key, reflect.TypeOf(typeValue).String())
	}

	return typeValue, nil
}

func getMapValueString(mapValue map[string]interface{}, key string) (string, error) {
	value, ok := mapValue[key]
	if !ok {
		return "", fmt.Errorf("map key [%s] not found", key)
	}
	typeValue, intOk := value.(string)
	if !intOk {
		return "", fmt.Errorf("%s is not string type: %s", key, reflect.TypeOf(typeValue).String())
	}

	return typeValue, nil
}

func getMapValueShopSpringDecimal(mapValue map[string]interface{}, key string) (decimal.Decimal, error) {
	value, ok := mapValue[key]
	if !ok {
		return decimal.Zero, fmt.Errorf("map key [%s] not found", key)
	}
	typeValue, decimalOk := value.(decimal.Decimal)
	if !decimalOk {
		return decimal.Zero, fmt.Errorf("%s is not shopspring decimal type: %s", key, reflect.TypeOf(typeValue).String())
	}

	return typeValue, nil
}

func TestValidateDataFromXML(t *testing.T) {
	docSchema := &DxDoc{
		Name:     "book",
		Revision: 1,
		Items: []DxItem{
			DxStr{
				Name:       "author",
				IsOptional: false,
			},
			DxInt{
				Name:       "year",
				IsOptional: false,
			},
			DxDecimal{
				Name:       "price",
				IsOptional: false,
				IsArray:    true,
				Precision:  2,
			},
		},
	}

	type args struct {
		dataXML   string
		docSchema *DxDoc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name: "test normal case",
			args: args{
				dataXML: `<book>
							<author>John</author>
							<year>1997</year>
							<price>12.45</price>
						</book>`,
				docSchema: docSchema,
			},
			wantErr: false,
		},
		{
			name: "test normal case",
			args: args{
				dataXML: `<book>
							<author>John</author>
							<year>1997</year>
							<price>12.45</price>
							<price>0.14</price>
						</book>`,
				docSchema: docSchema,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDataFromXML(tt.args.dataXML, tt.args.docSchema); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataFromXML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDataFromJSON(t *testing.T) {
	docSchema := &DxDoc{
		Name:     "book",
		Revision: 1,
		Items: []DxItem{
			DxStr{
				Name:       "author",
				IsOptional: false,
			},
			DxInt{
				Name:       "year",
				IsOptional: false,
			},
			DxDecimal{
				Name:       "price",
				IsOptional: false,
				IsArray:    true,
				Precision:  2,
			},
		},
	}

	type args struct {
		dataJSON  string
		docSchema *DxDoc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				dataJSON:  `{"author":"John", "year":1996, "price":[1.23,4.07]}`,
				docSchema: docSchema,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDataFromJSON(tt.args.dataJSON, tt.args.docSchema); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataFromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
