# Background Implementation Examples

This document lists the background implementations found in the codebase, serving as a reference for creating new backgrounds or themes.

Backgrounds are implemented using CSS properties, primarily `backgroundColor` and `backgroundImage`. They support CSS gradients, patterns, and image URLs (including Data URIs).

## Standalone Backgrounds (`src/utils/backgrounds.ts`)

These backgrounds can be selected independently of the current theme.

### Pure White
```typescript
{
  bgClass: 'bg-white',
  bgStyle: { backgroundColor: '#ffffff' },
}
```

### Pure Black
```typescript
{
  bgClass: 'bg-black',
  bgStyle: { backgroundColor: '#000000' },
}
```

### Light Gray
```typescript
{
  bgClass: 'bg-gray-100',
  bgStyle: { backgroundColor: '#f3f4f6' },
}
```

### Soft Blue
```typescript
{
  bgClass: 'bg-blue-50',
  bgStyle: { backgroundColor: '#eff6ff' },
}
```

### Soft Green
```typescript
{
  bgClass: 'bg-green-50',
  bgStyle: { backgroundColor: '#f0fdf4' },
}
```

### Soft Purple
```typescript
{
  bgClass: 'bg-purple-50',
  bgStyle: { backgroundColor: '#faf5ff' },
}
```

### Blue Gradient
```typescript
{
  bgClass: 'bg-gradient-to-br from-blue-50 to-indigo-100',
  bgStyle: {
    background: 'linear-gradient(to bottom right, #eff6ff, #e0e7ff)',
  },
}
```

### Sunset Gradient
```typescript
{
  bgClass: 'bg-gradient-to-br from-orange-50 to-pink-100',
  bgStyle: {
    background: 'linear-gradient(to bottom right, #fff7ed, #fce7f3)',
  },
}
```

### Dots Pattern
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#ffffff',
    backgroundImage: 'radial-gradient(#e5e7eb 1px, transparent 1px)',
    backgroundSize: '20px 20px',
  },
}
```

### Grid Pattern
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#ffffff',
    backgroundImage: `
      linear-gradient(#e5e7eb 1px, transparent 1px),
      linear-gradient(90deg, #e5e7eb 1px, transparent 1px)
    `,
    backgroundSize: '20px 20px',
  },
}
```

### Diagonal Lines
```typescript
{
  bgClass: 'bg-gray-50',
  bgStyle: {
    backgroundColor: '#f9fafb',
    backgroundImage: `repeating-linear-gradient(
      45deg,
      transparent,
      transparent 10px,
      #e5e7eb 10px,
      #e5e7eb 11px
    )`,
  },
}
```

### Noise Texture (SVG Data URI)
```typescript
{
  bgClass: 'bg-gray-50',
  bgStyle: {
    backgroundColor: '#f9fafb',
    backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='0.05'/%3E%3C/svg%3E")`,
  },
}
```

### Blurred Blue
```typescript
{
  bgClass: 'bg-blue-50',
  bgStyle: {
    backgroundColor: '#eff6ff',
    backgroundImage: 'radial-gradient(circle at 20% 30%, rgba(59, 130, 246, 0.15) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(99, 102, 241, 0.15) 0%, transparent 50%)',
  },
}
```

### Blurred Rainbow
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#ffffff',
    backgroundImage: `
      radial-gradient(circle at 10% 20%, rgba(239, 68, 68, 0.1) 0%, transparent 50%),
      radial-gradient(circle at 90% 30%, rgba(59, 130, 246, 0.1) 0%, transparent 50%),
      radial-gradient(circle at 50% 80%, rgba(16, 185, 129, 0.1) 0%, transparent 50%)
    `,
  },
}
```

### Pop Art
```typescript
{
  bgClass: 'bg-yellow-50',
  bgStyle: {
    backgroundColor: '#fefce8',
    backgroundImage: `
      repeating-linear-gradient(0deg, transparent, transparent 50px, rgba(251, 146, 60, 0.1) 50px, rgba(251, 146, 60, 0.1) 100px),
      repeating-linear-gradient(90deg, transparent, transparent 50px, rgba(236, 72, 153, 0.1) 50px, rgba(236, 72, 153, 0.1) 100px)
    `,
  },
}
```

### Memphis Style (SVG Data URI)
```typescript
{
  bgClass: 'bg-pink-50',
  bgStyle: {
    backgroundColor: '#fdf2f8',
    backgroundImage: `
      url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z' fill='%23ec4899' fill-opacity='0.1'/%3E%3C/g%3E%3C/svg%3E")
    `,
  },
}
```

### Cyberpunk
```typescript
{
  bgClass: 'bg-black',
  bgStyle: {
    backgroundColor: '#000000',
    backgroundImage: `
      linear-gradient(rgba(0, 255, 255, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 0, 255, 0.05) 1px, transparent 1px)
    `,
    backgroundSize: '50px 50px',
  },
}
```

### Cosmic Gradient
```typescript
{
  bgClass: 'bg-gradient-to-br',
  bgStyle: {
    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #00f2fe 100%)',
  },
}
```

### Kawaii
```typescript
{
  bgClass: 'bg-pink-100',
  bgStyle: {
    backgroundColor: '#fce7f3',
    backgroundImage: `
      radial-gradient(circle at 20% 30%, rgba(251, 207, 232, 0.6) 0%, transparent 30%),
      radial-gradient(circle at 80% 20%, rgba(254, 240, 138, 0.5) 0%, transparent 30%),
      radial-gradient(circle at 60% 80%, rgba(191, 219, 254, 0.5) 0%, transparent 30%)
    `,
  },
}
```

### Retro Wave
```typescript
{
  bgClass: 'bg-purple-900',
  bgStyle: {
    backgroundColor: '#581c87',
    backgroundImage: `
      repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(236, 72, 153, 0.3) 2px, rgba(236, 72, 153, 0.3) 4px),
      linear-gradient(180deg, rgba(88, 28, 135, 1) 0%, rgba(147, 51, 234, 0.8) 50%, rgba(236, 72, 153, 0.6) 100%)
    `,
  },
}
```

### Paper Texture (SVG Data URI)
```typescript
{
  bgClass: 'bg-amber-50',
  bgStyle: {
    backgroundColor: '#fffbeb',
    backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='2' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='0.08'/%3E%3C/svg%3E")`,
  },
}
```

### Doodle (SVG Data URI)
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#ffffff',
    backgroundImage: `url("data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3z' fill='%239333ea' fill-opacity='0.08' fill-rule='evenodd'/%3E%3C/svg%3E")`,
  },
}
```

### Minimalist Dark
```typescript
{
  bgClass: 'bg-gray-900',
  bgStyle: {
    backgroundColor: '#111827',
    backgroundImage: 'linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px), linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px)',
    backgroundSize: '40px 40px',
  },
}
```

### Aurora
```typescript
{
  bgClass: 'bg-indigo-950',
  bgStyle: {
    backgroundColor: '#1e1b4b',
    backgroundImage: `
      radial-gradient(ellipse at 30% 20%, rgba(59, 130, 246, 0.3) 0%, transparent 40%),
      radial-gradient(ellipse at 70% 60%, rgba(168, 85, 247, 0.3) 0%, transparent 40%),
      radial-gradient(ellipse at 50% 80%, rgba(34, 211, 238, 0.2) 0%, transparent 40%)
    `,
  },
}
```

### Circuit Board (SVG Data URI)
```typescript
{
  bgClass: 'bg-green-950',
  bgStyle: {
    backgroundColor: '#022c22',
    backgroundImage: `url("data:image/svg+xml,%3Csvg width='100' height='100' viewBox='0 0 100 100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M11 18c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm48 25c3.866 0 7-3.134 7-7s-3.134-7-7-7-7 3.134-7 7 3.134 7 7 7zm-43-7c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm63 31c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM34 90c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zm56-76c1.657 0 3-1.343 3-3s-1.343-3-3-3-3 1.343-3 3 1.343 3 3 3zM12 86c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm28-65c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm23-11c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-6 60c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 1.79 4 4 4zm29 22c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zM32 63c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm57-13c2.76 0 5-2.24 5-5s-2.24-5-5-5-5 2.24-5 5 2.24 5 5 5zm-9-21c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM60 91c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2zM35 41c1.105 0 2-.895 2-2s-.895-2-2-2-2 .895-2 2 .895 2 2 2z' fill='%2310b981' fill-opacity='0.15' fill-rule='evenodd'/%3E%3C/svg%3E")`,
  },
}
```

### Watercolor
```typescript
{
  bgClass: 'bg-blue-50',
  bgStyle: {
    backgroundColor: '#eff6ff',
    backgroundImage: `
      radial-gradient(ellipse at 15% 25%, rgba(191, 219, 254, 0.6) 0%, transparent 50%),
      radial-gradient(ellipse at 85% 35%, rgba(196, 181, 253, 0.5) 0%, transparent 50%),
      radial-gradient(ellipse at 50% 75%, rgba(254, 202, 202, 0.4) 0%, transparent 50%),
      radial-gradient(ellipse at 30% 80%, rgba(187, 247, 208, 0.4) 0%, transparent 50%)
    `,
  },
}
```

## Theme Backgrounds (`src/utils/themes.ts`)

These are the default backgrounds associated with specific themes.

### Linear Light
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundImage: 'radial-gradient(#e5e5e5 1px, transparent 1px)',
    backgroundSize: '20px 20px'
  }
}
```

### Notion
```typescript
{
  bgClass: 'bg-[#051423]',
  bgStyle: {
    backgroundColor: '#051423',
    backgroundImage: `
        linear-gradient(rgba(0, 242, 255, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(0, 242, 255, 0.03) 1px, transparent 1px),
        radial-gradient(circle at 50% 50%, rgba(0, 242, 255, 0.05), transparent 70%)
    `,
    backgroundSize: '40px 40px, 40px 40px, 100% 100%',
    backgroundBlendMode: 'screen'
  }
}
```

### Glassmorphism
```typescript
{
  bgClass: 'bg-gradient-to-br from-purple-50 via-pink-50 to-blue-50',
  bgStyle: {
    background: 'linear-gradient(135deg, #F3E8FF 0%, #FCE7F3 30%, #DBEAFE 60%, #F3E8FF 100%)',
    position: 'relative' as const,
  }
}
```

### Soft Pop
```typescript
{
  bgClass: 'bg-[#f5f5f5]',
  bgStyle: {
    backgroundImage: `
        linear-gradient(#d4d4d4 1px, transparent 1px),
        linear-gradient(90deg, #d4d4d4 1px, transparent 1px)
    `,
    backgroundSize: '20px 20px'
  }
}
```

### Hand Drawn
```typescript
{
  bgClass: 'bg-[#fffef9]',
  bgStyle: {
    backgroundColor: '#fffef9',
    backgroundImage: `
        radial-gradient(circle at 2px 2px, rgba(26, 26, 26, 0.03) 1px, transparent 1px)
    `,
    backgroundSize: '30px 30px'
  }
}
```

### Grafana
```typescript
{
  bgClass: 'bg-[#181B1F]',
  bgStyle: {
    backgroundColor: '#181B1F',
    backgroundImage: `
    linear-gradient(rgba(61, 67, 75, 0.15) 1px, transparent 1px),
    linear-gradient(90deg, rgba(61, 67, 75, 0.15) 1px, transparent 1px)
  `,
    backgroundSize: '24px 24px'
  }
}
```

### Memphis
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#FEFEFE',
    backgroundImage: `
    linear-gradient(45deg, rgba(255, 107, 107, 0.3) 25%, transparent 25%, transparent 75%, rgba(255, 107, 107, 0.3) 75%, rgba(255, 107, 107, 0.3)),
    linear-gradient(45deg, rgba(255, 107, 107, 0.3) 25%, transparent 25%, transparent 75%, rgba(255, 107, 107, 0.3) 75%, rgba(255, 107, 107, 0.3)),
    linear-gradient(45deg, transparent 25%, rgba(254, 202, 87, 0.08) 25%, rgba(254, 202, 87, 0.08) 50%, transparent 50%, transparent),
    linear-gradient(-45deg, rgba(72, 219, 251, 0.06) 25%, transparent 25%, transparent 75%, rgba(72, 219, 251, 0.06) 75%, rgba(72, 219, 251, 0.06)),
    linear-gradient(90deg, rgba(255, 159, 243, 0.04) 1px, transparent 1px),
    linear-gradient(rgba(255, 159, 243, 0.04) 1px, transparent 1px)
  `,
    backgroundSize: '50px 50px, 50px 50px, 25px 25px, 25px 25px, 10px 10px, 10px 10px',
  }
}
```

### Noir
```typescript
{
  bgClass: 'bg-[#0a0a0a]',
  bgStyle: {
    backgroundColor: '#0a0a0a',
    backgroundImage: `
    linear-gradient(135deg, #1a1a1a 0%, #0a0a0a 100%),
    radial-gradient(circle at 50% 20%, rgba(255, 255, 255, 0.05) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(255, 255, 255, 0.03) 0%, transparent 50%)
  `,
  }
}
```

### Material
```typescript
{
  bgClass: 'bg-white',
  bgStyle: {
    backgroundColor: '#ffffff',
  }
}
```

### Aurora
```typescript
{
  bgClass: 'bg-gradient-to-br from-[#667eea] via-[#764ba2] to-[#f093fb]',
  bgStyle: {
    backgroundColor: '#667eea',
    backgroundImage: `
    linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%),
    radial-gradient(circle at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(246, 147, 251, 0.15) 0%, transparent 50%)
  `,
  }
}
```

### Windows 95
```typescript
{
  bgClass: 'bg-[#008080]',
  bgStyle: {
    backgroundColor: '#008080',
    backgroundImage: `
    linear-gradient(45deg, #008080 25%, transparent 25%),
    linear-gradient(-45deg, #008080 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #008080 75%),
    linear-gradient(-45deg, transparent 75%, #008080 75%)
  `,
    backgroundSize: '2px 2px',
    backgroundPosition: '0 0, 0 1px, 1px -1px, -1px 0px',
  }
}
```

### Organic Natural
```typescript
{
  bgClass: 'bg-gradient-to-br from-[#c8d5b9] to-[#8ba888]',
  bgStyle: {
    background: 'linear-gradient(135deg, #c8d5b9 0%, #8ba888 100%)',
  }
}
```

### High Tech
```typescript
{
  bgClass: 'bg-[#0a0f1a]',
  bgStyle: {
    background: 'radial-gradient(circle at center, #0a0f1a, #000000)',
    boxShadow: '0 0 20px rgba(0,255,65,0.3), inset 0 0 20px rgba(0,255,65,0.1)',
  }
}
```

### Kawaii Cute
```typescript
{
  bgClass: 'bg-gradient-to-br from-[#ffe4f1] via-[#ffd4e5] to-[#ffe9f5]',
  bgStyle: {
    background: 'linear-gradient(135deg, #ffe4f1 0%, #ffd4e5 50%, #ffe9f5 100%)',
    position: 'relative' as const,
  }
}
```

### Geometric Collage
```typescript
{
  bgClass: 'bg-[#f5f5f0]',
  bgStyle: {
    backgroundColor: '#f5f5f0',
    backgroundImage: `
    repeating-linear-gradient(
      45deg,
      transparent,
      transparent 10px,
      rgba(74, 144, 226, 0.03) 10px,
      rgba(74, 144, 226, 0.03) 20px
    ),
    repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 10px,
      rgba(74, 144, 226, 0.03) 10px,
      rgba(74, 144, 226, 0.03) 20px
    )
  `,
  }
}
```


