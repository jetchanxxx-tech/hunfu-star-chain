<template>
  <div>
    <el-card>
      <template #header><span>角色权限管理</span></template>
      <el-table :data="[]" empty-text="暂无用户" style="width:100%">
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="status" label="状态" />
        <el-table-column label="操作" />
      </el-table>
    </el-card>

    <el-card style="margin-top:16px">
      <template #header><span>接口监控</span></template>
      <el-empty description="接口调用日志（P1阶段交付）" />
    </el-card>
  </div>
</template>
