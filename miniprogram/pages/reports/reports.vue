<template>
  <view class="reports-page">
    <view class="nav-back" @click="goBack">
      <text class="back-icon">←</text>
      <text>返回</text>
    </view>
    <view class="tabs">
      <text v-for="t in tabs" :key="t" :class="['tab', { active: activeTab === t }]"
        @click="switchTab(t)">{{ t }}</text>
    </view>
    <view v-for="r in filteredReports" :key="r.id" class="report-card" @click="toggleDetail(r)">
      <view class="report-header">
        <text class="report-type">{{ typeLabel(r.report_type) }}</text>
        <text class="report-date">{{ r.report_date }}</text>
        <text class="expand-icon">{{ expandedId === r.id ? '▲' : '▼' }}</text>
      </view>
      <text v-if="r.abnormal_flags" class="abnormal-badge">有异常指标</text>
      <view v-if="expandedId === r.id" class="report-detail">
        <text class="detail-text">报告编号: {{ r.id }}</text>
        <text v-if="r.summary" class="detail-text">摘要: {{ r.summary }}</text>
        <text v-if="r.report_no" class="detail-text">报告号: {{ r.report_no }}</text>
        <text v-if="r.hospital_code" class="detail-text">医院: {{ r.hospital_code }}</text>
      </view>
    </view>
    <view v-if="reports.length === 0" class="empty">暂无报告</view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { reports: [], activeTab: '全部', expandedId: null,
      tabs: ['全部', '检验', '检查', '出院小结'] }
  },
  onShow() { this.load() },
  methods: {
    goBack() { uni.navigateBack() },
    load() {
      const memberId = uni.getStorageSync('member_id')
      if (!memberId) return
      api.getReports(memberId, 50).then(reports => { this.reports = reports || [] })
    },
    typeLabel(t) {
      const map = { lab: '检验报告', imaging: '检查报告', discharge: '出院小结' }
      return map[t] || t
    },
    switchTab(t) {
      this.activeTab = t
      this.expandedId = null
    },
    toggleDetail(r) {
      this.expandedId = this.expandedId === r.id ? null : r.id
    }
  },
  computed: {
    filteredReports() {
      if (this.activeTab === '全部') return this.reports
      const typeMap = { '检验': 'lab', '检查': 'imaging', '出院小结': 'discharge' }
      const filterType = typeMap[this.activeTab]
      return this.reports.filter(r => r.report_type === filterType)
    }
  }
}
</script>

<style>
.reports-page { padding: 16px; }
.nav-back { display: flex; align-items: center; gap: 4px; margin-bottom: 12px;
  font-size: 14px; color: #2E75B6; cursor: pointer; }
.back-icon { font-size: 18px; font-weight: bold; }
.tabs { display: flex; gap: 10px; margin-bottom: 16px; flex-wrap: wrap; }
.tab { padding: 6px 14px; border-radius: 14px; font-size: 13px; background: #fff; color: #666; }
.tab.active { background: #2E75B6; color: #fff; }
.report-card { background: #fff; padding: 14px; border-radius: 8px; margin-bottom: 10px; }
.report-header { display: flex; justify-content: space-between; align-items: center; }
.report-type { font-size: 15px; color: #333; font-weight: 500; }
.report-date { font-size: 13px; color: #999; }
.expand-icon { font-size: 12px; color: #bbb; margin-left: 8px; }
.abnormal-badge { color: #E53E3E; font-size: 12px; margin-top: 4px; display: inline-block;
  background: #FFF5F5; padding: 2px 8px; border-radius: 4px; }
.report-detail { margin-top: 10px; padding-top: 10px; border-top: 1px solid #f0f0f0; }
.detail-text { display: block; font-size: 13px; color: #666; line-height: 1.8; }
.empty { text-align: center; color: #ccc; padding: 80px 0; }
</style>
