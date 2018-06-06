package gxschema

import "testing"

func TestDxInt_ValidateData(t *testing.T) {
	intDef := DxInt{Name: "qty", IsOptional: false, IsArray: false}

	input := make(map[string]interface{})
	input["koko"] = 12
	input["gogo"] = []interface{}{12, 4, 5}
	input["popo"] = []int{4, 5, 6}
	input["zozo"] = []float64{7, 8, 9.5}
	input["soso"] = 12.5

	if err := intDef.ValidateData(input, "koko"); err != nil {
		t.Error(err)
	}

	intDef.IsArray = true
	if err := intDef.ValidateData(input, "gogo"); err != nil {
		t.Errorf("gogo: %s", err.Error())
	}
	if err := intDef.ValidateData(input, "popo"); err != nil {
		t.Errorf("popo: %s", err.Error())
	}
	if err := intDef.ValidateData(input, "zozo"); err == nil {
		t.Errorf("zozo: expect fail but pass")
	}
	if err := intDef.ValidateData(input, "soso"); err == nil {
		t.Errorf("soso: expect fail but pass")
	}
}
