<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Knowledge Base</h1>
        <p class="text-gray-600 mt-1">Kelola pengetahuan AI untuk menjawab customer</p>
      </div>
      <div class="flex gap-3">
        <NuxtLink to="/dashboard/ai/config" class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Konfigurasi AI
        </NuxtLink>
        <button
          @click="showModal = true; editingItem = null"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Tambah Knowledge
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total</p>
            <p class="text-xl font-bold text-gray-900">{{ stats.total || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Aktif</p>
            <p class="text-xl font-bold text-gray-900">{{ stats.active || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Kategori</p>
            <p class="text-xl font-bold text-gray-900">{{ Object.keys(categoryCount).length }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl border border-gray-100 p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
            </svg>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total Digunakan</p>
            <p class="text-xl font-bold text-gray-900">{{ totalUsage }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Filter & Search -->
    <div class="bg-white rounded-xl border border-gray-100 p-4">
      <div class="flex flex-col md:flex-row gap-4">
        <!-- Search -->
        <div class="flex-1">
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari knowledge..."
              class="w-full pl-10 pr-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
          </div>
        </div>
        
        <!-- Category Filter -->
        <select
          v-model="selectedCategory"
          class="px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
        >
          <option value="">Semua Kategori</option>
          <option v-for="cat in categories" :key="cat" :value="cat">{{ formatCategory(cat) }}</option>
        </select>

        <!-- Active Filter -->
        <select
          v-model="activeFilter"
          class="px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
        >
          <option value="">Semua Status</option>
          <option value="true">Aktif</option>
          <option value="false">Nonaktif</option>
        </select>
      </div>
    </div>

    <!-- Knowledge List -->
    <div class="bg-white rounded-xl border border-gray-100 overflow-hidden">
      <div v-if="loading" class="p-12 text-center">
        <div class="inline-flex items-center gap-3 text-gray-500">
          <svg class="w-6 h-6 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Memuat data...
        </div>
      </div>

      <div v-else-if="filteredKnowledge.length === 0" class="p-12 text-center">
        <svg class="w-16 h-16 mx-auto text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">Belum ada knowledge</h3>
        <p class="mt-1 text-gray-500">Mulai tambahkan pengetahuan untuk AI</p>
        <button
          @click="showModal = true"
          class="mt-4 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
        >
          Tambah Knowledge Pertama
        </button>
      </div>

      <div v-else class="divide-y divide-gray-100">
        <div
          v-for="item in filteredKnowledge"
          :key="item.id"
          class="p-4 hover:bg-gray-50 transition-colors"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <h3 class="font-medium text-gray-900 truncate">{{ item.title }}</h3>
                <span
                  :class="item.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-600'"
                  class="px-2 py-0.5 text-xs font-medium rounded-full"
                >
                  {{ item.is_active ? 'Aktif' : 'Nonaktif' }}
                </span>
              </div>
              <p class="mt-1 text-sm text-gray-600 line-clamp-2">{{ item.content }}</p>
              <div class="mt-2 flex flex-wrap items-center gap-2">
                <span v-if="item.category" class="px-2 py-1 bg-indigo-50 text-indigo-600 text-xs font-medium rounded">
                  {{ formatCategory(item.category) }}
                </span>
                <span
                  v-for="keyword in (item.keywords || []).slice(0, 3)"
                  :key="keyword"
                  class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded"
                >
                  {{ keyword }}
                </span>
                <span v-if="(item.keywords || []).length > 3" class="text-xs text-gray-400">
                  +{{ item.keywords.length - 3 }} lainnya
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <div class="text-right mr-4">
                <p class="text-xs text-gray-500">Digunakan</p>
                <p class="font-semibold text-gray-900">{{ item.usage_count || 0 }}x</p>
              </div>
              <button
                @click="editItem(item)"
                class="p-2 text-gray-400 hover:text-indigo-600 transition-colors"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button
                @click="deleteItem(item)"
                class="p-2 text-gray-400 hover:text-red-600 transition-colors"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showModal" class="fixed inset-0 z-50 overflow-y-auto">
          <div class="flex min-h-screen items-center justify-center p-4">
            <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
            
            <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-2xl">
              <!-- Modal Header -->
              <div class="flex items-center justify-between p-6 border-b border-gray-100">
                <h2 class="text-xl font-bold text-gray-900">
                  {{ editingItem ? 'Edit Knowledge' : 'Tambah Knowledge' }}
                </h2>
                <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <!-- Modal Body -->
              <div class="p-6 space-y-4 max-h-[60vh] overflow-y-auto">
                <!-- Title -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Judul *</label>
                  <input
                    v-model="form.title"
                    type="text"
                    placeholder="Contoh: Jam Operasional Toko"
                    class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  >
                </div>

                <!-- Content -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Konten *</label>
                  <textarea
                    v-model="form.content"
                    rows="5"
                    placeholder="Tulis informasi yang akan digunakan AI untuk menjawab pertanyaan customer..."
                    class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
                  ></textarea>
                </div>

                <!-- Media Upload -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Lampiran Media (Opsional)</label>
                  <div v-if="!form.media_url" class="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center hover:border-indigo-400 transition cursor-pointer" @click="openKnowledgeFilePicker">
                    <svg class="w-8 h-8 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                    <p class="mt-2 text-sm text-gray-600">Klik untuk upload gambar atau dokumen</p>
                    <p class="text-xs text-gray-400">PNG, JPG, PDF max 10MB</p>
                  </div>
                  <div v-else class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                    <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
                      <svg v-if="form.media_type === 'image'" class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                      <svg v-else class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                      </svg>
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-medium text-gray-900 truncate">{{ uploadedFileName || 'File terlampir' }}</p>
                      <p class="text-xs text-gray-500">{{ form.media_type }}</p>
                    </div>
                    <button @click="removeMedia" class="p-1 hover:bg-gray-200 rounded" type="button">
                      <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                  <input 
                    ref="knowledgeFileInput"
                    type="file"
                    accept="image/*,.pdf,.doc,.docx"
                    class="hidden"
                    @change="handleKnowledgeFileUpload"
                  />
                  <div v-if="uploadingKnowledgeFile" class="mt-2 flex items-center gap-2 text-sm text-gray-500">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Mengunggah...</span>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <!-- Category -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Kategori</label>
                    <select
                      v-model="form.category"
                      class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    >
                      <option value="">Pilih kategori</option>
                      <option value="faq">FAQ</option>
                      <option value="product">Produk</option>
                      <option value="pricing">Harga</option>
                      <option value="shipping">Pengiriman</option>
                      <option value="hours">Jam Operasional</option>
                      <option value="location">Lokasi</option>
                      <option value="payment">Pembayaran</option>
                      <option value="policy">Kebijakan</option>
                    </select>
                  </div>

                  <!-- Priority -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Prioritas</label>
                    <select
                      v-model.number="form.priority"
                      class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                    >
                      <option :value="10">10 - Sangat Tinggi</option>
                      <option :value="8">8 - Tinggi</option>
                      <option :value="5">5 - Normal</option>
                      <option :value="3">3 - Rendah</option>
                      <option :value="1">1 - Sangat Rendah</option>
                    </select>
                  </div>
                </div>

                <!-- Keywords -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Keywords</label>
                  <input
                    v-model="keywordsInput"
                    type="text"
                    placeholder="Pisahkan dengan koma: harga, biaya, tarif"
                    class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  >
                  <p class="mt-1 text-xs text-gray-500">Keywords membantu AI mencocokkan pertanyaan customer</p>
                </div>

                <!-- Tags -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Tags</label>
                  <input
                    v-model="tagsInput"
                    type="text"
                    placeholder="Pisahkan dengan koma: promo, terbaru"
                    class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  >
                </div>

                <!-- Is Active -->
                <label class="flex items-center gap-3 cursor-pointer">
                  <input type="checkbox" v-model="form.is_active" class="w-4 h-4 text-indigo-600 rounded">
                  <span class="text-gray-700">Aktifkan knowledge ini</span>
                </label>
              </div>

              <!-- Modal Footer -->
              <div class="flex items-center justify-end gap-3 p-6 border-t border-gray-100">
                <button
                  @click="closeModal"
                  class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  Batal
                </button>
                <button
                  @click="saveKnowledge"
                  :disabled="saving || !form.title || !form.content"
                  class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ saving ? 'Menyimpan...' : (editingItem ? 'Update' : 'Simpan') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteConfirm" class="fixed inset-0 z-50 overflow-y-auto">
          <div class="flex min-h-screen items-center justify-center p-4">
            <div class="fixed inset-0 bg-black/50" @click="deleteConfirm = null"></div>
            
            <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
              <div class="text-center">
                <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto">
                  <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </div>
                <h3 class="mt-4 text-lg font-semibold text-gray-900">Hapus Knowledge?</h3>
                <p class="mt-2 text-gray-600">
                  Yakin ingin menghapus "{{ deleteConfirm?.title }}"? Aksi ini tidak dapat dibatalkan.
                </p>
              </div>
              <div class="mt-6 flex justify-center gap-3">
                <button
                  @click="deleteConfirm = null"
                  class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200"
                >
                  Batal
                </button>
                <button
                  @click="confirmDelete"
                  :disabled="deleting"
                  class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
                >
                  {{ deleting ? 'Menghapus...' : 'Hapus' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Toast Notification -->
    <Transition name="slide-up">
      <div v-if="toast.show" :class="[
        'fixed bottom-6 right-6 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 z-50',
        toast.type === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
      ]">
        <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        {{ toast.message }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

interface KnowledgeItem {
  id: string
  title: string
  content: string
  category: string
  keywords: string[]
  tags: string[]
  priority: number
  usage_count: number
  is_active: boolean
  media_url?: string
  media_type?: string
  created_at: string
}

const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)
const showModal = ref(false)
const editingItem = ref<KnowledgeItem | null>(null)
const deleteConfirm = ref<KnowledgeItem | null>(null)
const toast = ref({ show: false, message: '', type: 'success' })

const knowledge = ref<KnowledgeItem[]>([])
const stats = ref({ total: 0, active: 0 })

const searchQuery = ref('')
const selectedCategory = ref('')
const activeFilter = ref('')

const form = ref({
  title: '',
  content: '',
  category: '',
  keywords: [] as string[],
  tags: [] as string[],
  priority: 5,
  is_active: true,
  media_url: '',
  media_type: ''
})

const keywordsInput = ref('')
const tagsInput = ref('')
const knowledgeFileInput = ref<HTMLInputElement | null>(null)
const uploadingKnowledgeFile = ref(false)
const uploadedFileName = ref('')

const categories = computed(() => {
  const cats = new Set<string>()
  knowledge.value.forEach(k => {
    if (k.category) cats.add(k.category)
  })
  return Array.from(cats)
})

const categoryCount = computed(() => {
  const counts: Record<string, number> = {}
  knowledge.value.forEach(k => {
    if (k.category) {
      counts[k.category] = (counts[k.category] || 0) + 1
    }
  })
  return counts
})

const totalUsage = computed(() => {
  return knowledge.value.reduce((sum, k) => sum + (k.usage_count || 0), 0)
})

const filteredKnowledge = computed(() => {
  return knowledge.value.filter(k => {
    const matchSearch = !searchQuery.value || 
      k.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      k.content.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchCategory = !selectedCategory.value || k.category === selectedCategory.value
    
    const matchActive = !activeFilter.value || 
      (activeFilter.value === 'true' ? k.is_active : !k.is_active)
    
    return matchSearch && matchCategory && matchActive
  })
})

const formatCategory = (cat: string) => {
  const map: Record<string, string> = {
    faq: '❓ FAQ',
    product: '📦 Produk',
    pricing: '💰 Harga',
    shipping: '🚚 Pengiriman',
    hours: '⏰ Jam Operasional',
    location: '📍 Lokasi',
    payment: '💳 Pembayaran',
    policy: '📋 Kebijakan'
  }
  return map[cat] || cat
}

const { fetch: apiFetch } = useApi()

const loadKnowledge = async () => {
  loading.value = true
  try {
    const res = await apiFetch<any>('/api/knowledge')

    knowledge.value = res.knowledge || []
    stats.value.total = res.total || 0
    stats.value.active = knowledge.value.filter(k => k.is_active).length
  } catch (err) {
    console.error('Failed to load knowledge:', err)
    showToast('Gagal memuat data', 'error')
  } finally {
    loading.value = false
  }
}

const editItem = (item: KnowledgeItem) => {
  editingItem.value = item
  form.value = {
    title: item.title,
    content: item.content,
    category: item.category || '',
    keywords: item.keywords || [],
    tags: item.tags || [],
    priority: item.priority || 5,
    tags: item.tags || [],
    priority: item.priority || 5,
    is_active: item.is_active,
    media_url: item.media_url || '',
    media_type: item.media_type || ''
  }
  keywordsInput.value = (item.keywords || []).join(', ')
  tagsInput.value = (item.tags || []).join(', ')
  showModal.value = true
}

const deleteItem = (item: KnowledgeItem) => {
  deleteConfirm.value = item
}

const confirmDelete = async () => {
  if (!deleteConfirm.value) return
  
  deleting.value = true
  try {
    await apiFetch(`/api/knowledge/${deleteConfirm.value.id}`, {
      method: 'DELETE'
    })
    
    showToast('Knowledge berhasil dihapus', 'success')
    deleteConfirm.value = null
    await loadKnowledge()
  } catch (err) {
    showToast('Gagal menghapus knowledge', 'error')
  } finally {
    deleting.value = false
  }
}

const saveKnowledge = async () => {
  if (!form.value.title || !form.value.content) return
  
  // Parse keywords and tags
  form.value.keywords = keywordsInput.value
    .split(',')
    .map(k => k.trim())
    .filter(k => k)
  
  form.value.tags = tagsInput.value
    .split(',')
    .map(t => t.trim())
    .filter(t => t)
  
  saving.value = true
  try {
    if (editingItem.value) {
      await apiFetch(`/api/knowledge/${editingItem.value.id}`, {
        method: 'PUT',
        body: JSON.stringify(form.value)
      })
      showToast('Knowledge berhasil diupdate', 'success')
    } else {
      await apiFetch('/api/knowledge', {
        method: 'POST',
        body: JSON.stringify(form.value)
      })
      showToast('Knowledge berhasil ditambahkan', 'success')
    }
    
    closeModal()
    await loadKnowledge()
  } catch (err: any) {
    showToast(err.data?.error || 'Gagal menyimpan knowledge', 'error')
  } finally {
    saving.value = false
  }
}

const closeModal = () => {
  showModal.value = false
  editingItem.value = null
  form.value = {
    title: '',
    content: '',
    category: '',
    keywords: [],
    tags: [],
    priority: 5,
    priority: 5,
    is_active: true,
    media_url: '',
    media_type: ''
  }
  keywordsInput.value = ''
  tagsInput.value = ''
}

const openKnowledgeFilePicker = () => {
  knowledgeFileInput.value?.click()
}

const handleKnowledgeFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  
  uploadingKnowledgeFile.value = true
  try {
    const formDataUpload = new FormData()
    formDataUpload.append('file', file)
    
    const response = await apiFetch<any>('/api/upload', {
      method: 'POST',
      body: formDataUpload
    })
    
    if (response.success) {
      form.value.media_url = response.file_url
      form.value.media_type = response.file_type
      uploadedFileName.value = response.file_name
    } else {
      throw new Error(response.error || 'Upload failed')
    }
  } catch (err: any) {
    showToast(err.data?.error || err.message || 'Gagal upload file', 'error')
  } finally {
    uploadingKnowledgeFile.value = false
    target.value = ''
  }
}

const removeMedia = () => {
  form.value.media_url = ''
  form.value.media_type = ''
  uploadedFileName.value = ''
}

const showToast = (message: string, type: 'success' | 'error') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

onMounted(() => {
  loadKnowledge()
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

