<template>
  <div class="teams-center">
    <div class="card">
      <h3>我的协作团队</h3>
      <!-- 新增：搜索与统计工具栏 -->
      <div class="toolbar">
        <input class="search" v-model="keyword" placeholder="搜索团队名称..." />
        <span class="count">共 {{ filteredTeams.length }} 个团队</span>
      </div>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-else>
        <div v-if="filteredTeams.length === 0" class="empty">暂无团队</div>
        <ul v-else class="team-list">
          <li v-for="team in filteredTeams" :key="team.id" class="team-item">
            <div class="line">
              <span class="name">{{ team.name }}</span>
              <span class="meta">成员数：{{ team.members ? parseMembers(team.members).length : 0 }}</span>
              <span class="meta">负责人ID：{{ team.owner_id }}</span>
            </div>
            <div class="desc" v-if="team.description">{{ team.description }}</div>
            <div class="time">创建于 {{ formatDate(team.created_at) }}</div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/auth'

export default {
  name: 'TeamsCenter',
  data() {
    return {
      loading: true,
      teams: [],
      errorMessage: '',
      // 新增：搜索关键字
      keyword: ''
    }
  },
  mounted() {
    this.loadTeams()
  },
  computed: {
    // 新增：搜索过滤
    filteredTeams() {
      const kw = (this.keyword || '').trim().toLowerCase()
      if (!kw) return this.teams
      return this.teams.filter(t => (t.name || '').toLowerCase().includes(kw))
    }
  },
  methods: {
    async loadTeams() {
      this.loading = true
      this.errorMessage = ''
      try {
        const res = await api.get('/api/v1/teams/')
        this.teams = Array.isArray(res) ? res : []
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '获取团队列表失败'
      } finally {
        this.loading = false
      }
    },
    parseMembers(membersStr) {
      try {
        // 后端 Team.members 字段描述为“团队成员列表（用户名形式）”，采用逗号或 JSON 数组存储均做兼容
        if (!membersStr) return []
        if (membersStr.trim().startsWith('[')) {
          return JSON.parse(membersStr)
        }
        return membersStr.split(',').map(s => s.trim()).filter(Boolean)
      } catch (_) {
        return []
      }
    },
    formatDate(dt) {
      if (!dt) return '-'
      try { return new Date(dt).toLocaleString() } catch (_) { return dt }
    }
  }
}
</script>

<style scoped>
.teams-center { display: flex; flex-direction: column; gap: 16px; }
.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; padding: 16px; box-shadow: 0 6px 24px rgba(0,0,0,0.06); }
.card h3 { font-size: 1.1rem; margin-bottom: 12px; color: #333; }
/* 新增：工具栏样式 */
.toolbar { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 10px; }
.search { flex: 1; padding: 10px 12px; border: 1px solid #d9d9e3; border-radius: 8px; font-size: 14px; transition: border-color .2s ease, box-shadow .2s ease; }
.search:focus { outline: none; border-color: #667eea; box-shadow: 0 0 0 3px rgba(102,126,234,.2); }
.count { color: #666; font-size: 13px; }
.loading { color: #666; }
.error-message { background: #fdecea; color: #d93025; border: 1px solid #f2a19a; padding: 8px 10px; border-radius: 8px; }
.empty { color: #888; }
.team-list { list-style: none; padding: 0; margin: 0; }
.team-item { border-bottom: 1px dashed #eee; padding: 10px 0; transition: background .15s ease; }
.team-item:hover { background: #fafbff; }
.line { display: flex; gap: 8px; align-items: center; justify-content: space-between; }
.name { font-weight: 600; color: #333; }
.meta { color: #666; font-size: 13px; }
.desc { margin-top: 6px; color: #555; }
.time { margin-top: 6px; color: #888; font-size: 12px; }
</style>