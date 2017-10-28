package cachen

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//Helper
func findValue(dataset []string, key string) bool {
	for _, v := range dataset {
		if v == key {
			return true
		}
	}
	return false
}

func TestReusableRequestIsCachable(t *testing.T) {
	cached := New()
	data := cached.ReusableRequest(true).State()

	if data.cachable == noStore {
		t.Error("Reusable request cannot be " + noStore)
	}
}

func TestRevalidateEachTimeNotSetOnNonReusable(t *testing.T) {
	cached := New()
	data := cached.ReusableRequest(false).RevalidateEachTime(true).State()

	if data.cachable == noCache || data.cachable == "" {
		t.Error("Non Reusable must be " + noStore + " and it's " + data.cachable)
	}
}
func TestIntermediatesAllowedOrNot(t *testing.T) {
	cached := New()
	data := cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).State()

	if data.intermediate != public {
		t.Error("intermediate must be " + public + " and it's" + data.intermediate)
	}

	ndata := cached.ReusableRequest(true).RevalidateEachTime(true).IntermediatesAllowed(false).State()

	if ndata.intermediate != private {
		t.Error("intermediate must be " + private + " and it's" + data.intermediate)
	}
}
func TestMaxAge(t *testing.T) {
	cached := New()
	cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5 * SECONDS)

	data := cached.State()
	if data.maxAge != data.smaxAge || data.maxAge != 5*SECONDS {
		t.Error("Max age cannot be different than smaxage or is not well set")
	}
	cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5*SECONDS, 1*MINUTES)

	ndata := cached.State()
	if ndata.maxAge == ndata.smaxAge || ndata.maxAge != 5*SECONDS || ndata.smaxAge != 1*MINUTES {
		t.Error("Max age cannot be equal to smaxcache and must save values correctly")
	}
}

func TestStaleAllowed(t *testing.T) {
	cached := New()
	cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5 * SECONDS).StaleAllowed(true)

	ndata := cached.State()
	if !findValue(ndata.cacheControl, revalidate) {
		t.Error("Must appear inside cacheControl the property:" + revalidate)
	}
	if findValue(ndata.cacheControl, proxyRevalidate) {
		t.Error("Cannot appear inside cacheControl the property:" + proxyRevalidate)
	}

	cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5*SECONDS).StaleAllowed(false, true)

	ndata = cached.State()
	if findValue(ndata.cacheControl, revalidate) {
		t.Error("Cannot appear inside cacheControl the property:" + revalidate)
	}
	if !findValue(ndata.cacheControl, proxyRevalidate) {
		t.Error("Must appear inside cacheControl the property:" + proxyRevalidate)
	}
}

func TestEtag(t *testing.T) {
	cached := New()
	cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5 * SECONDS).StaleAllowed(true).Etag("rd2d")

	ndata := cached.State()
	if ndata.etag != "rd2d" {
		t.Error("Must have an etag correctly set")
	}
}

func TestHandlerReturn(t *testing.T) {

	cached := New()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(cached.ReusableRequest(true).
		RevalidateEachTime(true).
		IntermediatesAllowed(true).MaxAge(5 * SECONDS).StaleAllowed(false).Etag("rd2d").Bind)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	etag := rr.Header().Get("Etag")
	if etag != "rd2d" {
		t.Error("Etag is not returned as a header")
	}
	cc := rr.Header().Get("Cache-Control")
	ccs := strings.Split(cc, ",")
	if !findValue(ccs, "no-cache") || !findValue(ccs, "public") || !findValue(ccs, "maxage=5") || !findValue(ccs, "smaxage=5") {
		t.Errorf("One or more cache keys are missing: %v", cc)
	}
}
