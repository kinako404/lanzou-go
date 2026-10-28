package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Port int `json:"port"`
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取配置文件: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("解析配置文件: %w", err)
	}
	if config.Port < 1 || config.Port > 65535 {
		return Config{}, fmt.Errorf("端口必须在 1 到 65535 之间，当前为 %d", config.Port)
	}
	return config, nil
}

func newHandler(client *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			_, _ = w.Write([]byte(homePage))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
		defer cancel()
		finalURL, logs, err := parseLanzouURL(ctx, client, targetURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = errorPageTemplate.Execute(w, ErrorPageData{
				Message:   err.Error(),
				TargetURL: targetURL,
			})
			return
		}
		if r.URL.Query().Get("debug") != "" {
			_ = debugPageTemplate.Execute(w, DebugPageData{
				Logs:      logs,
				FinalURL:  finalURL,
				DirectURL: "/?url=" + url.QueryEscape(targetURL),
			})
			return
		}
		http.Redirect(w, r, finalURL, http.StatusFound)
	})
}

func main() {
	configPath := flag.String("config", "config.json", "配置文件路径")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	client := newHTTPClient()
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           newHandler(client),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("蓝奏云直链解析服务已监听 %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

const homePage = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>蓝奏云直链下载工具</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .container { background: #f8f9fa; padding: 30px; border-radius: 10px; }
        code { background: #eee; padding: 2px 5px; border-radius: 3px; }
        .example { background: #e8f4f8; padding: 15px; border-left: 4px solid #007bff; }
        .debug { background: #fff3cd; padding: 15px; border-left: 4px solid #ffc107; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔗 蓝奏云直链下载工具</h1>
        <p>使用方式：在URL后添加 <code>?url=你的蓝奏云链接</code></p>
        <div class="example">
            <strong>示例：</strong><br>
            <code>http://你的服务器地址:18080/?url=https://wwbnk.lanzouv.com/xxxxxxxx</code>
        </div>
        <p>📝 功能：</p>
        <ul>
            <li>自动解析蓝奏云分享链接</li>
            <li>直接重定向到文件下载</li>
        </ul>
        <div class="debug">
            <p>🔍 <strong>调试模式</strong>：添加 <code>&amp;debug=1</code> 参数查看解析过程</p>
            <code>http://你的服务器地址:18080/?url=蓝奏云链接&amp;debug=1</code>
        </div>
    </div>
</body>
</html>`
