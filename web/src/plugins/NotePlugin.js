// Note插件

import { noteService } from '../services/note.js'

class NotePlugin {
  constructor() {
    this.name = 'NotePlugin'
    this.version = '1.0.0'
    this.description = '一个简单的笔记插件'
    this.notes = []
  }

  // 初始化插件
  async initialize() {
    console.log('NotePlugin 初始化开始')
    // 从后端API加载笔记
    try {
      await this.loadNotesFromAPI()
      console.log('NotePlugin 初始化完成，当前笔记数量:', this.notes.length)
    } catch (error) {
      console.error('笔记插件初始化失败:', error)
    }
  }

  // 获取插件信息
  getInfo() {
    return {
      name: this.name,
      version: this.version,
      description: this.description,
      noteCount: this.notes.length
    }
  }

  // 添加笔记
  async addNote(title, content) {
    try {
      // noteService.createNote会返回经过auth.js响应拦截器处理后的response.data
      const result = await noteService.createNote({ title, content })
      // 更新本地笔记列表
      await this.loadNotesFromAPI()
      return result
    } catch (error) {
      console.error('添加笔记失败:', error)
      throw error
    }
  }

  // 获取所有笔记
  getAllNotes() {
    return this.notes
  }

  // 从后端API加载笔记
  async loadNotesFromAPI() {
    try {
      const data = await noteService.getAllNotes()
      // 检查响应格式，根据auth.js中响应拦截器的行为调整
      if (data && Array.isArray(data)) {
        // 处理直接返回的笔记数组
        this.notes = data
      } else if (data && data.notes && Array.isArray(data.notes)) {
        // 处理嵌套在notes字段中的笔记数组
        this.notes = data.notes
      } else {
        console.warn('Unexpected response format:', data)
        this.notes = []
      }
      console.log('Loaded notes:', this.notes)
    } catch (error) {
      console.error('从API加载笔记失败:', error)
      this.notes = []
    }
  }

  // 删除笔记
  async deleteNote(id) {
    try {
      await noteService.deleteNote(id)
      // 更新本地笔记列表
      this.loadNotesFromAPI()
    } catch (error) {
      console.error('删除笔记失败:', error)
      throw error
    }
  }



  // 渲染插件内容
  render() {
    return {
      template: `<div class="plugin-note">
                  <h3>📝 笔记插件</h3>
                  <div class="note-form">
                    <input v-model="newNoteTitle" placeholder="笔记标题" type="text">
                    <textarea v-model="newNoteContent" placeholder="笔记内容"></textarea>
                    <button @click="addNewNote">添加笔记</button>
                  </div>
                  <div class="notes-list">
                    <div v-for="note in notes" :key="note.id" class="note-item">
                      <h4>{{ note.title }}</h4>
                      <p>{{ note.content }}</p>
                      <small>{{ formatDate(note.created_time) }}</small>
                      <button @click="deleteNoteItem(note.id)">删除</button>
                    </div>
                  </div>
                </div>`,
      data: function() {
        // 使用外部this引用或直接从插件实例获取数据
        const pluginInstance = this.plugin || window.notePluginInstance || this;
        return {
          newNoteTitle: '',
          newNoteContent: '',
          notes: pluginInstance.getAllNotes ? pluginInstance.getAllNotes() : []
        }
      },
      methods: {
        addNewNote: async function() {
          if (this.newNoteTitle.trim()) {
            try {
              await this.addNote(this.newNoteTitle, this.newNoteContent)
              this.newNoteTitle = ''
              this.newNoteContent = ''
              this.notes = this.getAllNotes()
            } catch (error) {
              console.error('添加笔记失败:', error)
              alert('添加笔记失败，请稍后重试')
            }
          }
        },
        deleteNoteItem: async function(id) {
          try {
            await this.deleteNote(id)
            this.notes = this.getAllNotes()
          } catch (error) {
            console.error('删除笔记失败:', error)
            alert('删除笔记失败，请稍后重试')
          }
        },
        formatDate: function(dateString) {
          return new Date(dateString).toLocaleString()
        }
      },
      watch: {
        notes: function(newNotes) {
          this.notes = newNotes
        }
      },
      css: `.plugin-note {
              padding: 1rem;
              border-radius: 8px;
              background-color: #fff9e6;
              border: 1px solid #ddd;
            }
            .plugin-note h3 {
              margin-top: 0;
              color: #333;
            }
            .note-form {
              margin-bottom: 1rem;
            }
            .note-form input,
            .note-form textarea {
              width: 100%;
              padding: 0.5rem;
              margin-bottom: 0.5rem;
              border: 1px solid #ddd;
              border-radius: 4px;
            }
            .note-form textarea {
              height: 80px;
              resize: vertical;
            }
            .notes-list {
              max-height: 300px;
              overflow-y: auto;
            }
            .note-item {
              padding: 0.8rem;
              margin-bottom: 0.5rem;
              background-color: white;
              border: 1px solid #eee;
              border-radius: 4px;
            }
            .note-item h4 {
              margin-top: 0;
              margin-bottom: 0.3rem;
              color: #333;
            }
            .note-item small {
              color: #888;
              display: block;
              margin-bottom: 0.5rem;
            }`
    }
  }

  // 销毁插件
  destroy() {
    console.log('NotePlugin 已销毁')
    // 这里可以添加插件的清理逻辑
  }
}

export default NotePlugin