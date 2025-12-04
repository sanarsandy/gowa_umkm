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
  app: {
    head: {
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' }
      ],
      // Note: Cloudflare Insights script is automatically injected by Cloudflare
      // If blocked by tracking prevention, it's normal and not critical
      script: [
        // Add any custom scripts here if needed
      ]
    }
  },
  nitro: {
    preset: 'node-server',
    compressPublicAssets: true,
    minify: true
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

