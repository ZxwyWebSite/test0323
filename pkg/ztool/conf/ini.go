package conf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZxwyWebSite/ztool"
	"github.com/ZxwyWebSite/ztool/logs"
	"github.com/ZxwyWebSite/ztool/x/bytesconv"
	"github.com/go-ini/ini"
)

var (
	defCfg = struct{}{}
	Config = defCfg
	ipath  string
)

// 保存已修改配置
func Save() error {
	if ipath == `` {
		return errors.New(`请先载入配置`)
	}
	buf := ini.Empty()
	if err := buf.ReflectFrom(Config); err != nil {
		return err
	}
	return buf.SaveTo(ipath)
}

// 对象化部分
type (
	Confs struct {
		model  any       // 配置模板
		Cfg    *ini.File // 自定义操作接口
		config *Confg    // 配置
		ipath  string    // 载入时路径
	}
	Confg struct {
		// Path       string
		AutoFormat bool         // 自动格式化保存 (即使不是初始化运行也会检测文件是否需要更新)
		UseBuf     bool         // 保存时隔离运行 (有忽略未定义值的作用)
		UnPretty   bool         // 不使用自动对齐 (适用添加了大量注释的情况)
		Loger      *logs.Logger // 开启日志功能
		// MustInit   bool         // 初始化后退出程序
		// First func() // 第一次初始化回调
	}
)

var (
	ErrNotInit = errors.New(`请先载入配置`)
)

// 创建Confs对象 (模板, 配置)
func New(model any, config *Confg) (*Confs, error) {
	c := ini.Empty()
	if err := c.ReflectFrom(model); err != nil {
		return nil, fmt.Errorf(`无法载入默认配置: %s`, err)
	}
	return &Confs{
		model:  model,
		Cfg:    c,
		config: config,
	}, nil
}

// 载入文件
func (c *Confs) Init(path string) error {
	var fromEmpty bool
	if path == `` {
		return ErrNotInit
	}
	if ztool.Fbj_IsExists(path) {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf(`无法读取配置文件: %s`, err)
		}
		if err := c.Cfg.Append(data); err != nil {
			return fmt.Errorf(`无法解析配置文件: %s`, err)
		}
		c.ipath = path
	} else {
		fromEmpty = true
		// if dir := filepath.Dir(path); !ztool.Fbj_IsExists(dir) {
		// 	if err := os.MkdirAll(dir, 0644); err != nil {
		// 		return fmt.Errorf(`无法创建配置目录: %s`, err)
		// 	}
		// }
	}
	// 映射到结构体
	if err := c.Cfg.StrictMapTo(c.model); err != nil {
		return fmt.Errorf(`无法映射配置文件: %s`, err)
	}
	// 写入格式化后的配置
	if c.config.AutoFormat || fromEmpty {
		// if err := c.Cfg.SaveTo(path); err != nil {
		// 	return fmt.Errorf(`无法写入配置文件, %s`, err)
		// }
		if c.config.Loger != nil && fromEmpty {
			c.config.Loger.Info(`已初始化配置文件并以默认配置运行`)
			c.config.Loger.Info(`文件路径：%q，修改后重启程序生效`, filepath.Clean(path))
			c.config.Loger.Warn(`注：误删配置字段不要慌，先保存已修改的部分，下次启动会自动补齐 :)`)
			// if c.config.MustInit {
			// 	defer c.config.Loger.Fatal(`已开启MustInit，请修改配置文件后再次启动`)
			// }
		}
		// 检测是否需要写入文件 (防止重复写入消耗磁盘寿命)
		return c.Save(path)
	}
	return nil
}
func (c *Confs) MustInit(path string) {
	if err := c.Init(path); err != nil {
		if c.config.Loger != nil {
			c.config.Loger.Fatal(err.Error())
		} else {
			panic(err)
		}
	}
}

// 保存至
func (c *Confs) Save(to string) error {
	if c.config.UnPretty {
		o1, o2 := ini.PrettyFormat, ini.PrettyEqual                   // 保存当前值
		ini.PrettyFormat, ini.PrettyEqual = false, true               // 设置UnPretty
		defer func() { ini.PrettyFormat, ini.PrettyEqual = o1, o2 }() // 恢复当前值
	}
	if to == `` {
		if c.ipath == `` {
			return ErrNotInit
		}
		to = c.ipath
	}
	var buf *ini.File
	if c.config.UseBuf {
		buf = ini.Empty()
		if err := buf.ReflectFrom(c.model); err != nil {
			return err
		}
	} else {
		buf = c.Cfg
	}
	// 检测变化
	obuf := new(bytes.Buffer)
	_, err := buf.WriteTo(obuf)
	if err == nil {
		obyte := obuf.Bytes()
		var f []byte
		f, err = os.ReadFile(to)
		if err == nil {
			if bytesconv.BytesToString(f) == bytesconv.BytesToString(obyte) {
				// if c.config.Loger != nil {
				// 	c.config.Loger.Info(`配置文件无变化，跳过写入`)
				// }
				return nil
			}
		}
		err = ztool.Fbj_SaveFile(to, obyte)
	}
	// 写入文件
	if /*err := buf.SaveTo(to);*/ err != nil {
		return fmt.Errorf(`无法写入配置文件: %s`, err)
	}
	if c.config.Loger != nil {
		c.config.Loger.Debug(`配置文件格式化完成`)
	}
	c.ipath = to
	return nil
}
