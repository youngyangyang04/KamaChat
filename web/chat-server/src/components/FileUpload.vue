<template>
  <div class="file-upload-container">
    <!-- 文件选择按钮 -->
    <el-tooltip 
      effect="customized" 
      content="发送图片" 
      placement="top"
      hide-after="0"
      enterable="false"
    >
      <button 
        class="image-button"
        @click="triggerImageUpload"
        :disabled="uploading"
      >
        <svg
          t="1733503000000"
          class="upload-icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          width="128"
          height="128"
        >
          <path
            d="M896 128H128c-35.2 0-64 28.8-64 64v640c0 35.2 28.8 64 64 64h768c35.2 0 64-28.8 64-64V192c0-35.2-28.8-64-64-64z m0 704H128V192h768v640z"
            fill="#2c2c2c"
          ></path>
          <path
            d="M320 512c53 0 96-43 96-96s-43-96-96-96-96 43-96 96 43 96 96 96z"
            fill="#2c2c2c"
          ></path>
          <path
            d="M896 640l-192-192-256 256-96-96-224 160v64h768z"
            fill="#2c2c2c"
          ></path>
        </svg>
      </button>
    </el-tooltip>

    <el-tooltip 
      effect="customized" 
      content="发送文件" 
      placement="top"
      hide-after="0"
      enterable="false"
    >
      <button 
        class="image-button"
        @click="triggerFileUpload"
        :disabled="uploading"
      >
        <svg
          t="1733503100000"
          class="upload-icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          width="128"
          height="128"
        >
          <path
            d="M854.6 288.6L639.4 73.4c-6-6-14.1-9.4-22.6-9.4H192c-17.7 0-32 14.3-32 32v832c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V311.3c0-8.5-3.4-16.7-9.4-22.7zM790.2 326H602V137.8L790.2 326z m1.8 562H232V136h302v216c0 23.2 18.8 42 42 42h216v494z"
            fill="#2c2c2c"
          ></path>
          <path
            d="M544 472c0-4.4-3.6-8-8-8h-48c-4.4 0-8 3.6-8 8v108H372c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h108v108c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8V644h108c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8H544V472z"
            fill="#2c2c2c"
          ></path>
        </svg>
      </button>
    </el-tooltip>

    <!-- 隐藏的文件输入框 -->
    <input
      ref="imageInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleImageSelect"
    />
    <input
      ref="fileInput"
      type="file"
      style="display: none"
      @change="handleFileSelect"
    />

    <!-- 上传进度弹窗 -->
    <el-dialog
      v-model="showProgress"
      title="正在加密上传..."
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="upload-progress-content">
        <div class="file-info">
          <el-icon v-if="currentFile?.type?.startsWith('image/')" :size="32">
            <Picture />
          </el-icon>
          <el-icon v-else :size="32">
            <Document />
          </el-icon>
          <div class="file-details">
            <span class="file-name">{{ currentFile?.name }}</span>
            <span class="file-size">{{ formatSize(currentFile?.size) }}</span>
          </div>
        </div>
        <el-progress 
          :percentage="uploadProgress" 
          :status="uploadProgress === 100 ? 'success' : ''"
          :stroke-width="20"
        />
        <div class="progress-text">
          <span v-if="uploadProgress < 20">正在加密文件...</span>
          <span v-else-if="uploadProgress < 100">正在上传... {{ uploadProgress }}%</span>
          <span v-else>上传完成！</span>
        </div>
      </div>
    </el-dialog>

    <!-- 图片预览弹窗 -->
    <el-dialog
      v-model="showImagePreview"
      title="发送图片"
      width="500px"
    >
      <div class="image-preview-content">
        <img :src="previewImageUrl" class="preview-image" />
        <div class="preview-info">
          <span>{{ currentFile?.name }}</span>
          <span>{{ formatSize(currentFile?.size) }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="cancelUpload">取消</el-button>
        <el-button type="primary" @click="confirmImageUpload" :loading="uploading">
          发送
        </el-button>
      </template>
    </el-dialog>

    <!-- 文件确认弹窗 -->
    <el-dialog
      v-model="showFileConfirm"
      title="发送文件"
      width="400px"
    >
      <div class="file-confirm-content">
        <el-icon :size="48" color="#409EFF">
          <Document />
        </el-icon>
        <div class="file-info-detail">
          <span class="file-name">{{ currentFile?.name }}</span>
          <span class="file-size">{{ formatSize(currentFile?.size) }}</span>
          <span class="file-type">{{ currentFile?.type || '未知类型' }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="cancelUpload">取消</el-button>
        <el-button type="primary" @click="confirmFileUpload" :loading="uploading">
          发送
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { 
  uploadEncryptedFile, 
  uploadEncryptedImage,
  formatFileSize 
} from '@/crypto/ossUpload';

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  },
  maxImageSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB
  },
  maxFileSize: {
    type: Number,
    default: 100 * 1024 * 1024 // 100MB
  }
});

const emit = defineEmits(['upload-complete', 'upload-error']);

// 引用
const imageInput = ref(null);
const fileInput = ref(null);

// 状态
const uploading = ref(false);
const showProgress = ref(false);
const showImagePreview = ref(false);
const showFileConfirm = ref(false);
const uploadProgress = ref(0);
const currentFile = ref(null);
const previewImageUrl = ref('');

// 格式化文件大小
const formatSize = (bytes) => {
  if (!bytes) return '0 B';
  return formatFileSize(bytes);
};

// 触发图片选择
const triggerImageUpload = () => {
  imageInput.value?.click();
};

// 触发文件选择
const triggerFileUpload = () => {
  fileInput.value?.click();
};

// 处理图片选择
const handleImageSelect = (event) => {
  const file = event.target.files?.[0];
  if (!file) return;

  // 重置输入
  event.target.value = '';

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件');
    return;
  }

  // 检查文件大小
  if (file.size > props.maxImageSize) {
    ElMessage.error(`图片大小不能超过 ${formatSize(props.maxImageSize)}`);
    return;
  }

  currentFile.value = file;
  previewImageUrl.value = URL.createObjectURL(file);
  showImagePreview.value = true;
};

// 处理文件选择
const handleFileSelect = (event) => {
  const file = event.target.files?.[0];
  if (!file) return;

  // 重置输入
  event.target.value = '';

  // 检查文件大小
  if (file.size > props.maxFileSize) {
    ElMessage.error(`文件大小不能超过 ${formatSize(props.maxFileSize)}`);
    return;
  }

  currentFile.value = file;
  showFileConfirm.value = true;
};

// 取消上传
const cancelUpload = () => {
  showImagePreview.value = false;
  showFileConfirm.value = false;
  if (previewImageUrl.value) {
    URL.revokeObjectURL(previewImageUrl.value);
    previewImageUrl.value = '';
  }
  currentFile.value = null;
};

// 确认上传图片
const confirmImageUpload = async () => {
  showImagePreview.value = false;
  await doUpload('image');
};

// 确认上传文件
const confirmFileUpload = async () => {
  showFileConfirm.value = false;
  await doUpload('file');
};

// 执行上传
const doUpload = async (type) => {
  if (!currentFile.value) return;

  uploading.value = true;
  showProgress.value = true;
  uploadProgress.value = 0;

  try {
    let result;
    if (type === 'image') {
      result = await uploadEncryptedImage(currentFile.value, (percent) => {
        uploadProgress.value = percent;
      });
      result.type = 'image';
    } else {
      result = await uploadEncryptedFile(currentFile.value, (percent) => {
        uploadProgress.value = percent;
      });
      result.type = 'file';
    }

    // 通知父组件上传完成
    emit('upload-complete', result);
    ElMessage.success('文件加密上传成功');
  } catch (error) {
    console.error('上传失败:', error);
    ElMessage.error('上传失败: ' + error.message);
    emit('upload-error', error);
  } finally {
    uploading.value = false;
    showProgress.value = false;
    uploadProgress.value = 0;
    if (previewImageUrl.value) {
      URL.revokeObjectURL(previewImageUrl.value);
      previewImageUrl.value = '';
    }
    currentFile.value = null;
  }
};
</script>

<style scoped>
.file-upload-container {
  display: inline-flex;
  align-items: center;
  gap: 0;
}

.image-button {
  width: 30px;
  height: 30px;
  border: none;
  background: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  margin: 0 2px;
}

.image-button:hover {
  background-color: #f5f5f5;
  border-radius: 4px;
}

.image-button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.upload-icon {
  width: 22px;
  height: 22px;
}

.upload-progress-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-details {
  display: flex;
  flex-direction: column;
}

.file-name {
  font-weight: 500;
  word-break: break-all;
}

.file-size {
  color: #909399;
  font-size: 12px;
}

.progress-text {
  text-align: center;
  color: #606266;
}

.image-preview-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.preview-image {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  border-radius: 8px;
}

.preview-info {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 14px;
}

.file-confirm-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 20px;
}

.file-info-detail {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.file-info-detail .file-name {
  font-size: 16px;
  font-weight: 500;
  word-break: break-all;
  text-align: center;
}

.file-info-detail .file-size {
  color: #909399;
}

.file-info-detail .file-type {
  color: #C0C4CC;
  font-size: 12px;
}
</style>
