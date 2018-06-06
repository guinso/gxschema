package gxschema

import "testing"

func TestDxBool_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
		name  string
	}
	tests := []struct {
		name    string
		item    DxBool
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name:    "simple bool test",
			item:    DxBool{Name: "isMandatory", IsOptional: false, IsArray: false},
			args:    args{name: "koko", input: map[string]interface{}{"koko": true}},
			wantErr: false,
		},
		{
			name:    "simple bool array test",
			item:    DxBool{Name: "isMandatory", IsOptional: false, IsArray: true},
			args:    args{name: "koko", input: map[string]interface{}{"koko": []bool{true, false, true, true}}},
			wantErr: false,
		},
		{
			name:    "simple bool array test 2",
			item:    DxBool{Name: "isMandatory", IsOptional: false, IsArray: true},
			args:    args{name: "koko", input: map[string]interface{}{"koko": []interface{}{true, false, true, true}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.ValidateData(tt.args.input, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DxBool.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
