import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('admin_token')
      window.location.href = '/admin/login'
    }
    return Promise.reject(err)
  }
)

export const api = {
  login(data: { username: string; password: string }) {
    return http.post('/admin/login', data)
  },
  // Dashboard
  getDashboard() {
    return http.get('/admin/dashboard')
  },
  // Members
  getMembers(params: any) {
    return http.get('/admin/members', { params })
  },
  getMemberDetail(id: string) {
    return http.get(`/admin/members/${id}`)
  },
  // Service Packages
  getPackages() {
    return http.get('/admin/packages')
  },
  createPackage(data: any) {
    return http.post('/admin/packages', data)
  },
  updatePackage(id: string, data: any) {
    return http.put(`/admin/packages/${id}`, data)
  },
  // Timeline Node Config
  getNodeTemplates() {
    return http.get('/admin/node-templates')
  },
  upsertNodeTemplate(data: any) {
    return http.post('/admin/node-templates', data)
  },
  updateTemplateStatus(code: string, status: string) {
    return http.patch(`/admin/node-templates/${code}/status`, { status })
  },
  getNodeOverrides(hospitalCode: string) {
    return http.get(`/admin/node-overrides/${hospitalCode}`)
  },
  upsertNodeOverride(data: any) {
    return http.post('/admin/node-overrides', data)
  },
  deleteNodeOverride(hospitalCode: string, nodeCode: string) {
    return http.delete(`/admin/node-overrides/${hospitalCode}?node_code=${nodeCode}`)
  },
  // Followup Tasks
  getTasks(params: any) {
    return http.get('/tasks', { params })
  },
  getTask(id: number) {
    return http.get(`/tasks/${id}`)
  },
  createTask(data: any) {
    return http.post('/tasks', data)
  },
  assignTask(id: number, stewardId: number) {
    return http.patch(`/tasks/${id}/assign`, { steward_id: stewardId })
  },
  completeTask(id: number, notes: string) {
    return http.patch(`/tasks/${id}/complete`, { notes })
  },
  cancelTask(id: number, reason: string) {
    return http.patch(`/tasks/${id}/cancel`, { reason })
  },
  getTaskStats() {
    return http.get('/admin/task-stats')
  },
  // Verification
  getVerificationRecords(params: any) {
    return http.get('/admin/verification-records', { params })
  },
  // Authorization Audit
  getAuthAuditLogs(params: any) {
    return http.get('/admin/authorization-logs', { params })
  },
  // Followup (existing)
  getFollowupRules() {
    return http.get('/admin/followup/rules')
  },
  // System
  getUsers() {
    return http.get('/admin/users')
  }
}
