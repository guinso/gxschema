# gxschema
Data validation tool, a simplified version of XML schema.

# TODO
- Support Date and Time data format
- Export definition into XSD format

## Sample Format
### Schema (XML)
```xml
<?xml version="1.0"?>
<dxdoc name="order" revision="3">
    <dxstr name="order number" lenLimit="7"></dxstr>
    <dxint name="qty"></dxint>
    <dxdecimal name="rate" precision="2"></dxdecimal>
    <dxbool name="is member"></dxbool>
    <dxsection name="customer info">
        <dxstr name="name"></dxstr>
        <dxstr name="id" lenLimit="6"></dxstr>
    </dxsection>
    <dxsection name="items" isArray="true">
        <dxstr name="description"></dxstr>
        <dxint name="qty"></dxint>
        <dxdecimal name="unit price" precision="2"></dxdecimal>
    </dxsection>
</dxdoc>
```

### Data (JSON)
```json
{
    "order number": "ODR0001",
    "qty": 10,
    "rate": 12.56,
    "is member": true,
    "customer info":{
        "name":"John",
        "id":"cust01"
    },
    "items":[
        {"description": "Cap Kapak winter oil", "qty": 1, "unit price":3.50},
        {"description": "Lucky coffee powder", "qty": 3, "unit price":0.60},
    ]
}
```

## Example 1
load schema definition from XML and input data from JSON string
```go
//define XML schema
defRaw := `
<dxdoc name="invoice" revision="1">
    <dxstr name="invNo"></dxstr>
    <dxint name="totalQty" isOptional="true"></dxint>
    <dxdecimal name="price" precision="2"></dxdecimal>
</dxdoc>`

//parse XML schema into GO instance
dxdoc, dxErr := document.DecodeDxXML(defRaw)
if dxErr != nil {
    panic(dxErr) //show parsed failed message
}

//define input data (JSON string)
inputJSON := `
{
	"invNo": "abcd",
	"totalQty": 3,
	"price": 12.58
}`

//convert JSON string into GO generic MAP instance
rawInput := make(map[string]interface{})
jsonErr := json.unmarshal([]byte(inputJSON), &rawInput) 

//validate input data
validateErr := dxdoc.ValidateData(rawInput) 

if validateErr != nil {
    log.Println(validateErr.Error()) //show invalid data message
} else {
    log.PrintLn("input data is valid")
}
```

## Example 2
Direct define definition and input data
```go
//define definition
dxdoc := document.DxDoc{
    Name:     "invoice",
    Revision: 1,
    Items: []document.DxItem{
        document.DxStr{Name: "invNo", EnableLenLimit: true, LenLimit: 4},
        document.DxInt{Name: "totalQty", IsOptional: true},
        document.DxDecimal{Name: "price", Precision: 2},
    },
}

//define input data
rawInput := map[string]interface{}{
    "invNo": "abcd",
    "totalQty": 3,
    "price": 12.58,
}

//validate input data
validateErr := dxdoc.ValidateData(rawInput) 

if validateErr != nil {
    log.Println(validateErr.Error()) //show invalid data message
} else {
    log.PrintLn("input data is valid")
}
```