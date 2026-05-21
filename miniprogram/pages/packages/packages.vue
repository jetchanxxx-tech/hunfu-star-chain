<template>
  <view class="pkg-page">
    <view v-for="p in packages" :key="p.package_uuid" class="pkg-card">
      <view class="pkg-header">
        <text class="pkg-name">{{ p.name }}</text>
        <text class="pkg-level" :class="p.level">{{ p.level }}</text>
      </view>
      <text class="pkg-desc">{{ p.description }}</text>
      <view class="pkg-footer">
        <text class="pkg-price">¥{{ p.price }}</text>
        <button class="buy-btn">立即购买</button>
      </view>
    </view>
    <view v-if="packages.length === 0" class="empty">暂无可用服务包</view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { packages: [] }
  },
  onShow() { this.load() },
  methods: {
    load() {
      const memberId = uni.getStorageSync('member_id')
      if (memberId) {
        api.getPackages().then(res => { this.packages = res || [] })
      } else {
        // 未登录：从 demo 接口加载服务包
        uni.request({
          url: '/api/v1/demo/home',
          success: (res) => {
            this.packages = (res.data && res.data.packages) ? res.data.packages : []
          }
        })
      }
    }
  }
}
</script>

<style>
.pkg-page { padding: 16px; }
.pkg-card { background: #fff; border-radius: 10px; padding: 16px; margin-bottom: 14px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05); }
.pkg-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.pkg-name { font-size: 17px; font-weight: bold; color: #333; }
.pkg-level { font-size: 12px; padding: 2px 10px; border-radius: 10px; }
.pkg-level.VIP { background: #E8F4FD; color: #2E75B6; }
.pkg-level.VVIP { background: #FFF0E6; color: #E67E22; }
.pkg-desc { font-size: 13px; color: #666; line-height: 1.5; display: block; margin-bottom: 12px; }
.pkg-footer { display: flex; justify-content: space-between; align-items: center; }
.pkg-price { font-size: 20px; font-weight: bold; color: #E53E3E; }
.buy-btn { background: #2E75B6; color: #fff; border: none; border-radius: 16px;
  padding: 6px 20px; font-size: 14px; }
.empty { text-align: center; color: #ccc; padding: 80px 0; }
</style>
