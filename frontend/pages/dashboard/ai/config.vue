<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">AI Auto-Reply Configuration</h1>
        <p class="text-gray-600 mt-1">Konfigurasi AI untuk balasan otomatis customer</p>
      </div>
      <div class="flex gap-3">
        <NuxtLink to="/dashboard/ai/test" class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          Test AI
        </NuxtLink>
        <NuxtLink to="/dashboard/ai/knowledge" class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          Knowledge Base
        </NuxtLink>
      </div>
    </div>

    <!-- Status Card -->
    <div :class="[
      'rounded-2xl p-6 text-white relative overflow-hidden',
      config.enabled ? 'bg-gradient-to-r from-emerald-600 to-teal-600' : 'bg-gradient-to-r from-slate-600 to-slate-700'
    ]">
      <div class="absolute top-0 right-0 opacity-10">
        <svg class="w-64 h-64 transform translate-x-16 -translate-y-8" fill="currentColor" viewBox="0 0 24 24">
          <path d="M21 11.5a8.38 8.38 0 01-.9 3.8 8.5 8.5 0 01-7.6 4.7 8.38 8.38 0 01-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 01-.9-3.8 8.5 8.5 0 014.7-7.6 8.38 8.38 0 013.8-.9h.5a8.48 8.48 0 018 8v.5z"/>
        </svg>
      </div>
      <div class="relative z-10">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold mb-2">Status AI Auto-Reply</h2>
            <p class="text-white/80 text-sm">{{ selectedProviderInfo?.name || 'No Provider Selected' }}</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="config.enabled" @change="toggleAI" class="sr-only peer">
            <div class="w-14 h-7 bg-white/30 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[4px] after:bg-white after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-green-400"></div>
          </label>
        </div>
        <div class="mt-4 flex items-center gap-6 flex-wrap">
          <div>
            <p class="text-white/70 text-xs uppercase tracking-wide">Status</p>
            <p class="text-lg font-semibold">{{ config.enabled ? 'Aktif' : 'Nonaktif' }}</p>
          </div>
          <div>
            <p class="text-white/70 text-xs uppercase tracking-wide">Provider</p>
            <p class="text-lg font-semibold">{{ selectedProviderInfo?.name || '-' }}</p>
          </div>
          <div>
            <p class="text-white/70 text-xs uppercase tracking-wide">Model</p>
            <p class="text-lg font-semibold">{{ config.model || '-' }}</p>
          </div>
          <div>
            <p class="text-white/70 text-xs uppercase tracking-wide">Confidence</p>
            <p class="text-lg font-semibold">{{ (config.confidence_threshold * 100).toFixed(0) }}%</p>
          </div>
          <div v-if="config.api_key_set">
            <p class="text-white/70 text-xs uppercase tracking-wide">API Key</p>
            <p class="text-lg font-semibold flex items-center gap-1">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
              Configured
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Configuration Tabs -->
    <div class="bg-white rounded-xl border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-100">
        <nav class="flex -mb-px overflow-x-auto">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'px-6 py-4 text-sm font-medium border-b-2 transition-colors whitespace-nowrap',
              activeTab === tab.id 
                ? 'border-indigo-500 text-indigo-600' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <!-- Provider & API Key Tab -->
      <div v-show="activeTab === 'provider'" class="p-6 space-y-6">
        <!-- Provider Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">Pilih AI Provider</label>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <button
              v-for="provider in providers"
              :key="provider.id"
              @click="selectProvider(provider.id)"
              :class="[
                'p-4 rounded-xl border-2 text-left transition-all',
                config.ai_provider === provider.id 
                  ? 'border-indigo-500 bg-indigo-50' 
                  : 'border-gray-200 hover:border-gray-300'
              ]"
            >
              <div class="flex items-center gap-3 mb-2">
                <div :class="[
                  'w-10 h-10 rounded-lg flex items-center justify-center font-bold text-white',
                  getProviderColor(provider.id)
                ]">
                  {{ provider.name.charAt(0) }}
                </div>
                <div>
                  <p class="font-semibold text-gray-900">{{ provider.name }}</p>
                  <p v-if="provider.free_available" class="text-xs text-green-600">Free tier available</p>
                  <p v-else class="text-xs text-gray-500">Paid only</p>
                </div>
              </div>
              <p class="text-xs text-gray-600">{{ provider.description }}</p>
            </button>
          </div>
        </div>

        <!-- API Key Input -->
        <div class="p-4 bg-gray-50 rounded-lg border border-gray-200">
          <div class="flex items-start gap-3 mb-4">
            <svg class="w-5 h-5 text-gray-500 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 8a6 6 0 01-7.743 5.743L10 14l-1 1-1 1H6v2H2v-4l4.257-4.257A6 6 0 1118 8zm-6-4a1 1 0 100 2 2 2 0 012 2 1 1 0 102 0 4 4 0 00-4-4z" clip-rule="evenodd" />
            </svg>
            <div class="flex-1">
              <h3 class="font-medium text-gray-900">API Key Configuration</h3>
              <p class="text-sm text-gray-600 mt-1">
                Masukkan API key dari provider yang dipilih. API key akan dienkripsi dan disimpan dengan aman.
              </p>
            </div>
          </div>

          <div class="space-y-4">
            <!-- Use System Key Option -->
            <label class="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200 cursor-pointer hover:bg-gray-50">
              <input type="checkbox" v-model="config.use_system_key" @change="onUseSystemKeyChange" class="w-4 h-4 text-indigo-600 rounded">
              <div>
                <p class="font-medium text-gray-900">Gunakan System API Key</p>
                <p class="text-sm text-gray-500">Gunakan API key default sistem (jika tersedia)</p>
              </div>
            </label>

            <!-- Custom API Key Input -->
            <div v-if="!config.use_system_key">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ selectedProviderInfo?.name }} API Key
              </label>
              <div class="flex gap-3">
                <div class="flex-1 relative">
                  <input
                    :type="showApiKey ? 'text' : 'password'"
                    v-model="apiKeyInput"
                    :placeholder="config.api_key_set ? '••••••••••••••••••••••••' : `Masukkan ${selectedProviderInfo?.name || 'API'} key`"
                    class="w-full px-4 py-3 pr-12 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-mono text-sm"
                  >
                  <button
                    type="button"
                    @click="showApiKey = !showApiKey"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  >
                    <svg v-if="showApiKey" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                </div>
                <button
                  @click="testConnection"
                  :disabled="testingConnection || !apiKeyInput"
                  class="px-4 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="testingConnection" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  Test
                </button>
              </div>
              
              <!-- Connection Test Result -->
              <div v-if="connectionTestResult" :class="[
                'mt-3 p-3 rounded-lg flex items-center gap-2',
                connectionTestResult.success ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
              ]">
                <svg v-if="connectionTestResult.success" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
                {{ connectionTestResult.message }}
              </div>

              <!-- Get API Key Link -->
              <p class="mt-3 text-sm text-gray-500">
                Belum punya API key? Dapatkan di 
                <a :href="getProviderLink(config.ai_provider)" target="_blank" class="text-indigo-600 hover:underline">
                  {{ selectedProviderInfo?.name }} Console →
                </a>
              </p>
            </div>
          </div>
        </div>

        <!-- Model Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Pilih Model</label>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <button
              v-for="model in currentModels"
              :key="model.id"
              @click="config.model = model.id"
              :class="[
                'p-4 rounded-lg border-2 text-left transition-all',
                config.model === model.id 
                  ? 'border-indigo-500 bg-indigo-50' 
                  : 'border-gray-200 hover:border-gray-300'
              ]"
            >
              <p class="font-semibold text-gray-900">{{ model.name }}</p>
              <p class="text-xs text-gray-600 mt-1">{{ model.description }}</p>
              <p class="text-xs text-indigo-600 mt-2">
                ${{ model.input_cost }}/1M input • ${{ model.output_cost }}/1M output
              </p>
            </button>
          </div>
        </div>
      </div>

      <!-- General Settings Tab -->
      <div v-show="activeTab === 'general'" class="p-6 space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Confidence Threshold -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Confidence Threshold: {{ (config.confidence_threshold * 100).toFixed(0) }}%
            </label>
            <input
              type="range"
              v-model.number="config.confidence_threshold"
              min="0.5"
              max="0.95"
              step="0.05"
              class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
            >
            <div class="flex justify-between text-xs text-gray-500 mt-1">
              <span>50% (Auto-reply lebih sering)</span>
              <span>95% (Lebih banyak eskalasi)</span>
            </div>
          </div>

          <!-- Max Tokens -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Max Response Tokens</label>
            <input
              type="number"
              v-model.number="config.max_tokens"
              min="50"
              max="2000"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
            <p class="mt-1 text-xs text-gray-500">Jumlah token maksimum untuk respons (50-2000)</p>
          </div>

          <!-- Response Language -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Bahasa Respons</label>
            <select
              v-model="config.language"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
              <option value="id">Bahasa Indonesia</option>
              <option value="en">English</option>
              <option value="auto">Auto-detect</option>
            </select>
          </div>
        </div>

        <!-- Auto Escalation Info -->
        <div class="p-4 bg-amber-50 border border-amber-200 rounded-lg">
          <div class="flex items-start gap-3">
            <svg class="w-5 h-5 text-amber-600 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <div>
              <h3 class="font-medium text-amber-800">Auto Escalation</h3>
              <p class="text-sm text-amber-700 mt-1">
                Ketika confidence AI di bawah threshold ({{ (config.confidence_threshold * 100).toFixed(0) }}%), 
                pesan akan di-eskalasi ke Anda untuk ditangani secara manual.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- System Prompt Tab -->
      <div v-show="activeTab === 'prompt'" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">System Prompt</label>
          <textarea
            v-model="config.system_prompt"
            rows="8"
            placeholder="Masukkan instruksi untuk AI..."
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-mono text-sm resize-none"
          ></textarea>
          <p class="mt-2 text-xs text-gray-500">
            Prompt ini akan menjadi instruksi dasar untuk AI dalam menjawab customer.
          </p>
        </div>

        <!-- Preset Prompts -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Preset Prompt</label>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <button
              v-for="preset in promptPresets"
              :key="preset.id"
              @click="applyPreset(preset)"
              class="p-3 text-left border border-gray-200 rounded-lg hover:border-indigo-300 hover:bg-indigo-50 transition-colors"
            >
              <p class="font-medium text-gray-900">{{ preset.name }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ preset.description }}</p>
            </button>
          </div>
        </div>
      </div>

      <!-- Business Context Tab -->
      <div v-show="activeTab === 'context'" class="p-6 space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Business Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Nama Bisnis</label>
            <input
              type="text"
              v-model="config.business_name"
              placeholder="Toko ABC"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
          </div>

          <!-- Business Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Jenis Bisnis</label>
            <select
              v-model="config.business_type"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
              <option value="retail">Retail / Toko</option>
              <option value="fnb">F&B / Makanan</option>
              <option value="service">Jasa / Service</option>
              <option value="fashion">Fashion</option>
              <option value="handcraft">Kerajinan / Handcraft</option>
              <option value="other">Lainnya</option>
            </select>
          </div>
        </div>

        <!-- Business Hours -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Jam Operasional</label>
          <input
            type="text"
            v-model="config.business_hours"
            placeholder="Senin-Jumat 09:00-17:00, Sabtu 09:00-14:00"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          >
        </div>

        <!-- Business Description -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Deskripsi Bisnis</label>
          <textarea
            v-model="config.business_description"
            rows="4"
            placeholder="Deskripsi singkat tentang bisnis Anda, produk/jasa yang ditawarkan..."
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
          ></textarea>
        </div>

        <!-- Contact Info -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Alamat</label>
            <input
              type="text"
              v-model="config.business_address"
              placeholder="Jl. Contoh No. 123, Makassar"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Metode Pembayaran</label>
            <input
              type="text"
              v-model="config.payment_methods"
              placeholder="Transfer BCA, GoPay, Cash"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
          </div>
        </div>
      </div>

      <!-- Escalation Tab -->
      <div v-show="activeTab === 'escalation'" class="p-6 space-y-4">
        <div>
          <h3 class="font-medium text-gray-900 mb-4">Kondisi Eskalasi Otomatis</h3>
          <div class="space-y-3">
            <label class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100">
              <input type="checkbox" v-model="config.escalate_low_confidence" class="w-4 h-4 text-indigo-600 rounded">
              <div>
                <p class="font-medium text-gray-900">Confidence rendah</p>
                <p class="text-sm text-gray-500">Eskalasi jika confidence di bawah threshold</p>
              </div>
            </label>
            <label class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100">
              <input type="checkbox" v-model="config.escalate_complaint" class="w-4 h-4 text-indigo-600 rounded">
              <div>
                <p class="font-medium text-gray-900">Komplain terdeteksi</p>
                <p class="text-sm text-gray-500">Eskalasi jika AI mendeteksi komplain/keluhan</p>
              </div>
            </label>
            <label class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100">
              <input type="checkbox" v-model="config.escalate_order" class="w-4 h-4 text-indigo-600 rounded">
              <div>
                <p class="font-medium text-gray-900">Intent order</p>
                <p class="text-sm text-gray-500">Eskalasi jika customer ingin memesan</p>
              </div>
            </label>
            <label class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100">
              <input type="checkbox" v-model="config.escalate_urgent" class="w-4 h-4 text-indigo-600 rounded">
              <div>
                <p class="font-medium text-gray-900">Kata kunci urgent</p>
                <p class="text-sm text-gray-500">Eskalasi jika ada kata "urgent", "darurat", dll</p>
              </div>
            </label>
          </div>
        </div>

        <!-- Notification Settings -->
        <div class="mt-6">
          <h3 class="font-medium text-gray-900 mb-4">Notifikasi Eskalasi</h3>
          <div class="space-y-3">
            <label class="flex items-center gap-3">
              <input type="checkbox" v-model="config.notify_whatsapp" class="w-4 h-4 text-indigo-600 rounded">
              <span class="text-gray-700">Kirim notifikasi WhatsApp ke admin</span>
            </label>
            <label class="flex items-center gap-3">
              <input type="checkbox" v-model="config.notify_email" class="w-4 h-4 text-indigo-600 rounded">
              <span class="text-gray-700">Kirim notifikasi Email</span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Save Button -->
    <div class="flex justify-end gap-3">
      <button
        @click="resetConfig"
        class="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
      >
        Reset
      </button>
      <button
        @click="saveConfig"
        :disabled="saving"
        class="px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
      >
        <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        {{ saving ? 'Menyimpan...' : 'Simpan Konfigurasi' }}
      </button>
    </div>

    <!-- Toast Notification -->
    <Transition name="slide-up">
      <div v-if="toast.show" :class="[
        'fixed bottom-6 right-6 px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 z-50',
        toast.type === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
      ]">
        <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
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

interface ProviderInfo {
  id: string
  name: string
  description: string
  models: ModelInfo[]
  requires_key: boolean
  free_available: boolean
}

interface ModelInfo {
  id: string
  name: string
  description: string
  input_cost: number
  output_cost: number
}

const activeTab = ref('provider')
const saving = ref(false)
const testingConnection = ref(false)
const showApiKey = ref(false)
const apiKeyInput = ref('')
const connectionTestResult = ref<{ success: boolean; message: string } | null>(null)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })
const providers = ref<ProviderInfo[]>([])

const tabs = [
  { id: 'provider', label: 'Provider & API Key' },
  { id: 'general', label: 'Umum' },
  { id: 'prompt', label: 'System Prompt' },
  { id: 'context', label: 'Konteks Bisnis' },
  { id: 'escalation', label: 'Eskalasi' }
]

const config = ref({
  enabled: false,
  ai_provider: 'gemini',
  model: 'gemini-2.0-flash',
  api_key_set: false,
  use_system_key: true,
  confidence_threshold: 0.80,
  language: 'id',
  max_tokens: 200,
  system_prompt: `Kamu adalah asisten toko online yang ramah dan membantu.
Jawab pertanyaan customer dengan sopan, jelas, dan ringkas.
Jika tidak tahu jawabannya, katakan dengan jujur dan sarankan untuk menghubungi admin.
Selalu gunakan bahasa yang santun dan profesional.`,
  business_name: '',
  business_type: 'retail',
  business_hours: '',
  business_description: '',
  business_address: '',
  payment_methods: '',
  escalate_low_confidence: true,
  escalate_complaint: true,
  escalate_order: false,
  escalate_urgent: true,
  notify_whatsapp: true,
  notify_email: false
})

const selectedProviderInfo = computed(() => {
  return providers.value.find(p => p.id === config.value.ai_provider)
})

const currentModels = computed(() => {
  return selectedProviderInfo.value?.models || []
})

const promptPresets = [
  {
    id: 'friendly',
    name: '🤗 Ramah & Santai',
    description: 'Cocok untuk toko casual',
    prompt: `Kamu adalah asisten toko yang ramah dan santai.
Jawab dengan bahasa yang hangat dan bersahabat.
Gunakan emoji sesekali untuk kesan friendly.
Bantu customer dengan sepenuh hati!`
  },
  {
    id: 'professional',
    name: '💼 Profesional',
    description: 'Cocok untuk B2B/jasa',
    prompt: `Kamu adalah asisten profesional untuk layanan kami.
Jawab dengan bahasa yang sopan, jelas, dan to-the-point.
Berikan informasi yang akurat dan terstruktur.
Tunjukkan keahlian dan kredibilitas dalam setiap respons.`
  },
  {
    id: 'sales',
    name: '🛒 Sales Oriented',
    description: 'Fokus closing penjualan',
    prompt: `Kamu adalah sales assistant yang membantu customer berbelanja.
Tunjukkan value dan benefit produk dengan menarik.
Gunakan teknik soft selling yang tidak memaksa.
Bantu customer menemukan produk yang tepat dan dorong ke pembelian.`
  }
]

const getProviderColor = (providerId: string) => {
  const colors: Record<string, string> = {
    gemini: 'bg-blue-500',
    openai: 'bg-green-600',
    groq: 'bg-orange-500',
    anthropic: 'bg-amber-600'
  }
  return colors[providerId] || 'bg-gray-500'
}

const getProviderLink = (providerId: string) => {
  const links: Record<string, string> = {
    gemini: 'https://aistudio.google.com/app/apikey',
    openai: 'https://platform.openai.com/api-keys',
    groq: 'https://console.groq.com/keys',
    anthropic: 'https://console.anthropic.com/settings/keys'
  }
  return links[providerId] || '#'
}

const selectProvider = (providerId: string) => {
  config.value.ai_provider = providerId
  // Set default model for provider
  const provider = providers.value.find(p => p.id === providerId)
  if (provider && provider.models.length > 0) {
    config.value.model = provider.models[0].id
  }
  // Reset connection test
  connectionTestResult.value = null
  apiKeyInput.value = ''
}

const onUseSystemKeyChange = () => {
  if (config.value.use_system_key) {
    apiKeyInput.value = ''
    connectionTestResult.value = null
  }
}

const applyPreset = (preset: any) => {
  config.value.system_prompt = preset.prompt
  showToast('Preset diterapkan!', 'success')
}

const toggleAI = async () => {
  await saveConfig()
}

const { fetch: apiFetch } = useApi()

const loadProviders = async () => {
  try {
    const res = await apiFetch<{ providers: ProviderInfo[] }>('/api/ai/providers')
    if (res?.providers) {
      providers.value = res.providers
    }
  } catch (err) {
    console.error('Failed to load providers:', err)
    // Set default providers if API fails
    providers.value = [
      {
        id: 'gemini',
        name: 'Google Gemini',
        description: 'Google\'s multimodal AI model',
        requires_key: true,
        free_available: true,
        models: [
          { id: 'gemini-2.0-flash', name: 'Gemini 2.0 Flash', description: 'Fast', input_cost: 0.075, output_cost: 0.30 }
        ]
      }
    ]
  }
}

const loadConfig = async () => {
  try {
    const res = await apiFetch<any>('/api/ai/config')
    if (res) {
      config.value = { ...config.value, ...res }
    }
  } catch (err) {
    console.error('Failed to load config:', err)
  }
}

const testConnection = async () => {
  testingConnection.value = true
  connectionTestResult.value = null

  try {
    const res = await apiFetch<{ success: boolean; message?: string; error?: string }>('/api/ai/test-connection', {
      method: 'POST',
      body: JSON.stringify({
        provider: config.value.ai_provider,
        api_key: apiKeyInput.value,
        model: config.value.model
      })
    })

    if (res?.success) {
      connectionTestResult.value = { success: true, message: 'Koneksi berhasil! API key valid.' }
    } else {
      connectionTestResult.value = { success: false, message: res?.error || 'Koneksi gagal' }
    }
  } catch (err: any) {
    connectionTestResult.value = { 
      success: false, 
      message: err.data?.error || 'Gagal menguji koneksi' 
    }
  } finally {
    testingConnection.value = false
  }
}

const saveConfig = async () => {
  saving.value = true

  try {
    const payload: any = { ...config.value }
    
    // Only include API key if user entered a new one
    if (apiKeyInput.value && !config.value.use_system_key) {
      payload.api_key = apiKeyInput.value
    }

    await apiFetch('/api/ai/config', {
      method: 'PUT',
      body: JSON.stringify(payload)
    })

    // Update api_key_set if we saved a new key
    if (apiKeyInput.value) {
      config.value.api_key_set = true
    }

    showToast('Konfigurasi berhasil disimpan!', 'success')
  } catch (err: any) {
    showToast(err.data?.error || 'Gagal menyimpan konfigurasi', 'error')
  } finally {
    saving.value = false
  }
}

const resetConfig = () => {
  config.value = {
    enabled: false,
    ai_provider: 'gemini',
    model: 'gemini-2.0-flash',
    api_key_set: false,
    use_system_key: true,
    confidence_threshold: 0.80,
    language: 'id',
    max_tokens: 200,
    system_prompt: `Kamu adalah asisten toko online yang ramah dan membantu.
Jawab pertanyaan customer dengan sopan, jelas, dan ringkas.
Jika tidak tahu jawabannya, katakan dengan jujur dan sarankan untuk menghubungi admin.
Selalu gunakan bahasa yang santun dan profesional.`,
    business_name: '',
    business_type: 'retail',
    business_hours: '',
    business_description: '',
    business_address: '',
    payment_methods: '',
    escalate_low_confidence: true,
    escalate_complaint: true,
    escalate_order: false,
    escalate_urgent: true,
    notify_whatsapp: true,
    notify_email: false
  }
  apiKeyInput.value = ''
  connectionTestResult.value = null
  showToast('Konfigurasi direset', 'success')
}

const showToast = (message: string, type: 'success' | 'error') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

onMounted(async () => {
  await loadProviders()
  await loadConfig()
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

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #4f46e5;
  cursor: pointer;
}
</style>
