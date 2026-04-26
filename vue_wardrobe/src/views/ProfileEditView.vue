<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { showToast, showLoadingToast, closeToast, showConfirmDialog } from 'vant'

const router = useRouter()
const authStore = useAuthStore()

// Avatar
const avatarUrl = ref('')

// Nickname
const nickname = ref('')
const savingNickname = ref(false)

// Password
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const savingPassword = ref(false)

onMounted(async () => {
  await authStore.fetchUser()
  if (authStore.user) {
    nickname.value = authStore.user.nickname || ''
    avatarUrl.value = authStore.user.avatar || ''
  }
})

async function onUploadAvatar() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return

    showLoadingToast({ message: '上传头像中...', forbidClick: true })
    try {
      const result = await authApi.uploadAvatar(file)
      avatarUrl.value = result.avatar
      if (authStore.user) authStore.user.avatar = result.avatar
      closeToast()
      showToast('头像已更新')
    } catch {
      closeToast()
      showToast('头像上传失败')
    }
  }
  input.click()
}

async function saveNickname() {
  if (!nickname.value.trim()) {
    showToast('昵称不能为空')
    return
  }
  savingNickname.value = true
  showLoadingToast({ message: '保存中...', forbidClick: true })
  try {
    await authApi.updateProfile(nickname.value.trim())
    if (authStore.user) authStore.user.nickname = nickname.value.trim()
    closeToast()
    showToast('昵称已更新')
  } catch {
    closeToast()
    showToast('昵称更新失败')
  } finally {
    savingNickname.value = false
  }
}

async function savePassword() {
  if (!oldPassword.value || !newPassword.value) {
    showToast('请填写原密码和新密码')
    return
  }
  if (newPassword.value.length < 6) {
    showToast('新密码长度不能少于6位')
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    showToast('两次输入的新密码不一致')
    return
  }
  savingPassword.value = true
  showLoadingToast({ message: '修改中...', forbidClick: true })
  try {
    await authApi.changePassword(oldPassword.value, newPassword.value)
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    closeToast()
    showConfirmDialog({
      title: '密码已修改',
      message: '请重新登录',
      confirmButtonText: '重新登录',
    }).then(() => {
      authStore.logout()
      router.push('/login')
    }).catch(() => {})
  } catch (e: any) {
    closeToast()
    showToast(e.message || '密码修改失败')
  } finally {
    savingPassword.value = false
  }
}
</script>

<template>
  <div class="edit-page">
    <!-- Avatar -->
    <div class="avatar-section" @click="onUploadAvatar">
      <div class="avatar-wrapper">
        <img v-if="avatarUrl" :src="avatarUrl" class="avatar-img" />
        <van-icon v-else name="contact" size="48" color="#1989fa" />
        <div class="avatar-overlay">
          <van-icon name="photograph" size="24" color="#fff" />
        </div>
      </div>
      <span class="avatar-hint">点击更换头像</span>
    </div>

    <!-- Nickname -->
    <van-cell-group inset class="form-section">
      <div class="section-label">修改昵称</div>
      <van-field
        v-model="nickname"
        label="昵称"
        placeholder="输入新的昵称"
        clearable
        maxlength="20"
      />
      <div class="btn-wrapper">
        <van-button
          type="primary"
          size="small"
          round
          :loading="savingNickname"
          :disabled="savingNickname"
          @click="saveNickname"
        >
          保存昵称
        </van-button>
      </div>
    </van-cell-group>

    <!-- Password -->
    <van-cell-group inset class="form-section">
      <div class="section-label">修改密码</div>
      <van-field
        v-model="oldPassword"
        type="password"
        label="原密码"
        placeholder="输入原密码"
      />
      <van-field
        v-model="newPassword"
        type="password"
        label="新密码"
        placeholder="输入新密码（至少6位）"
      />
      <van-field
        v-model="confirmPassword"
        type="password"
        label="确认密码"
        placeholder="再次输入新密码"
      />
      <div class="btn-wrapper">
        <van-button
          type="primary"
          size="small"
          round
          :loading="savingPassword"
          :disabled="savingPassword"
          @click="savePassword"
        >
          修改密码
        </van-button>
      </div>
    </van-cell-group>
  </div>
</template>

<style scoped>
.edit-page {
  padding: 16px 0 40px;
}
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px;
  background: #fff;
  margin-bottom: 16px;
  cursor: pointer;
}
.avatar-wrapper {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #f0f8ff;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}
.avatar-wrapper:hover .avatar-overlay,
.avatar-wrapper:active .avatar-overlay {
  opacity: 1;
}
.avatar-hint {
  font-size: 13px;
  color: #999;
  margin-top: 8px;
}
.form-section {
  margin-bottom: 16px;
}
.section-label {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  padding: 16px 16px 0;
}
.btn-wrapper {
  padding: 12px 16px;
}
</style>
