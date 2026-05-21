<template>
  <div>
    <el-card class="mb20">
      <template #header><span>默认节点模板</span></template>
      <el-table :data="templates" stripe>
        <el-table-column prop="node_code" label="节点编码" width="130" />
        <el-table-column prop="node_name" label="节点名称" width="130" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="default_start" label="默认开始" width="120" />
        <el-table-column prop="default_end" label="默认结束" width="130" />
        <el-table-column prop="reminder_days" label="提醒提前(天)" width="120" />
        <el-table-column prop="sort_order" label="排序" width="70" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'enabled' ? 'success' : 'danger'">{{ row.status === 'enabled' ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="editTemplate(row)">编辑</el-button>
            <el-button size="small" :type="row.status==='enabled'?'warning':'success'"
              @click="toggleTemplate(row)">{{ row.status==='enabled'?'禁用':'启用' }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-button type="primary" style="margin-top:12px" @click="editTemplate(null)">新增节点</el-button>
    </el-card>

    <el-card>
      <template #header>
        <span>医院覆盖配置</span>
        <el-input v-model="hospitalCode" placeholder="输入医院编码" style="width:200px;margin-left:16px" />
        <el-button style="margin-left:8px" type="primary" @click="loadOverrides">查询</el-button>
      </template>
      <el-table :data="overrides" stripe>
        <el-table-column prop="node_code" label="节点编码" width="130" />
        <el-table-column prop="start_offset" label="覆盖开始" width="120" />
        <el-table-column prop="end_offset" label="覆盖结束" width="130" />
        <el-table-column prop="reminder_days" label="提醒提前(天)" width="130" />
        <el-table-column prop="is_enabled" label="启用" width="80">
          <template #default="{ row }"><el-tag :type="row.is_enabled?'success':'info'">{{ row.is_enabled?'是':'否' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="editOverride(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteOverride(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-button type="primary" style="margin-top:12px" @click="editOverride(null)">新增覆盖</el-button>
    </el-card>

    <!-- Template Dialog -->
    <el-dialog :title="tplForm.id?'编辑节点模板':'新增节点模板'" v-model="tplVisible" width="500px">
      <el-form :model="tplForm" label-width="110px">
        <el-form-item label="节点编码"><el-input v-model="tplForm.node_code" :disabled="!!tplForm.id" /></el-form-item>
        <el-form-item label="节点名称"><el-input v-model="tplForm.node_name" /></el-form-item>
        <el-form-item label="分类">
          <el-select v-model="tplForm.category">
            <el-option label="产前" value="prenatal" />
            <el-option label="产后" value="postpartum" />
            <el-option label="儿科" value="pediatrics" />
            <el-option label="疫苗" value="vaccine" />
          </el-select>
        </el-form-item>
        <el-form-item label="默认开始"><el-input v-model="tplForm.default_start" /></el-form-item>
        <el-form-item label="默认结束"><el-input v-model="tplForm.default_end" /></el-form-item>
        <el-form-item label="提醒提前天数"><el-input-number v-model="tplForm.reminder_days" :min="0" :max="90" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="tplForm.sort_order" :min="0" /></el-form-item>
        <el-form-item label="备注说明"><el-input v-model="tplForm.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tplVisible=false">取消</el-button>
        <el-button type="primary" @click="saveTemplate">保存</el-button>
      </template>
    </el-dialog>

    <!-- Override Dialog -->
    <el-dialog :title="ovrForm.id?'编辑覆盖配置':'新增覆盖配置'" v-model="ovrVisible" width="500px">
      <el-form :model="ovrForm" label-width="110px">
        <el-form-item label="医院编码"><el-input v-model="ovrForm.hospital_code" :disabled="!!ovrForm.id" /></el-form-item>
        <el-form-item label="节点编码"><el-input v-model="ovrForm.node_code" :disabled="!!ovrForm.id" /></el-form-item>
        <el-form-item label="覆盖开始"><el-input v-model="ovrForm.start_offset" /></el-form-item>
        <el-form-item label="覆盖结束"><el-input v-model="ovrForm.end_offset" /></el-form-item>
        <el-form-item label="提醒提前天数"><el-input-number v-model="ovrForm.reminder_days" :min="0" :max="90" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="ovrForm.is_enabled" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="ovrForm.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ovrVisible=false">取消</el-button>
        <el-button type="primary" @click="saveOverride">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const templates = ref<any[]>([])
const overrides = ref<any[]>([])
const hospitalCode = ref('')
const tplVisible = ref(false)
const ovrVisible = ref(false)
const tplForm = reactive<any>({})
const ovrForm = reactive<any>({})

async function loadTemplates() {
  const { data } = await api.getNodeTemplates()
  templates.value = data
}
async function loadOverrides() {
  if (!hospitalCode.value) return ElMessage.warning('请输入医院编码')
  const { data } = await api.getNodeOverrides(hospitalCode.value)
  overrides.value = data
}
function editTemplate(row: any) {
  Object.assign(tplForm, row ? { ...row } : { id: 0, node_code: '', node_name: '', category: 'prenatal', default_start: '', default_end: '', reminder_days: 7, sort_order: 0, description: '' })
  tplVisible.value = true
}
async function saveTemplate() {
  await api.upsertNodeTemplate(tplForm)
  ElMessage.success('保存成功')
  tplVisible.value = false
  loadTemplates()
}
async function toggleTemplate(row: any) {
  const newStatus = row.status === 'enabled' ? 'disabled' : 'enabled'
  await ElMessageBox.confirm(`确认${newStatus==='enabled'?'启用':'禁用'}该节点？`, '提示')
  await api.updateTemplateStatus(row.node_code, newStatus)
  ElMessage.success('状态已更新')
  loadTemplates()
}
function editOverride(row: any) {
  Object.assign(ovrForm, row ? { ...row } : { id: 0, hospital_code: hospitalCode.value || '', node_code: '', start_offset: '', end_offset: '', reminder_days: null, is_enabled: true, description: '' })
  ovrVisible.value = true
}
async function saveOverride() {
  await api.upsertNodeOverride(ovrForm)
  ElMessage.success('保存成功')
  ovrVisible.value = false
  loadOverrides()
}
async function deleteOverride(row: any) {
  await ElMessageBox.confirm('确认删除该覆盖配置？', '提示', { type: 'warning' })
  await api.deleteNodeOverride(row.hospital_code, row.node_code)
  ElMessage.success('已删除')
  loadOverrides()
}
onMounted(loadTemplates)
</script>
<style scoped>.mb20 { margin-bottom: 20px; }</style>
