package tenants

import "encoding/xml"

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}

func createView(results []Tenant) ([]byte, error) {

	docRoot := &Root{}

	for _, row := range results {
		r := &TenantXML{}
		r.Name = row.Name
		// Iterate over db configurations
		for _, d := range row.DbConf {
			dXML := &TenantDBConfXML{
				Store:    d.Store,
				Server:   d.Server,
				Port:     d.Port,
				Username: d.Username,
				Password: d.Password,
			}
			r.DbConf = append(r.DbConf, dXML)
		}
		// Iterate over users
		for _, u := range row.Users {
			uXML := &TenantUserXML{
				Name:   u.Name,
				Email:  u.Email,
				APIkey: u.APIkey,
			}
			r.Users = append(r.Users, uXML)
		}

		docRoot.Tenants = append(docRoot.Tenants, r)
	}
	output, err := xml.MarshalIndent(docRoot, "", " ")
	return output, err
}
