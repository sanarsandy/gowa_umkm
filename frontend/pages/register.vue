<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 font-sans relative overflow-hidden">
    <!-- Background Gradients -->
    <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] rounded-full bg-blue-100/50 blur-3xl"></div>
      <div class="absolute bottom-[10%] right-[10%] w-[40%] h-[40%] rounded-full bg-sky-100/50 blur-3xl"></div>
    </div>

    <div class="max-w-md w-full space-y-8 relative z-10 bg-white p-8 sm:p-10 rounded-3xl shadow-xl shadow-slate-200/50 border border-slate-100">
      <div class="text-center">
        <NuxtLink to="/" class="inline-flex mx-auto w-12 h-12 bg-gradient-to-br from-blue-600 to-sky-500 rounded-xl items-center justify-center mb-6 shadow-lg shadow-blue-500/20">
          <span class="text-white font-bold text-xl">G</span>
        </NuxtLink>
        <h2 class="text-3xl font-bold text-slate-900 tracking-tight">
          Buat akun baru
        </h2>
        <p class="mt-2 text-sm text-slate-500">
          Sudah punya akun?
          <NuxtLink to="/login" class="font-medium text-blue-600 hover:text-blue-500 transition">
            Masuk sekarang
          </NuxtLink>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister" novalidate>
        <div class="space-y-5">
          <div>
            <label for="full-name" class="block text-sm font-medium text-slate-700 mb-1">Nama Lengkap</label>
            <input
              id="full-name"
              v-model="fullName"
              name="name"
              type="text"
              required
              class="appearance-none block w-full px-4 py-3 border border-slate-200 placeholder-slate-400 text-slate-900 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition bg-slate-50 focus:bg-white"
              placeholder="Nama Lengkap Anda"
            />
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-slate-700 mb-1">Email</label>
            <input
              id="email"
              v-model="email"
              name="email"
              type="email"
              required
              class="appearance-none block w-full px-4 py-3 border border-slate-200 placeholder-slate-400 text-slate-900 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition bg-slate-50 focus:bg-white"
              placeholder="nama@perusahaan.com"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-slate-700 mb-1">Password</label>
            <input
              id="password"
              v-model="password"
              name="password"
              type="password"
              required
              minlength="6"
              class="appearance-none block w-full px-4 py-3 border border-slate-200 placeholder-slate-400 text-slate-900 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition bg-slate-50 focus:bg-white"
              placeholder="Minimal 6 karakter"
            />
            <p class="mt-1 text-xs text-slate-500">Password minimal 6 karakter</p>
          </div>
        </div>

        <div v-if="error" class="text-red-600 text-sm text-center bg-red-50 p-3 rounded-xl border border-red-100">
          <p class="font-medium">{{ error }}</p>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            @click="handleRegister"
            class="group relative w-full flex justify-center py-3.5 px-4 border border-transparent text-sm font-bold rounded-xl text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition shadow-lg shadow-blue-600/20 hover:shadow-blue-600/30 transform hover:-translate-y-0.5"
          >
            <span v-if="loading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Mendaftar...
            </span>
            <span v-else>Daftar</span>
          </button>
        </div>

        <div class="mt-6">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-slate-100" />
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-4 bg-white text-slate-400">Atau daftar dengan</span>
            </div>
          </div>

          <div class="mt-6">
            <button
              type="button"
              @click="handleGoogleLogin"
              class="w-full inline-flex justify-center items-center py-3.5 px-4 border border-slate-200 rounded-xl shadow-sm bg-white text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900 transition hover:border-slate-300"
            >
              <svg class="w-5 h-5 mr-2" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M12.545,10.239v3.821h5.445c-0.712,2.315-2.647,3.972-5.445,3.972c-3.332,0-6.033-2.701-6.033-6.032s2.701-6.032,6.033-6.032c1.498,0,2.866,0.549,3.921,1.453l2.814-2.814C17.503,2.988,15.139,2,12.545,2C7.021,2,2.543,6.477,2.543,12s4.478,10,10.002,10c8.396,0,10.249-7.85,9.426-11.748L12.545,10.239z"
                />
              </svg>
              Google
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'guest'
})

const { fetch } = useApi()
const authStore = useAuthStore()

const fullName = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleRegister = async (e?: Event) => {
  // Prevent default form submission
  if (e) {
    e.preventDefault()
  }
  
  // Validasi frontend
  if (!fullName.value || !email.value || !password.value) {
    error.value = 'Semua field harus diisi'
    return
  }

  // Basic email validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value)) {
    error.value = 'Format email tidak valid'
    return
  }

  if (password.value.length < 6) {
    error.value = 'Password minimal 6 karakter'
    return
  }

  loading.value = true
  error.value = ''

  console.log('=== REGISTER START ===')
  console.log('Full Name:', fullName.value)
  console.log('Email:', email.value)
  console.log('Loading:', loading.value)

  try {
    console.log('Attempting register for:', email.value)
    
    const response = await fetch<{ token: string; user: any }>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({
        full_name: fullName.value,
        email: email.value,
        password: password.value
      })
    })

    console.log('Register successful:', response)

    if (!response || !response.token || !response.user) {
      throw new Error('Invalid response from server')
    }

    console.log('Setting token and user...')
    authStore.setToken(response.token)
    authStore.setUser(response.user)
    
    console.log('Token set:', !!authStore.token)
    console.log('User set:', !!authStore.user)
    
    // Wait for state to update
    await nextTick()
    
    console.log('Navigating to dashboard...')
    
    // Try navigateTo first, with fallback
    try {
      await navigateTo('/dashboard', { replace: true })
      console.log('Navigation successful')
    } catch (navError) {
      console.error('Navigation error, using fallback:', navError)
      // Fallback to window.location
      if (process.client) {
        window.location.href = '/dashboard'
      }
    }
  } catch (err: any) {
    console.error('=== REGISTER ERROR ===')
    console.error('Register error:', err)
    console.error('Error type:', typeof err)
    console.error('Error keys:', Object.keys(err || {}))
    console.error('Error message:', err?.message)
    console.error('Error data:', err?.data)
    
    // Extract error message from various possible formats
    let errorMessage = 'Pendaftaran gagal. Silakan coba lagi.'
    
    // Check multiple error formats
    if (err?.data?.error) {
      errorMessage = err.data.error
    } else if (err?.response?._data?.error) {
      errorMessage = err.response._data.error
    } else if (err?.message) {
      errorMessage = err.message
    } else if (typeof err === 'string') {
      errorMessage = err
    }
    
    // Check for network errors
    const errorMsg = err?.message || err?.toString() || ''
    if (errorMsg.includes('fetch') || errorMsg.includes('network') || errorMsg.includes('Failed to fetch') || errorMsg.includes('ERR_') || errorMsg.includes('ECONNREFUSED')) {
      errorMessage = 'Tidak dapat terhubung ke server. Pastikan backend sedang berjalan di http://localhost:8080'
    }
    
    // Check for CORS errors
    if (errorMsg.includes('CORS') || errorMsg.includes('cors')) {
      errorMessage = 'Error CORS. Pastikan backend mengizinkan origin frontend.'
    }
    
    error.value = errorMessage
    console.log('Displayed error:', errorMessage)
  } finally {
    loading.value = false
    console.log('=== REGISTER END ===')
  }
}

const handleGoogleLogin = async () => {
  try {
    const response = await fetch<{ auth_url: string }>('/api/auth/google')
    window.location.href = response.auth_url
  } catch (err: any) {
    error.value = err.data?.error || 'Gagal memulai login Google.'
  }
}
</script>
