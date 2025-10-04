import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8081',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

console.log('Axios实例已创建，baseURL:', api.defaults.baseURL)

// 请求拦截器 - 添加token
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理token过期等情况
api.interceptors.response.use(
  response => {
    console.log('API请求成功:', response.config.url, response.data)
    return response.data
  },
  error => {
    if (error.response) {
      // 处理HTTP错误
      console.error('API请求错误:', error.config.url, error.response.status, error.response.data)
      switch (error.response.status) {
        case 401:
          // token过期或无效，清除localStorage中的token并跳转到登录页面
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          // 这里可以添加跳转到登录页面的逻辑
          break
        default:
          // 其他错误情况
          console.error('API请求错误:', error.response.data)
      }
    } else if (error.request) {
      // 请求发出但没有收到响应
      console.error('API请求无响应:', error.config.url)
      console.error('网络错误详情:', error.message)
    } else {
      // 设置请求时发生错误
      console.error('API请求配置错误:', error.message)
    }
    return Promise.reject(error)
  }
)

// 认证相关API方法
export const authService = {
  // 用户注册
  register: async (userData) => {
    try {
      const response = await api.post('/auth/register', userData)
      return response
    } catch (error) {
      console.error('注册失败:', error)
      throw error
    }
  },

  // 用户登录
  login: async (credentials) => {
    try {
      const response = await api.post('/auth/login', credentials)
      // 保存token和用户信息到localStorage
      if (response.token && response.user) {
        localStorage.setItem('token', response.token)
        localStorage.setItem('userInfo', JSON.stringify(response.user))
      }
      return response
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  },

  // 用户登出
  logout: () => {
    // 清除localStorage中的token和用户信息
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  },

  // 获取当前登录用户信息
  getCurrentUser: () => {
    const userInfo = localStorage.getItem('userInfo')
    return userInfo ? JSON.parse(userInfo) : null
  },

  // 检查用户是否已登录
  isAuthenticated: () => {
    const token = localStorage.getItem('token')
    return !!token
  }
}

export default api