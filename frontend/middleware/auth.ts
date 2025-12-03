export default defineNuxtRouteMiddleware((to, from) => {
    const authStore = useAuthStore()

    // Check authentication status
    if (!authStore.isAuthenticated) {
        return navigateTo('/login', { replace: true })
    }
})

