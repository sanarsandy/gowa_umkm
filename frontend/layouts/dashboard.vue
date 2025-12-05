<template>
  <div class="flex h-screen bg-slate-50 font-sans text-slate-900 overflow-hidden">
    <!-- Mobile Sidebar Overlay -->
    <div v-if="sidebarOpen" class="fixed inset-0 bg-slate-900/50 z-40 lg:hidden" @click="sidebarOpen = false"></div>

    <!-- Sidebar (Desktop & Mobile) -->
    <aside 
      class="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-white border-r border-slate-200 transform transition-transform duration-300 lg:transform-none flex flex-col flex-shrink-0"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <!-- Logo -->
      <div class="h-16 flex items-center px-6 border-b border-slate-100 flex-shrink-0">
        <NuxtLink to="/dashboard" class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">G</span>
          </div>
          <span class="text-lg font-bold text-slate-800">Gowa UMKM</span>
        </NuxtLink>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
        <NuxtLink to="/dashboard" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
          Dashboard
        </NuxtLink>

        <NuxtLink to="/dashboard/whatsapp" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/whatsapp') ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path.startsWith('/dashboard/whatsapp') ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          WhatsApp
        </NuxtLink>

        <NuxtLink to="/dashboard/customers" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/customers' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/customers' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          Pelanggan
        </NuxtLink>

        <NuxtLink to="/dashboard/customers/segments" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/customers/segments' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/customers/segments' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
          Segmentasi
        </NuxtLink>

        <NuxtLink to="/dashboard/analytics" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/analytics' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/analytics' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Analytics
        </NuxtLink>

        <NuxtLink to="/dashboard/broadcasts" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/broadcasts') ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path.startsWith('/dashboard/broadcasts') ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
          </svg>
          Broadcast
        </NuxtLink>

        <NuxtLink to="/dashboard/templates" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/templates') ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path.startsWith('/dashboard/templates') ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Template
        </NuxtLink>

        <div class="px-4 mt-6 mb-2 text-xs font-semibold text-slate-400 uppercase tracking-wider">AI</div>

        <NuxtLink to="/dashboard/ai/config" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/ai/config' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/ai/config' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          AI Auto-Reply
        </NuxtLink>

        <NuxtLink to="/dashboard/ai/knowledge" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/ai/knowledge' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/ai/knowledge' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          Knowledge Base
        </NuxtLink>

        <NuxtLink to="/dashboard/ai/test" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/ai/test' ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path === '/dashboard/ai/test' ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
          Test AI
        </NuxtLink>

        <div class="px-4 mt-6 mb-2 text-xs font-semibold text-slate-400 uppercase tracking-wider">Pengaturan</div>

        <NuxtLink to="/dashboard/settings/tenant" class="flex items-center px-4 py-2.5 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/settings') ? 'bg-blue-50 text-blue-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <svg class="w-5 h-5 mr-3" :class="$route.path.startsWith('/dashboard/settings') ? 'text-blue-600' : 'text-slate-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Pengaturan
        </NuxtLink>
      </nav>

      <!-- User Profile (Bottom) -->
      <div class="border-t border-slate-100 p-4 bg-white flex-shrink-0">
        <div class="flex items-center space-x-3">
          <div class="w-9 h-9 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-sm">
            {{ userInitials }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-slate-900 truncate">{{ userName }}</p>
            <p class="text-xs text-slate-500 truncate">Pro Plan</p>
          </div>
          <button @click="logout" class="text-slate-400 hover:text-red-500 transition">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content Wrapper -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top Header -->
      <header class="bg-white border-b border-slate-200 flex-shrink-0">
        <div class="flex items-center justify-between px-6 py-3">
          <div class="flex items-center">
            <button @click="sidebarOpen = !sidebarOpen" class="lg:hidden mr-4 text-slate-500 hover:text-slate-700">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <h1 class="text-lg font-semibold text-slate-800">Dashboard</h1>
          </div>
          
          <div class="flex items-center space-x-4">
            <button class="relative p-2 text-slate-400 hover:text-blue-600 transition">
              <span class="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </button>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 overflow-y-auto p-6 scroll-smooth">
        <div class="max-w-7xl mx-auto">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const sidebarOpen = ref(false)
const authStore = useAuthStore()

const userName = computed(() => authStore.user?.full_name || 'Pengguna')
const userEmail = computed(() => authStore.user?.email || 'user@example.com')
const userInitials = computed(() => {
  const name = userName.value
  return name.split(' ').map((n: string) => n[0]).join('').substring(0, 2).toUpperCase()
})

const logout = async () => {
  console.log('[Dashboard] Logout clicked')
  
  // Clear auth and wait for it to complete
  await authStore.logout()
  
  console.log('[Dashboard] Auth cleared, redirecting...')
  
  // Force full page navigation to clear all state
  // Use external: true to ensure full page reload
  if (process.client) {
    // Use setTimeout to ensure state is fully cleared
    setTimeout(() => {
      // Force hard redirect
      window.location.href = '/login'
    }, 200)
  }
}
</script>
