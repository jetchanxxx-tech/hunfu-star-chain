<template>
  <div>
    <el-card>
      <template #header>
        <el-row justify="space-between" align="middle">
          <span>会员列表</span>
          <el-input v-model="search" placeholder="搜索昵称/手机号" style="width:240px" clearable />
        </el-row>
      </template>
      <el-table :data="[]" empty-text="暂无会员数据" style="width:100%">
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="phone" label="手机号" />
        <el-table-column prop="family" label="所属家庭" />
        <el-table-column prop="package" label="服务包" />
        <el-table-column prop="status" label="状态" />
        <el-table-column label="操作" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const search = ref('')
</script>
