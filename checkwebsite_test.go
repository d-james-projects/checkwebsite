package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func assertResponse(t *testing.T, want bool, got bool, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("did not expect an error but got one %v", err)
	}

	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}

func makeServer(responseCode int, delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(responseCode)
	}))
}

func TestCheckWebsite(t *testing.T) {
	t.Run("test server response OK", func(t *testing.T) {
		myServer := makeServer(http.StatusOK, 50*time.Millisecond)
		defer myServer.Close()

		want := true
		got, err := checkWebsite(myServer.URL)
		assertResponse(t, want, got, err)
	})
	t.Run("test server response Service Unavailable", func(t *testing.T) {
		myServer := makeServer(http.StatusServiceUnavailable, 50*time.Millisecond)
		defer myServer.Close()

		want := false
		got, err := checkWebsite(myServer.URL)
		assertResponse(t, want, got, err)
	})
}

func mockCheckWebsite(URL string) (bool, error) {
	if URL == "thiswebsiteisdown" {
		return false, nil
	}
	return true, nil
}

func TestStartChecking(t *testing.T) {
	want := true
	got := startChecking(mockCheckWebsite, "websiteisup", time.Duration(1*time.Second))
	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
	got = startChecking(mockCheckWebsite, "thiswebsiteisdown", time.Duration(1*time.Second))
	if got != want {
		t.Errorf("got %t, want %t", got, want)
	}
}
