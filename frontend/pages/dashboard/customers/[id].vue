<template>
  <div class="space-y-6">
    <!-- Back Button & Header -->
    <div class="flex items-center gap-4">
      <button 
        @click="router.back()" 
        class="p-2 hover:bg-gray-100 rounded-lg transition"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
        </svg>
      </button>
      <h1 class="text-2xl font-bold text-gray-900">Detail Pelanggan</h1>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="animate-spin w-10 h-10 border-4 border-blue-600 border-t-transparent rounded-full"></div>
    </div>

    <!-- Content -->
    <div v-else-if="customer" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column - Customer Info -->
      <div class="lg:col-span-1 space-y-6">
        <!-- Profile Card -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <div class="text-center">
            <div class="w-24 h-24 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center mx-auto">
              <span class="text-white font-bold text-3xl">{{ getInitials(customerHelper.getDisplayName(customer)) }}</span>
            </div>
            <div class="flex items-center justify-center gap-2 mt-4">
              <h2 class="text-xl font-bold text-gray-900">{{ customerHelper.getDisplayName(customer) }}</h2>
              <button 
                @click="openEditNameModal"
                class="p-1 hover:bg-gray-100 rounded-lg transition"
                title="Edit Nama"
              >
                <svg class="w-4 h-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                </svg>
              </button>
            </div>
            <p class="text-gray-500">{{ customerHelper.formatPhoneFromJID(customer.customer_jid) }}</p>
            <span :class="customerHelper.getStatusColor(customer.status)" class="inline-block mt-3 px-3 py-1 text-sm font-medium rounded-full">
              {{ customerHelper.getStatusLabel(customer.status) }}
            </span>
          </div>

          <div class="mt-6 pt-6 border-t border-gray-100 space-y-4">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Total Pesan</span>
              <span class="font-semibold text-gray-900">{{ customer.message_count }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Terakhir Aktif</span>
              <span class="font-semibold text-gray-900">{{ customerHelper.formatRelativeTime(customer.last_message_at) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">Bergabung</span>
              <span class="font-semibold text-gray-900">{{ formatDate(customer.created_at) }}</span>
            </div>
          </div>
        </div>

        <!-- Lead Score Card -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">Lead Score</h3>
          <div class="flex items-center gap-4">
            <div class="w-20 h-20 rounded-full flex items-center justify-center" :class="getLeadScoreColor(leadScore)">
              <span class="text-2xl font-bold text-white">{{ leadScore }}</span>
            </div>
            <div class="flex-1">
              <input 
                type="range" 
                v-model.number="leadScore" 
                min="0" 
                max="100" 
                class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                @change="updateLeadScore"
              >
              <div class="flex justify-between text-xs text-gray-500 mt-1">
                <span>Cold</span>
                <span>Warm</span>
                <span>Hot</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Tags Card -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-gray-900">Tags</h3>
            <button 
              @click="showTagModal = true"
              class="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              + Tambah Tag
            </button>
          </div>
          
          <div v-if="customerTags.length === 0" class="text-center py-4 text-gray-500 text-sm">
            Belum ada tag
          </div>
          
          <div v-else class="flex flex-wrap gap-2">
            <span 
              v-for="tag in customerTags" 
              :key="tag.id"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium text-white"
              :style="{ backgroundColor: tag.color }"
            >
              {{ tag.name }}
              <button 
                @click="removeTag(tag.id)"
                class="hover:opacity-70 transition"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </span>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">Aksi Cepat</h3>
          <div class="space-y-3">
            <button 
              @click="toggleFollowUp"
              :class="customer.needs_follow_up ? 'bg-yellow-100 text-yellow-700' : 'bg-gray-100 text-gray-700'"
              class="w-full px-4 py-2.5 rounded-lg font-medium text-sm hover:opacity-80 transition flex items-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
              </svg>
              {{ customer.needs_follow_up ? 'Hapus dari Follow Up' : 'Tandai untuk Follow Up' }}
            </button>
            
            <select 
              v-model="selectedStatus" 
              @change="updateStatus"
              class="w-full px-4 py-2.5 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
            >
              <option value="" disabled>Ubah Status</option>
              <option value="new">Baru</option>
              <option value="hot_lead">Hot Lead</option>
              <option value="warm_lead">Warm Lead</option>
              <option value="cold_lead">Cold Lead</option>
              <option value="customer">Pelanggan</option>
              <option value="complaint">Komplain</option>
              <option value="spam">Spam</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Right Column - Chat & Notes -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Real-time Chat -->
        <div 
          :class="[
            'bg-white rounded-xl shadow-sm border border-gray-100 flex flex-col transition-all duration-300',
            isFullscreen ? 'fixed inset-0 z-50 h-screen rounded-none' : ''
          ]"
          :style="isFullscreen ? '' : 'height: 500px;'"
        >
          <!-- Chat Header -->
          <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-green-500 flex items-center justify-center">
                <svg class="w-5 h-5 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/>
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-gray-900">Chat WhatsApp</h3>
                <p class="text-xs text-gray-500">{{ customerHelper.formatPhoneFromJID(customer.customer_jid) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-1">
              <button 
                @click="loadMessages" 
                class="p-2 hover:bg-gray-100 rounded-lg transition"
                title="Refresh"
              >
                <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                </svg>
              </button>
              <button 
                @click="toggleFullscreen" 
                class="p-2 hover:bg-blue-50 rounded-lg transition"
                :title="isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'"
              >
                <svg v-if="!isFullscreen" class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"/>
                </svg>
                <svg v-else class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25"/>
                </svg>
              </button>
              <button 
                @click="showClearChatConfirm = true" 
                class="p-2 hover:bg-red-50 rounded-lg transition"
                title="Clear Chat"
              >
                <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Messages Container -->
          <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
            <div v-if="messages.length === 0" class="flex items-center justify-center h-full text-gray-500">
              Belum ada pesan
            </div>
            
            <div 
              v-for="msg in messages" 
              :key="msg.id"
              :class="[
                'max-w-[75%] rounded-2xl px-4 py-2.5 shadow-sm',
                msg.is_from_me 
                  ? 'bg-green-500 text-white ml-auto rounded-br-md' 
                  : 'bg-white text-gray-800 rounded-bl-md'
              ]"
            >
              <!-- Media Preview -->
              <div v-if="msg.media_url || msg.message_type === 'image'" class="mb-2">
                <img 
                  v-if="msg.message_type === 'image'" 
                  :src="getMediaUrl(msg.media_url)"
                  alt="Image"
                  class="max-w-full rounded-lg cursor-pointer hover:opacity-90"
                  style="max-height: 200px;"
                  @click="openMediaPreview(msg.media_url)"
                />
                <a 
                  v-else-if="msg.message_type === 'document'"
                  :href="getMediaUrl(msg.media_url)"
                  target="_blank"
                  class="flex items-center gap-2 p-2 bg-white/20 rounded-lg hover:bg-white/30"
                >
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                  <span class="text-sm">📄 Dokumen</span>
                </a>
              </div>
              <!-- Text Content -->
              <p v-if="msg.message_text || msg.body" class="text-sm whitespace-pre-wrap break-words">{{ msg.message_text || msg.body }}</p>
              <p v-else-if="msg.message_type === 'image'" class="text-sm italic opacity-75">📷 Gambar</p>
              <p v-else-if="msg.message_type === 'document'" class="text-sm italic opacity-75">📄 Dokumen</p>
              <p :class="[
                'text-xs mt-1 text-right',
                msg.is_from_me ? 'text-green-100' : 'text-gray-400'
              ]">
                {{ formatMessageTimeAuto(msg.timestamp) }}
              </p>
            </div>
          </div>

          <!-- Message Input -->
          <div class="p-4 border-t border-gray-100 bg-white rounded-b-xl">
            <div class="flex items-end gap-3">
              <!-- File Attachment Button -->
              <button 
                @click="openFilePicker"
                :disabled="sendingMessage"
                class="p-3 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-xl transition disabled:opacity-50 disabled:cursor-not-allowed"
                title="Kirim File/Gambar"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"/>
                </svg>
              </button>
              <input 
                ref="fileInput"
                type="file"
                accept="image/*,.pdf,.doc,.docx,.xls,.xlsx"
                class="hidden"
                @change="handleFileSelect"
              />
              <!-- Emoji Picker -->
              <div class="relative">
                <button 
                  type="button"
                  @click="showEmojiPicker = !showEmojiPicker"
                  class="p-3 text-gray-500 hover:text-yellow-500 hover:bg-gray-100 rounded-xl transition"
                  title="Emoji"
                >
                  <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-5-6c.78 2.34 2.72 4 5 4s4.22-1.66 5-4H7zm8-4c.55 0 1-.45 1-1s-.45-1-1-1-1 .45-1 1 .45 1 1 1zm-6 0c.55 0 1-.45 1-1s-.45-1-1-1-1 .45-1 1 .45 1 1 1z"/>
                  </svg>
                </button>
                <!-- Emoji Dropdown -->
                <div 
                  v-if="showEmojiPicker" 
                  class="absolute bottom-12 left-0 bg-white rounded-xl shadow-lg border border-gray-200 p-3 grid grid-cols-8 gap-1 z-50"
                  style="width: 280px;"
                >
                  <button 
                    v-for="emoji in emojis" 
                    :key="emoji"
                    @click="insertEmoji(emoji)"
                    class="text-xl p-1 hover:bg-gray-100 rounded transition"
                  >
                    {{ emoji }}
                  </button>
                </div>
              </div>
              <textarea 
                v-model="newMessage" 
                @keydown.enter.exact.prevent="sendMessage"
                @focus="showEmojiPicker = false"
                placeholder="Ketik pesan..."
                class="flex-1 px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent resize-none"
                rows="1"
                :disabled="sendingMessage"
              ></textarea>
              <button 
                @click="sendMessage"
                :disabled="!newMessage.trim() || sendingMessage"
                class="p-3 bg-green-500 text-white rounded-xl hover:bg-green-600 transition disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg v-if="sendingMessage" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
                </svg>
              </button>
            </div>
            <!-- Upload Progress -->
            <div v-if="uploadingFile" class="mt-2 flex items-center gap-2 text-sm text-gray-500">
              <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>Mengunggah file...</span>
            </div>
          </div>
        </div>

        <!-- Notes Section -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-gray-900">Catatan Internal</h3>
          </div>

          <!-- Add Note Form -->
          <div class="flex gap-3 mb-6">
            <textarea 
              v-model="newNote" 
              placeholder="Tambahkan catatan tentang pelanggan ini..."
              class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 resize-none"
              rows="2"
            ></textarea>
            <button 
              @click="addNote"
              :disabled="!newNote.trim() || addingNote"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed self-end"
            >
              <span v-if="addingNote">...</span>
              <span v-else>Simpan</span>
            </button>
          </div>

          <!-- Notes List -->
          <div v-if="notes.length === 0" class="text-center py-8 text-gray-500">
            Belum ada catatan
          </div>
          
          <div v-else class="space-y-4 max-h-64 overflow-y-auto">
            <div 
              v-for="note in notes" 
              :key="note.id"
              class="p-4 bg-gray-50 rounded-lg"
            >
              <p class="text-gray-800 whitespace-pre-wrap">{{ note.content }}</p>
              <p class="text-xs text-gray-500 mt-2">{{ formatDateTime(note.created_at) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tag Modal -->
    <Teleport to="body">
      <div v-if="showTagModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showTagModal = false">
        <div class="bg-white rounded-xl shadow-xl max-w-md w-full mx-4 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Tambah Tag</h3>
          
          <div v-if="availableTags.length === 0" class="text-center py-6 text-gray-500">
            Semua tag sudah diterapkan
          </div>
          
          <div v-else class="space-y-2 max-h-64 overflow-y-auto">
            <button 
              v-for="tag in availableTags" 
              :key="tag.id"
              @click="addTag(tag.id)"
              class="w-full flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 transition text-left"
            >
              <span 
                class="w-4 h-4 rounded-full" 
                :style="{ backgroundColor: tag.color }"
              ></span>
              <span class="font-medium text-gray-800">{{ tag.name }}</span>
            </button>
          </div>
          
          <div class="mt-6 pt-4 border-t border-gray-100 flex justify-between items-center">
            <NuxtLink 
              to="/dashboard/customers/segments"
              class="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              Kelola Tags
            </NuxtLink>
            <button 
              @click="showTagModal = false"
              class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Clear Chat Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showClearChatConfirm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showClearChatConfirm = false">
        <div class="bg-white rounded-xl shadow-xl max-w-sm w-full mx-4 p-6">
          <div class="text-center">
            <div class="w-12 h-12 rounded-full bg-red-100 flex items-center justify-center mx-auto mb-4">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Hapus Semua Chat?</h3>
            <p class="text-gray-500 text-sm mb-6">
              Semua pesan dengan pelanggan ini akan dihapus secara permanen. Tindakan ini tidak dapat dibatalkan.
            </p>
            <div class="flex gap-3">
              <button 
                @click="showClearChatConfirm = false"
                class="flex-1 px-4 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
              >
                Batal
              </button>
              <button 
                @click="clearChat"
                :disabled="clearingChat"
                class="flex-1 px-4 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition font-medium disabled:opacity-50"
              >
                <span v-if="clearingChat">Menghapus...</span>
                <span v-else>Ya, Hapus</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Edit Name Modal -->
    <Teleport to="body">
      <div v-if="showEditNameModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showEditNameModal = false">
        <div class="bg-white rounded-xl shadow-xl max-w-sm w-full mx-4 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Edit Nama Pelanggan</h3>
          
          <input 
            v-model="editingName"
            type="text"
            placeholder="Masukkan nama pelanggan..."
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            @keydown.enter="saveName"
          />
          
          <p class="text-xs text-gray-500 mt-2">
            Kosongkan untuk menggunakan nomor telepon sebagai nama.
          </p>
          
          <div class="flex gap-3 mt-6">
            <button 
              @click="showEditNameModal = false"
              class="flex-1 px-4 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
            >
              Batal
            </button>
            <button 
              @click="saveName"
              :disabled="savingName"
              class="flex-1 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-medium disabled:opacity-50"
            >
              <span v-if="savingName">Menyimpan...</span>
              <span v-else>Simpan</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import type { Customer } from '~/composables/useCustomers'
import { useRealtimeUpdates, type NewMessageData } from '~/composables/useRealtimeUpdates'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const customerHelper = useCustomers()
const { fetch: apiFetch } = useApi()

const customerId = computed(() => route.params.id as string)

// State
const loading = ref(true)
const customer = ref<Customer | null>(null)
const customerTags = ref<any[]>([])
const allTags = ref<any[]>([])
const notes = ref<any[]>([])
const messages = ref<any[]>([])
const leadScore = ref(0)
const selectedStatus = ref('')
const newNote = ref('')
const addingNote = ref(false)
const showTagModal = ref(false)
const newMessage = ref('')
const sendingMessage = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const showClearChatConfirm = ref(false)
const clearingChat = ref(false)
const showEditNameModal = ref(false)
const editingName = ref('')
const savingName = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const uploadingFile = ref(false)
const showEmojiPicker = ref(false)
const isFullscreen = ref(false)

// Common WhatsApp emojis
const emojis = [
  '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂',
  '🙂', '😉', '😊', '😇', '🥰', '😍', '🤩', '😘',
  '😗', '😚', '😙', '🥲', '😋', '😛', '😜', '🤪',
  '😝', '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨',
  '😐', '😑', '😶', '😏', '😒', '🙄', '😬', '🤥',
  '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕',
  '👍', '👎', '👋', '🙏', '💪', '❤️', '🔥', '✨',
  '🎉', '🎊', '💯', '✅', '❌', '⭐', '🌟', '💫'
]

// Real-time updates
useRealtimeUpdates({
  onNewMessage: (data: NewMessageData) => {
    console.log('[WS] New message received:', data)
    if (!customer.value) {
      console.log('[WS] No customer loaded, ignoring')
      return
    }
    
    console.log('[WS] Customer JID:', customer.value.customer_jid)
    console.log('[WS] Sender JID:', data.sender_jid)
    console.log('[WS] Chat JID:', data.chat_jid)
    
    // Check if message is related to current customer
    // - Incoming: sender_jid = customer JID, is_from_me = false
    // - Outgoing: chat_jid = customer JID, is_from_me = true
    const isFromCustomer = data.sender_jid === customer.value.customer_jid
    const isToCustomer = data.chat_jid === customer.value.customer_jid
    
    console.log('[WS] isFromCustomer:', isFromCustomer, ', isToCustomer:', isToCustomer)
    
    if (isFromCustomer || isToCustomer) {
      // Avoid duplicates
      const exists = messages.value.some(m => m.id === data.message_id)
      if (exists) {
        console.log('[WS] Message already exists, skipping')
        return
      }
      
      console.log('[WS] Adding message with media_url:', data.media_url)
      // Use is_from_me directly from backend - don't override it
      messages.value.push({
        id: data.message_id,
        message_text: data.message_text,
        message_type: data.message_type || 'text',
        media_url: data.media_url || '',
        is_from_me: data.is_from_me, // Trust backend value
        timestamp: new Date(data.timestamp * 1000).toISOString()
      })
      scrollToBottom()
    } else {
      console.log('[WS] Message not for current customer, ignoring')
    }
  }
})

// Computed
const availableTags = computed(() => {
  const assignedIds = new Set(customerTags.value.map(t => t.id))
  return allTags.value.filter(t => !assignedIds.has(t.id))
})

// Scroll to bottom of messages
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// Load data
const loadCustomer = async () => {
  try {
    // Get customer detail with messages in one call
    const data = await apiFetch(`/api/customers/${customerId.value}`)
    customer.value = data.customer
    messages.value = (data.messages || []).reverse() // Reverse to show oldest first
    selectedStatus.value = customer.value?.status || ''
    leadScore.value = customer.value?.lead_score || 0
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load customer:', error)
  }
}

const loadTags = async () => {
  try {
    const [allTagsData, customerTagsData] = await Promise.all([
      apiFetch('/api/tags'),
      apiFetch(`/api/customers/${customerId.value}/tags`)
    ])
    allTags.value = allTagsData || []
    customerTags.value = customerTagsData || []
  } catch (error) {
    console.error('Failed to load tags:', error)
  }
}

const loadNotes = async () => {
  try {
    const data = await apiFetch(`/api/customers/${customerId.value}/notes`)
    notes.value = data || []
  } catch (error) {
    console.error('Failed to load notes:', error)
    notes.value = []
  }
}

const loadMessages = async () => {
  try {
    // Messages are included in customer detail response
    const data = await apiFetch(`/api/customers/${customerId.value}`)
    messages.value = (data.messages || []).reverse() // Reverse to show oldest first
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load messages:', error)
  }
}

const loadAll = async () => {
  loading.value = true
  await loadCustomer() // This also loads messages
  await Promise.all([loadTags(), loadNotes()])
  loading.value = false
}

// Insert emoji into message
const insertEmoji = (emoji: string) => {
  newMessage.value += emoji
  showEmojiPicker.value = false
}

// Toggle fullscreen mode
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// Send message
const sendMessage = async () => {
  if (!newMessage.value.trim() || !customer.value || sendingMessage.value) return
  
  sendingMessage.value = true
  const messageText = newMessage.value.trim()
  newMessage.value = ''
  
  try {
    await apiFetch('/api/whatsapp/send', {
      method: 'POST',
      body: {
        recipient_jid: customer.value.customer_jid,
        message: messageText
      }
    })
    // Message will be added via WebSocket broadcast
    // Don't add locally to avoid duplicates
  } catch (error) {
    console.error('Failed to send message:', error)
    newMessage.value = messageText // Restore message on error
  } finally {
    sendingMessage.value = false
  }
}

// File upload
const openFilePicker = () => {
  fileInput.value?.click()
}

const handleFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file || !customer.value) return
  
  uploadingFile.value = true
  try {
    // First upload the file
    const formData = new FormData()
    formData.append('file', file)
    
    const config = useRuntimeConfig()
    const apiBase = config.public.apiBase || 'http://localhost:8080'
    const authStore = useAuthStore()
    
    const uploadResponse = await fetch(`${apiBase}/api/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: formData
    })
    
    const uploadResult = await uploadResponse.json()
    if (!uploadResult.success) {
      throw new Error(uploadResult.error || 'Upload failed')
    }
    
    // Then send as media message
    const mediaFormData = new FormData()
    mediaFormData.append('file', file)
    mediaFormData.append('recipient_jid', customer.value.customer_jid)
    mediaFormData.append('caption', file.name)
    
    const sendResponse = await fetch(`${apiBase}/api/whatsapp/send/media`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: mediaFormData
    })
    
    const sendResult = await sendResponse.json()
    if (!sendResult.success) {
      throw new Error(sendResult.error || 'Failed to send media')
    }
    
    // Clear input
    target.value = ''
  } catch (error: any) {
    console.error('Failed to send file:', error)
    alert('Gagal mengirim file: ' + (error?.message || 'Unknown error'))
  } finally {
    uploadingFile.value = false
  }
}

// Clear chat
const clearChat = async () => {
  if (!customer.value || clearingChat.value) return
  
  clearingChat.value = true
  const jid = customer.value.customer_jid
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080'
  const url = `${apiBase}/api/whatsapp/messages/${encodeURIComponent(jid)}`
  console.log('[Clear Chat] Clearing messages for:', jid)
  console.log('[Clear Chat] URL:', url)
  
  try {
    const authStore = useAuthStore()
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      }
    })
    
    const result = await response.json()
    console.log('[Clear Chat] Response status:', response.status)
    console.log('[Clear Chat] Result:', result)
    
    if (response.ok && result.success) {
      messages.value = []
      showClearChatConfirm.value = false
    } else {
      throw new Error(result.error || 'Failed to clear chat')
    }
  } catch (error: any) {
    console.error('[Clear Chat] Error:', error)
    alert('Gagal menghapus chat: ' + (error?.message || 'Unknown error'))
  } finally {
    clearingChat.value = false
  }
}

// Edit name
const openEditNameModal = () => {
  editingName.value = customer.value?.customer_name || ''
  showEditNameModal.value = true
}

const saveName = async () => {
  if (!customer.value || savingName.value) return
  
  savingName.value = true
  try {
    await customerHelper.updateCustomer(customerId.value, { 
      customer_name: editingName.value.trim() || null 
    })
    customer.value.customer_name = editingName.value.trim() || null
    showEditNameModal.value = false
  } catch (error) {
    console.error('Failed to save name:', error)
    alert('Gagal menyimpan nama')
  } finally {
    savingName.value = false
  }
}

// Actions
const addTag = async (tagId: string) => {
  try {
    await apiFetch(`/api/customers/${customerId.value}/tags`, { 
      method: 'POST',
      body: { tag_id: tagId }
    })
    await loadTags()
    showTagModal.value = false
  } catch (error) {
    console.error('Failed to add tag:', error)
  }
}

const removeTag = async (tagId: string) => {
  try {
    await apiFetch(`/api/customers/${customerId.value}/tags/${tagId}`, { method: 'DELETE' })
    await loadTags()
  } catch (error) {
    console.error('Failed to remove tag:', error)
  }
}

const addNote = async () => {
  if (!newNote.value.trim()) return
  addingNote.value = true
  try {
    await apiFetch(`/api/customers/${customerId.value}/notes`, {
      method: 'POST',
      body: { content: newNote.value }
    })
    newNote.value = ''
    await loadNotes()
  } catch (error) {
    console.error('Failed to add note:', error)
  } finally {
    addingNote.value = false
  }
}

const updateLeadScore = async () => {
  try {
    // Convert to number explicitly (slider might return string)
    const score = Number(leadScore.value)
    console.log('[Lead Score] Saving:', score)
    
    await apiFetch(`/api/customers/${customerId.value}/lead-score`, {
      method: 'PUT',
      body: { lead_score: score }
    })
    console.log('[Lead Score] Saved successfully')
  } catch (error) {
    console.error('[Lead Score] Failed to update:', error)
    alert('Gagal menyimpan lead score')
  }
}

const updateStatus = async () => {
  if (!selectedStatus.value || !customer.value) return
  try {
    await customerHelper.updateCustomer(customerId.value, { status: selectedStatus.value })
    customer.value.status = selectedStatus.value
  } catch (error) {
    console.error('Failed to update status:', error)
  }
}

const toggleFollowUp = async () => {
  if (!customer.value) return
  try {
    await customerHelper.updateCustomer(customerId.value, { 
      needs_follow_up: !customer.value.needs_follow_up 
    })
    customer.value.needs_follow_up = !customer.value.needs_follow_up
  } catch (error) {
    console.error('Failed to toggle follow up:', error)
  }
}

// Helper for media URLs
const getMediaUrl = (url: string) => {
  if (!url) return ''
  // If it's already a full URL, return as-is
  if (url.startsWith('http')) return url
  // If it's a relative path, prepend API base
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080'
  return `${apiBase}${url}`
}

const openMediaPreview = (url: string) => {
  const fullUrl = getMediaUrl(url)
  if (fullUrl) {
    window.open(fullUrl, '_blank')
  }
}

// Helpers
const getInitials = (name: string): string => {
  if (!name) return '?'
  const parts = name.trim().split(/[\s@]+/)
  if (parts.length === 1) {
    if (/^\d+$/.test(parts[0])) return parts[0].slice(-2)
    return parts[0].substring(0, 2).toUpperCase()
  }
  return (parts[0][0] + (parts[1]?.[0] || '')).toUpperCase()
}

const getLeadScoreColor = (score: number): string => {
  if (score >= 70) return 'bg-red-500'
  if (score >= 40) return 'bg-orange-500'
  return 'bg-blue-500'
}

const formatDate = (date: string): string => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { 
    day: 'numeric', month: 'short', year: 'numeric' 
  })
}

const formatDateTime = (date: string): string => {
  if (!date) return '-'
  return new Date(date).toLocaleString('id-ID', {
    day: 'numeric', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })
}

const formatMessageTime = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString('id-ID', {
    day: 'numeric', month: 'short',
    hour: '2-digit', minute: '2-digit'
  })
}

// Handle both unix timestamp (number) and ISO date string
const formatMessageTimeAuto = (timestamp: number | string): string => {
  let date: Date
  if (typeof timestamp === 'number') {
    date = new Date(timestamp * 1000)
  } else if (typeof timestamp === 'string') {
    date = new Date(timestamp)
  } else {
    return '-'
  }
  return date.toLocaleString('id-ID', {
    day: 'numeric', month: 'short',
    hour: '2-digit', minute: '2-digit'
  })
}

// Initialize
onMounted(() => {
  loadAll()
})
</script>

