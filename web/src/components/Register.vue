<template>
  <div class="auth-container">
    <div class="auth-form-wrapper">
      <h2>用户注册</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            type="text"
            v-model="formData.username"
            placeholder="请输入用户名（3-50个字符）"
            required
            minlength="3"
            maxlength="50"
            autofocus
          />
        </div>
        
        <div class="form-group">
          <label for="email">邮箱</label>
          <input
            id="email"
            type="email"
            v-model="formData.email"
            placeholder="请输入邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            type="password"
            v-model="formData.password"
            placeholder="请输入密码（至少6个字符）"
            required
            minlength="6"
          />
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            id="confirmPassword"
            type="password"
            v-model="formData.confirmPassword"
            placeholder="请再次输入密码"
            required
            minlength="6"
          />
        </div>
        
        <div class="form-actions">
          <button type="submit" :disabled="isLoading" class="btn btn-primary">
            {{ isLoading ? '注册中...' : '注册' }}
          </button>
          <p class="login-link">
            已有账号？<button type="button" @click="$emit('switch-to-login')">立即登录</button>
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
const emit = defineEmits(['switch-to-login', 'register-success'])

// 表单数据
const formData = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

// 状态管理
const isLoading = ref(false)
const error = ref('')
const successMessage = ref('')

// 处理注册
const handleRegister = async () => {
  // 重置错误和成功消息
  error.value = ''
  successMessage.value = ''
  
  // 表单验证
  if (!formData.value.username || !formData.value.email || !formData.value.password || !formData.value.confirmPassword) {
    error.value = '请填写所有必填字段'
    return
  }
  
  if (formData.value.username.length < 3 || formData.value.username.length > 50) {
    error.value = '用户名长度必须在3-50个字符之间'
    return
  }
  
  if (formData.value.password.length < 6) {
    error.value = '密码长度至少为6个字符'
    return
  }
  
  if (formData.value.password !== formData.value.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }
  
  // 简单的邮箱格式验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(formData.value.email)) {
    error.value = '请输入有效的邮箱地址'
    return
  }
  
  // 设置加载状态
  isLoading.value = true
  
  try {
    // 调用注册API
    const response = await authService.register({
      username: formData.value.username,
      email: formData.value.email,
      password: formData.value.password,
      confirm_password: formData.value.confirmPassword
    })
    
    // 显示成功消息
    successMessage.value = response.message || '注册成功'
    
    // 通知父组件注册成功
    emit('register-success')
    
    // 重置表单
    formData.value = {
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    }
    
  } catch (err) {
    // 处理注册错误
    if (err.response && err.response.data) {
      error.value = err.response.data.error || '注册失败，请稍后重试'
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

.login-link {
  text-align: center;
  margin-top: 1rem;
  color: #666;
}

.login-link button {
  background: none;
  border: none;
  color: #667eea;
  cursor: pointer;
  font-size: 1rem;
  text-decoration: underline;
}

.login-link button:hover {
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