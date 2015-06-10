package tenants

import "encoding/xml"

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}

func createView(results []Tenant) ([]byte, error) {

	docRoot := &RootXML{}
	docRoot.Tenants = &results
	output, err := xml.MarshalIndent(docRoot, "", " ")
	return output, err
}
