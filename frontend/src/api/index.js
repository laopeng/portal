const BASE = ''

export async function login(username, password) {
  const csrf = await fetch(BASE + '/portal/csrf').then(r => r.json())
  const form = new URLSearchParams()
  form.append('username', username)
  form.append('password', password)
  form.append('csrf_token', csrf.token)
  const r = await fetch(BASE + '/portal/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: form,
  })
  if (r.ok) return { success: true }
  const data = await r.json().catch(() => ({}))
  return { success: false, error: data.error || '登录失败' }
}

export async function logout() {
  await fetch(BASE + '/portal/logout', { method: 'POST' })
}

export async function fetchServices() {
  const r = await fetch(BASE + '/portal/services')
  return r.json()
}

export async function fetchHealth() {
  const r = await fetch(BASE + '/portal/health')
  return r.json()
}