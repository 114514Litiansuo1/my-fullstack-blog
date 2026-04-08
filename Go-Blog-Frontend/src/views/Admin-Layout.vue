<script setup>
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ChatLineRound, Document, Folder, SwitchButton } from '@element-plus/icons-vue'

const router = useRouter()

// 退出登录逻辑
const handleLogout = () => {
  localStorage.removeItem('token')
  ElMessage.success('已安全退出')
  router.push('/login')
}
</script>

<template>
  <el-container class="admin-container">
    <!-- 左侧侧边栏 -->
    <el-aside width="200px" class="sidebar">
      <div class="logo">博客后台</div>
      <el-menu
        active-text-color="#409EFF"
        background-color="#304156"
        text-color="#bfcbd9"
        :default-active="$route.path"
        router
        class="el-menu-vertical"
      >
        <el-menu-item index="/admin/articles">
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/categories">
          <el-icon><Folder /></el-icon>
          <span>分类管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/comments">
          <el-icon><ChatLineRound /></el-icon>
          <span>评论管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧内容区 -->
    <el-container>
      <!-- 顶部 Header -->
      <el-header class="header">
        <div class="header-left">
          <span>欢迎回来，超级管理员</span>
        </div>
        <div class="header-right">
          <el-button type="danger" size="small" plain @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>&nbsp; 退出登录</el-button
          >
        </div>
      </el-header>

      <!-- 子路由画框 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.admin-container {
  height: 100vh;
}
.sidebar {
  background-color: #304156;
  color: white;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #1f2d3d;
}
.el-menu-vertical {
  border-right: none;
}
.header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}
.main-content {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
