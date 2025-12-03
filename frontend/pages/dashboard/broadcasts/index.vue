<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Broadcast Pesan</h1>
        <p class="text-gray-500 mt-1">Kirim pesan ke banyak pelanggan sekaligus</p>
      </div>
      <button 
        @click="router.push('/dashboard/broadcasts/create')"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Buat Broadcast
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <p class="text-gray-500 text-sm">Total Broadcast</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total_broadcasts }}</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <p class="text-gray-500 text-sm">Pesan Terkirim</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.total_messages_sent }}</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <p class="text-gray-500 text-sm">Pesan Diterima</p>
        <p class="text-2xl font-bold text-green-600">{{ stats.total_delivered }}</p>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <p class="text-gray-500 text-sm">Gagal</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.total_failed }}</p>
      </div>
    </div>

    <!-- Status Filter -->
    <div class="flex gap-2 overflow-x-auto pb-2">
      <button
        v-for="status in statuses"
        :key="status.value"
        @click="selectedStatus = status.value"
        :class="selectedStatus === status.value ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
        class="px-4 py-2 rounded-lg border whitespace-nowrap transition"
      >
        {{ status.label }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>

    <!-- Broadcasts List -->
    <div v-else-if="broadcasts.length > 0" class="space-y-4">
      <div 
        v-for="broadcast in broadcasts" 
        :key="broadcast.id"
        class="bg-white rounded-xl border border-gray-100 p-5 hover:shadow-lg transition cursor-pointer"
        @click="router.push(`/dashboard/broadcasts/${broadcast.id}`)"
      >
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="font-semibold text-gray-900">{{ broadcast.name }}</h3>
            <p class="text-gray-500 text-sm">{{ formatDate(broadcast.created_at) }}</p>
          </div>
          <span 
            :class="getStatusClass(broadcast.status)"
            class="px-3 py-1 rounded-full text-xs font-medium"
          >
            {{ broadcastHelper.getStatusInfo(broadcast.status).label }}
          </span>
        </div>

        <p class="text-gray-600 text-sm line-clamp-2 mb-4">{{ broadcast.message_content }}</p>

        <!-- Scheduled Info -->
        <div v-if="broadcast.scheduled_at" class="mb-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
          <div class="flex items-start gap-2">
            <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
            <div class="flex-1">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-blue-900">
                  {{ broadcast.status === 'scheduled' ? 'Dijadwalkan' : 'Dijalankan' }}
                </p>
                <span v-if="broadcast.is_recurring" class="px-2 py-0.5 bg-purple-100 text-purple-700 rounded text-xs font-medium">
                  🔄 Berulang
                </span>
              </div>
              <p class="text-sm text-blue-700">
                {{ formatScheduledDate(broadcast.scheduled_at) }}
              </p>
              
              <!-- Countdown for scheduled broadcasts -->
              <div v-if="broadcast.status === 'scheduled'" class="mt-2">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  <span class="text-xs font-medium text-blue-800">
                    {{ getCountdown(broadcast.scheduled_at) }}
                  </span>
                </div>
              </div>
              
              <!-- Recurring info -->
              <div v-if="broadcast.is_recurring" class="mt-2 text-xs text-purple-700">
                <p>{{ getRecurringInfo(broadcast) }}</p>
                <p v-if="broadcast.status === 'active' && broadcast.last_executed_at" class="mt-1 text-purple-600">
                  Terakhir: {{ formatDate(broadcast.last_executed_at) }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Progress -->
        <div class="flex items-center gap-4 text-sm">
          <div class="flex items-center gap-1 text-gray-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
            </svg>
            {{ broadcast.total_recipients }} penerima
          </div>
          <div class="flex items-center gap-1 text-blue-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            {{ broadcast.sent_count }} terkirim
          </div>
          <div v-if="broadcast.failed_count > 0" class="flex items-center gap-1 text-red-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            {{ broadcast.failed_count }} gagal
          </div>
        </div>

        <!-- Progress Bar -->
        <div v-if="broadcast.status === 'sending' || broadcast.status === 'completed'" class="mt-3">
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div 
              :class="broadcast.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'"
              class="h-2 rounded-full transition-all duration-300"
              :style="{ width: broadcastHelper.formatProgress(broadcast) }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white rounded-xl border border-gray-100 p-12 text-center">
      <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Belum ada broadcast</h3>
      <p class="text-gray-500 mb-4">Kirim pesan ke banyak pelanggan sekaligus</p>
      <button 
        @click="router.push('/dashboard/broadcasts/create')"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        Buat Broadcast Pertama
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Broadcast, BroadcastStats } from '~/composables/useBroadcasts'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const router = useRouter()
const broadcastHelper = useBroadcasts()

// State
const loading = ref(true)
const broadcasts = ref<Broadcast[]>([])
const stats = ref<BroadcastStats>({
  total_broadcasts: 0,
  total_messages_sent: 0,
  total_delivered: 0,
  total_failed: 0
})
const selectedStatus = ref('all')

const statuses = [
  { value: 'all', label: 'Semua' },
  { value: 'draft', label: 'Draft' },
  { value: 'scheduled', label: 'Terjadwal' },
  { value: 'active', label: 'Aktif (Berulang)' },
  { value: 'sending', label: 'Mengirim' },
  { value: 'completed', label: 'Selesai' },
]

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const [broadcastsRes, statsRes] = await Promise.all([
      broadcastHelper.getBroadcasts(selectedStatus.value),
      broadcastHelper.getStats()
    ])
    broadcasts.value = broadcastsRes.broadcasts
    stats.value = statsRes
  } catch (error) {
    console.error('Failed to load broadcasts:', error)
  } finally {
    loading.value = false
  }
}

// Watch status change
watch(selectedStatus, () => {
  loadData()
})

// Initial load
onMounted(() => {
  loadData()
})

// Status class helper
const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-gray-100 text-gray-700',
    scheduled: 'bg-blue-100 text-blue-700',
    active: 'bg-purple-100 text-purple-700',
    sending: 'bg-yellow-100 text-yellow-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-red-100 text-red-700',
  }
  return classes[status] || classes.draft
}

// Format date
const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Format scheduled date
const formatScheduledDate = (date: string) => {
  return new Date(date).toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Get countdown
const getCountdown = (scheduledAt: string) => {
  const now = new Date()
  const scheduled = new Date(scheduledAt)
  const diff = scheduled.getTime() - now.getTime()
  
  if (diff <= 0) return 'Segera dijalankan...'
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (days > 0) return `${days} hari ${hours} jam lagi`
  if (hours > 0) return `${hours} jam ${minutes} menit lagi`
  return `${minutes} menit lagi`
}

// Get recurring info
const getRecurringInfo = (broadcast: Broadcast) => {
  if (!broadcast.is_recurring) return ''
  
  const type = broadcast.recurrence_type
  const interval = broadcast.recurrence_interval || 1
  
  let info = ''
  if (type === 'hourly') {
    info = `Setiap ${interval} jam`
  } else if (type === 'daily') {
    info = `Setiap ${interval} hari`
  } else if (type === 'weekly') {
    const days = broadcast.recurrence_days || []
    const dayNames: Record<string, string> = {
      sunday: 'Minggu', monday: 'Senin', tuesday: 'Selasa',
      wednesday: 'Rabu', thursday: 'Kamis', friday: 'Jumat', saturday: 'Sabtu'
    }
    const dayLabels = days.map((d: string) => dayNames[d] || d).join(', ')
    info = `Setiap ${dayLabels}`
  }
  
  // Add end condition
  if (broadcast.recurrence_end_date) {
    info += ` sampai ${new Date(broadcast.recurrence_end_date).toLocaleDateString('id-ID')}`
  } else if (broadcast.recurrence_count) {
    const remaining = broadcast.recurrence_count - (broadcast.execution_count || 0)
    info += ` (${remaining} kali lagi)`
  }
  
  return info
}

// Auto-refresh countdown every minute
let countdownInterval: NodeJS.Timeout | null = null

onMounted(() => {
  loadData()
  
  // Refresh countdown every minute
  countdownInterval = setInterval(() => {
    // Force re-render by updating a reactive property
    broadcasts.value = [...broadcasts.value]
  }, 60000)
})

onUnmounted(() => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
  }
})
</script>




