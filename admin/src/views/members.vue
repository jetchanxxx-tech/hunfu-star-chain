<template>
  <div>
    <el-card>
      <template #header>
        <el-row justify="space-between" align="middle">
          <span>会员列表 ({{ total }} 人)</span>
          <el-input v-model="search" placeholder="搜索昵称/家庭名" style="width:240px" clearable @input="load" />
        </el-row>
      </template>
      <el-table :data="members" empty-text="暂无会员数据" style="width:100%" v-loading="loading">
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="relation" label="关系" width="80">
          <template #default="{ row }">
            <el-tag size="small">{{ relMap[row.relation] || row.relation }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="60">
          <template #default="{ row }">{{ row.gender === 1 ? '男' : row.gender === 0 ? '女' : '' }}</template>
        </el-table-column>
        <el-table-column prop="birth_date" label="出生日期" width="120" />
        <el-table-column prop="family_name" label="所属家庭" />
        <el-table-column prop="package_name" label="服务包" width="140">
          <template #default="{ row }">
            <el-tag :type="row.package_name === '无' ? 'info' : 'success'" size="small">{{ row.package_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'warning'" size="small">{{ row.status === 'active' ? '活跃' : '非活跃' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const relMap: Record<string, string> = { self: '本人', spouse: '配偶', child: '子女', parent: '父母', other: '其他' }
const members = ref<any[]>([])
const total = ref(0)
const search = ref('')
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const params = search.value ? '?search=' + encodeURIComponent(search.value) : ''
    const res = await fetch('/api/v1/admin/members' + params, { headers: { Authorization: `Bearer ${token}` } })
    if (res.ok) {
      members.value = await res.json()
      total.value = members.value.length
    }
  } catch { /* ignore */ }
  loading.value = false
}

onMounted(() => load())
</script>
