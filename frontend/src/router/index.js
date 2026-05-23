import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'

const routes = [
  { path: '/login', name: 'Login', component: Login },
  { path: '/', name: 'Dashboard', component: Dashboard },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from) => {
  if (to.name === 'Login') return
  try {
    const r = await fetch('/portal/auth/status')
    const d = await r.json()
    if (!d.authenticated) return '/login'
  } catch {
    return '/login'
  }
})

export default router