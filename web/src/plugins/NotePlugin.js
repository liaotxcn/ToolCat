// Note插件

class NotePlugin {
  constructor() {
    this.name = 'NotePlugin'
    this.version = '1.0.0'
    this.description = '一个简单的笔记插件'
    this.notes = []
  }

  // 初始化插件
  initialize() {
    console.log('NotePlugin 初始化')
    // 从本地存储加载笔记
    this.loadNotes()
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
  addNote(title, content) {
    const note = {
      id: Date.now(),
      title: title,
      content: content,
      createdAt: new Date().toISOString()
    }
    this.notes.push(note)
    this.saveNotes()
    return note
  }

  // 获取所有笔记
  getAllNotes() {
    return this.notes
  }

  // 删除笔记
  deleteNote(id) {
    this.notes = this.notes.filter(note => note.id !== id)
    this.saveNotes()
  }

  // 保存笔记到本地存储
  saveNotes() {
    try {
      localStorage.setItem('toolcat-notes', JSON.stringify(this.notes))
    } catch (error) {
      console.error('保存笔记失败:', error)
    }
  }

  // 从本地存储加载笔记
  loadNotes() {
    try {
      const savedNotes = localStorage.getItem('toolcat-notes')
      if (savedNotes) {
        this.notes = JSON.parse(savedNotes)
      }
    } catch (error) {
      console.error('加载笔记失败:', error)
      this.notes = []
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
                      <small>{{ formatDate(note.createdAt) }}</small>
                      <button @click="deleteNoteItem(note.id)">删除</button>
                    </div>
                  </div>
                </div>`,
      data: function() {
        return {
          newNoteTitle: '',
          newNoteContent: '',
          notes: this.getAllNotes()
        }
      },
      methods: {
        addNewNote: function() {
          if (this.newNoteTitle.trim()) {
            this.addNote(this.newNoteTitle, this.newNoteContent)
            this.newNoteTitle = ''
            this.newNoteContent = ''
            this.notes = this.getAllNotes()
          }
        },
        deleteNoteItem: function(id) {
          this.deleteNote(id)
          this.notes = this.getAllNotes()
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