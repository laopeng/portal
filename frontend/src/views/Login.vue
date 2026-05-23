<template>
  <div class="login-page">
    <div class="login-box">
      <h1><span class="mark">■</span> Portal</h1>
      <p class="subtitle">本地服务门户</p>

      <div v-if="alertMsg" :class="['alert', alertType === 'error' ? 'alert-error' : 'alert-info']">
        {{ alertMsg }}
      </div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名</label>
          <input id="username" v-model.trim="username" type="text" required autocomplete="username" />
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="password" type="password" required autocomplete="current-password" minlength="4" />
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '登录中...' : '登 录' }}
        </button>
      </form>

      <p class="hint">首次使用请运行 <code>./portal user add &lt;用户名&gt;</code></p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { login } from '../api'

const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const loading = ref(false)
const alertMsg = ref('')
const alertType = ref('info')

onMounted(() => {
  const reason = route.query.reason
  if (reason === 'expired') {
    alertMsg.value = '会话已过期，请重新登录。'
    alertType.value = 'info'
  } else if (reason) {
    alertMsg.value = '请先登录。'
    alertType.value = 'error'
  }
})

async function handleLogin() {
  loading.value = true
  alertMsg.value = ''
  const res = await login(username.value, password.value)
  if (res.success) {
    router.push('/')
  } else {
    alertMsg.value = res.error || '登录失败'
    alertType.value = 'error'
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0d1117;
  color: #e6edf3;
  font-family: system-ui, -apple-system, sans-serif;
}

.login-box {
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 10px;
  padding: 40px;
  width: min(400px, calc(100% - 32px));
}

h1 {
  font-family: "JetBrains Mono", "SF Mono", monospace;
  font-size: 20px;
  text-align: center;
  margin-bottom: 4px;
}

.mark { color: #3fb950; }

.subtitle {
  font-family: "JetBrains Mono", "SF Mono", monospace;
  font-size: 12px;
  color: #8b949e;
  text-align: center;
  margin-bottom: 28px;
}

.alert {
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 13px;
  margin-bottom: 18px;
}

.alert-error {
  background: rgba(248,81,73,0.12);
  border: 1px solid #f85149;
  color: #f85149;
}

.alert-info {
  background: rgba(210,153,34,0.12);
  border: 1px solid #d29922;
  color: #d29922;
}

.form-group { margin-bottom: 16px; }

label {
  display: block;
  font-size: 13px;
  color: #8b949e;
  margin-bottom: 6px;
}

input {
  width: 100%;
  padding: 10px 12px;
  font-size: 14px;
  background: #0d1117;
  border: 1px solid #30363d;
  border-radius: 6px;
  color: #e6edf3;
  font-family: system-ui, -apple-system, sans-serif;
  transition: border-color 0.15s;
}

input:focus {
  outline: none;
  border-color: #3fb950;
}

.btn {
  width: 100%;
  padding: 12px;
  font-size: 14px;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: opacity 0.15s;
  min-height: 44px;
}

.btn:disabled { opacity: 0.6; cursor: not-allowed; }

.btn-primary {
  background: #238636;
  color: #fff;
}

.btn-primary:hover:not(:disabled) { background: #2ea043; }

.hint {
  text-align: center;
  font-size: 12px;
  color: #8b949e;
  margin-top: 20px;
}

.hint code {
  font-family: "JetBrains Mono", "SF Mono", monospace;
  font-size: 12px;
  background: #0d1117;
  padding: 2px 6px;
  border-radius: 4px;
}

@media (max-width: 640px) {
  .login-box { padding: 28px 20px; }
}
</style>