export interface Tenant {
    id: string
    user_id: string
    business_name: string
    business_type: string
    business_description: string
    business_phone: string
    business_address: string
    is_active: boolean
    created_at: string
    updated_at: string
}

export interface CreateTenantRequest {
    business_name: string
    business_type: string
    business_description: string
    business_phone: string
    business_address: string
}

export interface UpdateTenantRequest {
    business_name: string
    business_type: string
    business_description: string
    business_phone: string
    business_address: string
}

export const useTenant = () => {
    const { fetch } = useApi()

    // Get current tenant
    const getTenant = async (): Promise<Tenant> => {
        const response = await fetch<Tenant>('/api/tenant', {
            method: 'GET'
        })
        return response
    }

    // Create tenant
    const createTenant = async (data: CreateTenantRequest): Promise<Tenant> => {
        const response = await fetch<Tenant>('/api/tenant', {
            method: 'POST',
            body: JSON.stringify(data)
        })
        return response
    }

    // Update tenant
    const updateTenant = async (data: UpdateTenantRequest): Promise<Tenant> => {
        const response = await fetch<Tenant>('/api/tenant', {
            method: 'PUT',
            body: JSON.stringify(data)
        })
        return response
    }

    return {
        getTenant,
        createTenant,
        updateTenant
    }
}



