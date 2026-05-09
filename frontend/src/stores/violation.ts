// Layer 2: Pinia stores — depends on types (Layer 0) and api (Layer 1)
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Violation, Pagination } from '@/types'
import { violationApi } from '@/api/violation'

export const useViolationStore = defineStore('violation', () => {
  const violations = ref<Violation[]>([])
  const currentViolation = ref<Violation | null>(null)
  const pagination = ref<Pagination>({
    page: 1,
    page_size: 20,
    total_items: 0,
    total_pages: 0,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasMore = computed(() => pagination.value.page < pagination.value.total_pages)

  async function fetchViolations(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.list({ page, page_size: pageSize })
      if (res.data.success && res.data.data) {
        if (page === 1) {
          violations.value = res.data.data.data
        } else {
          violations.value.push(...res.data.data.data)
        }
        pagination.value = res.data.data.pagination
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch violations'
    } finally {
      loading.value = false
    }
  }

  async function fetchViolation(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.get(id)
      if (res.data.success && res.data.data) {
        currentViolation.value = res.data.data
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch violation'
    } finally {
      loading.value = false
    }
  }

  async function createViolation(data: Parameters<typeof violationApi.create>[0]) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.create(data)
      if (res.data.success && res.data.data) {
        violations.value.unshift(res.data.data)
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to create violation')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create violation'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateViolation(id: string, data: Parameters<typeof violationApi.update>[1]) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.update(id, data)
      if (res.data.success && res.data.data) {
        const idx = violations.value.findIndex((v) => v.id === id)
        if (idx !== -1) {
          violations.value[idx] = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to update violation')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update violation'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteViolation(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.delete(id)
      if (res.data.success) {
        violations.value = violations.value.filter((v) => v.id !== id)
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete violation'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function resolveViolation(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await violationApi.resolve(id, { status: 'resolved' })
      if (res.data.success && res.data.data) {
        const idx = violations.value.findIndex((v) => v.id === id)
        if (idx !== -1) {
          violations.value[idx] = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to resolve violation')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to resolve violation'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    violations,
    currentViolation,
    pagination,
    loading,
    error,
    hasMore,
    fetchViolations,
    fetchViolation,
    createViolation,
    updateViolation,
    deleteViolation,
    resolveViolation,
  }
})
