import axios from 'axios'

// 创建一个基本的axios实例
const api = axios.create()

// 本地转换功能实现作为后备方案
const localConverters = {
  // 本地JSON转YAML实现
  jsonToYaml: (jsonData) => {
    let yaml = ''
    const indentLevel = 0
    
    // 简单的JSON转YAML实现
    function jsonToYamlRecursive(obj, indent = 0) {
      const spaces = '  '.repeat(indent)
      
      if (obj === null) return spaces + 'null\n'
      if (typeof obj === 'boolean' || typeof obj === 'number') return spaces + obj + '\n'
      if (typeof obj === 'string') {
        // 如果字符串包含特殊字符或换行符，使用引号
        if (obj.match(/[\\n\r:\[\]{}#&*!|>\'"%@`]/) || obj.trim() !== obj) {
          return spaces + '"' + obj.replace(/"/g, '\\"').replace(/\\n/g, '\\n') + '"\n'
        }
        return spaces + obj + '\n'
      }
      
      if (Array.isArray(obj)) {
        if (obj.length === 0) return spaces + '[]\n'
        
        let result = ''
        for (const item of obj) {
          result += spaces + '- ' + jsonToYamlRecursive(item, indent + 1).slice(spaces.length + 2)
        }
        return result
      }
      
      // 对象处理
      const keys = Object.keys(obj)
      if (keys.length === 0) return spaces + '{}' + '\n'
      
      let result = ''
      for (const key of keys) {
        result += spaces + key + ': ' + jsonToYamlRecursive(obj[key], indent + 1).slice(spaces.length + key.length + 2)
      }
      return result
    }
    
    yaml = jsonToYamlRecursive(jsonData)
    return yaml
  },
  
  // 本地YAML转JSON实现
  yamlToJson: (yamlStr) => {
    // 简单的YAML解析器实现
    const lines = yamlStr.split('\n')
    const result = {}
    const stack = [result]
    const indentStack = [-2] // 初始缩进为-2表示根级别
    
    lines.forEach(line => {
      line = line.trim()
      if (!line || line.startsWith('#')) return
      
      // 处理数组项
      if (line.startsWith('- ')) {
        const content = line.substring(2).trim()
        const lastObj = stack[stack.length - 1]
        const lastIndent = indentStack[indentStack.length - 1]
        
        if (content.includes(': ')) {
          // 复杂数组项 {key: value}
          const [key, ...valueParts] = content.split(': ')
          const value = valueParts.join(': ')
          const newObj = {[key]: parseYamlValue(value)}
          
          if (Array.isArray(lastObj)) {
            lastObj.push(newObj)
            stack.push(newObj)
            indentStack.push(lastIndent + 2)
          } else {
            // 这应该是一个新的数组
            const arr = [newObj]
            const parent = stack[stack.length - 2]
            const parentKeys = Object.keys(parent)
            const lastKey = parentKeys[parentKeys.length - 1]
            parent[lastKey] = arr
            stack.push(newObj)
            indentStack.push(lastIndent + 2)
          }
        } else {
          // 简单数组项
          const value = parseYamlValue(content)
          if (Array.isArray(lastObj)) {
            lastObj.push(value)
          } else {
            // 这应该是一个新的数组
            const arr = [value]
            const parent = stack[stack.length - 2]
            const parentKeys = Object.keys(parent)
            const lastKey = parentKeys[parentKeys.length - 1]
            parent[lastKey] = arr
          }
        }
      } else if (line.includes(': ')) {
        // 处理键值对
        const indent = (line.match(/^ */) || [])[0].length
        const [key, ...valueParts] = line.split(': ')
        const value = valueParts.join(': ')
        
        // 找到对应的父对象
        while (indent <= indentStack[indentStack.length - 1]) {
          stack.pop()
          indentStack.pop()
        }
        
        const parent = stack[stack.length - 1]
        
        // 检查是否是嵌套对象的开始
        if (value === '') {
          const newObj = {}
          parent[key] = newObj
          stack.push(newObj)
          indentStack.push(indent)
        } else {
          // 普通键值对
          parent[key] = parseYamlValue(value)
        }
      }
    })
    
    // 解析YAML值
    function parseYamlValue(value) {
      if (value === 'true') return true
      if (value === 'false') return false
      if (value === 'null') return null
      if (!isNaN(value) && value !== '') return Number(value)
      // 处理引号字符串
      if ((value.startsWith('"') && value.endsWith('"')) || 
          (value.startsWith('\'') && value.endsWith('\''))) {
        return value.substring(1, value.length - 1)
      }
      return value
    }
    
    return result
  }
}

// 格式转换器相关API方法
export const formatConverterService = {
  // JSON转YAML
  jsonToYaml: async (jsonData) => {
    try {
      // 尝试通过API转换
      const response = await api.post('/plugins/format_converter/convert/json-to-yaml', jsonData, {
        headers: {
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      // 检查响应是否有效
      return response?.data || localConverters.jsonToYaml(jsonData)
    } catch (error) {
      console.warn('使用本地JSON转YAML作为后备方案:', error.message)
      // 使用本地转换作为后备
      return localConverters.jsonToYaml(jsonData)
    }
  },

  // YAML转JSON
  yamlToJson: async (yamlData) => {
    try {
      // 尝试通过API转换
      const response = await api.post('/plugins/format_converter/convert/yaml-to-json', yamlData, {
        headers: {
          'Content-Type': 'text/yaml; charset=utf-8'
        }
      })
      // 检查响应是否有效
      return response?.data || localConverters.yamlToJson(yamlData)
    } catch (error) {
      console.warn('使用本地YAML转JSON作为后备方案:', error.message)
      // 使用本地转换作为后备
      return localConverters.yamlToJson(yamlData)
    }
  },

  // JSON转Protobuf
  jsonToProtobuf: async (jsonData) => {
    try {
      // 尝试通过API转换
      const response = await api.post('/plugins/format_converter/convert/json-to-protobuf', jsonData, {
        headers: {
          'Content-Type': 'application/json; charset=utf-8'
        },
        responseType: 'arraybuffer' // 接收二进制响应
      })
      return response?.data || new ArrayBuffer(0)
    } catch (error) {
      console.warn('JSON转Protobuf失败，API可能不可用:', error.message)
      // 返回空ArrayBuffer作为后备
      return new ArrayBuffer(0)
    }
  },

  // Protobuf转JSON
  protobufToJson: async (protobufData) => {
    try {
      // 尝试通过API转换
      const response = await api.post('/plugins/format_converter/convert/protobuf-to-json', protobufData, {
        headers: {
          'Content-Type': 'application/x-protobuf'
        }
      })
      return response?.data || {}
    } catch (error) {
      console.warn('Protobuf转JSON失败，API可能不可用:', error.message)
      // 返回空对象作为后备
      return {}
    }
  },

  // 获取插件信息
  getPluginInfo: async () => {
    try {
      const response = await api.get('/plugins/format_converter/')
      return response?.data || { name: '格式转换器', version: '1.0.0', status: '本地模式' }
    } catch (error) {
      console.warn('获取插件信息失败，使用默认信息:', error.message)
      // 返回默认信息
      return { name: '格式转换器', version: '1.0.0', status: '本地模式' }
    }
  }
}