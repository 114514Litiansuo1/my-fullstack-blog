import axios from 'axios'
import router from '../router'
import { ElMessage } from 'element-plus'

const service = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 5000,
})

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

service.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('登录状态已失效，请重新登录')
          localStorage.removeItem('token')
          router.push('/login')
          break
        case 403:
          ElMessage.error('您没有权限执行此操作')
          break
        case 404:
          ElMessage.error('请求的接口不存在')
          break
        case 500:
          ElMessage.error('后端服务器异常')
          break
      }
    } else {
      ElMessage.error('网络连接失败，请检查后端服务是否启动')
    }
    return Promise.reject(error)
  },
)

export default service
