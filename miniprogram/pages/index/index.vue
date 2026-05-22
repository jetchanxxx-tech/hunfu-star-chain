<template>
  <view class="home">
    <!-- Top user greeting -->
    <view class="header">
      <text class="greeting">你好，{{ nickname || '星球居民' }}</text>
      <text class="subtitle">今天是你孕期的第 {{ week }} 周</text>
    </view>

    <!-- Quick actions -->
    <view class="actions">
      <view class="action-item" @click="navigate('/pages/timeline/timeline')">
        <view class="action-icon timeline"></view>
        <text>时光轴</text>
      </view>
      <view class="action-item" @click="navigate('/pages/family/family')">
        <view class="action-icon family"></view>
        <text>同心圆</text>
      </view>
      <view class="action-item" @click="navigate('/pages/packages/packages')">
        <view class="action-icon pkg"></view>
        <text>服务包</text>
      </view>
      <view class="action-item" @click="navigate('/pages/ai-chat/ai-chat')">
        <view class="action-icon ai"></view>
        <text>灵犀问答</text>
      </view>
    </view>

    <!-- Timeline preview -->
    <view class="section">
      <view class="section-title">
        <text>关键节点</text>
        <text class="more" @click="navigate('/pages/timeline/timeline')">全部 ></text>
      </view>
      <view v-if="events.length === 0" class="empty">暂无数据</view>
      <view v-for="e in events" :key="e.id" class="event-item">
        <text class="event-dot"></text>
        <text class="event-date">{{ e.event_date }}</text>
        <view class="event-info">
          <text class="event-type">{{ eventLabel(e.event_type) }}</text>
          <text v-if="e.event_data" class="event-desc">{{ e.event_data }}</text>
        </view>
      </view>
    </view>

    <!-- Reports preview -->
    <view class="section">
      <view class="section-title">
        <text>最近报告</text>
        <text class="more" @click="navigate('/pages/reports/reports')">全部 ></text>
      </view>
      <view v-if="reports.length === 0" class="empty">暂无数据</view>
      <view v-for="r in reports" :key="r.id" class="report-item" @click="navigate('/pages/reports/reports')">
        <text class="rpt-type">{{ reportLabel(r.report_type) }}</text>
        <text class="rpt-date">{{ r.report_date }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { week: 0, events: [], reports: [], nickname: '' }
  },
  onShow() {
    this.loadData()
  },
  methods: {
    navigate(url) {
      const tabPages = ['/pages/index/index', '/pages/timeline/timeline', '/pages/packages/packages', '/pages/ai-chat/ai-chat']
      if (tabPages.includes(url)) {
        uni.switchTab({ url })
      } else {
        uni.navigateTo({ url })
      }
    },
    eventLabel(t) {
      const map = { first_prenatal: '首次产检', nt: 'NT检查', early_tang: '早唐筛查', ogtt: '糖耐量', quad_d: '四维彩超',
        delivery: '分娩', '42day': '产后42天复查', vaccine_2m: '2月龄疫苗', vaccine_3m: '3月龄疫苗', vaccine: '疫苗接种' }
      return map[t] || t
    },
    reportLabel(t) {
      const map = { lab: '检验报告', imaging: '检查报告', discharge: '出院小结' }
      return map[t] || t
    },
    loadData() {
      const memberId = uni.getStorageSync('member_id')
      if (memberId) {
        api.getTimeline(memberId, 5).then(events => { this.events = events || [] })
        api.getReports(memberId, 3).then(reports => { this.reports = reports || [] })
      } else {
        this.loadDemo()
      }
    },
    loadDemo() {
      uni.request({
        url: '/api/v1/demo/home',
        success: (res) => {
          const d = res.data
          this.nickname = d.nickname || ''
          this.week = d.week || 0
          this.events = d.events || []
          this.reports = d.reports || []
        }
      })
    }
  }
}
</script>

<style>
.home { padding: 20px 16px; }
.header { background: linear-gradient(135deg, #2E75B6, #1A4F7E); padding: 30px 20px;
  border-radius: 12px; margin-bottom: 20px; }
.greeting { color: #fff; font-size: 20px; font-weight: bold; display: block; }
.subtitle { color: rgba(255,255,255,0.8); font-size: 14px; margin-top: 6px; display: block; }
.actions { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
.action-item { display: flex; flex-direction: column; align-items: center; background: #fff;
  padding: 16px 8px; border-radius: 8px; font-size: 12px; color: #333; }
.action-icon { width: 40px; height: 40px; border-radius: 50%; margin-bottom: 6px; }
.action-icon.timeline { background: #E8F4FD; }
.action-icon.family { background: #FFF0E6; }
.action-icon.pkg { background: #E6F7E6; }
.action-icon.ai { background: #F3E8FF; }
.section { background: #fff; border-radius: 10px; padding: 16px; margin-bottom: 16px; }
.section-title { display: flex; justify-content: space-between; font-size: 16px; font-weight: bold; margin-bottom: 12px; }
.more { font-size: 12px; color: #999; font-weight: normal; }
.empty { color: #ccc; text-align: center; padding: 20px; font-size: 14px; }
.event-item { display: flex; align-items: center; padding: 8px 0; border-bottom: 1px solid #f0f0f0; }
.event-dot { width: 8px; height: 8px; border-radius: 50%; background: #2E75B6; margin-right: 10px; }
.event-date { font-size: 13px; color: #999; margin-right: 10px; min-width: 85px; }
.event-info { display: flex; flex-direction: column; }
.event-type { font-size: 14px; color: #333; }
.event-desc { font-size: 11px; color: #2E75B6; margin-top: 2px; }
.report-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f0f0f0; }
.rpt-type { font-size: 14px; color: #333; }
.rpt-date { font-size: 13px; color: #999; }
</style>
