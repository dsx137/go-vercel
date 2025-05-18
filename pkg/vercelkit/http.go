package vercelkit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var HttpClient = &http.Client{}

func HttpTest(t *testing.T, method string, handler VercelHandler, params url.Values) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, "https://test.com?"+params.Encode(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Error("Status code: " + fmt.Sprintf("%d", w.Code) + ", Body: " + w.Body.String())
		return
	}
}

func HttpRequest(req *http.Request) ([]byte, error) {
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func HttpResponse(w http.ResponseWriter, status int, message any) error {
	h := w.Header()

	h.Del("Content-Length")
	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(status)
	jsonResponse, err := json.Marshal(map[string]any{"message": message})
	if err != nil {
		return err
	}
	w.Write([]byte(jsonResponse))
	return nil
}
