<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../api/request'
import { User, Warning } from '@element-plus/icons-vue' // 定制 Axios

// === 状态定义 ===
const comments = ref([])
const isLoadingList = ref(false)
const isSubmitting = ref(false)
const isBanned = ref(false) // 是否被封禁 (403状态)

const form = reactive({
  nickname: '',
  content: '',
})

// === 逻辑方法 ===

// 1. 获取留言列表
const fetchComments = async () => {
  isLoadingList.value = true
  try {
    const res = await request.get('/comments')
    // 根据后端实际返回的 JSON 结构调整
    comments.value = res.data.data || []
  } catch (error) {
    console.error('加载留言失败', error)
  } finally {
    isLoadingList.value = false
  }
}

// 2. 提交留言
const submitComment = async () => {
  // 简单的前端校验
  if (!form.nickname.trim()) {
    return ElMessage.warning('请输入昵称哦')
  }
  if (!form.content.trim()) {
    return ElMessage.warning('留言内容不能为空')
  }

  isSubmitting.value = true
  try {
    await request.post('/comment', {
      nickname: form.nickname,
      content: form.content,
    })

    ElMessage.success('留言发布成功！')
    form.content = ''
    fetchComments()
  } catch (error) {
    if (error.response) {
      const status = error.response.status
      if (status === 429) {
        ElMessage.warning('您留言太频繁啦，请稍作休息再试！')
      } else if (status === 403) {
        ElMessage.error('您的 IP 已被系统封禁，无法继续留言。')
        isBanned.value = true
      }
    }
  } finally {
    isSubmitting.value = false
  }
}

// 3. 日期格式化工具
const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  fetchComments()
})
</script>

<!-- src/components/CommentSection.vue -->
<template>
  <div class="comment-section">
    <el-divider><h3>访客留言区</h3></el-divider>

    <!-- 1. 留言发布表单 -->
    <el-card shadow="never" class="comment-form-card">
      <el-form :model="form" ref="formRef" label-position="top">
        <!-- 昵称输入限制长度 -->
        <el-form-item label="您的昵称" required>
          <el-input
            v-model="form.nickname"
            placeholder="请输入您的昵称（最多 15 个字符）"
            maxlength="15"
            show-word-limit
            style="max-width: 300px"
            :disabled="isBanned"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <!-- 留言内容：文本域 -->
        <el-form-item label="留言内容" required>
          <el-input
            type="textarea"
            v-model="form.content"
            placeholder="说点什么吧..."
            :rows="4"
            maxlength="500"
            show-word-limit
            :disabled="isBanned"
          />
        </el-form-item>

        <!-- 提交按钮与状态提示 -->
        <el-form-item>
          <el-button
            type="primary"
            @click="submitComment"
            :loading="isSubmitting"
            :disabled="isBanned"
          >
            发布留言
          </el-button>

          <!-- 如果被封禁，在按钮旁边显示醒目提示 -->
          <span v-if="isBanned" class="banned-tips">
            <el-icon><Warning /></el-icon> 您的 IP 已被限制留言
          </span>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 2. 留言展示列表 -->
    <div class="comment-list" v-loading="isLoadingList">
      <h4 style="margin-bottom: 20px">最新留言 ({{ comments.length }})</h4>

      <el-empty v-if="comments.length === 0" description="暂无留言，来做第一个发言的人吧！" />

      <div class="comment-item" v-for="item in comments" :key="item.id">
        <!-- 默认头像 -->
        <el-avatar
          :size="45"
          src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png"
          class="avatar"
        />

        <div class="comment-main">
          <div class="comment-header">
            <span class="nickname">{{ item.nickname }}</span>
            <!-- 格式化时间 -->
            <span class="time">{{ formatDate(item.created_at) }}</span>
          </div>
          <!-- 留言内容：使用纯文本展示，防范 XSS -->
          <div class="content">{{ item.content }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-section {
  margin-top: 50px;
}

.comment-form-card {
  margin-bottom: 30px;
  background-color: #fafafa;
  border-radius: 8px;
}

.banned-tips {
  margin-left: 15px;
  color: #f56c6c;
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.comment-item {
  display: flex;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px dashed #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.avatar {
  flex-shrink: 0;
  margin-right: 15px;
  border: 1px solid #ebeef5;
}

.comment-main {
  flex-grow: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.nickname {
  font-weight: bold;
  color: #303133;
  font-size: 15px;
}

.time {
  font-size: 12px;
  color: #909399;
}

.content {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  word-break: break-all;
}
</style>
