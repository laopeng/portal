// store 包管理服务状态的并发安全读写和 JSON 文件持久化。
//
// 所有公开方法都是线程安全的——读写操作通过 sync.RWMutex 保护。
// 持久化采用脏标记（dirty flag）优化：仅在状态实际变更时才写入磁盘。
//
// 自动发现的服务离线超过 7 天会被自动清理（CleanupZombies）。
// 手动添加的服务永不过期。

package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 持久化文件的版本号，用于未来迁移
const stateVersion = 1

// Service 表示一个被监控的服务
type Service struct {
	ID          string     `json:"id"`                    // 内部标识（port-{port}）
	Name        string     `json:"name"`                  // 服务显示名称
	URL         string     `json:"url"`                   // 服务访问地址
	Icon        string     `json:"icon"`                  // 显示图标（emoji）
	Description string     `json:"description"`           // 服务用途描述
	Category    string     `json:"category,omitempty"`    // 服务分类: ai, infra, data, dev
	Online      bool       `json:"online"`                // 当前是否在线
	LatencyMs   int64      `json:"latency_ms"`            // 最近一次探测延迟（毫秒）
	LastChecked time.Time  `json:"last_checked"`          // 最近一次探测时间
	LastOnline  *time.Time `json:"last_online,omitempty"` // 最近一次在线时间
	Error       string     `json:"error,omitempty"`       // 最近一次探测错误信息
	Manual      bool       `json:"manual"`                // 是否为手动添加的服务
	Port        int        `json:"port"`                  // 端口号
}

// ServiceStore 是服务状态的内存存储，支持并发安全访问和 JSON 持久化。
//
// 使用模式:
//
//	store, _ := New("/path/to/state.json")
//	store.UpsertProbed(port, url, name, icon, desc, cat, err, online, latency)
//	services := store.GetAll()
//	store.CleanupZombies(7 * 24 * time.Hour)
type ServiceStore struct {
	mu       sync.RWMutex // 保护 services 切片和 dirty 标记
	services []Service    // 当前所有服务
	dirty    bool         // 内存数据是否与磁盘不一致
	path     string       // 持久化文件路径
}

// 持久化文件结构（包含版本号用于向后兼容）
type stateFile struct {
	Version  int       `json:"version"`
	Services []Service `json:"services"`
}

// New 创建或加载 ServiceStore。
// 如果持久化文件存在则加载，不存在则创建空存储。
func New(path string) (*ServiceStore, error) {
	store := &ServiceStore{path: path, services: []Service{}}
	if err := store.load(); err != nil {
		if os.IsNotExist(err) {
			return store, nil // 文件不存在 -> 空存储，正常
		}
		return nil, err
	}
	return store, nil
}

// load 从磁盘加载 state.json
func (s *ServiceStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	var state stateFile
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("解析 state.json 失败: %w", err)
	}
	s.services = state.Services
	return nil
}

// save 将当前状态写入磁盘，含缩进格式便于人工阅读
func (s *ServiceStore) save() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("创建状态目录失败: %w", err)
	}
	state := stateFile{Version: stateVersion, Services: s.services}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化状态失败: %w", err)
	}
	return os.WriteFile(s.path, data, 0600)
}

// FlushIfDirty 仅在数据变更时将内存状态写入磁盘。
// 用于优雅关闭时确保数据不丢失。
func (s *ServiceStore) FlushIfDirty() error {
	if !s.dirty {
		return nil
	}
	if err := s.save(); err != nil {
		return err
	}
	s.dirty = false
	return nil
}

// GetAll 返回当前所有服务的快照拷贝。
// 拷贝操作在持锁期间完成，返回后可安全使用。
func (s *ServiceStore) GetAll() []Service {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Service, len(s.services))
	copy(result, s.services)
	return result
}

// UpsertProbed 根据探测结果更新或创建服务记录。
// 如果同一个端口（port）的服务已存在则更新，否则创建新记录。
// 并发安全，自动标记脏数据。
func (s *ServiceStore) UpsertProbed(port int, url, name, icon, desc, cat, errMsg string, online bool, latencyMs int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("port-%d", port)
	now := time.Now()

	// 查找已有记录并更新
	for i, svc := range s.services {
		if svc.ID == id {
			s.services[i].Name = name
			s.services[i].URL = url
			s.services[i].Icon = icon
			s.services[i].Description = desc
			s.services[i].Category = cat
			s.services[i].Online = online
			s.services[i].LatencyMs = latencyMs
			s.services[i].LastChecked = now
			s.services[i].Error = errMsg
			if online {
				t := now
				s.services[i].LastOnline = &t
			}
			s.dirty = true
			return
		}
	}

	// 新服务，追加到列表
	svc := Service{
		ID:          id,
		Name:        name,
		URL:         url,
		Icon:        icon,
		Description: desc,
		Category:    cat,
		Online:      online,
		LatencyMs:   latencyMs,
		LastChecked: now,
		Error:       errMsg,
		Manual:      false,
		Port:        port,
	}
	if online {
		t := now
		svc.LastOnline = &t
	}
	s.services = append(s.services, svc)
	s.dirty = true
}

// CleanupZombies 清理长时间离线的自动发现服务。
// ttl 指定离线多久后清理（如 7*24*time.Hour）。
// 手动添加的服务（Manual=true）永不过期。
func (s *ServiceStore) CleanupZombies(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-ttl)
	filtered := s.services[:0]
	removed := false

	for _, svc := range s.services {
		// 手动添加的服务永久保留
		if svc.Manual {
			filtered = append(filtered, svc)
			continue
		}
		// 长时间离线的自动发现服务 -> 清理
		if !svc.Online && svc.LastChecked.Before(cutoff) {
			removed = true
			continue
		}
		filtered = append(filtered, svc)
	}

	if removed {
		s.services = filtered
		s.dirty = true
	}
}

// CleanupNonConfigured 清理不在配置端口列表中的服务。
// 手动添加的服务（Manual=true）永久保留。
func (s *ServiceStore) CleanupNonConfigured(configuredPorts map[string]bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := s.services[:0]
	removed := false

	for _, svc := range s.services {
		// 手动添加的服务永久保留
		if svc.Manual {
			filtered = append(filtered, svc)
			continue
		}
		// 端口在配置中 -> 保留
		portStr := fmt.Sprintf("%d", svc.Port)
		if configuredPorts[portStr] {
			filtered = append(filtered, svc)
			continue
		}
		// 端口不在配置中 -> 清理
		removed = true
	}

	if removed {
		s.services = filtered
		s.dirty = true
	}
}
