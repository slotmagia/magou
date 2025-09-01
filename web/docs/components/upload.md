# BasicUpload 上传组件

BasicUpload 是基于 Naive UI Upload 组件封装的高级文件上传组件，支持图片上传、文件上传、拖拽上传、多文件上传等功能。

## 组件特性

### 🎯 核心功能
- **多种上传方式** - 支持点击上传、拖拽上传
- **文件类型检测** - 支持图片、文档、视频等文件类型
- **文件预览** - 图片预览、文件信息展示
- **上传进度** - 实时显示上传进度
- **文件大小限制** - 可配置文件大小限制
- **批量上传** - 支持多文件同时上传

### 📁 支持的文件类型
- **图片**: jpg, jpeg, png, gif, bmp, webp
- **文档**: pdf, doc, docx, xls, xlsx, ppt, pptx, txt
- **压缩包**: zip, rar, 7z
- **视频**: mp4, avi, mov, wmv
- **音频**: mp3, wav, flac

## 基础用法

### 图片上传

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>单图上传</h3>
      <BasicUpload
        :max-size="5"
        :max-number="1"
        file-type="image"
        @upload-change="handleImageChange"
        v-model:value="imageUrl"
      />
    </div>
    
    <div class="mb-4">
      <h3>多图上传</h3>
      <BasicUpload
        :max-size="10"
        :max-number="5"
        file-type="image"
        @upload-change="handleImagesChange"
        v-model:value="imageList"
      />
    </div>
    
    <div class="mt-4">
      <p>单图URL: {{ imageUrl }}</p>
      <p>多图列表: {{ imageList }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicUpload } from '@/components/Upload';

const imageUrl = ref('');
const imageList = ref<string[]>([]);

// 单图上传回调
const handleImageChange = (fileList: string[]) => {
  console.log('图片变化:', fileList);
  if (fileList.length > 0) {
    imageUrl.value = fileList[0];
  } else {
    imageUrl.value = '';
  }
};

// 多图上传回调
const handleImagesChange = (fileList: string[]) => {
  console.log('图片列表变化:', fileList);
  imageList.value = fileList;
};
</script>
```

### 文件上传

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>文档上传</h3>
      <BasicUpload
        :max-size="20"
        :max-number="10"
        file-type="file"
        :accept="['pdf', 'doc', 'docx', 'xls', 'xlsx']"
        @upload-change="handleFileChange"
        v-model:value="fileList"
      />
    </div>
    
    <div class="mb-4">
      <h3>所有类型文件上传</h3>
      <BasicUpload
        :max-size="50"
        :max-number="20"
        file-type="all"
        @upload-change="handleAllFileChange"
        v-model:value="allFileList"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicUpload } from '@/components/Upload';

const fileList = ref<string[]>([]);
const allFileList = ref<string[]>([]);

const handleFileChange = (files: string[]) => {
  console.log('文档变化:', files);
  fileList.value = files;
};

const handleAllFileChange = (files: string[]) => {
  console.log('所有文件变化:', files);
  allFileList.value = files;
};
</script>
```

## 高级用法

### 自定义上传逻辑

```vue
<template>
  <div class="p-4">
    <BasicUpload
      :max-size="10"
      :max-number="5"
      file-type="image"
      :custom-request="customUpload"
      @upload-change="handleChange"
      @upload-before="handleBefore"
      @upload-success="handleSuccess"
      @upload-error="handleError"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicUpload } from '@/components/Upload';
import { uploadFile } from '@/api/upload';

// 自定义上传函数
const customUpload = async (file: File, onProgress?: (percent: number) => void) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('folder', 'images'); // 指定上传文件夹
    
    const response = await uploadFile(formData, {
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(percent);
        }
      },
    });
    
    return response.data.url; // 返回文件URL
  } catch (error) {
    console.error('上传失败:', error);
    throw error;
  }
};

// 上传前检查
const handleBefore = (file: File) => {
  console.log('上传前检查:', file);
  
  // 检查文件大小
  if (file.size > 10 * 1024 * 1024) {
    window.$message.error('文件大小不能超过10MB');
    return false;
  }
  
  // 检查文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
  if (!allowedTypes.includes(file.type)) {
    window.$message.error('只支持 JPEG、PNG、GIF 格式的图片');
    return false;
  }
  
  return true;
};

// 上传成功
const handleSuccess = (fileUrl: string, file: File) => {
  console.log('上传成功:', fileUrl, file);
  window.$message.success(`${file.name} 上传成功`);
};

// 上传失败
const handleError = (error: Error, file: File) => {
  console.error('上传失败:', error, file);
  window.$message.error(`${file.name} 上传失败: ${error.message}`);
};

// 文件列表变化
const handleChange = (fileList: string[]) => {
  console.log('文件列表变化:', fileList);
};
</script>
```

### 拖拽上传

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>拖拽上传区域</h3>
      <BasicUpload
        :max-size="20"
        :max-number="10"
        file-type="all"
        :drag-upload="true"
        :show-file-list="true"
        @upload-change="handleDragUpload"
      >
        <template #drag-content>
          <div class="text-center py-8">
            <n-icon size="48" class="text-gray-400 mb-4">
              <CloudUploadOutlined />
            </n-icon>
            <div class="text-lg font-medium mb-2">拖拽文件到此区域上传</div>
            <div class="text-gray-500">支持单个或批量上传，严禁上传公司数据或其他违禁文件</div>
          </div>
        </template>
      </BasicUpload>
    </div>
    
    <div v-if="uploadedFiles.length > 0">
      <h4>已上传文件:</h4>
      <ul class="list-disc pl-6">
        <li v-for="file in uploadedFiles" :key="file">
          <a :href="file" target="_blank" class="text-blue-500 hover:underline">
            {{ getFileName(file) }}
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { CloudUploadOutlined } from '@vicons/antd';
import { BasicUpload } from '@/components/Upload';

const uploadedFiles = ref<string[]>([]);

const handleDragUpload = (fileList: string[]) => {
  uploadedFiles.value = fileList;
};

const getFileName = (url: string) => {
  return url.split('/').pop() || 'unknown';
};
</script>
```

### 带预览的图片上传

```vue
<template>
  <div class="p-4">
    <BasicUpload
      :max-size="5"
      :max-number="9"
      file-type="image"
      :preview="true"
      :show-preview-list="true"
      v-model:value="previewImages"
      @preview="handlePreview"
    />
    
    <!-- 图片预览模态框 -->
    <n-modal v-model:show="previewVisible" preset="card" title="图片预览">
      <img :src="previewUrl" alt="预览图片" class="w-full" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicUpload } from '@/components/Upload';

const previewImages = ref<string[]>([]);
const previewVisible = ref(false);
const previewUrl = ref('');

const handlePreview = (url: string) => {
  previewUrl.value = url;
  previewVisible.value = true;
};
</script>
```

## API 接口

### BasicUpload Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string \| string[]` | - | 绑定值，文件URL或URL数组 |
| fileType | `'image' \| 'file' \| 'all'` | `'image'` | 文件类型 |
| maxSize | `number` | `10` | 最大文件大小(MB) |
| maxNumber | `number` | `1` | 最大文件数量 |
| accept | `string[]` | - | 允许的文件类型 |
| width | `number \| string` | `104` | 上传区域宽度 |
| height | `number \| string` | `104` | 上传区域高度 |
| dragUpload | `boolean` | `false` | 是否支持拖拽上传 |
| showFileList | `boolean` | `true` | 是否显示文件列表 |
| preview | `boolean` | `true` | 是否支持预览 |
| showPreviewList | `boolean` | `false` | 是否显示预览列表 |
| customRequest | `Function` | - | 自定义上传函数 |
| disabled | `boolean` | `false` | 是否禁用 |

### 事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: string \| string[])` | 绑定值更新 |
| uploadChange | `(fileList: string[])` | 文件列表变化 |
| uploadBefore | `(file: File)` | 上传前检查 |
| uploadSuccess | `(url: string, file: File)` | 上传成功 |
| uploadError | `(error: Error, file: File)` | 上传失败 |
| preview | `(url: string)` | 预览文件 |
| remove | `(url: string, index: number)` | 删除文件 |

### 插槽

| 插槽名 | 说明 |
|--------|------|
| default | 自定义上传区域内容 |
| drag-content | 自定义拖拽区域内容 |
| file-list | 自定义文件列表 |

## 上传配置

### 文件类型配置

```typescript
// 图片类型
const imageTypes = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'];

// 文档类型
const documentTypes = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt'];

// 压缩包类型
const archiveTypes = ['zip', 'rar', '7z', 'tar', 'gz'];

// 视频类型
const videoTypes = ['mp4', 'avi', 'mov', 'wmv', 'flv', 'mkv'];

// 音频类型
const audioTypes = ['mp3', 'wav', 'flac', 'aac', 'ogg'];
```

### 上传接口配置

```typescript
// api/upload.ts
import { http } from '@/utils/http';

export interface UploadResponse {
  url: string;
  filename: string;
  size: number;
  type: string;
}

/**
 * 单文件上传
 */
export function uploadFile(
  formData: FormData,
  config?: AxiosRequestConfig
): Promise<UploadResponse> {
  return http.post({
    url: '/upload/file',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    ...config,
  });
}

/**
 * 多文件上传
 */
export function uploadFiles(
  formData: FormData,
  config?: AxiosRequestConfig
): Promise<UploadResponse[]> {
  return http.post({
    url: '/upload/files',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    ...config,
  });
}

/**
 * 图片上传（带压缩）
 */
export function uploadImage(
  file: File,
  options: {
    quality?: number; // 压缩质量 0-1
    maxWidth?: number; // 最大宽度
    maxHeight?: number; // 最大高度
  } = {}
): Promise<UploadResponse> {
  return new Promise((resolve, reject) => {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    const img = new Image();
    
    img.onload = () => {
      const { quality = 0.8, maxWidth = 1920, maxHeight = 1080 } = options;
      
      // 计算新尺寸
      let { width, height } = img;
      if (width > maxWidth) {
        height = (height * maxWidth) / width;
        width = maxWidth;
      }
      if (height > maxHeight) {
        width = (width * maxHeight) / height;
        height = maxHeight;
      }
      
      canvas.width = width;
      canvas.height = height;
      
      // 绘制并压缩
      ctx?.drawImage(img, 0, 0, width, height);
      canvas.toBlob(
        (blob) => {
          if (blob) {
            const formData = new FormData();
            formData.append('file', blob, file.name);
            uploadFile(formData).then(resolve).catch(reject);
          } else {
            reject(new Error('图片压缩失败'));
          }
        },
        'image/jpeg',
        quality
      );
    };
    
    img.onerror = () => reject(new Error('图片加载失败'));
    img.src = URL.createObjectURL(file);
  });
}
```

## 实用工具

### 文件处理工具

```typescript
// utils/fileUtils.ts

/**
 * 获取文件扩展名
 */
export function getFileExtension(filename: string): string {
  return filename.split('.').pop()?.toLowerCase() || '';
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * 检查文件类型
 */
export function checkFileType(file: File, allowedTypes: string[]): boolean {
  const extension = getFileExtension(file.name);
  return allowedTypes.includes(extension);
}

/**
 * 生成文件预览URL
 */
export function createPreviewUrl(file: File): string {
  return URL.createObjectURL(file);
}

/**
 * 清理预览URL
 */
export function revokePreviewUrl(url: string): void {
  URL.revokeObjectURL(url);
}

/**
 * 文件转Base64
 */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * Base64转文件
 */
export function base64ToFile(base64: string, filename: string): File {
  const arr = base64.split(',');
  const mime = arr[0].match(/:(.*?);/)?.[1] || '';
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  
  return new File([u8arr], filename, { type: mime });
}
```

### 上传进度管理

```typescript
// composables/useUploadProgress.ts
export function useUploadProgress() {
  const progressMap = ref(new Map<string, number>());
  
  const setProgress = (fileId: string, progress: number) => {
    progressMap.value.set(fileId, progress);
  };
  
  const getProgress = (fileId: string): number => {
    return progressMap.value.get(fileId) || 0;
  };
  
  const removeProgress = (fileId: string) => {
    progressMap.value.delete(fileId);
  };
  
  const clearProgress = () => {
    progressMap.value.clear();
  };
  
  return {
    progressMap: readonly(progressMap),
    setProgress,
    getProgress,
    removeProgress,
    clearProgress,
  };
}
```

## 最佳实践

### 1. 文件上传优化

```typescript
// 分片上传大文件
export class ChunkUploader {
  private chunkSize = 1024 * 1024; // 1MB
  
  async uploadLargeFile(file: File, onProgress?: (percent: number) => void) {
    const chunks = this.createChunks(file);
    const uploadedChunks: string[] = [];
    
    for (let i = 0; i < chunks.length; i++) {
      const chunk = chunks[i];
      const chunkData = new FormData();
      chunkData.append('chunk', chunk);
      chunkData.append('chunkIndex', i.toString());
      chunkData.append('totalChunks', chunks.length.toString());
      chunkData.append('fileName', file.name);
      
      try {
        const response = await uploadChunk(chunkData);
        uploadedChunks.push(response.chunkId);
        
        if (onProgress) {
          onProgress(Math.round(((i + 1) / chunks.length) * 100));
        }
      } catch (error) {
        throw new Error(`分片 ${i} 上传失败: ${error.message}`);
      }
    }
    
    // 合并分片
    return await mergeChunks({
      fileName: file.name,
      chunkIds: uploadedChunks,
    });
  }
  
  private createChunks(file: File): Blob[] {
    const chunks: Blob[] = [];
    let start = 0;
    
    while (start < file.size) {
      const end = Math.min(start + this.chunkSize, file.size);
      chunks.push(file.slice(start, end));
      start = end;
    }
    
    return chunks;
  }
}
```

### 2. 上传队列管理

```typescript
// composables/useUploadQueue.ts
export function useUploadQueue(maxConcurrent = 3) {
  const queue = ref<UploadTask[]>([]);
  const running = ref<UploadTask[]>([]);
  const completed = ref<UploadTask[]>([]);
  const failed = ref<UploadTask[]>([]);
  
  const addTask = (task: UploadTask) => {
    queue.value.push(task);
    processQueue();
  };
  
  const processQueue = async () => {
    while (queue.value.length > 0 && running.value.length < maxConcurrent) {
      const task = queue.value.shift();
      if (task) {
        running.value.push(task);
        executeTask(task);
      }
    }
  };
  
  const executeTask = async (task: UploadTask) => {
    try {
      task.status = 'uploading';
      await task.upload();
      task.status = 'completed';
      completed.value.push(task);
    } catch (error) {
      task.status = 'failed';
      task.error = error;
      failed.value.push(task);
    } finally {
      const index = running.value.findIndex(t => t.id === task.id);
      if (index > -1) {
        running.value.splice(index, 1);
      }
      processQueue();
    }
  };
  
  return {
    queue: readonly(queue),
    running: readonly(running),
    completed: readonly(completed),
    failed: readonly(failed),
    addTask,
  };
}
```

### 3. 文件安全检查

```typescript
// utils/fileSecurity.ts
export class FileSecurityChecker {
  private dangerousExtensions = [
    'exe', 'bat', 'cmd', 'com', 'pif', 'scr', 'vbs', 'js', 'jar',
  ];
  
  private maxFileSize = 100 * 1024 * 1024; // 100MB
  
  checkFile(file: File): { valid: boolean; message?: string } {
    // 检查文件扩展名
    const extension = getFileExtension(file.name);
    if (this.dangerousExtensions.includes(extension)) {
      return {
        valid: false,
        message: `不允许上传 ${extension.toUpperCase()} 文件`,
      };
    }
    
    // 检查文件大小
    if (file.size > this.maxFileSize) {
      return {
        valid: false,
        message: `文件大小不能超过 ${formatFileSize(this.maxFileSize)}`,
      };
    }
    
    // 检查MIME类型
    if (!this.isValidMimeType(file)) {
      return {
        valid: false,
        message: '文件类型不匹配',
      };
    }
    
    return { valid: true };
  }
  
  private isValidMimeType(file: File): boolean {
    const extension = getFileExtension(file.name);
    const expectedMimeTypes = this.getExpectedMimeTypes(extension);
    
    return expectedMimeTypes.includes(file.type);
  }
  
  private getExpectedMimeTypes(extension: string): string[] {
    const mimeMap: Record<string, string[]> = {
      jpg: ['image/jpeg'],
      jpeg: ['image/jpeg'],
      png: ['image/png'],
      gif: ['image/gif'],
      pdf: ['application/pdf'],
      // ... 更多映射
    };
    
    return mimeMap[extension] || [];
  }
}
```

---

下一步：[其他组件](./others.md)






