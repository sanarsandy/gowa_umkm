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
        // Clear cookies by setting maxAge to 0
        token.value = null
        user.value = null

        // Force remove cookies
        if (process.client) {
            document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT'
            document.cookie = 'user=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT'
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

