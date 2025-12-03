export interface Broadcast {
  id: string
  tenant_id: string
  name: string
  message_content: string
  template_id: string | null
  status: 'draft' | 'scheduled' | 'active' | 'sending' | 'completed' | 'cancelled'
  scheduled_at: string | null
  started_at: string | null
  completed_at: string | null
  total_recipients: number
  sent_count: number
  delivered_count: number
  failed_count: number
  is_recurring: boolean
  recurrence_type: 'hourly' | 'daily' | 'weekly' | null
  recurrence_interval: number | null
  recurrence_days: string[] | null
  recurrence_time: string | null
  recurrence_end_date: string | null
  recurrence_count: number | null
  last_executed_at: string | null
  execution_count: number | null
  created_at: string
  updated_at: string
}

export interface BroadcastRecipient {
  id: string
  broadcast_id: string
  customer_id: string
  customer_jid: string
  customer_name: string
  status: 'pending' | 'sent' | 'delivered' | 'failed'
  message_id: string | null
  sent_at: string | null
  delivered_at: string | null
  error_message: string | null
  created_at: string
}

export interface CreateBroadcastRequest {
  name: string
  message_content: string
  template_id?: string
  customer_ids: string[]
  scheduled_at?: string
}

export interface BroadcastStats {
  total_broadcasts: number
  total_messages_sent: number
  total_delivered: number
  total_failed: number
}

export const useBroadcasts = () => {
  const { fetch: apiFetch } = useApi()

  const statusLabels: Record<string, { label: string, color: string }> = {
    draft: { label: 'Draft', color: 'gray' },
    scheduled: { label: 'Terjadwal', color: 'blue' },
    active: { label: 'Aktif', color: 'purple' },
    sending: { label: 'Mengirim', color: 'yellow' },
    completed: { label: 'Selesai', color: 'green' },
    cancelled: { label: 'Dibatalkan', color: 'red' },
  }

  const recipientStatusLabels: Record<string, { label: string, color: string }> = {
    pending: { label: 'Menunggu', color: 'gray' },
    sent: { label: 'Terkirim', color: 'blue' },
    delivered: { label: 'Diterima', color: 'green' },
    failed: { label: 'Gagal', color: 'red' },
  }

  const getBroadcasts = async (status?: string): Promise<{ broadcasts: Broadcast[], total: number }> => {
    const params = status && status !== 'all' ? `?status=${status}` : ''
    return await apiFetch<{ broadcasts: Broadcast[], total: number }>(`/api/broadcasts${params}`)
  }

  const getBroadcast = async (id: string): Promise<{ broadcast: Broadcast, recipients: BroadcastRecipient[] }> => {
    return await apiFetch<{ broadcast: Broadcast, recipients: BroadcastRecipient[] }>(`/api/broadcasts/${id}`)
  }

  const getStats = async (): Promise<BroadcastStats> => {
    return await apiFetch<BroadcastStats>('/api/broadcasts/stats')
  }

  const createBroadcast = async (data: CreateBroadcastRequest): Promise<Broadcast> => {
    return await apiFetch<Broadcast>('/api/broadcasts', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }

  const sendBroadcast = async (id: string): Promise<{ message: string, status: string }> => {
    return await apiFetch<{ message: string, status: string }>(`/api/broadcasts/${id}/send`, {
      method: 'POST'
    })
  }

  const cancelBroadcast = async (id: string): Promise<{ message: string }> => {
    return await apiFetch<{ message: string }>(`/api/broadcasts/${id}/cancel`, {
      method: 'POST'
    })
  }

  const deleteBroadcast = async (id: string): Promise<void> => {
    await apiFetch<{ message: string }>(`/api/broadcasts/${id}`, {
      method: 'DELETE'
    })
  }

  const getStatusInfo = (status: string) => {
    return statusLabels[status] || statusLabels.draft
  }

  const getRecipientStatusInfo = (status: string) => {
    return recipientStatusLabels[status] || recipientStatusLabels.pending
  }

  const formatProgress = (broadcast: Broadcast): string => {
    if (broadcast.total_recipients === 0) return '0%'
    const progress = ((broadcast.sent_count + broadcast.failed_count) / broadcast.total_recipients) * 100
    return `${Math.round(progress)}%`
  }

  return {
    statusLabels,
    recipientStatusLabels,
    getBroadcasts,
    getBroadcast,
    getStats,
    createBroadcast,
    sendBroadcast,
    cancelBroadcast,
    deleteBroadcast,
    getStatusInfo,
    getRecipientStatusInfo,
    formatProgress
  }
}

