export interface BackgroundStyle {
  id: string;
  name: {
    'en': string;
    'zh-CN': string;
    'zh-TW': string;
    'ja': string;
    'es': string;
    'pt': string;
  };
  bgClass: string;
  bgStyle?: React.CSSProperties;
}

export const backgrounds: BackgroundStyle[] = [
  {
    id: 'default',
    name: {
      'en': 'Default (Theme)',
      'zh-CN': '默认（主题）',
      'zh-TW': '預設（主題）',
      'ja': 'デフォルト（テーマ）',
      'es': 'Por defecto (Tema)',
      'pt': 'Padrão (Tema)',
    },
    bgClass: '',
    bgStyle: undefined,
  },
  {
    id: 'white',
    name: {
      'en': 'Pure White',
      'zh-CN': '纯白色',
      'zh-TW': '純白色',
      'ja': '純白',
      'es': 'Blanco puro',
      'pt': 'Branco puro',
    },
    bgClass: 'bg-white',
    bgStyle: { backgroundColor: '#ffffff' },
  },
  {
    id: 'black',
    name: {
      'en': 'Pure Black',
      'zh-CN': '纯黑色',
      'zh-TW': '純黑色',
      'ja': '純黒',
      'es': 'Negro puro',
      'pt': 'Preto puro',
    },
    bgClass: 'bg-black',
    bgStyle: { backgroundColor: '#000000' },
  },
  {
    id: 'gray',
    name: {
      'en': 'Light Gray',
      'zh-CN': '浅灰色',
      'zh-TW': '淺灰色',
      'ja': 'ライトグレー',
      'es': 'Gris claro',
      'pt': 'Cinza claro',
    },
    bgClass: 'bg-gray-100',
    bgStyle: { backgroundColor: '#f3f4f6' },
  },
  {
    id: 'blue',
    name: {
      'en': 'Soft Blue',
      'zh-CN': '柔和蓝',
      'zh-TW': '柔和藍',
      'ja': 'ソフトブルー',
      'es': 'Azul suave',
      'pt': 'Azul suave',
    },
    bgClass: 'bg-blue-50',
    bgStyle: { backgroundColor: '#eff6ff' },
  },
  {
    id: 'green',
    name: {
      'en': 'Soft Green',
      'zh-CN': '柔和绿',
      'zh-TW': '柔和綠',
      'ja': 'ソフトグリーン',
      'es': 'Verde suave',
      'pt': 'Verde suave',
    },
    bgClass: 'bg-green-50',
    bgStyle: { backgroundColor: '#f0fdf4' },
  },
  {
    id: 'purple',
    name: {
      'en': 'Soft Purple',
      'zh-CN': '柔和紫',
      'zh-TW': '柔和紫',
      'ja': 'ソフトパープル',
      'es': 'Púrpura suave',
      'pt': 'Roxo suave',
    },
    bgClass: 'bg-purple-50',
    bgStyle: { backgroundColor: '#faf5ff' },
  },
  {
    id: 'gradient-blue',
    name: {
      'en': 'Blue Gradient',
      'zh-CN': '蓝色渐变',
      'zh-TW': '藍色漸變',
      'ja': '青グラデーション',
      'es': 'Degradado azul',
      'pt': 'Gradiente azul',
    },
    bgClass: 'bg-gradient-to-br from-blue-50 to-indigo-100',
    bgStyle: {
      background: 'linear-gradient(to bottom right, #eff6ff, #e0e7ff)',
    },
  },
  {
    id: 'gradient-sunset',
    name: {
      'en': 'Sunset Gradient',
      'zh-CN': '日落渐变',
      'zh-TW': '日落漸變',
      'ja': '夕焼けグラデーション',
      'es': 'Degradado atardecer',
      'pt': 'Gradiente pôr do sol',
    },
    bgClass: 'bg-gradient-to-br from-orange-50 to-pink-100',
    bgStyle: {
      background: 'linear-gradient(to bottom right, #fff7ed, #fce7f3)',
    },
  },
  {
    id: 'dots',
    name: {
      'en': 'Dots Pattern',
      'zh-CN': '圆点图案',
      'zh-TW': '圓點圖案',
      'ja': 'ドットパターン',
      'es': 'Patrón de puntos',
      'pt': 'Padrão de pontos',
    },
    bgClass: 'bg-white',
    bgStyle: {
      backgroundColor: '#ffffff',
      backgroundImage: 'radial-gradient(#e5e7eb 1px, transparent 1px)',
      backgroundSize: '20px 20px',
    },
  },
  {
    id: 'grid',
    name: {
      'en': 'Grid Pattern',
      'zh-CN': '网格图案',
      'zh-TW': '網格圖案',
      'ja': 'グリッドパターン',
      'es': 'Patrón de cuadrícula',
      'pt': 'Padrão de grade',
    },
    bgClass: 'bg-white',
    bgStyle: {
      backgroundColor: '#ffffff',
      backgroundImage: `
        linear-gradient(#e5e7eb 1px, transparent 1px),
        linear-gradient(90deg, #e5e7eb 1px, transparent 1px)
      `,
      backgroundSize: '20px 20px',
    },
  },
  {
    id: 'diagonal',
    name: {
      'en': 'Diagonal Lines',
      'zh-CN': '斜线图案',
      'zh-TW': '斜線圖案',
      'ja': '斜線パターン',
      'es': 'Líneas diagonales',
      'pt': 'Linhas diagonais',
    },
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
  },
  {
    id: 'noise',
    name: {
      'en': 'Noise Texture',
      'zh-CN': '噪点纹理',
      'zh-TW': '噪點紋理',
      'ja': 'ノイズテクスチャ',
      'es': 'Textura de ruido',
      'pt': 'Textura de ruído',
    },
    bgClass: 'bg-gray-50',
    bgStyle: {
      backgroundColor: '#f9fafb',
      backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='0.05'/%3E%3C/svg%3E")`,
    },
  },
  {
    id: 'blur-blue',
    name: {
      'en': 'Blurred Blue',
      'zh-CN': '模糊蓝色',
      'zh-TW': '模糊藍色',
      'ja': 'ぼかしブルー',
      'es': 'Azul difuminado',
      'pt': 'Azul desfocado',
    },
    bgClass: 'bg-blue-50',
    bgStyle: {
      backgroundColor: '#eff6ff',
      backgroundImage: 'radial-gradient(circle at 20% 30%, rgba(59, 130, 246, 0.15) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(99, 102, 241, 0.15) 0%, transparent 50%)',
    },
  },
  {
    id: 'blur-rainbow',
    name: {
      'en': 'Blurred Rainbow',
      'zh-CN': '模糊彩虹',
      'zh-TW': '模糊彩虹',
      'ja': 'ぼかしレインボー',
      'es': 'Arcoíris difuminado',
      'pt': 'Arco-íris desfocado',
    },
    bgClass: 'bg-white',
    bgStyle: {
      backgroundColor: '#ffffff',
      backgroundImage: `
        radial-gradient(circle at 10% 20%, rgba(239, 68, 68, 0.1) 0%, transparent 50%),
        radial-gradient(circle at 90% 30%, rgba(59, 130, 246, 0.1) 0%, transparent 50%),
        radial-gradient(circle at 50% 80%, rgba(16, 185, 129, 0.1) 0%, transparent 50%)
      `,
    },
  },
];

