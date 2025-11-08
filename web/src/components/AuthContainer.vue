<template>
  <div class="auth-container">
    <!-- 背景装饰元素 -->
    <div class="auth-bg-pattern"></div>
    
    <div class="auth-card">
      <!-- 品牌展示区域 -->
      <div class="brand">
        <div class="logo-container">
          <div class="logo">
            <svg viewBox="0 0 40 40" width="40" height="40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="5" y="5" width="30" height="30" rx="6" stroke="currentColor" stroke-width="2"/>
              <path d="M12 20H16L18 16L20 20H24L28 12V28H12V20Z" fill="currentColor" opacity="0.8"/>
            </svg>
          </div>
          <h1>Weave</h1>
        </div>
        <p class="brand-subtitle">欢迎使用，请先登录或注册</p>
      </div>
      
      <!-- 标签切换区域 -->
      <div class="tabs">
        <button :class="['tab', showLogin ? 'active' : '']" @click="switchToLogin">
          <span class="tab-icon">
            <svg viewBox="0 0 20 20" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 12a6 6 0 11-12 0 6 6 0 0112 0z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 10h2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 10h2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
          登录
        </button>
        <button :class="['tab', !showLogin ? 'active' : '']" @click="switchToRegister">
          <span class="tab-icon">
            <svg viewBox="0 0 20 20" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M10 5v10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M5 10h10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
          注册
        </button>
        <div class="tab-indicator" :style="{ transform: showLogin ? 'translateX(0%)' : 'translateX(100%)' }"></div>
      </div>
      
      <!-- 表单区域 -->
      <div class="form-area">
        <transition name="form-switch" mode="out-in" appear>
          <Login
            v-if="showLogin"
            @switch-to-register="switchToRegister"
            @login-success="handleLoginSuccess"
          />
          <Register
            v-else
            @switch-to-login="switchToLogin"
            @register-success="handleRegisterSuccess"
          />
        </transition>
      </div>
      
      <!-- 页脚信息 -->
      <div class="auth-footer">
        <p>© {{ new Date().getFullYear() }} Weave - 插件开发/服务聚合平台</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Login from './Login.vue'
import Register from './Register.vue'
import { authService } from '../services/auth.js'
import api from '../services/auth.js'

// 定义props和emits
const emit = defineEmits(['auth-success'])

// 状态管理
const showLogin = ref(true) // 默认显示登录表单

// 初始化时检查用户是否已登录，并预热CSRF Cookie
onMounted(() => {
  // 预热CSRF：访问健康检查接口以便后端设置XSRF-TOKEN Cookie
  api.get('/health', { withCredentials: true }).catch(() => {})

  if (authService.isAuthenticated()) {
    // 用户已登录，通知父组件
    emit('auth-success')
  }
})

// 切换到注册表单
const switchToRegister = () => {
  showLogin.value = false
}

// 切换到登录表单
const switchToLogin = () => {
  showLogin.value = true
}

// 处理登录成功
const handleLoginSuccess = () => {
  // 登录成功后，通知父组件
  emit('auth-success')
}

// 处理注册成功
const handleRegisterSuccess = () => {
  // 注册成功后，自动切换到登录表单
  showLogin.value = true
}
</script>

<style scoped>
/* 背景容器 */
.auth-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-700) 100%);
  padding: 24px;
  position: relative;
  overflow: hidden;
}

/* 背景装饰图案 */
.auth-bg-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 20% 30%, rgba(255, 255, 255, 0.1) 0%, transparent 25%),
    radial-gradient(circle at 80% 70%, rgba(255, 255, 255, 0.1) 0%, transparent 30%);
  z-index: 1;
}

/* 主卡片 */
.auth-card {
  width: 100%;
  max-width: 480px;
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  overflow: hidden;
  position: relative;
  z-index: 2;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.auth-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

/* 品牌区域 */
.brand {
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-700) 100%);
  color: var(--bg-primary);
  padding: 32px 32px 24px;
  text-align: center;
}

.logo-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 12px;
}

.logo {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: var(--radius-lg);
  color: var(--bg-primary);
}

.brand h1 {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  margin: 0;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--primary-100) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-subtitle {
  font-size: var(--font-size-sm);
  opacity: 0.95;
  margin: 0;
  font-weight: var(--font-weight-medium);
}

/* 标签切换区域 */
.tabs {
  position: relative;
  display: flex;
  gap: 8px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-light);
  background: var(--bg-secondary);
}

.tab {
  flex: 1;
  padding: 12px 16px;
  background: var(--bg-primary);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-lg);
  cursor: pointer;
  color: var(--text-secondary);
  font-weight: var(--font-weight-medium);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all var(--transition-normal);
}

.tab:hover:not(.active) {
  background: var(--bg-secondary);
  border-color: var(--primary-300);
  color: var(--primary-600);
}

.tab.active {
  background: var(--primary-600);
  border-color: var(--primary-600);
  color: var(--bg-primary);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
  transform: translateY(-1px);
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 滑动指示条 */
.tab-indicator {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 50%;
  height: 3px;
  background: var(--primary-600);
  border-radius: var(--radius-full);
  transition: transform var(--transition-slow) cubic-bezier(0.22, 0.61, 0.36, 1);
}

/* 表单区域 */
.form-area {
  padding: 24px 24px;
  min-height: 360px;
}

/* 表单切换过渡动画 */
.form-switch-enter-active,
.form-switch-leave-active {
  transition: opacity var(--transition-normal), transform var(--transition-normal);
}

.form-switch-enter-from {
  opacity: 0;
  transform: translateX(10px) scale(0.98);
}

.form-switch-enter-to {
  opacity: 1;
  transform: translateX(0) scale(1);
}

.form-switch-leave-from {
  opacity: 1;
  transform: translateX(0) scale(1);
}

.form-switch-leave-to {
  opacity: 0;
  transform: translateX(-10px) scale(0.98);
}

/* 表单统一风格 */
:deep(form) {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

:deep(.form-group) {
  display: flex;
  flex-direction: column;
  gap: 6px;
  position: relative;
}

:deep(.form-group label) {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  letter-spacing: 0.02em;
}

:deep(.form-group input) {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--border-light);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-base);
  transition: all var(--transition-normal);
  background: var(--bg-primary);
  color: var(--text-primary);
}

:deep(.form-group input:focus) {
  outline: none;
  border-color: var(--primary-500);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
  background: var(--bg-primary);
}

:deep(.form-group input:placeholder-shown) {
  color: var(--text-muted);
}

:deep(.form-group input::placeholder) {
  color: var(--text-muted);
  opacity: 1;
}

/* 提交按钮 */
:deep(button[type="submit"]) {
  padding: 12px 16px;
  border: none;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-700) 100%);
  color: var(--bg-primary);
  cursor: pointer;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-normal);
  position: relative;
  overflow: hidden;
}

:deep(button[type="submit"])::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  transform: translate(-50%, -50%);
  transition: width 0.6s ease, height 0.6s ease;
}

:deep(button[type="submit"]:hover::before) {
  width: 300px;
  height: 300px;
}

:deep(button[type="submit"]:hover) {
  transform: translateY(-1px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

:deep(button[type="submit"]:active) {
  transform: translateY(0);
}

:deep(button[type="submit"]:disabled) {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 错误消息 */
:deep(.error-message) {
  background: var(--error-100);
  color: var(--error-700);
  border: 1px solid var(--error);
  padding: 10px 12px;
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

:deep(.error-message)::before {
  content: '';
  width: 16px;
  height: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='%23dc2626'%3E%3Cpath fill-rule='evenodd' d='M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z' clip-rule='evenodd'/%3E%3C/svg%3E") no-repeat center center;
  background-size: contain;
}

/* 切换提示 */
:deep(.switch-tip) {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-top: 16px;
}

:deep(.link-btn) {
  font-size: var(--font-size-sm);
  background: none;
  border: none;
  color: var(--primary-600);
  cursor: pointer;
  font-weight: var(--font-weight-medium);
  padding: 2px 6px;
  border-radius: var(--radius);
  transition: all var(--transition-fast);
  text-decoration: none;
}

:deep(.link-btn:hover) {
  background: var(--primary-50);
  color: var(--primary-700);
}

/* 页脚 */
.auth-footer {
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-700) 100%);
  border-top: 3px solid var(--primary-400);
  text-align: center;
  position: relative;
  overflow: hidden;
}

/* 页脚装饰元素 */
.auth-footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 20% 80%, rgba(255, 255, 255, 0.15) 0%, transparent 20%),
              radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 15%);
  z-index: 1;
}

.auth-footer p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--bg-primary);
  font-weight: var(--font-weight-medium);
  letter-spacing: 0.02em;
  position: relative;
  z-index: 2;
  display: inline-block;
  padding: 4px 8px;
  border-radius: var(--radius);
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(4px);
  transition: all var(--transition-normal);
}

.auth-footer p:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

/* 响应式设计 */
@media (max-width: 480px) {
  .auth-container {
    padding: 16px;
  }
  
  .brand {
    padding: 24px 20px 20px;
  }
  
  .form-area {
    padding: 20px 20px;
  }
  
  .auth-card {
    max-width: 100%;
    border-radius: var(--radius-lg);
  }
}
</style>