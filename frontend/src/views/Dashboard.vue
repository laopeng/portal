<template>
  <div class="dashboard">
    <header>
      <h1><span class="portal-mark">■</span>本地服务门户</h1>
      <div class="meta">
        <span class="header-stats">{{ onlineCount }} / {{ services.length }} 在线</span>
        <a class="logout-btn" href="/logout">退出</a>
      </div>
    </header>

    <main v-if="services.length === 0 && firstLoad" class="grid">
      <div v-for="i in 2" :key="i" class="card skeleton skeleton-card"></div>
    </main>

    <main v-else-if="services.length === 0" class="grid">
      <div class="empty-state">
        <div class="empty-icon">📭</div>
        <h2>没有配置任何服务</h2>
        <p>请在 <code>~/.portal/config.json</code> 的 <code>port_hints</code> 中添加端口。</p>
      </div>
    </main>

    <main v-else class="grid">
      <template v-for="cat in categories" :key="cat.key">
        <div class="cat-title">{{ catLabel(cat.key) }}</div>
        <a
          v-for="s in cat.services"
          :key="s.id"
          class="card"
          :href="`/proxy/${s.port}/`"
          target="_blank"
          rel="noopener"
          :aria-label="`${s.name} — ${s.online ? '在线' : '离线'}`"
        >
          <div class="card-header">
            <span class="card-icon" aria-hidden="true">{{ s.icon || '❓' }}</span>
            <span class="card-name" :title="s.name">{{ s.name }}</span>
            <span :class="['status-dot', s.online ? 'online' : 'offline']" aria-hidden="true"></span>
            <span class="sr-only">{{ s.online ? '在线' : '离线' }}</span>
          </div>
          <div class="card-body">
            <span class="card-url">{{ s.url }}</span>
            <div v-if="s.description" class="card-desc">{{ s.description }}</div>
            <span v-if="s.online" :class="['card-latency', latencyClass(s.latency_ms)]">{{ s.latency_ms }}ms</span>
            <div v-if="s.error" class="card-error">{{ s.error }}</div>
            <div class="card-time">上次检查：{{ fmtTime(s.last_checked) }}</div>
          </div>
        </a>
      </template>
    </main>

    <div v-if="toasts.length" id="toast-container">
      <div v-for="t in toasts" :key="t.id" :class="['toast', t.leaving ? 'removing' : '']">
        <span class="toast-icon">{{ t.icon }}</span>
        <span>{{ t.name }} {{ t.status }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { fetchServices, fetchHealth } from '../api'

const CAT_LABELS = { ai: '🤖 AI 服务', infra: '🖥️ 基础设施', data: '📊 数据服务', dev: '🛠 开发工具' }

const services = ref([])
const onlineCount = ref(0)
const firstLoad = ref(true)
const toasts = ref([])
let prevOnline = {}
let timer = null
let toastId = 0

const categories = computed(() => {
  const groups = {}
  const uncat = []
  for (const s of services.value) {
    if (s.category) {
      if (!groups[s.category]) groups[s.category] = []
      groups[s.category].push(s)
    } else {
      uncat.push(s)
    }
  }
  const result = Object.keys(groups).sort().map(k => ({ key: k, services: groups[k] }))
  if (uncat.length) result.push({ key: '', services: uncat })
  return result
})

function catLabel(key) {
  return CAT_LABELS[key] || key
}

function latencyClass(ms) {
  if (ms < 100) return 'fast'
  if (ms < 500) return 'medium'
  return 'slow'
}

function fmtTime(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleTimeString()
}

function addToast(icon, name, status) {
  const id = ++toastId
  toasts.value.push({ id, icon, name, status, leaving: false })
  setTimeout(() => {
    const t = toasts.value.find(t => t.id === id)
    if (t) t.leaving = true
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, 250)
  }, 3000)
}

async function refresh() {
  try {
    const r1 = await fetchServices()
    if (r1.status === 401) {
      window.location.href = '/login?reason=expired'
      return
    }
    const svcData = r1
    const svcs = svcData.services || []
    services.value = svcs

    const health = await fetchHealth().catch(() => null)
    if (health) {
      onlineCount.value = health.services_online
    }

    if (!firstLoad.value) {
      for (const s of svcs) {
        const prev = prevOnline[s.id]
        if (prev !== undefined && prev !== s.online) {
          addToast(s.icon || '❓', s.name, s.online ? '已上线' : '已离线')
        }
        prevOnline[s.id] = s.online
      }
    } else {
      for (const s of svcs) {
        prevOnline[s.id] = s.online
      }
    }

    firstLoad.value = false
  } catch {
    // keep last rendered state
  }
}

onMounted(() => {
  refresh()
  timer = setInterval(refresh, 10000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.dashboard {
  --bg: #0d1117;
  --card-bg: #161b22;
  --card-border: #30363d;
  --card-hover-border: #484f58;
  --text-primary: #e6edf3;
  --text-secondary: #8b949e;
  --accent-online: #3fb950;
  --accent-offline: #f85149;
  --accent-warn: #d29922;
  --toast-bg: #21262d;
  --radius: 8px;
  --grid-gap: 20px;
  --max-width: 1200px;
  --font-body: system-ui, -apple-system, sans-serif;
  --font-mono: "JetBrains Mono", "SF Mono", "Cascadia Code", "Fira Code", monospace;

  min-height: 100vh;
  background: var(--bg);
  color: var(--text-primary);
  font-family: var(--font-body);
}

header {
  padding: 24px 32px 0;
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}

h1 {
  font-family: var(--font-mono);
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.portal-mark { color: var(--accent-online); margin-right: 6px; }

.meta {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 12px;
}

.logout-btn {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
  text-decoration: none;
  padding: 4px 10px;
  border: 1px solid var(--card-border);
  border-radius: 4px;
  transition: border-color 0.15s;
}

.logout-btn:hover { border-color: var(--accent-offline); color: var(--accent-offline); }

.grid {
  max-width: var(--max-width);
  margin: 24px auto;
  padding: 0 32px 32px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--grid-gap);
}

.card {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: var(--radius);
  padding: 20px;
  text-decoration: none;
  color: inherit;
  display: block;
  transition: border-color 0.15s, transform 0.15s;
}

.card:hover { border-color: var(--card-hover-border); transform: translateY(-2px); }

.card:focus-visible { outline: 2px solid var(--accent-online); outline-offset: 2px; }

.card-header { display: flex; align-items: center; gap: 12px; margin-bottom: 10px; }

.card-icon { font-size: 28px; line-height: 1; flex-shrink: 0; }

.card-name {
  font-size: 15px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-dot { width: 10px; height: 10px; border-radius: 50%; margin-left: auto; flex-shrink: 0; }

.status-dot.online {
  background: var(--accent-online);
  box-shadow: 0 0 6px var(--accent-online);
  animation: pulse 2s infinite;
}

.status-dot.offline { background: var(--accent-offline); }

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 4px var(--accent-online); }
  50% { box-shadow: 0 0 12px var(--accent-online); }
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0,0,0,0);
  white-space: nowrap;
}

.card-body { font-size: 14px; color: var(--text-secondary); display: flex; flex-direction: column; gap: 4px; }

.card-url {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-desc { font-size: 13px; color: var(--text-secondary); margin-top: 2px; opacity: 0.85; }

.card-latency { font-size: 12px; }

.card-latency.fast { color: var(--accent-online); }

.card-latency.medium { color: var(--accent-warn); }

.card-latency.slow { color: var(--accent-offline); }

.card-error {
  font-size: 12px;
  color: var(--accent-offline);
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-time {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 6px;
  opacity: 0.7;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.empty-icon { font-size: 48px; margin-bottom: 16px; }

.empty-state h2 {
  font-family: var(--font-mono);
  font-size: 18px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-state p { font-size: 14px; margin-bottom: 4px; }

.empty-state code {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 4px;
  padding: 2px 6px;
  font-family: var(--font-mono);
  font-size: 12px;
}

.skeleton {
  background: linear-gradient(90deg, var(--card-bg) 25%, #21262d 50%, var(--card-bg) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: var(--radius);
}

.skeleton-card { height: 130px; }

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.cat-title {
  grid-column: 1 / -1;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 8px 0 -8px;
  padding-top: 8px;
}

#toast-container {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
}

.toast {
  background: var(--toast-bg);
  border: 1px solid var(--card-border);
  border-radius: 6px;
  padding: 10px 16px;
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.4);
  animation: toastIn 0.3s ease-out;
  pointer-events: auto;
  max-width: 360px;
}

.toast.removing { animation: toastOut 0.25s ease-in forwards; }

.toast-icon { font-size: 16px; flex-shrink: 0; }

@keyframes toastIn {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}

@keyframes toastOut {
  from { opacity: 1; transform: translateX(0); }
  to { opacity: 0; transform: translateX(20px); }
}

@media (max-width: 640px) {
  header, .grid { padding-left: 16px; padding-right: 16px; }
  .grid { grid-template-columns: 1fr; }
  #toast-container { left: 16px; right: 16px; }
  .toast { max-width: none; }
}
</style>