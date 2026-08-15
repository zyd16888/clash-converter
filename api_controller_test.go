package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSubscriptionEndpointCustomNamesAndFilename(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/nodes.yaml":
			w.Header().Set("Content-Disposition", `attachment; filename="provider.yaml"`)
			w.Header().Set("Subscription-Userinfo", "upload=0; download=1073741824; total=10737418240; expire=1788220800")
			_, _ = w.Write([]byte("proxies:\n  - name: HK 01\n    type: ss\n    server: 127.0.0.1\n    port: 10001\n    cipher: aes-128-gcm\n    password: test\n"))
		case "/script.js":
			_, _ = w.Write([]byte("function rulesets(register) {}\nfunction buildConfig(config, legacyRelay) { config['proxy-groups'] = []; config['rules'] = config['rules'].concat(['MATCH,DIRECT']); }"))
		case "/template.yaml":
			_, _ = w.Write([]byte("mode: rule\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()

	previousToken := Token
	Token = "test-token"
	t.Cleanup(func() { Token = previousToken })

	query := url.Values{}
	query.Add("sub", upstream.URL+"/nodes.yaml")
	query.Add("subName", "Premium A")
	query.Set("script", upstream.URL+"/script.js")
	query.Set("template", upstream.URL+"/template.yaml")
	query.Set("token", Token)
	query.Set("filename", "My Combined.yml")

	request := httptest.NewRequest(http.MethodGet, "/sub?"+query.Encode(), nil)
	response := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Disposition"); got != "attachment; filename*=UTF-8''My%20Combined.yaml" {
		t.Fatalf("unexpected Content-Disposition: %q", got)
	}
	for _, fragment := range []string{"[Premium A] HK 01", "Premium A | 已用 1.0 GB / 10.0 GB", "MATCH,DIRECT"} {
		if !strings.Contains(response.Body.String(), fragment) {
			t.Errorf("response does not contain %q:\n%s", fragment, response.Body.String())
		}
	}

	for _, endpoint := range []string{"/example/script.js", "/example/template.yaml"} {
		exampleRequest := httptest.NewRequest(http.MethodGet, endpoint, nil)
		exampleResponse := httptest.NewRecorder()
		router.ServeHTTP(exampleResponse, exampleRequest)
		if exampleResponse.Code != http.StatusOK || exampleResponse.Body.Len() == 0 {
			t.Errorf("embedded example %s is unavailable", endpoint)
		}
	}
}
