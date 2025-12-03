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

    function logout() {
        // Clear cookies by setting to null
        token.value = null
        user.value = null

        // Force remove cookies with all possible variations
        if (process.client) {
            // Clear with matching attributes
            const cookieOptions = `path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax${isSecure ? '; Secure' : ''}`
            document.cookie = `token=; ${cookieOptions}`
            document.cookie = `user=; ${cookieOptions}`

            // Also try without SameSite (for older browsers)
            document.cookie = `token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT${isSecure ? '; Secure' : ''}`
            document.cookie = `user=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT${isSecure ? '; Secure' : ''}`

            // Clear without Secure flag as fallback
            document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax'
            document.cookie = 'user=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax'

            // Clear Pinia state from localStorage
            localStorage.removeItem('auth')
            localStorage.removeItem('pinia')
        }
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

