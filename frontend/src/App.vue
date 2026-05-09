<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <div id="app">
    <nav class="navbar">
      <div class="nav-brand">
        <RouterLink to="/">🏠 宿舍管理系统</RouterLink>
      </div>
      <ul class="nav-menu">
        <li><RouterLink to="/students">学生管理</RouterLink></li>
        <li><RouterLink to="/buildings">楼宇管理</RouterLink></li>
        <li><RouterLink to="/rooms">宿舍管理</RouterLink></li>
        <li><RouterLink to="/allocations">分配管理</RouterLink></li>
        <li><RouterLink to="/violations">违规管理</RouterLink></li>
      </ul>
    </nav>
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  padding: 1rem 2rem;
  background: #2c3e50;
  color: white;
}
.nav-brand a {
  color: white;
  text-decoration: none;
  font-size: 1.25rem;
  font-weight: bold;
}
.nav-menu {
  display: flex;
  list-style: none;
  gap: 1.5rem;
  margin-left: 2rem;
}
.nav-menu a {
  color: #ecf0f1;
  text-decoration: none;
}
.nav-menu a:hover,
.nav-menu a.router-link-active {
  color: #42b983;
}
.main-content {
  padding: 2rem;
}
</style>
