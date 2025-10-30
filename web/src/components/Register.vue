<template>
  <div class="register-container">
    <h2 class="form-title">注册</h2>
    <form class="auth-form" @submit.prevent="handleRegister">
      <div class="form-group">
        <label for="username">用户名</label>
        <input v-model="username" type="text" id="username" required minlength="3" maxlength="50" placeholder="3-50个字符" @input="clearError" />
      </div>
      <div class="form-group">
        <label for="email">邮箱</label>
        <input v-model="email" type="email" id="email" required placeholder="name@example.com" @input="clearError" />
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <div class="password-wrap">
          <input :type="showPassword ? 'text' : 'password'" v-model="password" id="password" required minlength="6" placeholder="至少6个字符" @input="clearError" />
          <button type="button" class="toggle-psw" @click="showPassword = !showPassword" :aria-pressed="showPassword" :title="showPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
        <div class="pw-strength">
          <div class="bar" :class="passwordLevel"></div>
          <span class="level">{{ passwordLabel }}</span>
        </div>
      </div>
      <div class="form-group">
        <label for="confirmPassword">确认密码</label>
        <div class="password-wrap">
          <input :type="showConfirmPassword ? 'text' : 'password'" v-model="confirmPassword" id="confirmPassword" required minlength="6" placeholder="再次输入密码" @input="clearError" />
          <button type="button" class="toggle-psw" @click="showConfirmPassword = !showConfirmPassword" :aria-pressed="showConfirmPassword" :title="showConfirmPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="20" height="20" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showConfirmPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
      </div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <button class="primary-btn" type="submit" :disabled="loading || !canRegister">{{ loading ? '注册中...' : '注册' }}</button>
    </form>
    <p class="switch-tip">
      已有账号？
      <button class="link-btn" type="button" @click="switchToLogin">返回登录</button>
    </p>
  </div>
</template>

<script>
import { authService } from '../services/auth'

export default {
  name: 'Register',
  emits: ['register-success', 'switch-to-login'],
  data() {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      showPassword: false,
      showConfirmPassword: false,
      loading: false,
      errorMessage: ''
    }
  },
  computed: {
    usernameInvalid() {
      return !(this.username && this.username.trim().length >= 3)
    },
    emailInvalid() {
      if (!this.email) return false
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return !re.test(this.email)
    },
    passwordLevel() {
      const n = this.password || ''
      let score = 0
      if (n.length >= 6) score++
      if (/[A-Z]/.test(n)) score++
      if (/[a-z]/.test(n)) score++
      if (/\d/.test(n)) score++
      if (/[^\w]/.test(n)) score++
      if (score <= 2) return 'weak'
      if (score === 3 || score === 4) return 'medium'
      return 'strong'
    },
    passwordLabel() {
      return this.passwordLevel === 'strong' ? '强' : (this.passwordLevel === 'medium' ? '中' : '弱')
    },
    passwordMismatch() {
      return !!this.confirmPassword && this.password !== this.confirmPassword
    },
    canRegister() {
      return !this.usernameInvalid && !this.emailInvalid && !!this.password && this.password.length >= 6 && !this.passwordMismatch
    }
  },
  methods: {
    clearError() {
      this.errorMessage = ''
    },
    async handleRegister() {
      if (!this.canRegister) return
      this.errorMessage = ''
      try {
        this.loading = true
        const payload = {
          username: this.username.trim(),
          email: this.email.trim(),
          password: this.password,
          confirm_password: this.confirmPassword
        }

        const response = await authService.register(payload)

        if (response && response.user) {
          this.$emit('register-success', response.user)
        } else {
          this.errorMessage = response?.message || '注册失败，请稍后重试'
        }
      } catch (error) {
        const data = error?.response?.data || {}
        this.errorMessage = data?.message || '注册失败，请检查输入或网络'
      } finally {
        this.loading = false
      }
    },
    switchToLogin() {
      this.$emit('switch-to-login')
    }
  }
}
</script>

<style scoped>
.register-container { display: flex; flex-direction: column; gap: 12px; }
.form-title { margin: 0; }
.auth-form { display: flex; flex-direction: column; gap: 12px; }
.form-group label { font-weight: 500; color: #333; margin-bottom: 6px; }
.form-group input { width: 100%; padding: 10px 12px; border: 1px solid #d9d9e3; border-radius: 8px; font-size: 14px; transition: border-color .2s ease, box-shadow .2s ease; }
.form-group input:focus { outline: none; border-color: #667eea; box-shadow: 0 0 0 3px rgba(102,126,234,.2); }
.password-wrap { position: relative; }
.toggle-psw { position: absolute; right: 8px; top: 50%; transform: translateY(-50%); border: 1px solid #e6e6f2; background: #fff; color: #555; padding: 4px; border-radius: 6px; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; }
.toggle-psw:hover { background: #f8f9fa; }
.eye-icon { display: block; }
.pw-strength { display: flex; align-items: center; gap: 8px; }
.pw-strength .bar { height: 6px; width: 80px; border-radius: 10px; background: #eee; position: relative; overflow: hidden; }
.pw-strength .bar::after { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 33%; background: #f59e0b; transition: width .2s ease, background .2s ease; }
.pw-strength .bar.medium::after { width: 66%; background: #fbbf24; }
.pw-strength .bar.strong::after { width: 100%; background: #22c55e; }
.pw-strength .level { color: #666; font-size: 12px; }
.primary-btn { padding: 10px 14px; border: none; border-radius: 8px; background: linear-gradient(135deg, #667eea 0%, #5a67d8 100%); color: #fff; cursor: pointer; font-size: 15px; transition: transform .05s ease, box-shadow .2s ease; }
.primary-btn:hover { box-shadow: 0 8px 20px rgba(102,126,234,.35); }
.primary-btn:active { transform: translateY(1px); }
.primary-btn:disabled { opacity: .7; cursor: not-allowed; }
.error-message { background: #fdecea; color: #d93025; border: 1px solid #f2a19a; padding: 8px 10px; border-radius: 8px; }
.switch-tip { text-align: center; margin-top: 8px; color: #666; }
.link-btn { background: none; border: none; color: #667eea; cursor: pointer; text-decoration: underline; padding: 0 4px; }
.link-btn:hover { color: #5a67d8; }
</style>