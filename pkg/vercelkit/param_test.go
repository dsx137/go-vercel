package vercelkit_test

import (
	"net/url"
	"testing"

	"github.com/dsx137/go-vercel/pkg/vercelkit"
)

type TestParams struct {
	Name    string `query:"name"`
	Age     int    `query:"age"`
	Numbers []int  `query:"numbers"`
}

func TestParam(t *testing.T) {
	query := url.Values{}
	query.Set("name", "John")
	query.Set("age", "30")
	if _, err := vercelkit.ReadParamsFromQuery[TestParams](query); err == nil {
		t.Errorf("expected error, got nil")
	}

	query.Set("numbers", "1,2,3")

	params, err := vercelkit.ReadParamsFromQuery[TestParams](query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if params.Name != "John" {
		t.Errorf("expected name to be 'John', got '%s'", params.Name)
	}
	if params.Age != 30 {
		t.Errorf("expected age to be 30, got %d", params.Age)
	}
	if len(params.Numbers) != 3 || params.Numbers[0] != 1 || params.Numbers[1] != 2 || params.Numbers[2] != 3 {
		t.Errorf("expected numbers to be [1, 2, 3], got %v", params.Numbers)
	}
}
