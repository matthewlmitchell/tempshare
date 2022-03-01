package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/gorilla/securecookie"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *application {

	templateCache, err := initTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New(securecookie.GenerateRandomKey(32))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteLaxMode

	// TODO: Add database support
	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		templateCache: templateCache,
	}
}

func newTestServer(t *testing.T, handler http.Handler, allowRedirect bool) *testServer {

	testServ := httptest.NewTLSServer(handler)

	// Add cookie support to the test server
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	testServ.Client().Jar = jar

	// If we don't want the client to be redirected, define the CheckRedirect
	// handler to simply return http.ErrUseLastResponse
	if !allowRedirect {
		testServ.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &testServer{testServ}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	response, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	return response.StatusCode, response.Header, responseBody
}

func (ts *testServer) post(t *testing.T, urlPath string, contentType string) (int, http.Header, []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "POST", urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	response, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	return response.StatusCode, response.Header, responseBody
}

func (ts *testServer) postForm(t *testing.T, urlPath string, data url.Values) (int, http.Header, []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postData := strings.NewReader(data.Encode())
	request, err := http.NewRequestWithContext(ctx, "POST", urlPath, postData)
	if err != nil {
		t.Fatal(err)
	}

	response, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	return response.StatusCode, response.Header, responseBody

}
