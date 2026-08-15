package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
)

func TestBuildTemplateNormalizesRulesetBehaviors(t *testing.T) {
	rulesets := []*Ruleset{
		{tag: "Domain", url: "domain.list", behavior: "domain", content: "example.com\n+.example.org\n"},
		{tag: "IP", url: "ip.list", behavior: "ipcidr", content: "192.0.2.0/24\n2001:db8::/32\n"},
		{tag: "Classic", url: "classic.list", behavior: "classical", content: "DOMAIN-SUFFIX,example.net\n"},
	}

	config, err := BuildTemplate("mode: rule\n", SubscriptionData{}, rulesets)
	if err != nil {
		t.Fatalf("BuildTemplate returned an error: %v", err)
	}

	want := []string{
		"DOMAIN,example.com,Domain",
		"DOMAIN-SUFFIX,example.org,Domain",
		"IP-CIDR,192.0.2.0/24,IP",
		"IP-CIDR6,2001:db8::/32,IP",
		"DOMAIN-SUFFIX,example.net,Classic",
	}
	rules, ok := config["rules"].([]string)
	if !ok {
		t.Fatalf("rules has unexpected type %T", config["rules"])
	}
	if strings.Join(rules, "\n") != strings.Join(want, "\n") {
		t.Fatalf("unexpected rules:\n%v\nwant:\n%v", rules, want)
	}
}

func TestExampleConfigHasSelfContainedGroupsWithoutRelay(t *testing.T) {
	scriptPath := filepath.Join("example", "script.js")
	templatePath := filepath.Join("example", "template.yaml")
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	template, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatal(err)
	}

	vm := goja.New()
	if _, err = vm.RunString(string(script)); err != nil {
		t.Fatalf("example script is not accepted by goja: %v", err)
	}

	var config map[string]any
	if err = yaml.Unmarshal(template, &config); err != nil {
		t.Fatalf("invalid example template: %v", err)
	}
	config["proxies"] = []map[string]any{
		{"name": "[Premium A] HK 01", "type": "ss", "server": "127.0.0.1", "port": 10001},
		{"name": "[流量优选] SG 01", "type": "ss", "server": "127.0.0.1", "port": 10003},
		{"name": "套餐剩余 100 GB", "type": "ss", "server": "127.0.0.1", "port": 10002},
	}
	config["rules"] = []string{"DOMAIN,example.com,PROXY"}

	var buildConfig func(map[string]any, bool)
	if err = vm.ExportTo(vm.Get("buildConfig"), &buildConfig); err != nil {
		t.Fatalf("cannot export buildConfig: %v", err)
	}
	buildConfig(config, true)

	output, err := Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(output), "relay") {
		t.Fatalf("generated config contains relay:\n%s", output)
	}
	if strings.Contains(output, "rule-providers") || strings.Contains(output, "RULE-SET") {
		t.Fatalf("generated config is not self-contained")
	}
	if strings.Contains(output, "套餐剩余") {
		t.Fatalf("information node was not filtered")
	}
	if !strings.Contains(output, "[流量优选] SG 01") {
		t.Fatalf("custom subscription name was mistaken for an information node")
	}

	groups := proxyGroupNames(t, config["proxy-groups"])
	required := []string{
		"🚀 节点选择", "⚡ 自动选择", "🛟 故障转移", "🛑 广告拦截",
		"🤖 AI 服务", "📹 YouTube", "🔍 Google 服务", "Ⓜ️ Microsoft 服务",
		"🍏 Apple 服务", "📲 Telegram", "💬 社交媒体", "🎬 国际流媒体",
		"🎮 游戏平台", "🛠️ 开发工具", "☁️ 云服务", "💳 金融支付",
		"🏠 私有网络", "🇨🇳 国内直连", "🐟 漏网之鱼",
	}
	if len(groups) != len(required) {
		t.Fatalf("unexpected proxy group count: got %d, want %d", len(groups), len(required))
	}
	for _, name := range required {
		if !groups[name] {
			t.Errorf("missing proxy group %q", name)
		}
	}

	var definitions [][]string
	if err = vm.ExportTo(vm.Get("RULESET_DEFINITIONS"), &definitions); err != nil {
		t.Fatalf("cannot export ruleset definitions: %v", err)
	}
	targets := make(map[string]string, len(definitions))
	indexes := make(map[string]int, len(definitions))
	for definitionIndex, definition := range definitions {
		if len(definition) != 3 {
			t.Fatalf("invalid ruleset definition: %v", definition)
		}
		if !groups[definition[0]] {
			t.Errorf("ruleset %q targets missing group %q", definition[2], definition[0])
		}
		if _, exists := indexes[definition[2]]; !exists {
			targets[definition[2]] = definition[0]
			indexes[definition[2]] = definitionIndex
		}
	}
	for _, source := range []string{
		"google-gemini", "google-deepmind", "github-copilot", "openai", "anthropic", "category-ai-chat-!cn",
	} {
		if targets[source] != "🤖 AI 服务" {
			t.Errorf("AI ruleset %q targets %q", source, targets[source])
		}
	}
	if indexes["google-gemini"] >= indexes["google"] || indexes["google-deepmind"] >= indexes["google"] {
		t.Error("Google AI rulesets must precede the general Google ruleset")
	}
	if indexes["github-copilot"] >= indexes["github"] {
		t.Error("GitHub Copilot ruleset must precede the general GitHub ruleset")
	}
}

func TestExampleConfigFullGenerationPipeline(t *testing.T) {
	ruleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/geoip/") {
			_, _ = w.Write([]byte("192.0.2.0/24\n2001:db8::/32\n"))
			return
		}
		if r.URL.Path == "/classic.list" {
			_, _ = w.Write([]byte("DOMAIN-SUFFIX,legacy.example\n"))
			return
		}
		_, _ = w.Write([]byte("example.com\n+.example.org\n"))
	}))
	defer ruleServer.Close()

	DbPath = filepath.Join(t.TempDir(), "rules.db")
	InitDb()
	sqlDB, err := orm.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	script, err := os.ReadFile(filepath.Join("example", "script.js"))
	if err != nil {
		t.Fatal(err)
	}
	template, err := os.ReadFile(filepath.Join("example", "template.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	realRulesets := os.Getenv("REAL_RULESETS_TEST") != ""
	if !realRulesets {
		script = []byte(strings.Replace(
			string(script),
			"https://cdn.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@meta/geo/",
			ruleServer.URL+"/geo/",
			1,
		))
	}

	result, err := ExecJs(string(script), string(template), SubscriptionData{
		Proxies: []map[string]any{{
			"name": "[Premium A] HK 01", "type": "ss", "server": "127.0.0.1",
			"port": 10001, "cipher": "aes-128-gcm", "password": "test",
		}},
	}, true)
	if err != nil {
		t.Fatalf("full example generation failed: %v", err)
	}
	if strings.Contains(result, "rule-providers") {
		t.Fatalf("full output contains a client-side rule dependency")
	}
	var generatedConfig map[string]any
	if err = yaml.Unmarshal([]byte(result), &generatedConfig); err != nil {
		t.Fatalf("full output is invalid YAML: %v", err)
	}
	proxyGroupNames(t, generatedConfig["proxy-groups"])
	rules := configRuleStrings(t, generatedConfig["rules"])
	if !containsString(rules, "MATCH,🐟 漏网之鱼") {
		t.Fatalf("full output is missing normalized rules")
	}
	if !realRulesets && (!containsString(rules, "DOMAIN-SUFFIX,example.org,🛑 广告拦截") ||
		!containsString(rules, "IP-CIDR,192.0.2.0/24,🏠 私有网络")) {
		t.Fatalf("full output is missing normalized fixture rules")
	}
	if realRulesets && (!containsRuleTarget(rules, "🛑 广告拦截") || !containsRuleTarget(rules, "🏠 私有网络")) {
		t.Fatalf("full output is missing real MetaCubeX rules")
	}
	if mihomoPath := os.Getenv("MIHOMO_TEST_BINARY"); mihomoPath != "" {
		configPath := filepath.Join(t.TempDir(), "config.yaml")
		if err = os.WriteFile(configPath, []byte(result), 0600); err != nil {
			t.Fatal(err)
		}
		command := exec.Command(mihomoPath, "-t", "-f", configPath, "-d", t.TempDir())
		if output, commandErr := command.CombinedOutput(); commandErr != nil {
			t.Fatalf("mihomo rejected generated config: %v\n%s", commandErr, output)
		}
	}

	legacyScript := "function rulesets(r) { r('PROXY', '" + ruleServer.URL + "/classic.list'); }"
	legacyResult, err := ExecJs(legacyScript, "mode: rule\n", SubscriptionData{}, false)
	if err != nil {
		t.Fatalf("legacy two-argument rulesets callback failed: %v", err)
	}
	if !strings.Contains(legacyResult, "DOMAIN-SUFFIX,legacy.example,PROXY") {
		t.Fatalf("legacy classical rule was not preserved: %s", legacyResult)
	}
}

func proxyGroupNames(t *testing.T, value any) map[string]bool {
	t.Helper()
	data, err := yaml.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var groups []map[string]any
	if err = yaml.Unmarshal(data, &groups); err != nil {
		t.Fatalf("invalid proxy groups: %v", err)
	}
	result := make(map[string]bool, len(groups))
	for _, group := range groups {
		name, _ := group["name"].(string)
		groupType, _ := group["type"].(string)
		if groupType == "relay" {
			t.Errorf("group %q uses relay", name)
		}
		if result[name] {
			t.Errorf("duplicate proxy group %q", name)
		}
		result[name] = true
	}
	return result
}

func configRuleStrings(t *testing.T, value any) []string {
	t.Helper()
	data, err := yaml.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var rules []string
	if err = yaml.Unmarshal(data, &rules); err != nil {
		t.Fatalf("invalid rules: %v", err)
	}
	return rules
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func containsRuleTarget(rules []string, target string) bool {
	for _, rule := range rules {
		if strings.HasSuffix(rule, ","+target) {
			return true
		}
	}
	return false
}
