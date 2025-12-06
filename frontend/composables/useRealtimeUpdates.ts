import { useAuthStore } from '~/stores/auth'

// Event types from backend
export type WSEventType =
  | 'new_message'
  | 'new_customer'
  | 'customer_updated'
  | 'message_sent'
  | 'connection_status'

export interface WSMessage {
  event: WSEventType
  data: any
}

export interface NewMessageData {
  message_id: string
  sender_jid: string
  chat_jid?: string // For outgoing messages, this is the customer JID
  message_text: string
  message_type: string
  media_url?: string // For image/document messages
  timestamp: number
  is_from_me: boolean
  type?: string // For jid_mapping events
  old_jid?: string
  new_jid?: string
}

export interface UseRealtimeUpdatesOptions {
  onNewMessage?: (data: NewMessageData) => void
  onNewCustomer?: (data: any) => void
  onCustomerUpdated?: (data: any) => void
  onMessageSent?: (data: any) => void
  onConnectionChange?: (connected: boolean) => void
  onJidMappingUpdate?: (oldJid: string, newJid: string) => void
}

export const useRealtimeUpdates = (options: UseRealtimeUpdatesOptions = {}) => {
  const authStore = useAuthStore()
  const { apiUrl } = useApi()

  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = 3000 // 3 seconds

  let socket: WebSocket | null = null
  let reconnectTimeout: NodeJS.Timeout | null = null
  let pingInterval: NodeJS.Timeout | null = null

  const connect = () => {
    if (!process.client) return
    if (!authStore.token) {
      console.warn('[WebSocket] No auth token, skipping connection')
      return
    }

    // Build WebSocket URL
    const wsProtocol = apiUrl.startsWith('https') ? 'wss' : 'ws'
    const wsHost = apiUrl.replace(/^https?:\/\//, '')
    const wsUrl = `${wsProtocol}://${wsHost}/api/ws?token=${authStore.token}`

    console.log('[WebSocket] Connecting to:', wsUrl.replace(authStore.token, '***'))

    try {
      socket = new WebSocket(wsUrl)

      socket.onopen = () => {
        console.log('[WebSocket] Connected')
        isConnected.value = true
        reconnectAttempts.value = 0
        options.onConnectionChange?.(true)

        // Start ping interval to keep connection alive
        pingInterval = setInterval(() => {
          if (socket?.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: 'ping' }))
          }
        }, 30000) // 30 seconds
      }

      socket.onmessage = (event) => {
        try {
          // Handle multiple messages in one event (separated by newlines)
          const messages = event.data.split('\n').filter(Boolean)

          for (const msgStr of messages) {
            const message: WSMessage = JSON.parse(msgStr)
            console.log('[WebSocket] Received:', message.event, message.data)

            switch (message.event) {
              case 'new_message':
                const msgData = message.data as NewMessageData
                // Check if this is a jid_mapping update event
                if (msgData.type === 'jid_mapping' && msgData.old_jid && msgData.new_jid) {
                  console.log('[WebSocket] JID mapping update:', msgData.old_jid, '->', msgData.new_jid)
                  options.onJidMappingUpdate?.(msgData.old_jid, msgData.new_jid)
                } else if (msgData.message_text || msgData.message_type === 'image' || msgData.message_type === 'document' || msgData.media_url) {
                  // Call onNewMessage if there's text content OR if it's a media message
                  options.onNewMessage?.(msgData)
                }
                break
              case 'new_customer':
                options.onNewCustomer?.(message.data)
                break
              case 'customer_updated':
                options.onCustomerUpdated?.(message.data)
                break
              case 'message_sent':
                options.onMessageSent?.(message.data)
                break
            }
          }
        } catch (err) {
          console.error('[WebSocket] Failed to parse message:', err)
        }
      }

      socket.onclose = (event) => {
        console.log('[WebSocket] Disconnected:', event.code, event.reason)
        isConnected.value = false
        options.onConnectionChange?.(false)

        if (pingInterval) {
          clearInterval(pingInterval)
          pingInterval = null
        }

        // Attempt to reconnect
        if (reconnectAttempts.value < maxReconnectAttempts) {
          reconnectAttempts.value++
          console.log(`[WebSocket] Reconnecting in ${reconnectDelay}ms (attempt ${reconnectAttempts.value}/${maxReconnectAttempts})`)
          reconnectTimeout = setTimeout(connect, reconnectDelay)
        } else {
          console.warn('[WebSocket] Max reconnect attempts reached')
        }
      }

      socket.onerror = (error) => {
        console.error('[WebSocket] Error:', error)
      }
    } catch (err) {
      console.error('[WebSocket] Failed to create connection:', err)
    }
  }

  const disconnect = () => {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout)
      reconnectTimeout = null
    }

    if (pingInterval) {
      clearInterval(pingInterval)
      pingInterval = null
    }

    if (socket) {
      socket.close()
      socket = null
    }

    isConnected.value = false
    reconnectAttempts.value = maxReconnectAttempts // Prevent auto-reconnect
  }

  // Auto-connect on mount, disconnect on unmount
  onMounted(() => {
    if (authStore.isAuthenticated) {
      connect()
    }
  })

  onUnmounted(() => {
    disconnect()
  })

  // Watch for auth changes
  watch(() => authStore.isAuthenticated, (authenticated) => {
    if (authenticated) {
      connect()
    } else {
      disconnect()
    }
  })

  return {
    isConnected: readonly(isConnected),
    connect,
    disconnect,
    reconnectAttempts: readonly(reconnectAttempts)
  }
}

