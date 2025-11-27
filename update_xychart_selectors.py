#!/usr/bin/env python3
import re

# Read the themes file
with open('src/utils/themes.ts', 'r', encoding='utf-8') as f:
    content = f.read()

# Replace patterns
replacements = [
    (r'svg\[aria-roledescription="xychart"\] \.plot-line-0', '.line-plot-0 path'),
    (r'svg\[aria-roledescription="xychart"\] \.plot-line-1', '.line-plot-1 path'),
    (r'svg\[aria-roledescription="xychart"\] \.plot-line-2', '.line-plot-2 path'),
    (r'svg\[aria-roledescription="xychart"\] \.plot-bar-0', '.bar-plot-0 rect'),
    (r'svg\[aria-roledescription="xychart"\] \.plot-bar-1', '.bar-plot-1 rect'),
    (r'svg\[aria-roledescription="xychart"\] \.plot-bar-2', '.bar-plot-2 rect'),
    (r'svg\[aria-roledescription="xychart"\] \.chart-title', '.chart-title text'),
    (r'svg\[aria-roledescription="xychart"\] \.axis-label', '.left-axis .title text, .bottom-axis .title text'),
    (r'svg\[aria-roledescription="xychart"\] \.legend text', '.legend text'),
    (r'svg\[aria-roledescription="xychart"\] \.tick text', '.left-axis .label text, .bottom-axis .label text'),
    (r'svg\[aria-roledescription="xychart"\] \.tick', '.ticks path'),
]

for old, new in replacements:
    content = re.sub(old, new, content)

# Write back
with open('src/utils/themes.ts', 'w', encoding='utf-8') as f:
    f.write(content)

print("Updated all XYChart selectors!")

