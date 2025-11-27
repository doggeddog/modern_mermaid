import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    host: '0.0.0.0', // 监听所有网络接口，允许外部访问
    port: 5173, // 默认端口
    strictPort: false, // 端口被占用时自动尝试下一个
  },
})
