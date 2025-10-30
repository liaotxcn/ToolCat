<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="brand">
        <h1>ToolCat</h1>
        <p>欢迎使用，请先登录或注册</p>
      </div>
      <div class="tabs">
        <button :class="['tab', showLogin ? 'active' : '']" @click="switchToLogin">登录</button>
        <button :class="['tab', !showLogin ? 'active' : '']" @click="switchToRegister">注册</button>
        <div class="tab-indicator" :style="{ transform: showLogin ? 'translateX(0%)' : 'translateX(100%)' }"></div>
      </div>
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
.auth-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

.auth-card {
  width: 100%;
  max-width: 460px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.15);
  overflow: hidden;
}

.brand {
  background: linear-gradient(135deg, #6b8cff 0%, #7a57ff 100%);
  color: #fff;
  padding: 20px 24px;
}
.brand h1 { font-size: 24px; margin-bottom: 4px; }
.brand p { opacity: 0.9; font-size: 14px; }

.tabs {
  position: relative;
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
}
.tab {
  flex: 1;
  padding: 10px 0;
  background: #f7f7fb;
  border: 1px solid #e6e6f2;
  border-radius: 8px;
  cursor: pointer;
  color: #555;
  transition: background-color .2s ease, color .2s ease, box-shadow .2s ease;
}
.tab.active {
  background: #667eea;
  border-color: #667eea;
  color: #fff;
  box-shadow: 0 4px 12px rgba(102,126,234,.35);
}
.tab:hover { filter: brightness(0.98); }

/* 滑动指示条（在两个tab之间平滑移动） */
.tab-indicator {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 50%;
  height: 3px;
  background: #667eea;
  border-radius: 3px;
  transition: transform .25s cubic-bezier(.22,.61,.36,1);
}

.form-area { padding: 16px 20px 20px; min-height: 340px; }

/* 表单切换过渡动画：淡入淡出 + 轻微位移 */
.form-switch-enter-active, .form-switch-leave-active {
  transition: opacity .18s ease-out, transform .18s ease-out;
}
.form-switch-enter-from { opacity: 0; transform: translateX(6px); }
.form-switch-enter-to   { opacity: 1; transform: translateX(0); }
.form-switch-leave-from { opacity: 1; transform: translateX(0); }
.form-switch-leave-to   { opacity: 0; transform: translateX(-6px); }

/* 表单统一风格（影响子组件） */
:deep(form) { display: flex; flex-direction: column; gap: 12px; }
:deep(.form-group label) { font-weight: 500; color: #333; margin-bottom: 6px; }
:deep(.form-group input) {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d9d9e3;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color .2s ease, box-shadow .2s ease;
}
:deep(.form-group input:focus) {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102,126,234,.2);
}
:deep(button[type="submit"]) {
  padding: 10px 14px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #5a67d8 100%);
  color: #fff;
  cursor: pointer;
  font-size: 15px;
  transition: transform .05s ease, box-shadow .2s ease;
}
:deep(button[type="submit"]:hover) { box-shadow: 0 8px 20px rgba(102,126,234,.35); }
:deep(button[type="submit"]:active) { transform: translateY(1px); }

:deep(.error-message) {
  background: #fdecea;
  color: #d93025;
  border: 1px solid #f2a19a;
  padding: 8px 10px;
  border-radius: 8px;
}

:deep(.switch-tip) { text-align: center; }
:deep(.link-btn) { font-size: 14px; }
</style>