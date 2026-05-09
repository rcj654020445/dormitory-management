<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { Building } from '@/types'

const buildings = ref<Building[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  // TODO: Fetch buildings from API
  loading.value = false
})
</script>

<template>
  <div class="building-list">
    <header class="page-header">
      <h1>楼宇管理</h1>
      <button class="btn-primary">新增楼宇</button>
    </header>

    <div v-if="loading" class="loading">加载中...</div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>名称</th>
          <th>性别</th>
          <th>楼层数</th>
          <th>每层房间数</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="building in buildings" :key="building.id">
          <td>{{ building.name }}</td>
          <td>{{ building.gender === 'male' ? '男' : '女' }}</td>
          <td>{{ building.floor_count }}</td>
          <td>{{ building.room_per_floor }}</td>
          <td>{{ building.status }}</td>
          <td>
            <button class="btn-link">编辑</button>
            <button class="btn-link">查看房间</button>
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
