export interface Customer {
    id: string
    tenant_id: string
    customer_jid: string
    customer_name: string | null
    customer_phone: string | null
    status: string
    sentiment: string | null
    intent: string | null
    product_interest: string | null
    last_message_summary: string | null
    message_count: number
    last_message_at: string | null
    first_message_at: string | null
    needs_follow_up: boolean
    tags: string | null
    lead_score: number
    created_at: string
    updated_at: string
}

export interface CustomerListResponse {
    customers: Customer[]
    total: number
    page: number
    limit: number
    total_pages: number
}

export interface CustomerMessage {
    id: string
    message_text: string
    message_type: string
    is_from_me: boolean
    timestamp: string
}

export interface CustomerDetailResponse {
    customer: Customer
    messages: CustomerMessage[]
}

export interface CustomerStats {
    total: number
    new: number
    hot_leads: number
    warm_leads: number
    cold_leads: number
    customers: number
    complaints: number
    need_follow_up: number
}

export interface CustomerFilters {
    page?: number
    limit?: number
    search?: string
    status?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
}

export const useCustomers = () => {
    const { fetch } = useApi()

    // Get customers with pagination and filters
    const getCustomers = async (filters: CustomerFilters = {}): Promise<CustomerListResponse> => {
        const params = new URLSearchParams()
        
        if (filters.page) params.set('page', filters.page.toString())
        if (filters.limit) params.set('limit', filters.limit.toString())
        if (filters.search) params.set('search', filters.search)
        if (filters.status && filters.status !== 'all') params.set('status', filters.status)
        if (filters.sort_by) params.set('sort_by', filters.sort_by)
        if (filters.sort_order) params.set('sort_order', filters.sort_order)

        const queryString = params.toString()
        const endpoint = `/api/customers${queryString ? '?' + queryString : ''}`
        
        return await fetch<CustomerListResponse>(endpoint, {
            method: 'GET'
        })
    }

    // Get customer detail with chat history
    const getCustomerDetail = async (customerId: string): Promise<CustomerDetailResponse> => {
        return await fetch<CustomerDetailResponse>(`/api/customers/${customerId}`, {
            method: 'GET'
        })
    }

    // Get customer by ID (alias for getCustomerDetail, returns just the customer)
    const getCustomerById = async (customerId: string): Promise<Customer> => {
        const response = await fetch<CustomerDetailResponse>(`/api/customers/${customerId}`, {
            method: 'GET'
        })
        return response.customer
    }

    // Update customer
    const updateCustomer = async (customerId: string, data: Partial<{
        customer_name: string
        status: string
        needs_follow_up: boolean
        tags: string
    }>): Promise<{ message: string }> => {
        return await fetch<{ message: string }>(`/api/customers/${customerId}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        })
    }

    // Get customer stats
    const getCustomerStats = async (): Promise<CustomerStats> => {
        return await fetch<CustomerStats>('/api/customers/stats', {
            method: 'GET'
        })
    }

    // Helper: Format phone from JID
    const formatPhoneFromJID = (jid: string): string => {
        if (!jid) return ''
        const phone = jid.split('@')[0].split(':')[0]
        if (phone.startsWith('62')) {
            return '+' + phone
        }
        return phone
    }

    // Helper: Get display name
    const getDisplayName = (customer: Customer): string => {
        if (customer.customer_name) return customer.customer_name
        return formatPhoneFromJID(customer.customer_jid)
    }

    // Helper: Get status color
    const getStatusColor = (status: string): string => {
        const colors: Record<string, string> = {
            'new': 'bg-gray-100 text-gray-800',
            'hot_lead': 'bg-red-100 text-red-800',
            'warm_lead': 'bg-orange-100 text-orange-800',
            'cold_lead': 'bg-blue-100 text-blue-800',
            'customer': 'bg-green-100 text-green-800',
            'complaint': 'bg-yellow-100 text-yellow-800',
            'spam': 'bg-gray-100 text-gray-500',
        }
        return colors[status] || colors['new']
    }

    // Helper: Get status label
    const getStatusLabel = (status: string): string => {
        const labels: Record<string, string> = {
            'new': 'Baru',
            'hot_lead': 'Hot Lead 🔥',
            'warm_lead': 'Warm Lead',
            'cold_lead': 'Cold Lead',
            'customer': 'Pelanggan',
            'complaint': 'Komplain',
            'spam': 'Spam',
        }
        return labels[status] || status
    }

    // Helper: Get sentiment color
    const getSentimentColor = (sentiment: string | null): string => {
        if (!sentiment) return 'text-gray-400'
        const colors: Record<string, string> = {
            'positive': 'text-green-600',
            'neutral': 'text-gray-600',
            'negative': 'text-red-600',
            'mixed': 'text-yellow-600',
        }
        return colors[sentiment] || 'text-gray-400'
    }

    // Helper: Format relative time
    const formatRelativeTime = (timestamp: string | null): string => {
        if (!timestamp) return '-'
        
        const now = new Date()
        const date = new Date(timestamp)
        const diffMs = now.getTime() - date.getTime()
        const diffMins = Math.floor(diffMs / 60000)
        const diffHours = Math.floor(diffMs / 3600000)
        const diffDays = Math.floor(diffMs / 86400000)

        if (diffMins < 1) return 'Baru saja'
        if (diffMins < 60) return `${diffMins} menit lalu`
        if (diffHours < 24) return `${diffHours} jam lalu`
        if (diffDays < 7) return `${diffDays} hari lalu`
        
        return date.toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'short',
            year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
        })
    }

    return {
        getCustomers,
        getCustomerDetail,
        getCustomerById,
        updateCustomer,
        getCustomerStats,
        formatPhoneFromJID,
        getDisplayName,
        getStatusColor,
        getStatusLabel,
        getSentimentColor,
        formatRelativeTime
    }
}


