<template>
  <div class="auth-container-wrapper">
    <!-- 根据当前显示的表单类型渲染对应的组件 -->
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Login from './Login.vue'
import Register from './Register.vue'
import { authService } from '../services/auth.js'

// 定义props和emits
const emit = defineEmits(['auth-success'])

// 状态管理
const showLogin = ref(true) // 默认显示登录表单

// 初始化时检查用户是否已登录
onMounted(() => {
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
  
  // 这里可以添加一个短暂的提示，告知用户注册成功，请登录
  // 为了简化，这里直接切换表单，实际项目中可以添加toast提示
}
</script>

<style scoped>
.auth-container-wrapper {
  width: 100%;
  min-height: 100vh;
}
</style>