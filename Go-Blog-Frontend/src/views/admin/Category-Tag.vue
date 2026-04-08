<script setup>
import { ref, onMounted } from 'vue'
import request from '../../api/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const activeTab = ref('category')
const categories = ref([])
const tags = ref([])
const newCategoryName = ref('')
const newTagName = ref('')

const fetchData = async () => {
  const [resCat, resTag] = await Promise.all([request.get('/categories'), request.get('/tags')])
  categories.value = resCat.data || []
  tags.value = resTag.data || []
}

const addCategory = async () => {
  if (!newCategoryName.value) return
  await request.post('/category', { name: newCategoryName.value })
  newCategoryName.value = ''
  fetchData()
  ElMessage.success('分类添加成功')
}

const addTag = async () => {
  if (!newTagName.value) return
  await request.post('/tag', { name: newTagName.value })
  newTagName.value = ''
  fetchData()
  ElMessage.success('标签添加成功')
}

const deleteItem = (type, id) => {
  ElMessageBox.confirm('确定删除吗？删除后相关文章的关联将失效。', '警告').then(async () => {
    await request.delete(`/${type}/${id}`)
    fetchData()
    ElMessage.success('删除成功')
  })
}

onMounted(fetchData)
</script>

<template>
  <el-card>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="分类管理" name="category">
        <div style="margin-bottom: 15px">
          <el-input
            v-model="newCategoryName"
            placeholder="新分类名称"
            style="width: 200px; margin-right: 10px"
          />
          <el-button type="primary" @click="addCategory">添加分类</el-button>
        </div>
        <el-table :data="categories" border size="small">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="Name" label="分类名称" />
          <el-table-column label="操作" width="100">
            <template #default="scope">
              <el-button type="danger" link @click="deleteItem('category', scope.row.ID)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="标签管理" name="tag">
        <div style="margin-bottom: 15px">
          <el-input
            v-model="newTagName"
            placeholder="新标签名称"
            style="width: 200px; margin-right: 10px"
          />
          <el-button type="primary" @click="addTag">添加标签</el-button>
        </div>
        <div style="display: flex; flex-wrap: wrap; gap: 10px">
          <el-tag v-for="tag in tags" :key="tag.ID" closable @close="deleteItem('tag', tag.ID)">
            {{ tag.Name }}
          </el-tag>
        </div>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<style scoped></style>
