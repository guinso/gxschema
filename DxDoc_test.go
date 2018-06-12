package gxschema

import (
	"strings"
	"testing"
)

func TestXML(t *testing.T) {
	doc := DxDoc{
		Name:     "invoice",
		Revision: 3,
		ID:       "733bee1b-f79a-4cb7-b675-842317b994b5",
		Items: []DxItem{
			DxDecimal{Name: "total", Precision: 2},
			DxStr{Name: "doc no", LenLimit: 6},
			DxSection{Name: "items", IsArray: true, Items: []DxItem{
				DxStr{Name: "description"},
				DxDecimal{Name: "price", Precision: 2},
				DxStr{Name: "remark", IsOptional: true},
			}},
		},
	}

	xmlStr, xmlErr := doc.XML()
	if xmlErr != nil {
		t.Error(xmlErr)
	}

	expectedXML := `<?xml version="1.0"?>
<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
	<dxdecimal name="total" precision="2"></dxdecimal>
	<dxstr name="doc no"></dxstr>
	<dxsection name="items" isArray="true">
		<dxstr name="description"></dxstr>
		<dxdecimal name="price" precision="2"></dxdecimal>
		<dxstr name="remark" isOptional="true"></dxstr>
	</dxsection>
</dxdoc>`

	if strings.Compare(xmlStr, expectedXML) != 0 {
		t.Errorf("XML output not tally with [output]: \n%s\n\n[expected]:\n%s", xmlStr, expectedXML)
	}
}

func TestDxDoc_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
	}
	tests := []struct {
		name    string
		doc     DxDoc
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "simple DxDoc test",
			doc: DxDoc{Name: "invoice", Revision: 2, Items: []DxItem{
				DxStr{Name: "docNo", EnableLenLimit: true, LenLimit: 4},
				DxSection{Name: "items", IsArray: true, Items: []DxItem{
					DxStr{Name: "description"},
					DxInt{Name: "qty"},
					DxDecimal{Name: "unit price", Precision: 2},
				}},
			}},
			args: args{input: map[string]interface{}{
				"docNo": "abcd",
				"items": []map[string]interface{}{
					map[string]interface{}{"description": "dfgfghfh", "qty": 3, "unit price": 3.50},
					map[string]interface{}{"description": "nrrty erte", "qty": 2, "unit price": 12.05},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.doc.ValidateData(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("DxDoc.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
