import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      meta: { title: '登录', guest: true, showTabBar: false },
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      meta: { title: '注册', guest: true, showTabBar: false },
      component: () => import('@/views/RegisterView.vue'),
    },
    {
      path: '/',
      name: 'wardrobe',
      meta: { title: '衣柜', showTabBar: true },
      component: () => import('@/views/WardrobeView.vue'),
    },
    {
      path: '/upload',
      name: 'upload',
      meta: { title: '上传衣物', showTabBar: false },
      component: () => import('@/views/UploadView.vue'),
    },
    {
      path: '/outfit',
      name: 'outfit',
      meta: { title: '搭配', showTabBar: true },
      component: () => import('@/views/OutfitView.vue'),
    },
    {
      path: '/outfit/new',
      name: 'outfit-new',
      meta: { title: '新建搭配', showTabBar: false },
      component: () => import('@/views/OutfitEditorView.vue'),
    },
    {
      path: '/outfit/:id',
      name: 'outfit-detail',
      meta: { title: '搭配详情', showTabBar: false },
      component: () => import('@/views/OutfitDetailView.vue'),
    },
    {
      path: '/my',
      name: 'profile',
      meta: { title: '我的', showTabBar: true },
      component: () => import('@/views/ProfileView.vue'),
    },
    {
      path: '/my/edit',
      name: 'profile-edit',
      meta: { title: '个人编辑', showTabBar: false },
      component: () => import('@/views/ProfileEditView.vue'),
    },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('wardrobe_token')
  const isGuestRoute = to.meta?.guest === true
  if (!token && !isGuestRoute) {
    return '/login'
  }
})

export default router
