<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useViolationStore } from '@/stores/violation'

const store = useViolationStore()
const showCreateModal = ref(false)

const typeLabels: Record<string, string> = {
  late_return: '晚归',
  noise: '噪音',
  damage: '损坏公物',
  property_violation: '财产违规',
  other: '其他',
}

const statusClass: Record<string, string> = {
  pending: 'status-pending',
  resolved: 'status-resolved',
}

const statusLabels: Record<string, string> = {
  pending: '待处理',
  resolved: '已处理',
}

onMounted(() => {
  store.fetchViolations()
})

function loadMore() {
  store.fetchViolations(store.pagination.page + 1)
}

async function handleResolve(id: string) {
  if (!confirm('确认将该违规标记为已处理？')) return
  try {
    await store.resolveViolation(id)
  } catch {
    alert('操作失败：' + store.error)
  }
}

async function handleDelete(id: string) {
  if (!confirm('确认删除该违规记录？')) return
  try {
    await store.deleteViolation(id)
  } catch {
    alert('删除失败：' + store.error)
  }
}
</script>

<template>
  <div class="violation-list">
    <header class="page-header">
      <h1>违规管理</h1>
      <button class="btn-primary" @click="showCreateModal = true">登记违规</button>
    </header>

    <div v-if="store.loading && store.violations.length === 0" class="loading">
      加载中...
    </div>

    <div v-else-if="store.error" class="error">
      {{ store.error }}
    </div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>学生ID</th>
          <th>类型</th>
          <th>描述</th>
          <th>扣分</th>
          <th>处理人</th>
          <th>处理时间</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="violation in store.violations" :key="violation.id">
          <td>{{ violation.student_id }}</td>
          <td>{{ typeLabels[violation.type] || violation.type }}</td>
          <td>{{ violation.description }}</td>
          <td>{{ violation.points }}分</td>
          <td>{{ violation.handled_by }}</td>
          <td>{{ new Date(violation.handled_at).toLocaleDateString() }}</td>
          <td>
            <span :class="['status', statusClass[violation.status]]">
              {{ statusLabels[violation.status] || violation.status }}
            </span>
          </td>
          <td>
            <button
              v-if="violation.status === 'pending'"
              class="btn-link"
              @click="handleResolve(violation.id)"
            >
              标记已处理
            </button>
            <button class="btn-link btn-danger" @click="handleDelete(violation.id)">
              删除
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="store.hasMore" class="load-more">
      <button @click="loadMore">加载更多</button>
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
.data-table th {
  font-weight: 600;
  color: #2c3e50;
}
.status {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
}
.status-pending {
  background: #f39c12;
  color: white;
}
.status-resolved {
  background: #27ae60;
  color: white;
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
.load-more {
  text-align: center;
  margin-top: 1rem;
}
.load-more button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}
</style>
