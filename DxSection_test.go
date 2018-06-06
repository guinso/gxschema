package gxschema

import "testing"

func TestDxSection_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
		name  string
	}
	tests := []struct {
		name    string
		item    DxSection
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name: "simple section test",
			item: DxSection{Name: "items", IsOptional: false, IsArray: false, Items: []DxItem{
				DxStr{Name: "description"},
				DxDecimal{Name: "price", Precision: 2},
				DxBool{Name: "mandatory"},
			}},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]interface{}{
					"description": "Cap Kapak winter oil",
					"price":       12.40,
					"mandatory":   false,
				},
			}},
			wantErr: false,
		},
		{
			name: "simple section array test",
			item: DxSection{Name: "items", IsOptional: false, IsArray: true, Items: []DxItem{
				DxStr{Name: "description"},
				DxDecimal{Name: "price", Precision: 2},
				DxBool{Name: "mandatory"},
			}},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": []map[string]interface{}{
					map[string]interface{}{
						"description": "Cap Kapak winter oil",
						"price":       12.40,
						"mandatory":   false,
					},
					map[string]interface{}{
						"description": "Mamee",
						"price":       0.50,
						"mandatory":   true,
					},
					map[string]interface{}{
						"description": "Cap Keluarga rice 10KG",
						"price":       20.30,
						"mandatory":   false,
					},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.ValidateData(tt.args.input, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DxSection.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
