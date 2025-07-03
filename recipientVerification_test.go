package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getUsersTest(a, b, c string) ([]User, error) {
	return []User{
		{
			Username: "foo",
		},
		{
			Username: "bar",
		},
		{
			Username: "baz",
		},
	}, nil
}

func Test_handleAdressenliste(t *testing.T) {
	expected := "#RZ-EV-BEGIN\nfoo\nbar\nbaz\n#RZ-EV-END\n"

	srv := httptest.NewServer(http.HandlerFunc(handleAdressenliste("", "", "", getUsersTest)))
	defer srv.Close()
	c, err := http.Get(srv.URL)
	defer c.Body.Close()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if c.StatusCode != http.StatusOK {
		t.Log("Status code:", c.StatusCode)
		t.FailNow()
	}
	out, err := io.ReadAll(c.Body)
	if string(out) != expected {
		t.Log("Body:", out)
		t.FailNow()
	}
}
