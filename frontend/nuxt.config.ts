export default defineNuxtConfig({
  devtools: { enabled: process.env.NODE_ENV !== 'production' },
  css: ['~/assets/css/main.css'],
  modules: ['@pinia/nuxt', 'nuxt-icon'],
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  runtimeConfig: {
    apiInternal: process.env.NUXT_API_INTERNAL || 'http://localhost:8080',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  },
  nitro: {
    preset: 'node-server'
  },
  vite: {
    server: {
      host: '0.0.0.0',
      allowedHosts: process.env.VITE_ALLOWED_HOSTS 
        ? process.env.VITE_ALLOWED_HOSTS.split(',')
        : ['localhost', '127.0.0.1']
    }
  }
})

