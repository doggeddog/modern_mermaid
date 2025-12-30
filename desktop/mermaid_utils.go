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
// Matches quoted strings, handling escaped quotes
var existingQuoteRegex = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)
var tokenRegex = regexp.MustCompile(`^__MQ_\d+__$`)

func init() {
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

	prefixPattern := `(?:^|\s|&|[-=>])`

	for _, def := range definitions {
		openEsc := regexp.QuoteMeta(def.Open)
		closeEsc := regexp.QuoteMeta(def.Close)
		pattern := fmt.Sprintf(`(%s[A-Za-z0-9_-]+)\s*%s(.*?)%s`, prefixPattern, openEsc, closeEsc)
		
		rule := &ShapeRule{
			Open:  def.Open,
			Close: def.Close,
			re:    regexp.MustCompile(pattern),
		}
		shapeRules = append(shapeRules, rule)
	}
}

func NormalizeMermaid(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" || strings.HasPrefix(trimLine, "%%") {
			result = append(result, line)
			continue
		}

		maskMap := make(map[string]string)
		maskCounter := 0
		
		maskFunc := func(s string) string {
			token := fmt.Sprintf("__MQ_%d__", maskCounter)
			maskCounter++
			maskMap[token] = s
			return token
		}

		// Helper to unmask a string fully
		unmaskStr := func(s string) string {
			// Simple iterative unmasking
			// Since our tokens don't overlap in the string content (they are replaced),
			// and we are unmasking a specific substring 's', we can just loop.
			// But careful about order if 's' contains multiple tokens.
			// strings.ReplaceAll for each token in the map?
			// The map might be large, but usually small per line.
			res := s
			// Iterate in reverse order of creation to handle nesting
			for i := maskCounter - 1; i >= 0; i-- {
				token := fmt.Sprintf("__MQ_%d__", i)
				if val, ok := maskMap[token]; ok {
					res = strings.ReplaceAll(res, token, val)
				}
			}
			return res
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
				content := groups[2]

				// Logic:
				// If content IS exactly one Token, it means it's fully quoted -> Skip.
				// If content contains Token mixed with text -> Unmask, Escape, Quote.
				// If content has no Token but needs quotes -> Quote.
				
				isFullToken := tokenRegex.MatchString(strings.TrimSpace(content))
				
				if isFullToken {
					// Already fully quoted, e.g. ["Text"] -> [__MQ_0__]
					// Preserve as is, but mask the whole shape
					newShape := fmt.Sprintf(`%s%s%s`, rule.Open, content, rule.Close)
					token := maskFunc(newShape)
					return fmt.Sprintf(`%s%s`, fullPrefixWithID, token)
				}
				
				// Check if needs quotes (mixed content or special chars)
				// Note: If content contains `__MQ_`, it implies special chars/quotes, so we definitely need to quote the outer wrapper.
				// But we must UNMASK first to get the raw text with original quotes.
				
				rawContent := content
				if strings.Contains(content, "__MQ_") {
					rawContent = unmaskStr(content)
				}
				
				if needsQuotes(rawContent) || strings.Contains(content, "__MQ_") {
					// Replace existing quotes with backticks ` to avoid syntax errors
					// (Mermaid flowcharts don't always support \" escaping well)
					escapedContent := strings.ReplaceAll(rawContent, `"`, "`")
					
					// Wrap in quotes
					newShape := fmt.Sprintf(`%s"%s"%s`, rule.Open, escapedContent, rule.Close)
					token := maskFunc(newShape)
					return fmt.Sprintf(`%s%s`, fullPrefixWithID, token)
				}
				
				// No quotes needed (simple text)
				newShape := fmt.Sprintf(`%s%s%s`, rule.Open, content, rule.Close)
				token := maskFunc(newShape)
				return fmt.Sprintf(`%s%s`, fullPrefixWithID, token)
			})
		}

		// 3. Final Unmask of the line
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
