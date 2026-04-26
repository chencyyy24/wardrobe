<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { showToast } from 'vant'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function onLogin() {
  if (!username.value || !password.value) {
    showToast('请输入用户名和密码')
    return
  }
  loading.value = true
  const ok = await authStore.login(username.value, password.value)
  loading.value = false
  if (ok) {
    router.replace('/')
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-header">
      <van-icon name="gem-o" size="56" color="#1989fa" />
      <h2>个人衣柜</h2>
      <p>登录你的衣橱，管理每日穿搭</p>
    </div>

    <div class="auth-form">
      <van-field
        v-model="username"
        label="用户名"
        placeholder="请输入用户名"
        clearable
        :disabled="loading"
      />
      <van-field
        v-model="password"
        type="password"
        label="密码"
        placeholder="请输入密码"
        clearable
        :disabled="loading"
      />

      <div class="auth-btn">
        <van-button
          type="primary"
          block
          round
          size="large"
          :loading="loading"
          @click="onLogin"
        >
          登录
        </van-button>
      </div>

      <div class="auth-link">
        还没有账号？
        <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  background: #fff;
  display: flex;
  flex-direction: column;
}
.auth-header {
  text-align: center;
  padding: 60px 20px 40px;
}
.auth-header h2 {
  font-size: 24px;
  color: #333;
  margin: 16px 0 8px;
}
.auth-header p {
  font-size: 14px;
  color: #999;
}
.auth-form {
  padding: 0 24px;
}
.auth-btn {
  margin-top: 32px;
}
.auth-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #666;
}
.auth-link a {
  color: #1989fa;
  text-decoration: none;
}
</style>
