<template>
  <el-card>
    <template #header><span>授权操作审计日志</span></template>
    <el-table :data="logs" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="action" label="操作类型" width="100">
        <template #default="{ row }">
          <el-tag :type="actionColor(row.action)">{{ actionText(row.action) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="actor_id" label="操作人ID" width="100" />
      <el-table-column prop="target_id" label="目标ID" width="100" />
      <el-table-column prop="detail" label="详情" min-width="200">
        <template #default="{ row }">{{ row.detail?.String || row.detail }}</template>
      </el-table-column>
      <el-table-column prop="created_at" label="操作时间" width="170" />
    </el-table>
    <el-pagination layout="prev,next" :total="100" :page-size="20" @current-change="onPage" style="margin-top:12px" />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api'
const logs = ref<any[]>([])
let page = 1

function actionText(a: string) {
  return { request: '发起请求', approve: '同意授权', reject: '拒绝授权', revoke: '撤销授权', expire: '到期失效' }[a] || a
}
function actionColor(a: string) {
  return { request: 'info', approve: 'success', reject: 'danger', revoke: 'warning', expire: 'info' }[a] || ''
}
async function loadLogs() {
  const { data } = await api.getAuthAuditLogs({ offset: (page-1)*20, limit: 20 })
  logs.value = data || []
}
function onPage(p: number) { page = p; loadLogs() }
onMounted(loadLogs)
</script>
