// config 包负责 Portal 配置文件的读取、校验、生成和默认值管理。
//
// 配置文件格式为 JSON，存储在 ~/.portal/config.json。
// 首次运行若配置文件不存在，会自动生成默认配置。
//
// 配置项说明:
//   port            - Portal 自身监听端口（默认 8747）
//   scan_interval   - 服务探测间隔，秒（默认 60，最小 5）
//   probe_timeout   - 单次探测超时，秒（默认 5，最小 1）
//   max_concurrency - 最大并发探测数（默认 20，最小 1）
//   port_hints      - 端口号到服务名称/图标/分类的映射表
//   keyword_hints   - HTML title 关键字到服务身份的映射表

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PortHint 定义已知端口的服务身份信息
type PortHint struct {
	Name string `json:"name"`           // 服务显示名称
	Icon string `json:"icon"`           // 服务图标（emoji）
	Desc string `json:"desc,omitempty"` // 服务用途描述
	Cat  string `json:"cat,omitempty"`  // 服务分类: ai, infra, data, dev
}

// KeywordHint 定义 HTML title 关键字的服务匹配规则。
// 当探测到端口的页面标题包含指定关键字时，自动归类。
type KeywordHint struct {
	Keyword string `json:"keyword"` // 匹配关键字（不区分大小写）
	Name    string `json:"name"`    // 匹配后设置的服务名
	Icon    string `json:"icon"`    // 匹配后设置的图标
}

// Config 是 Portal 的完整配置结构
type Config struct {
	Port           int                 `json:"port"`            // 监听端口
	ScanInterval   int                 `json:"scan_interval"`   // 扫描间隔（秒）
	ProbeTimeout   int                 `json:"probe_timeout"`   // 探测超时（秒）
	MaxConcurrency int                 `json:"max_concurrency"` // 最大并发探测数
	PortHints      map[string]PortHint `json:"port_hints"`      // 端口到服务映射
	KeywordHints   []KeywordHint       `json:"keyword_hints"`   // 关键字到服务映射
}

// Default 返回 Portal 的默认配置。
// 包含常见开发端口和 AI 基础设施端口。
func Default() Config {
	return Config{
		Port:           8747,
		ScanInterval:   60,
		ProbeTimeout:   5,
		MaxConcurrency: 20,
		PortHints: map[string]PortHint{
			"5000":  {Name: "资金管理系统", Icon: "💰", Desc: "个人/企业资金收支管理", Cat: "data"},
			"8787":  {Name: "Hermes Web", Icon: "📡", Desc: "AI 网关管理平台，管理 API 路由与密钥", Cat: "ai"},
			"20128": {Name: "9Router", Icon: "🔀", Desc: "AI 基础设施管理，统一管理所有 AI 供应商端点", Cat: "ai"},
		},
		KeywordHints: []KeywordHint{
			{Keyword: "Grafana", Name: "Grafana", Icon: "📈"},
			{Keyword: "Prometheus", Name: "Prometheus", Icon: "🔥"},
			{Keyword: "Nginx", Name: "Nginx", Icon: "🌐"},
			{Keyword: "Jupyter", Name: "Jupyter", Icon: "📓"},
		},
	}
}

// Dir 返回 Portal 配置目录路径（~/.portal/）
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".portal"), nil
}

// Path 返回配置文件完整路径（~/.portal/config.json）
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// GenerateDefault 在指定路径生成默认配置文件。
// 如果目录不存在则自动创建。
func GenerateDefault(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	cfg := Default()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// Load 从文件加载配置，并校验参数合法性。
// 非法值会被自动修正为安全的默认值。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	cfg := Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", path, err)
	}

	// 参数合法性校验与自动修正
	if cfg.Port < 1 || cfg.Port > 65535 {
		cfg.Port = 8747
	}
	if cfg.ScanInterval < 5 {
		cfg.ScanInterval = 60
	}
	if cfg.ProbeTimeout < 1 {
		cfg.ProbeTimeout = 5
	}
	if cfg.MaxConcurrency < 1 {
		cfg.MaxConcurrency = 20
	}

	return &cfg, nil
}
