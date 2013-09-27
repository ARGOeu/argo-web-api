package main

import (
	"encoding/xml"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
//	"strings"
	"time"
	"fmt"
)

type Site struct {
	SiteScope	string	"ss"
	Scope		string 	"sc"
	Date 		int	"dt"
	Namespace	string	"ns"	
	Profile		string	"p"
	Production	string	"pr"
	Monitored	string 	"m"
	Ngi		string 	"n"
	Site 		string 	"s"
	Infastructure 	string 	"i"
	CertStatus	string 	"cs"
	Availability 	float64 "a"
	Reliability 	float64 "r"
}


func createSiteXMLResponse(results []Site) ([]byte, error) {
		
	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
		}
	
	type Site struct {
		Site		string `xml:"site,attr"`
		Ngi		string `xml:"NGI,attr"`
		Infastructure	string `xml:"infastructure,attr"`
		Scope		string `xml:"scope,attr"`
		SiteScope	string `xml:"site_scope,attr"`
		Production	string `xml:"production,attr"`
		Monitored	string `xml:"monitored,attr"`
		CertStatus	string `xml:"certification_status"`
		Availability   []Availability
	}
	
	type Profile struct {
		XMLName   xml.Name 	`xml:"Profile"`
		Name      string 	`xml:"name,attr"`
		Namespace string	`xml:"namespace,attr"`
		Sites   []Site
	}
	
		type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []Profile
	}
	
	
	v := &Root{}
	
	//v.Profile = make([]Profile,len(results))
	v.Profile = make([]Profile,0)
	site := Site{}	

	//total := len(results)
	
	prevProfile := ""
	prevSite := ""
	for cur, result := range results {
		timestamp, _ := time.Parse(ymdForm, strconv.Itoa(result.Date))
		
		if prevProfile != result.Profile {
			 
			prevProfile = result.Profile
			v.Profile = append(v.Profile,
					Profile {
						Name: result.Profile,
						Namespace: result.Namespace })
		}
			
		if prevSite != result.Site {
			prevSite = result.Site
			if cur>0 {
			v.Profile[len(v.Profile)-1].Sites = append(
				v.Profile[len(v.Profile)-1].Sites, site)
			}
			site = Site{
				Site: 		result.Site,
				Ngi: 		result.Ngi,
				Infastructure: 	result.Infastructure,
				Scope: 		result.Scope,
				SiteScope: 	result.SiteScope,
				Production: 	result.Production,
				Monitored: 	result.Monitored,
				CertStatus: 	result.CertStatus }
		}
		site.Availability = append(site.Availability,
					Availability{
						Timestamp: 	timestamp.Format(zuluForm),
						Availability: 	fmt.Sprintf("%g",result.Availability),
						Reliability: 	fmt.Sprintf("%g",result.Reliability) })
	}
	if (len(v.Profile)>0){
	v.Profile[len(v.Profile)-1].Sites = append(v.Profile[len(v.Profile)-1].Sites, site)
	}
	//v.Profile = v.Profile[1:len(v.Profile)-1]

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}

//const zuluForm = "2006-01-02T15:04:05Z"
//const ymdForm = "20060102"

func SitesAvailabilityInProfile(w http.ResponseWriter, r *http.Request) string {

	// This is the input we will receive from the API

	type ApiServiceAvailabilityInProfileInput struct {
		// mandatory values
		start_time          string   // UTC time in W3C format
		end_time            string   // UTC time in W3C format
		profile_name        []string // may appear more than once. (eg: CMS_CRITICAL)
		group_type          []string // may appear more than once. (eg: CMS_Site)
		availability_period string   // availability period; possible values: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`
		// optional values
		output           string   // default XML; possible values are: XML, JSON
		namespace        []string // profile namespace; may appear more than once. (eg: ch.cern.sam)
		group_name       []string // site name; may appear more than once
	}

	urlValues := r.URL.Query()

	input := ApiServiceAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
	}

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("AR").C("sites")
	results := []Site{}
	q := bson.M{
		"dt":  bson.M{"$gte": tsYMD, "$lte": teYMD},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		q["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		// TODO: We do not have the site name in the timeline
	}


	err = c.Find(q).Sort("p", "n", "s", "dt").All(&results)
	fmt.Println(q)
	fmt.Println(len(results))
	if err != nil {
		return ("<root><error>" + err.Error() + "</error></root>")
	}

	output, err := createSiteXMLResponse(results)

	return string(output)
}
