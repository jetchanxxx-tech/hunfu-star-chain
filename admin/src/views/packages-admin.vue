<template>
  <div>
    <el-card>
      <template #header>
        <el-row justify="space-between" align="middle">
          <span>服务包配置</span>
          <el-button type="primary" @click="dialogVisible = true">新增服务包</el-button>
        </el-row>
      </template>
      <el-table :data="[]" empty-text="暂无服务包" style="width:100%">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="level" label="等级" />
        <el-table-column prop="price" label="价格" />
        <el-table-column prop="status" label="状态" />
        <el-table-column label="操作" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增服务包" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="等级">
          <el-select v-model="form.level">
            <el-option value="VIP" label="VIP" />
            <el-option value="VVIP" label="VVIP" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格"><el-input-number v-model="form.price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
const dialogVisible = ref(false)
const form = reactive({ name: '', level: 'VIP', price: 0, description: '' })
</script>
