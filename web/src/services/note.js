import { api } from './auth.js'

// 笔记相关API方法
export const noteService = {
  // 获取所有笔记
  getAllNotes: async () => {
    try {
      const response = await api.get('/plugins/note/notes')
      return response
    } catch (error) {
      console.error('获取笔记失败:', error)
      throw error
    }
  },

  // 获取单个笔记
  getNote: async (id) => {
    try {
      const response = await api.get(`/plugins/note/notes/${id}`)
      return response
    } catch (error) {
      console.error('获取笔记失败:', error)
      throw error
    }
  },

  // 创建新笔记
  createNote: async (noteData) => {
    try {
      const response = await api.post('/plugins/note/notes', noteData)
      return response
    } catch (error) {
      console.error('创建笔记失败:', error)
      throw error
    }
  },

  // 更新笔记
  updateNote: async (id, noteData) => {
    try {
      const response = await api.put(`/plugins/note/notes/${id}`, noteData)
      return response
    } catch (error) {
      console.error('更新笔记失败:', error)
      throw error
    }
  },

  // 删除笔记
  deleteNote: async (id) => {
    try {
      const response = await api.delete(`/plugins/note/notes/${id}`)
      return response
    } catch (error) {
      console.error('删除笔记失败:', error)
      throw error
    }
  }
}

export default noteService