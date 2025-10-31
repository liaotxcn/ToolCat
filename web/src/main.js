import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import axios from 'axios'

// 导入插件
import FormatConverterPlugin from './plugins/FormatConverterPlugin.js'
import pluginManager from './pluginManager.js'

const app = createApp(App)
app.config.globalProperties.$axios = axios

// 注册插件
pluginManager.registerPlugin('FormatConverterPlugin', new FormatConverterPlugin())

app.mount('#app')
