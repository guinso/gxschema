package gxschema

import (
	"strings"
	"testing"
)

func TestParseSchemaFromXML_expectKeywordIDFail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="id" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured, keyword 'id' shall not be accepted by schema parser")
	}
}

func TestParseSchemaFromXML_expectKeywordParentIDFail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="parent_id" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured, keyword 'id' shall not be accepted by schema parser")
	}
}

func TestParseSchemaFromXML_expectKeywordFilenameFail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="filename" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured, keyword 'id' shall not be accepted by schema parser")
	}
}

func TestParseSchemaFromXML_expectKeywordFilepathFail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="filepath" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured, keyword 'id' shall not be accepted by schema parser")
	}
}

func TestParseSchemaFromXML_expectNamePatternFail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="123Koko" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured due to invalid name pattern '123koko'")
	}
}

func TestParseSchemaFromXML_expectNamePattern2Fail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="Koko Jelly" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured due to invalid name pattern 'koko jelly'")
	}
}

func TestParseSchemaFromXML_expectNamePattern3Fail(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
		<dxint name="age"></dxint>
		<dxdecimal name="koko#$jelly" precision="2"></dxdecimal>
		<dxstr name="createdBy" isOptional="false"></dxstr>
		<dxbool name="isUrgent" isArray="true"></dxbool>
		<dxstr name="docNo" lenLimit="6"></dxstr>
		<dxsection name="items" isArray="true">
			<dxstr name="description"></dxstr>
			<dxint name="quantity"></dxint>
			<dxdecimal name="price" precision="2"></dxdecimal>
		</dxsection>
	</dxdoc>`

	_, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr == nil {
		t.Error("Expect error occured due to invalid name pattern 'koko#$jelly'")
	}
}

func TestWalkDxDoc(t *testing.T) {
	rawXML := `<?xml version="1.0"?>
	<dxdoc name="invoice" revision="3" id="733bee1b-f79a-4cb7-b675-842317b994b5">
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

	dx, dxErr := ParseSchemaFromXML(rawXML)
	if dxErr != nil {
		t.Error(dxErr)
		return
	}

	if dx.Revision != 3 {
		t.Errorf("expect revision value is 3 but get %d instead", dx.Revision)
	}

	if strings.Compare(dx.ID, "733bee1b-f79a-4cb7-b675-842317b994b5") != 0 {
		t.Errorf("expect ID value is 733bee1b-f79a-4cb7-b675-842317b994b5 but get %s instead", dx.ID)
	}

	if strings.Compare(dx.Name, "invoice") != 0 {
		t.Errorf("expect name value is 'invoice' but get '%s' instead", dx.Name)
	}

	if len(dx.Items) != 6 {
		t.Errorf("Expect has 6 items definition but get %d instead", len(dx.Items))
	}
}
