<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { showToast } from 'vant'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const nickname = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

async function onRegister() {
  if (!username.value || !password.value) {
    showToast('请填写用户名和密码')
    return
  }
  if (password.value.length < 6) {
    showToast('密码长度不能少于6位')
    return
  }
  if (password.value !== confirmPassword.value) {
    showToast('两次密码不一致')
    return
  }
  loading.value = true
  const ok = await authStore.register(username.value, password.value, nickname.value || undefined)
  loading.value = false
  if (ok) {
    router.replace('/')
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-header">
      <van-icon name="friends-o" size="56" color="#1989fa" />
      <h2>创建账号</h2>
      <p>注册后即可管理你的个人衣橱</p>
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
        v-model="nickname"
        label="昵称"
        placeholder="选填，默认为用户名"
        clearable
        :disabled="loading"
      />
      <van-field
        v-model="password"
        type="password"
        label="密码"
        placeholder="至少6位密码"
        clearable
        :disabled="loading"
      />
      <van-field
        v-model="confirmPassword"
        type="password"
        label="确认密码"
        placeholder="再次输入密码"
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
          @click="onRegister"
        >
          注册
        </van-button>
      </div>

      <div class="auth-link">
        已有账号？
        <router-link to="/login">去登录</router-link>
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
  padding: 50px 20px 30px;
}
.auth-header h2 {
  font-size: 22px;
  color: #333;
  margin: 12px 0 6px;
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
