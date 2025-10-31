<template>
  <div class="user-center">
    <div class="cards">
      <section class="card">
        <h3>基本资料</h3>
        <div class="profile-row"><span class="label">用户名</span><span class="value">{{ user?.username }}</span></div>
        <div class="profile-row"><span class="label">邮箱</span>
          <input v-model="email" type="email" placeholder="name@example.com" @input="clearStatus" />
        </div>
        <!-- 新增：邮箱格式校验提示 -->
        <div v-if="emailInvalid" class="field-tip error">请输入有效邮箱地址</div>
        <div class="profile-row"><span class="label">创建时间</span><span class="value">{{ formatDate(user?.created_at) }}</span></div>
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
        <button class="primary-btn" @click="updateProfile" :disabled="updating || !canSave">{{ updating ? '保存中...' : '保存资料' }}</button>
      </section>

      <section class="card">
        <h3>安全设置</h3>
        <div class="profile-row"><span class="label">新密码</span>
          <input v-model="newPassword" type="password" placeholder="至少6个字符" @input="clearStatus" />
        </div>
        <!-- 新增：密码强度提示条 -->
        <div class="pw-strength" v-if="newPassword">
          <div class="bar" :class="passwordLevel"></div>
          <span class="level">{{ passwordLabel }}</span>
        </div>
        <div class="profile-row"><span class="label">确认密码</span>
          <input v-model="confirmNewPassword" type="password" placeholder="再次输入新密码" @input="clearStatus" />
        </div>
        <!-- 新增：两次密码不一致提示 -->
        <div v-if="passwordMismatch" class="field-tip error">两次输入的密码不一致</div>
        <button class="secondary-btn" @click="updatePassword" :disabled="updating || !canUpdatePassword">{{ updating ? '提交中...' : '更新密码' }}</button>
      </section>
    </div>

    <section class="card">
      <h3>近期活动（审计日志）</h3>
      <div v-if="auditLogs.length === 0" class="empty">暂无记录</div>
      <ul v-else class="audit-list">
        <li v-for="log in auditLogs" :key="log.id" class="audit-item">
          <div class="audit-line">
            <span class="tag">{{ log.action }}</span>
            <span class="resource">{{ log.resource_type }} #{{ log.resource_id }}</span>
            <span class="time">{{ formatDate(log.created_at) }}</span>
          </div>
          <div class="detail" v-if="log.new_value || log.old_value">
            <span v-if="log.old_value">变更前: {{ short(log.old_value) }}</span>
            <span v-if="log.new_value">变更后: {{ short(log.new_value) }}</span>
          </div>
        </li>
      </ul>
    </section>
  </div>
</template>

<script>
import api, { authService } from '../services/auth'

export default {
  name: 'UserCenter',
  props: {
    currentUser: { type: Object, default: null }
  },
  data() {
    return {
      user: null,
      email: '',
      newPassword: '',
      confirmNewPassword: '',
      auditLogs: [],
      updating: false,
      errorMessage: '',
      successMessage: ''
    }
  },
  mounted() {
    this.loadUser()
  },
  computed: {
    // 新增：是否允许保存资料
    canSave() {
      return !!this.email && !this.emailInvalid
    },
    // 新增：邮箱是否非法
    emailInvalid() {
      if (!this.email) return false
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return !re.test(this.email)
    },
    // 新增：密码强度等级样式
    passwordLevel() {
      const n = this.newPassword || ''
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
    // 新增：密码强度文案
    passwordLabel() {
      return this.passwordLevel === 'strong' ? '强' : (this.passwordLevel === 'medium' ? '中' : '弱')
    },
    // 新增：两次密码是否不一致
    passwordMismatch() {
      return !!this.confirmNewPassword && this.newPassword !== this.confirmNewPassword
    },
    // 新增：是否允许更新密码
    canUpdatePassword() {
      return !!this.newPassword && this.newPassword.length >= 6 && !this.passwordMismatch
    }
  },
  methods: {
    // 新增：输入时清除顶部成功/错误提示，避免阻碍表单操作
    clearStatus() {
      this.errorMessage = ''
      this.successMessage = ''
    },
    async loadUser() {
      try {
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) {
          this.errorMessage = '未获取到当前用户信息，请重新登录'
          return
        }
        const res = await api.get(`/api/v1/users/${cur.id}`)
        this.user = res
        this.email = res.email || ''
        this.auditLogs = Array.isArray(res.audit_logs) ? res.audit_logs.slice(0, 20) : []
        this.errorMessage = ''
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '加载用户信息失败'
      }
    },
    async updateProfile() {
      try {
        this.updating = true
        this.successMessage = ''
        this.errorMessage = ''
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) throw new Error('未登录')
        const payload = { email: this.email }
        const updated = await api.put(`/api/v1/users/${cur.id}`, payload)
        this.user = updated
        this.$emit('updated-user', updated)
        this.successMessage = '资料已更新'
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '更新失败'
      } finally {
        this.updating = false
      }
    },
    async updatePassword() {
      try {
        this.updating = true
        this.successMessage = ''
        this.errorMessage = ''
        if (!this.newPassword || this.newPassword.length < 6) {
          this.errorMessage = '新密码至少6个字符'
          return
        }
        if (this.newPassword !== this.confirmNewPassword) {
          this.errorMessage = '两次输入的密码不一致'
          return
        }
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) throw new Error('未登录')
        const payload = { password: this.newPassword }
        const updated = await api.put(`/api/v1/users/${cur.id}`, payload)
        this.user = updated
        this.$emit('updated-user', updated)
        this.newPassword = ''
        this.confirmNewPassword = ''
        this.successMessage = '密码已更新'
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '密码更新失败'
      } finally {
        this.updating = false
      }
    },
    formatDate(dt) {
      if (!dt) return '-'
      try { return new Date(dt).toLocaleString() } catch (_) { return dt }
    },
    short(text) {
      if (!text) return ''
      const s = String(text)
      return s.length > 180 ? s.slice(0, 180) + '…' : s
    }
  }
}
</script>

<style scoped>
.user-center { display: flex; flex-direction: column; gap: 16px; }
.cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 16px; }
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; padding: 16px; box-shadow: 0 6px 24px rgba(0,0,0,0.06); }
.card h3 { font-size: 1.1rem; margin-bottom: 12px; color: #333; }
.profile-row { display: grid; grid-template-columns: 100px 1fr; align-items: center; gap: 8px; margin-bottom: 10px; }
.label { color: #666; }
.value { color: #333; }
input { width: 100%; padding: 10px 12px; border: 1px solid #d9d9e3; border-radius: 8px; font-size: 14px; transition: border-color .2s ease, box-shadow .2s ease; }
input:focus { outline: none; border-color: #667eea; box-shadow: 0 0 0 3px rgba(102,126,234,.2); }
.primary-btn { padding: 10px 14px; border: none; border-radius: 8px; background: linear-gradient(135deg, #667eea 0%, #5a67d8 100%); color: #fff; cursor: pointer; font-size: 15px; transition: transform .05s ease, box-shadow .2s ease; }
.primary-btn:hover { box-shadow: 0 8px 20px rgba(102,126,234,.35); }
.primary-btn:active { transform: translateY(1px); }
.secondary-btn { padding: 10px 14px; border: 1px solid #d9d9e3; border-radius: 8px; background: #fff; color: #333; cursor: pointer; font-size: 15px; transition: all .2s ease; }
.secondary-btn:hover { background: #f8f9fa; }
.error-message { background: #fdecea; color: #d93025; border: 1px solid #f2a19a; padding: 8px 10px; border-radius: 8px; margin-top: 8px; }
.success-message { background: #e7f5ee; color: #2f9e44; border: 1px solid #b6e3c6; padding: 8px 10px; border-radius: 8px; margin-top: 8px; }
.audit-list { list-style: none; padding: 0; margin: 0; }
.audit-item { border-bottom: 1px dashed #eee; padding: 8px 0; }
.audit-line { display: flex; gap: 8px; align-items: center; justify-content: space-between; }
.tag { background: #eef2ff; color: #4f46e5; border: 1px solid #e0e7ff; padding: 2px 8px; border-radius: 12px; font-size: 12px; }
.resource { color: #555; }
.time { color: #888; font-size: 12px; }
.detail { margin-top: 6px; color: #666; font-size: 13px; }
.empty { color: #888; }
/* 新增：字段提示与密码强度样式 */
.field-tip { margin: -4px 0 8px 100px; font-size: 12px; }
.field-tip.error { color: #d93025; }
.pw-strength { display: flex; align-items: center; gap: 8px; margin: -4px 0 8px 100px; }
.pw-strength .bar { height: 6px; width: 80px; border-radius: 10px; background: #eee; position: relative; overflow: hidden; }
.pw-strength .bar::after { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 33%; background: #f59e0b; transition: width .2s ease, background .2s ease; }
.pw-strength .bar.medium::after { width: 66%; background: #fbbf24; }
.pw-strength .bar.strong::after { width: 100%; background: #22c55e; }
.pw-strength .level { color: #666; font-size: 12px; }
</style>