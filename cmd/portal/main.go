// Portal v2 — 本地服务门户 + 认证
//
// v2 新增: 用户名+密码认证（SHA-256 + salt）、Session Cookie、
// 反向代理 /proxy/:port/*、用户管理 CLI。
//
// 构建: go build -o portal ./cmd/portal/
// 运行: ./portal

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"portal/internal/config"
	"portal/internal/prober"
	"portal/internal/store"

	"portal/internal/auth"
	authstore "portal/internal/auth/store"
	"portal/internal/proxy"
)

var (
	version = "2.0.0"
	started time.Time
)

var dashboardPath, loginPath string

func main() {
	started = time.Now()

	dashboardFlag := flag.String("dashboard", "", "自定义 Dashboard HTML 路径")
	loginFlag := flag.String("login", "", "自定义登录页 HTML 路径")
	portFlag := flag.Int("port", 0, "Portal 监听端口")
	configFlag := flag.String("config", "", "配置文件路径")
	flag.Parse()

	// CLI 子命令: ./portal user add/passwd/list/remove
	if len(flag.Args()) > 0 && flag.Args()[0] == "user" {
		handleUserCommand(flag.Args()[1:])
		return
	}

	// HTML 路径
	resolve := func(name string) string {
		exe, _ := os.Executable()
		if exe == "" {
			return name
		}
		return filepath.Join(filepath.Dir(exe), name)
	}
	dashboardPath = resolve("dashboard.html")
	loginPath = resolve("login.html")
	if *dashboardFlag != "" {
		dashboardPath = *dashboardFlag
	}
	if *loginFlag != "" {
		loginPath = *loginFlag
	}

	// 配置
	cfgPath := *configFlag
	if cfgPath == "" {
		var err error
		cfgPath, err = config.Path()
		if err != nil {
			log.Fatalf("无法确定配置文件路径: %v", err)
		}
	}
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Printf("未找到配置文件，正在生成默认配置: %s", cfgPath)
		if err := config.GenerateDefault(cfgPath); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Printf("警告: 解析 %s 失败: %v — 使用默认配置", cfgPath, err)
		defaultCfg := config.Default()
		cfg = &defaultCfg
	}
	if *portFlag != 0 {
		cfg.Port = *portFlag
	}

	// 认证存储
	portalDir, err := config.Dir()
	if err != nil {
		log.Fatalf("无法确定 Portal 数据目录: %v", err)
	}
	authStore, err := authstore.New(portalDir)
	if err != nil {
		log.Fatalf("初始化认证存储失败: %v", err)
	}
	if !authStore.UserExists() {
		log.Println("⚠ 尚未创建任何用户账号。请运行: ./portal user add <用户名>")
	}

	// 服务状态存储
	statePath := filepath.Join(portalDir, "state.json")
	serviceStore, err := store.New(statePath)
	if err != nil {
		log.Fatalf("初始化状态存储失败: %v", err)
	}

	// 探测引擎
	probeEngine := prober.New(cfg, serviceStore)
	probeEngine.ProbeAll()
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.ScanInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			probeEngine.ProbeAll()
		}
	}()

	// 认证中间件
	authMw := auth.Auth(authStore)

	// 路由
	mux := http.NewServeMux()

	// 无需认证
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			if authStore.Authenticate(username, password) {
				// 轮转：销毁旧会话，创建新会话
				authStore.UpdateLastLogin(username)
				sid := authStore.CreateSession(username, 24*time.Hour)

				http.SetCookie(w, &http.Cookie{
					Name:     auth.SessionCookie,
					Value:    sid,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   86400,
				})
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			http.Redirect(w, r, "/login?error=auth", http.StatusFound)
			return
		}
		serveHTMLFile(w, loginPath)
	})
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		svcs := serviceStore.GetAll()
		online := 0
		for _, s := range svcs {
			if s.Online {
				online++
			}
		}
		writeJSON(w, 200, map[string]interface{}{
			"status":              "ok",
			"version":             version,
			"uptime":              time.Since(started).String(),
			"services_configured": len(svcs),
			"services_online":     online,
			"last_scan":           probeEngine.LastScan().Format(time.RFC3339),
		})
	})

	// /api/ 代理转发：浏览器 XHR/fetch 发到 /api/xxx，转发到下游服务
	// 类似 Nginx 的 location /api/ { proxy_pass http://localhost:5000/; }
	// Go ServeMux longest-prefix 优先：/api/health 等具体路由先匹配
	mux.Handle("/api/", proxy.Handler("5000"))

	// 需要认证
	mux.Handle("/", authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serveHTMLFile(w, dashboardPath)
			return
		}
		http.NotFound(w, r)
	})))
	mux.Handle("/api/services", authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]interface{}{"services": serviceStore.GetAll()})
	})))
	mux.Handle("/api/probe", authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, map[string]string{"error": "仅支持 POST 方法"})
			return
		}
		go probeEngine.ProbeAll()
		writeJSON(w, 200, map[string]string{"status": "扫描已触发"})
	})))
	mux.Handle("/proxy/", authMw(proxy.Handler("5000")))
	mux.Handle("/logout", authMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(auth.SessionCookie); err == nil {
			authStore.DestroySession(c.Value)
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})))
	mux.HandleFunc("/api/me", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.SessionCookie)
		if err != nil {
			writeJSON(w, 401, map[string]string{"error": "未登录"})
			return
		}
		username, valid := authStore.ValidateSession(cookie.Value)
		if !valid {
			writeJSON(w, 401, map[string]string{"error": "会话无效"})
			return
		}
		writeJSON(w, 200, map[string]string{"username": username})
	})

	// CSP middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; img-src 'self' data:; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; connect-src 'self' http://localhost:* https://localhost:*")
		mux.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)
	srv := &http.Server{Addr: addr, Handler: handler}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("正在关闭...")
		serviceStore.FlushIfDirty()
		srv.Close()
	}()

	log.Printf("Portal v%s 已启动: http://%s", version, addr)
	log.Printf("认证模式: 密码登录")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("服务器错误: %v", err)
	}
}

func serveHTMLFile(w http.ResponseWriter, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "文件未找到", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ─── CLI 用户管理 ────────────────────────────────────────

func handleUserCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("用法: ./portal user <add|passwd|list|remove> [参数]")
		os.Exit(1)
	}

	portalDir, err := config.Dir()
	if err != nil {
		log.Fatalf("无法确定 Portal 数据目录: %v", err)
	}
	as, err := authstore.New(portalDir)
	if err != nil {
		log.Fatalf("打开用户存储失败: %v", err)
	}

	switch args[0] {
	case "add":
		var username string
		r := bufio.NewReader(os.Stdin)

		if len(args) > 1 {
			username = args[1]
		} else {
			fmt.Print("用户名: ")
			line, _ := r.ReadString('\n')
			username = strings.TrimSpace(line)
		}
		if username == "" {
			log.Fatal("用户名不能为空")
		}

		fmt.Print("密码: ")
		line, _ := r.ReadString('\n')
		password := strings.TrimSpace(line)
		if password == "" {
			log.Fatal("密码不能为空")
		}

		fmt.Print("确认密码: ")
		line2, _ := r.ReadString('\n')
		confirm := strings.TrimSpace(line2)
		if password != confirm {
			log.Fatal("两次密码不一致")
		}

		if err := as.CreateUser(username, password); err != nil {
			log.Fatalf("创建用户失败: %v", err)
		}
		fmt.Printf("用户 '%s' 创建成功\n", username)

	case "list":
		// 简单实现：直接读 users.json
		data, _ := os.ReadFile(filepath.Join(portalDir, "users.json"))
		fmt.Println(string(data))

	default:
		fmt.Printf("未知子命令: %s\n用法: ./portal user <add|passwd|list|remove>\n", args[0])
		os.Exit(1)
	}
}
