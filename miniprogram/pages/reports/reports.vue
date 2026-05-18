<template>
  <view class="reports-page">
    <view class="tabs">
      <text v-for="t in tabs" :key="t" :class="['tab', { active: activeTab === t }]"
        @click="activeTab = t">{{ t }}</text>
    </view>
    <view v-for="r in reports" :key="r.id" class="report-card" @click="viewReport(r)">
      <view class="report-header">
        <text class="report-type">{{ typeLabel(r.report_type) }}</text>
        <text class="report-date">{{ r.report_date }}</text>
      </view>
      <text v-if="r.abnormal_flags" class="abnormal-badge">有异常指标</text>
    </view>
    <view v-if="reports.length === 0" class="empty">暂无报告</view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { reports: [], activeTab: '全部', tabs: ['全部', '检验', '检查', '出院小结'] }
  },
  onShow() { this.load() },
  methods: {
    load() {
      const memberId = uni.getStorageSync('member_id')
      if (!memberId) return
      api.getReports(memberId, 50).then(reports => { this.reports = reports || [] })
    },
    typeLabel(t) {
      const map = { lab: '检验报告', imaging: '检查报告', discharge: '出院小结' }
      return map[t] || t
    },
    viewReport(r) {
      uni.showToast({ title: '查看报告 #' + r.id, icon: 'none' })
    }
  }
}
</script>

<style>
.reports-page { padding: 16px; }
.tabs { display: flex; gap: 10px; margin-bottom: 16px; }
.tab { padding: 6px 14px; border-radius: 14px; font-size: 13px; background: #fff; color: #666; }
.tab.active { background: #2E75B6; color: #fff; }
.report-card { background: #fff; padding: 14px; border-radius: 8px; margin-bottom: 10px; }
.report-header { display: flex; justify-content: space-between; align-items: center; }
.report-type { font-size: 15px; color: #333; font-weight: 500; }
.report-date { font-size: 13px; color: #999; }
.abnormal-badge { color: #E53E3E; font-size: 12px; margin-top: 4px; display: inline-block;
  background: #FFF5F5; padding: 2px 8px; border-radius: 4px; }
.empty { text-align: center; color: #ccc; padding: 80px 0; }
</style>
