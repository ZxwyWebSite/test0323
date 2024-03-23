## ZxwyWebSite/zTool
### 简介
+ Golang常用工具
+ 暂未正式发布，可在近期项目中找到beta版
+ 默认以MIT协议开源
+ 

### 分类
base/ 基础
   + 基础库，包含常用操作
   - Str_ 高性能字符串
   - Var_ 变量操作
   - Cmd_ 终端执行 (beta)
   - Cyp_ 数据加密 (beta)
   - Pak_ 文件打包 (beta)
   - Tme_ 定时相关 (beta)

x/ 扩展(ext)
   + 简单功能扩展库
   - cookie (来自Alist)
   - bytesconv (来自Gin)
   - json (来自Gin)
      * 支持Tag: "go_json", "json"(std), "jsoniter", "sonic"

mod/ 模块
   + 经过修改的外部模块
   - v1.16.0-color

cache/ 缓存
   + 缓存支持库，可使用以下驱动
   - memo 内存缓存
   - redis

conf/ 配置
   + 配置文件处理库，暂时只支持ini格式

json/ JSON
   + 第三方Json处理库，须使用Tag开启 (感谢Gin提供源码)
   - go_json
   - json (std)
   - jsoniter
   - sonic

logs/ 日志
   + 简单日志输出库，未完善 (始于lx-source项目)

menu/ 菜单
   + 创建简单的命令行交互程序 (始于ngconf项目)

storage/ 存储
   + 简单存储驱动，提供统一调用接口
   + 自开发
      - local 本机存储
      - cloudreve 平步云端
      - ...
   + Alist移植
      - (用到时也许就做了)

zcypt/ 加密
   + 封装加解密库，兼容Node标准 (始于lx-sync项目)

### 更新
2023-12-09
   + 整合文件，创建Package

### 其它
+ None
+ 