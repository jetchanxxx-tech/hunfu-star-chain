<template>
  <view class="login-page">
    <view class="brand">
      <text class="brand-name">惠福星链</text>
      <text class="brand-desc">全病程健康协同平台</text>
    </view>
    <button class="wx-btn" @click="handleLogin">微信一键登录</button>
    <text class="tip">登录即同意《用户协议》与《隐私政策》</text>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { loading: false }
  },
  methods: {
    handleLogin() {
      this.loading = true
      uni.login({
        provider: 'weixin',
        success: (loginRes) => {
          api.wxLogin(loginRes.code).then((res) => {
            uni.setStorageSync('token', res.token)
            if (res.is_new_user) {
              uni.navigateTo({ url: '/pages/family/family?action=register' })
            } else {
              uni.switchTab({ url: '/pages/index/index' })
            }
          }).catch((err) => {
            uni.showToast({ title: '登录失败: ' + (err.error || err.errMsg), icon: 'none' })
          }).finally(() => { this.loading = false })
        },
        fail: () => {
          this.loading = false
          uni.showToast({ title: '微信授权失败', icon: 'none' })
        }
      })
    }
  }
}
</script>

<style>
.login-page {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  min-height: 100vh; padding: 60px 30px;
}
.brand { text-align: center; margin-bottom: 60px; }
.brand-name { font-size: 32px; font-weight: bold; color: #2E75B6; }
.brand-desc { font-size: 14px; color: #999; margin-top: 8px; display: block; }
.wx-btn { width: 280px; height: 48px; line-height: 48px; background: #07C160; color: #fff;
  border-radius: 24px; font-size: 16px; border: none; }
.tip { font-size: 12px; color: #bbb; margin-top: 20px; }
</style>
