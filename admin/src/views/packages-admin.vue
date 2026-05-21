<template>
  <div>
    <el-card>
      <template #header>
        <el-row justify="space-between" align="middle">
          <span>服务包配置</span>
          <el-button type="primary" @click="openDialog()">新增服务包</el-button>
        </el-row>
      </template>
      <el-table :data="packages" empty-text="暂无服务包" style="width:100%" v-loading="loading">
        <el-table-column prop="name" label="名称" min-width="160" />
        <el-table-column prop="level" label="等级" width="80">
          <template #default="{ row }">
            <el-tag :type="row.level === 'VVIP' ? 'warning' : 'primary'" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">¥{{ row.price.toLocaleString() }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">{{ row.status === 'online' ? '上线' : row.status === 'draft' ? '草稿' : '下线' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDialog(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑服务包' : '新增服务包'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="等级">
          <el-select v-model="form.level">
            <el-option value="VIP" label="VIP" />
            <el-option value="VVIP" label="VVIP" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格"><el-input-number v-model="form.price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status">
            <el-option value="draft" label="草稿" />
            <el-option value="online" label="上线" />
            <el-option value="offline" label="下线" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="save" :loading="saving">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const packages = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({ name: '', level: 'VIP', price: 0, description: '', status: 'draft' })

function openDialog(row?: any) {
  if (row) {
    editingId.value = row.id
    form.name = row.name; form.level = row.level; form.price = row.price
    form.description = row.description; form.status = row.status
  } else {
    editingId.value = null
    form.name = ''; form.level = 'VIP'; form.price = 0; form.description = ''; form.status = 'draft'
  }
  dialogVisible.value = true
}

async function save() {
  saving.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const url = editingId.value ? `/api/v1/admin/packages/${editingId.value}` : '/api/v1/admin/packages'
    const method = editingId.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(form)
    })
    if (res.ok) {
      ElMessage.success(editingId.value ? '已更新' : '已创建')
      dialogVisible.value = false
      load()
    } else {
      const err = await res.json()
      ElMessage.error(err.error || '保存失败')
    }
  } catch { ElMessage.error('保存失败') }
  saving.value = false
}

async function load() {
  loading.value = true
  try {
    const token = localStorage.getItem('admin_token')
    const res = await fetch('/api/v1/admin/packages', { headers: { Authorization: `Bearer ${token}` } })
    if (res.ok) packages.value = await res.json()
  } catch { /* ignore */ }
  loading.value = false
}

onMounted(() => load())
</script>
