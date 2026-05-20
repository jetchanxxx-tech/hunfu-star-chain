<template>
  <el-container class="layout">
    <el-aside width="220px">
      <div class="logo">惠福星链 · 管理后台</div>
      <el-menu :default-active="route.path" router background-color="#1a1f2e" text-color="#aeb5c4" active-text-color="#409EFF">
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据驾驶舱</span>
        </el-menu-item>
        <el-menu-item index="/members">
          <el-icon><User /></el-icon>
          <span>会员管理</span>
        </el-menu-item>
        <el-menu-item index="/packages">
          <el-icon><Goods /></el-icon>
          <span>服务包配置</span>
        </el-menu-item>
        <el-menu-item index="/followup">
          <el-icon><Clock /></el-icon>
          <span>随访规则</span>
        </el-menu-item>
        <el-menu-item index="/tasks">
          <el-icon><List /></el-icon>
          <span>任务管理</span>
        </el-menu-item>
        <el-menu-item index="/timeline-config">
          <el-icon><Timer /></el-icon>
          <span>时间轴配置</span>
        </el-menu-item>
        <el-menu-item index="/verification">
          <el-icon><CircleCheck /></el-icon>
          <span>核销记录</span>
        </el-menu-item>
        <el-menu-item index="/auth-audit">
          <el-icon><Document /></el-icon>
          <span>授权审计</span>
        </el-menu-item>
        <el-menu-item index="/system">
          <el-icon><Setting /></el-icon>
          <span>系统管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <span>{{ $route.meta.title }}</span>
        <el-dropdown @command="onCmd">
          <span class="user-info">管理员 ▼</span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main><router-view /></el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
const route = useRoute()
const router = useRouter()
function onCmd(cmd: string) {
  if (cmd === 'logout') {
    localStorage.removeItem('admin_token')
    router.push('/login')
  }
}
</script>

<style>
.layout { min-height: 100vh; }
.el-aside { background: #1a1f2e; }
.logo { color: #fff; font-size: 15px; font-weight: bold; padding: 20px 16px; text-align: center; border-bottom: 1px solid #2a2f3e; }
.el-header { background: #fff; display: flex; align-items: center; justify-content: space-between;
  border-bottom: 1px solid #eee; padding: 0 20px; font-size: 16px; font-weight: 500; }
.user-info { cursor: pointer; color: #666; font-size: 14px; }
.el-main { background: #f0f2f5; padding: 20px; }
</style>
