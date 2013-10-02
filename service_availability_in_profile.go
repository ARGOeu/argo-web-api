package main

import (
	"encoding/xml"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
//	"fmt"
)

type Timeline struct {
	Profile       string "p"
	ServiceFlavor string "sf"
	Host          string "h"
	Timeline      string "tm"
	VO            string "vo"
	Date          int    "d"
	Namespace     string "ns"
}

func createXMLResponse(results []Timeline) ([]byte, error) {
	type Availability struct {
		XMLName      xml.Name `xml:"Availability"`
		Timestamp    string   `xml:"timestamp,attr"`
		Availability string   `xml:"availability,attr"`
		Reliability  string   `xml:"reliability,attr"`
		Maintenance  string   `xml:"maintenance,attr"`
	}

	type Service struct {
		Hostname       string `xml:"hostname,attr"`
		Service_Type   string `xml:"type,attr"`
		Service_Flavor string `xml:"flavor,attr"`
		Availability   []Availability
	}

	type Profile struct {
		XMLName   xml.Name `xml:"Profile"`
		Name      string   `xml:"name,attr"`
		Namespace string   `xml:"namespace,attr"`
		VO        string   `xml:"defined_by_vo_name,attr"`
		Service   []Service
	}

	type Root struct {
		XMLName xml.Name `xml:"root"`
		Profile []Profile
	}

	v := &Root{}
	v.Profile = make([]Profile, 0)

	service := Service{}

	prevProfile := ""
	prevService := ""

	for cur, result := range results {
		timestamp, _ := time.Parse(ymdForm, strconv.Itoa(result.Date))
		timeline := strings.Split(strings.Trim(result.Timeline, "[\n]"), ", ")
		if prevProfile != result.Profile {
			prevProfile = result.Profile
			v.Profile = append(v.Profile,
				Profile{
					Name:      result.Profile,
					Namespace: result.Namespace,
					VO:        result.VO})
		}
		if prevService != result.Host+result.ServiceFlavor {
			prevService = result.Host + result.ServiceFlavor
			if cur > 0 {
				v.Profile[len(v.Profile)-1].Service = append(v.Profile[len(v.Profile)-1].Service, service)
			}
			service = Service{
				Hostname:       result.Host,
				Service_Type:   result.ServiceFlavor,
				Service_Flavor: result.ServiceFlavor}

		}

		for _, timeslot := range timeline {
			ar := strings.Split(timeslot, ":")
			if len(ar) != 3 {
				return []byte("<root><error>500: Internal server error (Malformed timeslot)</error></root>"), nil
			}

			service.Availability = append(service.Availability,
				Availability{
					Timestamp:    timestamp.Format(zuluForm),
					Availability: ar[0],
					Reliability:  ar[1],
					Maintenance:  ar[2]})
			timestamp = timestamp.Add(time.Duration(60*60) * time.Second)
		}

	}
	if (len(v.Profile)>0) {
	v.Profile[len(v.Profile)-1].Service = append(v.Profile[len(v.Profile)-1].Service, service)
	}

	output, err := xml.MarshalIndent(v, " ", "  ")

	return output, err
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func ServiceAvailabilityInProfile(w http.ResponseWriter, r *http.Request) string {

	// This is the input we will receive from the API

	type ApiServiceAvailabilityInProfileInput struct {
		// mandatory values
		start_time          string   // UTC time in W3C format
		end_time            string   // UTC time in W3C format
		vo_name             []string // may appear more than once. (eg: ops)
		profile_name        []string // may appear more than once. (eg: CMS_CRITICAL)
		group_type          []string // may appear more than once. (eg: CMS_Site)
		availability_period string   // availability period; possible values: 'HOURLY', 'DAILY', 'WEEKLY', 'MONTHLY'
		// optional values
		output           string   // default XML; possible values are: XML, JSON
		namespace        []string // profile namespace; may appear more than once. (eg: ch.cern.sam)
		group_name       []string // site name; may appear more than once
		service_flavour  []string // service flavour name; may appear more than once. (eg: SRMv2)
		service_hostname []string // service hostname; may appear more than once. (eg: ce202.cern.ch)
	}

	urlValues := r.URL.Query()

	input := ApiServiceAvailabilityInProfileInput{
		urlValues.Get("start_time"),
		urlValues.Get("end_time"),
		urlValues["vo_name"],
		urlValues["profile_name"],
		urlValues["group_type"],
		urlValues.Get("type"),
		urlValues.Get("output"),
		urlValues["namespace"],
		urlValues["group_name"],
		urlValues["service_flavour"],
		urlValues["service_hostname"],
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
	c := session.DB("AR").C("timelines")
	results := []Timeline{}
	q := bson.M{
		"d":  bson.M{"$gte": tsYMD, "$lte": teYMD},
		"vo": bson.M{"$in": input.vo_name},
		"p":  bson.M{"$in": input.profile_name},
	}

	if len(input.namespace) > 0 {
		q["ns"] = bson.M{"$in": input.namespace}
	}

	if len(input.group_name) > 0 {
		// TODO: We do not have the site name in the timeline
	}

	if len(input.service_flavour) > 0 {
		q["ns"] = bson.M{"$in": input.service_flavour}
	}

	if len(input.service_hostname) > 0 {
		q["h"] = bson.M{"$in": input.service_hostname}
	}
	query := []bson.M{{"$match" : q}, {"$sort" : bson.D{{"p" ,1},{"h",1},{"sf",1}}}}
	err = c.Pipe(query).All(&results)
	
	//err = c.Find(q).Sort("p", "h", "sf").All(&results)
	if err != nil {
		return ("<root><error>" + err.Error() + "</error></root>")
	}
	
//	fmt.Println(results)
	output, err := createXMLResponse(results)

	return string(output)
}
