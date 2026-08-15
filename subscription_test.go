package main

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestSetSubscriptionNamePrefixesNodesAndReferences(t *testing.T) {
	nodes := SubscriptionData{
		Proxies: []map[string]any{
			{"name": " HK  01 "},
			{"name": "HK  01", "dialer-proxy": " HK  01 "},
		},
		SubInfos: []*SubscriptionMeta{{Name: "provider-name"}},
	}

	setSubscriptionName(&nodes, "Premium A")

	if got := nodes.Proxies[0]["name"]; got != "[Premium A] HK 01" {
		t.Fatalf("unexpected first node name: %v", got)
	}
	if got := nodes.Proxies[1]["name"]; got != "[Premium A] HK 01 #2" {
		t.Fatalf("unexpected duplicate node name: %v", got)
	}
	if got := nodes.Proxies[1]["dialer-proxy"]; got != "[Premium A] HK 01" {
		t.Fatalf("dialer-proxy was not updated: %v", got)
	}
	if nodes.SubInfos[0].Name != "Premium A" {
		t.Fatalf("subscription display name was not updated")
	}
}

func TestOutputFilenameAndSubInfoFormatting(t *testing.T) {
	if got := normalizeOutputFilename("folder\\My Combined.yml"); got != "My Combined.yaml" {
		t.Fatalf("unexpected output filename: %q", got)
	}

	base, err := yaml.Marshal(map[string]any{
		"proxies":      []any{},
		"proxy-groups": []any{},
	})
	if err != nil {
		t.Fatal(err)
	}
	expire := time.Date(2026, 9, 1, 0, 0, 0, 0, time.Local).Unix()
	result, err := addSubInfoGroup(string(base), []*SubscriptionMeta{{
		Name:     "Airport",
		Download: 25 * 1024 * 1024 * 1024,
		Total:    100 * 1024 * 1024 * 1024,
		Expire:   expire,
	}})
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"Airport", "已用 25.0 GB / 100.0 GB", "剩余 75.0 GB", "到期 2026-09-01"} {
		if !strings.Contains(result, fragment) {
			t.Errorf("Sub Info does not contain %q:\n%s", fragment, result)
		}
	}
}
