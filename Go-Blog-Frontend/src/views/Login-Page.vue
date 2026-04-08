<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import request from '../api/request.js'
import { ElMessage } from 'element-plus'

const router = useRouter()

// 定义表单绑定的数据
const loginForm = reactive({
  username: '',
  password: '',
})

// 定义表单验证规则
const rules = {
  username: [{ required: true, message: '账号不能为空', trigger: 'blur' }],
  password: [{ required: true, message: '密码不能为空', trigger: 'blur' }],
}

const formRef = ref(null)
const loading = ref(false)

// 处理登录逻辑
const handleLogin = async () => {
  // 先验证表单填没填完
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 向后端发送 POST 请求
        const response = await request.post('http://localhost:8080/api/v1/login', {
          username: loginForm.username,
          password: loginForm.password,
        })

        // 登录成功！
        ElMessage.success('登录成功，欢迎回来！')

        // 把后端发来的 JWT Token 存到浏览器的本地存储中
        localStorage.setItem('token', response.data.token)

        // 跳转到后台管理页面
        await router.push('/admin')
      } catch (error) {
        // 处理密码错误或账号不存在的情况
        if (error.response && error.response.status === 401) {
          ElMessage.error('账号或密码错误！')
        } else {
          ElMessage.error('服务器异常，请检查后端是否启动')
        }
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <h2 class="login-title">博客后台登录</h2>

      <!-- Element Plus 的表单组件 -->
      <el-form :model="loginForm" :rules="rules" ref="formRef" label-width="0">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入管理员账号"
            prefix-icon="user"
            size="large"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="lock"
            size="large"
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%"
            size="large"
            :loading="loading"
            @click="handleLogin"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #2b3137;
}

.login-card {
  width: 400px;
  padding: 20px;
  border-radius: 10px;
}

.login-title {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}
</style>
