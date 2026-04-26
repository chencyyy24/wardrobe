<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoryStore } from '@/stores/category'
import { useClothingStore } from '@/stores/clothing'
import { clothingApi } from '@/api/clothing'
import type { Category, Subcategory } from '@/types/category'
import { showToast, showLoadingToast, closeToast } from 'vant'

const router = useRouter()
const categoryStore = useCategoryStore()
const clothingStore = useClothingStore()

// Form data
const imageFile = ref<File | null>(null)
const imagePreview = ref<string>('')
const name = ref('')
const selectedCategory = ref<Category | null>(null)
const selectedSubcategory = ref<Subcategory | null>(null)

// Action sheet
const showCategorySheet = ref(false)
const showSubcategorySheet = ref(false)
const uploading = ref(false)

const subcategories = computed(() => {
  if (!selectedCategory.value) return []
  return selectedCategory.value.subcategories
})

const isValid = computed(() => {
  return imageFile.value && selectedCategory.value
})

onMounted(() => {
  categoryStore.load()
})

function onFileSelect(file: File) {
  imageFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)

  if (!name.value) {
    name.value = file.name.replace(/\.[^/.]+$/, '')
  }
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    onFileSelect(file)
  }
  // 允许重复选择同一文件
  input.value = ''
}

function selectCategory(cat: Category) {
  selectedCategory.value = cat
  selectedSubcategory.value = null
  showCategorySheet.value = false
  if (cat.subcategories.length > 0) {
    showSubcategorySheet.value = true
  }
}

function selectSubcategory(sub: Subcategory) {
  selectedSubcategory.value = sub
  showSubcategorySheet.value = false
}

async function onSubmit() {
  if (!isValid.value || !imageFile.value) {
    showToast('请选择图片和分类')
    return
  }

  uploading.value = true
  showLoadingToast({ message: '上传中...', forbidClick: true })

  try {
    const formData = new FormData()
    formData.append('image', imageFile.value)
    formData.append('name', name.value || imageFile.value.name.replace(/\.[^/.]+$/, ''))
    formData.append('category_id', String(selectedCategory.value!.id))
    if (selectedSubcategory.value) {
      formData.append('subcategory_id', String(selectedSubcategory.value!.id))
    }

    const item = await clothingApi.upload(formData)
    clothingStore.addItem(item)
    closeToast()
    showToast('上传成功')
    router.push('/')
  } catch (e) {
    closeToast()
    showToast('上传失败，请重试')
    console.error(e)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="upload-page">
    <!-- 图片选择 -->
    <div class="section">
      <div class="section-title">衣物照片</div>
      <div class="upload-area">
        <img v-if="imagePreview" :src="imagePreview" class="preview-img" />
        <div v-else class="upload-placeholder">
          <van-icon name="photograph" size="48" color="#c8c9cc" />
          <p>点击选择图片</p>
        </div>
        <!-- file input 覆盖在区域上方，直接接收点击 -->
        <input
          type="file"
          accept="image/*"
          class="file-input-overlay"
          @change="onFileChange"
        />
      </div>
    </div>

    <!-- 名称 -->
    <div class="section">
      <div class="section-title">衣物名称</div>
      <van-field
        v-model="name"
        placeholder="输入衣物名称（可选）"
        clearable
        class="name-field"
      />
    </div>

    <!-- 分类选择 -->
    <div class="section">
      <div class="section-title">分类</div>
      <van-field
        :model-value="selectedCategory?.name || '请选择大类'"
        is-link
        readonly
        placeholder="请选择大类"
        @click="showCategorySheet = true"
      />
      <van-field
        v-if="selectedCategory"
        :model-value="selectedSubcategory?.name || '请选择小类（可选）'"
        is-link
        readonly
        :placeholder="'请选择小类'"
        @click="showSubcategorySheet = true"
      />
    </div>

    <!-- 分类 ActionSheet -->
    <van-action-sheet v-model:show="showCategorySheet" title="选择大类">
      <div class="sheet-list">
        <div
          v-for="cat in categoryStore.categories"
          :key="cat.id"
          class="sheet-item"
          :class="{ active: selectedCategory?.id === cat.id }"
          @click="selectCategory(cat)"
        >
          {{ cat.name }}
        </div>
      </div>
    </van-action-sheet>

    <van-action-sheet v-model:show="showSubcategorySheet" title="选择小类">
      <div class="sheet-list">
        <div
          v-for="sub in subcategories"
          :key="sub.id"
          class="sheet-item"
          :class="{ active: selectedSubcategory?.id === sub.id }"
          @click="selectSubcategory(sub)"
        >
          {{ sub.name }}
        </div>
      </div>
    </van-action-sheet>

    <!-- 提交按钮 -->
    <div class="submit-area">
      <van-button
        type="primary"
        block
        size="large"
        round
        :disabled="!isValid || uploading"
        :loading="uploading"
        @click="onSubmit"
      >
        上传衣物
      </van-button>
    </div>
  </div>
</template>

<style scoped>
.upload-page {
  padding: 16px;
}
.section {
  margin-bottom: 16px;
}
.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}
.upload-area {
  position: relative;
  background: #fff;
  border: 2px dashed #ddd;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  min-height: 180px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.file-input-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
  font-size: 0;
}
.upload-placeholder p {
  color: #bbb;
  font-size: 14px;
  margin-top: 8px;
}
.preview-img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
}
.name-field {
  background: #fff;
  border-radius: 8px;
}
.sheet-list {
  padding: 0 16px 20px;
}
.sheet-item {
  padding: 12px 16px;
  font-size: 15px;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  border-radius: 8px;
}
.sheet-item.active {
  color: #1989fa;
  font-weight: 600;
}
.sheet-item:active {
  background: #f5f5f5;
}
.submit-area {
  margin-top: 32px;
  padding-bottom: 20px;
}
</style>
