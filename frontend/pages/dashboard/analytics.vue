<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Analytics Dashboard</h1>
        <p class="text-gray-600 mt-1">Analisis performa bisnis dan interaksi customer</p>
      </div>
      <div class="flex gap-2">
        <button
          v-for="p in periods"
          :key="p.value"
          @click="period = p.value; loadData()"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
            period === p.value 
              ? 'bg-indigo-600 text-white' 
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          {{ p.label }}
        </button>
      </div>
    </div>

    <!-- Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total Pesan</p>
            <p class="text-2xl font-bold text-gray-900">{{ formatNumber(overview.total_messages) }}</p>
            <p class="text-xs text-green-600">+{{ formatNumber(overview.today_messages) }} hari ini</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total Customer</p>
            <p class="text-2xl font-bold text-gray-900">{{ formatNumber(overview.total_customers) }}</p>
            <p class="text-xs text-green-600">+{{ overview.new_customers_week }} minggu ini</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">AI Responses</p>
            <p class="text-2xl font-bold text-gray-900">{{ formatNumber(overview.ai_responses_total) }}</p>
            <p class="text-xs text-gray-500">{{ formatCurrency(overview.ai_total_cost) }} total</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl p-5 border border-gray-100 shadow-sm">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-amber-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Customer Aktif</p>
            <p class="text-2xl font-bold text-gray-900">{{ formatNumber(overview.active_customers) }}</p>
            <p class="text-xs text-gray-500">7 hari terakhir</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Message Trend Chart -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Tren Pesan</h3>
        <div class="h-64">
          <canvas ref="messageChartRef"></canvas>
        </div>
      </div>

      <!-- Customer Growth Chart -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Pertumbuhan Customer</h3>
        <div class="h-64">
          <canvas ref="customerChartRef"></canvas>
        </div>
      </div>
    </div>

    <!-- AI Analytics & Hourly Distribution -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- AI Performance -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Performa AI</h3>
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="p-3 bg-green-50 rounded-lg">
            <p class="text-xs text-green-600 uppercase">Auto Reply Rate</p>
            <p class="text-2xl font-bold text-green-700">{{ aiStats.totals?.auto_reply_rate?.toFixed(1) || 0 }}%</p>
          </div>
          <div class="p-3 bg-blue-50 rounded-lg">
            <p class="text-xs text-blue-600 uppercase">Avg Confidence</p>
            <p class="text-2xl font-bold text-blue-700">{{ ((aiStats.totals?.avg_confidence || 0) * 100).toFixed(1) }}%</p>
          </div>
          <div class="p-3 bg-purple-50 rounded-lg">
            <p class="text-xs text-purple-600 uppercase">Total Tokens</p>
            <p class="text-2xl font-bold text-purple-700">{{ formatNumber(aiStats.totals?.total_tokens || 0) }}</p>
          </div>
          <div class="p-3 bg-amber-50 rounded-lg">
            <p class="text-xs text-amber-600 uppercase">Total Cost</p>
            <p class="text-2xl font-bold text-amber-700">{{ formatCurrency(aiStats.totals?.total_cost || 0) }}</p>
          </div>
        </div>
        <div class="h-40">
          <canvas ref="aiChartRef"></canvas>
        </div>
      </div>

      <!-- Hourly Distribution -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Distribusi Jam Pesan</h3>
        <p class="text-sm text-gray-500 mb-4">Waktu paling aktif customer menghubungi</p>
        <div class="h-52">
          <canvas ref="hourlyChartRef"></canvas>
        </div>
      </div>
    </div>

    <!-- Intent Distribution & Top Customers -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Intent Distribution -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Distribusi Intent</h3>
        <div class="space-y-3">
          <div v-for="intent in intents" :key="intent.intent" class="flex items-center gap-3">
            <div class="flex-1">
              <div class="flex justify-between text-sm mb-1">
                <span class="text-gray-700">{{ formatIntentName(intent.intent) }}</span>
                <span class="text-gray-500">{{ intent.count }}</span>
              </div>
              <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                <div 
                  class="h-full rounded-full"
                  :style="{ 
                    width: `${(intent.count / maxIntentCount) * 100}%`,
                    backgroundColor: getIntentColor(intent.intent)
                  }"
                ></div>
              </div>
            </div>
          </div>
          <p v-if="intents.length === 0" class="text-gray-500 text-sm text-center py-4">
            Belum ada data intent
          </p>
        </div>
      </div>

      <!-- Top Customers -->
      <div class="bg-white rounded-xl border border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Top Customers</h3>
        <div class="space-y-3">
          <div 
            v-for="(customer, index) in topCustomers" 
            :key="customer.id"
            class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50"
          >
            <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600 font-bold text-sm">
              {{ index + 1 }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="font-medium text-gray-900 truncate">{{ customer.name }}</p>
              <p class="text-xs text-gray-500">{{ customer.phone_number }}</p>
            </div>
            <div class="text-right">
              <p class="font-semibold text-gray-900">{{ customer.message_count }}</p>
              <p class="text-xs text-gray-500">pesan</p>
            </div>
          </div>
          <p v-if="topCustomers.length === 0" class="text-gray-500 text-sm text-center py-4">
            Belum ada data customer
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import Chart from 'chart.js/auto'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch: apiFetch } = useApi()

const period = ref('30d')
const periods = [
  { value: '7d', label: '7 Hari' },
  { value: '30d', label: '30 Hari' },
  { value: '90d', label: '90 Hari' }
]

const overview = ref<any>({})
const messageStats = ref<any[]>([])
const customerStats = ref<any[]>([])
const aiStats = ref<any>({ data: [], totals: {} })
const hourlyStats = ref<any[]>([])
const intents = ref<any[]>([])
const topCustomers = ref<any[]>([])

const messageChartRef = ref<HTMLCanvasElement | null>(null)
const customerChartRef = ref<HTMLCanvasElement | null>(null)
const aiChartRef = ref<HTMLCanvasElement | null>(null)
const hourlyChartRef = ref<HTMLCanvasElement | null>(null)

let messageChart: Chart | null = null
let customerChart: Chart | null = null
let aiChart: Chart | null = null
let hourlyChart: Chart | null = null

const maxIntentCount = computed(() => {
  if (intents.value.length === 0) return 1
  return Math.max(...intents.value.map(i => i.count))
})

const formatNumber = (num: number) => {
  if (!num) return '0'
  return num.toLocaleString('id-ID')
}

const formatCurrency = (num: number) => {
  if (!num) return '$0.00'
  return `$${num.toFixed(4)}`
}

const formatIntentName = (intent: string) => {
  const names: Record<string, string> = {
    'price_inquiry': '💰 Tanya Harga',
    'location_inquiry': '📍 Tanya Lokasi',
    'hours_inquiry': '🕐 Tanya Jam Buka',
    'availability_inquiry': '📦 Tanya Stok',
    'order_intent': '🛒 Mau Order',
    'complaint': '😤 Komplain',
    'shipping_inquiry': '🚚 Tanya Pengiriman',
    'payment_inquiry': '💳 Tanya Pembayaran',
    'general_inquiry': '❓ Pertanyaan Umum'
  }
  return names[intent] || intent
}

const getIntentColor = (intent: string) => {
  const colors: Record<string, string> = {
    'price_inquiry': '#f59e0b',
    'location_inquiry': '#3b82f6',
    'hours_inquiry': '#8b5cf6',
    'availability_inquiry': '#10b981',
    'order_intent': '#ef4444',
    'complaint': '#dc2626',
    'shipping_inquiry': '#06b6d4',
    'payment_inquiry': '#ec4899',
    'general_inquiry': '#6b7280'
  }
  return colors[intent] || '#6b7280'
}

const loadData = async () => {
  try {
    const [overviewRes, messagesRes, customersRes, aiRes, hourlyRes, intentsRes, topRes] = await Promise.all([
      apiFetch<any>('/api/analytics/overview'),
      apiFetch<any>(`/api/analytics/messages?period=${period.value}`),
      apiFetch<any>(`/api/analytics/customers?period=${period.value}`),
      apiFetch<any>(`/api/analytics/ai?period=${period.value}`),
      apiFetch<any>('/api/analytics/hourly'),
      apiFetch<any>('/api/analytics/intents'),
      apiFetch<any>('/api/analytics/top-customers')
    ])

    overview.value = overviewRes || {}
    messageStats.value = messagesRes?.data || []
    customerStats.value = customersRes?.data || []
    aiStats.value = aiRes || { data: [], totals: {} }
    hourlyStats.value = hourlyRes || []
    intents.value = intentsRes || []
    topCustomers.value = topRes || []

    await nextTick()
    renderCharts()
  } catch (err) {
    console.error('Failed to load analytics:', err)
  }
}

const renderCharts = () => {
  // Message Chart
  if (messageChartRef.value) {
    if (messageChart) messageChart.destroy()
    messageChart = new Chart(messageChartRef.value, {
      type: 'line',
      data: {
        labels: messageStats.value.map(s => formatDate(s.date)),
        datasets: [
          {
            label: 'Masuk',
            data: messageStats.value.map(s => s.received),
            borderColor: '#3b82f6',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            fill: true,
            tension: 0.4
          },
          {
            label: 'Keluar',
            data: messageStats.value.map(s => s.sent),
            borderColor: '#10b981',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            fill: true,
            tension: 0.4
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { position: 'bottom' } },
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  // Customer Chart
  if (customerChartRef.value) {
    if (customerChart) customerChart.destroy()
    customerChart = new Chart(customerChartRef.value, {
      type: 'line',
      data: {
        labels: customerStats.value.map(s => formatDate(s.date)),
        datasets: [
          {
            label: 'Total Customer',
            data: customerStats.value.map(s => s.total_active),
            borderColor: '#8b5cf6',
            backgroundColor: 'rgba(139, 92, 246, 0.1)',
            fill: true,
            tension: 0.4
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }

  // AI Chart
  if (aiChartRef.value) {
    if (aiChart) aiChart.destroy()
    aiChart = new Chart(aiChartRef.value, {
      type: 'bar',
      data: {
        labels: aiStats.value.data?.map((s: any) => formatDate(s.date)) || [],
        datasets: [
          {
            label: 'Responses',
            data: aiStats.value.data?.map((s: any) => s.responses) || [],
            backgroundColor: '#10b981'
          },
          {
            label: 'Escalations',
            data: aiStats.value.data?.map((s: any) => s.escalations) || [],
            backgroundColor: '#ef4444'
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { position: 'bottom' } },
        scales: {
          y: { beginAtZero: true, stacked: false },
          x: { stacked: false }
        }
      }
    })
  }

  // Hourly Chart
  if (hourlyChartRef.value) {
    if (hourlyChart) hourlyChart.destroy()
    hourlyChart = new Chart(hourlyChartRef.value, {
      type: 'bar',
      data: {
        labels: hourlyStats.value.map(s => `${s.hour}:00`),
        datasets: [{
          label: 'Pesan',
          data: hourlyStats.value.map(s => s.count),
          backgroundColor: hourlyStats.value.map(s => 
            s.hour >= 9 && s.hour <= 17 ? '#3b82f6' : '#94a3b8'
          ),
          borderRadius: 4
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          y: { beginAtZero: true }
        }
      }
    })
  }
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

onMounted(() => {
  loadData()
})
</script>

