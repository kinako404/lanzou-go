# 蓝奏云直链解析工具（Go 版）

使用 Go 标准库实现的蓝奏云分享链接解析服务，支持新旧分享页面、直接重定向和调试模式，无需 Node.js 或 Cloudflare。

## 下载

从 [Releases](https://github.com/kinako404/lanzou-go/releases/latest) 下载 Linux 或 Windows amd64 预编译二进制，以及示例配置文件。

## 配置

默认读取运行目录下的 `config.json`：

```json
{
  "port": 18080
}
```

修改端口后重启服务生效。可通过 `-config` 指定其他配置文件：

```bash
./bin/lanzou-linux-amd64 -config /path/to/config.json
```

服务监听所有网卡的指定端口，端口范围为 1–65535。

## 运行

Linux：

```bash
./bin/lanzou-linux-amd64
```

Windows：

```powershell
.\bin\lanzou-windows-amd64.exe -config config.json
```

## 使用

- 首页：`http://服务器地址:18080/`
- 直接下载：`http://服务器地址:18080/?url=蓝奏云分享链接`
- 调试模式：`http://服务器地址:18080/?url=蓝奏云分享链接&debug=1`

需要提取码的私密链接或已失效链接可能无法解析。

## 从源码构建

```bash
go build -o bin/lanzou-linux-amd64 .
```

交叉编译：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/lanzou-linux-amd64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/lanzou-windows-amd64.exe .
```

## 许可证

MIT，详见 [LICENSE](LICENSE)。
