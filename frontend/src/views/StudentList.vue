<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useStudentStore } from '@/stores/student'

const store = useStudentStore()
const showCreateModal = ref(false)

onMounted(() => {
  store.fetchStudents()
})

function loadMore() {
  store.fetchStudents(store.pagination.page + 1)
}
</script>

<template>
  <div class="student-list">
    <header class="page-header">
      <h1>学生管理</h1>
      <button class="btn-primary" @click="showCreateModal = true">新增学生</button>
    </header>

    <div v-if="store.loading && store.students.length === 0" class="loading">
      加载中...
    </div>

    <div v-else-if="store.error" class="error">
      {{ store.error }}
    </div>

    <table v-else class="data-table">
      <thead>
        <tr>
          <th>学号</th>
          <th>姓名</th>
          <th>性别</th>
          <th>专业</th>
          <th>年级</th>
          <th>宿舍</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="student in store.students" :key="student.id">
          <td>{{ student.student_id }}</td>
          <td>{{ student.name }}</td>
          <td>{{ student.gender === 'male' ? '男' : '女' }}</td>
          <td>{{ student.major }}</td>
          <td>{{ student.grade }}</td>
          <td>{{ student.room_id || '未分配' }}</td>
          <td>
            <span :class="['status', student.status]">
              {{ student.status }}
            </span>
          </td>
          <td>
            <button class="btn-link">编辑</button>
            <button class="btn-link">分配</button>
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
.status.pending { background: #f39c12; color: white; }
.status.checked_in { background: #27ae60; color: white; }
.btn-link {
  background: none;
  border: none;
  color: #3498db;
  cursor: pointer;
  margin-right: 0.5rem;
}
</style>
