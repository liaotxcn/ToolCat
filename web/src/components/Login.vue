<template>
  <div class="auth-container">
    <div class="auth-form-wrapper">
      <h2>用户登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            type="text"
            v-model="formData.username"
            placeholder="请输入用户名"
            required
            autofocus
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            type="password"
            v-model="formData.password"
            placeholder="请输入密码"
            required
          />
        </div>
        
        <div class="form-actions">
          <button type="submit" :disabled="isLoading" class="btn btn-primary">
            {{ isLoading ? '登录中...' : '登录' }}
          </button>
          <p class="register-link">
            还没有账号？<button type="button" @click="$emit('switch-to-register')">立即注册</button>
          </p>
        </div>
      </form>
      
      <!-- 错误提示 -->
      <div v-if="error" class="error-message">
        {{ error }}
      </div>
      
      <!-- 成功提示 -->
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { authService } from '../services/auth.js'

// 定义props和emits
const emit = defineEmits(['switch-to-register', 'login-success'])

// 表单数据
const formData = ref({
  username: '',
  password: ''
})

// 状态管理
const isLoading = ref(false)
const error = ref('')
const successMessage = ref('')

// 处理登录
const handleLogin = async () => {
  // 重置错误和成功消息
  error.value = ''
  successMessage.value = ''
  
  // 表单验证
  if (!formData.value.username || !formData.value.password) {
    error.value = '请填写所有必填字段'
    return
  }
  
  // 设置加载状态
  isLoading.value = true
  
  try {
    // 调用登录API
    const response = await authService.login({
      username: formData.value.username,
      password: formData.value.password
    })
    
    // 显示成功消息
    successMessage.value = response.message || '登录成功'
    
    // 通知父组件登录成功
    emit('login-success')
    
  } catch (err) {
    // 处理登录错误
    if (err.response && err.response.data) {
      error.value = err.response.data.error || '登录失败，请检查用户名和密码'
    } else {
      error.value = '网络错误，请稍后重试'
    }
  } finally {
    // 重置加载状态
    isLoading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 1rem;
}

.auth-form-wrapper {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.auth-form-wrapper h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-actions {
  margin-top: 1.5rem;
}

.btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 100%;
}

.btn-primary {
  background-color: #667eea;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #5a67d8;
}

.btn-primary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 1rem;
  color: #666;
}

.register-link button {
  background: none;
  border: none;
  color: #667eea;
  cursor: pointer;
  font-size: 1rem;
  text-decoration: underline;
}

.register-link button:hover {
  color: #5a67d8;
}

.error-message {
  margin-top: 1rem;
  padding: 0.8rem;
  background-color: #fed7d7;
  color: #c53030;
  border-radius: 4px;
  text-align: center;
}

.success-message {
  margin-top: 1rem;
  padding: 0.8rem;
  background-color: #c6f6d5;
  color: #276749;
  border-radius: 4px;
  text-align: center;
}
</style>