package util

import (
	"encoding/xml"
	"encoding/json"
	"io"
)

// XMLStringMap 可以把下列形式xml解析为 map[string]string
// <map>
//   <k1>v1</k1>
//   <k2>v2</k2>
// </map>
type XMLStringMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// MarshalXML 实现 xml.Marshaler
func (m XMLStringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

// UnmarshalXML 实现 xml.Unmarshaler
func (m *XMLStringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = XMLStringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

// EncodedJSONNode 可以保存任意的Json结构，只保存原始数据不进行解析
type EncodedJSONNode struct {
	Data []byte `json:"-"`
}

// MarshalJSON 实现 json.Marshaler
func (p *EncodedJSONNode) MarshalJSON() ([]byte, error) {
	if p.Data == nil {
		return []byte("null"), nil
	}
	return p.Data, nil
}

// UnmarshalJSON 实现 json.Unmarshaler
func (p *EncodedJSONNode) UnmarshalJSON(data []byte) error {
	if p.Data == nil {
		buffer := make([]byte, 0)
		p.Data = buffer
	}
	p.Data = append(p.Data, data...)
	return nil
}

func (p *EncodedJSONNode) Marshal(v interface{}) error {
	var err error
	p.Data, err = json.Marshal(v)
	return err
}

func (p *EncodedJSONNode) Unmarshal(v interface{}) error {
	return json.Unmarshal(p.Data, v)
}