package cachen

import (
	"testing"
)

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

// func TestHandlerReturn(t *testing.T) {

// 	cached := New()

// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(cached.Handler)

// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// // Check the response body is what we expect.
// 	// expected := `{"alive": true}`
// 	// if rr.Body.String() != expected {
// 	// 	t.Errorf("handler returned unexpected body: got %v want %v",
// 	// 		rr.Body.String(), expected)
// 	// }
// }
