<template>
  <div class="space-y-6">
    <!-- Header with Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
      <div 
        v-for="stat in statusFilters" 
        :key="stat.value"
        @click="filterStatus = stat.value"
        :class="[
          'p-4 rounded-xl border cursor-pointer transition-all',
          filterStatus === stat.value 
            ? 'bg-blue-50 border-blue-300 shadow-md' 
            : 'bg-white border-gray-200 hover:border-gray-300'
        ]"
      >
        <p class="text-2xl font-bold" :class="stat.color">{{ getStatCount(stat.value) }}</p>
        <p class="text-xs text-gray-500 mt-1">{{ stat.label }}</p>
      </div>
    </div>

    <!-- Search and Filter Bar -->
    <div class="bg-white rounded-xl shadow-sm p-4 border border-gray-100">
      <div class="flex flex-col md:flex-row md:items-center gap-4">
        <!-- Search -->
        <div class="flex-1">
          <div class="relative">
            <input
              type="text"
              v-model="searchQuery"
              @input="debouncedSearch"
              placeholder="Cari nama, nomor HP, atau pesan..."
              class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <svg class="w-5 h-5 text-gray-400 absolute left-3 top-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
          </div>
        </div>

        <!-- Sort -->
        <div class="flex items-center gap-2">
          <select 
            v-model="sortBy" 
            @change="loadCustomers"
            class="px-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          >
            <option value="last_message_at">Terakhir Chat</option>
            <option value="message_count">Jumlah Pesan</option>
            <option value="created_at">Tanggal Dibuat</option>
            <option value="customer_name">Nama</option>
          </select>
          
          <button 
            @click="toggleSortOrder"
            class="p-2.5 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
            :title="sortOrder === 'desc' ? 'Terbaru dulu' : 'Terlama dulu'"
          >
            <svg v-if="sortOrder === 'desc'" class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"/>
            </svg>
            <svg v-else class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4"/>
            </svg>
          </button>

          <button 
            @click="loadCustomers"
            :disabled="loading"
            class="px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 flex items-center gap-2"
          >
            <svg class="w-4 h-4" :class="{'animate-spin': loading}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Customer List -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="p-12 text-center">
        <div class="animate-spin w-10 h-10 border-4 border-blue-600 border-t-transparent rounded-full mx-auto"></div>
        <p class="text-gray-500 mt-4">Memuat data pelanggan...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="customers.length === 0" class="text-center py-16">
        <div class="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">Belum ada pelanggan</h3>
        <p class="text-gray-500 mb-6">Pelanggan akan muncul otomatis saat ada pesan masuk dari WhatsApp</p>
        <NuxtLink to="/dashboard/whatsapp" class="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition">
          <svg class="w-5 h-5 mr-2" viewBox="0 0 24 24" fill="currentColor">
            <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
          </svg>
          Hubungkan WhatsApp
        </NuxtLink>
      </div>

      <!-- Customer Table -->
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Pelanggan</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Pesan Terakhir</th>
              <th class="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">Total Pesan</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Terakhir Aktif</th>
              <th class="px-6 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr 
              v-for="customer in customers" 
              :key="customer.id" 
              class="hover:bg-gray-50 transition cursor-pointer"
              @click="openCustomerDetail(customer.id)"
            >
              <!-- Customer Info -->
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center flex-shrink-0">
                    <span class="text-white font-bold text-sm">{{ getInitials(customerHelper.getDisplayName(customer)) }}</span>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-semibold text-gray-900">{{ customerHelper.getDisplayName(customer) }}</div>
                    <div class="text-xs text-gray-500">{{ customerHelper.formatPhoneFromJID(customer.customer_jid) }}</div>
                  </div>
                  <span v-if="customer.needs_follow_up" class="ml-2 px-2 py-0.5 bg-yellow-100 text-yellow-800 text-xs rounded-full">
                    Follow Up
                  </span>
                </div>
              </td>

              <!-- Status -->
              <td class="px-6 py-4">
                <span :class="customerHelper.getStatusColor(customer.status)" class="px-2.5 py-1 text-xs font-medium rounded-full">
                  {{ customerHelper.getStatusLabel(customer.status) }}
                </span>
              </td>

              <!-- Last Message -->
              <td class="px-6 py-4">
                <p class="text-sm text-gray-600 truncate max-w-xs">
                  {{ customer.last_message_summary || '-' }}
                </p>
              </td>

              <!-- Message Count -->
              <td class="px-6 py-4 text-center">
                <span class="text-sm font-medium text-gray-900">{{ customer.message_count }}</span>
              </td>

              <!-- Last Active -->
              <td class="px-6 py-4">
                <span class="text-sm text-gray-500">{{ customerHelper.formatRelativeTime(customer.last_message_at) }}</span>
              </td>

              <!-- Actions -->
              <td class="px-6 py-4 text-center" @click.stop>
                <div class="flex items-center justify-center gap-2">
                  <button 
                    @click="openCustomerDetail(customer.id)"
                    class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition"
                    title="Lihat Detail"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                  <button 
                    @click="toggleFollowUp(customer)"
                    :class="customer.needs_follow_up ? 'text-yellow-600 bg-yellow-50' : 'text-gray-400 hover:bg-gray-50'"
                    class="p-2 rounded-lg transition"
                    :title="customer.needs_follow_up ? 'Hapus Follow Up' : 'Tandai Follow Up'"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-100 flex items-center justify-between">
        <p class="text-sm text-gray-500">
          Menampilkan {{ (currentPage - 1) * limit + 1 }} - {{ Math.min(currentPage * limit, total) }} dari {{ total }} pelanggan
        </p>
        <div class="flex items-center gap-2">
          <button 
            @click="goToPage(currentPage - 1)" 
            :disabled="currentPage <= 1"
            class="px-3 py-1.5 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Sebelumnya
          </button>
          <span class="px-3 py-1.5 text-sm text-gray-600">
            Halaman {{ currentPage }} dari {{ totalPages }}
          </span>
          <button 
            @click="goToPage(currentPage + 1)" 
            :disabled="currentPage >= totalPages"
            class="px-3 py-1.5 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Customer, CustomerStats } from '~/composables/useCustomers'
import { useRealtimeUpdates, type NewMessageData } from '~/composables/useRealtimeUpdates'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const router = useRouter()
const customerHelper = useCustomers()

// State
const loading = ref(false)
const customers = ref<Customer[]>([])
const stats = ref<CustomerStats | null>(null)
const total = ref(0)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 20

// Filters
const searchQuery = ref('')
const filterStatus = ref('all')
const sortBy = ref('last_message_at')
const sortOrder = ref<'asc' | 'desc'>('desc')

// Real-time updates via WebSocket
useRealtimeUpdates({
  onNewMessage: (data: NewMessageData) => {
    console.log('[Customers] New message from:', data.sender_jid)
    
    // Find existing customer and update their last message
    const customerIndex = customers.value.findIndex(
      c => c.customer_jid === data.sender_jid
    )
    
    if (customerIndex !== -1) {
      // Update existing customer
      const customer = customers.value[customerIndex]
      customer.last_message_summary = data.message_text
      customer.message_count = (customer.message_count || 0) + 1
      customer.last_message_at = new Date(data.timestamp * 1000).toISOString()
      
      // Move to top of list if sorted by last_message_at
      if (sortBy.value === 'last_message_at') {
        customers.value.splice(customerIndex, 1)
        customers.value.unshift(customer)
      }
    } else {
      // New customer - reload the list
      loadCustomers()
    }
  },
  onNewCustomer: () => {
    // Reload stats and customer list
    loadCustomers()
    loadStats()
  }
})

// Status filter options
const statusFilters = [
  { value: 'all', label: 'Semua', color: 'text-gray-900' },
  { value: 'new', label: 'Baru', color: 'text-gray-600' },
  { value: 'hot_lead', label: 'Hot Lead', color: 'text-red-600' },
  { value: 'warm_lead', label: 'Warm Lead', color: 'text-orange-600' },
  { value: 'cold_lead', label: 'Cold Lead', color: 'text-blue-600' },
  { value: 'customer', label: 'Pelanggan', color: 'text-green-600' },
  { value: 'complaint', label: 'Komplain', color: 'text-yellow-600' },
  { value: 'spam', label: 'Spam', color: 'text-gray-400' },
]

// Get stat count for each status
const getStatCount = (status: string): number => {
  if (!stats.value) return 0
  if (status === 'all') return stats.value.total
  const key = status === 'new' ? 'new' : 
              status === 'hot_lead' ? 'hot_leads' :
              status === 'warm_lead' ? 'warm_leads' :
              status === 'cold_lead' ? 'cold_leads' :
              status === 'customer' ? 'customers' :
              status === 'complaint' ? 'complaints' : 'total'
  return (stats.value as any)[key] || 0
}

// Load customers
const loadCustomers = async () => {
  loading.value = true
  try {
    const response = await customerHelper.getCustomers({
      page: currentPage.value,
      limit,
      search: searchQuery.value || undefined,
      status: filterStatus.value !== 'all' ? filterStatus.value : undefined,
      sort_by: sortBy.value,
      sort_order: sortOrder.value
    })
    
    customers.value = response.customers
    total.value = response.total
    totalPages.value = response.total_pages
  } catch (error) {
    console.error('Failed to load customers:', error)
  } finally {
    loading.value = false
  }
}

// Load stats
const loadStats = async () => {
  try {
    stats.value = await customerHelper.getCustomerStats()
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

// Debounced search
let searchTimeout: NodeJS.Timeout
const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadCustomers()
  }, 300)
}

// Toggle sort order
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
  loadCustomers()
}

// Go to page
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadCustomers()
  }
}

// Open customer detail
const openCustomerDetail = (customerId: string) => {
  router.push(`/dashboard/customers/${customerId}`)
}

// Toggle follow up
const toggleFollowUp = async (customer: Customer) => {
  try {
    await customerHelper.updateCustomer(customer.id, {
      needs_follow_up: !customer.needs_follow_up
    })
    customer.needs_follow_up = !customer.needs_follow_up
  } catch (error) {
    console.error('Failed to update follow up:', error)
  }
}

// Get initials from name
const getInitials = (name: string): string => {
  if (!name) return '?'
  const parts = name.trim().split(/[\s@]+/)
  if (parts.length === 1) {
    if (/^\d+$/.test(parts[0])) {
      return parts[0].slice(-2)
    }
    return parts[0].substring(0, 2).toUpperCase()
  }
  return (parts[0][0] + (parts[1]?.[0] || '')).toUpperCase()
}

// Watch filter changes
watch(filterStatus, () => {
  currentPage.value = 1
  loadCustomers()
})

// Initial load
onMounted(() => {
  loadCustomers()
  loadStats()
})
</script>
