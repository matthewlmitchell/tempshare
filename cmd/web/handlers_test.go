package main

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	statusCode, _, responseBody := testServ.get(t, "/")

	if statusCode != http.StatusOK {
		t.Errorf("Expected status %d, received %d", http.StatusOK, statusCode)
	}

	strResponse := string(responseBody)

	expectedValue := "Home - TempShare"
	if !strings.Contains(strResponse, expectedValue) {
		t.Errorf("Expected %s in response body, received %s", expectedValue, strResponse)
	}
}

func TestCreateTempShareForm(t *testing.T) {
	app := newTestApplication(t)

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	statusCode, _, responseBody := testServ.get(t, "/create")

	if statusCode != http.StatusOK {
		t.Errorf("Expected status %d, received %d", http.StatusOK, statusCode)
	}

	strResponse := string(responseBody)

	expectedValue := "Generate link"
	if !strings.Contains(strResponse, expectedValue) {
		t.Errorf("Expected %s in response body, received %s", expectedValue, strResponse)
	}

}

func TestCreateTempShare(t *testing.T) {
	app := newTestApplication(t)

	err := app.initializeClient("./../../tls/cert.pem")
	if err != nil {
		t.Fatal(err)
	}

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	_, _, responseBody := testServ.get(t, "/create")
	csrfToken := extractCSRFToken(t, responseBody)

	testCases := []struct {
		name                    string
		tokenCSRF               string
		inputTempShareText      string
		inputTempShareExpires   string
		inputTempShareViewLimit string
		expectedStatusCode      int
		expectedResponse        []byte
	}{
		{
			name:                    "Valid Submission",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("Your TempShare link: https://"),
		},
		{
			name:                    "Empty CSRF Token",
			tokenCSRF:               "",
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusForbidden,
			expectedResponse:        []byte("Forbidden - CSRF token not found in request"),
		},
		{
			name:                    "Invalid CSRF Token",
			tokenCSRF:               "INVALID",
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusForbidden,
			expectedResponse:        []byte("Forbidden - CSRF token invalid"),
		},
		{
			name:                    "Text too short",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "A",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field must contain more than 2 characters"),
		},
		{
			name:                    "Text too long",
			tokenCSRF:               csrfToken,
			inputTempShareText:      strings.Repeat("A", 1025),
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field must contain less than 1024 characters"),
		},
		{
			name:                    "No input",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field must not be blank"),
		},
		{
			name:                    "No expiry",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field must not be blank"),
		},
		{
			name:                    "No view limit",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field must not be blank"),
		},
		{
			name:                    "Invalid expiry",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "100",
			inputTempShareViewLimit: "1",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field is invalid."),
		},
		{
			name:                    "Invalid view limit",
			tokenCSRF:               csrfToken,
			inputTempShareText:      "Hello World",
			inputTempShareExpires:   "1",
			inputTempShareViewLimit: "100",
			expectedStatusCode:      http.StatusOK,
			expectedResponse:        []byte("This field is invalid."),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("gorilla.csrf.Token", testCase.tokenCSRF)
			form.Add("text", testCase.inputTempShareText)
			form.Add("expires", testCase.inputTempShareExpires)
			form.Add("viewlimit", testCase.inputTempShareViewLimit)
			form.Add("g-recaptcha-response", "this-value-doesnt-matter-for-test-servers")

			statusCode, _, responseBody := testServ.postForm(t, "/create", form)

			if statusCode != testCase.expectedStatusCode {
				t.Errorf("Expected status %d, received status %d", testCase.expectedStatusCode, statusCode)
			}

			if !bytes.Contains(responseBody, testCase.expectedResponse) {
				t.Errorf("Expected %s, received %s", string(testCase.expectedResponse), string(responseBody))
			}
		})
	}
}

func TestViewTempShareForm(t *testing.T) {
	app := newTestApplication(t)

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	statusCode, _, responseBody := testServ.get(t, "/view")

	if statusCode != http.StatusOK {
		t.Errorf("Expected status %d, received %d", http.StatusOK, statusCode)
	}

	strResponse := string(responseBody)

	expectedValue := "View - TempShare"
	if !strings.Contains(strResponse, expectedValue) {
		t.Errorf("Expected %s in response body, received %s", expectedValue, strResponse)
	}
}

func TestViewTempShare(t *testing.T) {
	app := newTestApplication(t)

	err := app.initializeClient("./../../tls/cert.pem")
	if err != nil {
		t.Fatal(err)
	}

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	// Order of operations:
	//  GET on /view -> grab CSRF Token from response body
	// -> construct url.Values with csrf token and other necessary form fields
	// -> POST to /view
	// Check response for tempshare text.

	_, _, responseBody := testServ.get(t, "/view")
	csrfToken := extractCSRFToken(t, responseBody)

	testCases := []struct {
		name               string
		tokenCSRF          string
		tokenTempShare     string
		expectedStatusCode int
		expectedResponse   []byte
	}{
		{
			name:               "Valid token",
			tokenCSRF:          csrfToken,
			tokenTempShare:     "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []byte("This link has 0 uses remaining."),
		},
		{
			name:               "Empty CSRF Token",
			tokenCSRF:          "",
			tokenTempShare:     "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   []byte("Forbidden - CSRF token not found in request"),
		},
		{
			name:               "Invalid CSRF Token",
			tokenCSRF:          "INVALID",
			tokenTempShare:     "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   []byte("Forbidden - CSRF token invalid"),
		},
		{
			name:               "TempShare token too short",
			tokenCSRF:          csrfToken,
			tokenTempShare:     "INVALID",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []byte("Invalid token"),
		},
		{
			name:               "TempShare token too long",
			tokenCSRF:          csrfToken,
			tokenTempShare:     strings.Repeat("A", 53),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []byte("Invalid token"),
		},
		{
			name:               "TempShare token is empty",
			tokenCSRF:          csrfToken,
			tokenTempShare:     "",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []byte("Invalid token"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("gorilla.csrf.Token", testCase.tokenCSRF)
			form.Add("token", testCase.tokenTempShare)
			form.Add("g-recaptcha-response", "this-value-doesnt-matter-with-test-key")

			statusCode, _, responseBody := testServ.postForm(t, "/view", form)

			if statusCode != testCase.expectedStatusCode {
				t.Errorf("Expected status %d, received status %d", testCase.expectedStatusCode, statusCode)
			}

			if !bytes.Contains(responseBody, testCase.expectedResponse) {
				t.Errorf("Expected body %s to contain %s", responseBody, testCase.expectedResponse)
			}
		})
	}
}

func TestAbout(t *testing.T) {
	app := newTestApplication(t)

	testServ := newTestServer(t, app.routes(), false)
	defer testServ.Close()

	statusCode, _, responseBody := testServ.get(t, "/about")

	if statusCode != http.StatusOK {
		t.Errorf("Expected status %d, received %d", http.StatusOK, statusCode)
	}

	strResponse := string(responseBody)

	expectedValue := "About - TempShare"
	if !strings.Contains(strResponse, expectedValue) {
		t.Errorf("Expected %s in response body, received %s", expectedValue, strResponse)
	}
}
