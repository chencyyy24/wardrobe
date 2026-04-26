<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoryStore } from '@/stores/category'
import { useClothingStore } from '@/stores/clothing'
import { useImageUrl } from '@/composables/useImageUrl'
import type { ClothingItem } from '@/types/clothing'
import { showDialog, showToast } from 'vant'
import { formatDate } from '@/utils/format'

const router = useRouter()
const categoryStore = useCategoryStore()
const clothingStore = useClothingStore()
const { resolve } = useImageUrl()

const refreshing = ref(false)
const detailVisible = ref(false)
const detailItem = ref<ClothingItem | null>(null)

onMounted(async () => {
  await categoryStore.load()
  await clothingStore.fetchList()
})

function onCategoryChange(categoryId: number | undefined) {
  clothingStore.activeCategoryId = categoryId
}

async function onRefresh() {
  refreshing.value = true
  await clothingStore.fetchList(clothingStore.activeCategoryId)
  refreshing.value = false
}

function showDetail(item: ClothingItem) {
  detailItem.value = item
  detailVisible.value = true
}

function confirmDelete(item: ClothingItem) {
  detailVisible.value = false
  showDialog({
    title: '删除衣物',
    message: `确定要删除「${item.name}」吗？`,
    confirmButtonColor: '#ee0a24',
  }).then(async () => {
    try {
      await clothingStore.deleteItem(item.id)
      showToast('删除成功')
    } catch {
      showToast('删除失败')
    }
  }).catch(() => {})
}
</script>

<template>
  <div class="wardrobe-page">
    <!-- 分类筛选 -->
    <CategoryTabs
      :categories="categoryStore.categories"
      :active-id="clothingStore.activeCategoryId"
      @change="onCategoryChange"
    />

    <!-- 列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <div class="items-container">
        <div
          v-if="clothingStore.filteredItems.length === 0 && !clothingStore.loading"
          class="empty-state"
        >
          <van-empty description="还没有衣物，快去上传吧" />
        </div>

        <div v-else class="card-grid">
          <div
            v-for="item in clothingStore.filteredItems"
            :key="item.id"
            class="card-wrapper"
            @click="showDetail(item)"
          >
            <div class="card-image-wrapper">
              <van-image
                :src="resolve(item.original_image)"
                fit="cover"
                class="card-image"
                loading-icon="photo-o"
              />
            </div>
            <div class="item-name">{{ item.name }}</div>
          </div>
        </div>

        <div v-if="clothingStore.loading" class="loading-state">
          <van-loading size="24">加载中...</van-loading>
        </div>
      </div>
    </van-pull-refresh>

    <!-- 衣物详情弹窗 -->
    <van-popup
      v-model:show="detailVisible"
      position="center"
      round
      closeable
      style="width: 85%; border-radius: 12px; overflow: hidden;"
    >
      <div v-if="detailItem" class="detail-card">
        <van-image
          :src="resolve(detailItem.original_image)"
          fit="contain"
          class="detail-image"
        />
        <div class="detail-info">
          <h3 class="detail-name">{{ detailItem.name }}</h3>
          <p class="detail-meta">
            <span>分类：{{ detailItem.category?.name || '未分类' }}</span>
            <span v-if="detailItem.subcategory"> / {{ detailItem.subcategory.name }}</span>
          </p>
          <p class="detail-meta">上传时间：{{ formatDate(detailItem.created_at) }}</p>
        </div>
        <div class="detail-actions">
          <van-button
            type="danger"
            block
            round
            @click="confirmDelete(detailItem)"
          >
            删除衣物
          </van-button>
        </div>
      </div>
    </van-popup>

    <!-- 悬浮上传按钮 -->
    <van-button
      type="primary"
      round
      icon="plus"
      class="fab-button"
      @click="router.push('/upload')"
    >
      上传衣物
    </van-button>
  </div>
</template>

<style scoped>
.wardrobe-page {
  min-height: 100vh;
}
.items-container {
  padding-bottom: 80px;
}
.card-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  padding: 12px;
}
.card-wrapper {
  background: #fff;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  cursor: pointer;
}
.card-image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%;
  background: #fafafa;
}
.card-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
.item-name {
  padding: 8px;
  font-size: 13px;
  color: #333;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.empty-state {
  padding-top: 40px;
}
.loading-state {
  text-align: center;
  padding: 20px;
}
.fab-button {
  position: fixed;
  bottom: 70px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
  width: auto;
  padding: 0 24px;
  box-shadow: 0 4px 12px rgba(25, 137, 250, 0.4);
}
/* Detail popup */
.detail-card {
  background: #fff;
}
.detail-image {
  width: 100%;
  aspect-ratio: 1;
  background: #fafafa;
}
.detail-info {
  padding: 16px 20px 12px;
}
.detail-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px;
}
.detail-meta {
  font-size: 13px;
  color: #888;
  margin: 4px 0;
  line-height: 1.5;
}
.detail-actions {
  padding: 0 20px 20px;
}
</style>
