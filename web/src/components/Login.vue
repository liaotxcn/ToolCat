<template>
  <div class="login-container">
    <h2 class="form-title">登录</h2>
    <form class="auth-form" @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">用户名</label>
        <input v-model="username" type="text" id="username" required placeholder="请输入用户名" autofocus @input="clearError" />
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <div class="password-wrap">
          <input :type="showPassword ? 'text' : 'password'" v-model="password" id="password" required placeholder="请输入密码" @input="clearError" />
          <button type="button" class="toggle-psw" @click="showPassword = !showPassword" :aria-pressed="showPassword" :title="showPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
      </div>
      <div class="assist">
        <label class="remember"><input type="checkbox" v-model="rememberMe" /> 记住我</label>
        <button class="link-btn" type="button" title="暂未实现接口">忘记密码？</button>
      </div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <button class="primary-btn" type="submit" :disabled="loading || !canLogin">{{ loading ? '登录中...' : '登录' }}</button>
    </form>
    <p class="switch-tip">
      还没有账号？
      <button class="link-btn" type="button" @click="switchToRegister">立即注册</button>
    </p>
  </div>
</template>

<script>
import { authService } from '../services/auth'

export default {
  name: 'Login',
  emits: ['login-success', 'switch-to-register'],
  data() {
    return {
      username: '',
      password: '',
      rememberMe: true,
      showPassword: false,
      loading: false,
      errorMessage: ''
    }
  },
  computed: {
    usernameInvalid() {
      return !(this.username && this.username.trim().length > 0)
    },
    passwordInvalid() {
      return !(this.password && this.password.length > 0)
    },
    canLogin() {
      return !this.usernameInvalid && !this.passwordInvalid
    }
  },
  methods: {
    clearError() {
      this.errorMessage = ''
    },
    async handleLogin() {
      if (!this.canLogin) return
      this.errorMessage = ''
      try {
        this.loading = true
        const payload = {
          username: this.username.trim(),
          password: this.password,
          remember_me: this.rememberMe
        }

        const response = await authService.login(payload)

        if (response && response.user) {
          this.$emit('login-success', response.user)
        } else {
          this.errorMessage = response?.message || '登录失败，请稍后重试'
        }
      } catch (error) {
        const data = error?.response?.data || {}
        this.errorMessage = data?.message || '登录失败，请检查账号或网络'
      } finally {
        this.loading = false
      }
    },
    switchToRegister() {
      this.$emit('switch-to-register')
    }
  }
}
</script>

<style scoped>
.login-container { display: flex; flex-direction: column; gap: 12px; }
.form-title { margin: 0; }
.auth-form { display: flex; flex-direction: column; gap: 12px; }
.form-group label { font-weight: 500; color: #333; margin-bottom: 6px; }
.form-group input { width: 100%; padding: 10px 12px; border: 1px solid #d9d9e3; border-radius: 8px; font-size: 14px; transition: border-color .2s ease, box-shadow .2s ease; }
.form-group input:focus { outline: none; border-color: #667eea; box-shadow: 0 0 0 3px rgba(102,126,234,.2); }
.password-wrap { position: relative; }
.toggle-psw { position: absolute; right: 8px; top: 50%; transform: translateY(-50%); border: 1px solid #e6e6f2; background: #fff; color: #555; padding: 4px; border-radius: 6px; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; }
.toggle-psw:hover { background: #f8f9fa; }
.eye-icon { display: block; }
.assist { display: flex; align-items: center; justify-content: space-between; }
.remember { color: #555; font-size: 14px; }
.primary-btn { padding: 10px 14px; border: none; border-radius: 8px; background: linear-gradient(135deg, #667eea 0%, #5a67d8 100%); color: #fff; cursor: pointer; font-size: 15px; transition: transform .05s ease, box-shadow .2s ease; }
.primary-btn:hover { box-shadow: 0 8px 20px rgba(102,126,234,.35); }
.primary-btn:active { transform: translateY(1px); }
.primary-btn:disabled { opacity: .7; cursor: not-allowed; }
.error-message { background: #fdecea; color: #d93025; border: 1px solid #f2a19a; padding: 8px 10px; border-radius: 8px; }
.switch-tip { text-align: center; margin-top: 8px; color: #666; }
.link-btn { background: none; border: none; color: #667eea; cursor: pointer; text-decoration: underline; padding: 0 4px; }
.link-btn:hover { color: #5a67d8; }
</style>