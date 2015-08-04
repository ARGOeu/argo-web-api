// ListMetricTimelines returns a list of metric timelines
import (
	"net/http"

	"github.com/ARGOeu/argo-web-api/utils/config"
)

func ListMetricTimelines(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "application/xml"
	charset := "utf-8"

	return code, h, output, err
}
