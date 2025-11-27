# SEO 优化指南

本文档介绍了如何优化网站的 SEO，以及如何让搜索引擎更快收录。

## 已完成的优化

### 1. HTML Meta 标签优化
- ✅ 添加了完整的 `<title>` 和 `<meta description>`
- ✅ 添加了关键词标签（keywords）
- ✅ 添加了 Open Graph 标签（社交媒体分享）
- ✅ 添加了 Twitter Card 标签
- ✅ 添加了多语言支持标签（hreflang）
- ✅ 添加了 canonical URL

### 2. 搜索引擎文件
- ✅ 创建了 `robots.txt` - 告诉搜索引擎哪些页面可以抓取
- ✅ 创建了 `sitemap.xml` - 提供网站结构地图

### 3. 结构化数据
- ✅ 添加了 JSON-LD 结构化数据（Schema.org）
- ✅ 标记为 WebApplication 类型
- ✅ 包含了功能列表和评分信息

### 4. 图片优化
- ✅ 创建了自定义 favicon
- ✅ 创建了 Open Graph 图片（用于社交媒体分享）

## 需要手动完成的步骤

### 1. 更新域名信息

在以下文件中，将 `https://yourdomain.com/` 替换为你的实际域名：

- `index.html` - 所有 meta 标签和 canonical URL
- `public/sitemap.xml` - 所有 URL
- `public/robots.txt` - Sitemap URL

**快速替换命令：**
```bash
# 替换为你的实际域名
DOMAIN="https://your-actual-domain.com"

# 在 index.html 中替换
sed -i "s|https://yourdomain.com/|${DOMAIN}/|g" index.html

# 在 sitemap.xml 中替换
sed -i "s|https://yourdomain.com/|${DOMAIN}/|g" public/sitemap.xml

# 在 robots.txt 中替换
sed -i "s|https://yourdomain.com/|${DOMAIN}/|g" public/robots.txt
```

### 2. 更新 sitemap.xml 日期

在 `public/sitemap.xml` 中，将 `<lastmod>` 更新为当前日期：

```xml
<lastmod>2024-12-20</lastmod>  <!-- 使用 YYYY-MM-DD 格式 -->
```

### 3. 提交到搜索引擎

#### Google Search Console
1. 访问 [Google Search Console](https://search.google.com/search-console)
2. 添加你的网站
3. 验证所有权（推荐使用 HTML 标签方法）
4. 提交 sitemap：`https://your-domain.com/sitemap.xml`
5. 请求索引（在 URL 检查工具中）

#### Bing Webmaster Tools
1. 访问 [Bing Webmaster Tools](https://www.bing.com/webmasters)
2. 添加你的网站
3. 验证所有权
4. 提交 sitemap

#### 百度站长平台（针对中国市场）
1. 访问 [百度站长平台](https://ziyuan.baidu.com/)
2. 添加网站
3. 验证所有权
4. 提交 sitemap
5. 使用主动推送功能（API）

### 4. 创建和提交 ping

**Google Ping：**
```bash
curl "https://www.google.com/ping?sitemap=https://your-domain.com/sitemap.xml"
```

**Bing Ping：**
```bash
curl "https://www.bing.com/ping?sitemap=https://your-domain.com/sitemap.xml"
```

### 5. 生成真实的 OG 图片

建议使用工具生成真实的 PNG/JPG 图片替代 SVG：

**推荐尺寸：**
- Open Graph: 1200 x 630 px
- Twitter Card: 1200 x 675 px

**在线工具：**
- [Canva](https://www.canva.com/) - 免费图片编辑
- [Figma](https://www.figma.com/) - 专业设计工具
- [OG Image Generator](https://og-image.vercel.app/) - 快速生成

将生成的图片命名为 `og-image.png` 并放在 `public/` 目录。

### 6. 创建 screenshot.png

为结构化数据创建一个应用截图（推荐尺寸：1920 x 1080 px），放在 `public/` 目录。

## 进阶 SEO 优化

### 1. 创建博客/文档页面

搜索引擎喜欢内容丰富的网站。考虑添加：

- 使用教程
- 示例集合
- 最佳实践
- 常见问题（FAQ）
- 博客文章

### 2. 性能优化

**网站速度是 SEO 排名因素：**

```bash
# 安装分析工具
npm install -D vite-plugin-compression

# 在 vite.config.ts 中启用压缩
import compression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    compression({ algorithm: 'gzip' }),
    compression({ algorithm: 'brotliCompress', ext: '.br' })
  ]
})
```

### 3. 添加 Prerender/SSR

对于 SEO，静态生成或服务器端渲染更友好：

**选项 1: 使用 Vite SSG**
```bash
npm install -D vite-ssg
```

**选项 2: 使用 Prerender**
```bash
npm install -D vite-plugin-prerender
```

### 4. 添加 manifest.json (PWA)

创建 `public/manifest.json`：

```json
{
  "name": "Mermaid Advanced",
  "short_name": "Mermaid",
  "description": "Online Diagram Editor",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#6366f1",
  "icons": [
    {
      "src": "/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icon-512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

然后在 `index.html` 中添加：
```html
<link rel="manifest" href="/manifest.json">
<meta name="theme-color" content="#6366f1">
```

### 5. 设置 CDN 和缓存

**推荐 CDN：**
- Cloudflare（免费）
- Vercel Edge Network
- AWS CloudFront

**缓存策略：**
```
# .htaccess 或 nginx 配置
<IfModule mod_expires.c>
  ExpiresActive On
  ExpiresByType image/* "access plus 1 year"
  ExpiresByType text/css "access plus 1 month"
  ExpiresByType application/javascript "access plus 1 month"
</IfModule>
```

### 6. 实现 OpenSearch

创建 `public/opensearch.xml`：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">
  <ShortName>Mermaid Advanced</ShortName>
  <Description>Search Mermaid diagram examples</Description>
  <Url type="text/html" template="https://your-domain.com/?q={searchTerms}"/>
</OpenSearchDescription>
```

### 7. 添加面包屑导航

在 JSON-LD 中添加面包屑：

```json
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [{
    "@type": "ListItem",
    "position": 1,
    "name": "Home",
    "item": "https://your-domain.com"
  }]
}
```

## 内容优化建议

### 关键词策略

**主要关键词：**
- mermaid diagram editor
- flowchart maker online
- sequence diagram generator
- free diagram tool
- online chart editor

**长尾关键词：**
- how to create flowchart online free
- mermaid js online editor
- best free diagram software
- real-time diagram preview

### 内容更新频率

- 每周添加新的示例
- 每月发布教程或博客文章
- 及时更新文档

### 外部链接建设

1. 在 GitHub 创建项目页面
2. 发布到 Product Hunt
3. 在 Reddit、Hacker News 分享
4. 写技术博客并链接回网站
5. 在社交媒体分享

## 监控和分析

### 工具推荐

1. **Google Analytics** - 已配置（见 GOOGLE_ANALYTICS_README.md）
2. **Google Search Console** - 监控搜索表现
3. **Bing Webmaster Tools** - Bing 搜索优化
4. **PageSpeed Insights** - 性能分析
5. **Lighthouse** - 综合评分

### 定期检查

- 每周检查 Search Console 的索引状态
- 每月分析流量来源和关键词
- 季度性能审计
- 修复所有爬虫错误

## 快速检查清单

部署前检查：

- [ ] 所有域名已更新为实际域名
- [ ] sitemap.xml 日期已更新
- [ ] 已创建真实的 OG 图片
- [ ] 已创建应用截图
- [ ] robots.txt 已配置正确
- [ ] 已构建并测试生产版本
- [ ] 已提交到 Google Search Console
- [ ] 已提交到 Bing Webmaster Tools
- [ ] 已设置 Google Analytics
- [ ] 已测试社交媒体分享预览

## 测试工具

**测试你的 SEO：**

1. [Google Rich Results Test](https://search.google.com/test/rich-results)
2. [Facebook Sharing Debugger](https://developers.facebook.com/tools/debug/)
3. [Twitter Card Validator](https://cards-dev.twitter.com/validator)
4. [Schema Markup Validator](https://validator.schema.org/)
5. [PageSpeed Insights](https://pagespeed.web.dev/)

## 预期效果

**时间线：**
- **1-3天：** Google 开始爬取
- **1-2周：** 网站被索引
- **2-4周：** 开始出现在搜索结果中
- **1-3个月：** 排名逐步提升

**加速收录的方法：**
1. 主动提交 URL 到 Search Console
2. 在社交媒体分享
3. 从其他网站链接
4. 创建高质量内容
5. 保持网站活跃更新

## 需要帮助？

如果需要进一步的 SEO 咨询或实施帮助，可以考虑：

1. 聘请 SEO 专家
2. 使用 SEO 工具（Ahrefs, SEMrush, Moz）
3. 参考 Google 的 SEO 入门指南
4. 加入网站管理员社区

## 相关资源

- [Google SEO 入门指南](https://developers.google.com/search/docs/beginner/seo-starter-guide)
- [Moz 初学者 SEO 指南](https://moz.com/beginners-guide-to-seo)
- [Schema.org 文档](https://schema.org/)
- [Open Graph 协议](https://ogp.me/)

