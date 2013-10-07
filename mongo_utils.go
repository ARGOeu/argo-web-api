package main

import (
	"encoding/xml"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type MongoProfile struct {
	Name      string "p"
	Namespace string "ns"
}

func CreateProfileNameXmlResponse(results []MongoProfile) ([]byte, error) {
	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []Profile
	}

	v := &Root{}

	for _, result := range results {
		v.Profile = append(v.Profile,
			Profile{
				Name:      result.Name,
				Namespace: result.Namespace})
	}

	output, err := xml.MarshalIndent(v, " ", "  ")
	return output, err

}

func GetProfileNames(w http.ResponseWriter, r *http.Request) string {
	var results []MongoProfile
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("AR").C("sites")
	err = c.Pipe([]bson.M{{"$group": bson.M{"_id": bson.M{"ns": "$ns", "p": "$p"}}}, {"$project": bson.M{"ns": "$_id.ns", "p": "$_id.p"}}}).All(&results)
	if err != nil {
		return ("<root><error>" + err.Error() + "</error></root>")
	}

	fmt.Println(results)
	output, err := CreateProfileNameXmlResponse(results)
	return string(output)

}
