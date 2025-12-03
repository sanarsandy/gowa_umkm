<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <button 
          @click="router.push('/dashboard/broadcasts')" 
          class="p-2 hover:bg-gray-100 rounded-lg transition"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
        </button>
        <div>
          <h1 class="text-xl font-bold text-gray-900">{{ broadcast?.name || 'Detail Broadcast' }}</h1>
          <p class="text-gray-500 text-sm">{{ formatDate(broadcast?.created_at) }}</p>
        </div>
      </div>
      
      <div class="flex items-center gap-2">
        <span 
          :class="getStatusClass(broadcast?.status || 'draft')"
          class="px-3 py-1 rounded-full text-sm font-medium"
        >
          {{ broadcastHelper.getStatusInfo(broadcast?.status || 'draft').label }}
        </span>
        
        <!-- Cancel Button -->
        <button 
          v-if="['draft', 'scheduled', 'active'].includes(broadcast?.status || '')"
          @click="cancelBroadcast"
          :disabled="cancelling"
          class="px-4 py-2 border border-red-300 text-red-700 rounded-lg hover:bg-red-50 disabled:opacity-50 flex items-center gap-2"
        >
          <svg v-if="cancelling" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
          {{ cancelling ? 'Membatalkan...' : 'Batalkan' }}
        </button>
        
        <!-- Send Button -->
        <button 
          v-if="broadcast?.status === 'draft'"
          @click="sendBroadcast"
          :disabled="sending"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 flex items-center gap-2"
        >
          <svg v-if="sending" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ sending ? 'Mengirim...' : 'Kirim Sekarang' }}
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>

    <!-- Content -->
    <template v-else-if="broadcast">
      <!-- Stats -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="bg-white rounded-xl border border-gray-100 p-4">
          <p class="text-gray-500 text-sm">Total Penerima</p>
          <p class="text-2xl font-bold text-gray-900">{{ broadcast.total_recipients }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4">
          <p class="text-gray-500 text-sm">Terkirim</p>
          <p class="text-2xl font-bold text-blue-600">{{ broadcast.sent_count }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4">
          <p class="text-gray-500 text-sm">Diterima</p>
          <p class="text-2xl font-bold text-green-600">{{ broadcast.delivered_count }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-100 p-4">
          <p class="text-gray-500 text-sm">Gagal</p>
          <p class="text-2xl font-bold text-red-600">{{ broadcast.failed_count }}</p>
        </div>
      </div>

      <!-- Progress Bar -->
      <div v-if="broadcast.status === 'sending' || broadcast.status === 'completed'" class="bg-white rounded-xl border border-gray-100 p-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-700">Progress Pengiriman</span>
          <span class="text-sm text-gray-500">{{ broadcastHelper.formatProgress(broadcast) }}</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-3">
          <div 
            :class="broadcast.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'"
            class="h-3 rounded-full transition-all duration-300"
            :style="{ width: broadcastHelper.formatProgress(broadcast) }"
          ></div>
        </div>
      </div>

      <!-- Message Content -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h2 class="font-semibold text-gray-900 mb-4">Isi Pesan</h2>
        <div class="bg-green-50 rounded-lg p-4">
          <p class="text-gray-700 whitespace-pre-wrap">{{ broadcast.message_content }}</p>
        </div>
      </div>

      <!-- Recipients -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h2 class="font-semibold text-gray-900 mb-4">Daftar Penerima ({{ recipients.length }})</h2>
        
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b">
                <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">Pelanggan</th>
                <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">Status</th>
                <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">Waktu Kirim</th>
                <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">Keterangan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="recipient in recipients" :key="recipient.id" class="border-b last:border-0 hover:bg-gray-50">
                <td class="py-3 px-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                      <span class="text-white font-medium text-xs">{{ getInitials(recipient.customer_name) }}</span>
                    </div>
                    <div>
                      <p class="font-medium text-gray-900 text-sm">{{ recipient.customer_name }}</p>
                      <p class="text-gray-500 text-xs">{{ formatPhone(recipient.customer_jid) }}</p>
                    </div>
                  </div>
                </td>
                <td class="py-3 px-4">
                  <span 
                    :class="getRecipientStatusClass(recipient.status)"
                    class="px-2 py-1 rounded-full text-xs font-medium"
                  >
                    {{ broadcastHelper.getRecipientStatusInfo(recipient.status).label }}
                  </span>
                </td>
                <td class="py-3 px-4 text-sm text-gray-500">
                  {{ recipient.sent_at ? formatTime(recipient.sent_at) : '-' }}
                </td>
                <td class="py-3 px-4 text-sm text-gray-500">
                  <span v-if="recipient.error_message" class="text-red-600">{{ recipient.error_message }}</span>
                  <span v-else-if="recipient.message_id" class="text-green-600">ID: {{ recipient.message_id.substring(0, 10) }}...</span>
                  <span v-else>-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <p v-if="recipients.length === 0" class="text-center text-gray-500 py-8">
          Tidak ada penerima
        </p>
      </div>
    </template>

    <!-- Error -->
    <div v-else class="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
      <p class="text-red-800">Broadcast tidak ditemukan</p>
      <button 
        @click="router.push('/dashboard/broadcasts')"
        class="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
      >
        Kembali ke Daftar
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Broadcast, BroadcastRecipient } from '~/composables/useBroadcasts'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const broadcastHelper = useBroadcasts()

// State
const loading = ref(true)
const sending = ref(false)
const cancelling = ref(false)
const broadcast = ref<Broadcast | null>(null)
const recipients = ref<BroadcastRecipient[]>([])

// Load broadcast
const loadBroadcast = async () => {
  const id = route.params.id as string
  loading.value = true
  
  try {
    const response = await broadcastHelper.getBroadcast(id)
    broadcast.value = response.broadcast
    recipients.value = response.recipients
  } catch (error) {
    console.error('Failed to load broadcast:', error)
    broadcast.value = null
  } finally {
    loading.value = false
  }
}

// Send broadcast
const sendBroadcast = async () => {
  if (!broadcast.value) return
  
  sending.value = true
  try {
    await broadcastHelper.sendBroadcast(broadcast.value.id)
    // Reload to get updated status
    await loadBroadcast()
  } catch (error) {
    console.error('Failed to send broadcast:', error)
  } finally {
    sending.value = false
  }
}

// Cancel broadcast
const cancelBroadcast = async () => {
  if (!broadcast.value) return
  
  const confirmed = confirm(
    broadcast.value.is_recurring 
      ? 'Apakah Anda yakin ingin membatalkan broadcast berulang ini? Semua jadwal berikutnya akan dibatalkan.'
      : 'Apakah Anda yakin ingin membatalkan broadcast ini?'
  )
  
  if (!confirmed) return
  
  cancelling.value = true
  try {
    await broadcastHelper.cancelBroadcast(broadcast.value.id)
    // Reload to get updated status
    await loadBroadcast()
  } catch (error) {
    console.error('Failed to cancel broadcast:', error)
    alert('Gagal membatalkan broadcast')
  } finally {
    cancelling.value = false
  }
}

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

const getRecipientStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-gray-100 text-gray-700',
    sent: 'bg-blue-100 text-blue-700',
    delivered: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
  }
  return classes[status] || classes.pending
}

// Format helpers
const formatDate = (date: string | undefined) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatTime = (date: string) => {
  return new Date(date).toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const formatPhone = (jid: string) => {
  if (!jid) return ''
  const phone = jid.split('@')[0].split(':')[0]
  if (phone.startsWith('62')) {
    return '+' + phone
  }
  return phone
}

const getInitials = (name: string | null): string => {
  if (!name) return '?'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

// Auto-refresh when sending
let refreshInterval: NodeJS.Timeout | null = null

watch(() => broadcast.value?.status, (status) => {
  if (status === 'sending') {
    // Refresh every 2 seconds while sending
    refreshInterval = setInterval(loadBroadcast, 2000)
  } else if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})

// Initial load
onMounted(() => {
  loadBroadcast()
})
</script>




