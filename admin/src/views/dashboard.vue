<template>
  <div>
    <el-row :gutter="16">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card class="stat-card">
          <div class="stat-label">{{ c.label }}</div>
          <div class="stat-value">{{ c.value }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="12">
        <el-card title="成员状态分布">
          <div style="padding:20px">
            <el-progress :percentage="activePct" :color="'#67C23A'" :stroke-width="20">
              <span>活跃 {{ stats.active_members }} 人</span>
            </el-progress>
            <el-progress :percentage="inactivePct" :color="'#F56C6C'" :stroke-width="20" style="margin-top:20px">
              <span>非活跃 {{ stats.inactive_members }} 人</span>
            </el-progress>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card title="服务概览">
          <div style="padding:10px">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="家庭总数">{{ stats.total_families }}</el-descriptions-item>
              <el-descriptions-item label="活跃服务包">{{ stats.active_packages }}</el-descriptions-item>
              <el-descriptions-item label="待处理任务">{{ stats.pending_tasks }}</el-descriptions-item>
              <el-descriptions-item label="随访完成率">{{ stats.task_complete_rate.toFixed(1) }}%</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const stats = ref({
  total_members: 0, active_members: 0, inactive_members: 0,
  total_families: 0, new_members_month: 0, active_packages: 0,
  pending_tasks: 0, task_complete_rate: 0
})

const activePct = computed(() => stats.value.total_members ? Math.round(stats.value.active_members / stats.value.total_members * 100) : 0)
const inactivePct = computed(() => stats.value.total_members ? Math.round(stats.value.inactive_members / stats.value.total_members * 100) : 0)

const cards = computed(() => [
  { label: '星球居民总数', value: stats.value.total_members },
  { label: '本月新增', value: stats.value.new_members_month },
  { label: '活跃服务包', value: stats.value.active_packages },
  { label: '随访完成率', value: stats.value.task_complete_rate.toFixed(1) + '%' }
])

onMounted(async () => {
  try {
    const token = localStorage.getItem('admin_token')
    const res = await fetch('/api/v1/admin/dashboard', { headers: { Authorization: `Bearer ${token}` } })
    if (res.ok) stats.value = await res.json()
  } catch { /* ignore */ }
})
</script>

<style scoped>
.stat-card { text-align: center; }
.stat-label { font-size: 13px; color: #999; margin-bottom: 8px; }
.stat-value { font-size: 28px; font-weight: bold; color: #333; }
.stat-trend { font-size: 12px; margin-top: 4px; }
.chart-placeholder { height: 280px; display: flex; align-items: center; justify-content: center; color: #ccc; font-size: 14px; background: #fafafa; border-radius: 4px; }
</style>
