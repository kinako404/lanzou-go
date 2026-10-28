package main

import "html/template"

type DebugPageData struct {
	Logs      []LogEntry
	FinalURL  string
	DirectURL string
}

type ErrorPageData struct {
	Message   string
	TargetURL string
}

var debugPageTemplate = template.Must(template.New("debug").Parse(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>蓝奏云解析调试模式</title>
    <style>
        body { font-family: 'Courier New', monospace; background: #000; color: #00ff00; padding: 20px; }
        .log-entry { margin: 10px 0; padding: 5px; border-left: 3px solid transparent; }
        .info { border-left-color: #17a2b8; color: #17a2b8; }
        .success { border-left-color: #28a745; color: #28a745; }
        .warning { border-left-color: #ffc107; color: #ffc107; }
        .error { border-left-color: #dc3545; color: #dc3545; }
        .timestamp { color: #6c757d; font-size: 12px; margin-right: 10px; }
        .result { background: #155724; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .final-url { color: #00ffff; word-break: break-all; }
        .action-buttons { margin: 20px 0; }
        button { background: #007bff; color: white; border: none; padding: 10px 15px; border-radius: 5px; cursor: pointer; margin-right: 10px; }
        button:hover { background: #0056b3; }
    </style>
</head>
<body>
    <h1>🔍 蓝奏云解析调试模式</h1>
    <div id="logs">
        {{range .Logs}}
        <div class="log-entry {{.Type}}">
            <span class="timestamp">{{.Time}}</span> {{.Message}}
        </div>
        {{end}}
    </div>
    <div class="result">
        <h3>✅ 解析成功！</h3>
        <p><strong>最终直链：</strong></p>
        <p class="final-url">{{.FinalURL}}</p>
    </div>
    <div class="action-buttons">
        <button onclick="window.open('{{.FinalURL}}')">⬇️ 直接下载</button>
        <button onclick="window.location.href='{{.DirectURL}}'">🚀 直接访问（非调试模式）</button>
        <button onclick="copyToClipboard('{{.FinalURL}}')">📋 复制直链</button>
    </div>
    <script>
        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(() => alert('已复制到剪贴板！'));
        }
    </script>
</body>
</html>`))

var errorPageTemplate = template.Must(template.New("error").Parse(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>解析失败</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .error { color: #dc3545; background: #f8d7da; padding: 20px; border-radius: 5px; }
        .debug-link { color: #007bff; text-decoration: underline; cursor: pointer; }
    </style>
</head>
<body>
    <h1>❌ 解析失败</h1>
    <div class="error">
        <p><strong>错误信息:</strong> {{.Message}}</p>
        <p><strong>目标链接:</strong> {{.TargetURL}}</p>
    </div>
    <p>💡 请检查链接是否正确</p>
    <p>
        <a href="/">返回首页</a> |
        <span class="debug-link" onclick="location.href=location.href + '&debug=1'">查看详细错误信息</span>
    </p>
</body>
</html>`))
