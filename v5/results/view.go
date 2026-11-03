package resultsV5

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/ARGOeu/argo-web-api/app/reports"
)

func createGroupResultView(results []EndpointGroupInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

	prevSuperGroup := ""
	prevEndpointGroup := ""
	endpointGroup := &Group{}
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.GetGroupType(),
			}
			docRoot.Result = append(docRoot.Result, superGroup)
			prevEndpointGroup = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevEndpointGroup != row.Name {
			prevEndpointGroup = row.Name
			endpointGroup = &Group{
				Name: row.Name,
				Type: report.GetEndpointGroupType(),
			}
			superGroup.Groups = append(superGroup.Groups, endpointGroup)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpointGroup.Results = append(endpointGroup.Results,
			&Result{
				Timestamp:    prepDate,
				Availability: RoundedFloat(row.Availability),
				Reliability:  RoundedFloat(row.Reliability),
				Unknown:      RoundedFloat(row.Unknown),
				Uptime:       RoundedFloat(row.Up),
				Downtime:     RoundedFloat(row.Down),
			})
	}

	return json.MarshalIndent(docRoot, " ", "  ")

}

func createSuperGroupView(results []SuperGroupInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

	prevSuperGroup := ""
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], row.Date)

		//if new superGroup does not match the previous superGroup value
		//we create a new superGroup entry in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.GetGroupType(),
			}
			docRoot.Result = append(docRoot.Result, superGroup)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		superGroup.Results = append(superGroup.Results,
			&ShortResult{
				Timestamp:    prepDate,
				Availability: RoundedFloat(row.Availability),
				Reliability:  RoundedFloat(row.Reliability)})
	}

	return json.MarshalIndent(docRoot, " ", "  ")

}

func createEndpointResultView(results []EndpointInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {
	docRoot := &root{}

	prevServiceFlavorGroup := ""
	prevServiceFlavor := ""
	prevEndpoint := ""
	serviceFlavorGroup := &ServiceFlavorGroup{}
	serviceEndpointGroup := &ServiceEndpointGroup{}
	endpoint := &Endpoint{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {

		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevServiceFlavorGroup != row.SuperGroup {
			prevServiceFlavorGroup = row.SuperGroup
			serviceFlavorGroup = &ServiceFlavorGroup{
				Name: row.SuperGroup,
				Type: report.GetEndpointGroupType(), // Endpoint groups are parents of SFs
			}
			docRoot.Result = append(docRoot.Result, serviceFlavorGroup)
			prevServiceFlavor = ""
		}
		//if new service flavor does not match the previous service value
		//we create a new service flavor entry in the xml/json output
		if prevServiceFlavor != row.Service {
			prevServiceFlavor = row.Service
			serviceEndpointGroup = &ServiceEndpointGroup{
				Name: row.Service,
				Type: "service",
			}
			serviceFlavorGroup.ServiceFlavor = append(serviceFlavorGroup.ServiceFlavor, serviceEndpointGroup)
			prevEndpoint = ""
		}

		if prevEndpoint != row.Name {
			prevEndpoint = row.Name
			endpoint = &Endpoint{
				Name: row.Name,
				Type: "endpoint",
				Info: row.Info,
			}
			serviceEndpointGroup.Endpoints = append(serviceEndpointGroup.Endpoints, endpoint)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		endpoint.Results = append(endpoint.Results,
			&Result{
				Timestamp:    prepDate,
				Availability: RoundedFloat(row.Availability),
				Reliability:  RoundedFloat(row.Reliability),
				Unknown:      RoundedFloat(row.Unknown),
				Uptime:       RoundedFloat(row.Up),
				Downtime:     RoundedFloat(row.Down),
			})
	}

	return json.MarshalIndent(docRoot, " ", "  ")

}

func createErrorMessage(message string, code int, format string) ([]byte, error) {

	var output []byte
	err := error(nil)
	docRoot := &errorMessage{}

	docRoot.Message = message
	docRoot.Code = code
	if strings.EqualFold(format, "application/json") {
		output, err = json.MarshalIndent(docRoot, " ", "  ")
	} else {
		output, err = xml.MarshalIndent(docRoot, " ", "  ")
	}
	return output, err
}

func createServiceFlavorResultView(results []ServiceFlavorInterface, report reports.MongoInterface, format string, custom bool) ([]byte, error) {

	docRoot := &root{}

	prevServiceFlavorGroup := ""
	prevServiceFlavor := ""
	serviceFlavor := &ServiceFlavor{}
	serviceFlavorGroup := &ServiceFlavorGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row
	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevServiceFlavorGroup != row.SuperGroup {
			prevServiceFlavorGroup = row.SuperGroup
			serviceFlavorGroup = &ServiceFlavorGroup{
				Name: row.SuperGroup,
				Type: report.GetEndpointGroupType(), // Endpoint groups are parents of SFs
			}
			docRoot.Result = append(docRoot.Result, serviceFlavorGroup)
			prevServiceFlavor = ""
		}
		//if new service flavor does not match the previous service value
		//we create a new service flavor entry in the xml/json output
		if prevServiceFlavor != row.Name {
			prevServiceFlavor = row.Name
			serviceFlavor = &ServiceFlavor{
				Name: row.Name,
				Type: "service",
			}
			serviceFlavorGroup.ServiceFlavor = append(serviceFlavorGroup.ServiceFlavor, serviceFlavor)
		}
		//we append the new availability values
		prepDate := timestamp.Format(customForm[1])
		if custom {
			prepDate = ""
		}
		serviceFlavor.Availability = append(serviceFlavor.Availability,
			&Result{
				Timestamp:    prepDate,
				Availability: RoundedFloat(row.Availability),
				Reliability:  RoundedFloat(row.Reliability),
				Unknown:      RoundedFloat(row.Unknown),
				Uptime:       RoundedFloat(row.Up),
				Downtime:     RoundedFloat(row.Down),
			})
	}

	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	}
	return xml.MarshalIndent(docRoot, " ", "  ")

}
