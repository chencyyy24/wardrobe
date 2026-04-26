**个人衣柜 - 开发文档**

-----
**1. 项目概述**

**项目名称**：个人衣柜 (Personal Wardrobe)

**项目定位**：一款帮助用户数字化管理衣橱、快速创建穿搭方案并生成搭配卡片的双端应用（Web + iOS）。核心解决两个痛点：**衣物资产可视化** 与 **高效搭配复用**。

**主要功能**：

- 上传衣物照片，自动完成抠图，生成白底单品图
- 按大类/小类分级管理衣物
- 自由选择衣物组合搭配，实时预览并生成穿搭卡片
- 提供历史搭配记录，支持基于已有搭配快速修改

**技术目标**：功能实用、流程闭环、双端可用、开发量可控，适合作为结课作业展示。

-----
**2. 功能需求详述**

**2.1 衣物管理模块**

**2.1.1 上传衣物**

- 用户选择本地图片（原图），点击上传后直接返回列表页
- 图片上传至服务器后，后端**异步**执行抠图处理，生成白底图
- 前端展示原图，待抠图完成后自动替换为白底图；支持“处理中”状态展示
- 衣物信息包含：
  - 名称（自定义）
  - 图片（原图 + 白底图）
  - 分类：**大类 → 小类**，数据固定如下：
    - 外套：风衣、牛仔夹克、西装、针织开衫
    - 上衣：T恤、衬衫、卫衣、毛衣
    - 裤子：牛仔裤、休闲裤、短裤
    - 裙子：JK裙、百褶裙、连衣裙
    - 鞋子：运动鞋、帆布鞋、靴子
    - 配饰：帽子、项链、耳环、包包
- 上传后即生成一条衣物记录，状态为“处理中”，后台抠图完成后更新状态为“已完成”

**2.1.2 衣物列表**

- 支持按**大类**筛选（外套/上衣/裤子/裙子/鞋子/配饰）
- 卡片式展示衣物缩略图（优先使用白底图，未完成则显示原图）
- 长按或滑动可删除衣物

**2.2 搭配模块**

**2.2.1 新建搭配**

- 搭配规则：**外套（可选）+ 上衣（必选）+ 裤子（必选）/ 裙子（可选，与裤子互斥）+ 鞋子（可选）+ 配饰（可选）**
- 每个部位区域展示对应大类下已上传的全部衣物（白底图），横向滑动选择
- 实时预览区：选中衣物按穿搭层次摆放，直观展示组合效果
- 支持为搭配命名（或自动生成如“春日通勤风”）
- 点击“生成卡片”后，将预览区域截图或使用 Canvas 合成一张固定版式的穿搭卡片图片，上传至服务器
- 保存搭配数据，卡片图作为搭配的封面展示

**2.2.2 历史搭配**

- 搭配主页面包含两个 Tab：**“历史记录”** 与 **“新建搭配”**
- 历史记录以卡片流展示所有已保存搭配（卡片缩略图 + 名称 + 创建日期）
- 支持点击查看大图、删除搭配
- 新建搭配页提供“从历史复制”快捷入口，选择某套历史搭配后，自动填充衣物部位，用户可微调后保存为新搭配，避免重复操作

**2.3 其他**

- 基础设置/关于页面（可选）
-----
**3. 技术架构**

**3.1 整体架构图**

text

┌───────────┐       HTTP        ┌──────────────────┐      异步调用       ┌─────────────────┐

│  Vue Web  │ ◄───────────────► │   Go 后端 (8080)  │ ────────────────► │ Python 抠图微服务 │

│  (SPA)    │                   │  + 静态文件服务    │ ◄──────────────── │ (rembg HTTP API) │

│           │                   │  + 数据库/文件存储 │                   │                 │

│ iOS 壳    │                   └────────┬─────────┘                   └─────────────────┘

│(WKWebView)│                            │

└───────────┘                            ▼

`                                  `┌───────────────┐

`                                  `│  文件系统      │

`                                  `│  uploads/     │

`                                  `│  origin/      │

`                                  `│  masked/      │

`                                  `│  cards/       │

`                                  `└───────────────┘

**3.2 技术选型**

|层次|技术栈|说明|
| :- | :- | :- |
|**Web 前端**|Vue 3 + Vant UI / 其他移动端组件库|单页面应用，双端统一界面|
|**iOS 端**|Swift + WKWebView|本地加载 Web 资源或在线地址，最低成本实现双端|
|**后端**|Go (Gin/标准库)|提供 RESTful API，处理上传、搭配、状态查询|
|**抠图服务**|Python + rembg|独立微服务，接收图片返回白底抠图|
|**数据存储**|JSON 文件 / SQLite|开发阶段使用 JSON 文件，部署后可切换|
|**前端抠图后备**|@imgly/background-removal|已弃用此方案，改为后端静默处理|

**3.3 数据存储说明**

- **开发阶段**：数据存储于本地 JSON 文件（data.json），简单易调试
- **部署阶段**：可继续使用 JSON 或引入 SQLite，利用 Go 的 database/sql 快速切换
- 图片文件统一存放于 uploads/ 目录，按类型分 origin/、masked/、cards/
-----
**4. 数据模型**

**4.1 衣物 Clothing**

json

{

`  `"id": "uuid",

`  `"name": "优衣库条纹衬衫",

`  `"category": "上衣",

`  `"subcategory": "衬衫",

`  `"original\_image": "/uploads/origin/abc.jpg",

`  `"masked\_image": "/uploads/masked/abc.png",   *// 为空表示未处理*

`  `"status": "pending",  *// pending | done | failed*

`  `"created\_at": "2025-04-26T12:00:00Z"

}

**4.2 搭配 Outfit**

json

{

`  `"id": "uuid",

`  `"name": "通勤简约风",

`  `"items": {

`    `"outer": "clothing\_id or null",

`    `"top": "clothing\_id",

`    `"bottom": "clothing\_id",

`    `"skirt": "clothing\_id or null",

`    `"shoes": "clothing\_id or null",

`    `"accessory": "clothing\_id or null"

`  `},

`  `"card\_image": "/uploads/cards/xyz.png",

`  `"created\_at": "2025-04-26T13:00:00Z"

}

**说明**：items 中使用固定字段代表穿搭部位，比数组更直观，前端可据此精确渲染位置。

-----
**5. API 设计**

**5.1 衣物相关**

|方法|路径|功能|请求/参数|
| :- | :- | :- | :- |
|POST|/api/clothing|上传衣物|multipart/form-data：image(原图文件)、name、category、subcategory|
|GET|/api/clothing|获取衣物列表|?category=上衣 筛选|
|DELETE|/api/clothing/:id|删除衣物|同时删除关联文件|
|GET|/api/clothing/status|批量查询状态|?ids=id1,id2，返回各衣物 status 和 masked\_image 最新值|
|POST|/api/clothing/:id/retry|重新触发抠图|用于 failed 状态的衣物|

**5.2 搭配相关**

|方法|路径|功能|请求体|
| :- | :- | :- | :- |
|POST|/api/outfit|创建搭配|{ name, items: {outer, top, bottom, skirt, shoes, accessory}, card\_image: File } (multipart)|
|GET|/api/outfit|获取所有搭配|返回列表，含卡片缩略图 URL|
|DELETE|/api/outfit/:id|删除搭配|同时删除卡片文件|
|GET|/api/outfit/:id|获取搭配详情|返回完整数据，用于“从历史复制”|

**5.3 上传卡片专用**

| POST | /api/upload/card | 上传生成好的搭配卡片图 | 返回图片 URL |

所有静态图片资源通过 /uploads/ 路径直接访问（Nginx/Go 静态文件服务）。

-----
**6. 前端设计**

**6.1 页面路由**

text

/              衣柜（衣物列表，分类筛选）

/upload        新增衣物

/outfit        搭配主页（子 Tab：历史记录 | 新建搭配）

/outfit/new    搭配编辑器（也可作为搭配主页的“新建”Tab内容）

/my            我的（设置/关于）

**6.2 核心组件**

**6.2.1 文件上传组件 UploadView**

- 使用 <input type="file"> 选择图片，支持拍照或相册选取
- 表单：名称、大类下拉、小类联动
- 选择后显示原图预览，点击“上传”将数据以 FormData POST 至 /api/clothing
- 上传后返回列表，卡片展示原图 + “处理中”动画（如轻微闪烁）
- 隔 5 秒轮询 /api/clothing/status 获取最新状态，白底图生成后替换显示

**6.2.2 搭配编辑器 OutfitEditor**

- 部位选择区：外套、上衣、裤子、裙子、鞋子、配饰
- 每个区域水平滚动列表，展示对应大类下所有已上传衣物的白底图，点击选中（高亮边框）
- 实时预览区：使用 CSS 绝对定位按层次排列选中的白底图（外套在上，上衣在下，裤子/裙子更下，鞋子在底部，配饰可浮动），模拟真实穿搭效果
- 必选校验：上衣、裤子必须选中，否则“保存”按钮置灰
- “生成卡片”按钮：调用 html2canvas 截取预览区内容，转换为 Blob，上传至 /api/upload/card 获取 URL，再提交搭配数据
- 搭配名称可自动生成（如当前日期 + 风格词）或手动输入

**6.2.3 历史搭配页面 HistoryOutfits**

- 使用卡片网格展示所有搭配，每张卡片显示搭配卡片缩略图、名称、日期
- 点击进入详情（大图 + 单品列表）
- 每张卡片提供“复制此搭配”按钮，跳转至搭配编辑器并预填选中项

**6.2.4 公共工具函数**

- 状态轮询 Hook：useClothingStatus 管理待处理衣物轮询
- 图片懒加载与错误处理

**6.3 UI/UX 要点**

- 整体风格采用移动端优先设计，参考衣橱类 App 的简洁卡片风
- 搭配预览区采用模拟真实穿着的纵向排列，背景为浅色，衣物抠图居中排列
- 搭配卡片生成时提供 2~3 种模板（如竖排经典、田字格、杂志风格），通过 Canvas 动态绘制
-----
**7. 抠图异步处理方案**

**7.1 流程说明**

1. 前端上传原图 → Go 后端保存至 uploads/origin/，衣物记录状态 pending，立即返回成功
1. Go 后端通过 goroutine 异步调用 Python 抠图微服务 (HTTP POST /api/remove)，传入原图路径或文件
1. Python 微服务使用 rembg 库处理，返回白底 PNG 数据
1. Go 后端将结果保存至 uploads/masked/，更新衣物记录的 masked\_image 字段，状态改为 done
1. 若失败，状态改为 failed，前端可手动重试

**7.2 Python 抠图微服务搭建**

- 依赖：rembg + flask (或直接使用 rembg s 启动内置 HTTP 服务)
- 推荐使用 rembg s --host 0.0.0.0 --port 5000 直接启动服务，无需额外编写 Flask 代码
- 接口调用示例 (rembg 命令启动后的默认端点):

text

POST http://127.0.0.1:5000/api/remove

Content-Type: multipart/form-data

字段: input\_image (图片文件)

返回: image/png

- 部署时在服务器上使用 nohup 或 systemd 保持服务运行

**7.3 Go 异步任务实现**

go

func processImageAsync(clothingID, originPath string) {

`    `go func() {

`        `*// 调用 rembg 服务*

`        `resp, err := http.Post("http://127.0.0.1:5000/api/remove", ...)

`        `if err != nil {

`            `updateClothingStatus(clothingID, "failed")

`            `return

`        `}

`        `*// 保存 masked 图片*

`        `saveMaskedImage(clothingID, resp.Body)

`        `updateClothingStatus(clothingID, "done")

`    `}()

}

-----
**8. iOS 端实现方案**

**策略**：使用 **WKWebView** 加载 Web 页面，实现最小原生开发工作。

**8.1 工程结构**

- 创建 Single View App (Swift)
- 主 ViewController 内嵌一个全屏 WKWebView
- 配置 Info.plist 允许 HTTP 加载（开发阶段）及相册/相机权限

**8.2 开发与部署**

- **开发阶段**：Web 项目运行在本地电脑 http://192.168.x.x:5173，iOS 模拟器/真机通过局域网访问
- **上线/演示阶段**：
  - 将 Vue 打包后的 dist 文件夹拖入 Xcode 项目包内
  - WKWebView 加载本地 index.html（Bundle.main.url）
  - 或者直接加载已部署的线上域名 https://你的域名，App 变为纯浏览器壳，免去后续更新重新打包

**8.3 拍照增强（可选）**

- 在 Web 页面上添加一个“拍照”按钮
- 通过 WKWebView 的 WKScriptMessageHandler 与原生通信
- 原生调用 UIImagePickerController 拍照，将照片转换为 base64 通过 evaluateJavaScript 回传 Web
- Web 端拿到照片后直接走上传流程

这样 iOS 端代码量控制在 200 行以内，无需深入学习 SwiftUI。

-----
**9. 部署方案**

**9.1 服务器环境**

- 云服务器（已购买），安装 Linux (Ubuntu/CentOS)、Nginx、Go 环境、Python 环境

**9.2 后端部署**

- 编译 Go 项目：go build -o wardrobe-server
- 上传二进制文件及静态资源目录
- 使用 systemd 管理服务：

ini

[Unit]

Description=Personal Wardrobe Server

After=network.target

[Service]

ExecStart=/path/to/wardrobe-server

WorkingDirectory=/path/to/

Restart=always

[Install]

WantedBy=multi-user.target

- 确保 Python 抠图服务也在后台运行（ rembg s 守护进程）

**9.3 Nginx 配置**

nginx

server {

`    `listen 80;

`    `server\_name your-domain.com;

`    `*# 前端静态文件*

`    `root /var/www/wardrobe-web;

`    `index index.html;

`    `location / {

`        `try\_files $uri $uri/ /index.html;

`    `}

`    `location /api/ {

`        `proxy\_pass http://127.0.0.1:8080;

`        `proxy\_set\_header Host $host;

`    `}

`    `location /uploads/ {

`        `alias /path/to/server/uploads/;

`    `}

`    `*# 可选：限制上传大小*

`    `client\_max\_body\_size 20m;

}

- 使用 Let’s Encrypt 的 certbot 配置 HTTPS

**9.4 Web 前端部署**

- 本地执行 npm run build，产生 dist 目录
- 上传 dist 内容至 /var/www/wardrobe-web
- 访问域名即可

**9.5 iOS 端发布**

- 真机安装测试即可（无需 App Store 上架），通过 Xcode 直接运行到手机演示
-----
**10. 开发流程与排期建议**

**第一阶段：基础服务搭建（2~3天）**

- 初始化 Go 项目，定义数据模型、JSON 文件读写
- 实现衣物上传接口（仅保存原图，暂不处理抠图）
- 配置静态文件服务

**第二阶段：衣物管理功能（2天）**

- 前端搭建 Vue 项目，路由配置
- 实现“衣柜”页面，衣物列表与分类筛选
- 实现上传页面，对接上传接口
- 轮询状态并更新显示

**第三阶段：抠图异步处理（1~2天）**

- 部署 rembg 服务
- Go 后端添加异步调用逻辑
- 测试完整链路（上传→ pending → done）

**第四阶段：搭配核心功能（3天）**

- 实现搭配编辑器：部位选择、预览区、必选校验
- 实现卡片生成 (html2canvas) 与上传
- 对接搭配创建 API

**第五阶段：历史搭配与辅助功能（1天）**

- 历史搭配列表与“复制搭配”功能
- 搭配删除、查看大图

**第六阶段：iOS 壳与双端联调（1~2天）**

- 创建 iOS 工程，嵌入 WKWebView
- 真机测试拍照增强（可选）
- 调整移动端适配细节

**第七阶段：部署与文档（1天）**

- 配置服务器环境，部署后端与前端
- 绑定域名，设置 HTTPS
- 撰写项目说明/演示视频录制
-----
**11. 补充建议（可扩展但非必需）**

- **多模板卡片**：提供“极简”、“杂志风”搭配卡片版式，使用 Canvas 切换
- **搭配评分/收藏**：本地化的打分，增加互动性
- **衣物统计**：简单统计衣橱数量与总消费（若增加价格字段）
- **数据导入**：支持从电商订单截图导入信息（OCR识别），但技术复杂度过高，暂不纳入
- **备份导出**：后端提供一键导出所有数据为 JSON 的接口，便于备份
-----
**12. 总结**

本文档完整定义了“个人衣柜”项目的功能边界、技术架构、数据流以及开发路径。核心聚焦于 **上传静默抠图** 与 **灵活搭配生成** 两大闭环，通过 Vue + Go + rembg 的低成本技术栈实现双端覆盖，既满足结课作业的实用性要求，又保证了演示效果的完整性。开发过程中可根据实际情况微调优先级，但总体方向清晰，可直接指导编码实现。

