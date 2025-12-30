package main

import (
	"fmt"
	"regexp"
	"strings"
)

// ShapeRule defines a mermaid node shape configuration
type ShapeRule struct {
	Open  string
	Close string
	re    *regexp.Regexp
}

var shapeRules []*ShapeRule
// Matches quoted strings, handling escaped quotes: " text \" inner "
var existingQuoteRegex = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)

func init() {
	// Initialize rules. Order matters: Process longer symbols first.
	definitions := []struct {
		Open, Close string
	}{
		{"[[", "]]"}, // Subroutine
		{"{{", "}}"}, // Hexagon
		{"([", "])"}, // Stadium
		{"[(", ")]"}, // DB
		{"[", "]"},   // Rect
		{"((", "))"}, // Circle
		{"(", ")"},   // Round
		{">", "]"},   // Asymmetric
		{"{", "}"},   // Rhombus
	}

	// Prefix pattern: ensures ID is preceded by Start-of-line, Whitespace, Ampersand, or Arrow/Line chars
	prefixPattern := `(?:^|\s|&|[-=>])`

	for _, def := range definitions {
		openEsc := regexp.QuoteMeta(def.Open)
		closeEsc := regexp.QuoteMeta(def.Close)
		
		// Pattern: Prefix + ID + Open + Content + Close
		// We use non-greedy matching .*? for content
		pattern := fmt.Sprintf(`(%s[A-Za-z0-9_-]+)\s*%s(.*?)%s`, prefixPattern, openEsc, closeEsc)
		
		rule := &ShapeRule{
			Open:  def.Open,
			Close: def.Close,
			re:    regexp.MustCompile(pattern),
		}
		shapeRules = append(shapeRules, rule)
	}
}

// NormalizeMermaid scans the input mermaid code and adds quotes to node labels.
func NormalizeMermaid(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" || strings.HasPrefix(trimLine, "%%") {
			result = append(result, line)
			continue
		}

		// Masking State
		maskMap := make(map[string]string)
		maskCounter := 0
		
		maskFunc := func(s string) string {
			token := fmt.Sprintf("__MQ_%d__", maskCounter)
			maskCounter++
			maskMap[token] = s
			return token
		}

		// 1. Mask existing quotes
		processedLine := existingQuoteRegex.ReplaceAllStringFunc(line, func(match string) string {
			return maskFunc(match)
		})

		// 2. Process all shape rules
		for _, rule := range shapeRules {
			processedLine = rule.re.ReplaceAllStringFunc(processedLine, func(match string) string {
				groups := rule.re.FindStringSubmatch(match)
				if len(groups) < 3 {
					return match
				}

				fullPrefixWithID := groups[1]
				content := groups[2] // This content might contain Tokens

				// Logic:
				// 1. If content contains a Mask Token (meaning it had quotes), we SKIP modifying it.
				//    We assume if the user put quotes, they handled it.
				//    But we MUST still mask the whole shape to protect it from shorter rules.
				// 2. If content needs quotes, we Add Quotes.
				// 3. In both cases, we mask the "Shape Part" (ID + Open + Content + Close) -> ID + Token
				//    Wait, we match Prefix+ID.
				//    If we replace the whole match with Prefix + Token, we hide ID.
				//    But hiding ID is fine for subsequent rules (they won't match inside).
				//    Actually, we should mask the whole Node Declaration?
				//    Yes. `A[B]` -> `__MQ__`.
				//    This is safest.
				
				// Reconstruct the shape string (without prefix logic, which is just lookbehind context in a way)
				// The `match` is the full string.
				// `groups[1]` is Prefix+ID.
				// We want to replace `match` with `groups[1] + MaskedShape`.
				// But `match` includes the prefix.
				// So we take the `Open + Content + Close` part.
				// Wait, `match` = `PrefixID` + `Open` + `Content` + `Close`.
				// We can reconstruct it.
				
				var newShape string
				
				if strings.Contains(content, "__MQ_") {
					// Already has quotes, preserve content
					newShape = fmt.Sprintf(`%s%s%s`, rule.Open, content, rule.Close)
				} else if needsQuotes(content) {
					// Add quotes
					newShape = fmt.Sprintf(`%s"%s"%s`, rule.Open, content, rule.Close)
				} else {
					// No quotes needed
					newShape = fmt.Sprintf(`%s%s%s`, rule.Open, content, rule.Close)
				}
				
				// Mask the shape part
				token := maskFunc(newShape)
				
				// Return Prefix+ID + Token
				return fmt.Sprintf(`%s%s`, fullPrefixWithID, token)
			})
		}

		// 3. Unmask
		// We iterate in REVERSE order because later tokens (e.g. from shape matches)
		// might contain earlier tokens (e.g. from existing quotes).
		for i := maskCounter - 1; i >= 0; i-- {
			token := fmt.Sprintf("__MQ_%d__", i)
			if val, ok := maskMap[token]; ok {
				processedLine = strings.ReplaceAll(processedLine, token, val)
			}
		}
		
		result = append(result, processedLine)
	}

	return strings.Join(result, "\n")
}

func needsQuotes(s string) bool {
	if s == "" {
		return true
	}
	specialChars := `()[]{}:?` 
	if strings.ContainsAny(s, specialChars) {
		return true
	}
	if strings.Contains(s, " ") {
		return true
	}
	
	isASCII := true
	for _, c := range s {
		if c > 127 {
			isASCII = false
			break
		}
	}
	if !isASCII {
		return true
	}

	return false
}
