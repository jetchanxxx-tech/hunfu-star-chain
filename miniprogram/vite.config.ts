import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'

// uni-app CLI 默认从 src/ 读取，此处指定为项目根目录
process.env.UNI_INPUT_DIR = process.env.UNI_INPUT_DIR || '.'

export default defineConfig({
  base: '/',
  plugins: [uni()],
  server: {
    port: 3091,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
