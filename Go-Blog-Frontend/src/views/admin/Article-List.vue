<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../../api/request' // 定制axios

// Markdown
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { ArrowRight } from '@element-plus/icons-vue'

// photo
const onUploadImg = async (files, callback) => {
  try {
    const resUrls = await Promise.all(
      files.map(async (file) => {
        const formData = new FormData()
        formData.append('file', file)

        const res = await request.post('/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })

        return res.data.url
      }),
    )

    callback(resUrls)
  } catch (error) {
    ElMessage.error('图片上传失败')
    console.error(error)
  }
}

// --- 数据列表相关状态 ---
const articles = ref([])
const loading = ref(false)
const categories = ref([])
const tags = ref([])

// --- 弹窗与表单相关状态 ---
const dialogVisible = ref(false)
const dialogTitle = ref('发布文章')
const submitLoading = ref(false)
const formRef = ref(null)

// --- 游标分页核心逻辑 ---
const currentCursor = ref('0')
const nextCursor = ref('')
const hasMore = ref(false)
const cursorStack = ref([])
const pagNum = ref(parseInt(currentCursor.value + 1, 10))

const form = reactive({
  ID: null,
  title: '',
  summary: '',
  content: '',
  category_id: 0,
  tag_ids: [],
})

// 表单验证规则
const rules = {
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  content: [{ required: true, message: '正文不能为空', trigger: 'blur' }],
}

// ================== 方法定义 ==================

// 获取文章列表 (Read)
const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await request.get('/articles', {
      params: {
        cursor: currentCursor.value,
        limit: 10,
      },
    })
    articles.value = res.data.data || []
    nextCursor.value = res.data.next_cursor
    hasMore.value = res.data.has_more
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 点击下一页
const goNext = () => {
  cursorStack.value.push(currentCursor.value)
  currentCursor.value = nextCursor.value
  fetchArticles()
}

// 点击上一页
const goPrev = () => {
  const prev = cursorStack.value.pop()
  currentCursor.value = prev
  fetchArticles()
}

// 刷新（重置到第一页）
const refreshList = () => {
  currentCursor.value = '0'
  cursorStack.value = []
  fetchArticles()
}

// 获取分类和标签列表
const fetchOptions = async () => {
  try {
    const [resCat, resTag] = await Promise.all([request.get('/categories'), request.get('/tags')])
    categories.value = resCat.data || []
    tags.value = resTag.data || []
  } catch (error) {
    console.error('加载选项失败', error)
  }
}

// 点击“发布新文章”按钮
const handleAdd = () => {
  dialogTitle.value = '发布新文章'
  form.ID = null
  dialogVisible.value = true
}

// 点击“编辑”按钮
const handleEdit = (row) => {
  dialogTitle.value = '编辑文章'
  form.ID = row.ID
  form.title = row.Title
  form.summary = row.Summary
  form.content = row.Content
  form.category_id = row.CategoryId || 0
  form.tag_ids = row.Tags ? row.Tags.map((t) => t.ID) : []
  dialogVisible.value = true
}

// 提交表单 (Create / Update)
const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (form.ID) {
          // 有 ID，说明是修改操作
          await request.put(`/article/${form.ID}`, form)
          ElMessage.success('修改成功')
        } else {
          // 没有 ID，说明是新增操作
          await request.post('/article', form)
          ElMessage.success('发布成功')
        }
        dialogVisible.value = false
        fetchArticles()
      } catch (error) {
      } finally {
        submitLoading.value = false
      }
    }
  })
}

// 点击“删除”按钮 (Delete)
const handleDelete = (ID) => {
  ElMessageBox.confirm('确认将这篇文章移入回收站吗？（软删除）', '危险操作提示', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      try {
        await request.delete(`/article/${ID}`)
        ElMessage.success('删除成功')
        fetchArticles()
      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {
    })
}

// 重置表单 (每次关闭弹窗时触发)
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  form.ID = null
}

// --- 初始化 ---
onMounted(() => {
  fetchArticles()
  fetchOptions()
})
</script>

<template>
  <div class="app-container">
    <!-- 1. 顶部操作区 -->
    <div style="margin-bottom: 20px; display: flex; justify-content: space-between">
      <el-button type="primary" icon="Plus" @click="handleAdd">发布新文章</el-button>
      <el-button icon="Refresh" @click="fetchArticles">刷新</el-button>
    </div>

    <!-- 2. 数据表格 -->
    <el-table :data="articles" v-loading="loading" border stripe style="width: 100%">
      <el-table-column prop="ID" label="ID" width="200" align="center" />
      <el-table-column prop="Title" label="文章标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="Summary" label="摘要" min-width="250" show-overflow-tooltip />

      <el-table-column label="发布时间" width="180">
        <template #default="scope">
          <!-- 格式化 GORM 自带的 CreatedAt 时间 -->
          {{ new Date(scope.row.CreatedAt).toLocaleString() }}
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" link @click="handleEdit(scope.row)"
            >编辑</el-button
          >
          <el-button size="small" type="danger" link @click="handleDelete(scope.row.ID)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <!--    分页-->
    <div
      style="
        margin-top: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 20px;
      "
    >
      <el-button :disabled="cursorStack.length === 0" icon="ArrowLeft" @click="goPrev">
        上一页
      </el-button>

      <span style="color: #606266; font-size: 14px">当前位置：{{ pagNum }}</span>
      <el-button @click="refreshList" icon="Refresh">刷新</el-button>

      <el-button :disabled="!hasMore" @click="goNext">
        下一页 <el-icon class="el-icon--right"><ArrowRight /></el-icon>
      </el-button>
    </div>

    <!-- 3. 弹窗表单 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="80%"
      top="5vh"
      @close="resetForm"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="摘要" prop="summary">
          <el-input v-model="form.summary" placeholder="请输入文章简短摘要" />
        </el-form-item>

        <!-- 分类  标签   -->

        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择文章分类" style="width: 100%">
            <el-option
              v-for="item in categories"
              :key="item.ID"
              :label="item.Name"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="标签">
          <el-select
            v-model="form.tag_ids"
            multiple
            placeholder="请选择文章标签"
            style="width: 100%"
          >
            <el-option v-for="item in tags" :key="item.ID" :label="item.Name" :value="item.ID" />
          </el-select>
        </el-form-item>

        <!--  Markdown 编辑器 -->
        <el-form-item label="正文" prop="content">
          <MdEditor v-model="form.content" style="height: 500px" @onUploadImg="onUploadImg" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitLoading">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.app-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>
