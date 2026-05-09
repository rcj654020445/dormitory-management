<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRepairStore } from '@/stores/repair'
import type { RepairStatus } from '@/types'

const store = useRepairStore()
const showForm = ref(false)
const showStatusModal = ref(false)
const editingRepair = ref<string | null>(null)
const currentStatusRepair = ref<string | null>(null)

const formData = ref({
  room_id: '',
  type: 'facility' as 'facility' | 'plumbing' | 'electrical' | 'network' | 'cleaning' | 'other',
  description: '',
  priority: 'normal' as 'urgent' | 'normal' | 'low',
})

const statusForm = ref({
  status: '' as RepairStatus,
  repairer_id: '',
  cost: null as number | null,
})

const typeOptions = [
  { label: '设施维修', value: 'facility' },
  { label: '水管/漏水', value: 'plumbing' },
  { label: '电路/灯具', value: 'electrical' },
  { label: '网络故障', value: 'network' },
  { label: '保洁', value: 'cleaning' },
  { label: '其他', value: 'other' },
]

const priorityOptions = [
  { label: '紧急', value: 'urgent' },
  { label: '普通', value: 'normal' },
  { label: '低', value: 'low' },
]

const statusOptions: { label: string; value: RepairStatus }[] = [
  { label: '待派单', value: 'pending' },
  { label: '已派单', value: 'assigned' },
  { label: '维修中', value: 'repairing' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' },
]

onMounted(() => {
  store.fetchRepairs()
})

function openCreateForm() {
  editingRepair.value = null
  formData.value = {
    room_id: '',
    type: 'facility',
    description: '',
    priority: 'normal',
  }
  showForm.value = true
}

function openEditForm(repair: typeof store.repairs[0]) {
  editingRepair.value = repair.id
  formData.value = {
    room_id: repair.room_id,
    type: repair.type,
    description: repair.description,
    priority: repair.priority,
  }
  showForm.value = true
}

function openStatusForm(repair: typeof store.repairs[0]) {
  currentStatusRepair.value = repair.id
  statusForm.value = {
    status: repair.status === 'pending' ? 'assigned' : repair.status,
    repairer_id: repair.repairer_id || '',
    cost: repair.cost,
  }
  showStatusModal.value = true
}

async function handleSubmit() {
  try {
    if (editingRepair.value) {
      await store.updateRepair(editingRepair.value, formData.value)
    } else {
      await store.createRepair(formData.value)
    }
    showForm.value = false
  } catch (e) {
    // Error is handled in store
  }
}

async function handleStatusSubmit() {
  if (!currentStatusRepair.value) return
  try {
    await store.updateStatus(currentStatusRepair.value, {
      status: statusForm.value.status,
      repairer_id: statusForm.value.repairer_id || undefined,
      cost: statusForm.value.cost,
    })
    showStatusModal.value = false
  } catch (e) {
    // Error is handled in store
  }
}

async function handleDelete(id: string) {
  if (confirm('确定要取消此维修单吗？')) {
    try {
      await store.deleteRepair(id)
    } catch (e) {
      // Error is handled in store
    }
  }
}

function getStatusClass(status: RepairStatus) {
  switch (status) {
    case 'pending': return 'status-pending'
    case 'assigned': return 'status-assigned'
    case 'repairing': return 'status-repairing'
    case 'completed': return 'status-completed'
    case 'cancelled': return 'status-cancelled'
    default: return ''
  }
}

function getPriorityClass(priority: string) {
  switch (priority) {
    case 'urgent': return 'priority-urgent'
    case 'normal': return 'priority-normal'
    case 'low': return 'priority-low'
    default: return ''
  }
}

function formatType(type: string) {
  const found = typeOptions.find((t) => t.value === type)
  return found ? found.label : type
}

function formatStatus(status: RepairStatus) {
  const found = statusOptions.find((s) => s.value === status)
  return found ? found.label : status
}

function handleFilterChange(e: Event) {
  const value = (e.target as HTMLSelectElement).value as RepairStatus | ''
  store.fetchRepairs({ status: value || undefined })
}
</script>

<template>
  <div class="repair-list">
    <header class="page-header">
      <h1>维修管理</h1>
      <button class="btn-primary" @click="openCreateForm">新增报修</button>
    </header>

    <div class="filter-bar">
      <label>状态筛选：</label>
      <select @change="handleFilterChange">
        <option value="">全部</option>
        <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>
    </div>

    <div v-if="store.loading" class="loading">加载中...</div>
    <div v-else-if="store.error" class="error">{{ store.error }}</div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>房间号</th>
          <th>楼栋</th>
          <th>类型</th>
          <th>优先级</th>
          <th>状态</th>
          <th>报修人</th>
          <th>维修人</th>
          <th>创建时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="repair in store.repairs" :key="repair.id">
          <td>{{ repair.room_number }}</td>
          <td>{{ repair.building_name }}</td>
          <td>{{ formatType(repair.type) }}</td>
          <td>
            <span :class="['priority-badge', getPriorityClass(repair.priority)]">
              {{ repair.priority }}
            </span>
          </td>
          <td>
            <span :class="['status-badge', getStatusClass(repair.status)]">
              {{ formatStatus(repair.status) }}
            </span>
          </td>
          <td>{{ repair.reporter_name }}</td>
          <td>{{ repair.repairer_name || '-' }}</td>
          <td>{{ new Date(repair.created_at).toLocaleString() }}</td>
          <td>
            <button class="btn-link" @click="openEditForm(repair)">编辑</button>
            <button
              v-if="repair.status === 'pending' || repair.status === 'assigned'"
              class="btn-link"
              @click="openStatusForm(repair)"
            >
              状态
            </button>
            <button
              v-if="repair.status === 'pending' || repair.status === 'assigned'"
              class="btn-link btn-danger"
              @click="handleDelete(repair.id)"
            >
              取消
            </button>
          </td>
        </tr>
        <tr v-if="store.repairs.length === 0">
          <td colspan="9" class="empty-state">暂无维修记录</td>
        </tr>
      </tbody>
    </table>

    <div v-if="store.pagination.total_pages > 1" class="pagination">
      <button
        class="btn-secondary"
        :disabled="store.pagination.page <= 1"
        @click="store.fetchRepairs({ page: store.pagination.page - 1, status: store.pagination.page === 1 ? undefined : undefined })"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ store.pagination.page }} 页，共 {{ store.pagination.total_pages }} 页
      </span>
      <button
        class="btn-secondary"
        :disabled="!store.hasMore"
        @click="store.fetchRepairs({ page: store.pagination.page + 1 })"
      >
        下一页
      </button>
    </div>

    <!-- Create/Edit Form Modal -->
    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <div class="modal">
        <h2>{{ editingRepair ? '编辑维修单' : '新增报修' }}</h2>
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label>房间ID</label>
            <input v-model="formData.room_id" type="text" required />
          </div>
          <div class="form-group">
            <label>维修类型</label>
            <select v-model="formData.type" required>
              <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>优先级</label>
            <select v-model="formData.priority" required>
              <option v-for="opt in priorityOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>故障描述</label>
            <textarea v-model="formData.description" rows="4" required></textarea>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showForm = false">取消</button>
            <button type="submit" class="btn-primary">{{ editingRepair ? '保存' : '创建' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Status Update Modal -->
    <div v-if="showStatusModal" class="modal-overlay" @click.self="showStatusModal = false">
      <div class="modal">
        <h2>更新状态</h2>
        <form @submit.prevent="handleStatusSubmit">
          <div class="form-group">
            <label>状态</label>
            <select v-model="statusForm.status" required>
              <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>维修工人ID（可选）</label>
            <input v-model="statusForm.repairer_id" type="text" />
          </div>
          <div class="form-group">
            <label>费用（可选）</label>
            <input v-model.number="statusForm.cost" type="number" step="0.01" min="0" />
          </div>
          <div class="form-actions">
            <button type="button" class="btn-secondary" @click="showStatusModal = false">取消</button>
            <button type="submit" class="btn-primary">更新</button>
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
.filter-bar {
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.filter-bar select {
  padding: 0.4rem 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
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
.status-badge, .priority-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
}
.status-pending {
  background: #fff3cd;
  color: #856404;
}
.status-assigned {
  background: #d1ecf1;
  color: #0c5460;
}
.status-repairing {
  background: #e2e3ff;
  color: #383d41;
}
.status-completed {
  background: #d4edda;
  color: #155724;
}
.status-cancelled {
  background: #f8d7da;
  color: #721c24;
}
.priority-urgent {
  background: #f8d7da;
  color: #721c24;
}
.priority-normal {
  background: #d1ecf1;
  color: #0c5460;
}
.priority-low {
  background: #e2e3ff;
  color: #383d41;
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
  width: 450px;
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
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 1.5rem;
}
</style>