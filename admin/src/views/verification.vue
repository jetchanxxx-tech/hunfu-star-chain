<template>
  <el-card>
    <template #header><span>核销记录</span></template>
    <el-form :inline="true">
      <el-form-item label="会员ID"><el-input-number v-model="filter.member_id" :min="1" /></el-form-item>
      <el-form-item><el-button type="primary" @click="loadRecords">查询</el-button></el-form-item>
    </el-form>
    <el-table :data="records" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="member_id" label="会员ID" width="80" />
      <el-table-column prop="benefit_type" label="权益类型" width="130" />
      <el-table-column prop="verify_count" label="核销次数" width="90" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status==='success'?'success':row.status==='failed'?'danger':'primary'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="verified_at" label="核销时间" width="170" />
      <el-table-column prop="qr_nonce" label="Nonce" width="160" />
      <el-table-column prop="fail_reason" label="失败原因" min-width="150" />
    </el-table>
    <el-pagination layout="prev,next" :total="100" :page-size="20" @current-change="onPage" style="margin-top:12px" />
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '@/api'
const records = ref<any[]>([])
const filter = reactive({ member_id: '' })
let page = 1

async function loadRecords() {
  const params: any = { offset: (page-1)*20, limit: 20 }
  if (filter.member_id) params.member_id = filter.member_id
  const { data } = await api.getVerificationRecords(params)
  records.value = data || []
}
function onPage(p: number) { page = p; loadRecords() }
onMounted(loadRecords)
</script>
