export interface Template {
  id: string
  tenant_id: string
  name: string
  category: string
  content: string
  variables: string
  is_active: boolean
  usage_count: number
  created_at: string
  updated_at: string
}

export interface CreateTemplateRequest {
  name: string
  category: string
  content: string
  variables?: string[]
}

export interface UpdateTemplateRequest {
  name?: string
  category?: string
  content?: string
  variables?: string[]
  is_active?: boolean
}

export const useTemplates = () => {
  const { fetch: apiFetch } = useApi()

  const categories = [
    { value: 'general', label: 'Umum', icon: '💬' },
    { value: 'greeting', label: 'Sapaan', icon: '👋' },
    { value: 'promo', label: 'Promo', icon: '🎉' },
    { value: 'follow_up', label: 'Follow Up', icon: '⏰' },
    { value: 'closing', label: 'Penutup', icon: '🙏' },
  ]

  const getTemplates = async (category?: string): Promise<{ templates: Template[], total: number }> => {
    const params = category && category !== 'all' ? `?category=${category}` : ''
    return await apiFetch<{ templates: Template[], total: number }>(`/api/templates${params}`)
  }

  const createTemplate = async (data: CreateTemplateRequest): Promise<Template> => {
    return await apiFetch<Template>('/api/templates', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  const updateTemplate = async (id: string, data: UpdateTemplateRequest): Promise<Template> => {
    return await apiFetch<Template>(`/api/templates/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  }

  const deleteTemplate = async (id: string): Promise<void> => {
    await apiFetch<{ message: string }>(`/api/templates/${id}`, {
      method: 'DELETE'
    })
  }

  const incrementUsage = async (id: string): Promise<void> => {
    await apiFetch<{ message: string }>(`/api/templates/${id}/use`, {
      method: 'POST'
    })
  }

  // Helper to parse variables from content (e.g., {{nama}})
  const extractVariables = (content: string): string[] => {
    const matches = content.match(/\{\{(\w+)\}\}/g)
    if (!matches) return []
    return matches.map(m => m.replace(/\{\{|\}\}/g, ''))
  }

  // Helper to replace variables in content
  const replaceVariables = (content: string, values: Record<string, string>): string => {
    let result = content
    for (const [key, value] of Object.entries(values)) {
      result = result.replace(new RegExp(`\\{\\{${key}\\}\\}`, 'g'), value)
    }
    return result
  }

  const getCategoryInfo = (category: string) => {
    return categories.find(c => c.value === category) || categories[0]
  }

  return {
    categories,
    getTemplates,
    createTemplate,
    updateTemplate,
    deleteTemplate,
    incrementUsage,
    extractVariables,
    replaceVariables,
    getCategoryInfo
  }
}

