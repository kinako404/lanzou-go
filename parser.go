package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	mobileUserAgent = "Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36"
	maxPageSize     = 2 << 20
)

var (
	legacyDownURLPattern = regexp.MustCompile(`(?i)<a href="([^"]+)"[^>]*id="downurl"`)
	scriptDownURLPattern = regexp.MustCompile(`(?i)\blink\.href\s*=\s*['"]([^'"]+)['"]`)
	prefixPattern        = regexp.MustCompile(`(?i)(?:var\s+)?vkjxld\s*=\s*['"]([^'"]+)`)
	tokenPattern         = regexp.MustCompile(`(?i)(?:var\s+)?hyggid\s*=\s*['"]([^'"]+)`)
)

type LogEntry struct {
	Time    string
	Type    string
	Message string
}

func parseLanzouURL(ctx context.Context, client *http.Client, targetURL string) (string, []LogEntry, error) {
	logs := make([]LogEntry, 0, 12)
	addLog := func(message, logType string) {
		logs = append(logs, LogEntry{
			Time:    time.Now().Format("15:04:05"),
			Type:    logType,
			Message: message,
		})
	}
	fail := func(err error) (string, []LogEntry, error) {
		addLog("解析失败: "+err.Error(), "error")
		return "", logs, err
	}

	addLog("开始解析蓝奏云链接", "info")
	addLog("目标URL: "+targetURL, "info")
	baseURL, err := url.Parse(targetURL)
	if err != nil || baseURL.Host == "" || (baseURL.Scheme != "http" && baseURL.Scheme != "https") {
		return fail(fmt.Errorf("目标链接必须是有效的 HTTP 或 HTTPS 地址"))
	}

	addLog("步骤1: 获取初始页面", "info")
	firstPage, err := fetchPage(ctx, client, targetURL, "https://www.lanzou.com/", "初始页面")
	if err != nil {
		return fail(err)
	}
	addLog("初始页面获取成功", "success")

	match := legacyDownURLPattern.FindStringSubmatch(firstPage)
	if match == nil {
		match = scriptDownURLPattern.FindStringSubmatch(firstPage)
	}
	if match == nil {
		return fail(fmt.Errorf("无法从页面提取下载链接"))
	}
	downReference, err := url.Parse(match[1])
	if err != nil {
		return fail(fmt.Errorf("下载链接格式无效: %w", err))
	}
	downURL := baseURL.ResolveReference(downReference).String()
	addLog("提取到下载链接: "+downURL, "success")

	addLog("步骤2: 获取第二页内容", "info")
	secondPage, err := fetchPage(ctx, client, downURL, targetURL, "第二页")
	if err != nil {
		return fail(err)
	}
	addLog("第二页获取成功", "success")

	addLog("步骤3: 提取URL参数", "info")
	prefixMatch := prefixPattern.FindStringSubmatch(secondPage)
	tokenMatch := tokenPattern.FindStringSubmatch(secondPage)
	if prefixMatch == nil || tokenMatch == nil {
		return fail(fmt.Errorf("无法提取URL参数"))
	}
	addLog(fmt.Sprintf("提取参数: vkjxld=%s..., hyggid=%s...", prefixMatch[1][:min(20, len(prefixMatch[1]))], tokenMatch[1][:min(20, len(tokenMatch[1]))]), "success")

	finalURL := prefixMatch[1] + tokenMatch[1]
	parsedFinalURL, err := url.Parse(finalURL)
	if err != nil || parsedFinalURL.Host == "" || (parsedFinalURL.Scheme != "http" && parsedFinalURL.Scheme != "https") {
		return fail(fmt.Errorf("最终直链格式无效"))
	}
	addLog("拼接最终URL: "+finalURL, "success")
	addLog("✅ 解析完成！", "success")
	return finalURL, logs, nil
}

func fetchPage(ctx context.Context, client *http.Client, pageURL, referer, step string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("%s请求创建失败: %w", step, err)
	}
	request.Header.Set("User-Agent", mobileUserAgent)
	request.Header.Set("Referer", referer)

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("%s请求失败: %w", step, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("%s请求失败: %d", step, response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxPageSize+1))
	if err != nil {
		return "", fmt.Errorf("%s读取失败: %w", step, err)
	}
	if len(body) > maxPageSize {
		return "", fmt.Errorf("%s内容超过大小限制", step)
	}
	return string(body), nil
}
