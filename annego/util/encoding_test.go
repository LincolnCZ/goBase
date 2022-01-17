package util

import (
	"encoding/xml"
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestXMLStringMap(t *testing.T) {
	s := `<XMLStringMap><k1>abc</k1><k2>def</k2></XMLStringMap>`
	sm := map[string]string {
		"k1": "abc",
		"k2": "def",
	}

	m := &XMLStringMap{}
	if assert.NoError(t, xml.Unmarshal([]byte(s), m)) {
		assert.EqualValues(t, map[string]string(*m), sm)
	}

	_, err := xml.Marshal(m)
	assert.NoError(t, err) 
}

func TestEncodedJSONNode(t *testing.T) {
	type Test struct {
		S string `json:"s"`
		Elm EncodedJSONNode `json:"elm"`
	}
	s := `
	{
		"s": "test",
		"elm": ["abc","def"]
	}
	`

	test := &Test{}
	if assert.NoError(t, json.Unmarshal([]byte(s), test)) {
		assert.EqualValues(t, string(test.Elm.Data), `["abc","def"]`)
	}

	var err error
	var d []byte
	test.Elm.Data = nil
	d, err = json.Marshal(test)
	if assert.NoError(t, err) {
		assert.Equal(t, string(d), `{"s":"test","elm":null}`)
	}

	elm := 123
	test.Elm.Marshal(&elm)
	d, err = json.Marshal(test)
	if assert.NoError(t, err) {
		assert.Equal(t, string(d), `{"s":"test","elm":123}`)
	}
}
