package cachen

import (
	"testing"
)

func TestReusableRequestIsCachable(t *testing.T) {
	cached := New()
	cached.ReusableRequest(true)
	if cached.cachable == noStore {
		t.Error("Reusable request cannot have a nostore")
	}
}

// func TestRevalidateEachTimeNotSetOnNonReusable(t *testing.T) {
// 	cached := New()
// 	cached.ReusableRequest(false).RevalidateEachTime(true)
// 	if cached.cachable == noCache || cached.cachable == "" {
// 		t.Error("Non Reusable must be " + noStore)
// 	}
// }

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
