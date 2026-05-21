// #ifdef H5
const BASE_URL = '/api/v1'
// #endif
// #ifdef MP-WEIXIN
const BASE_URL = 'https://huifu.pangu-cloud.com/api/v1'
// #endif

function request(method, path, data) {
  const token = uni.getStorageSync('token')
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + path,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success(res) {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
        } else {
          reject(res.data)
        }
      },
      fail(err) {
        reject(err)
      }
    })
  })
}

export const api = {
  // Auth
  wxLogin(code) {
    return request('POST', '/auth/wx-login', { code })
  },
  bindWechat(memberId, code) {
    return request('POST', '/auth/bind-wechat', { member_id: memberId, code })
  },

  // Family
  createFamily(data) {
    return request('POST', '/families', data)
  },
  getFamily(id) {
    return request('GET', `/families/${id}`)
  },
  addMember(familyId, data) {
    return request('POST', `/families/${familyId}/members`, data)
  },

  // Timeline
  getTimeline(memberId, limit = 50) {
    return request('GET', `/members/${memberId}/timeline?limit=${limit}`)
  },

  // Reports
  getReports(memberId, limit = 20, offset = 0) {
    return request('GET', `/members/${memberId}/reports?limit=${limit}&offset=${offset}`)
  },

  // Packages
  getPackages() {
    return request('GET', '/packages')
  },

  // AI
  chat(message, sessionId) {
    return request('POST', '/ai/chat', { message, session_id: sessionId })
  },
  searchFAQ(query) {
    return request('GET', `/ai/faq?q=${encodeURIComponent(query)}`)
  }
}
