export interface WhatsAppStatus {
    is_connected: boolean
    jid: string
    status: string
}

export interface WhatsAppConnectionResponse {
    status: string
    message: string
    stream_url?: string
}

export interface QRCodeEvent {
    event: string
    code?: string
    qr_image?: string
    error?: string
}

export interface SendMessageRequest {
    recipient_jid: string
    message: string
}

export interface SendMessageResponse {
    success: boolean
    message_id: string
    status: string
}

export const useWhatsApp = () => {
    const { fetch: apiFetch, apiUrl } = useApi()
    const authStore = useAuthStore()

    // Connect WhatsApp
    const connect = async (): Promise<WhatsAppConnectionResponse> => {
        const response = await apiFetch<WhatsAppConnectionResponse>('/api/whatsapp/connect', {
            method: 'POST'
        })
        return response
    }

    // Disconnect WhatsApp
    const disconnect = async (): Promise<{ status: string; message: string }> => {
        const response = await apiFetch<{ status: string; message: string }>('/api/whatsapp/disconnect', {
            method: 'DELETE'
        })
        return response
    }

    // Get WhatsApp status
    const getStatus = async (): Promise<WhatsAppStatus> => {
        const response = await apiFetch<WhatsAppStatus>('/api/whatsapp/status', {
            method: 'GET'
        })
        return response
    }

    // Listen to QR code stream via SSE
    const listenToQRStream = (
        onQRCode: (qrCode: string) => void,
        onSuccess: () => void,
        onError: (error: string) => void,
        onTimeout: () => void
    ): (() => void) => {
        console.log('[QR Stream] Starting QR stream listener')

        if (!process.client) {
            console.log('[QR Stream] Not on client side, skipping')
            return () => { }
        }

        const token = authStore.token
        if (!token) {
            console.error('[QR Stream] No auth token found')
            onError('Tidak ada token autentikasi')
            return () => { }
        }

        console.log('[QR Stream] Token found, setting up stream')
        let abortController: AbortController | null = null
        let isClosed = false

        const close = () => {
            if (isClosed) return
            console.log('[QR Stream] Closing stream')
            isClosed = true
            if (abortController) {
                abortController.abort()
            }
        }

        // Use NATIVE browser fetch API for SSE (not the custom useApi fetch)
        const startSSE = async () => {
            try {
                abortController = new AbortController()

                // Build the full URL for the SSE stream
                // apiUrl is the backend URL (e.g., http://localhost:8080)
                const streamUrl = `${apiUrl}/api/whatsapp/qr/stream`

                console.log('[QR Stream] API URL:', apiUrl)
                console.log('[QR Stream] Stream URL:', streamUrl)
                console.log('[QR Stream] Initiating native fetch request...')

                // Use globalThis.fetch (native browser fetch) - NOT the custom fetch from useApi
                const response = await globalThis.fetch(streamUrl, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'text/event-stream'
                    },
                    signal: abortController.signal
                })

                console.log('[QR Stream] Response status:', response.status)
                console.log('[QR Stream] Response OK:', response.ok)

                if (!response.ok) {
                    let errorMsg = 'Gagal memulai QR stream'
                    try {
                        const errorData = await response.json()
                        errorMsg = errorData.error || errorMsg
                    } catch {
                        // If response is not JSON, use status text
                        errorMsg = `${response.status} ${response.statusText}`
                    }
                    console.error('[QR Stream] Request failed:', errorMsg)
                    onError(errorMsg)
                    close()
                    return
                }

                console.log('[QR Stream] Response OK, setting up reader...')

                const reader = response.body?.getReader()
                const decoder = new TextDecoder()

                if (!reader) {
                    console.error('[QR Stream] No reader available')
                    onError('Tidak dapat membaca stream')
                    close()
                    return
                }

                console.log('[QR Stream] Reader obtained, starting read loop...')

                let buffer = ''

                while (true) {
                    const { done, value } = await reader.read()

                    if (done) {
                        console.log('[QR Stream] Stream ended')
                        break
                    }

                    buffer += decoder.decode(value, { stream: true })
                    const lines = buffer.split('\n')
                    buffer = lines.pop() || ''

                    let currentEvent = ''
                    let currentData = ''

                    for (const line of lines) {
                        if (line.startsWith('event: ')) {
                            currentEvent = line.substring(7).trim()
                        } else if (line.startsWith('data: ')) {
                            currentData = line.substring(6).trim()
                        } else if (line === '') {
                            // Empty line indicates end of event
                            if (currentData) {
                                try {
                                    console.log('[QR Stream] Received event:', currentEvent, 'data length:', currentData.length)
                                    const data = JSON.parse(currentData) as QRCodeEvent

                                    if (data.event === 'code' && data.qr_image) {
                                        console.log('[QR Stream] QR code received!')
                                        onQRCode(data.qr_image)
                                    } else if (data.event === 'success') {
                                        console.log('[QR Stream] Success!')
                                        onSuccess()
                                        close()
                                        return
                                    } else if (data.event === 'error') {
                                        console.log('[QR Stream] Error:', data.error)
                                        onError(data.error || 'Terjadi kesalahan')
                                        close()
                                        return
                                    } else if (data.event === 'timeout') {
                                        console.log('[QR Stream] Timeout')
                                        onTimeout()
                                        close()
                                        return
                                    }
                                } catch (e) {
                                    console.error('[QR Stream] Failed to parse SSE data:', e, currentData)
                                }
                            }
                            currentEvent = ''
                            currentData = ''
                        }
                    }
                }
            } catch (error: any) {
                if (error.name === 'AbortError') {
                    console.log('[QR Stream] Aborted')
                    return
                }
                console.error('[QR Stream] Error:', error)
                onError(error.message || 'Gagal membaca QR stream')
                close()
            }
        }

        startSSE()

        return close
    }

    // Send a WhatsApp message
    const sendMessage = async (recipientJID: string, message: string): Promise<SendMessageResponse> => {
        const response = await apiFetch<SendMessageResponse>('/api/whatsapp/send', {
            method: 'POST',
            body: JSON.stringify({
                recipient_jid: recipientJID,
                message: message
            })
        })
        return response
    }

    return {
        connect,
        disconnect,
        getStatus,
        listenToQRStream,
        sendMessage
    }
}

