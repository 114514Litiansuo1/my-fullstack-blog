<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../../api/request'

const comments = ref([])
const loading = ref(false)

// 获取所有留言
const fetchComments = async () => {
  loading.value = true
  try {
    const res = await request.get('/comments')
    comments.value = res.data.data || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 删除留言
const handleDelete = (id) => {
  ElMessageBox.confirm('确定要永久删除这条评论吗？此操作不可撤销。', '警告', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      try {
        await request.delete(`/comment/${id}`)
        ElMessage.success('评论已成功删除')
        fetchComments() // 刷新列表
      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {})
}

// 封禁 IP 逻辑
const handleBanIP = (ip) => {
  if (!ip) {
    return ElMessage.warning('该留言没有记录 IP 地址，无法封禁')
  }

  ElMessageBox.confirm(
    `确定要将 IP【${ip}】加入黑名单吗？封禁后该 IP 将被彻底禁止留言。`,
    '高危操作提示',
    {
      confirmButtonText: '确定封禁',
      cancelButtonText: '取消',
      type: 'warning',
    },
  )
    .then(async () => {
      try {
        await request.post('/ip/ban', { ip: ip })

        ElMessage.success(`IP: ${ip} 已成功封禁！`)

      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {
    })
}

onMounted(() => {
  fetchComments()
})
</script>

<template>
  <div class="app-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>评论管理</span>
          <el-button icon="Refresh" @click="fetchComments">刷新列表</el-button>
        </div>
      </template>

      <el-table :data="comments" v-loading="loading" border stripe style="width: 100%">
        <el-table-column prop="ID" label="ID" width="180" align="center" />

        <el-table-column prop="nickname" label="昵称" width="120" />

        <el-table-column prop="content" label="留言内容" min-width="250" show-overflow-tooltip />

        <el-table-column prop="ip" label="发布IP" width="140" align="center" />

        <el-table-column label="发布时间" width="180" align="center">
          <template #default="scope">
            {{ new Date(scope.row.CreatedAt).toLocaleString() }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="scope">
            <el-button type="warning" link icon="Lock" @click="handleBanIP(scope.row.ip)">
              封禁IP
            </el-button>
            <el-button type="danger" link icon="Delete" @click="handleDelete(scope.row.ID)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.app-container {
  padding: 20px;
}
</style>
