package statusDetail

import "encoding/xml"

func createView(results []StatusDetailOutput) ([]byte, error) {

	docRoot := &ReadRoot{}

	if len(results) == 0 {
		output, err := xml.MarshalIndent(docRoot, " ", "  ")
		return output, err
	}

	profile := &Profile{}
	profile.Name = "ch.cern.sam.ROC_CRITICAL"

	vo := &Group{}
	vo.Type = "vo"
	vo.Name = "ops"

	//Removed loop call

	profile.Group = append(profile.Group, vo)

	docRoot.Profile = profile

	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err

}
