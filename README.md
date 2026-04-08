# 🚀 My Fullstack Blog (全栈个人博客系统)

基于 **Go (Gin) + Vue 3** 开发的现代化、响应式个人博客系统。
本项目从零开始构建，采用了前后端分离架构，并深度整合了多种企业级安全与高并发处理策略。

## ✨ 核心特性 (Features)

- **📝 沉浸式写作**：集成 `md-editor-v3`，支持 Markdown 实时预览、代码高亮与拖拽/粘贴本地图片上传。
- **🚀 高性能分页**：彻底摒弃传统 Offset 深度分页，采用**雪花算法 (Snowflake)** 生成分布式 ID，结合**游标分页 (Cursor Pagination)** 和瀑布流加载，应对海量数据。
- **🛡️ 坚固的安全防线**：
  - **IP 阶梯式限流防刷**：基于 Redis 实现 `429 Too Many Requests` (高频警告) 与 `403 Forbidden` (黑名单封禁) 机制。
  - **XSS 彻底免疫**：后端接入 `bluemonday` 严格过滤 UGC 内容中的恶意 HTML/Script 标签。
  - **文件上传校验**：通过读取文件头 (Magic Bytes) 真实校验图片类型，防止 Web Shell 木马上传。
- **💬 互动留言板**：前台支持访客留言，后台支持实时审核与恶意 IP 一键封禁。
- **🔍 智能检索**：支持对文章标题与摘要进行 SQL 级模糊搜索。

---

## 🛠️ 技术栈与核心依赖 (Tech Stack)

### 后端 (Backend - Go)
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) 
- **ORM**: [GORM](https://gorm.io/gorm) (搭配 MySQL 驱动)
- **缓存与限流**: [go-redis/v9](https://github.com/redis/go-redis/v9)
- **分布式 ID**: [snowflake](https://github.com/bwmarrin/snowflake)
- **安全认证**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) (JWT Token 鉴权)
- **XSS 净化**: [bluemonday](https://github.com/microcosm-cc/bluemonday)

### 前端 (Frontend - Vue 3)
- **核心框架**: Vue 3 (Composition API) + Vite
- **路由管理**: Vue Router 4
- **UI 组件库**: [Element Plus](https://element-plus.org/)
- **网络请求**: Axios (封装全局拦截器)
- **Markdown 引擎**: [md-editor-v3](https://github.com/imzbf/md-editor-v3) (支持暗黑模式与定制化解析)

---

## ⚙️ 环境准备与本地运行 (Quick Start)

在运行本项目之前，请确保你的本地环境已安装以下依赖：
- **Go** (推荐 1.20+)
- **Node.js** (推荐 18+)
- **MySQL** (推荐 8.0+)
- **Redis**: 
  - *开发环境测试版本为 `cygwin-8.6.2` (Windows 环境下的移植版)。*
  - *生产环境推荐使用官方标准版 Redis 7.x 或 8.x (建议通过 Docker 部署)。*

### 1. 克隆项目
```bash
git@github.com:114514Litiansuo1/my-fullstack-blog.git](https://github.com/114514Litiansuo1/my-fullstack-blog.git)
cd my-fullstack-blog
```

### 2. 后端配置与启动
进入 `backend` 目录，进行核心配置：

## ⚠️ 注意事项 1：MySQL 数据库配置
你需要提前在 MySQL 中创建一个数据库（例如 `go_blog` ），并在后端的数据库初始化代码中填入你自己的 MySQL 用户名和密码。

## ⚠️ 注意事项 2：超级管理员密码 (环境变量)
为了防止密码硬编码泄露，系统在首次启动创建超级管理员时，会读取名为 `password` 的环境变量。请在启动前配置它：
```bash
# Linux / macOS
export password="your_secure_password"

# Windows (CMD)
set password="your_secure_password"

# Windows (PowerShell)
$env:password="your_secure_password"
```

## ⚠️ 注意事项 3：雪花算法节点配置
如果你打算在多台服务器上分布式部署此博客，请务必在代码中修改 Snowflake 的初始节点 ID (Node ID)，以防止生成重复的文章/留言 ID。

启动后端服务：
```bash
go mod tidy
go run main.go
```

后端默认运行在 `http://localhost:8080`

### 3. 前端配置与启动
打开一个新的终端，进入 `frontend` 目录：
```bash
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173` 。访问 /admin 进入后台管理系统。

### 📜 开源协议 (License)
本项目基于 MIT License 开源，欢迎自由使用、修改和分发，但请保留原作者版权声明。
