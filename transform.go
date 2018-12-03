package ago

import (
	"bytes"
	"encoding/xml"
	"html/template"
)

// Transform is ...
func Transform(gox string, state interface{}) Element {
	tmpl := template.Must(template.New("").Parse(gox))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, state)
	htmlStr := buf.String()
	element := createElementFromXML([]byte(htmlStr))
	return element
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node
	return d.DecodeElement((*node)(n), &start)
}

func createElementFromXML(data []byte) Element {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		panic(err)
	}
	return produceElement(n)
}

func produceElement(n Node) Element {
	childElements := []Element{}
	for _, child := range n.Nodes {
		childElements = append(childElements, produceElement(child))
	}
	if len(n.Nodes) == 0 && string(n.Content) != "" {
		childElements = CreateElementContent(string(n.Content))
	}
	props := make(map[string]interface{})
	if len(n.Attrs) > 0 {
		for _, att := range n.Attrs {
			props[att.Name.Local] = att.Value
		}
	}

	x := CreateElement(n.XMLName.Local, props, childElements)
	return x
}
