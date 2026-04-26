<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useOutfitStore } from '@/stores/outfit'
import { useImageUrl } from '@/composables/useImageUrl'
import { formatDate } from '@/utils/format'
import { showDialog, showToast, showConfirmDialog } from 'vant'

const router = useRouter()
const outfitStore = useOutfitStore()
const { resolve } = useImageUrl()

const activeTab = ref(0)
const refreshing = ref(false)

onMounted(() => {
  outfitStore.fetchSummaries()
})

async function onRefresh() {
  refreshing.value = true
  await outfitStore.fetchSummaries()
  refreshing.value = false
}

function viewDetail(id: string) {
  router.push(`/outfit/${id}`)
}

async function cloneOutfit(id: string, name: string) {
  showDialog({
    title: '复制搭配',
    message: `复制「${name}」并创建新搭配？`,
    confirmButtonColor: '#1989fa',
  }).then(async () => {
    try {
      await outfitStore.cloneOutfit(id)
      showToast('复制成功')
    } catch {
      showToast('复制失败')
    }
  }).catch(() => {})
}

async function deleteOutfit(id: string, name: string) {
  showConfirmDialog({
    title: '删除搭配',
    message: `确定要删除「${name}」吗？`,
    confirmButtonColor: '#ee0a24',
  }).then(async () => {
    try {
      await outfitStore.deleteOutfit(id)
      showToast('删除成功')
    } catch {
      showToast('删除失败')
    }
  }).catch(() => {})
}
</script>

<template>
  <div class="outfit-page">
    <van-tabs v-model:active="activeTab">
      <!-- Tab 1: History -->
      <van-tab title="历史记录">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <div class="history-grid">
            <div
              v-for="outfit in outfitStore.summaries"
              :key="outfit.id"
              class="history-card"
            >
              <div class="card-image" @click="viewDetail(outfit.id)">
                <van-image
                  :src="resolve(outfit.card_image)"
                  fit="cover"
                  class="thumb"
                  loading-icon="photo-o"
                >
                  <template #error>
                    <div class="img-placeholder">
                      <van-icon name="gem-o" size="32" color="#c8c9cc" />
                    </div>
                  </template>
                </van-image>
              </div>
              <div class="card-info">
                <div class="card-name">{{ outfit.name || '未命名搭配' }}</div>
                <div class="card-meta">
                  <span>{{ formatDate(outfit.created_at) }}</span>
                  <span>{{ outfit.item_count }}件</span>
                </div>
                <div class="card-actions">
                  <van-button size="mini" plain type="primary" @click="cloneOutfit(outfit.id, outfit.name)">
                    复制
                  </van-button>
                  <van-button size="mini" plain type="danger" @click="deleteOutfit(outfit.id, outfit.name)">
                    删除
                  </van-button>
                </div>
              </div>
            </div>

            <div v-if="outfitStore.summaries.length === 0 && !outfitStore.loading" class="empty-state">
              <van-empty description="暂无搭配记录" />
            </div>
            <div v-if="outfitStore.loading" class="loading-state">
              <van-loading size="20" />
            </div>
          </div>
        </van-pull-refresh>
      </van-tab>

      <!-- Tab 2: New Outfit -->
      <van-tab title="新建搭配">
        <div class="new-outfit-tab">
          <div class="new-outfit-hero">
            <van-icon name="gem-o" size="64" color="#1989fa" />
            <h3>开始搭配</h3>
            <p>选择不同部位的衣物，组合你的专属穿搭</p>
            <van-button
              type="primary"
              round
              size="large"
              class="start-btn"
              @click="router.push('/outfit/new')"
            >
              开始搭配
            </van-button>
          </div>

          <!-- Quick clone from history -->
          <div v-if="outfitStore.summaries.length > 0" class="quick-clone">
            <div class="section-title">从历史复制</div>
            <div class="recent-list">
              <div
                v-for="outfit in outfitStore.summaries.slice(0, 4)"
                :key="outfit.id"
                class="recent-card"
                @click="cloneOutfit(outfit.id, outfit.name)"
              >
                <van-image
                  :src="resolve(outfit.card_image)"
                  fit="cover"
                  class="recent-thumb"
                >
                  <template #error>
                    <div class="recent-placeholder">
                      <van-icon name="gem-o" size="20" color="#c8c9cc" />
                    </div>
                  </template>
                </van-image>
                <div class="recent-name">{{ outfit.name || '未命名' }}</div>
              </div>
            </div>
          </div>
        </div>
      </van-tab>
    </van-tabs>
  </div>
</template>

<style scoped>
.outfit-page {
  min-height: 100vh;
}
.history-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 12px;
}
.history-card {
  background: #fff;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}
.card-image {
  width: 100%;
  aspect-ratio: 3/4;
  overflow: hidden;
  cursor: pointer;
}
.thumb, .img-placeholder {
  width: 100%;
  height: 100%;
}
.img-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
.card-info {
  padding: 8px 10px;
}
.card-name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-meta {
  font-size: 11px;
  color: #999;
  margin: 4px 0 6px;
  display: flex;
  justify-content: space-between;
}
.card-actions {
  display: flex;
  gap: 6px;
}
.new-outfit-tab {
  padding: 20px 16px;
}
.new-outfit-hero {
  text-align: center;
  padding: 40px 0;
}
.new-outfit-hero h3 {
  font-size: 20px;
  color: #333;
  margin: 16px 0 8px;
}
.new-outfit-hero p {
  font-size: 14px;
  color: #999;
  margin-bottom: 24px;
}
.start-btn {
  width: 200px;
}
.quick-clone {
  margin-top: 20px;
}
.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}
.recent-list {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}
.recent-card {
  cursor: pointer;
  text-align: center;
}
.recent-thumb {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}
.recent-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
.recent-name {
  font-size: 11px;
  color: #666;
  margin-top: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.empty-state, .loading-state {
  grid-column: 1 / -1;
  padding: 20px 0;
}
</style>
