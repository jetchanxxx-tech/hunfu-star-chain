<template>
  <view class="timeline-page">
    <view class="member-picker" v-if="members.length > 1" @click="showPicker = !showPicker">
      <text>当前: {{ currentMember?.nickname || '本人' }}</text>
      <text class="arrow">▼</text>
    </view>
    <view class="timeline">
      <view v-for="e in events" :key="e.id" class="tl-item">
        <view class="tl-dot"></view>
        <view class="tl-line"></view>
        <view class="tl-content">
          <text class="tl-date">{{ e.event_date }}</text>
          <text class="tl-type">{{ eventLabel(e.event_type) }}</text>
          <text v-if="e.event_data" class="tl-meta">{{ e.event_data }}</text>
        </view>
      </view>
      <view v-if="events.length === 0" class="empty">暂无时间轴数据</view>
    </view>
  </view>
</template>

<script>
import { api } from '@/api/index.js'

export default {
  data() {
    return { events: [], members: [], currentMember: null, showPicker: false }
  },
  onShow() { this.load() },
  methods: {
    load() {
      const memberId = uni.getStorageSync('member_id')
      if (!memberId) return
      api.getTimeline(memberId, 50).then(events => { this.events = events || [] })
      api.getFamily(uni.getStorageSync('family_id')).then(res => {
        this.members = res.members || []
        this.currentMember = this.members.find(m => m.member_uuid === memberId)
      })
    },
    eventLabel(t) {
      const map = { first_prenatal: '首次产检', nt: 'NT检查', ogtt: '糖耐量检测',
        delivery: '分娩', '42day': '42天复查', vaccine: '疫苗接种' }
      return map[t] || t
    }
  }
}
</script>

<style>
.timeline-page { padding: 16px; }
.member-picker { background: #fff; padding: 12px 16px; border-radius: 8px;
  display: flex; justify-content: space-between; margin-bottom: 16px; font-size: 14px; }
.arrow { color: #999; }
.timeline { padding-left: 30px; position: relative; }
.timeline::before { content: ''; position: absolute; left: 14px; top: 0; bottom: 0;
  width: 2px; background: #E5E5E5; }
.tl-item { position: relative; padding-bottom: 24px; }
.tl-dot { position: absolute; left: -18px; top: 6px; width: 10px; height: 10px;
  border-radius: 50%; background: #2E75B6; border: 2px solid #fff; }
.tl-content { display: flex; flex-direction: column; }
.tl-date { font-size: 13px; color: #999; }
.tl-type { font-size: 15px; color: #333; margin-top: 2px; }
.tl-meta { font-size: 12px; color: #2E75B6; margin-top: 2px; }
.empty { text-align: center; color: #ccc; padding: 60px 0; }
</style>
