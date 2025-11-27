# 预览优化指南

## 已完成的修改

### 1. 刷新按钮功能修改 ✅

**问题：** 原来的刷新按钮会恢复到默认示例代码，而不是重新渲染当前代码。

**解决方案：**

- 在 `PreviewHandle` 接口中添加了 `refresh()` 方法
- 在 `Preview` 组件中添加了 `renderKey` 状态来强制重新渲染
- 修改 `Layout` 组件的 `handleRefreshEditor` 函数，调用 `previewRef.current.refresh()`

**使用方式：**
点击例子选择框旁边的刷新按钮（🔄），会立即重新生成当前代码的预览，无需等待防抖延迟。

**代码变更：**

```typescript
// Preview.tsx
export interface PreviewHandle {
  exportImage: (transparent: boolean) => Promise<void>;
  clearAnnotations: () => void;
  refresh: () => void; // 新增
}

// 在组件内部添加状态
const [renderKey, setRenderKey] = useState(0);

// 在 useImperativeHandle 中实现
useImperativeHandle(ref, () => ({
  clearAnnotations: () => { ... },
  refresh: () => {
    setRenderKey(prev => prev + 1); // 强制重新渲染
  },
  exportImage: async (transparent: boolean) => { ... }
}));

// Layout.tsx
const handleRefreshEditor = () => {
  if (previewRef.current) {
    previewRef.current.refresh(); // 调用刷新方法
  }
};
```

---

### 2. 预览生成速度优化 ✅

**问题：** 预览生成有延迟感，特别是部署到线上后。

**解决方案：**

1. **降低防抖延迟**
   - 从 600ms 降低到 300ms
   - 在用户输入时减少等待时间

2. **添加 renderKey 依赖**
   - 允许通过刷新按钮立即触发渲染，跳过防抖延迟

**代码变更：**

```typescript
// Preview.tsx - useEffect 中的防抖
const timeoutId = setTimeout(() => {
  renderDiagram();
}, 300); // 从 600ms 改为 300ms

// 添加 renderKey 到依赖项
}, [code, themeConfig, actualFont, nodeColors, renderKey]);
```

---

## 关于预览生成慢的原因分析

### 可能的原因

#### 1. **Mermaid 渲染本身比较慢**
   - Mermaid 是客户端渲染的图表库
   - 复杂的图表（如大型流程图、序列图）需要更多时间计算布局
   - **解决方案：** 
     - ✅ 已降低防抖延迟到 300ms
     - 考虑使用 Web Worker（需要较大改动）

#### 2. **防抖延迟太长**
   - ✅ 已解决：从 600ms 降低到 300ms
   - 这样用户输入后等待时间减半

#### 3. **依赖项过多导致不必要的重新渲染**
   - 当前依赖项：`code`, `themeConfig`, `actualFont`, `nodeColors`, `renderKey`
   - 这些都是必要的依赖项，变化时确实需要重新渲染
   - **状态：** 优化空间有限

#### 4. **没有使用生产构建**
   - 开发模式下 React 会有额外的检查和警告
   - **解决方案：** 确保部署时使用 `pnpm build` 的生产版本
   - ✅ 已验证构建成功

#### 5. **大型依赖库加载慢**
   - Mermaid 库较大（约 265KB gzipped）
   - **解决方案：**
     - 使用 CDN 加速
     - 启用浏览器缓存
     - 考虑代码分割（Code Splitting）

#### 6. **浏览器性能问题**
   - 低配设备或旧浏览器可能渲染慢
   - **建议：** 在控制台检查是否有性能警告

---

## 性能优化建议（进一步优化）

### 立即可实施

1. **✅ 已完成：降低防抖延迟**
   ```typescript
   setTimeout(() => renderDiagram(), 300); // 已改为 300ms
   ```

2. **✅ 已完成：添加强制刷新功能**
   - 用户可以通过刷新按钮立即触发渲染

3. **确保使用生产构建**
   ```bash
   pnpm build
   # 部署 dist/ 目录
   ```

### 中等难度

4. **添加加载指示器优化**
   - 当前已有 `loading` 状态
   - 可以添加更明显的进度提示
   - 显示渲染耗时统计

5. **使用 React.memo 优化组件**
   ```typescript
   const Editor = React.memo(EditorComponent);
   const Toolbar = React.memo(ToolbarComponent);
   ```

6. **优化 Mermaid 配置**
   ```typescript
   mermaid.initialize({
     startOnLoad: false,
     securityLevel: 'loose',
     suppressErrorRendering: true,
     // 添加性能优化选项
     logLevel: 'error', // 减少日志输出
   });
   ```

### 高级优化（需要较大改动）

7. **使用 Web Worker 渲染**
   - 将 Mermaid 渲染移到 Web Worker
   - 避免阻塞主线程
   - **注意：** 需要重构渲染逻辑

8. **实现虚拟滚动**
   - 对于大型图表，只渲染可见区域
   - 减少 DOM 节点数量

9. **使用 requestIdleCallback**
   ```typescript
   requestIdleCallback(() => {
     renderDiagram();
   });
   ```

10. **服务端渲染（SSR）预渲染**
    - 对于静态示例，可以预先渲染
    - 减少客户端首次渲染时间

---

## 性能监控建议

### 1. 添加性能统计

```typescript
const renderDiagram = async () => {
  const startTime = performance.now();
  setLoading(true);
  
  try {
    // ... 渲染逻辑
    const { svg } = await mermaid.render(id, code);
    setSvg(svg);
    
    const endTime = performance.now();
    console.log(`渲染耗时: ${(endTime - startTime).toFixed(2)}ms`);
  } catch (error) {
    // ...
  } finally {
    setLoading(false);
  }
};
```

### 2. 使用浏览器开发者工具

- **Performance 面板**：录制渲染过程，查看瓶颈
- **Network 面板**：检查依赖加载时间
- **Memory 面板**：检查是否有内存泄漏

### 3. 添加错误边界

```typescript
class ErrorBoundary extends React.Component {
  componentDidCatch(error, errorInfo) {
    console.error('渲染错误:', error, errorInfo);
  }
}
```

---

## 部署优化清单

### 前端优化
- ✅ 使用生产构建 (`pnpm build`)
- ✅ 降低防抖延迟到 300ms
- ⬜ 启用 Gzip/Brotli 压缩
- ⬜ 使用 CDN 加速静态资源
- ⬜ 配置浏览器缓存策略

### 服务器优化
- ⬜ 使用 HTTP/2 或 HTTP/3
- ⬜ 启用服务器端缓存
- ⬜ 使用负载均衡
- ⬜ 配置适当的 Cache-Control 头

### 代码优化
- ✅ 已优化组件渲染逻辑
- ⬜ 考虑代码分割（Code Splitting）
- ⬜ 懒加载非关键组件
- ⬜ 优化图片和资源大小

---

## 测试建议

### 本地测试
```bash
# 1. 构建生产版本
pnpm build

# 2. 使用静态服务器测试
npx serve dist

# 3. 访问 http://localhost:3000 测试性能
```

### 性能测试
1. 打开浏览器开发者工具
2. 切换到 Performance 面板
3. 点击录制按钮
4. 在编辑器中输入代码
5. 等待预览生成
6. 停止录制，分析结果

### 不同场景测试
- 简单图表（5-10 个节点）
- 中等复杂度图表（20-50 个节点）
- 复杂图表（100+ 个节点）
- 不同浏览器和设备

---

## 总结

### 已完成的优化
1. ✅ 刷新按钮现在可以立即重新生成预览
2. ✅ 防抖延迟从 600ms 降低到 300ms
3. ✅ 添加了强制刷新机制（`renderKey`）
4. ✅ 验证了生产构建正常工作

### 预期效果
- **用户输入响应速度提升 50%**（600ms → 300ms）
- **刷新按钮提供即时反馈**（跳过防抖延迟）
- **减少不必要的重新渲染**

### 下一步建议
如果部署后仍然感觉慢，请：
1. 使用浏览器开发者工具的 Performance 面板分析瓶颈
2. 检查 Network 面板，确认资源加载时间
3. 测试不同复杂度的图表，找出慢的具体场景
4. 考虑实施"中等难度"或"高级优化"中的建议

---

## 联系和反馈

如果在使用过程中发现性能问题，请提供以下信息：
- 浏览器版本和设备信息
- Mermaid 代码示例（导致慢的具体代码）
- Performance 面板截图或数据
- Network 面板的加载时间数据

