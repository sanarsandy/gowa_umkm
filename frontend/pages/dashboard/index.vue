<template>
  <div class="space-y-6">
    <!-- Welcome Section -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">
          Selamat Datang, <span class="text-blue-600">{{ userName }}</span> 👋
        </h1>
        <p class="text-slate-500 mt-1">Ringkasan performa bisnis Anda hari ini.</p>
      </div>
      <div class="flex items-center space-x-3">
        <button 
          @click="refreshData"
          :disabled="loading"
          class="px-4 py-2 bg-white border border-slate-200 text-slate-600 font-medium rounded-lg hover:bg-slate-50 transition shadow-sm disabled:opacity-50"
        >
          <span v-if="loading">Memuat...</span>
          <span v-else>🔄 Refresh</span>
        </button>
        <NuxtLink to="/dashboard/whatsapp" class="px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition shadow-sm flex items-center">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Hubungkan Device
        </NuxtLink>
      </div>
    </div>

    <!-- Connection Status Banner -->
    <div v-if="!stats.is_connected" class="bg-amber-50 border border-amber-200 rounded-lg p-4 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 bg-amber-100 rounded-full flex items-center justify-center">
          <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
        </div>
        <div>
          <p class="font-medium text-amber-800">WhatsApp belum terhubung</p>
          <p class="text-sm text-amber-600">Hubungkan WhatsApp untuk mulai menerima pesan</p>
        </div>
      </div>
      <NuxtLink to="/dashboard/whatsapp" class="px-4 py-2 bg-amber-600 text-white font-medium rounded-lg hover:bg-amber-700 transition">
        Hubungkan Sekarang
      </NuxtLink>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Card 1: Total Messages -->
      <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center text-blue-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"/></svg>
          </div>
          <span v-if="stats.messages_today > 0" class="text-emerald-600 text-xs font-bold bg-emerald-50 px-2 py-1 rounded-full">
            +{{ stats.messages_today }} hari ini
          </span>
        </div>
        <p class="text-sm font-medium text-slate-500">Total Pesan</p>
        <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ formatNumber(stats.total_messages) }}</h3>
      </div>

      <!-- Card 2: Hot Leads -->
      <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-red-50 rounded-lg flex items-center justify-center text-red-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z"/></svg>
          </div>
          <span v-if="stats.hot_leads > 0" class="text-red-600 text-xs font-bold bg-red-50 px-2 py-1 rounded-full animate-pulse">
            🔥 Hot!
          </span>
        </div>
        <p class="text-sm font-medium text-slate-500">Hot Leads</p>
        <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ stats.hot_leads }}</h3>
      </div>

      <!-- Card 3: Total Customers -->
      <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-purple-50 rounded-lg flex items-center justify-center text-purple-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
          </div>
          <span v-if="stats.new_customers_today > 0" class="text-purple-600 text-xs font-bold bg-purple-50 px-2 py-1 rounded-full">
            +{{ stats.new_customers_today }} baru
          </span>
        </div>
        <p class="text-sm font-medium text-slate-500">Pelanggan</p>
        <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ stats.total_customers }}</h3>
      </div>

      <!-- Card 4: Positive Sentiment -->
      <div class="bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-emerald-50 rounded-lg flex items-center justify-center text-emerald-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <span :class="[
            'text-xs font-bold px-2 py-1 rounded-full',
            stats.positive_sentiment >= 70 ? 'text-emerald-600 bg-emerald-50' : 
            stats.positive_sentiment >= 40 ? 'text-amber-600 bg-amber-50' : 'text-red-600 bg-red-50'
          ]">
            {{ getSentimentLabel(stats.positive_sentiment) }}
          </span>
        </div>
        <p class="text-sm font-medium text-slate-500">Sentimen Positif</p>
        <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ stats.positive_sentiment.toFixed(0) }}%</h3>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Quick Actions (Grid) -->
      <div class="lg:col-span-2">
        <h2 class="text-lg font-bold text-slate-800 mb-4">Aksi Cepat</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <NuxtLink to="/dashboard/whatsapp" class="group p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:border-green-300 hover:shadow-md transition-all duration-200 flex items-start space-x-4">
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center text-green-600">
              <svg class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
              </svg>
            </div>
            <div>
              <h3 class="font-bold text-slate-800 group-hover:text-green-600 transition-colors">WhatsApp</h3>
              <p class="text-sm text-slate-500 mt-1">
                {{ stats.is_connected ? 'Terhubung ✓' : 'Hubungkan sekarang' }}
              </p>
            </div>
          </NuxtLink>

          <NuxtLink to="/dashboard/customers" class="group p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:border-blue-300 hover:shadow-md transition-all duration-200 flex items-start space-x-4">
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center text-blue-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/></svg>
            </div>
            <div>
              <h3 class="font-bold text-slate-800 group-hover:text-blue-600 transition-colors">Pelanggan</h3>
              <p class="text-sm text-slate-500 mt-1">{{ stats.total_customers }} pelanggan</p>
            </div>
          </NuxtLink>

          <NuxtLink to="/dashboard/settings/tenant" class="group p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:border-purple-300 hover:shadow-md transition-all duration-200 flex items-start space-x-4">
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center text-purple-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
            </div>
            <div>
              <h3 class="font-bold text-slate-800 group-hover:text-purple-600 transition-colors">Pengaturan Bisnis</h3>
              <p class="text-sm text-slate-500 mt-1">Kelola profil bisnis</p>
            </div>
          </NuxtLink>
          
          <div class="group p-6 bg-slate-50 rounded-xl border border-slate-200 border-dashed flex items-center justify-center text-center cursor-not-allowed opacity-60">
            <div>
              <div class="w-8 h-8 mx-auto bg-slate-200 rounded-full flex items-center justify-center text-slate-500 mb-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
              </div>
              <span class="text-sm font-medium text-slate-500">AI Assistant</span>
              <p class="text-xs text-slate-400 mt-1">Segera hadir</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Messages (Chat Style) -->
      <div class="bg-white rounded-xl border border-slate-200 shadow-sm flex flex-col overflow-hidden">
        <div class="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50">
          <h3 class="font-bold text-slate-800">Pesan Terbaru</h3>
          <NuxtLink to="/dashboard/messages" class="text-xs font-bold text-blue-600 hover:text-blue-700">Lihat Semua</NuxtLink>
        </div>
        <div class="flex-1 overflow-y-auto max-h-[400px]">
          <!-- Loading State -->
          <div v-if="loadingMessages" class="p-8 text-center">
            <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto"></div>
            <p class="text-sm text-slate-500 mt-2">Memuat pesan...</p>
          </div>
          
          <!-- Empty State -->
          <div v-else-if="recentMessages.length === 0" class="p-8 text-center">
            <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
            </div>
            <p class="text-slate-500 font-medium">Belum ada pesan</p>
            <p class="text-sm text-slate-400 mt-1">Pesan dari pelanggan akan muncul di sini</p>
          </div>
          
          <!-- Messages List -->
          <div v-else v-for="msg in recentMessages" :key="msg.id" class="p-4 hover:bg-slate-50 transition-colors border-b border-slate-50 last:border-0 cursor-pointer group">
            <div class="flex items-start space-x-3">
              <div class="relative">
                <div class="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-bold text-xs">
                  {{ getInitials(msg.customer_name) }}
                </div>
                <div class="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-baseline mb-1">
                  <h4 class="text-sm font-bold text-slate-800 group-hover:text-blue-600 transition-colors truncate">
                    {{ msg.customer_name }}
                  </h4>
                  <span class="text-xs text-slate-400 ml-2 flex-shrink-0">
                    {{ dashboard.formatRelativeTime(msg.timestamp) }}
                  </span>
                </div>
                <p class="text-sm text-slate-500 truncate">{{ msg.message_text || '(media)' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardStats, RecentMessage } from '~/composables/useDashboard'
import { useRealtimeUpdates, type NewMessageData } from '~/composables/useRealtimeUpdates'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const authStore = useAuthStore()
const dashboard = useDashboard()
const userName = computed(() => authStore.user?.full_name || 'Pengguna')

// State
const loading = ref(false)
const loadingMessages = ref(false)
const stats = ref<DashboardStats>({
  total_messages: 0,
  total_customers: 0,
  hot_leads: 0,
  positive_sentiment: 0,
  messages_today: 0,
  new_customers_today: 0,
  is_connected: false
})
const recentMessages = ref<RecentMessage[]>([])

// Real-time updates via WebSocket
const { isConnected: wsConnected } = useRealtimeUpdates({
  onNewMessage: (data: NewMessageData) => {
    console.log('[Dashboard] New message received:', data)
    
    // Add new message to the top of the list
    recentMessages.value.unshift({
      id: data.message_id,
      sender_jid: data.sender_jid,
      sender_name: data.sender_jid.split('@')[0], // Extract phone number
      message_text: data.message_text,
      message_type: data.message_type,
      timestamp: new Date(data.timestamp * 1000).toISOString(),
      is_from_me: data.is_from_me
    })
    
    // Keep only the last 10 messages
    if (recentMessages.value.length > 10) {
      recentMessages.value.pop()
    }
    
    // Update stats
    stats.value.total_messages++
    stats.value.messages_today++
  },
  onNewCustomer: () => {
    console.log('[Dashboard] New customer')
    stats.value.total_customers++
    stats.value.new_customers_today++
  },
  onConnectionChange: (connected) => {
    console.log('[Dashboard] WebSocket connection:', connected ? 'connected' : 'disconnected')
  }
})

// Load data on mount
onMounted(async () => {
  await refreshData()
})

// Refresh all data
const refreshData = async (silent = false) => {
  if (!silent) loading.value = true
  loadingMessages.value = !silent
  
  try {
    // Fetch stats and messages in parallel
    const [statsData, messagesData] = await Promise.all([
      dashboard.getStats(),
      dashboard.getRecentMessages()
    ])
    
    stats.value = statsData
    recentMessages.value = messagesData
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  } finally {
    loading.value = false
    loadingMessages.value = false
  }
}

// Helper functions
const formatNumber = (num: number): string => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

const getSentimentLabel = (percentage: number): string => {
  if (percentage >= 70) return 'Bagus 👍'
  if (percentage >= 40) return 'Netral'
  return 'Perlu Perhatian'
}

const getInitials = (name: string): string => {
  if (!name) return '?'
  const parts = name.split(/[@\s]+/)
  if (parts.length === 1) {
    // If it's a phone number, return last 2 digits
    if (/^\d+$/.test(parts[0])) {
      return parts[0].slice(-2)
    }
    return parts[0].substring(0, 2).toUpperCase()
  }
  return (parts[0][0] + (parts[1]?.[0] || '')).toUpperCase()
}
</script>
