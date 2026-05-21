<template>
  <div>
    <el-card>
      <template #header><span>管理员用户</span></template>
      <el-table :data="users" empty-text="暂无用户" style="width:100%" v-loading="loading">
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="row.role === 'super_admin' ? 'danger' : row.role === 'doctor' ? 'primary' : 'success'" size="small">{{ roleMap[row.role] || row.role }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="real_name" label="姓名" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'warning'" size="small">{{ row.status === 'active' ? '正常' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card style="margin-top:16px">
      <template #header><span>接口监控</span></template>
      <el-empty description="接口调用日志（P1阶段交付）" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const roleMap: Record<string,string> = { super_admin:'超级管理员', steward:'健康管家', doctor:'医生', operator:'运营专员' }
const users = ref<any[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const res = await fetch('/api/v1/admin/users', { headers: { Authorization: `Bearer ${token}` } })
    if (res.ok) users.value = await res.json()
  } catch { /* ignore */ }
  loading.value = false
})
</script>
