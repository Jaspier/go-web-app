package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have any required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(postedData) // r.PostForm and postedData are of the same type

	if !form.Has("a") {
		t.Error("Shows form does not have field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	postedData := url.Values{}
	postedData.Add("invalid", "a")
	if form.MinLength("invalid", 3) {
		t.Error("Shows form field value has passed the minlength when it hasn't")
	}

	isError := form.Errors.Get("invalid")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("valid", "abc")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(postedData)

	if !form.MinLength("valid", 3) {
		t.Error("Shows form field value hasn't passed the minlength when it has")
	}

	isError = form.Errors.Get("valid")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	postedData = url.Values{}
	postedData.Add("invalid", "a@a")
	if form.IsEmail("invalid") {
		t.Error("Shows form field has valid email when it doesn't")
	}

	postedData = url.Values{}
	postedData.Add("valid", "a@a.com")
	form = New(postedData)

	if !form.IsEmail("valid") {
		t.Error("Shows form field doesn't have valid email when it does")
	}

}