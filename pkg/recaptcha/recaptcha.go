package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type RecaptchaResponse struct {
	Success     bool
	Score       float32
	Action      string
	ChallengeTS time.Time
	Hostname    string
	ErrorCodes  []string
}

func VerifyRecaptcha(r *http.Request, gRecaptchaResponse string) (bool, error) {
	googleAPIEndpoint := "https://google.com/recaptcha/api/siteverify"

	requestData := url.Values{
		"secret":   {os.Getenv("TEMPSHARE_reCAPTCHA_SECRET")},
		"response": {gRecaptchaResponse},
		"remoteip": {r.RemoteAddr},
	}

	response, err := http.PostForm(googleAPIEndpoint, requestData)
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
