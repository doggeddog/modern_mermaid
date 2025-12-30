package main

import (
	"strings"
	"testing"
)

func TestNormalizeMermaid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Rect",
			input:    "A[Hello]",
			expected: "A[Hello]", // No special chars, no change needed (actually spaces trigger quotes, but here no space)
		},
		{
			name:     "Rect with Space",
			input:    "A[Hello World]",
			expected: `A["Hello World"]`,
		},
		{
			name:     "Rect with Special Chars",
			input:    "A[Is valid?]",
			expected: `A["Is valid?"]`,
		},
		{
			name:     "Rect with Chinese",
			input:    "A[用户操作]",
			expected: `A["用户操作"]`,
		},
		{
			name:     "Already Quoted",
			input:    `A["Already Quoted"]`,
			expected: `A["Already Quoted"]`, // Should not double quote
		},
		{
			name:     "Nested Parens Protection",
			input:    "A[Text(inner)]",
			expected: `A["Text(inner)"]`, // Should be quoted because ( is special
		},
		{
			name:     "Multiple Nodes Line",
			input:    "A[Start] --> B[End Process]",
			expected: `A[Start] --> B["End Process"]`,
		},
		{
			name:     "Arrow Prefix",
			input:    "A-->B(Process: 1)",
			expected: `A-->B("Process: 1")`,
		},
		{
			name:     "Round Node",
			input:    "B(Click?)",
			expected: `B("Click?")`,
		},
		{
			name:     "Rhombus Node",
			input:    "C{Is correct?}",
			expected: `C{"Is correct?"}`,
		},
		{
			name:     "Subroutine Node",
			input:    "D[[Sub Process]]",
			expected: `D[["Sub Process"]]`,
		},
		{
			name:     "Database Node",
			input:    "E[(My DB)]",
			expected: `E[("My DB")]`,
		},
		{
			name:     "Stadium Node",
			input:    "F([End State])",
			expected: `F(["End State"])`,
		},
		{
			name:     "Prefix Protection (False Match)",
			input:    `A["Text(inner)"]`,
			expected: `A["Text(inner)"]`, // Inner ( should not trigger Round Node match because match is inside quote
		},
		{
			name:     "Complex Line",
			input:    `User[用户] --> API{验证?}`,
			expected: `User["用户"] --> API{"验证?"}`,
		},
		{
			name:     "User Reported Case 1",
			input:    `MermaidRender[mermaid.render (解析语法)]`,
			expected: `MermaidRender["mermaid.render (解析语法)"]`,
		},
		{
			name:     "User Reported Case 2 (Already quoted inside)",
			input:    `MermaidRender[mermaid.render ("解析语法")]`,
			expected: `MermaidRender["mermaid.render (` + "`" + `解析语法` + "`" + `)"]`, // Quotes replaced by backticks
		},
		{
			name:     "User Reported Case 3 (Parens inside Rect)",
			input:    `EmitEvent --> Bridge[Wails Bridge (DOM Ready)]`,
			expected: `EmitEvent --> Bridge["Wails Bridge (DOM Ready)"]`,
		},
		{
			name:     "Sequence Diagram (Should Skip)",
			input:    `Alice->>John: Hello John, how are you?`,
			expected: `Alice->>John: Hello John, how are you?`,
		},
		{
			name: "Class Diagram (Should Skip Body)",
			input: `class BankAccount {
    +String owner
    +BigDecimal balance
}`,
			expected: `class BankAccount {
    +String owner
    +BigDecimal balance
}`,
		},
		{
			name:     "ER Diagram",
			input:    `CUSTOMER ||--o{ ORDER : places`,
			expected: `CUSTOMER ||--o{ ORDER : places`,
		},
		{
			name:     "Pie Chart",
			input:    `title What makes a good diagram`,
			expected: `title What makes a good diagram`,
		},
		{
			name:     "Quotes inside Text",
			input:    `A[Say "Hello" please]`,
			expected: `A["Say ` + "`" + `Hello` + "`" + ` please"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeMermaid(tt.input)
			if strings.TrimSpace(got) != strings.TrimSpace(tt.expected) {
				t.Errorf("NormalizeMermaid() = %v, want %v", got, tt.expected)
			}
		})
	}
}
