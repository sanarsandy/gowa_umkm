<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Customer Segments</h1>
        <p class="text-gray-600 mt-1">Kelola tag dan segmentasi customer</p>
      </div>
      <button
        @click="showTagModal = true; editingTag = null; tagForm = { name: '', color: '#6366f1', description: '' }"
        class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Buat Tag Baru
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-500">Total Tags</p>
        <p class="text-2xl font-bold text-gray-900">{{ tags.length }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-500">Total Customer</p>
        <p class="text-2xl font-bold text-gray-900">{{ totalCustomers }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-500">Tagged Customers</p>
        <p class="text-2xl font-bold text-gray-900">{{ taggedCustomers }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-gray-100">
        <p class="text-sm text-gray-500">Untagged</p>
        <p class="text-2xl font-bold text-gray-900">{{ totalCustomers - taggedCustomers }}</p>
      </div>
    </div>

    <!-- Tags Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="tag in tags"
        :key="tag.id"
        class="bg-white rounded-xl border border-gray-100 overflow-hidden hover:shadow-md transition-shadow"
      >
        <div class="p-5">
          <div class="flex items-start justify-between">
            <div class="flex items-center gap-3">
              <div 
                class="w-4 h-4 rounded-full"
                :style="{ backgroundColor: tag.color }"
              ></div>
              <div>
                <h3 class="font-semibold text-gray-900">{{ tag.name }}</h3>
                <p class="text-sm text-gray-500">{{ tag.customer_count }} customer</p>
              </div>
            </div>
            <div class="flex gap-1">
              <button 
                @click="editTag(tag)"
                class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
              </button>
              <button 
                @click="confirmDeleteTag(tag)"
                class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
          <p v-if="tag.description" class="text-sm text-gray-600 mt-3">{{ tag.description }}</p>
          
          <!-- View Customers Button -->
          <button
            @click="viewTagCustomers(tag)"
            class="mt-4 w-full py-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
          >
            Lihat Customer →
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="tags.length === 0" class="col-span-full text-center py-12 bg-white rounded-xl border border-gray-100">
        <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
        <p class="text-gray-500 mb-4">Belum ada tag customer</p>
        <button
          @click="showTagModal = true"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
        >
          Buat Tag Pertama
        </button>
      </div>
    </div>

    <!-- Tag Modal -->
    <Teleport to="body">
      <div v-if="showTagModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl w-full max-w-md mx-4 overflow-hidden">
          <div class="p-6 border-b border-gray-100">
            <h2 class="text-lg font-semibold text-gray-900">
              {{ editingTag ? 'Edit Tag' : 'Buat Tag Baru' }}
            </h2>
          </div>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nama Tag</label>
              <input
                v-model="tagForm.name"
                type="text"
                placeholder="Contoh: VIP, Hot Lead, Repeat Customer"
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Warna</label>
              <div class="flex gap-2">
                <button
                  v-for="color in tagColors"
                  :key="color"
                  @click="tagForm.color = color"
                  :class="[
                    'w-8 h-8 rounded-full border-2 transition-transform',
                    tagForm.color === color ? 'border-gray-900 scale-110' : 'border-transparent'
                  ]"
                  :style="{ backgroundColor: color }"
                ></button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Deskripsi (opsional)</label>
              <textarea
                v-model="tagForm.description"
                rows="2"
                placeholder="Deskripsi singkat tentang tag ini..."
                class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
              ></textarea>
            </div>
          </div>
          <div class="p-6 border-t border-gray-100 flex justify-end gap-3">
            <button
              @click="showTagModal = false"
              class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
            >
              Batal
            </button>
            <button
              @click="saveTag"
              :disabled="!tagForm.name || savingTag"
              class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
            >
              {{ savingTag ? 'Menyimpan...' : (editingTag ? 'Update' : 'Simpan') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Customers by Tag Modal -->
    <Teleport to="body">
      <div v-if="showCustomersModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl w-full max-w-2xl mx-4 max-h-[80vh] overflow-hidden flex flex-col">
          <div class="p-6 border-b border-gray-100 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div 
                class="w-4 h-4 rounded-full"
                :style="{ backgroundColor: selectedTag?.color }"
              ></div>
              <h2 class="text-lg font-semibold text-gray-900">
                Customer dengan tag "{{ selectedTag?.name }}"
              </h2>
            </div>
            <button @click="showCustomersModal = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto p-6">
            <div class="space-y-3">
              <div
                v-for="customer in tagCustomers"
                :key="customer.id"
                class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
              >
                <div>
                  <p class="font-medium text-gray-900">{{ customer.name }}</p>
                  <p class="text-sm text-gray-500">{{ customer.phone_number }}</p>
                </div>
                <div class="flex items-center gap-3">
                  <span :class="[
                    'px-2 py-1 text-xs rounded-full',
                    getLeadStatusColor(customer.lead_status)
                  ]">
                    {{ customer.lead_status }}
                  </span>
                  <NuxtLink 
                    :to="`/dashboard/customers/${customer.id}`"
                    class="text-indigo-600 hover:text-indigo-700 text-sm"
                  >
                    Detail →
                  </NuxtLink>
                </div>
              </div>
              <p v-if="tagCustomers.length === 0" class="text-center text-gray-500 py-8">
                Belum ada customer dengan tag ini
              </p>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
        <div class="bg-white rounded-xl w-full max-w-md mx-4 p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-2">Hapus Tag?</h2>
          <p class="text-gray-600 mb-6">
            Tag "{{ deletingTag?.name }}" akan dihapus dari {{ deletingTag?.customer_count }} customer. Aksi ini tidak dapat dibatalkan.
          </p>
          <div class="flex justify-end gap-3">
            <button
              @click="showDeleteModal = false"
              class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
            >
              Batal
            </button>
            <button
              @click="deleteTag"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
            >
              Hapus
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" :class="[
        'fixed bottom-6 right-6 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 z-50',
        toast.type === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
      ]">
        {{ toast.message }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch: apiFetch } = useApi()

interface Tag {
  id: string
  name: string
  color: string
  description: string
  customer_count: number
}

interface Customer {
  id: string
  name: string
  phone_number: string
  lead_score: number
  lead_status: string
}

const tags = ref<Tag[]>([])
const totalCustomers = ref(0)
const showTagModal = ref(false)
const showCustomersModal = ref(false)
const showDeleteModal = ref(false)
const editingTag = ref<Tag | null>(null)
const deletingTag = ref<Tag | null>(null)
const selectedTag = ref<Tag | null>(null)
const tagCustomers = ref<Customer[]>([])
const savingTag = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const tagForm = ref({
  name: '',
  color: '#6366f1',
  description: ''
})

const tagColors = [
  '#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#6366f1', 
  '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16', '#f97316'
]

const taggedCustomers = computed(() => {
  return tags.value.reduce((sum, tag) => sum + tag.customer_count, 0)
})

const getLeadStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'new': 'bg-blue-100 text-blue-700',
    'hot_lead': 'bg-red-100 text-red-700',
    'warm_lead': 'bg-orange-100 text-orange-700',
    'cold_lead': 'bg-gray-100 text-gray-700',
    'customer': 'bg-green-100 text-green-700',
    'complaint': 'bg-amber-100 text-amber-700',
    'spam': 'bg-slate-100 text-slate-700'
  }
  return colors[status] || 'bg-gray-100 text-gray-700'
}

const loadTags = async () => {
  try {
    const res = await apiFetch<Tag[]>('/api/tags')
    tags.value = res || []
  } catch (err) {
    console.error('Failed to load tags:', err)
  }
}

const loadCustomerCount = async () => {
  try {
    const res = await apiFetch<{ total: number }>('/api/customers/stats')
    totalCustomers.value = res?.total || 0
  } catch (err) {
    console.error('Failed to load customer count:', err)
  }
}

const editTag = (tag: Tag) => {
  editingTag.value = tag
  tagForm.value = {
    name: tag.name,
    color: tag.color,
    description: tag.description
  }
  showTagModal.value = true
}

const saveTag = async () => {
  if (!tagForm.value.name) return

  savingTag.value = true
  try {
    if (editingTag.value) {
      await apiFetch(`/api/tags/${editingTag.value.id}`, {
        method: 'PUT',
        body: JSON.stringify(tagForm.value)
      })
      showToast('Tag berhasil diupdate', 'success')
    } else {
      await apiFetch('/api/tags', {
        method: 'POST',
        body: JSON.stringify(tagForm.value)
      })
      showToast('Tag berhasil dibuat', 'success')
    }
    showTagModal.value = false
    await loadTags()
  } catch (err: any) {
    showToast(err.data?.error || 'Gagal menyimpan tag', 'error')
  } finally {
    savingTag.value = false
  }
}

const confirmDeleteTag = (tag: Tag) => {
  deletingTag.value = tag
  showDeleteModal.value = true
}

const deleteTag = async () => {
  if (!deletingTag.value) return

  try {
    await apiFetch(`/api/tags/${deletingTag.value.id}`, { method: 'DELETE' })
    showToast('Tag berhasil dihapus', 'success')
    showDeleteModal.value = false
    await loadTags()
  } catch (err: any) {
    showToast(err.data?.error || 'Gagal menghapus tag', 'error')
  }
}

const viewTagCustomers = async (tag: Tag) => {
  selectedTag.value = tag
  try {
    const res = await apiFetch<Customer[]>(`/api/tags/${tag.id}/customers`)
    tagCustomers.value = res || []
    showCustomersModal.value = true
  } catch (err) {
    console.error('Failed to load customers:', err)
  }
}

const showToast = (message: string, type: 'success' | 'error') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

onMounted(() => {
  loadTags()
  loadCustomerCount()
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>

