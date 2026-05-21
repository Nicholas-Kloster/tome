package fingerprint

import (
	"testing"

	"github.com/Nicholas-Kloster/tome/internal/corpus"
)

var weaviatePlatform = corpus.Platform{
	Platform:     "weaviate",
	DefaultPorts: []int{8080, 50051},
	AuthDefault:  "none",
	Fingerprint: corpus.Fingerprint{
		Passive: []string{
			"product:Weaviate",
			"http.html:\"/v1/graphql\"",
			"port:8080",
		},
	},
}

func TestMatchPassiveFullMatch(t *testing.T) {
	host := ShodanHost{
		IPStr: "1.2.3.4",
		Ports: []int{8080, 50051},
		Data: []ShodanData{
			{
				Port:    8080,
				Product: "Weaviate",
				HTTP: ShodanHTTP{
					HTML:  "visit /v1/graphql for the explorer",
					Title: "Weaviate",
				},
			},
		},
	}
	conf := MatchPassive(weaviatePlatform, host)
	if conf != 1.0 {
		t.Errorf("full match confidence = %.2f, want 1.0", conf)
	}
}

func TestMatchPassivePartialMatch(t *testing.T) {
	host := ShodanHost{
		IPStr: "1.2.3.4",
		Ports: []int{8080},
		Data: []ShodanData{
			{Port: 8080, Product: "Weaviate"},
		},
	}
	conf := MatchPassive(weaviatePlatform, host)
	// product:Weaviate and port:8080 match; http.html does not
	if conf < 0.5 || conf >= 1.0 {
		t.Errorf("partial match confidence = %.2f, want ~0.67", conf)
	}
}

func TestMatchPassiveNoMatch(t *testing.T) {
	host := ShodanHost{
		IPStr: "1.2.3.4",
		Ports: []int{3000},
		Data:  []ShodanData{{Port: 3000, Product: "nginx"}},
	}
	conf := MatchPassive(weaviatePlatform, host)
	if conf != 0.0 {
		t.Errorf("no match confidence = %.2f, want 0.0", conf)
	}
}

func TestMatchFilterPort(t *testing.T) {
	host := ShodanHost{Ports: []int{8080, 443}}
	if !matchFilter("port:8080", host) {
		t.Error("port:8080 should match host with port 8080")
	}
	if matchFilter("port:9999", host) {
		t.Error("port:9999 should not match host without that port")
	}
}

func TestMatchFilterProduct(t *testing.T) {
	host := ShodanHost{Data: []ShodanData{{Product: "Weaviate"}}}
	if !matchFilter("product:Weaviate", host) {
		t.Error("product:Weaviate should match")
	}
	if matchFilter("product:ChromaDB", host) {
		t.Error("product:ChromaDB should not match Weaviate host")
	}
}

func TestMatchFilterHTMLCaseInsensitive(t *testing.T) {
	host := ShodanHost{Data: []ShodanData{{HTTP: ShodanHTTP{HTML: "Visit /V1/GRAPHQL for explorer"}}}}
	if !matchFilter("http.html:\"/v1/graphql\"", host) {
		t.Error("html match should be case-insensitive")
	}
}
