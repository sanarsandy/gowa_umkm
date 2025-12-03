export interface DashboardStats {
    total_messages: number
    total_customers: number
    hot_leads: number
    positive_sentiment: number
    messages_today: number
    new_customers_today: number
    is_connected: boolean
    connected_number?: string
}

export interface RecentMessage {
    id: string
    customer_name: string
    customer_jid: string
    message_text: string
    timestamp: string
    is_from_me: boolean
}

export interface RecentCustomer {
    id: string
    customer_name: string
    customer_jid: string
    status: string
    last_message_at: string
    message_count: number
    last_message_text: string
}

export const useDashboard = () => {
    const { fetch } = useApi()

    // Get dashboard statistics
    const getStats = async (): Promise<DashboardStats> => {
        return await fetch<DashboardStats>('/api/dashboard/stats', {
            method: 'GET'
        })
    }

    // Get recent messages
    const getRecentMessages = async (): Promise<RecentMessage[]> => {
        return await fetch<RecentMessage[]>('/api/dashboard/messages/recent', {
            method: 'GET'
        })
    }

    // Get recent customers
    const getRecentCustomers = async (): Promise<RecentCustomer[]> => {
        return await fetch<RecentCustomer[]>('/api/dashboard/customers/recent', {
            method: 'GET'
        })
    }

    // Format phone number from JID
    const formatPhoneFromJID = (jid: string): string => {
        if (!jid) return ''
        // Extract phone number from JID (format: 628xxxx@s.whatsapp.net)
        const phone = jid.split('@')[0]
        if (phone.startsWith('62')) {
            return '+' + phone
        }
        return phone
    }

    // Format relative time
    const formatRelativeTime = (timestamp: string): string => {
        if (!timestamp) return ''
        
        const now = new Date()
        const date = new Date(timestamp)
        const diffMs = now.getTime() - date.getTime()
        const diffMins = Math.floor(diffMs / 60000)
        const diffHours = Math.floor(diffMs / 3600000)
        const diffDays = Math.floor(diffMs / 86400000)

        if (diffMins < 1) return 'Baru saja'
        if (diffMins < 60) return `${diffMins}m`
        if (diffHours < 24) return `${diffHours}j`
        if (diffDays < 7) return `${diffDays}h`
        
        return date.toLocaleDateString('id-ID', {
            day: 'numeric',
            month: 'short'
        })
    }

    // Get status color
    const getStatusColor = (status: string): string => {
        switch (status) {
            case 'hot_lead':
                return 'text-red-600 bg-red-50'
            case 'warm_lead':
                return 'text-orange-600 bg-orange-50'
            case 'cold_lead':
                return 'text-blue-600 bg-blue-50'
            case 'customer':
                return 'text-green-600 bg-green-50'
            case 'complaint':
                return 'text-yellow-600 bg-yellow-50'
            default:
                return 'text-gray-600 bg-gray-50'
        }
    }

    // Get status label
    const getStatusLabel = (status: string): string => {
        switch (status) {
            case 'hot_lead':
                return 'Hot Lead'
            case 'warm_lead':
                return 'Warm Lead'
            case 'cold_lead':
                return 'Cold Lead'
            case 'customer':
                return 'Pelanggan'
            case 'complaint':
                return 'Komplain'
            case 'new':
                return 'Baru'
            default:
                return status
        }
    }

    return {
        getStats,
        getRecentMessages,
        getRecentCustomers,
        formatPhoneFromJID,
        formatRelativeTime,
        getStatusColor,
        getStatusLabel
    }
}


