package gxschema

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var preservedPropertyNames = [4]string{"id", "parent_id", "filename", "filepath"}
var propertyNamePattern = regexp.MustCompile(`^[_a-zA-Z][a-zA-Z0-9_\-]*$`)

//XMLNode raw XML node definition
//source: https://github.com/golang/go/issues/3633
type XMLNode struct {
	XMLName    xml.Name
	Attributes []xml.Attr
	Data       string
	Nodes      []XMLNode
}

//UnmarshalXML unmarshall XML
//source: https://github.com/golang/go/issues/3633
func (e *XMLNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var nodes []XMLNode
	var done bool
	for !done {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.CharData:
			e.Data = strings.TrimSpace(string(t))
		case xml.StartElement:
			e := &XMLNode{}
			e.UnmarshalXML(d, t)
			nodes = append(nodes, *e)
		case xml.EndElement:
			done = true
		}
	}
	e.XMLName = start.Name
	e.Attributes = start.Attr
	e.Nodes = nodes
	return nil
}

//MarshalXML marshall XML
//source: https://github.com/golang/go/issues/3633
func (e *XMLNode) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = e.XMLName
	start.Attr = e.Attributes
	return enc.EncodeElement(struct {
		Data  string `xml:",chardata"`
		Nodes []XMLNode
	}{
		Data:  e.Data,
		Nodes: e.Nodes,
	}, start)
}

//ParseSchemaFromXML parse document schema (DxDoc) from XML string
func ParseSchemaFromXML(rawXML string) (*DxDoc, error) {
	var n XMLNode

	marshallErr := xml.Unmarshal([]byte(rawXML), &n)
	if marshallErr != nil {
		return nil, marshallErr
	}

	dxdoc, errr := walkDxDoc(&n)
	if errr != nil {
		return nil, fmt.Errorf("failed to schema dxdoc: %s", errr.Error())
	}

	//travel all sub XML nodes
	for index, node := range n.Nodes {
		if strings.Compare(node.XMLName.Local, "dxbool") == 0 {
			dxbool, boolErr := walkDxBool(&node)
			if boolErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxbool at path dxdoc>dxbool(%d): %s",
					index, boolErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxbool)
		} else if strings.Compare(node.XMLName.Local, "dxint") == 0 {
			dxint, intErr := walkDxInt(&node)
			if intErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxint at path dxdoc>dxint(%d): %s",
					index, intErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxint)
		} else if strings.Compare(node.XMLName.Local, "dxdecimal") == 0 {
			dxdecimal, decimalErr := walkDxDecimal(&node)
			if decimalErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxdecimal at path dxdoc>dxdecimal(%d): %s",
					index, decimalErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxdecimal)
		} else if strings.Compare(node.XMLName.Local, "dxstr") == 0 {
			dxstr, strErr := walkDxStr(&node)
			if strErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxstr at path dxdoc.dxstr(%d): %s",
					index, strErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxstr)
		} else if strings.Compare(node.XMLName.Local, "dxsection") == 0 {
			dxsection, xmllPath, sectionErr := walkDxSection(&node, fmt.Sprintf("dxdoc>dxsection(%d)", index))
			if sectionErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxsection at path %s: %s",
					xmllPath, sectionErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxsection)
		} else if strings.Compare(node.XMLName.Local, "dxfile") == 0 {
			dxfile, fileErr := walkDxFile(&node)
			if fileErr != nil {
				return nil, fmt.Errorf(
					"failed to parse dxfile at path dxdoc.dxfile(%d): %s",
					index, fileErr.Error())
			}

			dxdoc.Items = append(dxdoc.Items, dxfile)
		} else {
			return nil, fmt.Errorf("unknown XML node %s found", node.XMLName.Local)
		}
	}

	if len(dxdoc.Items) == 0 {
		return nil, fmt.Errorf("DxDoc must atleast declare one data type definition")
	}

	return dxdoc, nil
}

func walkDxDoc(root *XMLNode) (*DxDoc, error) {
	if root.XMLName.Local != "dxdoc" {
		return nil, fmt.Errorf("expect tag name is dxdoc but get %s instead", root.XMLName.Local)
	}

	var err error

	hasRevision := false
	revision := 0

	hasName := false
	name := ""

	hasID := false
	var id string

	for _, attribute := range root.Attributes {
		//check revision attribute
		if isAttributeNameMatch(&attribute, "revision") {
			revision, err = parseAttributeInt(&attribute)
			if err != nil {
				return nil, fmt.Errorf("<dxdoc> tag %s", err.Error())
			}

			hasRevision = true
		}

		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "id") {
			id = attribute.Value
			hasID = true
		}
	}

	if !hasRevision {
		return nil, fmt.Errorf("missing attribute 'revision'")
	}

	if !hasName {
		return nil, fmt.Errorf("missing attribute 'name'")
	}

	if !hasID {
		return nil, fmt.Errorf("missing attribute 'id'")
	}

	//must has attribute 'revision', 'name' and child node(s) 'items'
	return &DxDoc{Revision: revision, Name: name, ID: id, Items: nil}, nil
}

func walkDxBool(node *XMLNode) (*DxBool, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}
	}

	if !hasName {
		return nil, fmt.Errorf("missing 'name' attribute")
	}

	return &DxBool{Name: name, IsOptional: optional, IsArray: array}, nil
}

func walkDxInt(node *XMLNode) (*DxInt, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}
	}

	if !hasName {
		return nil, fmt.Errorf("missing 'name' attribute")
	}

	return &DxInt{Name: name, IsOptional: optional, IsArray: array}, nil
}

func walkDxDecimal(node *XMLNode) (*DxDecimal, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	hasPrecision := false
	precision := 0

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "precision") {
			precision, err = parseAttributeInt(&attribute)
			if err != nil {
				return nil, err
			}

			hasPrecision = true
		}
	}

	if !hasName {
		return nil, fmt.Errorf("missing 'name' attribute")
	}

	if !hasPrecision {
		return nil, fmt.Errorf("missing attribute 'precision'")
	}

	return &DxDecimal{Name: name, IsOptional: optional, IsArray: array, Precision: precision}, nil
}

func walkDxStr(node *XMLNode) (*DxStr, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	limit := false
	len := 0

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "lenLimit") {
			limit = true
			len, err = parseAttributeInt(&attribute)
			if err != nil {
				return nil, err
			}
		}
	}

	if !hasName {
		return nil, fmt.Errorf("missing 'name' attribute")
	}

	return &DxStr{Name: name, IsOptional: optional, IsArray: array,
		EnableLenLimit: limit, LenLimit: len}, nil
}

func walkDxSection(node *XMLNode, xmlPath string) (*DxSection, string, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, "", err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, xmlPath, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, xmlPath, err
			}
		}
	}

	if !hasName {
		return nil, xmlPath, fmt.Errorf("missing 'name' attribute")
	}

	var items []DxItem

	for index, subNode := range node.Nodes {
		if strings.Compare(node.XMLName.Local, "dxbool") == 0 {
			dxbool, boolErr := walkDxBool(&subNode)
			if boolErr != nil {
				return nil, fmt.Sprintf("%s>dxbool(%d)", xmlPath, index), fmt.Errorf(
					"failed to parse dxbool, %s", boolErr.Error())
			}

			items = append(items, dxbool)
		} else if strings.Compare(node.XMLName.Local, "dxint") == 0 {
			dxint, intErr := walkDxInt(&subNode)
			if intErr != nil {
				return nil, fmt.Sprintf("%s>dxint(%d)", xmlPath, index), fmt.Errorf(
					"failed to parse dxint, %s", intErr.Error())
			}

			items = append(items, dxint)
		} else if strings.Compare(node.XMLName.Local, "dxdecimal") == 0 {
			dxdecimal, decimalErr := walkDxDecimal(&subNode)
			if decimalErr != nil {
				return nil, fmt.Sprintf("%s>dxdecimal(%d)", xmlPath, index), fmt.Errorf(
					"failed to parse dxdecimal, %s", decimalErr.Error())
			}

			items = append(items, dxdecimal)
		} else if strings.Compare(node.XMLName.Local, "dxstr") == 0 {
			dxstr, strErr := walkDxStr(&subNode)
			if strErr != nil {
				return nil, fmt.Sprintf("%s>dxstr(%d)", xmlPath, index), fmt.Errorf(
					"failed to parse dxstr, %s", strErr.Error())
			}

			items = append(items, dxstr)
		} else if strings.Compare(node.XMLName.Local, "dxsection") == 0 {
			dxSection, xmllPath, sectionErr := walkDxSection(
				&subNode, fmt.Sprintf("%s>dxsection(%d)", xmlPath, index))

			if sectionErr != nil {
				return nil, xmllPath, sectionErr
			}

			items = append(items, dxSection)
		} else {
			return nil, xmlPath, fmt.Errorf("unknown XML node")
		}
	}

	return &DxSection{Name: name, IsOptional: optional, IsArray: array, Items: items}, xmlPath, nil
}

func walkDxFile(node *XMLNode) (*DxFile, error) {
	var err error

	hasName := false
	name := ""

	optional := false
	array := false

	for _, attribute := range node.Attributes {
		if isAttributeNameMatch(&attribute, "name") {
			if err := validatePropertyName(attribute.Value); err != nil {
				return nil, err
			}

			name = attribute.Value
			hasName = true
		}

		if isAttributeNameMatch(&attribute, "isOptional") {
			optional, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}

		if isAttributeNameMatch(&attribute, "isArray") {
			array, err = parseAttributeBool(&attribute)
			if err != nil {
				return nil, err
			}
		}
	}

	if !hasName {
		return nil, fmt.Errorf("missing 'name' attribute")
	}

	return &DxFile{Name: name, IsOptional: optional, IsArray: array}, nil
}

//parseAttributeInt validate XML attribute is matching with provided name and it is integer type
func parseAttributeInt(attr *xml.Attr) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(attr.Value))
	if err != nil {
		return 0, fmt.Errorf("unable to parse int from attribute %s, attribute value: %s",
			attr.Name.Local, strings.TrimSpace(attr.Value))
	}

	return value, nil
}

func parseAttributeBool(attr *xml.Attr) (bool, error) {
	rawStr := strings.ToLower(attr.Value)

	if strings.Compare(rawStr, "true") == 0 {
		return true, nil
	} else if strings.Compare(rawStr, "false") == 0 {
		return false, nil
	}

	return false, fmt.Errorf("unable to parse bool from attribute %s, attribute value: %s",
		attr.Name.Local, attr.Value)
}

func isAttributeNameMatch(attr *xml.Attr, name string) bool {
	return strings.Compare(attr.Name.Local, name) == 0
}

func validatePropertyName(propertyName string) error {
	//validate is clash with preserved keyword or not
	for _, keyword := range preservedPropertyNames {
		if strings.Compare(keyword, strings.ToLower(propertyName)) == 0 {
			return fmt.Errorf("'%s' is preserved keyword, please avoid %v",
				propertyName, preservedPropertyNames)
		}
	}

	//validate string pattern
	if !propertyNamePattern.MatchString(propertyName) {
		return fmt.Errorf("'%s' is not a valid property name, "+
			"first charactor must begin with letter or underscore; "+
			"no white scpace or symbols allowed",
			propertyName)
	}

	return nil
}
