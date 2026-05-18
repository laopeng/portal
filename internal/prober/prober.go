// prober 包负责对配置的本地端口进行 HTTP 健康探测，并识别服务身份。
//
// 探测流程（三层识别）:
//   1. 端口映射匹配 — 根据 config.json 的 port_hints 直接查找
//   2. 页面标题提取 — 读取 HTML <title> 作为服务描述
//   3. 关键字规则匹配 — 根据 config.json 的 keyword_hints 重新分类
//
// 并发控制通过信号量（buffered channel）实现，限制同时探测数。
// 错误信息会翻译为中文（如 "connection refused" → "连接被拒绝"）。

package prober

import (
	"context"
	"crypto/tls"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"portal/internal/config"
	"portal/internal/store"
)

// Prober 是服务探测引擎，管理定时扫描和状态更新
type Prober struct {
	cfg      *config.Config      // 配置引用
	store    *store.ServiceStore // 状态存储引用
	client   *http.Client        // 复用的 HTTP 客户端（带超时和 TLS 配置）
	sema     chan struct{}       // 并发控制信号量
	mu       sync.Mutex          // 保护 lastScan
	lastScan time.Time           // 最近一次全量扫描完成时间
}

// New 创建探测引擎
func New(cfg *config.Config, store *store.ServiceStore) *Prober {
	return &Prober{
		cfg:   cfg,
		store: store,
		// HTTP 客户端配置:
		// - 全局超时防止单个请求永久阻塞
		// - TLS 跳过证书验证（本地开发环境可接受）
		// - 跟随重定向以获取最终页面标题
		client: &http.Client{
			Timeout: time.Duration(cfg.ProbeTimeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		sema: make(chan struct{}, cfg.MaxConcurrency),
	}
}

// ProbeAll 对所有配置的端口执行一次全量探测。
// 每个端口在独立的 goroutine 中探测，通过信号量限制并发数。
// 探测完成后将变更写入磁盘，并清理长期离线的僵尸服务。
func (p *Prober) ProbeAll() {
	p.mu.Lock()
	p.lastScan = time.Now()
	p.mu.Unlock()

	var wg sync.WaitGroup

	// 并行为每个配置的端口启动探测 goroutine
	for portStr := range p.cfg.PortHints {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		wg.Add(1)
		go p.probePort(port, &wg)
	}

	wg.Wait()

	// 探测完成后持久化并清理
	p.store.FlushIfDirty()
	p.store.CleanupZombies(7 * 24 * time.Hour) // 离线 7 天的服务自动清理
}

// probePort 探测单个端口并更新 store
func (p *Prober) probePort(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 获取信号量，限制并发探测数
	p.sema <- struct{}{}
	defer func() { <-p.sema }()

	url := fmt.Sprintf("http://localhost:%d", port)
	name, icon, desc, cat := p.identifyService(port)
	online := false
	var latencyMs int64
	var errMsg string

	start := time.Now()

	// 创建带超时的请求上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cfg.ProbeTimeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		p.store.UpsertProbed(port, url, name, icon, desc, cat,
			fmt.Sprintf("无效 URL: %v", err), false, 0)
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		// 错误翻译为中文，方便 Dashboard 显示
		p.store.UpsertProbed(port, url, name, icon, desc, cat,
			translateError(err), false, 0)
		return
	}
	defer resp.Body.Close()

	latencyMs = time.Since(start).Milliseconds()
	online = true

	// 读取页面内容用于标题提取和关键字匹配
	// 限制 64KB 防止大页面内存占用
	bodyPrefix := make([]byte, 65536)
	n, _ := io.ReadFull(resp.Body, bodyPrefix)
	if n == 0 {
		n, _ = resp.Body.Read(bodyPrefix)
	}
	body := string(bodyPrefix[:n])

	// 提取 HTML <title> 作为描述
	extractedTitle := extractTitle(body)
	if extractedTitle != "" {
		// 如果配置中未指定描述，使用提取的标题
		if desc == "" {
			desc = extractedTitle
		}
		// 尝试通过关键字匹配获得更准确的服务名和图标
		if newName, newIcon := p.matchKeyword(extractedTitle); newName != "" {
			name = newName
			icon = newIcon
		}
	}

	p.store.UpsertProbed(port, url, name, icon, desc, cat, errMsg, online, latencyMs)
}

// identifyService 根据端口号识别服务身份。
// 优先查找 port_hints 配置，未找到则返回通用名称。
func (p *Prober) identifyService(port int) (name, icon, desc, cat string) {
	portStr := strconv.Itoa(port)
	if hint, ok := p.cfg.PortHints[portStr]; ok {
		return hint.Name, hint.Icon, hint.Desc, hint.Cat
	}
	return fmt.Sprintf("未识别服务 :%d", port), "❓", "", ""
}

// matchKeyword 根据 HTML 标题中的关键字匹配服务身份。
// 返回匹配到的名称和图标，未匹配返回空字符串。
func (p *Prober) matchKeyword(title string) (name, icon string) {
	titleLower := strings.ToLower(title)
	for _, kw := range p.cfg.KeywordHints {
		if strings.Contains(titleLower, strings.ToLower(kw.Keyword)) {
			return kw.Name, kw.Icon
		}
	}
	return "", ""
}

// 预编译的正则表达式，用于提取 HTML <title> 标签内容
var titleRegex = regexp.MustCompile(`(?is)<title[^>]*>(.+?)</title>`)

// extractTitle 从 HTML 内容中提取 <title> 标签的文本。
// 限制 200 字符，防止标题过长。
func extractTitle(htmlBody string) string {
	matches := titleRegex.FindStringSubmatch(htmlBody)
	if len(matches) < 2 {
		return ""
	}
	title := strings.TrimSpace(matches[1])
	title = collapseWhitespace(title)  // 合并多余空白字符
	title = html.UnescapeString(title) // 解码 HTML 实体
	if len(title) > 200 {
		title = title[:200] + "..."
	}
	return title
}

// collapseWhitespace 将连续空白字符合并为单个空格
func collapseWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// translateError 将 Go HTTP 客户端错误信息翻译为中文
func translateError(err error) string {
	msg := err.Error()
	translations := map[string]string{
		"connection refused":          "连接被拒绝",
		"connect: connection refused": "连接被拒绝",
		"i/o timeout":                 "连接超时",
		"no such host":                "无法解析主机名",
		"EOF":                         "连接意外关闭",
		"connection reset by peer":    "连接被重置",
		"no route to host":            "无法到达主机",
	}
	for eng, chn := range translations {
		if strings.Contains(msg, eng) {
			return chn
		}
	}
	return fmt.Sprintf("连接失败: %s", msg)
}

// LastScan 返回最近一次全量扫描完成的时间
func (p *Prober) LastScan() time.Time {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lastScan
}
