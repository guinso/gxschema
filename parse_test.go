package gxschema

import (
	"strings"
	"testing"
)

func TestWalkDxDoc(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3">
		<dxint name="age"></dxint>
		<dxdecimal name="price" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	dx, dxErr := DecodeDxXML(rawXML)
	if dxErr != nil {
		t.Error(dxErr)
		return
	}

	if dx.Revision != 3 {
		t.Errorf("expect revision value is 3 but get %d instead", dx.Revision)
	}

	if strings.Compare(dx.Name, "invoice") != 0 {
		t.Errorf("expect name value is 'invoice' but get '%s' instead", dx.Name)
	}

	if len(dx.Items) != 6 {
		t.Errorf("Expect has 6 items definition but get %d instead", len(dx.Items))
	}
}
