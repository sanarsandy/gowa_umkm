<template>
  <div>
    <!-- Connection Status Card -->
    <div class="bg-white rounded-xl shadow-sm p-8 border border-gray-100 mb-8">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Status WhatsApp</h2>
          <p class="text-gray-600 mt-1">Kelola koneksi WhatsApp bisnis Anda</p>
        </div>
        <div class="flex items-center space-x-3">
          <div class="flex items-center space-x-2">
            <div :class="['w-3 h-3 rounded-full', isConnected ? 'bg-green-500 animate-pulse' : 'bg-gray-300']"></div>
            <span :class="['font-medium', isConnected ? 'text-green-600' : 'text-gray-500']">
              {{ isConnected ? 'Terhubung' : 'Tidak Terhubung' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Connected State -->
      <div v-if="isConnected" class="bg-green-50 border border-green-200 rounded-lg p-6">
        <div class="flex items-start space-x-4">
          <div class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center flex-shrink-0">
            <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
            </svg>
          </div>
          <div class="flex-1">
            <h3 class="font-semibold text-gray-900 mb-1">WhatsApp Terhubung</h3>
            <p class="text-sm text-gray-600 mb-3">Nomor: {{ connectedNumber }}</p>
            <p class="text-sm text-gray-600">Terakhir terhubung: {{ lastConnected }}</p>
          </div>
          <button
            @click="handleDisconnect"
            :disabled="loading"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition disabled:opacity-50"
          >
            {{ loading ? 'Memutuskan...' : 'Putuskan Koneksi' }}
          </button>
        </div>
      </div>

      <!-- Disconnected State -->
      <div v-else class="text-center">
        <div v-if="!showQR" class="py-8">
          <div class="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-10 h-10 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
            </svg>
          </div>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">WhatsApp Belum Terhubung</h3>
          <p class="text-gray-600 mb-6">Hubungkan WhatsApp Anda untuk mulai menggunakan fitur AI</p>
          <button
            @click="handleConnect"
            :disabled="loading"
            class="px-6 py-3 bg-gradient-to-r from-green-600 to-emerald-600 text-white rounded-lg font-semibold hover:shadow-lg transition transform hover:scale-105 disabled:opacity-50"
          >
            {{ loading ? 'Menghubungkan...' : 'Hubungkan WhatsApp' }}
          </button>
        </div>

        <!-- QR Code Display -->
        <div v-else class="py-8">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">Scan QR Code</h3>
          <p class="text-gray-600 mb-6">Buka WhatsApp di ponsel Anda dan scan QR code di bawah ini</p>
          
          <div class="max-w-sm mx-auto">
            <div v-if="qrCode" class="bg-white p-6 rounded-xl border-2 border-gray-200 mb-4">
              <img :src="qrCode" alt="QR Code" class="w-full h-auto" />
            </div>
            <div v-else class="bg-gray-100 p-12 rounded-xl border-2 border-gray-200 mb-4">
              <div class="animate-spin w-12 h-12 border-4 border-indigo-600 border-t-transparent rounded-full mx-auto"></div>
              <p class="text-gray-600 mt-4">Generating QR code...</p>
            </div>

            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 text-left">
              <h4 class="font-semibold text-blue-900 mb-2">Cara Scan:</h4>
              <ol class="text-sm text-blue-800 space-y-1 list-decimal list-inside">
                <li>Buka WhatsApp di ponsel Anda</li>
                <li>Tap Menu (⋮) atau Settings</li>
                <li>Pilih "Linked Devices"</li>
                <li>Tap "Link a Device"</li>
                <li>Scan QR code ini</li>
              </ol>
            </div>

            <button
              @click="cancelQR"
              class="mt-4 px-4 py-2 text-gray-600 hover:text-gray-900 transition"
            >
              Batal
            </button>
          </div>
        </div>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
        <p class="text-red-800 mb-2">{{ error }}</p>
        <NuxtLink 
          v-if="error.includes('Tenant tidak ditemukan')" 
          to="/dashboard/settings/tenant" 
          class="text-sm text-red-600 hover:text-red-800 underline font-medium"
        >
          Buka Halaman Pengaturan →
        </NuxtLink>
      </div>
    </div>

    <!-- Features Info -->
    <div class="grid md:grid-cols-3 gap-6">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <h3 class="font-semibold text-gray-900 mb-2">Auto-Reply 24/7</h3>
        <p class="text-sm text-gray-600">AI akan membalas pesan pelanggan secara otomatis kapan saja</p>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h3 class="font-semibold text-gray-900 mb-2">Deteksi Lead Otomatis</h3>
        <p class="text-sm text-gray-600">Sistem mendeteksi pelanggan potensial dari percakapan</p>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="w-12 h-12 bg-slate-100 rounded-lg flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
          </svg>
        </div>
        <h3 class="font-semibold text-gray-900 mb-2">Analisis Sentimen</h3>
        <p class="text-sm text-gray-600">Pahami mood pelanggan untuk layanan yang lebih baik</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const whatsapp = useWhatsApp()
const { fetch } = useApi()

const isConnected = ref(false)
const showQR = ref(false)
const qrCode = ref('')
const loading = ref(false)
const error = ref('')
const connectedNumber = ref('')
const lastConnected = ref('')
let qrStreamClose: (() => void) | null = null

// Load initial status
onMounted(async () => {
  await checkStatus()
  
  // Poll status every 10 seconds if not connected
  const statusInterval = setInterval(async () => {
    if (!isConnected.value && !showQR.value) {
      await checkStatus()
    }
  }, 10000)
  
  onUnmounted(() => {
    clearInterval(statusInterval)
  })
})

// Check WhatsApp connection status
const checkStatus = async () => {
  try {
    const status = await whatsapp.getStatus()
    isConnected.value = status.is_connected
    connectedNumber.value = status.jid || ''
    
    // Format last connected time if available
    if (status.is_connected) {
      lastConnected.value = new Date().toLocaleString('id-ID', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        timeZoneName: 'short'
      })
    }
  } catch (err: any) {
    console.error('Failed to get status:', err)
    // Don't show error for initial load
  }
}

const handleConnect = async () => {
  console.log('[WhatsApp Page] handleConnect called')
  loading.value = true
  error.value = ''
  showQR.value = false
  qrCode.value = ''
  
  try {
    console.log('[WhatsApp Page] Calling whatsapp.connect()...')
    const response = await whatsapp.connect()
    console.log('[WhatsApp Page] Connect response:', response)
    
    if (response.status === 'already_connected' || response.status === 'connected') {
      // Already connected
      console.log('[WhatsApp Page] Already connected')
      isConnected.value = true
      await checkStatus()
    } else if (response.status === 'pairing_required') {
      // Need to show QR code
      console.log('[WhatsApp Page] Pairing required, showing QR')
      showQR.value = true
      
      // Start listening to QR stream
      if (response.stream_url) {
        console.log('[WhatsApp Page] Starting QR stream with URL:', response.stream_url)
        startQRStream()
      } else {
        console.error('[WhatsApp Page] No stream_url in response!')
      }
    } else {
      console.log('[WhatsApp Page] Unknown status:', response.status)
      error.value = response.message || 'Status tidak diketahui'
    }
  } catch (err: any) {
    console.error('[WhatsApp Page] Connect error:', err)
    const errorMsg = err.data?.error || err.message || 'Gagal menghubungkan WhatsApp'
    
    // Check if error is about tenant not found
    if (errorMsg.includes('Tenant not found') || errorMsg.includes('tenant')) {
      error.value = 'Tenant tidak ditemukan. Silakan buat tenant terlebih dahulu di halaman Pengaturan.'
    } else {
      error.value = errorMsg
    }
    
    console.error('Connect error:', err)
  } finally {
    loading.value = false
  }
}

const startQRStream = () => {
  console.log('[WhatsApp Page] startQRStream called')
  
  // Close existing stream if any
  if (qrStreamClose) {
    console.log('[WhatsApp Page] Closing existing stream')
    qrStreamClose()
    qrStreamClose = null
  }

  console.log('[WhatsApp Page] Setting up QR stream listeners')
  qrStreamClose = whatsapp.listenToQRStream(
    (qrCodeData: string) => {
      // QR code received
      console.log('[WhatsApp Page] QR code received, length:', qrCodeData.length)
      qrCode.value = qrCodeData
    },
    () => {
      // Connection successful
      console.log('[WhatsApp Page] Connection successful!')
      showQR.value = false
      qrCode.value = ''
      isConnected.value = true
      checkStatus()
      qrStreamClose = null
    },
    (errorMsg: string) => {
      // Error occurred
      console.error('[WhatsApp Page] QR stream error:', errorMsg)
      error.value = errorMsg
      showQR.value = false
      qrCode.value = ''
      qrStreamClose = null
    },
    () => {
      // Timeout
      console.log('[WhatsApp Page] QR stream timeout')
      error.value = 'QR code expired. Silakan coba lagi.'
      showQR.value = false
      qrCode.value = ''
      qrStreamClose = null
    }
  )
  console.log('[WhatsApp Page] QR stream listeners set up complete')
}

const handleDisconnect = async () => {
  if (!confirm('Apakah Anda yakin ingin memutuskan koneksi WhatsApp?')) return
  
  loading.value = true
  error.value = ''
  
  try {
    await whatsapp.disconnect()
    
    isConnected.value = false
    showQR.value = false
    qrCode.value = ''
    connectedNumber.value = ''
    lastConnected.value = ''
    
    // Close QR stream if active
    if (qrStreamClose) {
      qrStreamClose()
      qrStreamClose = null
    }
  } catch (err: any) {
    error.value = err.data?.error || err.message || 'Gagal memutuskan koneksi'
    console.error('Disconnect error:', err)
  } finally {
    loading.value = false
  }
}

const cancelQR = () => {
  showQR.value = false
  qrCode.value = ''
  
  // Close QR stream
  if (qrStreamClose) {
    qrStreamClose()
    qrStreamClose = null
  }
}

// Cleanup on unmount
onUnmounted(() => {
  if (qrStreamClose) {
    qrStreamClose()
  }
})
</script>
