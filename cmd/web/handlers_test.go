package main

import (
	"net/http"
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

// TODO: Add tests for form submission and spinning up a mock MySQL DB
