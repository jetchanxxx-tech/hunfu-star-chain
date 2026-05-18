<template>
  <div>
    <el-row :gutter="16">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card class="stat-card">
          <div class="stat-label">{{ c.label }}</div>
          <div class="stat-value">{{ c.value }}</div>
          <div class="stat-trend" :style="{ color: c.trend > 0 ? '#67C23A' : '#F56C6C' }">
            {{ c.trend > 0 ? '↑' : '↓' }} {{ Math.abs(c.trend) }}%
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="12">
        <el-card title="会员增长趋势">
          <div class="chart-placeholder">会员增长趋势图（预留 ECharts）</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card title="服务包销售排行">
          <div class="chart-placeholder">销售排行榜（预留 ECharts）</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="12">
        <el-card title="随访完成率">
          <div class="chart-placeholder">随访完成率趋势（预留 ECharts）</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card title="转化漏斗">
          <div class="chart-placeholder">注册→购买→续费转化漏斗（预留 ECharts）</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const cards = ref([
  { label: '星球居民总数', value: '--', trend: 0 },
  { label: '本月新增家庭', value: '--', trend: 0 },
  { label: '活跃服务包', value: '--', trend: 0 },
  { label: '随访完成率', value: '--%', trend: 0 }
])

onMounted(async () => {
  // Fetch real data when API ready
})
</script>

<style scoped>
.stat-card { text-align: center; }
.stat-label { font-size: 13px; color: #999; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: bold; color: #333; }
.stat-trend { font-size: 12px; margin-top: 4px; }
.chart-placeholder { height: 280px; display: flex; align-items: center; justify-content: center; color: #ccc; font-size: 14px; background: #fafafa; border-radius: 4px; }
</style>
