# Mermaid 代码正规化 (Normalization) 设计方案

## 1. 背景与目标

用户从外部来源（如 GitHub Issue、文档）复制的 Mermaid 代码可能存在语法不严谨的情况，最典型的问题是**节点或连线标签中包含特殊字符（如 `()`, `[]`, `?`, `:`）但未被双引号包裹**。这会导致 Mermaid 解析器报错。

本方案旨在后端提取代码后，增加一个 **"正规化 (Normalization)"** 步骤，自动检测并修复这些潜在的语法问题。

## 2. 技术方案：基于掩码保护 (Masking) 的正则替换

为了解决正则匹配容易误伤（如误匹配嵌套结构或已引用的文本）的问题，我们最终采用了 **掩码保护 + 上下文感知正则** 的混合策略。

### 2.1 核心策略

1.  **掩码保护 (Masking Strategy)**:
    *   在进行正则匹配前，先将代码中 **已存在的引号内容** 替换为临时 Token (如 `__MQ_0__`)。
    *   当正则匹配并添加新引号后，也立即将新生成的引号内容替换为 Token。
    *   **作用**: 防止后续规则（或短符号规则）错误地匹配到已引用字符串内部的内容。

2.  **转义策略**:
    *   如果需要为混合内容的节点（如 `A[Say "Hello"]`）添加外层双引号，我们会先将**内部的双引号替换为反引号 (` ` `)**。
    *   **结果**: `A[Say "Hello"]` -> `A["Say `Hello`"]`。这是为了在保证 Flowchart 语法合法的同时，最大程度保留原义且避免复杂的转义序列 (`\"`) 可能带来的解析兼容性问题。

3.  **严格的前缀断言**: 利用正则，限定“节点ID”的前面必须是合法的语法边界（行首、空格、箭头等），且排除 `.`。

### 2.2 伪代码逻辑 (Go 实现)

```go
func NormalizeMermaid(input string) string {
    // 1. 初始化掩码状态
    
    // 2. 预处理：保护已存在的引号
    input = MaskExistingQuotes(input)

    // 3. 遍历形状规则 (双字符优先)
    for _, rule := range shapeRules {
        input = rule.re.ReplaceAllStringFunc(input, func(match) {
            // 如果 Content 需要引号:
            // 1. 将内部的 " 替换为 `
            // 2. 添加外层双引号
            // 3. 掩码保护
            return Prefix + Token
        })
    }

    // 4. 反向还原 (Reverse Unmasking)
    return Unmask(input)
}
```

## 3. 测试验证 (Verified Cases)

已通过以下复杂场景的单元测试：

1.  **嵌套括号**: `Bridge[Wails Bridge (DOM Ready)]` -> `Bridge["Wails Bridge (DOM Ready)"]`
2.  **已引用的特殊字符**: `MermaidRender[mermaid.render ("解析语法")]` -> `MermaidRender["mermaid.render (`解析语法`)"]`
3.  **混合引号**: `A[Say "Hello"]` -> `A["Say `Hello`"]`
4.  **多图表类型兼容**: 自动跳过 Sequence Diagram, Class Diagram Body 等不适用的语法结构。

## 4. 结论

该方案通过 **正则扫描** 与 **状态掩码** 的结合，并采用 **反引号替换** 策略处理内部引号，实现了高鲁棒性的代码正规化。
