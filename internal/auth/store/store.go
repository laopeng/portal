// auth/store 包管理用户和会话的 JSON 文件持久化存储。
//
// 使用纯 Go 标准库实现，零外部依赖。
// 用户数据: ~/.portal/users.json
// 会话数据: 内存 map（重启后需重新登录）

package store

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// User 表示一个本地账号
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"` // SHA-256 hex
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

// Session 表示一个活跃的登录会话
type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
}

// Store 管理用户和会话
type Store struct {
	mu       sync.RWMutex
	users    []User
	sessions map[string]*Session
	path     string
}

type userFile struct {
	Version int    `json:"version"`
	Users   []User `json:"users"`
}

// New 打开或创建用户存储
func New(portalDir string) (*Store, error) {
	if err := os.MkdirAll(portalDir, 0700); err != nil {
		return nil, fmt.Errorf("创建 Portal 目录失败: %w", err)
	}

	path := filepath.Join(portalDir, "users.json")
	s := &Store{path: path, sessions: make(map[string]*Session)}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil // 无用户则空存储
		}
		return nil, err
	}

	var uf userFile
	if err := json.Unmarshal(data, &uf); err != nil {
		return nil, fmt.Errorf("解析 users.json 失败: %w", err)
	}
	s.users = uf.Users
	return s, nil
}

func (s *Store) save() error {
	uf := userFile{Version: 1, Users: s.users}
	data, err := json.MarshalIndent(uf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// HashPassword 用 SHA-256 + salt 哈希密码
func HashPassword(plaintext string) (hash, salt string) {
	salt = GenerateID(16)
	h := sha256.Sum256([]byte(salt + plaintext))
	return fmt.Sprintf("%x", h), salt
}

// VerifyPassword 校验密码
func VerifyPassword(plaintext, hash, salt string) bool {
	h := sha256.Sum256([]byte(salt + plaintext))
	expected := fmt.Sprintf("%x", h)
	return subtle.ConstantTimeCompare([]byte(hash), []byte(expected)) == 1
}

// CreateUser 创建新用户
func (s *Store) CreateUser(username, plaintext string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if u.Username == username {
			return fmt.Errorf("用户 '%s' 已存在", username)
		}
	}

	hash, salt := HashPassword(plaintext)
	s.users = append(s.users, User{
		Username:  username,
		Password:  hash,
		Salt:      salt,
		CreatedAt: time.Now(),
	})
	return s.save()
}

// Authenticate 验证用户名和密码，成功返回 true
func (s *Store) Authenticate(username, plaintext string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Username == username {
			return VerifyPassword(plaintext, u.Password, u.Salt)
		}
	}
	return false
}

// UserExists 检查是否存在至少一个用户
func (s *Store) UserExists() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users) > 0
}

// ─── 会话操作（纯内存）───────────────────────────────────────────

// CreateSession 创建新会话
func (s *Store) CreateSession(username string, ttl time.Duration) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sid := GenerateID(32)
	s.sessions[sid] = &Session{
		ID:        sid,
		Username:  username,
		ExpiresAt: time.Now().Add(ttl),
	}
	return sid
}

// ValidateSession 验证会话
func (s *Store) ValidateSession(sessionID string) (username string, valid bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[sessionID]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return "", false
	}
	return sess.Username, true
}

// DestroySession 删除会话
func (s *Store) DestroySession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// DestroyUserSessions 删除用户的所有会话
func (s *Store) DestroyUserSessions(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, sess := range s.sessions {
		if sess.Username == username {
			delete(s.sessions, id)
		}
	}
}

// CleanupExpiredSessions 清理过期会话
func (s *Store) CleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for id, sess := range s.sessions {
		if now.After(sess.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
}

// UpdateLastLogin 更新最后登录时间
func (s *Store) UpdateLastLogin(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.users {
		if u.Username == username {
			s.users[i].LastLogin = time.Now()
			s.save()
			return
		}
	}
}
