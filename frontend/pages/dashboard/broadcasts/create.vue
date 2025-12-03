<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center gap-4">
      <button 
        @click="router.back()" 
        class="p-2 hover:bg-gray-100 rounded-lg transition"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
        </svg>
      </button>
      <h1 class="text-xl font-bold text-gray-900">Buat Broadcast Baru</h1>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left: Form -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Broadcast Name -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <h2 class="font-semibold text-gray-900 mb-4">Detail Broadcast</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Nama Broadcast</label>
              <input 
                v-model="form.name"
                type="text"
                placeholder="Contoh: Promo Tahun Baru 2024"
                class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        <!-- Message Content -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="font-semibold text-gray-900">Isi Pesan</h2>
            <button 
              @click="showTemplateModal = true"
              class="text-blue-600 hover:text-blue-800 text-sm flex items-center gap-1"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
              </svg>
              Gunakan Template
            </button>
          </div>
          
          <textarea 
            v-model="form.message_content"
            rows="8"
            placeholder="Ketik pesan broadcast Anda di sini...

Contoh:
Halo {{nama}}, terima kasih sudah menjadi pelanggan kami! 🎉"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
          ></textarea>
          
          <!-- Placeholder Info -->
          <div class="mt-3 p-3 bg-blue-50 rounded-lg">
            <p class="text-sm text-blue-800 font-medium mb-1">💡 Personalisasi Pesan</p>
            <p class="text-xs text-blue-600">
              Gunakan <code class="px-1.5 py-0.5 bg-blue-100 rounded font-mono">{{nama}}</code> 
              untuk otomatis diganti dengan nama pelanggan.
            </p>
            <p class="text-xs text-blue-500 mt-1">
              Contoh: "Halo <strong>{{nama}}</strong>, ada promo spesial untuk Anda!" → 
              "Halo <strong>Budi</strong>, ada promo spesial untuk Anda!"
            </p>
          </div>
          
          <p class="text-xs text-gray-400 mt-2">
            Pesan akan dikirim ke {{ selectedCustomers.length }} pelanggan yang dipilih
          </p>
        </div>

        <!-- Scheduling Options -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="font-semibold text-gray-900">Jadwalkan Broadcast</h2>
              <p class="text-sm text-gray-500 mt-1">Atur waktu pengiriman otomatis</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                v-model="scheduling.enabled" 
                class="sr-only peer"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>

          <div v-if="scheduling.enabled" class="space-y-4 pt-4 border-t">
            <!-- Date & Time -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal</label>
                <input 
                  type="date"
                  v-model="scheduling.date"
                  :min="minDate"
                  class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Waktu</label>
                <input 
                  type="time"
                  v-model="scheduling.time"
                  class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
              </div>
            </div>

            <!-- Timezone -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Zona Waktu</label>
              <select 
                v-model="scheduling.timezone"
                class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="Asia/Jakarta">WIB (Jakarta)</option>
                <option value="Asia/Makassar">WITA (Makassar)</option>
                <option value="Asia/Jayapura">WIT (Jayapura)</option>
              </select>
            </div>

            <!-- Recurring Toggle -->
            <div class="pt-4 border-t">
              <label class="flex items-center gap-3 cursor-pointer">
                <input 
                  type="checkbox" 
                  v-model="scheduling.isRecurring"
                  class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                >
                <div>
                  <span class="font-medium text-gray-900">Broadcast Berulang</span>
                  <p class="text-sm text-gray-500">Kirim otomatis secara berkala</p>
                </div>
              </label>
            </div>

            <!-- Recurring Options -->
            <div v-if="scheduling.isRecurring" class="space-y-4 pl-7">
              <!-- Recurrence Type -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Jenis Pengulangan</label>
                <div class="grid grid-cols-3 gap-2">
                  <button
                    v-for="type in recurrenceTypes"
                    :key="type.value"
                    @click="scheduling.recurrenceType = type.value"
                    :class="[
                      'px-4 py-2 rounded-lg border-2 transition',
                      scheduling.recurrenceType === type.value
                        ? 'border-blue-500 bg-blue-50 text-blue-700'
                        : 'border-gray-200 hover:border-gray-300'
                    ]"
                  >
                    <div class="text-center">
                      <div class="text-lg mb-1">{{ type.icon }}</div>
                      <div class="text-sm font-medium">{{ type.label }}</div>
                    </div>
                  </button>
                </div>
              </div>

              <!-- Interval -->
              <div v-if="scheduling.recurrenceType !== 'weekly'">
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Setiap {{ scheduling.recurrenceType === 'hourly' ? 'Jam' : 'Hari' }}
                </label>
                <input 
                  type="number"
                  v-model.number="scheduling.recurrenceInterval"
                  min="1"
                  :max="scheduling.recurrenceType === 'hourly' ? 24 : 30"
                  class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
              </div>

              <!-- Weekly Days Selection -->
              <div v-if="scheduling.recurrenceType === 'weekly'">
                <label class="block text-sm font-medium text-gray-700 mb-2">Pilih Hari</label>
                <div class="grid grid-cols-7 gap-2">
                  <button
                    v-for="day in weekDays"
                    :key="day.value"
                    @click="toggleWeekDay(day.value)"
                    :class="[
                      'px-2 py-2 rounded-lg border-2 transition text-sm font-medium',
                      scheduling.recurrenceDays.includes(day.value)
                        ? 'border-blue-500 bg-blue-50 text-blue-700'
                        : 'border-gray-200 hover:border-gray-300'
                    ]"
                  >
                    {{ day.short }}
                  </button>
                </div>
              </div>

              <!-- Recurrence Time (for daily/weekly) -->
              <div v-if="scheduling.recurrenceType !== 'hourly'">
                <label class="block text-sm font-medium text-gray-700 mb-1">Jam Pengiriman</label>
                <input 
                  type="time"
                  v-model="scheduling.recurrenceTime"
                  class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
              </div>

              <!-- End Condition -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Berhenti Setelah</label>
                <div class="space-y-2">
                  <label class="flex items-center gap-2">
                    <input 
                      type="radio" 
                      value="date"
                      v-model="scheduling.endType"
                      class="w-4 h-4 text-blue-600"
                    >
                    <span class="text-sm">Tanggal tertentu</span>
                  </label>
                  <input 
                    v-if="scheduling.endType === 'date'"
                    type="date"
                    v-model="scheduling.endDate"
                    :min="scheduling.date"
                    class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ml-6"
                  >
                  
                  <label class="flex items-center gap-2">
                    <input 
                      type="radio" 
                      value="count"
                      v-model="scheduling.endType"
                      class="w-4 h-4 text-blue-600"
                    >
                    <span class="text-sm">Jumlah pengiriman</span>
                  </label>
                  <input 
                    v-if="scheduling.endType === 'count'"
                    type="number"
                    v-model.number="scheduling.endCount"
                    min="1"
                    max="100"
                    class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ml-6"
                  >
                </div>
              </div>
            </div>

            <!-- Schedule Preview -->
            <div class="mt-4 p-4 bg-blue-50 rounded-lg border border-blue-200">
              <p class="text-sm font-medium text-blue-900 mb-1">📅 Preview Jadwal</p>
              <p class="text-sm text-blue-700">{{ schedulePreview }}</p>
            </div>

            <!-- Validation Errors -->
            <div v-if="validationErrors.length > 0" class="mt-4 p-4 bg-red-50 rounded-lg border border-red-200">
              <p class="text-sm font-medium text-red-900 mb-2">⚠️ Error Validasi:</p>
              <ul class="list-disc list-inside space-y-1">
                <li v-for="(error, index) in validationErrors" :key="index" class="text-sm text-red-700">
                  {{ error }}
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Select Customers -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="font-semibold text-gray-900">Pilih Penerima</h2>
            <div class="flex items-center gap-2">
              <button 
                @click="selectAll"
                class="text-blue-600 hover:text-blue-800 text-sm"
              >
                Pilih Semua
              </button>
              <span class="text-gray-300">|</span>
              <button 
                @click="deselectAll"
                class="text-gray-500 hover:text-gray-700 text-sm"
              >
                Batalkan Semua
              </button>
            </div>
          </div>

          <!-- Filters -->
          <div class="mb-4 space-y-3">
            <!-- Status Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Filter berdasarkan Status</label>
              <select 
                v-model="statusFilter"
                class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="all">Semua Status</option>
                <option value="hot_lead">🔥 Hot Lead</option>
                <option value="warm_lead">🌡️ Warm Lead</option>
                <option value="cold_lead">❄️ Cold Lead</option>
                <option value="customer">✅ Customer</option>
                <option value="new">🆕 New</option>
                <option value="complaint">⚠️ Complaint</option>
              </select>
            </div>

            <!-- Search -->
            <div>
              <input 
                v-model="searchQuery"
                type="text"
                placeholder="Cari pelanggan..."
                class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <!-- Customer List -->
          <div v-if="loadingCustomers" class="py-8 text-center">
            <div class="animate-spin w-6 h-6 border-4 border-blue-600 border-t-transparent rounded-full mx-auto"></div>
          </div>
          
          <div v-else class="max-h-80 overflow-y-auto space-y-2">
            <label 
              v-for="customer in filteredCustomers" 
              :key="customer.id"
              class="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg cursor-pointer"
            >
              <input 
                type="checkbox"
                :value="customer.id"
                v-model="selectedCustomers"
                class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
              />
              <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                <span class="text-white font-medium text-sm">{{ getInitials(customer.customer_name || customer.customer_phone) }}</span>
              </div>
              <div class="flex-1">
                <p class="font-medium text-gray-900">{{ customer.customer_name || customer.customer_phone }}</p>
                <p class="text-gray-500 text-sm">{{ customer.customer_phone }}</p>
              </div>
            </label>
          </div>

          <p v-if="filteredCustomers.length === 0 && !loadingCustomers" class="text-center text-gray-500 py-8">
            Tidak ada pelanggan ditemukan
          </p>
        </div>
      </div>

      <!-- Right: Summary & Actions -->
      <div class="space-y-6">
        <!-- Summary -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <h2 class="font-semibold text-gray-900 mb-4">Ringkasan</h2>
          
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-500">Penerima</span>
              <span class="font-medium">{{ selectedCustomers.length }} pelanggan</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">Panjang Pesan</span>
              <span class="font-medium">{{ form.message_content.length }} karakter</span>
            </div>
          </div>
        </div>

        <!-- Preview -->
        <div class="bg-white rounded-xl border border-gray-100 p-6">
          <h2 class="font-semibold text-gray-900 mb-4">Preview Pesan</h2>
          
          <div class="bg-green-50 rounded-lg p-4">
            <p class="text-gray-700 text-sm whitespace-pre-wrap">{{ form.message_content || 'Pesan akan muncul di sini...' }}</p>
          </div>
        </div>

        <!-- Actions -->
        <div class="space-y-3">
          <button 
            @click="sendNow"
            :disabled="!canSend || sending"
            class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="sending" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
            </svg>
            {{ sending ? 'Memproses...' : (scheduling.enabled ? 'Jadwalkan Broadcast' : 'Kirim Sekarang') }}
          </button>

          <button 
            @click="saveAsDraft"
            :disabled="!form.name || saving"
            class="w-full px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 disabled:opacity-50"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan sebagai Draft' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Template Modal -->
    <div v-if="showTemplateModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl w-full max-w-lg max-h-[80vh] flex flex-col">
        <div class="p-6 border-b flex items-center justify-between">
          <h2 class="text-xl font-bold">Pilih Template</h2>
          <button @click="showTemplateModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="flex-1 overflow-y-auto p-6">
          <div v-if="loadingTemplates" class="py-8 text-center">
            <div class="animate-spin w-6 h-6 border-4 border-blue-600 border-t-transparent rounded-full mx-auto"></div>
          </div>
          
          <div v-else class="space-y-3">
            <div 
              v-for="template in templates"
              :key="template.id"
              @click="useTemplate(template)"
              class="p-4 border border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 cursor-pointer transition"
            >
              <div class="flex items-center gap-2 mb-2">
                <span>{{ templateHelper.getCategoryInfo(template.category).icon }}</span>
                <span class="font-medium">{{ template.name }}</span>
              </div>
              <p class="text-gray-600 text-sm line-clamp-2">{{ template.content }}</p>
            </div>
          </div>

          <p v-if="templates.length === 0 && !loadingTemplates" class="text-center text-gray-500 py-8">
            Belum ada template
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Template } from '~/composables/useTemplates'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const router = useRouter()
const broadcastHelper = useBroadcasts()
const templateHelper = useTemplates()
const customerHelper = useCustomers()

// State
const sending = ref(false)
const saving = ref(false)
const loadingCustomers = ref(true)
const loadingTemplates = ref(false)
const showTemplateModal = ref(false)
const searchQuery = ref('')
const statusFilter = ref('all')

// Form
const form = ref({
  name: '',
  message_content: ''
})

// Scheduling
const scheduling = ref({
  enabled: false,
  date: '',
  time: '',
  timezone: 'Asia/Jakarta',
  isRecurring: false,
  recurrenceType: 'daily',
  recurrenceInterval: 1,
  recurrenceDays: [] as string[],
  recurrenceTime: '09:00',
  endType: 'count',
  endDate: '',
  endCount: 10
})

// Recurrence types
const recurrenceTypes = [
  { value: 'hourly', label: 'Per Jam', icon: '⏰' },
  { value: 'daily', label: 'Harian', icon: '📅' },
  { value: 'weekly', label: 'Mingguan', icon: '📆' }
]

// Week days
const weekDays = [
  { value: 'sunday', short: 'Min', full: 'Minggu' },
  { value: 'monday', short: 'Sen', full: 'Senin' },
  { value: 'tuesday', short: 'Sel', full: 'Selasa' },
  { value: 'wednesday', short: 'Rab', full: 'Rabu' },
  { value: 'thursday', short: 'Kam', full: 'Kamis' },
  { value: 'friday', short: 'Jum', full: 'Jumat' },
  { value: 'saturday', short: 'Sab', full: 'Sabtu' }
]

// Customers
const customers = ref<any[]>([])
const selectedCustomers = ref<string[]>([])

// Templates
const templates = ref<Template[]>([])

// Computed - Min date (today)
const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

// Computed - Schedule Preview
const schedulePreview = computed(() => {
  if (!scheduling.value.enabled) return ''
  
  const date = scheduling.value.date
  const time = scheduling.value.time
  
  if (!date || !time) return 'Pilih tanggal dan waktu terlebih dahulu'
  
  let preview = `Akan dikirim pada ${new Date(date + 'T' + time).toLocaleDateString('id-ID', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })} pukul ${time} ${scheduling.value.timezone.split('/')[1]}`
  
  if (scheduling.value.isRecurring) {
    const type = scheduling.value.recurrenceType
    
    if (type === 'hourly') {
      preview += `, kemudian setiap ${scheduling.value.recurrenceInterval} jam`
    } else if (type === 'daily') {
      preview += `, kemudian setiap ${scheduling.value.recurrenceInterval} hari pada pukul ${scheduling.value.recurrenceTime}`
    } else if (type === 'weekly') {
      const days = scheduling.value.recurrenceDays.map(d => 
        weekDays.find(wd => wd.value === d)?.full || d
      ).join(', ')
      preview += `, kemudian setiap ${days} pada pukul ${scheduling.value.recurrenceTime}`
    }
    
    if (scheduling.value.endType === 'date' && scheduling.value.endDate) {
      preview += ` sampai ${new Date(scheduling.value.endDate).toLocaleDateString('id-ID')}`
    } else if (scheduling.value.endType === 'count') {
      preview += ` sebanyak ${scheduling.value.endCount}x`
    }
  }
  
  return preview
})

// Computed
const filteredCustomers = computed(() => {
  let filtered = customers.value
  
  // Filter by status
  if (statusFilter.value !== 'all') {
    filtered = filtered.filter(c => c.status === statusFilter.value)
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(c => 
      c.customer_name?.toLowerCase().includes(query) ||
      c.customer_phone?.toLowerCase().includes(query)
    )
  }
  
  return filtered
})

// Validation errors
const validationErrors = ref<string[]>([])

// Validate scheduling
const validateScheduling = (): boolean => {
  validationErrors.value = []
  
  if (!scheduling.value.enabled) return true
  
  // Check date and time
  if (!scheduling.value.date || !scheduling.value.time) {
    validationErrors.value.push('Tanggal dan waktu harus diisi')
    return false
  }
  
  // Check if date is in the future
  const scheduledDateTime = new Date(`${scheduling.value.date}T${scheduling.value.time}:00`)
  const now = new Date()
  if (scheduledDateTime <= now) {
    validationErrors.value.push('Waktu penjadwalan harus di masa depan')
    return false
  }
  
  // Check recurring options
  if (scheduling.value.isRecurring) {
    // Check interval
    if (scheduling.value.recurrenceInterval < 1) {
      validationErrors.value.push('Interval harus minimal 1')
      return false
    }
    
    // Check weekly days
    if (scheduling.value.recurrenceType === 'weekly' && scheduling.value.recurrenceDays.length === 0) {
      validationErrors.value.push('Pilih minimal 1 hari untuk pengulangan mingguan')
      return false
    }
    
    // Check recurrence time for daily/weekly
    if (scheduling.value.recurrenceType !== 'hourly' && !scheduling.value.recurrenceTime) {
      validationErrors.value.push('Jam pengiriman harus diisi')
      return false
    }
    
    // Check end conditions
    if (scheduling.value.endType === 'date') {
      if (!scheduling.value.endDate) {
        validationErrors.value.push('Tanggal akhir harus diisi')
        return false
      }
      const endDate = new Date(scheduling.value.endDate)
      if (endDate <= scheduledDateTime) {
        validationErrors.value.push('Tanggal akhir harus setelah tanggal mulai')
        return false
      }
    } else if (scheduling.value.endType === 'count') {
      if (scheduling.value.endCount < 1) {
        validationErrors.value.push('Jumlah pengiriman harus minimal 1')
        return false
      }
    }
  }
  
  return true
}

const canSend = computed(() => {
  return form.value.name && form.value.message_content && selectedCustomers.value.length > 0
})

// Load customers
const loadCustomers = async () => {
  loadingCustomers.value = true
  try {
    const response = await customerHelper.getCustomers({ limit: 100 })
    customers.value = response.customers
  } catch (error) {
    console.error('Failed to load customers:', error)
  } finally {
    loadingCustomers.value = false
  }
}

// Load templates
const loadTemplates = async () => {
  loadingTemplates.value = true
  try {
    const response = await templateHelper.getTemplates()
    templates.value = response.templates
  } catch (error) {
    console.error('Failed to load templates:', error)
  } finally {
    loadingTemplates.value = false
  }
}

// Select all customers
const selectAll = () => {
  selectedCustomers.value = filteredCustomers.value.map(c => c.id)
}

// Deselect all
const deselectAll = () => {
  selectedCustomers.value = []
}

// Use template
const useTemplate = (template: Template) => {
  form.value.message_content = template.content
  showTemplateModal.value = false
  
  // Increment usage
  templateHelper.incrementUsage(template.id)
}

// Get initials
const getInitials = (name: string | null): string => {
  if (!name) return '?'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

// Toggle week day
const toggleWeekDay = (day: string) => {
  const index = scheduling.value.recurrenceDays.indexOf(day)
  if (index > -1) {
    scheduling.value.recurrenceDays.splice(index, 1)
  } else {
    scheduling.value.recurrenceDays.push(day)
  }
}

// Save as draft
const saveAsDraft = async () => {
  // Validate scheduling if enabled
  if (scheduling.value.enabled && !validateScheduling()) {
    return
  }
  
  saving.value = true
  try {
    const payload: any = {
      name: form.value.name,
      message_content: form.value.message_content,
      customer_ids: selectedCustomers.value
    }
    
    // Add scheduling data if enabled
    if (scheduling.value.enabled && scheduling.value.date && scheduling.value.time) {
      const scheduledAt = new Date(`${scheduling.value.date}T${scheduling.value.time}:00`)
      payload.scheduled_at = scheduledAt.toISOString()
      
      if (scheduling.value.isRecurring) {
        payload.is_recurring = true
        payload.recurrence_type = scheduling.value.recurrenceType
        payload.recurrence_interval = scheduling.value.recurrenceInterval
        
        if (scheduling.value.recurrenceType === 'weekly') {
          payload.recurrence_days = scheduling.value.recurrenceDays
        }
        
        if (scheduling.value.recurrenceType !== 'hourly') {
          payload.recurrence_time = scheduling.value.recurrenceTime
        }
        
        if (scheduling.value.endType === 'date' && scheduling.value.endDate) {
          payload.recurrence_end_date = new Date(scheduling.value.endDate).toISOString()
        } else if (scheduling.value.endType === 'count') {
          payload.recurrence_count = scheduling.value.endCount
        }
      }
    }
    
    await broadcastHelper.createBroadcast(payload)
    router.push('/dashboard/broadcasts')
  } catch (error) {
    console.error('Failed to save broadcast:', error)
  } finally {
    saving.value = false
  }
}

// Send now
const sendNow = async () => {
  // Validate scheduling if enabled
  if (scheduling.value.enabled && !validateScheduling()) {
    return
  }
  
  sending.value = true
  try {
    const payload: any = {
      name: form.value.name,
      message_content: form.value.message_content,
      customer_ids: selectedCustomers.value
    }
    
    // Add scheduling data if enabled
    if (scheduling.value.enabled && scheduling.value.date && scheduling.value.time) {
      const scheduledAt = new Date(`${scheduling.value.date}T${scheduling.value.time}:00`)
      payload.scheduled_at = scheduledAt.toISOString()
      
      if (scheduling.value.isRecurring) {
        payload.is_recurring = true
        payload.recurrence_type = scheduling.value.recurrenceType
        payload.recurrence_interval = scheduling.value.recurrenceInterval
        
        if (scheduling.value.recurrenceType === 'weekly') {
          payload.recurrence_days = scheduling.value.recurrenceDays
        }
        
        if (scheduling.value.recurrenceType !== 'hourly') {
          payload.recurrence_time = scheduling.value.recurrenceTime
        }
        
        if (scheduling.value.endType === 'date' && scheduling.value.endDate) {
          payload.recurrence_end_date = new Date(scheduling.value.endDate).toISOString()
        } else if (scheduling.value.endType === 'count') {
          payload.recurrence_count = scheduling.value.endCount
        }
      }
    }
    
    // Create broadcast
    const broadcast = await broadcastHelper.createBroadcast(payload)
    
    // Only send immediately if not scheduled
    if (!scheduling.value.enabled) {
      await broadcastHelper.sendBroadcast(broadcast.id)
    }
    
    router.push(`/dashboard/broadcasts/${broadcast.id}`)
  } catch (error) {
    console.error('Failed to send broadcast:', error)
  } finally {
    sending.value = false
  }
}

// Watch template modal
watch(showTemplateModal, (show) => {
  if (show && templates.value.length === 0) {
    loadTemplates()
  }
})

// Initial load
onMounted(() => {
  loadCustomers()
})
</script>

