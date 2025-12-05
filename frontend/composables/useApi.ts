export const useApi = () => {
    const config = useRuntimeConfig()
    const authStore = useAuthStore()

    // Get base URL and ensure it doesn't have trailing slash
    let baseUrl = process.server
        ? config.apiInternal
        : config.public.apiBase

    // Normalize URL: remove trailing slash and ensure it's just the base URL
    if (baseUrl) {
        baseUrl = baseUrl.replace(/\/+$/, '') // Remove trailing slashes
        // Remove any duplicate base URLs (safety check)
        // If baseUrl contains the same URL twice, extract only the first one
        const urlMatch = baseUrl.match(/^(https?:\/\/[^\/]+)/)
        if (urlMatch && baseUrl.includes(urlMatch[1] + urlMatch[1])) {
            baseUrl = urlMatch[1]
        }
    }

    const apiUrl = baseUrl

    if (process.client) {
        console.log('[DEBUG] useApi - apiUrl:', apiUrl)
        console.log('[DEBUG] useApi - config.public.apiBase:', config.public.apiBase)
    }

    const fetch = async <T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> => {
        const url = `${apiUrl}${endpoint}`

        if (process.client) {
            console.log('API Request:', {
                url,
                method: options.method || 'GET',
                endpoint
            })
        }

        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...options.headers,
        }

        // Add auth token if available
        const token = authStore.token
        if (token) {
            headers['Authorization'] = `Bearer ${token}`
            if (process.client) {
                console.log('[DEBUG] useApi - Token found, adding Authorization header')
            }
        } else {
            if (process.client) {
                console.warn('[DEBUG] useApi - No token found in authStore!')
                console.warn('[DEBUG] useApi - authStore.isAuthenticated:', authStore.isAuthenticated)
            }
        }

        try {
            const response = await $fetch<T>(url, {
                ...options,
                headers,
                onResponseError({ response }) {
                    // This will be called on error responses
                    console.error('Response Error:', {
                        url,
                        status: response.status,
                        statusText: response.statusText,
                        data: response._data
                    })
                },
                onRequestError({ error }) {
                    // This will be called on request errors (network, etc)
                    console.error('Request Error:', {
                        url,
                        error: error.message
                    })
                }
            })

            if (process.client) {
                console.log('API Response:', {
                    url,
                    success: true
                })
            }

            return response
        } catch (error: any) {
            // Log error for debugging
            if (process.client) {
                console.error('API Error:', {
                    url,
                    status: error.statusCode || error.status || error.response?.status,
                    message: error.message,
                    data: error.data || error.response?._data,
                    fullError: error
                })
            }

            // Handle 401 Unauthorized - auto logout (only for protected routes)
            const statusCode = error.statusCode || error.status || error.response?.status
            if (statusCode === 401) {
                // Only logout if we have a token (means we were authenticated)
                if (authStore.token) {
                    authStore.logout()
                    if (process.client) {
                        // Force redirect to login
                        await navigateTo('/login', { replace: true })
                        // Force reload if still on protected route
                        if (window.location.pathname.startsWith('/dashboard')) {
                            window.location.href = '/login'
                        }
                    }
                }
            }

            // Handle error response - check multiple possible error formats
            // Nuxt $fetch error format - check error.data first
            if (error.data) {
                // Ensure error.data has error property
                if (typeof error.data === 'object' && error.data.error) {
                    throw error
                } else if (typeof error.data === 'string') {
                    throw {
                        data: { error: error.data },
                        statusCode: statusCode || 500,
                        message: error.data
                    }
                } else {
                    throw error
                }
            }

            // Check if error has response data (Nuxt 3 format)
            if (error.response?._data) {
                const responseData = error.response._data
                if (typeof responseData === 'object' && responseData.error) {
                    throw {
                        data: responseData,
                        statusCode: error.response.status,
                        message: error.message
                    }
                } else if (typeof responseData === 'string') {
                    throw {
                        data: { error: responseData },
                        statusCode: error.response.status,
                        message: responseData
                    }
                } else {
                    throw {
                        data: responseData,
                        statusCode: error.response.status,
                        message: error.message
                    }
                }
            }

            // Check if error is a FetchError or network error
            const errorMsg = error.message || error.toString() || ''
            if (errorMsg.includes('fetch') || errorMsg.includes('network') || errorMsg.includes('Failed to fetch') || errorMsg.includes('ERR_') || errorMsg.includes('ECONNREFUSED')) {
                throw {
                    data: {
                        error: 'Tidak dapat terhubung ke server. Silakan coba lagi nanti.'
                    },
                    statusCode: 0,
                    message: errorMsg
                }
            }

            // If error doesn't have data, wrap it with proper message
            const errorMessage = error.message || error.statusMessage || error.toString() || 'Terjadi kesalahan'
            throw {
                data: {
                    error: errorMessage
                },
                statusCode: statusCode || 500,
                message: errorMessage
            }
        }
    }

    return {
        apiUrl,
        fetch,
    }
}

