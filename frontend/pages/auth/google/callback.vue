<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <h2 class="text-xl font-semibold text-slate-900">Memproses Login...</h2>
      <p class="text-slate-500 mt-2">Mohon tunggu sebentar</p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'guest'
})

const route = useRoute()
const authStore = useAuthStore()

onMounted(async () => {
  const token = route.query.token as string
  const userStr = route.query.user as string

  if (token && userStr) {
    try {
      const user = JSON.parse(userStr)
      authStore.setToken(token)
      authStore.setUser(user)
      
      // Wait for state to update
      await nextTick()
      // Redirect to dashboard
      await navigateTo('/dashboard', { replace: true })
    } catch (e) {
      console.error('Failed to parse user data', e)
      await navigateTo('/login?error=invalid_data', { replace: true })
    }
  } else {
    await navigateTo('/login?error=no_token', { replace: true })
  }
})
</script>
