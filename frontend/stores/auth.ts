import { defineStore } from 'pinia'

export interface User {
    id: string
    email: string
    full_name: string
    role?: string
    // Add other user fields as needed
}

export const useAuthStore = defineStore('auth', () => {
    // Only set secure flag if using HTTPS
    const isSecure = process.client ? window.location.protocol === 'https:' : false

    const token = useCookie<string | null>('token', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: isSecure, // Only secure if HTTPS
        httpOnly: false
    })

    const user = useCookie<User | null>('user', {
        maxAge: 60 * 60 * 24 * 3, // 3 days
        sameSite: 'lax',
        secure: isSecure, // Only secure if HTTPS
        httpOnly: false
    })

    const isAuthenticated = computed(() => {
        return !!token.value && token.value !== 'null' && token.value !== ''
    })

    function setToken(newToken: string) {
        if (newToken) {
            token.value = newToken
        }
    }

    function setUser(newUser: User) {
        if (newUser) {
            user.value = newUser
        }
    }

    async function logout(): Promise<void> {
        // Call backend to clear HttpOnly cookies first
        if (process.client) {
            try {
                const config = useRuntimeConfig()
                const apiBase = config.public.apiBase || 'http://localhost:8080'
                await $fetch(`${apiBase}/api/auth/logout`, {
                    method: 'POST',
                    credentials: 'include' // Important: include cookies in request
                })
                console.log('[Auth] Backend logout API called successfully')
            } catch (error) {
                console.error('[Auth] Backend logout API call failed:', error)
                // Continue with client-side cleanup even if backend fails
            }
        }

        // Clear cookies by setting to null
        token.value = null
        user.value = null

        return new Promise((resolve) => {
            // Force remove cookies with all possible variations
            if (process.client) {
                console.log('[Auth] Logout: clearing cookies client-side...')

                // Get current domain for cookie clearing
                const domain = window.location.hostname
                const isSecureNow = window.location.protocol === 'https:'

                // Clear all possible cookie variations
                const cookiesToClear = ['token', 'user']
                const paths = ['/', '']
                const domains = ['', domain, `.${domain}`]

                cookiesToClear.forEach(name => {
                    paths.forEach(path => {
                        domains.forEach(d => {
                            // Without Secure
                            const base = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=${path || '/'}`
                            document.cookie = d ? `${base}; domain=${d}` : base
                            document.cookie = d ? `${base}; domain=${d}; SameSite=Lax` : `${base}; SameSite=Lax`

                            // With Secure (for HTTPS)
                            if (isSecureNow) {
                                document.cookie = d ? `${base}; domain=${d}; Secure` : `${base}; Secure`
                                document.cookie = d ? `${base}; domain=${d}; Secure; SameSite=Lax` : `${base}; Secure; SameSite=Lax`
                                document.cookie = d ? `${base}; domain=${d}; Secure; SameSite=None` : `${base}; Secure; SameSite=None`
                            }
                        })
                    })
                })

                // Clear Pinia state from localStorage
                localStorage.removeItem('auth')
                localStorage.removeItem('pinia')

                // Clear sessionStorage too
                sessionStorage.clear()

                console.log('[Auth] Logout: cookies cleared, remaining:', document.cookie)

                // Delay to ensure browser processes cookie removal
                setTimeout(resolve, 150)
            } else {
                resolve()
            }
        })
    }

    return {
        token,
        user,
        isAuthenticated,
        setToken,
        setUser,
        logout
    }
})

