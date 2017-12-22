package csrf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	expected := &Handler{
		success:       GetTestHandler(),
		failure:       http.HandlerFunc(defaultFailureHandler),
		baseCookie:    nil,
		exemptMethods: ExemptMethods,
	}

	actual := New(GetTestHandler())

	if cmp.Equal(expected, actual, cmp.AllowUnexported(Handler{})) {
		t.Error("Expected New() to return correct Handler but it did not")
	}
}

func TestExemptMethods(t *testing.T) {
	handler := New(GetTestHandler())

	SetDummyCookie(handler)

	for _, method := range ExemptMethods {
		r, err := http.NewRequest(method, "http://dummy.url", nil)

		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		expected := 200

		if w.Code != expected {
			t.Errorf("An exempt method: %s did not pass CSRF check", method)
		}

		w.Flush()
	}
}

func TestContextIsAvailable(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		token := GetReqToken(r)

		if token == "" {
			t.Error("Could not find token in request")
		}
	}

	handler := New(http.HandlerFunc(handlerFunc))
	SetDummyCookie(handler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "https://dummy.url", nil)

	if err != nil {
		t.Error(err.Error())
	}

	handler.ServeHTTP(w, r)
}
