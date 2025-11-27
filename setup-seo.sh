#!/bin/bash

# SEO 设置脚本
# 用于快速配置网站的 SEO 相关设置

set -e

echo "🚀 Mermaid Advanced - SEO 设置向导"
echo "=================================="
echo ""

# 检查是否有参数
if [ "$#" -eq 0 ]; then
    echo "请输入你的网站域名（不带尾部斜杠）："
    echo "例如: https://mermaid-advanced.com"
    read -r DOMAIN
else
    DOMAIN=$1
fi

# 验证域名格式
if [[ ! $DOMAIN =~ ^https?:// ]]; then
    echo "❌ 错误：域名必须以 http:// 或 https:// 开头"
    exit 1
fi

# 移除尾部斜杠（如果有）
DOMAIN=${DOMAIN%/}

echo ""
echo "📝 使用域名: $DOMAIN"
echo ""

# 备份原始文件
echo "📦 备份原始文件..."
cp index.html index.html.backup
cp public/sitemap.xml public/sitemap.xml.backup
cp public/robots.txt public/robots.txt.backup

# 替换 index.html 中的域名
echo "🔄 更新 index.html..."
sed -i.tmp "s|https://yourdomain.com|${DOMAIN}|g" index.html
rm -f index.html.tmp

# 替换 sitemap.xml 中的域名
echo "🔄 更新 sitemap.xml..."
sed -i.tmp "s|https://yourdomain.com|${DOMAIN}|g" public/sitemap.xml
rm -f public/sitemap.xml.tmp

# 替换 robots.txt 中的域名
echo "🔄 更新 robots.txt..."
sed -i.tmp "s|https://yourdomain.com|${DOMAIN}|g" public/robots.txt
rm -f public/robots.txt.tmp

# 更新 sitemap 日期为今天
CURRENT_DATE=$(date +%Y-%m-%d)
echo "📅 更新 sitemap 日期为: $CURRENT_DATE"
sed -i.tmp "s|<lastmod>.*</lastmod>|<lastmod>${CURRENT_DATE}</lastmod>|g" public/sitemap.xml
rm -f public/sitemap.xml.tmp

echo ""
echo "✅ SEO 配置完成！"
echo ""
echo "📋 后续步骤："
echo "1. 生成 OG 图片（1200x630px）并保存为 public/og-image.png"
echo "2. 生成应用截图（1920x1080px）并保存为 public/screenshot.png"
echo "3. 运行 'pnpm build' 构建生产版本"
echo "4. 部署后，提交 sitemap 到:"
echo "   - Google: https://search.google.com/search-console"
echo "   - Bing: https://www.bing.com/webmasters"
echo "5. Ping 搜索引擎:"
echo "   curl 'https://www.google.com/ping?sitemap=${DOMAIN}/sitemap.xml'"
echo "   curl 'https://www.bing.com/ping?sitemap=${DOMAIN}/sitemap.xml'"
echo ""
echo "📖 详细指南请查看: SEO_GUIDE.md"
echo ""
echo "💡 提示：备份文件已保存为 *.backup"
echo ""

