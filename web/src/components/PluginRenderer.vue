<template>
  <div class="plugin-renderer">
    <div v-if="loading" class="plugin-loading">
      加载插件中...
    </div>
    <div v-else-if="error" class="plugin-error">
      <p>{{ error }}</p>
    </div>
    <div v-else-if="pluginComponent" class="plugin-content">
      <!-- 使用动态组件渲染插件 -->
      <component :is="pluginComponent" :plugin="plugin" ref="pluginContainer"></component>
    </div>
    <div v-else class="plugin-empty">
      请选择一个插件
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, defineComponent } from 'vue'

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

// 辅助函数：更新组件数据
const updateComponentData = async (pluginInstance, component) => {
  // 如果有loadNotesFromAPI方法，调用它更新数据
  if (pluginInstance.loadNotesFromAPI && typeof pluginInstance.loadNotesFromAPI === 'function') {
    await pluginInstance.loadNotesFromAPI()
    // 更新组件数据
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
    if (renderResult.template) {
      // 创建方法映射，确保它们能访问插件实例
      const pluginMethods = {};
      if (renderResult.methods) {
        Object.keys(renderResult.methods).forEach(key => {
          if (typeof renderResult.methods[key] === 'function') {
            // 包装原始方法，确保它能访问插件实例
            pluginMethods[key] = function(...args) {
              const result = renderResult.methods[key].apply(this, args);
              // 如果是Promise，在完成后更新数据
              if (result && typeof result.then === 'function') {
                return result.then(async (resolvedResult) => {
                  // 异步方法完成后更新数据
                  await updateComponentData(plugin.value, this);
                  return resolvedResult;
                });
              } else {
                // 同步方法完成后更新数据
                updateComponentData(plugin.value, this);
                return result;
              }
            };
          }
        });
      }

      // 创建一个新的Vue组件
      pluginComponent.value = defineComponent({
        name: `${props.pluginName}-component`,
        props: {
          plugin: Object
        },
        template: renderResult.template,
        data() {
          // 调用插件的data函数，传入插件实例作为上下文
          if (renderResult.data && typeof renderResult.data === 'function') {
            const dataResult = renderResult.data.call(plugin.value);
            // 确保笔记数据正确初始化
            if (plugin.value.getAllNotes && !dataResult.notes) {
              dataResult.notes = [...(plugin.value.getAllNotes() || [])];
            }
            return dataResult;
          }
          return {};
        },
        methods: {
          // 添加辅助函数
          escapeHtml: (text) => {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
          },
          formatDate: (dateString) => {
            try {
              return new Date(dateString).toLocaleString();
            } catch (e) {
              return dateString;
            }
          },
          // 合并包装后的插件方法
          ...pluginMethods
        },
        computed: renderResult.computed || {},
        watch: renderResult.watch || {},
        created() {
          // 将插件实例的方法绑定到组件上
          const pluginInstance = this.plugin;
          // 将插件实例的关键方法直接添加到组件实例
          if (pluginInstance) {
            ['addNote', 'deleteNote', 'getAllNotes', 'loadNotesFromAPI'].forEach(methodName => {
              if (typeof pluginInstance[methodName] === 'function') {
                this[methodName] = pluginInstance[methodName].bind(pluginInstance);
              }
            });
          }
        },
        mounted() {
          // 组件挂载后加载初始数据
          if (this.loadNotesFromAPI && typeof this.loadNotesFromAPI === 'function') {
            this.loadNotesFromAPI().then(() => {
              if (this.getAllNotes && typeof this.getAllNotes === 'function') {
                this.notes = [...(this.getAllNotes() || [])];
              }
            });
          }
        }
      });
    }

    // 处理插件样式
    if (renderResult.css) {
      loadPluginCSS(renderResult.css)
    }

  } catch (err) {
    console.error('加载插件失败:', err)
    error.value = `加载插件失败: ${err.message || '未知错误'}`
  } finally {
    loading.value = false
  }
}





// 加载插件CSS
const loadPluginCSS = (css) => {
  if (!css) return

  // 创建style元素
  const styleId = `plugin-css-${props.pluginName}`
  let styleElement = document.getElementById(styleId)

  if (!styleElement) {
    styleElement = document.createElement('style')
    styleElement.id = styleId
    document.head.appendChild(styleElement)
  }

  styleElement.textContent = css
}

// 辅助函数：HTML转义
const escapeHtml = (text) => {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

// 辅助函数：格式化日期
const formatDate = (dateString) => {
  try {
    return new Date(dateString).toLocaleString()
  } catch (e) {
    return dateString
  }
}

// 监听pluginName变化
watch(
  () => props.pluginName,
  () => {
    loadPlugin()
  },
  {
    immediate: true
  }
)

// 导出方法供父组件使用
defineExpose({
  refreshPlugin: loadPlugin
})

// 组件挂载时加载插件
onMounted(() => {
  loadPlugin()
})
</script>

<style scoped>
.plugin-renderer {
  width: 100%;
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
}

.plugin-loading,
.plugin-error,
.plugin-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
}

.plugin-error {
  color: #f44336;
}

.plugin-content {
  width: 100%;
  height: 100%;
  overflow: auto;
}
</style>