<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import request from '../api/request'

// 引入 md-editor-v3 的预览组件和预览专属样式
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { Calendar, Folder } from '@element-plus/icons-vue'
import CommentSection from '@/views/Comment-Section.vue'

const route = useRoute()
const article = ref({})
const loading = ref(false)

// 获取单个文章的数据
const fetchArticleDetail = async () => {
  loading.value = true
  try {
    const res = await request.get(`/article/${route.params.id}`)
    article.value = res.data.data || res.data
  } catch (error) {
    console.error('获取文章详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 简单的日期格式化工具
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

onMounted(() => {
  fetchArticleDetail()
})
</script>

<template>
  <div class="article-detail-container">
    <!-- 顶部导航栏  -->
    <el-menu
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
      router
    >
      <el-menu-item index="/">返回首页</el-menu-item>
      <!--      <el-button style="flex-grow: 1">我是按钮</el-button>-->
    </el-menu>

    <div class="main-content" v-loading="loading">
      <!-- 如果文章还在加载或者不存在，显示骨架屏或空状态 -->
      <div v-if="article.Title">
        <!-- 1. 文章头部信息 -->
        <div class="article-header">
          <h1 class="title">{{ article.Title }}</h1>
          <div class="meta-info">
            <span class="time"
              ><el-icon><Calendar /></el-icon> {{ formatDate(article.CreatedAt) }}</span
            >
            <span class="category" v-if="article.Category">
              <el-icon><Folder /></el-icon> {{ article.Category.Name }}
            </span>
          </div>

          <div class="tags" v-if="article.Tags && article.Tags.length > 0">
            <el-tag
              v-for="tag in article.Tags"
              :key="tag.ID"
              size="small"
              type="info"
              class="tag-item"
            >
              {{ tag.Name }}
            </el-tag>
          </div>
        </div>

        <el-divider />

        <!-- 2. 正文渲染区 -->
        <!-- MdPreview 会自动将 article.Content 里的 Markdown 渲染为带排版和代码高亮的网页 -->
        <div class="article-body">
          <MdPreview :modelValue="article.Content" />
        </div>
      </div>

      <el-empty v-else-if="!loading" description="文章找不到了..." />

<!--      留言-->
      <CommentSection />

    </div>
  </div>
</template>

<style scoped>
.article-detail-container {
  background-color: #f4f5f7;
  min-height: 100vh;
}

.main-content {
  max-width: 800px;
  margin: 30px auto;
  padding: 40px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.article-header {
  text-align: center;
  margin-bottom: 30px;
}

.title {
  font-size: 28px;
  color: #303133;
  margin-bottom: 15px;
}

.meta-info {
  color: #909399;
  font-size: 14px;
  margin-bottom: 15px;
  display: flex;
  justify-content: center;
  gap: 20px;
}

.tag-item {
  margin: 0 5px;
}

:deep(.md-editor-preview-wrapper) {
  padding: 0;
}
</style>
