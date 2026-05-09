// Layer 2: Pinia stores — depends on types (Layer 0) and api (Layer 1)
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Room, Pagination } from '@/types'
import { roomApi } from '@/api/room'

export const useRoomStore = defineStore('room', () => {
  const rooms = ref<Room[]>([])
  const currentRoom = ref<Room | null>(null)
  const pagination = ref<Pagination>({
    page: 1,
    page_size: 20,
    total_items: 0,
    total_pages: 0,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasMore = computed(() => pagination.value.page < pagination.value.total_pages)

  async function fetchRooms(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const res = await roomApi.list(page, pageSize)
      if (res.data.success && res.data.data) {
        if (page === 1) {
          rooms.value = res.data.data.data
        } else {
          rooms.value.push(...res.data.data.data)
        }
        pagination.value = res.data.data.pagination
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch rooms'
    } finally {
      loading.value = false
    }
  }

  async function fetchRoom(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await roomApi.get(id)
      if (res.data.success && res.data.data) {
        currentRoom.value = res.data.data
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch room'
    } finally {
      loading.value = false
    }
  }

  async function createRoom(data: Parameters<typeof roomApi.create>[0]) {
    loading.value = true
    error.value = null
    try {
      const res = await roomApi.create(data)
      if (res.data.success && res.data.data) {
        rooms.value.unshift(res.data.data)
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to create room')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create room'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateRoom(id: string, data: Parameters<typeof roomApi.update>[1]) {
    loading.value = true
    error.value = null
    try {
      const res = await roomApi.update(id, data)
      if (res.data.success && res.data.data) {
        const idx = rooms.value.findIndex((r) => r.id === id)
        if (idx !== -1) {
          rooms.value[idx] = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to update room')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update room'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteRoom(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await roomApi.delete(id)
      if (res.data.success) {
        rooms.value = rooms.value.filter((r) => r.id !== id)
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete room'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    rooms,
    currentRoom,
    pagination,
    loading,
    error,
    hasMore,
    fetchRooms,
    fetchRoom,
    createRoom,
    updateRoom,
    deleteRoom,
  }
})