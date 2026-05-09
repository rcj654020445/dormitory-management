<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoomStore } from '@/stores/room'

const store = useRoomStore()
const showForm = ref(false)
const editingRoom = ref<string | null>(null)

const formData = ref({
  building_id: '',
  number: '',
  floor: 1,
  type: 'standard' as 'standard' | 'suite' | 'triple',
  capacity: 4,
  has_bathroom: false,
  has_ac: false,
})

const roomTypeOptions = [
  { label: 'Standard (标准间)', value: 'standard' },
  { label: 'Suite (套间)', value: 'suite' },
  { label: 'Triple (三人间)', value: 'triple' },
]

onMounted(() => {
  store.fetchRooms()
})

function openCreateForm() {
  editingRoom.value = null
  formData.value = {
    building_id: '',
    number: '',
    floor: 1,
    type: 'standard',
    capacity: 4,
    has_bathroom: false,
    has_ac: false,
  }
  showForm.value = true
}

function openEditForm(room: typeof store.rooms[0]) {
  editingRoom.value = room.id
  formData.value = {
    building_id: room.building_id,
    number: room.number,
    floor: room.floor,
    type: room.type,
    capacity: room.beds_total,
    has_bathroom: room.has_bathroom,
    has_ac: room.has_ac,
  }
  showForm.value = true
}

async function handleSubmit() {
  try {
    if (editingRoom.value) {
      await store.updateRoom(editingRoom.value, formData.value)
    } else {
      await store.createRoom(formData.value)
    }
    showForm.value = false
  } catch (e) {
    // Error is handled in store
  }
}

async function handleDelete(id: string) {
  if (confirm('Are you sure you want to delete this room?')) {
    try {
      await store.deleteRoom(id)
    } catch (e) {
      // Error is handled in store
    }
  }
}

function getStatusClass(status: string) {
  switch (status) {
    case 'available': return 'status-available'
    case 'full': return 'status-full'
    case 'maintenance': return 'status-maintenance'
    default: return ''
  }
}

function formatType(type: string) {
  switch (type) {
    case 'standard': return 'Standard'
    case 'suite': return 'Suite'
    case 'triple': return 'Triple'
    default: return type
  }
}
</script>

<template>
  <div class="room-list">
    <header class="page-header">
      <h1>宿舍管理</h1>
      <button class="btn-primary" @click="openCreateForm">新增宿舍</button>
    </header>

    <div v-if="store.loading" class="loading">加载中...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>楼号</th>
          <th>房间号</th>
          <th>楼层</th>
          <th>类型</th>
          <th>总床位</th>
          <th>已用</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="room in store.rooms" :key="room.id">
          <td>{{ room.building_id }}</td>
          <td>{{ room.number }}</td>
          <td>{{ room.floor }}</td>
          <td>{{ formatType(room.type) }}</td>
          <td>{{ room.beds_total }}</td>
          <td>{{ room.beds_used }}</td>
          <td>
            <span :class="['status-badge', getStatusClass(room.status)]">
              {{ room.status }}
            </span>
          </td>
          <td>
            <button class="btn-link" @click="openEditForm(room)">编辑</button>
            <button class="btn-link btn-danger" @click="handleDelete(room.id)">删除</button>
          </td>
        </tr>
        <tr v-if="store.rooms.length === 0">
          <td colspan="8" class="empty-state">No rooms found</td>
        </tr>
      </tbody>
    </table>

    <div v-if="store.pagination.total_pages > 1" class="pagination">
      <button
        class="btn-secondary"
        :disabled="store.pagination.page <= 1"
        @click="store.fetchRooms(store.pagination.page - 1)"
      >
        Previous
      </button>
      <span class="page-info">
        Page {{ store.pagination.page }} of {{ store.pagination.total_pages }}
      </span>
      <button
        class="btn-secondary"
        :disabled="!store.hasMore"
        @click="store.fetchRooms(store.pagination.page + 1)"
      >
        Next
      </button>
    </div>

    <!-- Modal Form -->
    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <div class="modal">
        <h2>{{ editingRoom ? '编辑宿舍' : '新增宿舍' }}</h2>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>楼号</label>
            <input v-model="formData.building_id" type="text" required />
          </div>
          <div class="form-group">
            <label>房间号</label>
            <input v-model="formData.number" type="text" required />
          </div>
          <div class="form-group">
            <label>楼层</label>
            <input v-model.number="formData.floor" type="number" min="1" required />
          </div>
          <div class="form-group">
            <label>类型</label>
            <select v-model="formData.type" required>
              <option v-for="opt in roomTypeOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>床位容量</label>
            <input v-model.number="formData.capacity" type="number" min="1" max="8" required />
          </div>
          <div class="form-group checkbox-group">
            <label>
              <input v-model="formData.has_bathroom" type="checkbox" />
              有卫生间
            </label>
          </div>
          <div class="form-group checkbox-group">
            <label>
              <input v-model="formData.has_ac" type="checkbox" />
              有空调
            </label>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showForm = false">取消</button>
            <button type="submit" class="btn-primary">{{ editingRoom ? '保存' : '创建' }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}
.btn-primary {
  padding: 0.5rem 1rem;
  background: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.btn-secondary {
  padding: 0.5rem 1rem;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
.btn-secondary:disabled {
  background: #ccc;
  cursor: not-allowed;
}
.btn-link {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  margin-right: 0.5rem;
}
.btn-danger {
  color: #e74c3c;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th,
.data-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}
.loading, .error, .empty-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}
.error {
  color: #e74c3c;
}
.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
}
.status-available {
  background: #d4edda;
  color: #155724;
}
.status-full {
  background: #f8d7da;
  color: #721c24;
}
.status-maintenance {
  background: #fff3cd;
  color: #856404;
}
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}
.page-info {
  color: #666;
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
}
.modal h2 {
  margin-top: 0;
}
.form-group {
  margin-bottom: 1rem;
}
.form-group label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
}
.form-group input,
.form-group select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 1.5rem;
}
</style>