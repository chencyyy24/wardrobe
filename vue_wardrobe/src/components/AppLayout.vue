<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const showTabBar = computed(() => route.meta?.showTabBar !== false)
const showNavBar = computed(() => route.meta?.guest !== true)
const title = computed(() => (route.meta?.title as string) || '个人衣柜')
const showBack = computed(() => !showTabBar.value && route.meta?.guest !== true)

const tabRoutes = ['/', '/outfit', '/my']
const activeTab = computed(() => {
  const idx = tabRoutes.indexOf(route.path)
  return idx >= 0 ? idx : 0
})

function onTabChange(index: number) {
  router.push(tabRoutes[index])
}
</script>

<template>
  <div class="app-container">
    <van-nav-bar
      v-if="showNavBar"
      :title="title"
      :left-arrow="showBack"
      @click-left="router.back()"
      fixed
      placeholder
    />
    <div class="app-content">
      <router-view />
    </div>
    <van-tabbar
      v-if="showTabBar"
      :model-value="activeTab"
      @change="onTabChange"
      route
      safe-area-inset-bottom
    >
      <van-tabbar-item icon="cluster-o" to="/">衣柜</van-tabbar-item>
      <van-tabbar-item icon="gem-o" to="/outfit">搭配</van-tabbar-item>
      <van-tabbar-item icon="contact" to="/my">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
  background: #f7f8fa;
}
.app-content {
  padding-bottom: 50px;
}
</style>
