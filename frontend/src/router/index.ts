import { createRouter, createWebHistory } from 'vue-router'
import StudentList from '@/views/StudentList.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/students',
    },
    {
      path: '/students',
      name: 'students',
      component: StudentList,
    },
    {
      path: '/buildings',
      name: 'buildings',
      component: () => import('@/views/BuildingList.vue'),
    },
    {
      path: '/rooms',
      name: 'rooms',
      component: () => import('@/views/RoomList.vue'),
    },
    {
      path: '/allocations',
      name: 'allocations',
      component: () => import('@/views/AllocationList.vue'),
    },
    {
      path: '/repairs',
      name: 'repairs',
      component: () => import('@/views/RepairList.vue'),
    },
    {
      path: '/violations',
      name: 'violations',
      component: () => import('@/views/ViolationList.vue'),
    },
  ],
})

export default router
