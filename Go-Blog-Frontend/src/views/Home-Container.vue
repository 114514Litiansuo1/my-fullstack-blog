<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '../api/request.js'
import CommentSection from '@/views/Comment-Section.vue'

// 存放文章数据的响应式数组
const articles = ref([])
const cursor = ref('0')
const hasMore = ref(true)
const loading = ref(false)
// 获取路由实例，用于页面跳转
const router = useRouter()

// 绑定搜索框的数据
const searchKeyword = ref('')

// 去后端拉取数据的函数
const fetchArticles = async (isLoadMore = false) => {
  if (loading.value || !hasMore.value) return (loading.value = true)
  try {
    // 假设你后端接收游标的参数名为 cursor，每页条数为 limit
    const response = await request.get('/articles', {
      params: {
        cursor: cursor.value,
        limit: 10, // 每次请求 10 条
        keyword: searchKeyword.value,
      },
    })

    // 提取后端返回的三件套
    const newArticles = response.data.data || []
    const nextCursor = response.data.next_cursor
    const moreFlag = response.data.has_more

    if (isLoadMore) {
      // 如果是加载更多，将新数据推入原有数组中 (ES6 扩展运算符)
      articles.value.push(...newArticles)
    } else {
      // 如果是首次加载，直接赋值
      articles.value = newArticles
    }

    // 更新游标和是否有更多数据的状态
    cursor.value = nextCursor
    hasMore.value = moreFlag
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 触发搜索的函数
const handleSearch = () => {
  cursor.value = '0'
  hasMore.value = true
  fetchArticles(false)
}

// 触发加载更多
const loadMore = () => {
  fetchArticles(true)
}

// 点击后台登录按钮触发的跳转
const goToLogin = () => {
  // 假设你的登录页面的路由路径是 /login
  alert("你好")
}

const goToDetail = (id) => {
  // 跳转到详情页，并把文章 ID 传过去
  router.push(`/article/${id}`)
}

// 当组件被挂载到页面上时，自动执行此函数拉取数据
onMounted(() => {
  cursor.value = '0'
  hasMore.value = true
  fetchArticles()
})
</script>

<template>
  <div class="home-container">
    <!-- 顶部导航栏 -->
    <el-menu
      mode="horizontal"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b"
    >
      <el-menu-item index="1" style="font-weight: bold; font-size: 18px"
        >我的个人博客
      </el-menu-item>
      <div style="flex-grow: 1"></div>

      <!-- 搜索框 -->
      <div style="display: flex; align-items: center; height: 60px; margin-right: 20px">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文章标题或摘要..."
          clearable
          @keyup.enter="handleSearch"
          style="width: 250px; margin-top: 10px; margin-bottom: 10px"
          size="default"
        >
          <template #append>
            <el-button icon="Search" @click="handleSearch" />
          </template>
        </el-input>
      </div>

      <!-- 这是一个占位符，把登录按钮挤到最右边 -->
      <el-menu-item index="2" @click="goToLogin">后台登录</el-menu-item>
    </el-menu>

    <!-- 文章列表展示区 -->
    <div class="main-content" style="max-width: 800px; margin: 30px auto; padding: 0 20px">
      <!-- v-for 循环渲染卡片 -->
      <el-card
        v-for="article in articles"
        :key="article.ID"
        style="margin-bottom: 20px; border-radius: 8px"
        shadow="hover"
      >
        <h2 style="margin-top: 0; color: #303133">{{ article.Title }}</h2>
        <!-- 如果没有摘要，就显示默认文字 -->
        <p style="color: #606266; line-height: 1.6">
          {{ article.Summary || '这篇干货文章还没有写摘要...' }}
        </p>

        <div
          style="
            margin-top: 15px;
            display: flex;
            align-items: center;
            justify-content: space-between;
          "
        >
          <div>
            <!-- 如果文章有关联分类，就展示出来 -->
            <el-tag type="info" size="small" v-if="article.Category"
              >{{ article.Category.Name || '没有分类' }}
            </el-tag>
          </div>
          <!-- 纯文本样式的按钮 -->
          <el-button type="primary" text @click="goToDetail(article.ID)">阅读全文 &gt;</el-button>
        </div>
      </el-card>

      <div class="load-more-container">
        <el-button
          v-if="hasMore"
          type="primary"
          plain
          :loading="loading"
          @click="loadMore"
          style="width: 100%"
        >
          {{ loading ? '加载中...' : '加载更多' }}
        </el-button>

        <el-divider v-else-if="articles.length > 0">
          <span style="color: #999">到底啦，没有更多文章了</span>
        </el-divider>
      </div>
      <!-- 数据为空时的提示 -->
      <el-empty v-if="articles.length === 0" description="博主很懒，还没有发布任何文章QAQ" />

      <!--      留言-->
      <CommentSection />
    </div>
  </div>
</template>

<style scoped>
/* scoped 表示这里的样式只在当前文件生效 */
.home-container {
  background-color: #f4f5f7;
  min-height: 100vh;
}
</style>
