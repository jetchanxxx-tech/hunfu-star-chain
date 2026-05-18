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
      window.location.href = '/login'
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
  // Followup
  getFollowupRules() {
    return http.get('/admin/followup/rules')
  },
  // System
  getUsers() {
    return http.get('/admin/users')
  }
}
