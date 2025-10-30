<template>
  <div class="plugin-renderer">
    <div v-if="loading" class="plugin-loading fade-in">
      <div class="spinner"></div>
      <span>加载插件中...</span>
    </div>
    <div v-else-if="error" class="plugin-error fade-in">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button @click="refreshPlugin" class="retry-btn">重试</button>
    </div>
    <div v-else-if="pluginComponent" class="plugin-content fade-in">
      <div class="plugin-toolbar">
        <div class="plugin-meta">
          <span class="plugin-name">{{ pluginInfo?.name || '插件' }}</span>
          <span v-if="pluginInfo?.description" class="plugin-sep">·</span>
          <span v-if="pluginInfo?.description" class="plugin-desc">{{ pluginInfo.description }}</span>
          <span v-if="pluginInfo?.version" class="plugin-version">v{{ pluginInfo.version }}</span>
        </div>
        <div class="plugin-actions">
          <button class="toolbar-btn" @click="refreshPlugin">刷新</button>
        </div>
      </div>
      <div class="plugin-body">
        <component :is="pluginComponent" :plugin="plugin" ref="pluginContainer"></component>
      </div>
    </div>
    <div v-else class="plugin-empty fade-in">
      <div class="empty-icon">📦</div>
      <p>请选择一个插件/服务</p>
    </div>
  </div>
</template>

<script setup>
import * as VueRuntimeDOM from 'vue'
import { ref, watch, onMounted, defineComponent } from 'vue'
import { compile as compileTemplate } from '@vue/compiler-dom'

const props = defineProps({
  pluginName: {
    type: String,
    required: true
  },
  pluginManager: {
    type: Object,
    required: true
  }
})

const plugin = ref(null)
const pluginComponent = ref(null)
const pluginContainer = ref(null)
const loading = ref(false)
const error = ref(null)
const pluginInfo = ref(null)

// 辅助函数：更新组件数据
const updateComponentData = async (pluginInstance, component) => {
  if (pluginInstance.loadNotesFromAPI && typeof pluginInstance.loadNotesFromAPI === 'function') {
    await pluginInstance.loadNotesFromAPI()
    if (component.$data.notes !== undefined && pluginInstance.getAllNotes) {
      component.$data.notes = [...(pluginInstance.getAllNotes() || [])]
    }
  }
}

// 加载插件
const loadPlugin = async () => {
  if (!props.pluginName || !props.pluginManager) {
    plugin.value = null
    pluginComponent.value = null
    return
  }

  loading.value = true
  error.value = null

  try {
    // 从pluginManager获取插件实例
    plugin.value = props.pluginManager.getPlugin(props.pluginName)
    if (!plugin.value) {
      throw new Error(`Plugin ${props.pluginName} not found`)
    }

    // 调用插件的初始化方法
    if (typeof plugin.value.initialize === 'function') {
      await plugin.value.initialize()
    }

    // 获取插件的渲染结果
    const renderResult = plugin.value.render()

    // 动态创建Vue组件
    if (renderResult && renderResult.template) {
      // 创建方法映射
      const pluginMethods = {}
      if (renderResult.methods) {
        Object.keys(renderResult.methods).forEach(key => {
          if (typeof renderResult.methods[key] === 'function') {
            pluginMethods[key] = function(...args) {
              const result = renderResult.methods[key].apply(this, args)
              if (result && typeof result.then === 'function') {
                return result.then(async (resolvedResult) => {
                  await updateComponentData(plugin.value, this)
                  return resolvedResult
                })
              } else {
                updateComponentData(plugin.value, this)
                return result
              }
            }
          }
        })
      }

      // 运行时编译模板（使用 function 模式），失败时回退到 template
      let renderFn = null
      try {
        if (compileTemplate) {
          const { code } = compileTemplate(renderResult.template, { mode: 'function' })
          renderFn = new Function('Vue', code)(VueRuntimeDOM)
        }
      } catch (e) {
        console.warn('Runtime compile failed, fallback to template option:', e)
      }

      pluginComponent.value = defineComponent({
        name: `${props.pluginName}-component`,
        props: { plugin: Object },
        // 兼容：优先使用 render 函数，否则使用 template 字符串
        ...(renderFn ? { render: renderFn } : { template: renderResult.template }),
        data() {
          if (renderResult.data && typeof renderResult.data === 'function') {
            const dataResult = renderResult.data.call(plugin.value)
            if (plugin.value.getAllNotes && !dataResult.notes) {
              dataResult.notes = [...(plugin.value.getAllNotes() || [])]
            }
            return dataResult
          }
          return {}
        },
        methods: {
          escapeHtml: (text) => { const div = document.createElement('div'); div.textContent = text; return div.innerHTML },
          formatDate: (dateString) => { try { return new Date(dateString).toLocaleString() } catch (e) { return dateString } },
          ...pluginMethods
        },
        computed: renderResult.computed || {},
        watch: renderResult.watch || {},
        created() {
          const pluginInstance = this.plugin
          if (pluginInstance) {
            ;['addNote', 'updateNote', 'deleteNote', 'getAllNotes', 'loadNotesFromAPI'].forEach(methodName => {
              if (typeof pluginInstance[methodName] === 'function') {
                this[methodName] = pluginInstance[methodName].bind(pluginInstance)
              }
            })
          }
        },
        mounted() {
          if (this.loadNotesFromAPI && typeof this.loadNotesFromAPI === 'function') {
            this.loadNotesFromAPI().then(() => {
              if (this.getAllNotes && typeof this.getAllNotes === 'function') {
                this.notes = [...(this.getAllNotes() || [])]
              }
            })
          }
        }
      })
    } else {
      throw new Error('插件未返回有效的模板')
    }

    // 处理插件样式
    if (renderResult.css) {
      loadPluginCSS(renderResult.css)
    }
  } catch (err) {
    console.error('加载插件失败:', err)
    error.value = `加载插件失败: ${err.message || '未知错误'}`
    pluginComponent.value = null
  } finally {
    loading.value = false
  }
}

// 加载插件CSS
const loadPluginCSS = (css) => {
  if (!css) return
  const styleId = `plugin-css-${props.pluginName}`
  let styleElement = document.getElementById(styleId)
  if (!styleElement) {
    styleElement = document.createElement('style')
    styleElement.id = styleId
    document.head.appendChild(styleElement)
  }
  styleElement.textContent = css
}

// 监听pluginName变化
watch(() => props.pluginName, () => { loadPlugin() }, { immediate: true })

// 导出方法供父组件使用
defineExpose({ refreshPlugin: loadPlugin })

// 组件挂载时加载插件
onMounted(() => { loadPlugin() })
const refreshPlugin = () => { loadPlugin() }
</script>

<style scoped>
.plugin-renderer { 
  width: 100%; 
  min-height: 100%; 
  padding: var(--space-4, 16px); 
  box-sizing: border-box;
  position: relative;
}

/* 加载状态 */
.plugin-loading { 
  display: flex; 
  flex-direction: column;
  align-items: center; 
  justify-content: center; 
  min-height: 300px;
  color: var(--text-tertiary, #666); 
  gap: var(--space-3, 12px);
  padding: var(--space-6, 24px);
  background: var(--bg-secondary, #f8fafc);
  border-radius: var(--radius-lg, 10px);
  transition: var(--transition-all, all 0.3s ease);
}

.spinner { 
  width: 40px; 
  height: 40px; 
  border: 3px solid var(--border-light, rgba(102,126,234,0.25)); 
  border-top-color: var(--primary-500, #667eea); 
  border-radius: 50%; 
  animation: spin 1s linear infinite; 
}

@keyframes spin { 
  to { transform: rotate(360deg); } 
}

/* 错误状态 */
.plugin-error { 
  display: flex; 
  flex-direction: column;
  align-items: center; 
  justify-content: center; 
  min-height: 300px;
  color: var(--error-700, #f44336); 
  gap: var(--space-3, 12px);
  padding: var(--space-6, 24px);
  text-align: center;
  background: var(--error-100, #fef2f2);
  border: 1px solid var(--error-200, #fecaca);
  border-radius: var(--radius-lg, 10px);
  transition: var(--transition-all, all 0.3s ease);
}

.error-icon {
  font-size: 3rem;
  filter: drop-shadow(0 4px 3px rgba(239, 68, 68, 0.1));
}

.plugin-error h3 {
  margin: 0;
  font-size: var(--font-size-xl, 1.25rem);
  font-weight: var(--font-weight-semibold, 600);
  color: var(--error-700, #c53030);
}

.plugin-error p {
  margin: 0;
  max-width: 600px;
  line-height: 1.6;
}

.retry-btn {
  background: var(--error-700, #c53030);
  color: white;
  border: none;
  padding: var(--space-2, 8px) var(--space-4, 16px);
  border-radius: var(--radius-md, 6px);
  cursor: pointer;
  transition: var(--transition-all, all 0.2s ease);
  margin-top: var(--space-2, 8px);
  font-weight: var(--font-weight-medium, 500);
  box-shadow: var(--shadow, 0 2px 4px rgba(0,0,0,0.1));
}

.retry-btn:hover {
  background: var(--error-800, #991b1b);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md, 0 4px 6px rgba(0,0,0,0.1));
}

/* 插件内容区域 */
.plugin-content { 
  width: 100%; 
  min-height: 300px;
  overflow: auto; 
  background: var(--bg-primary, #ffffff);
  border-radius: var(--radius-lg, 10px);
  border: 1px solid var(--border-light, #e5e7eb);
  transition: var(--transition-all, all 0.3s ease);
  display: flex;
  flex-direction: column;
}

/* 扁平化样式调整 */
.plugin-toolbar { 
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
  padding: 12px 16px; 
  border-bottom: 1px solid var(--border-color, #e5e7eb); 
  background: var(--toolbar-bg, #f8fafc); 
  position: sticky; 
  top: 0; 
  z-index: 1;
  transition: var(--transition-all, all 0.3s ease);
}

.plugin-meta .plugin-name { 
  color: var(--primary-800, #3949ab); 
  font-weight: 600; 
}

.plugin-meta .plugin-version { 
  margin-left: 8px; 
  color: #64748b; 
  font-size: 12px; 
  background: var(--badge-bg, #f1f5f9); 
  border-radius: 6px; 
  padding: 2px 6px; 
}

.plugin-meta .plugin-desc { 
  color: var(--muted, #6b7280); 
  margin-left: 6px; 
}

.plugin-sep { 
  color: var(--muted, #9ca3af); 
  margin: 0 6px; 
}

.plugin-actions .toolbar-btn { 
  background: var(--primary-600, #667eea); 
  color: #fff; 
  border: none; 
  padding: 6px 10px; 
  border-radius: 8px; 
  font-size: 13px; 
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm, 0 1px 3px rgba(0,0,0,0.1));
}

.toolbar-btn:hover { 
  background: var(--primary-700, #5a67d8); 
  transform: translateY(-1px);
  box-shadow: var(--shadow, 0 2px 4px rgba(0,0,0,0.1));
}

.toolbar-btn:active { 
  transform: translateY(0);
  box-shadow: var(--shadow-sm, 0 1px 3px rgba(0,0,0,0.1));
}

.plugin-body { 
  padding: 16px; 
  overflow: auto; 
  height: 100%; 
}

/* 空状态 */
.plugin-empty { 
  display: flex;
  flex-direction: column;
  align-items: center; 
  justify-content: center; 
  min-height: 300px;
  color: var(--text-muted, #6b7280); 
  gap: var(--space-3, 12px);
  padding: var(--space-6, 24px);
  background: var(--bg-secondary, #f8fafc);
  border-radius: var(--radius-lg, 10px);
  transition: var(--transition-all, all 0.3s ease);
}

.empty-icon {
  font-size: 3.5rem;
  opacity: 0.6;
}

.plugin-empty p {
  margin: 0;
  font-size: var(--font-size-lg, 1.125rem);
  font-weight: var(--font-weight-medium, 500);
}

/* 滚动条样式 */
.plugin-content::-webkit-scrollbar, 
.plugin-body::-webkit-scrollbar { 
  width: 8px; 
  height: 8px; 
}

.plugin-content::-webkit-scrollbar-thumb, 
.plugin-body::-webkit-scrollbar-thumb { 
  background: rgba(102,126,234,0.25); 
  border-radius: 8px; 
}

/* 淡入动画 */
.fade-in { 
  animation: fadeIn 0.25s ease-out both; 
}

@keyframes fadeIn { 
  from { 
    opacity: 0; 
    transform: translateY(4px); 
  } 
  to { 
    opacity: 1; 
    transform: translateY(0); 
  } 
}

/* 响应式设计 */
@media (max-width: 768px) {
  .plugin-loading,
  .plugin-error,
  .plugin-empty {
    min-height: 250px;
    padding: var(--space-4, 16px);
  }
  
  .spinner {
    width: 32px;
    height: 32px;
  }
  
  .error-icon,
  .empty-icon {
    font-size: 2.5rem;
  }
  
  .plugin-error h3 {
    font-size: var(--font-size-lg, 1.125rem);
  }
}

@media (max-width: 480px) {
  .plugin-renderer {
    padding: var(--space-3, 12px);
  }
  
  .plugin-content {
    border-radius: var(--radius-md, 8px);
  }
  
  .plugin-toolbar {
    padding: 10px 12px;
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .plugin-body {
    padding: 12px;
  }
}
</style>