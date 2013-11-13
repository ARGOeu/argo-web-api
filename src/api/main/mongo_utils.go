package main

import (
	"encoding/xml"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func ErrorXML(s string) string {
	return "<root><error>" + s + "</error></root>"
}

type PoemProfile struct {
	Name           string "p"
	Namespace      string "ns"
	Group          string "g"
	Service_flavor string "sf"
}

func CreatePoemProfileNameXmlResponse(results []PoemProfile) ([]byte, error) {
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

func GetProfileNames(w http.ResponseWriter, r *http.Request) string {
	var results []PoemProfile
	session, err := mgo.Dial(cfg.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(cfg.MongoDB.Db).C("sites")
	err = c.Pipe([]bson.M{{"$group": bson.M{"_id": bson.M{"ns": "$ns", "p": "$p"}}}, {"$project": bson.M{"ns": "$_id.ns", "p": "$_id.p"}}}).All(&results)
	if err != nil {
		return ("<root><error>" + err.Error() + "</error></root>")
	}

	fmt.Println(results)
	output, err := CreatePoemProfileNameXmlResponse(results)
	return string(output)

}
