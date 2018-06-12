package gxschema

import (
	"fmt"
	"strings"
)

//DxMap shorthand of map[string]interface{}
//type DxMap map[string]interface{}

//DxDoc Document schema
type DxDoc struct {
	Name     string   //Name document name
	ID       string   //ID document unique identifier (suggest UUID)
	Revision int      //Revision document revision, each changes of document structure revision value shall increament by 1
	Items    []DxItem //Items document contents, each item represent single field of document
}

//XML generate document definition into XML format
func (doc DxDoc) XML() (string, error) {
	result := fmt.Sprintf(
		"<?xml version=\"1.0\"?>\n<dxdoc name=\"%s\" revision=\"%d\" id=\"%s\">",
		doc.Name, doc.Revision, doc.ID)

	for _, item := range doc.Items {
		result += "\n" + item.XML(1)
	}

	return result + "\n</dxdoc>", nil
}

//ValidateData check input data integration with present DxDoc definition instance
func (doc DxDoc) ValidateData(input map[string]interface{}) error {
	checkMark := make([]int, len(doc.Items))

	for key := range input {
		tmpIndex, tmpItem := doc.findItem(key)
		if tmpItem == nil {
			continue
		}

		if err := tmpItem.ValidateData(input, key); err != nil {
			return err
		}

		checkMark[tmpIndex]++
	}

	for i := 0; i < len(doc.Items); i++ {
		if checkMark[i] == 0 && !doc.Items[i].IsValueOptional() {
			return fmt.Errorf("'%s' not found in %s", doc.Items[i].GetName(), doc.Name)
		}
	}

	return nil
}

func (doc DxDoc) findItem(name string) (int, DxItem) {
	for i := 0; i < len(doc.Items); i++ {
		if strings.Compare(name, doc.Items[i].GetName()) == 0 {
			return i, doc.Items[i]
		}
	}

	return -1, nil
}

//DxItem document item's interface
type DxItem interface {
	GetName() string                                              //GetName get item's name
	XML(indentLevel int) string                                   //XML generate into XML format
	ValidateData(input map[string]interface{}, name string) error //ValidateData check input data is matching with definition
	IsValueOptional() bool                                        //IsOptional is value optional
}
