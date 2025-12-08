<template>
  <div class="file-message" :class="{ 'is-image': isImage, 'is-file': !isImage }">
    <!-- 图片消息 -->
    <div v-if="isImage" class="image-message">
      <!-- 缩略图/加载中状态 -->
      <div 
        v-if="!imageLoaded" 
        class="image-placeholder"
        :style="{ width: placeholderWidth + 'px', height: placeholderHeight + 'px' }"
      >
        <img 
          v-if="thumbnailUrl" 
          :src="thumbnailUrl" 
          class="thumbnail-image"
          @click="loadFullImage"
        />
        <div v-else class="loading-placeholder" @click="loadFullImage">
          <el-icon class="is-loading" :size="24">
            <Loading />
          </el-icon>
          <span>点击加载图片</span>
        </div>
      </div>
      
      <!-- 完整图片 -->
      <img 
        v-else
        :src="fullImageUrl" 
        class="full-image"
        @click="previewImage"
      />
      
      <!-- 加密标识 -->
      <div class="encrypt-badge" v-if="isEncrypted">
        <el-icon :size="12"><Lock /></el-icon>
      </div>
      
      <!-- 下载进度 -->
      <div v-if="downloading" class="download-progress">
        <el-progress 
          type="circle" 
          :percentage="downloadProgress" 
          :width="60"
        />
      </div>
    </div>

    <!-- 文件消息 -->
    <div v-else class="file-card" @click="downloadFile">
      <div class="file-icon">
        <el-icon :size="36">
          <Document />
        </el-icon>
      </div>
      <div class="file-info">
        <span class="file-name" :title="fileInfo.fileName">
          {{ truncateFileName(fileInfo.fileName) }}
        </span>
        <span class="file-size">{{ formatSize(fileInfo.fileSize) }}</span>
      </div>
      <div class="file-action">
        <el-icon v-if="downloading" class="is-loading" :size="20">
          <Loading />
        </el-icon>
        <el-icon v-else :size="20">
          <Download />
        </el-icon>
      </div>
      <!-- 加密标识 -->
      <div class="encrypt-badge" v-if="isEncrypted">
        <el-icon :size="12"><Lock /></el-icon>
      </div>
    </div>

    <!-- 图片预览对话框 -->
    <el-dialog
      v-model="showPreview"
      :title="fileInfo.fileName"
      width="80%"
      class="image-preview-dialog"
    >
      <div class="preview-container">
        <img :src="fullImageUrl" class="preview-full-image" />
      </div>
      <template #footer>
        <el-button @click="saveImage">保存图片</el-button>
        <el-button @click="showPreview = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { Document, Download, Loading, Lock } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { 
  downloadAndDecryptFileAsURL, 
  downloadAndSaveFile,
  formatFileSize
} from '@/crypto/ossUpload';
import { decryptThumbnail } from '@/crypto/fileEncryption';

const props = defineProps({
  // 文件信息对象
  fileInfo: {
    type: Object,
    required: true
    // 结构: {
    //   type: 'image' | 'file',
    //   ossKey: string,
    //   fileKey: string,
    //   fileIv: string,
    //   fileHash: string,
    //   fileName: string,
    //   fileSize: number,
    //   mimeType: string,
    //   width?: number,
    //   height?: number,
    //   thumbnail?: { data, key, iv, width, height }
    // }
  },
  isEncrypted: {
    type: Boolean,
    default: true
  }
});

// 状态
const imageLoaded = ref(false);
const fullImageUrl = ref('');
const thumbnailUrl = ref('');
const downloading = ref(false);
const downloadProgress = ref(0);
const showPreview = ref(false);

// 计算属性
const isImage = computed(() => {
  return props.fileInfo?.type === 'image' || 
         props.fileInfo?.mimeType?.startsWith('image/');
});

const placeholderWidth = computed(() => {
  if (props.fileInfo?.thumbnail?.width) {
    return Math.min(props.fileInfo.thumbnail.width, 200);
  }
  if (props.fileInfo?.width) {
    return Math.min(props.fileInfo.width, 200);
  }
  return 150;
});

const placeholderHeight = computed(() => {
  if (props.fileInfo?.thumbnail?.height) {
    return Math.min(props.fileInfo.thumbnail.height, 200);
  }
  if (props.fileInfo?.height && props.fileInfo?.width) {
    const ratio = props.fileInfo.height / props.fileInfo.width;
    return Math.min(placeholderWidth.value * ratio, 200);
  }
  return 150;
});

// 格式化文件大小
const formatSize = (bytes) => {
  if (!bytes) return '未知大小';
  return formatFileSize(bytes);
};

// 截断文件名
const truncateFileName = (name) => {
  if (!name) return '未知文件';
  if (name.length <= 20) return name;
  const ext = name.split('.').pop();
  const baseName = name.substring(0, name.length - ext.length - 1);
  return baseName.substring(0, 15) + '...' + ext;
};

// 加载缩略图
const loadThumbnail = async () => {
  if (!props.fileInfo?.thumbnail) return;
  
  try {
    const thumb = props.fileInfo.thumbnail;
    const blob = await decryptThumbnail(thumb.data, thumb.key, thumb.iv);
    thumbnailUrl.value = URL.createObjectURL(blob);
  } catch (error) {
    console.error('加载缩略图失败:', error);
  }
};

// 加载完整图片
const loadFullImage = async () => {
  if (downloading.value || imageLoaded.value) return;
  
  downloading.value = true;
  downloadProgress.value = 0;
  
  try {
    fullImageUrl.value = await downloadAndDecryptFileAsURL(props.fileInfo, (percent) => {
      downloadProgress.value = percent;
    });
    imageLoaded.value = true;
  } catch (error) {
    console.error('加载图片失败:', error);
    ElMessage.error('加载图片失败: ' + error.message);
  } finally {
    downloading.value = false;
  }
};

// 预览图片
const previewImage = () => {
  if (imageLoaded.value) {
    showPreview.value = true;
  }
};

// 保存图片
const saveImage = async () => {
  try {
    await downloadAndSaveFile(props.fileInfo);
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message);
  }
};

// 下载文件
const downloadFile = async () => {
  if (downloading.value) return;
  
  downloading.value = true;
  downloadProgress.value = 0;
  
  try {
    await downloadAndSaveFile(props.fileInfo, (percent) => {
      downloadProgress.value = percent;
    });
    ElMessage.success('文件下载完成');
  } catch (error) {
    console.error('下载文件失败:', error);
    ElMessage.error('下载失败: ' + error.message);
  } finally {
    downloading.value = false;
  }
};

// 生命周期
onMounted(() => {
  if (isImage.value) {
    loadThumbnail();
  }
});

onUnmounted(() => {
  // 清理 Blob URL
  if (fullImageUrl.value) {
    URL.revokeObjectURL(fullImageUrl.value);
  }
  if (thumbnailUrl.value) {
    URL.revokeObjectURL(thumbnailUrl.value);
  }
});
</script>

<style scoped>
.file-message {
  position: relative;
  display: inline-block;
  max-width: 300px;
}

/* 图片消息样式 */
.image-message {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  min-width: 100px;
  min-height: 100px;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: blur(2px);
  transition: filter 0.3s;
}

.thumbnail-image:hover {
  filter: blur(0);
}

.loading-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #909399;
  font-size: 12px;
}

.full-image {
  max-width: 100%;
  max-height: 300px;
  cursor: pointer;
  transition: transform 0.2s;
}

.full-image:hover {
  transform: scale(1.02);
}

.download-progress {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  padding: 5px;
}

/* 文件消息样式 */
.file-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  min-width: 200px;
}

.file-card:hover {
  background: #e8eaed;
}

.file-icon {
  color: #409EFF;
  flex-shrink: 0;
}

.file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-info .file-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-info .file-size {
  font-size: 12px;
  color: #909399;
}

.file-action {
  color: #606266;
  flex-shrink: 0;
}

/* 加密标识 */
.encrypt-badge {
  position: absolute;
  bottom: 4px;
  right: 4px;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  border-radius: 4px;
  padding: 2px 4px;
  display: flex;
  align-items: center;
}

.file-card .encrypt-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  bottom: auto;
}

/* 预览对话框 */
.image-preview-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.preview-container {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000;
  max-height: 70vh;
  overflow: auto;
}

.preview-full-image {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

/* 加载动画 */
.is-loading {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
