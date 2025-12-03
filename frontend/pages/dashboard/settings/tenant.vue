<template>
  <div>
    <div class="bg-white rounded-xl shadow-sm p-8 border border-gray-100">
      <div class="mb-6">
        <h2 class="text-2xl font-bold text-gray-900">Informasi Bisnis</h2>
        <p class="text-gray-600 mt-1">Kelola informasi bisnis Anda</p>
      </div>

      <div v-if="loading && !tenant" class="text-center py-12">
        <div class="animate-spin w-12 h-12 border-4 border-indigo-600 border-t-transparent rounded-full mx-auto"></div>
        <p class="text-gray-600 mt-4">Memuat data...</p>
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid md:grid-cols-2 gap-6">
          <div>
            <label for="business_name" class="block text-sm font-medium text-gray-700 mb-1">
              Nama Bisnis <span class="text-red-500">*</span>
            </label>
            <input
              id="business_name"
              v-model="form.business_name"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              placeholder="Nama bisnis Anda"
            />
          </div>

          <div>
            <label for="business_type" class="block text-sm font-medium text-gray-700 mb-1">
              Jenis Bisnis <span class="text-red-500">*</span>
            </label>
            <select
              id="business_type"
              v-model="form.business_type"
              required
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            >
              <option value="">Pilih jenis bisnis</option>
              <option value="UMKM">UMKM</option>
              <option value="Retail">Retail</option>
              <option value="F&B">Food & Beverage</option>
              <option value="Service">Jasa</option>
              <option value="E-commerce">E-commerce</option>
              <option value="Manufacturing">Manufaktur</option>
              <option value="Other">Lainnya</option>
            </select>
          </div>
        </div>

        <div>
          <label for="business_description" class="block text-sm font-medium text-gray-700 mb-1">
            Deskripsi Bisnis
          </label>
          <textarea
            id="business_description"
            v-model="form.business_description"
            rows="4"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="Jelaskan tentang bisnis Anda"
          ></textarea>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
          <div>
            <label for="business_phone" class="block text-sm font-medium text-gray-700 mb-1">
              Nomor Telepon
            </label>
            <input
              id="business_phone"
              v-model="form.business_phone"
              type="tel"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              placeholder="+62 812 3456 7890"
            />
          </div>

          <div>
            <label for="business_address" class="block text-sm font-medium text-gray-700 mb-1">
              Alamat
            </label>
            <input
              id="business_address"
              v-model="form.business_address"
              type="text"
              class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              placeholder="Alamat bisnis Anda"
            />
          </div>
        </div>

        <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
          <p class="text-red-800 text-sm">{{ error }}</p>
        </div>

        <div v-if="success" class="bg-green-50 border border-green-200 rounded-lg p-4">
          <p class="text-green-800 text-sm">{{ success }}</p>
        </div>

        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="handleCancel"
            class="px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="saving" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Menyimpan...
            </span>
            <span v-else>Simpan Perubahan</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const tenantApi = useTenant()
const router = useRouter()

const tenant = ref<any>(null)
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const success = ref('')

const form = ref({
  business_name: '',
  business_type: '',
  business_description: '',
  business_phone: '',
  business_address: ''
})

// Load tenant data
onMounted(async () => {
  await loadTenant()
})

const loadTenant = async () => {
  loading.value = true
  error.value = ''
  
  try {
    tenant.value = await tenantApi.getTenant()
    form.value = {
      business_name: tenant.value.business_name || '',
      business_type: tenant.value.business_type || '',
      business_description: tenant.value.business_description || '',
      business_phone: tenant.value.business_phone || '',
      business_address: tenant.value.business_address || ''
    }
  } catch (err: any) {
    console.error('Failed to load tenant:', err)
    if (err.data?.error === 'Tenant not found' || err.statusCode === 404) {
      error.value = 'Tenant tidak ditemukan. Silakan buat tenant terlebih dahulu.'
    } else {
      error.value = err.data?.error || 'Gagal memuat data tenant'
    }
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  saving.value = true
  error.value = ''
  success.value = ''

  try {
    if (tenant.value) {
      // Update existing tenant
      await tenantApi.updateTenant(form.value)
      success.value = 'Data bisnis berhasil diperbarui'
    } else {
      // Create new tenant
      await tenantApi.createTenant(form.value)
      success.value = 'Data bisnis berhasil dibuat'
    }
    
    // Reload tenant data
    await loadTenant()
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (err: any) {
    console.error('Failed to save tenant:', err)
    error.value = err.data?.error || 'Gagal menyimpan data bisnis'
  } finally {
    saving.value = false
  }
}

const handleCancel = () => {
  router.push('/dashboard')
}
</script>



