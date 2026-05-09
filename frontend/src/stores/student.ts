// Layer 2: Pinia stores — depends on types (Layer 0) and api (Layer 1)
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Student, Pagination } from '@/types'
import { studentApi } from '@/api/student'

export const useStudentStore = defineStore('student', () => {
  const students = ref<Student[]>([])
  const currentStudent = ref<Student | null>(null)
  const pagination = ref<Pagination>({
    page: 1,
    page_size: 20,
    total_items: 0,
    total_pages: 0,
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasMore = computed(() => pagination.value.page < pagination.value.total_pages)

  async function fetchStudents(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const res = await studentApi.list(page, pageSize)
      if (res.data.success && res.data.data) {
        if (page === 1) {
          students.value = res.data.data.data
        } else {
          students.value.push(...res.data.data.data)
        }
        pagination.value = res.data.data.pagination
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch students'
    } finally {
      loading.value = false
    }
  }

  async function fetchStudent(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await studentApi.get(id)
      if (res.data.success && res.data.data) {
        currentStudent.value = res.data.data
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch student'
    } finally {
      loading.value = false
    }
  }

  async function createStudent(data: Parameters<typeof studentApi.create>[0]) {
    loading.value = true
    error.value = null
    try {
      const res = await studentApi.create(data)
      if (res.data.success && res.data.data) {
        students.value.unshift(res.data.data)
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to create student')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create student'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateStudent(id: string, data: Parameters<typeof studentApi.update>[1]) {
    loading.value = true
    error.value = null
    try {
      const res = await studentApi.update(id, data)
      if (res.data.success && res.data.data) {
        const idx = students.value.findIndex((s) => s.id === id)
        if (idx !== -1) {
          students.value[idx] = res.data.data
        }
        return res.data.data
      }
      throw new Error(res.data.error?.message || 'Failed to update student')
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update student'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteStudent(id: string) {
    loading.value = true
    error.value = null
    try {
      const res = await studentApi.delete(id)
      if (res.data.success) {
        students.value = students.value.filter((s) => s.id !== id)
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete student'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    students,
    currentStudent,
    pagination,
    loading,
    error,
    hasMore,
    fetchStudents,
    fetchStudent,
    createStudent,
    updateStudent,
    deleteStudent,
  }
})
