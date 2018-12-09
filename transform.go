package ago

import (
	"bytes"
	"encoding/xml"
	"html/template"
)

// Transform is ...
func Transform(gox string, state interface{}, createRealDom bool) Element {
	tmpl := template.Must(template.New("").Parse(gox))
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, state)
	htmlStr := buf.String()
	element := createElementFromXML([]byte(htmlStr), createRealDom)
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

func nodeToElement(n Node) Element {
	childElements := []Element{}
	for _, child := range n.Nodes {
		childElements = append(childElements, nodeToElement(child))
	}
	props := make(map[string]interface{})
	if len(n.Attrs) > 0 {
		for _, att := range n.Attrs {
			props[att.Name.Local] = att.Value
		}
	}
	element := Element{
		DomType:    n.XMLName.Local,
		Props:      props,
		DomContent: string(n.Content),
		Children:   &childElements,
	}
	return element
}

func createElementFromXML(data []byte, createRealDom bool) Element {
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		panic(err)
	}
	element := nodeToElement(n)
	return CreateElementRecursive(element.DomType, element.DomType, element.Props, *element.Children, createRealDom)
}
