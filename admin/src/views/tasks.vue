<template>
  <div>
    <!-- Stats Cards -->
    <el-row :gutter="16" class="mb20">
      <el-col :span="5"><el-card><div class="stat-label">待处理</div><div class="stat-val">{{ stats.pending || 0 }}</div></el-card></el-col>
      <el-col :span="5"><el-card><div class="stat-label">进行中</div><div class="stat-val">{{ stats.in_progress || 0 }}</div></el-card></el-col>
      <el-col :span="5"><el-card><div class="stat-label">今日完成</div><div class="stat-val">{{ stats.completed_today || 0 }}</div></el-card></el-col>
      <el-col :span="5"><el-card><div class="stat-label">超时未处理</div><div class="stat-val" style="color:#f56c6c">{{ stats.overdue || 0 }}</div></el-card></el-col>
    </el-row>

    <el-card>
      <template #header>
        <span>任务列表</span>
        <el-select v-model="filter.status" placeholder="状态筛选" clearable style="width:140px;margin-left:16px" @change="loadTasks">
          <el-option label="待处理" value="pending" /><el-option label="进行中" value="in_progress" />
          <el-option label="已完成" value="completed" /><el-option label="已取消" value="cancelled" />
        </el-select>
        <el-select v-model="filter.trigger_type" placeholder="触发类型" clearable style="width:140px;margin-left:8px" @change="loadTasks">
          <el-option label="孕周提醒" value="gestation_week" /><el-option label="月龄提醒" value="age_month" />
          <el-option label="检验异常" value="lab_abnormal" /><el-option label="未到诊" value="no_show" />
          <el-option label="咨询升级" value="consult_escalation" /><el-option label="投诉处理" value="complaint" />
        </el-select>
        <el-button style="margin-left:8px" @click="loadTasks">查询</el-button>
        <el-button style="float:right" type="primary" @click="createVisible=true">创建任务</el-button>
      </template>

      <el-table :data="tasks" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="任务标题" min-width="200" />
        <el-table-column prop="trigger_type" label="触发类型" width="110" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusColor(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="assigned_to" label="指派人" width="80" />
        <el-table-column prop="due_date" label="截止日期" width="110" />
        <el-table-column label="操作" width="260">
          <template #default="{ row }">
            <el-button v-if="row.status==='pending'" size="small" type="primary" @click="assignTask(row)">指派</el-button>
            <el-button v-if="row.status==='in_progress'" size="small" type="success" @click="completeTask(row)">完成</el-button>
            <el-button v-if="row.status==='pending'||row.status==='in_progress'" size="small" type="warning" @click="cancelTask(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination layout="prev,next" :total="100" :page-size="20" @current-change="onPage" style="margin-top:12px" />
    </el-card>

    <!-- Assign Dialog -->
    <el-dialog title="指派任务" v-model="assignVisible" width="350px">
      <el-form><el-form-item label="管家ID"><el-input-number v-model="stewardId" :min="1" /></el-form-item></el-form>
      <template #footer>
        <el-button @click="assignVisible=false">取消</el-button>
        <el-button type="primary" @click="doAssign">确认指派</el-button>
      </template>
    </el-dialog>

    <!-- Complete Dialog -->
    <el-dialog title="完成任务" v-model="completeVisible" width="400px">
      <el-form><el-form-item label="备注"><el-input v-model="completeNotes" type="textarea" /></el-form-item></el-form>
      <template #footer>
        <el-button @click="completeVisible=false">取消</el-button>
        <el-button type="primary" @click="doComplete">确认完成</el-button>
      </template>
    </el-dialog>

    <!-- Cancel Dialog -->
    <el-dialog title="取消任务" v-model="cancelVisible" width="400px">
      <el-form><el-form-item label="取消原因"><el-input v-model="cancelReason" type="textarea" /></el-form-item></el-form>
      <template #footer>
        <el-button @click="cancelVisible=false">取消</el-button>
        <el-button type="danger" @click="doCancel">确认取消</el-button>
      </template>
    </el-dialog>

    <!-- Create Dialog -->
    <el-dialog title="创建任务" v-model="createVisible" width="450px">
      <el-form :model="createForm" label-width="90px">
        <el-form-item label="标题"><el-input v-model="createForm.title" /></el-form-item>
        <el-form-item label="会员ID"><el-input-number v-model="createForm.member_id" :min="1" /></el-form-item>
        <el-form-item label="触发类型">
          <el-select v-model="createForm.trigger_type">
            <el-option label="孕周提醒" value="gestation_week" /><el-option label="投诉处理" value="complaint" />
            <el-option label="活动执行" value="event_execution" /><el-option label="精准营销" value="marketing" />
          </el-select>
        </el-form-item>
        <el-form-item label="截止日期"><el-date-picker v-model="createForm.due_date" type="date" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible=false">取消</el-button>
        <el-button type="primary" @click="doCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { api } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const tasks = ref<any[]>([])
const stats = ref<any>({})
const filter = reactive({ status: '', trigger_type: '' })
let page = 1
let currentTaskId = 0
const assignVisible = ref(false), stewardId = ref(1)
const completeVisible = ref(false), completeNotes = ref('')
const cancelVisible = ref(false), cancelReason = ref('')
const createVisible = ref(false)
const createForm = reactive({ title: '', member_id: 0, trigger_type: 'gestation_week', due_date: '' })

function statusColor(s: string) { return { pending: 'warning', in_progress: 'primary', completed: 'success', cancelled: 'info' }[s] || '' }
function statusText(s: string) { return { pending: '待处理', in_progress: '进行中', completed: '已完成', cancelled: '已取消' }[s] || s }

async function loadTasks() {
  const params: any = { offset: (page-1)*20, limit: 20 }
  if (filter.status) params.status = filter.status
  if (filter.trigger_type) params.trigger_type = filter.trigger_type
  const { data } = await api.getTasks(params)
  tasks.value = data || []
}
async function loadStats() {
  const { data } = await api.getTaskStats()
  stats.value = data || {}
}
function onPage(p: number) { page = p; loadTasks() }
function assignTask(row: any) { currentTaskId = row.id; assignVisible.value = true }
async function doAssign() {
  await api.assignTask(currentTaskId, stewardId.value)
  ElMessage.success('指派成功'); assignVisible.value = false; loadTasks(); loadStats()
}
function completeTask(row: any) { currentTaskId = row.id; completeNotes.value = ''; completeVisible.value = true }
async function doComplete() {
  await api.completeTask(currentTaskId, completeNotes.value)
  ElMessage.success('已完成'); completeVisible.value = false; loadTasks(); loadStats()
}
function cancelTask(row: any) { currentTaskId = row.id; cancelReason.value = ''; cancelVisible.value = true }
async function doCancel() {
  if (!cancelReason.value) return ElMessage.warning('请输入取消原因')
  await api.cancelTask(currentTaskId, cancelReason.value)
  ElMessage.success('已取消'); cancelVisible.value = false; loadTasks(); loadStats()
}
async function doCreate() {
  await api.createTask({...createForm, due_date: createForm.due_date ? new Date(createForm.due_date).toISOString().slice(0,10) : ''})
  ElMessage.success('创建成功'); createVisible.value = false; loadTasks(); loadStats()
}
onMounted(() => { loadTasks(); loadStats() })
</script>
<style scoped>
.mb20 { margin-bottom: 20px; }
.stat-label { font-size: 13px; color: #999; }
.stat-val { font-size: 28px; font-weight: bold; color: #303133; }
</style>
