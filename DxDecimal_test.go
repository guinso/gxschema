package gxschema

import "testing"

func TestDxDecimal_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
		name  string
	}
	tests := []struct {
		name    string
		item    DxDecimal
		args    args
		wantErr bool
	}{
		//Add test cases.
		{
			name:    "simple decimal test",
			item:    DxDecimal{Name: "price", IsOptional: false, IsArray: false, Precision: 2},
			args:    args{name: "koko", input: map[string]interface{}{"koko": 12.34}},
			wantErr: false,
		},
		{
			name:    "incorrect decimal precision test",
			item:    DxDecimal{Name: "price", IsOptional: false, IsArray: false, Precision: 2},
			args:    args{name: "koko", input: map[string]interface{}{"koko": 12.3401}},
			wantErr: true,
		},
		{
			name:    "simple decimal array test",
			item:    DxDecimal{Name: "price", IsOptional: false, IsArray: true, Precision: 3},
			args:    args{name: "koko", input: map[string]interface{}{"koko": []float64{12.345, 0.078}}},
			wantErr: false,
		},
		{
			name:    "simple decimal array test 2",
			item:    DxDecimal{Name: "price", IsOptional: false, IsArray: true, Precision: 3},
			args:    args{name: "koko", input: map[string]interface{}{"koko": []interface{}{12.345, 0.078}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.ValidateData(tt.args.input, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DxDecimal.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
