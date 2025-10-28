package formatconverter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"toolcat/pkg"
	"toolcat/plugins/core"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

var _ core.Plugin = &FormatConverterPlugin{}

type FormatConverterPlugin struct {
	mu            sync.RWMutex
	pluginManager *core.PluginManager
	version       string
}

func (p *FormatConverterPlugin) Name() string { return "format_converter" }
func (p *FormatConverterPlugin) Description() string {
	return "格式转换插件：JSON↔YAML, JSON↔Protobuf"
}
func (p *FormatConverterPlugin) Version() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.version
}

func (p *FormatConverterPlugin) Init() error {
	p.mu.Lock()
	p.version = "1.0.0"
	p.mu.Unlock()
	pkg.Info(fmt.Sprintf("%s: 初始化完成", p.Name()))
	return nil
}
func (p *FormatConverterPlugin) Shutdown() error                          { return nil }
func (p *FormatConverterPlugin) OnEnable() error                          { return nil }
func (p *FormatConverterPlugin) OnDisable() error                         { return nil }
func (p *FormatConverterPlugin) GetDependencies() []string                { return []string{} }
func (p *FormatConverterPlugin) GetConflicts() []string                   { return []string{} }
func (p *FormatConverterPlugin) SetPluginManager(m *core.PluginManager)   { p.pluginManager = m }
func (p *FormatConverterPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return []gin.HandlerFunc{} }

func (p *FormatConverterPlugin) GetRoutes() []core.Route {
	return []core.Route{
		{Path: "/", Method: "GET", Handler: func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin": p.Name(), "description": p.Description(), "version": p.Version(),
				"endpoints": []string{
					"POST /convert/json-to-yaml",
					"POST /convert/yaml-to-json",
					"POST /convert/json-to-protobuf",
					"POST /convert/protobuf-to-json",
				},
			})
		}, Description: "获取插件信息", AuthRequired: false, Tags: []string{"info"}},
		{Path: "/convert/json-to-yaml", Method: "POST", Handler: func(c *gin.Context) {
			data, err := c.GetRawData()
			if err != nil {
				c.JSON(400, gin.H{"error": fmt.Sprintf("读取请求体失败: %v", err)})
				return
			}
			var obj interface{}
			if err = json.Unmarshal(data, &obj); err != nil {
				c.JSON(400, gin.H{"error": fmt.Sprintf("解析JSON失败: %v", err)})
				return
			}
			out, err := yaml.Marshal(obj)
			if err != nil {
				c.JSON(500, gin.H{"error": fmt.Sprintf("转换为YAML失败: %v", err)})
				return
			}
			c.Data(200, "text/yaml; charset=utf-8", out)
		}, Description: "将JSON转换为YAML（请求体为原始JSON）", AuthRequired: false, Tags: []string{"convert"}},
		{Path: "/convert/yaml-to-json", Method: "POST", Handler: func(c *gin.Context) {
			data, err := c.GetRawData()
			if err != nil {
				c.JSON(400, gin.H{"error": fmt.Sprintf("读取请求体失败: %v", err)})
				return
			}
			var obj interface{}
			if err = yaml.Unmarshal(data, &obj); err != nil {
				c.JSON(400, gin.H{"error": fmt.Sprintf("解析YAML失败: %v", err)})
				return
			}
			norm := normalizeYaml(obj)
			out, err := json.Marshal(norm)
			if err != nil {
				c.JSON(500, gin.H{"error": fmt.Sprintf("转换为JSON失败: %v", err)})
				return
			}
			c.Data(200, "application/json; charset=utf-8", out)
		}, Description: "将YAML转换为JSON（请求体为原始YAML）", AuthRequired: false, Tags: []string{"convert"}},
		{Path: "/convert/json-to-protobuf", Method: "POST", Handler: p.jsonToProtobufHandler,
			Description: "将JSON转换为Protobuf（使用DynamicMessage）", AuthRequired: false, Tags: []string{"convert"}},
		{Path: "/convert/protobuf-to-json", Method: "POST", Handler: p.protobufToJsonHandler,
			Description: "将Protobuf转换为JSON（使用DynamicMessage）", AuthRequired: false, Tags: []string{"convert"}},
	}
}

func (p *FormatConverterPlugin) RegisterRoutes(router *gin.Engine) {
	pkg.Info(fmt.Sprintf("%s: 使用GetRoutes进行路由注册", p.Name()))
}

func normalizeYaml(i interface{}) interface{} {
	switch v := i.(type) {
	case map[string]interface{}:
		for k, val := range v {
			v[k] = normalizeYaml(val)
		}
		return v
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(v))
		for k, val := range v {
			m[fmt.Sprint(k)] = normalizeYaml(val)
		}
		return m
	case []interface{}:
		for i := range v {
			v[i] = normalizeYaml(v[i])
		}
		return v
	default:
		return v
	}
}

// JSON到Protobuf的处理函数
func (p *FormatConverterPlugin) jsonToProtobufHandler(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("读取请求体失败: %v", err)})
		return
	}

	// 将JSON转换为Structpb.Struct
	var obj interface{}
	if err = json.Unmarshal(data, &obj); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("解析JSON失败: %v", err)})
		return
	}

	// 使用structpb将interface{}转换为protobuf兼容的结构
	structObj, err := structpb.NewValue(obj)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("转换为Protobuf结构失败: %v", err)})
		return
	}

	// 转换为二进制格式
	binaryData, err := proto.Marshal(structObj)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Protobuf序列化失败: %v", err)})
		return
	}

	c.Data(200, "application/x-protobuf", binaryData)
}

// Protobuf到JSON的处理函数
func (p *FormatConverterPlugin) protobufToJsonHandler(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("读取请求体失败: %v", err)})
		return
	}

	// 创建一个新的Structpb.Value作为接收容器
	value := &structpb.Value{}
	if err = proto.Unmarshal(data, value); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("解析Protobuf失败: %v", err)})
		return
	}

	// 将Protobuf转换为JSON
	jsonData, err := protojson.Marshal(value)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("转换为JSON失败: %v", err)})
		return
	}

	c.Data(200, "application/json; charset=utf-8", jsonData)
}

func (p *FormatConverterPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	action, _ := params["action"].(string)
	input, _ := params["input"].(string)
	switch action {
	case "json_to_yaml":
		if input == "" {
			return nil, errors.New("缺少输入: input")
		}
		var obj interface{}
		if err := json.Unmarshal([]byte(input), &obj); err != nil {
			return nil, fmt.Errorf("解析JSON失败: %w", err)
		}
		out, err := yaml.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("转换为YAML失败: %w", err)
		}
		return string(out), nil
	case "yaml_to_json":
		if input == "" {
			return nil, errors.New("缺少输入: input")
		}
		var obj interface{}
		if err := yaml.Unmarshal([]byte(input), &obj); err != nil {
			return nil, fmt.Errorf("解析YAML失败: %w", err)
		}
		norm := normalizeYaml(obj)
		out, err := json.Marshal(norm)
		if err != nil {
			return nil, fmt.Errorf("转换为JSON失败: %w", err)
		}
		return string(out), nil
	case "json_to_protobuf":
		if input == "" {
			return nil, errors.New("缺少输入: input")
		}
		var obj interface{}
		if err := json.Unmarshal([]byte(input), &obj); err != nil {
			return nil, fmt.Errorf("解析JSON失败: %w", err)
		}
		structObj, err := structpb.NewValue(obj)
		if err != nil {
			return nil, fmt.Errorf("转换为Protobuf结构失败: %w", err)
		}
		binaryData, err := proto.Marshal(structObj)
		if err != nil {
			return nil, fmt.Errorf("protobuf序列化失败: %w", err)
		}
		return binaryData, nil
	case "protobuf_to_json":
		if input == "" {
			return nil, errors.New("缺少输入: input")
		}
		value := &structpb.Value{}
		if err := proto.Unmarshal([]byte(input), value); err != nil {
			return nil, fmt.Errorf("解析Protobuf失败: %w", err)
		}
		jsonData, err := protojson.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("转换为JSON失败: %w", err)
		}
		return string(jsonData), nil
	default:
		return map[string]interface{}{
			"plugin":            p.Name(),
			"version":           p.Version(),
			"supported_actions": []string{"json_to_yaml", "yaml_to_json", "json_to_protobuf", "protobuf_to_json"},
		}, nil
	}
}
