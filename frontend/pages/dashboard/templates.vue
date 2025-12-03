<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Template Pesan</h1>
        <p class="text-gray-500 mt-1">Kelola template pesan untuk balasan cepat</p>
      </div>
      <button 
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Buat Template
      </button>
    </div>

    <!-- Category Filter -->
    <div class="flex gap-2 overflow-x-auto pb-2">
      <button
        @click="selectedCategory = 'all'"
        :class="selectedCategory === 'all' ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
        class="px-4 py-2 rounded-lg border whitespace-nowrap transition"
      >
        Semua
      </button>
      <button
        v-for="cat in templateHelper.categories"
        :key="cat.value"
        @click="selectedCategory = cat.value"
        :class="selectedCategory === cat.value ? 'bg-blue-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
        class="px-4 py-2 rounded-lg border whitespace-nowrap transition flex items-center gap-2"
      >
        <span>{{ cat.icon }}</span>
        {{ cat.label }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>

    <!-- Templates Grid -->
    <div v-else-if="templates.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div 
        v-for="template in templates" 
        :key="template.id"
        class="bg-white rounded-xl border border-gray-100 p-5 hover:shadow-lg transition group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="flex items-center gap-2">
            <span class="text-xl">{{ templateHelper.getCategoryInfo(template.category).icon }}</span>
            <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">
              {{ templateHelper.getCategoryInfo(template.category).label }}
            </span>
          </div>
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition">
            <button 
              @click="editTemplate(template)"
              class="p-1.5 hover:bg-gray-100 rounded"
              title="Edit"
            >
              <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
            </button>
            <button 
              @click="confirmDelete(template)"
              class="p-1.5 hover:bg-red-50 rounded"
              title="Hapus"
            >
              <svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>

        <h3 class="font-semibold text-gray-900 mb-2">{{ template.name }}</h3>
        <p class="text-gray-600 text-sm line-clamp-3 mb-4">{{ template.content }}</p>

        <div class="flex items-center justify-between text-xs text-gray-400">
          <span>Digunakan {{ template.usage_count }}x</span>
          <button 
            @click="copyToClipboard(template.content)"
            class="text-blue-600 hover:text-blue-800"
          >
            Salin
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white rounded-xl border border-gray-100 p-12 text-center">
      <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Belum ada template</h3>
      <p class="text-gray-500 mb-4">Buat template pesan untuk membalas pelanggan dengan cepat</p>
      <button 
        @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        Buat Template Pertama
      </button>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl w-full max-w-lg">
        <div class="p-6 border-b">
          <h2 class="text-xl font-bold">{{ showEditModal ? 'Edit Template' : 'Buat Template Baru' }}</h2>
        </div>

        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Template</label>
            <input 
              v-model="form.name"
              type="text"
              placeholder="Contoh: Sapaan Pembuka"
              class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Kategori</label>
            <select 
              v-model="form.category"
              class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option v-for="cat in templateHelper.categories" :key="cat.value" :value="cat.value">
                {{ cat.icon }} {{ cat.label }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Isi Pesan</label>
            <textarea 
              v-model="form.content"
              rows="5"
              placeholder="Ketik isi pesan template..."
              class="w-full px-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            ></textarea>
            <p class="text-xs text-gray-400 mt-1">Gunakan {"{{nama}}"} untuk variabel yang bisa diganti</p>
          </div>
        </div>

        <div class="p-6 border-t flex justify-end gap-3">
          <button 
            @click="closeModals"
            class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
          >
            Batal
          </button>
          <button 
            @click="saveTemplate"
            :disabled="saving || !form.name || !form.content"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl w-full max-w-sm">
        <div class="p-6 text-center">
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Hapus Template?</h3>
          <p class="text-gray-500 mb-6">Template "{{ templateToDelete?.name }}" akan dihapus permanen.</p>
          <div class="flex gap-3">
            <button 
              @click="showDeleteModal = false"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
            >
              Batal
            </button>
            <button 
              @click="deleteTemplateConfirm"
              :disabled="deleting"
              class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
            >
              {{ deleting ? 'Menghapus...' : 'Hapus' }}
            </button>
          </div>
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

const templateHelper = useTemplates()

// State
const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)
const templates = ref<Template[]>([])
const selectedCategory = ref('all')

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const templateToDelete = ref<Template | null>(null)
const editingTemplate = ref<Template | null>(null)

// Form
const form = ref({
  name: '',
  category: 'general',
  content: ''
})

// Load templates
const loadTemplates = async () => {
  loading.value = true
  try {
    const response = await templateHelper.getTemplates(selectedCategory.value)
    templates.value = response.templates
  } catch (error) {
    console.error('Failed to load templates:', error)
  } finally {
    loading.value = false
  }
}

// Watch category change
watch(selectedCategory, () => {
  loadTemplates()
})

// Initial load
onMounted(() => {
  loadTemplates()
})

// Edit template
const editTemplate = (template: Template) => {
  editingTemplate.value = template
  form.value = {
    name: template.name,
    category: template.category,
    content: template.content
  }
  showEditModal.value = true
}

// Confirm delete
const confirmDelete = (template: Template) => {
  templateToDelete.value = template
  showDeleteModal.value = true
}

// Save template
const saveTemplate = async () => {
  saving.value = true
  try {
    if (showEditModal.value && editingTemplate.value) {
      await templateHelper.updateTemplate(editingTemplate.value.id, form.value)
    } else {
      await templateHelper.createTemplate(form.value)
    }
    closeModals()
    loadTemplates()
  } catch (error) {
    console.error('Failed to save template:', error)
  } finally {
    saving.value = false
  }
}

// Delete template
const deleteTemplateConfirm = async () => {
  if (!templateToDelete.value) return
  
  deleting.value = true
  try {
    await templateHelper.deleteTemplate(templateToDelete.value.id)
    showDeleteModal.value = false
    templateToDelete.value = null
    loadTemplates()
  } catch (error) {
    console.error('Failed to delete template:', error)
  } finally {
    deleting.value = false
  }
}

// Close modals
const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingTemplate.value = null
  form.value = {
    name: '',
    category: 'general',
    content: ''
  }
}

// Copy to clipboard
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    // Could add toast notification here
  } catch (error) {
    console.error('Failed to copy:', error)
  }
}
</script>




