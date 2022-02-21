package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// RecaptchaResponse is a custom struct used for unmarshalling the json response
// from a POST request to the reCAPTCHA API endpoint.
// Note: Score and Action are values used in reCAPTCHA v3 and are not used here at the moment.
type RecaptchaResponse struct {
	Success     bool `json:"success"`
	Score       float32
	Action      string
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string
}

// VerifyRecaptcha submits a POST request to the google reCAPTCHA API endpoint containing
// the server-side reCAPTCHA secret key along with the client's "g-recaptcha-response" and
// the client's IP address. The json response is then unmarshalled into our RecaptchaResponse
// struct and we return whether true if the recaptcha challenge was successful, false otherwise.
func VerifyRecaptcha(client *http.Client, r *http.Request, gRecaptchaResponse string) (bool, error) {
	googleAPIEndpoint := "https://google.com/recaptcha/api/siteverify"

	requestData := url.Values{
		"secret":   {os.Getenv("TEMPSHARE_reCAPTCHA_SECRET")},
		"response": {gRecaptchaResponse},
		"remoteip": {r.RemoteAddr},
	}

	response, err := client.PostForm(googleAPIEndpoint, requestData)
	if err != nil {
		return false, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	parsedResponse := RecaptchaResponse{}
	err = json.Unmarshal(responseData, &parsedResponse)
	if err != nil {
		return false, err
	}

	return parsedResponse.validate(), nil
}

func (response *RecaptchaResponse) validate() bool {
	return response.Success
}
