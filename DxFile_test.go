package gxschema

import "testing"

func TestDxFile_ValidateData(t *testing.T) {
	type args struct {
		input map[string]interface{}
		name  string
	}
	tests := []struct {
		name    string
		item    DxFile
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name: "simple file node test",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]string{"filename": "asdw.pdf", "filepath": "/sfgfd/dfgd/"}}},
			wantErr: false,
		},
		{
			name: "simple file node test 2",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]interface{}{"filename": "asdw.pdf", "filepath": "/sfgfd/dfgd/"}}},
			wantErr: false,
		},
		{
			name: "file node missing filename",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]string{"filepath": "/sfgfd/dfgd/"}}},
			wantErr: true,
		},
		{
			name: "file node missing filepath",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]string{"filename": "roth.xml"}}},
			wantErr: true,
		},
		{
			name: "file node missing filename 2",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]interface{}{"filepath": "/sfgfd/dfgd/"}}},
			wantErr: true,
		},
		{
			name: "file node missing filepath 2",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: false},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": map[string]interface{}{"filename": "roth.xml"}}},
			wantErr: true,
		},
		{
			name: "simple file node array test",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: true},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": []map[string]interface{}{
					map[string]interface{}{"filename": "asdw.pdf", "filepath": "/sfgfd/dfgd/45644.pdf"},
					map[string]interface{}{"filename": "bhu.pdf", "filepath": "/bnhg/yuhr/bfh5675.pdf"},
				}},
			},
			wantErr: false,
		},
		{
			name: "simple file node array test 2",
			item: DxFile{Name: "attachment", IsOptional: false, IsArray: true},
			args: args{name: "koko", input: map[string]interface{}{
				"koko": []map[string]string{
					map[string]string{"filename": "asdw.pdf", "filepath": "/sfgfd/dfgd/45644.pdf"},
					map[string]string{"filename": "bhu.pdf", "filepath": "/bnhg/yuhr/bfh5675.pdf"},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.ValidateData(tt.args.input, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DxFile.ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
