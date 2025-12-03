<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">AI Auto-Reply Test</h1>
      <p class="text-gray-600 mt-1">Test AI auto-reply dengan Gemini API</p>
    </div>

    <!-- Test Form -->
    <div class="bg-white rounded-xl border border-gray-100 p-6">
      <h2 class="font-semibold text-gray-900 mb-4">Test AI Response</h2>

      <div class="space-y-4">
        <!-- Message Input -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Customer Message
          </label>
          <textarea
            v-model="testMessage"
            rows="4"
            placeholder="Contoh: Halo, berapa harga produk X?"
            class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
          ></textarea>
        </div>

        <!-- Test Button -->
        <button
          @click="testAI"
          :disabled="!testMessage || testing"
          class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          <svg v-if="testing" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ testing ? 'Testing...' : 'Test AI Response' }}
        </button>
      </div>
    </div>

    <!-- Response -->
    <div v-if="response" class="bg-white rounded-xl border border-gray-100 p-6">
      <h2 class="font-semibold text-gray-900 mb-4">AI Response</h2>

      <!-- Response Text -->
      <div class="mb-4 p-4 bg-green-50 rounded-lg">
        <p class="text-gray-800">{{ response.response }}</p>
      </div>

      <!-- Metadata -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!-- Confidence -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Confidence</p>
          <p class="text-lg font-semibold" :class="getConfidenceColor(response.confidence)">
            {{ (response.confidence * 100).toFixed(1) }}%
          </p>
        </div>

        <!-- Intent -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Intent</p>
          <p class="text-sm font-medium text-gray-900">{{ formatIntent(response.detected_intent) }}</p>
        </div>

        <!-- Tokens -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Tokens Used</p>
          <p class="text-lg font-semibold text-gray-900">{{ response.tokens_used }}</p>
          <p class="text-xs text-gray-500">
            In: {{ response.input_tokens }} | Out: {{ response.output_tokens }}
          </p>
        </div>

        <!-- Cost -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Cost</p>
          <p class="text-lg font-semibold text-gray-900">
            ${{ response.cost_usd.toFixed(6) }}
          </p>
          <p class="text-xs text-gray-500">≈ Rp {{ (response.cost_usd * 15600).toFixed(0) }}</p>
        </div>

        <!-- Response Time -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Response Time</p>
          <p class="text-lg font-semibold text-gray-900">{{ response.response_time_ms }}ms</p>
        </div>

        <!-- Model -->
        <div class="p-3 bg-gray-50 rounded-lg">
          <p class="text-xs text-gray-500 mb-1">Model</p>
          <p class="text-sm font-medium text-gray-900">{{ response.model }}</p>
        </div>

        <!-- Escalation -->
        <div class="p-3 bg-gray-50 rounded-lg col-span-2">
          <p class="text-xs text-gray-500 mb-1">Escalation</p>
          <div v-if="response.should_escalate" class="flex items-center gap-2">
            <span class="text-orange-600 font-semibold">⚠️ Yes</span>
            <span class="text-xs text-gray-600">{{ response.escalation_reason }}</span>
          </div>
          <div v-else class="text-green-600 font-semibold">✓ No</div>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-800 font-medium">Error:</p>
      <p class="text-red-600 text-sm mt-1">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch: apiFetch } = useApi()

const testMessage = ref('')
const testing = ref(false)
const response = ref<any>(null)
const error = ref('')

const testAI = async () => {
  if (!testMessage.value) return

  testing.value = true
  error.value = ''
  response.value = null

  try {
    const res = await apiFetch<any>('/api/ai/test', {
      method: 'POST',
      body: JSON.stringify({
        message: testMessage.value
      })
    })

    response.value = res
  } catch (err: any) {
    error.value = err.data?.error || err.message || 'Failed to test AI'
  } finally {
    testing.value = false
  }
}

const getConfidenceColor = (confidence: number) => {
  if (confidence >= 0.8) return 'text-green-600'
  if (confidence >= 0.6) return 'text-yellow-600'
  return 'text-red-600'
}

const formatIntent = (intent: string) => {
  const intents: Record<string, string> = {
    'price_inquiry': '💰 Price Inquiry',
    'location_inquiry': '📍 Location',
    'hours_inquiry': '⏰ Business Hours',
    'availability_inquiry': '📦 Availability',
    'order_intent': '🛒 Order',
    'complaint': '⚠️ Complaint',
    'shipping_inquiry': '🚚 Shipping',
    'payment_inquiry': '💳 Payment',
    'general_inquiry': '💬 General'
  }
  return intents[intent] || intent
}
</script>
