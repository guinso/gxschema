package gxschema

//DxItem document item's interface
type DxItem interface {
	GetName() string                                              //GetName get item's name
	XML(indentLevel int) string                                   //XML generate into XML format
	ValidateData(input map[string]interface{}, name string) error //ValidateData check input data is matching with definition
	IsValueOptional() bool                                        //IsValueOptional is value optional
	IsValueArray() bool                                           //IsValueArray is the item allow to store more than 1 record
}
