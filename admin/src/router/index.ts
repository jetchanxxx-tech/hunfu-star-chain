import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard.vue'), meta: { title: '数据驾驶舱' } },
      { path: 'members', name: 'Members', component: () => import('@/views/members.vue'), meta: { title: '会员管理' } },
      { path: 'packages', name: 'Packages', component: () => import('@/views/packages-admin.vue'), meta: { title: '服务包配置' } },
      { path: 'followup', name: 'Followup', component: () => import('@/views/followup.vue'), meta: { title: '随访规则' } },
      { path: 'tasks', name: 'Tasks', component: () => import('@/views/tasks.vue'), meta: { title: '任务管理' } },
      { path: 'timeline-config', name: 'TimelineConfig', component: () => import('@/views/timeline-config.vue'), meta: { title: '时间轴配置' } },
      { path: 'verification', name: 'Verification', component: () => import('@/views/verification.vue'), meta: { title: '核销记录' } },
      { path: 'auth-audit', name: 'AuthAudit', component: () => import('@/views/auth-audit.vue'), meta: { title: '授权审计' } },
      { path: 'system', name: 'System', component: () => import('@/views/system.vue'), meta: { title: '系统管理' } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('admin_token')
  if (!token && !to.meta.guest) {
    next('/login')
  } else {
    next()
  }
})

export default router
