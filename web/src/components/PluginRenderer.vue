<template>
  <div class="plugin-renderer">
    <div v-if="plugin" class="plugin-container">
      <!-- 使用v-html渲染插件的模板内容 -->
      <div v-html="pluginTemplate"></div>
    </div>
    <div v-else class="plugin-error">
      插件未加载
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'

export default {
  name: 'PluginRenderer',
  props: {
    pluginName: {
      type: String,
      required: true
    },
    pluginManager: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const plugin = ref(null)
    const pluginTemplate = ref('')
    const pluginCSS = ref('')
    const pluginData = ref({})
    const pluginMethods = ref({})

    // 加载插件
    const loadPlugin = () => {
      const loadedPlugin = props.pluginManager.getPlugin(props.pluginName)
      if (loadedPlugin) {
        plugin.value = loadedPlugin
        const renderResult = loadedPlugin.render()
        
        if (renderResult) {
          // 处理插件模板
          if (renderResult.template) {
            pluginTemplate.value = renderResult.template
          }
          
          // 处理插件样式
          if (renderResult.css) {
            loadPluginCSS(renderResult.css)
          }
          
          // 处理插件数据
          if (renderResult.data && typeof renderResult.data === 'function') {
            try {
              pluginData.value = renderResult.data.call(loadedPlugin)
            } catch (error) {
              console.error('加载插件数据失败:', error)
            }
          }
          
          // 处理插件方法
          if (renderResult.methods) {
            pluginMethods.value = renderResult.methods
            // 绑定方法到插件实例
            Object.keys(pluginMethods.value).forEach(methodName => {
              const originalMethod = pluginMethods.value[methodName]
              if (typeof originalMethod === 'function') {
                pluginMethods.value[methodName] = function(...args) {
                  return originalMethod.apply(loadedPlugin, args)
                }
              }
            })
          }
          
          // 处理插件计算属性
          if (renderResult.computed) {
            // 这里可以实现计算属性的处理
          }
          
          // 处理插件监听器
          if (renderResult.watch) {
            // 这里可以实现监听器的处理
          }
        }
      } else {
        console.error(`插件 ${props.pluginName} 不存在`)
        plugin.value = null
        pluginTemplate.value = ''
      }
    }

    // 加载插件CSS
    const loadPluginCSS = (css) => {
      // 清除之前的样式
      const existingStyle = document.getElementById(`plugin-${props.pluginName}-style`)
      if (existingStyle) {
        existingStyle.remove()
      }
      
      // 添加新样式
      const style = document.createElement('style')
      style.id = `plugin-${props.pluginName}-style`
      style.textContent = css
      document.head.appendChild(style)
    }

    // 当插件名称变化时重新加载插件
    watch(() => props.pluginName, () => {
      loadPlugin()
    })

    // 当插件管理器变化时重新加载插件
    watch(() => props.pluginManager, () => {
      loadPlugin()
    })

    // 组件挂载时加载插件
    onMounted(() => {
      loadPlugin()
    })

    return {
      plugin,
      pluginTemplate,
      pluginCSS,
      pluginData,
      pluginMethods
    }
  }
}
</script>

<style scoped>
.plugin-renderer {
  width: 100%;
  height: 100%;
}

.plugin-container {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.plugin-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #ff4d4f;
  font-size: 14px;
}
</style>