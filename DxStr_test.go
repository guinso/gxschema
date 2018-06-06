package gxschema

import "testing"

func TestDxStr_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
		name  string
	}
	tests := []struct {
		name    string
		item    DxStr
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name:    "simple string test",
			item:    DxStr{Name: "customer", IsOptional: false, IsArray: false, EnableLenLimit: false, LenLimit: 0},
			args:    args{name: "nono", input: map[string]interface{}{"nono": "qwe asd"}},
			wantErr: false},
		{
			name:    "simple string array test",
			item:    DxStr{Name: "customer", IsOptional: false, IsArray: true, EnableLenLimit: false, LenLimit: 0},
			args:    args{name: "nono", input: map[string]interface{}{"nono": []string{"qwe", "asd", "zxc"}}},
			wantErr: false},
		{
			name:    "string length limit test",
			item:    DxStr{Name: "customer", IsOptional: false, IsArray: false, EnableLenLimit: true, LenLimit: 4},
			args:    args{name: "nono", input: map[string]interface{}{"nono": "qasd"}},
			wantErr: false},
		{
			name:    "string array length limit test",
			item:    DxStr{Name: "customer", IsOptional: false, IsArray: true, EnableLenLimit: true, LenLimit: 6},
			args:    args{name: "nono", input: map[string]interface{}{"nono": []string{"qwerty", "123456", "zxcghj"}}},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.ValidateData(tt.args.input, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DxStr.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
