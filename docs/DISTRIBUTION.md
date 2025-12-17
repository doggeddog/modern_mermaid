# 应用发布指南 (Distribution Guide)

## 1. 构建 Release 版本

使用以下命令构建生产版本：

```bash
make build
```

这会在 `desktop/build/bin/` 目录下生成可执行文件（在 macOS 上是 `.app` 包）。

## 2. macOS 签名与公证 (Signing & Notarization)

要分发 macOS 应用以免受 Gatekeeper 拦截，您需要拥有 Apple Developer Program 账号，并对应用进行签名和公证。

### 准备工作

1.  登录 [Apple Developer](https://developer.apple.com/) 并创建一个 "Developer ID Application" 证书。
2.  将证书下载并安装到您的 Keychain 中。
3.  获取您的 Team ID (在开发者账号右上角可以找到)。
4.  在 Apple ID 管理页面创建一个应用专用密码 (App-Specific Password) 用于公证工具登录。

### 签名步骤 (命令行方式)

构建完成后，在项目根目录执行以下步骤：

#### 1. 签名 (Code Sign)

```bash
# 替换 "Developer ID Application: Your Name (TEAMID)" 为您的证书名称（在 Keychain Access 中查看）
codesign --force --options runtime --deep --sign "Developer ID Application: Your Name (TEAMID)" "desktop/build/bin/Modern Mermaid.app"
```

验证签名是否成功：
```bash
codesign -dv --verbose=4 "desktop/build/bin/Modern Mermaid.app"
```

#### 2. 打包 (Zip)

公证服务需要上传 ZIP 文件：

```bash
/usr/bin/ditto -c -k --keepParent "desktop/build/bin/Modern Mermaid.app" "ModernMermaid.zip"
```

#### 3. 公证 (Notarize)

使用 `notarytool` 提交公证请求。这可能需要几分钟。

```bash
# 替换 <apple-id>, <password>, <team-id>
xcrun notarytool submit "ModernMermaid.zip" --apple-id "your@email.com" --password "your-app-specific-password" --team-id "YOUR_TEAM_ID" --wait
```

如果成功，终端会显示 "Accepted"。

#### 4. 盖章 (Staple)

将公证票据附加到原始应用文件上，这样离线用户也能验证。

```bash
xcrun stapler staple "desktop/build/bin/Modern Mermaid.app"
```

验证盖章：
```bash
spctl --assess --verbose "desktop/build/bin/Modern Mermaid.app"
```

### 分发

现在，您可以分发 `desktop/build/bin/Modern Mermaid.app`，或者将其制作成 DMG 镜像。
可以使用 `create-dmg` 等工具制作 DMG。

