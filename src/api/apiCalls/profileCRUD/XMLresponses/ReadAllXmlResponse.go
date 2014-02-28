package XMLresponses

import (
	"encoding/xml"
)

type PoemProfile struct {
	Name           string "p"
	Namespace      string "ns"
	Group          string "g"
	Service_flavor string "sf"
}

func ReadAllXmlResponse(results []PoemProfile) ([]byte, error) {
	type Profile struct {
		XMLName        xml.Name `xml:"Profile"`
		Name           string   `xml:"name,attr"`
		Namespace      string   `xml:"namespace,attr"`
		Group          string   `xml:"group,attr"`
		Service_flavor string   `xml:"service_flavor,attr"`
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []Profile
	}

	v := &Root{}

	for _, result := range results {
		v.Profile = append(v.Profile,
			Profile{
				Name:           result.Name,
				Namespace:      result.Namespace,
				Group:          result.Group,
				Service_flavor: result.Service_flavor,
			})
	}
	output, err := xml.MarshalIndent(v, " ", "  ")
	return output, err

}