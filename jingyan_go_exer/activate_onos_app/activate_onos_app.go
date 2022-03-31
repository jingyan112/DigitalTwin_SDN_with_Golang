package activate_onos_app

import (
	"log"
	"net/http"
)

//
// Function name: Activate_onos_app
//
// Function description: activate the applications of ONOS container based on the Paramters through RESTAPI provided by ONOS,
//                       and return the status code of HTTP POST request,
//                       eg: curl -X POST -u onos:rocks http://127.0.0.1:8181/onos/v1/applications/org.onosproject.fwd/active
//
// Paramters: app_name (string), used to activate the applications of ONOS container
//
// Returns: resp.StatusCode, status code of HTTP POST request
//
func Activate_onos_app(app_name string) int {
	url := "http://127.0.0.1:8181/onos/v1/applications/" + app_name + "/active"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Print("HTTP Request establishment failed as the following reason: ")
		panic(err)
	}
	req.SetBasicAuth("onos", "rocks")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("Sending HTTP Request failed as the following reason: ")
		panic(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
