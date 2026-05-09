<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Allocation } from '@/types'

const allocations = ref<Allocation[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  // TODO: Fetch allocations from API
  loading.value = false
})
</script>

<template>
  <div class="allocation-list">
    <header class="page-header">
      <h1>分配管理</h1>
      <button class="btn-primary">自动分配</button>
    </header>

    <div v-if="loading" class="loading">加载中...</div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>学生ID</th>
          <th>房间ID</th>
          <th>床位号</th>
          <th>状态</th>
          <th>入住时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="alloc in allocations" :key="alloc.id">
          <td>{{ alloc.student_id }}</td>
          <td>{{ alloc.room_id }}</td>
          <td>{{ alloc.bed_number }}</td>
          <td>{{ alloc.status }}</td>
          <td>{{ alloc.check_in_at }}</td>
          <td>
            <button class="btn-link">详情</button>
            <button class="btn-link">退宿</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
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
.btn-link {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  margin-right: 0.5rem;
}
</style>
