// Layer 2: Pinia stores — depends on types (Layer 0) and api (Layer 1)
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Repair, RepairWithRoom, Pagination, RepairStatus } from '@/types'
import type { CreateRepairDto, UpdateRepairDto, UpdateStatusDto, ListRepairsParams } from '@/api/repair'
import { repairApi } from '@/api/repair'

export const useRepairStore = defineStore('repair', () => {
  const repairs = ref<RepairWithRoom[]>([])
  const currentRepair = ref<RepairWithRoom | null>(null)
  const pagination = ref<Pagination>({
    page: 1,
    page_size: 20,
    total_items: 0,
    total_pages: 0,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasMore = computed(() => pagination.value.page < pagination.value.total_pages)

  async function fetchRepairs(params: ListRepairsParams = {}) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.list(params)
      if (res.data.success && res.data.data) {
        const newPage = params.page || 1
        if (newPage === 1) {
          repairs.value = res.data.data.data
        } else {
          repairs.value.push(...res.data.data.data)
        }
        pagination.value = res.data.data.pagination
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch repairs'
    } finally {
      loading.value = false
    }
  }

  async function fetchRepair(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.get(id)
      if (res.data.success && res.data.data) {
        currentRepair.value = res.data.data
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch repair'
    } finally {
      loading.value = false
    }
  }

  async function createRepair(data: CreateRepairDto) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.create(data)
      if (res.data.success && res.data.data) {
        repairs.value.unshift(res.data.data)
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to create repair')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create repair'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateRepair(id: string, data: UpdateRepairDto) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.update(id, data)
      if (res.data.success && res.data.data) {
        const idx = repairs.value.findIndex((r) => r.id === id)
        if (idx !== -1) {
          repairs.value[idx] = res.data.data
        }
        if (currentRepair.value?.id === id) {
          currentRepair.value = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to update repair')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update repair'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateStatus(id: string, data: UpdateStatusDto) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.updateStatus(id, data)
      if (res.data.success && res.data.data) {
        const idx = repairs.value.findIndex((r) => r.id === id)
        if (idx !== -1) {
          repairs.value[idx] = res.data.data
        }
        if (currentRepair.value?.id === id) {
          currentRepair.value = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to update status')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update status'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteRepair(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await repairApi.delete(id)
      if (res.data.success) {
        repairs.value = repairs.value.filter((r) => r.id !== id)
        if (currentRepair.value?.id === id) {
          currentRepair.value = null
        }
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete repair'
      throw e
    } finally {
      loading.value = false
    }
  }

  function setFilter(status: RepairStatus | undefined) {
    pagination.value.page = 1
    fetchRepairs({ status })
  }

  return {
    repairs,
    currentRepair,
    pagination,
    loading,
    error,
    hasMore,
    fetchRepairs,
    fetchRepair,
    createRepair,
    updateRepair,
    updateStatus,
    deleteRepair,
    setFilter,
  }
})