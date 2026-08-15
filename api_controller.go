package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed ui.html
var uiHTML string

//go:embed example/script.js
var exampleScript string

//go:embed example/template.yaml
var exampleTemplate string

// setupRouter 配置路由
func setupRouter() (r *gin.Engine) {
	r = gin.Default()
	err := r.SetTrustedProxies([]string{"127.0.0.1", "172.17.0.0/16"})
	if err != nil {
		panic(err)
	}

	r.GET("/ping", handlePing)
	r.GET("/sub", handleSubscription)
	r.GET("/ui", handleUI)
	r.GET("/example/script.js", handleExampleScript)
	r.GET("/example/template.yaml", handleExampleTemplate)

	return
}

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func handleUI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, uiHTML)
}

func handleExampleScript(c *gin.Context) {
	c.Data(http.StatusOK, "application/javascript; charset=utf-8", []byte(exampleScript))
}

func handleExampleTemplate(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", []byte(exampleTemplate))
}

// handleSubscription 处理订阅转换请求
// 支持多个订阅合并、流量统计、用量信息显示
func handleSubscription(c *gin.Context) {
	subs := c.QueryArray("sub")
	subNames := c.QueryArray("subName")
	scriptUrl := c.Query("script")
	templateUrl := c.Query("template")
	outputFilename := c.Query("filename")
	userToken := c.Query("token")
	legacyRelayRaw := c.Query("legacyRelay")
	if legacyRelayRaw == "" {
		legacyRelayRaw = c.Query("legacy_relay")
	}
	legacyRelay := false

	// 鉴权
	if userToken != Token {
		L().Warn("Unauthorized request received")
		c.String(http.StatusUnauthorized, "Unauthorized request")
		return
	}

	if legacyRelayRaw != "" {
		parsedLegacyRelay, err := strconv.ParseBool(legacyRelayRaw)
		if err != nil {
			c.String(http.StatusBadRequest, "legacyRelay must be a boolean")
			return
		}
		legacyRelay = parsedLegacyRelay
	}

	// 参数校验
	if len(subs) == 0 || scriptUrl == "" || templateUrl == "" {
		c.String(http.StatusBadRequest, "sub, script and template are required")
		return
	}

	// 提取所有订阅的节点
	allProxies := make([]SubscriptionData, 0, len(subs))
	for i, sub := range subs {
		name := fmt.Sprintf("订阅%02d", i+1)
		proxies, err := ExtractProxies(sub, name)
		if err != nil {
			L().Error(err.Error())
			c.String(http.StatusInternalServerError, fmt.Sprintf("%s:\n%s", sub, err.Error()))
			return
		}
		if i < len(subNames) {
			setSubscriptionName(&proxies, subNames[i])
		}
		allProxies = append(allProxies, proxies)
	}

	// 合并订阅数据
	mergedProxies := mergeProxies(allProxies)

	// 获取模板和脚本
	template, err := FetchString(templateUrl)
	if err != nil {
		L().Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	script, err := FetchString(scriptUrl)
	if err != nil {
		L().Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 执行 JS 脚本生成配置
	result, err := ExecJs(script, template, mergedProxies, legacyRelay)
	if err != nil {
		L().Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 添加用量信息节点组
	finalResult, err := addSubInfoGroup(result, mergedProxies.SubInfos)
	if err != nil {
		L().Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 设置响应头
	for h, v := range mergedProxies.TransparentHeaders {
		c.Header(h, v)
	}
	if filename := normalizeOutputFilename(outputFilename); filename != "" {
		c.Header("Content-Disposition", fmt.Sprintf(
			"attachment; filename*=UTF-8''%s", url.PathEscape(filename),
		))
	}
	c.String(http.StatusOK, finalResult)
}

func normalizeOutputFilename(filename string) string {
	filename = strings.TrimSpace(strings.ReplaceAll(filename, "\\", "/"))
	filename = path.Base(filename)
	if filename == "." || filename == "/" || filename == "" {
		return ""
	}
	if strings.EqualFold(path.Ext(filename), ".yaml") || strings.EqualFold(path.Ext(filename), ".yml") {
		filename = strings.TrimSuffix(filename, path.Ext(filename))
	}
	return filename + ".yaml"
}
