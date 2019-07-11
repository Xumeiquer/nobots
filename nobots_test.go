package nobots

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

var t1 = `nobots "nobots.go" { "Googlebot" }`

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("http", t1)
	err := setup(c)

	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}

	mids := httpserver.GetConfig(c).Middleware()
	if len(mids) == 0 {
		t.Fatal("Expected middleware, got 0 instead")
	}

	handler := mids[0](httpserver.EmptyNext)
	myHandler, ok := handler.(BotUA)
	if !ok {
		t.Fatalf("Expected handler to be type BotUA, got: %#v", handler)
	}

	if !httpserver.SameNext(myHandler.Next, httpserver.EmptyNext) {
		t.Error("'Next' field of handler was not set properly")
	}
	tests := []struct {
		input     string
		shouldErr bool
	}{
		// Bomb exists so plugin initiates correctly
		{`nobots "nobots.go" { "Googlebot" }`, false},
		// Bomb exists and regexp keyword is valid
		{`nobots "nobots.go" { regexp "Googlebot" }`, false},
		/* Bomb exists and regexp keyword is not valid even though
		   Nobots take it as UA */
		{`nobots "nobots.go" { regex "Googlebot" }`, false},
		// Bomb exists and regexp valid
		{`nobots "nobots.go" { regexp "^Googlebot$" }`, false},
		// Bomb exists and regexp not valid
		{`nobots "nobots.go" { regexp "(?P<name>re" }`, true},
	}

	for i, test := range tests {
		_, err := parseUA(caddy.NewTestController("http", test.input))
		if test.shouldErr && err == nil {
			t.Errorf("Test %v: Expected error but found nil", i)
		} else if !test.shouldErr && err != nil {
			t.Errorf("Test %v: Expected no error but found error: %v", i, err)
		}
	}
}

var t2 = `nobots "no_exist.go" { "Googlebot" }`

// Bomb does not exist so the plugin must throw an error
func TestSetup1(t *testing.T) {
	c := caddy.NewTestController("http", t2)
	err := setup(c)

	if err == nil {
		t.Errorf("Expected error but found nil")
	}
}

func TestNobotsWithPublic(t *testing.T) {
	funcName := "TestNobotsWithPublic"
	myHandler := func(w http.ResponseWriter, r *http.Request) (int, error) {

		return http.StatusOK, nil
	}

	filename := "nobots.go"

	rws := []BotUA{
		{
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    []string{"Bot"},
				re:     nil,
				public: []*regexp.Regexp{regexp.MustCompile("/public")},
			},
		}, {
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    nil,
				re:     []*regexp.Regexp{regexp.MustCompile("^Bot")},
				public: []*regexp.Regexp{regexp.MustCompile("/public")},
			},
		}, {
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    nil,
				re:     nil,
				public: nil,
			},
		},
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {

	}

	fileSize := strconv.Itoa(len(file))

	type headerType struct {
		type_    string
		encoding string
		length   string
	}

	type testType struct {
		path   string
		ua     string
		result int
		header headerType
	}

	tests := []testType{
		{
			path:   "/private",
			ua:     "Bot",
			result: http.StatusOK,
			header: headerType{
				type_:    "text/html; charset=UTF-8",
				encoding: "gzip",
				length:   fileSize,
			},
		}, {
			path:   "/this/is/public",
			ua:     "Bot",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/public",
			ua:     "Bot",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/private",
			ua:     "Got",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/public",
			ua:     "Got",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/private",
			ua:     "",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		},
	}

	for i, rw := range rws {
		for j, test := range tests {
			req, err := http.NewRequest("GET", test.path, nil)
			if err != nil {
				t.Fatalf("Test %d: Could not create HTTP request: %v", j, err)
			}

			req.Header.Set("User-Agent", test.ua)
			rec := httptest.NewRecorder()
			result, err := rw.ServeHTTP(rec, req)

			if err != nil {
				t.Fatalf("Test %d: Could not ServeHTTP: %v", j, err)
			}

			if result != test.result {
				t.Errorf("Test %d: Expected status code %d but was %d",
					j, test.result, result)
			}

			if len(rw.UA.uas) > 0 || len(rw.UA.re) > 0 {
				if rec.HeaderMap.Get("Content-Type") != test.header.type_ {
					t.Errorf("Test %d-%d (%s): Expected Content-Type '%s' but found '%s'",
						i, j, funcName, test.header.type_, rec.HeaderMap.Get("Content-Type"))
				}
				if rec.HeaderMap.Get("Content-Encoding") != test.header.encoding {
					t.Errorf("Test %d-%d (%s): Expected Content-Encoding '%s' but found '%s'",
						i, j, funcName, test.header.encoding, rec.HeaderMap.Get("Content-Encoding"))
				}
				if rec.HeaderMap.Get("Content-Length") != test.header.length {
					t.Errorf("Test %d-%d (%s): Expected Content-Length '%s' but found '%s'",
						i, j, funcName, test.header.length, rec.HeaderMap.Get("Content-Length"))
				}
			}
		}
	}

}

func TestNobots(t *testing.T) {
	funcName := "TestNobots"
	myHandler := func(w http.ResponseWriter, r *http.Request) (int, error) {

		return http.StatusOK, nil
	}

	filename := "nobots.go"

	rws := []BotUA{
		{
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    []string{"Bot"},
				re:     nil,
				public: nil,
			},
		}, {
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    nil,
				re:     []*regexp.Regexp{regexp.MustCompile("^Bot")},
				public: nil,
			},
		}, {
			Next: httpserver.HandlerFunc(myHandler),
			UA: &botUA{
				bomb:   filename,
				uas:    nil,
				re:     nil,
				public: nil,
			},
		},
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {

	}

	fileSize := strconv.Itoa(len(file))

	type headerType struct {
		type_    string
		encoding string
		length   string
	}

	type testType struct {
		path   string
		ua     string
		result int
		header headerType
	}

	tests := []testType{
		{
			path:   "/private",
			ua:     "Bot",
			result: http.StatusOK,
			header: headerType{
				type_:    "text/html; charset=UTF-8",
				encoding: "gzip",
				length:   fileSize,
			},
		}, {
			path:   "/public",
			ua:     "Bot",
			result: http.StatusOK,
			header: headerType{
				type_:    "text/html; charset=UTF-8",
				encoding: "gzip",
				length:   fileSize,
			},
		}, {
			path:   "/private",
			ua:     "Got",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/public",
			ua:     "Got",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		}, {
			path:   "/private",
			ua:     "",
			result: http.StatusOK,
			header: headerType{
				type_:    "",
				encoding: "",
				length:   "",
			},
		},
	}

	for i, rw := range rws {
		for j, test := range tests {
			req, err := http.NewRequest("GET", test.path, nil)
			if err != nil {
				t.Fatalf("Test %d: Could not create HTTP request: %v", j, err)
			}

			req.Header.Set("User-Agent", test.ua)
			rec := httptest.NewRecorder()
			result, err := rw.ServeHTTP(rec, req)

			if err != nil {
				t.Fatalf("Test %d: Could not ServeHTTP: %v", j, err)
			}

			if result != test.result {
				t.Errorf("Test %d: Expected status code %d but was %d",
					j, test.result, result)
			}

			if len(rw.UA.uas) > 0 || len(rw.UA.re) > 0 {
				if rec.HeaderMap.Get("Content-Type") != test.header.type_ {
					t.Errorf("Test %d-%d (%s): Expected Content-Type '%s' but found '%s'",
						i, j, funcName, test.header.type_, rec.HeaderMap.Get("Content-Type"))
				}
				if rec.HeaderMap.Get("Content-Encoding") != test.header.encoding {
					t.Errorf("Test %d-%d (%s): Expected Content-Encoding '%s' but found '%s'",
						i, j, funcName, test.header.encoding, rec.HeaderMap.Get("Content-Encoding"))
				}
				if rec.HeaderMap.Get("Content-Length") != test.header.length {
					t.Errorf("Test %d-%d (%s): Expected Content-Length '%s' but found '%s'",
						i, j, funcName, test.header.length, rec.HeaderMap.Get("Content-Length"))
				}
			}
		}
	}

}
