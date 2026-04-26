<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { showDialog, showToast } from 'vant'

const authStore = useAuthStore()
const router = useRouter()

function onLogout() {
  showDialog({
    title: '退出登录',
    message: '确定要退出当前账号吗？',
    showCancelButton: true,
    confirmButtonColor: '#ee0a24',
  }).then(() => {
    authStore.logout()
    router.push('/login')
  }).catch(() => {})
}
</script>

<template>
  <div class="profile-page">
    <!-- User info -->
    <div class="user-section" @click="router.push('/my/edit')">
      <div class="avatar">
        <img v-if="authStore.user?.avatar" :src="authStore.user.avatar" class="avatar-img" />
        <van-icon v-else name="contact" size="48" color="#1989fa" />
      </div>
      <div class="user-info">
        <h3>{{ authStore.user?.nickname || authStore.user?.username || '用户' }}</h3>
        <p>{{ authStore.user?.username }}</p>
      </div>
      <van-icon name="arrow" color="#c8c9cc" />
    </div>

    <!-- Menu -->
    <van-cell-group inset class="menu-section">
      <van-cell title="功能介绍" is-link>
        <template #icon><van-icon name="info-o" class="cell-icon" /></template>
      </van-cell>
      <van-cell title="上传衣物" is-link @click="router.push('/upload')">
        <template #icon><van-icon name="photograph" class="cell-icon" /></template>
      </van-cell>
      <van-cell title="开始搭配" is-link @click="router.push('/outfit/new')">
        <template #icon><van-icon name="gem-o" class="cell-icon" /></template>
      </van-cell>
    </van-cell-group>

    <van-cell-group inset class="menu-section">
      <van-cell title="关于项目">
        <template #icon><van-icon name="description-o" class="cell-icon" /></template>
        <template #value>
          <span class="about-text">
            基于 Vue 3 + Go 的个人衣橱管理应用，支持衣物上传管理、智能搭配组合和穿搭卡片生成。
          </span>
        </template>
      </van-cell>
    </van-cell-group>

    <!-- Logout -->
    <div class="logout-section">
      <van-button block round plain type="danger" @click="onLogout">
        退出登录
      </van-button>
    </div>

    <div class="footer">
      <p>v1.0.0 · 结课作业演示</p>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 0 0 20px;
}
.user-section {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 20px;
  background: #fff;
  margin-bottom: 12px;
  cursor: pointer;
}
.user-section:active {
  background: #fafafa;
}
.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}
.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #f0f8ff;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
}
.user-info h3 {
  font-size: 18px;
  color: #333;
  margin: 0;
}
.user-info p {
  font-size: 13px;
  color: #999;
  margin: 4px 0 0;
}
.menu-section {
  margin-bottom: 16px;
}
.cell-icon {
  margin-right: 8px;
  font-size: 18px;
}
.about-text {
  font-size: 12px;
  color: #999;
  line-height: 1.5;
}
.logout-section {
  padding: 20px 24px;
}
.footer {
  text-align: center;
  padding: 20px 0;
  color: #ccc;
  font-size: 12px;
}
</style>
