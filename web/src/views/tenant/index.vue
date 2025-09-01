<template>
  <div>
    <n-card :bordered="false" class="proCard" title="租户管理">
      <n-space vertical :size="12">
        <n-space>
          <n-button type="primary" @click="addTable">
            <template #icon>
              <n-icon>
                <PlusOutlined />
              </n-icon>
            </template>
            添加租户
          </n-button>
          <n-button type="info" @click="refreshTable">
            <template #icon>
              <n-icon>
                <ReloadOutlined />
              </n-icon>
            </template>
            刷新
          </n-button>
        </n-space>

        <BasicTable
          :columns="columns"
          :request="loadDataTable"
          :row-key="(row) => row.id"
          ref="actionRef"
          :actionColumn="actionColumn"
          :scroll-x="1400"
          :resizeHeightOffset="-10000"
        >
          <template #tableTitle>
            <n-alert type="info" style="margin-bottom: 16px;">
              租户管理用于管理系统中的所有租户，包括创建、编辑、删除租户以及配置租户相关权限。
            </n-alert>
          </template>
        </BasicTable>
      </n-space>
    </n-card>

    <!-- 编辑租户弹窗 -->
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      :title="formParams.id > 0 ? '编辑租户 #' + formParams.id : '添加租户'"
      :style="{ width: dialogWidth }"
    >
      <n-scrollbar style="max-height: 87vh" class="pr-5">
        <n-spin :show="loading" description="请稍候...">
          <n-form
            :model="formParams"
            :rules="rules"
            ref="formRef"
            label-placement="left"
            :label-width="120"
            class="py-4"
          >
            <n-form-item label="租户名称" path="name">
              <n-input placeholder="请输入租户名称" v-model:value="formParams.name" />
            </n-form-item>
            <n-form-item label="租户编码" path="code">
              <n-input placeholder="请输入租户编码（全局唯一）" v-model:value="formParams.code" :disabled="formParams.id > 0" />
            </n-form-item>
            <n-form-item label="租户域名" path="domain">
              <n-input placeholder="请输入租户域名（可选）" v-model:value="formParams.domain" />
            </n-form-item>
            <n-form-item label="最大用户数" path="maxUsers">
              <n-input-number v-model:value="formParams.maxUsers" :min="1" :max="10000" style="width: 100%" />
            </n-form-item>
            <n-form-item label="存储限制" path="storageLimit">
              <n-input-number v-model:value="formParams.storageLimit" :min="0" style="width: 100%" />
              <template #feedback>
                单位：字节，1GB = 1073741824字节
              </template>
            </n-form-item>
            <n-form-item label="过期时间" path="expireAt">
              <n-date-picker
                v-model:value="formParams.expireAt"
                type="datetime"
                clearable
                style="width: 100%"
                placeholder="请选择过期时间（可选）"
              />
            </n-form-item>
            <template v-if="formParams.id === 0">
              <n-form-item label="管理员用户名" path="adminName">
                <n-input placeholder="请输入管理员用户名" v-model:value="formParams.adminName" />
              </n-form-item>
              <n-form-item label="管理员邮箱" path="adminEmail">
                <n-input placeholder="请输入管理员邮箱" v-model:value="formParams.adminEmail" />
              </n-form-item>
              <n-form-item label="管理员密码" path="adminPassword">
                <n-input type="password" placeholder="请输入管理员密码" v-model:value="formParams.adminPassword" />
              </n-form-item>
            </template>
            <n-form-item label="备注" path="remark">
              <n-input type="textarea" placeholder="请输入备注信息" v-model:value="formParams.remark" />
            </n-form-item>
          </n-form>
        </n-spin>
      </n-scrollbar>
      <template #action>
        <n-space>
          <n-button @click="closeForm">取消</n-button>
          <n-button type="primary" :loading="formBtnLoading" @click="confirmForm">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 租户详情弹窗 -->
    <n-modal
      v-model:show="showDetailModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      title="租户详情"
      :style="{ width: dialogWidth }"
    >
      <n-scrollbar style="max-height: 87vh" class="pr-5">
        <n-spin :show="detailLoading" description="加载中...">
          <n-descriptions bordered :column="2" v-if="tenantDetail">
            <n-descriptions-item label="租户ID">{{ tenantDetail.id }}</n-descriptions-item>
            <n-descriptions-item label="租户名称">{{ tenantDetail.name }}</n-descriptions-item>
            <n-descriptions-item label="租户编码">{{ tenantDetail.code }}</n-descriptions-item>
            <n-descriptions-item label="租户域名">{{ tenantDetail.domain || '未设置' }}</n-descriptions-item>
            <n-descriptions-item label="状态">
              <n-tag :type="tenantDetail.status === 1 ? 'success' : 'error'">
                {{ tenantDetail.statusName }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="最大用户数">{{ tenantDetail.maxUsers }}</n-descriptions-item>
            <n-descriptions-item label="存储限制">{{ formatFileSize(tenantDetail.storageLimit) }}</n-descriptions-item>
            <n-descriptions-item label="过期时间">{{ tenantDetail.expireAt || '永不过期' }}</n-descriptions-item>
            <n-descriptions-item label="管理员ID">{{ tenantDetail.adminUserId }}</n-descriptions-item>
            <n-descriptions-item label="管理员名称">{{ tenantDetail.adminName }}</n-descriptions-item>
            <n-descriptions-item label="创建时间">{{ tenantDetail.createdAt }}</n-descriptions-item>
            <n-descriptions-item label="更新时间">{{ tenantDetail.updatedAt }}</n-descriptions-item>
            <n-descriptions-item label="备注" :span="2">{{ tenantDetail.remark || '无' }}</n-descriptions-item>
            
            <!-- 统计信息 -->
            <template v-if="tenantDetail.stats">
              <n-descriptions-item label="用户数量">{{ tenantDetail.stats.userCount }}</n-descriptions-item>
              <n-descriptions-item label="角色数量">{{ tenantDetail.stats.roleCount }}</n-descriptions-item>
              <n-descriptions-item label="菜单数量">{{ tenantDetail.stats.menuCount }}</n-descriptions-item>
              <n-descriptions-item label="存储使用">{{ tenantDetail.stats.storageUsedText }}</n-descriptions-item>
              <n-descriptions-item label="最后活跃时间" :span="2">{{ tenantDetail.stats.lastActiveTime }}</n-descriptions-item>
            </template>
          </n-descriptions>
        </n-spin>
      </n-scrollbar>
      <template #action>
        <n-space>
          <n-button @click="showDetailModal = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { h, onMounted, reactive, ref, computed } from 'vue';
  import { NButton, NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { PlusOutlined, ReloadOutlined, EyeOutlined, EditOutlined, DeleteOutlined } from '@vicons/antd';
  import { 
    getTenantList, 
    createTenant, 
    updateTenant, 
    deleteTenant, 
    getTenantDetail,
    updateTenantStatus,
    type TenantInfo,
    type CreateTenantParams,
    type UpdateTenantParams
  } from '@/api/tenant';
  import { adaModalWidth } from '@/utils/hotgo';
  import { formatFileSize } from '@/utils';

  const message = useMessage();
  const dialog = useDialog();
  const loading = ref(false);
  const detailLoading = ref(false);
  const showModal = ref(false);
  const showDetailModal = ref(false);
  const formBtnLoading = ref(false);
  const actionRef = ref();
  const formRef = ref();
  const tenantDetail = ref<TenantInfo | null>(null);

  const dialogWidth = computed(() => {
    return adaModalWidth(800);
  });

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      key: 'id',
      width: 80,
    },
    {
      title: '租户名称',
      key: 'name',
      width: 150,
      render(row: TenantInfo) {
        return h(NTag, { type: 'info' }, { default: () => row.name });
      },
    },
    {
      title: '租户编码',
      key: 'code',
      width: 120,
    },
    {
      title: '域名',
      key: 'domain',
      width: 180,
      render(row: TenantInfo) {
        return row.domain || '未设置';
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render(row: TenantInfo) {
        return h(
          NTag,
          {
            type: row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'error',
          },
          { default: () => row.statusName }
        );
      },
    },
    {
      title: '最大用户数',
      key: 'maxUsers',
      width: 120,
    },
    {
      title: '存储限制',
      key: 'storageLimit',
      width: 120,
      render(row: TenantInfo) {
        return formatFileSize(row.storageLimit);
      },
    },
    {
      title: '过期时间',
      key: 'expireAt',
      width: 180,
      render(row: TenantInfo) {
        return row.expireAt || '永不过期';
      },
    },
    {
      title: '管理员',
      key: 'adminName',
      width: 120,
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: 180,
    },
  ];

  // 操作列
  const actionColumn = reactive({
    width: 240,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: TenantInfo) {
      return h(TableAction, {
        style: 'button',
        actions: [
          {
            label: '详情',
            icon: EyeOutlined,
            onClick: handleDetail.bind(null, record),
          },
          {
            label: '编辑',
            icon: EditOutlined,
            onClick: handleEdit.bind(null, record),
            ifShow: () => record.id !== 1, // 系统租户不可编辑
          },
          {
            label: record.status === 1 ? '禁用' : '启用',
            onClick: handleToggleStatus.bind(null, record),
            ifShow: () => record.id !== 1, // 系统租户不可禁用
          },
          {
            label: '删除',
            icon: DeleteOutlined,
            onClick: handleDelete.bind(null, record),
            ifShow: () => record.id !== 1, // 系统租户不可删除
          },
        ],
      });
    },
  });

  // 表单数据
  const defaultFormParams: CreateTenantParams = {
    name: '',
    code: '',
    domain: '',
    maxUsers: 100,
    storageLimit: 1073741824, // 1GB
    expireAt: '',
    adminName: '',
    adminEmail: '',
    adminPassword: '',
    remark: '',
  };

  const formParams = ref<CreateTenantParams & { id?: number }>({ ...defaultFormParams });

  // 表单验证规则
  const rules = {
    name: {
      required: true,
      trigger: ['blur', 'input'],
      message: '请输入租户名称',
    },
    code: {
      required: true,
      trigger: ['blur', 'input'],
      message: '请输入租户编码',
    },
    maxUsers: {
      required: true,
      type: 'number',
      trigger: ['blur', 'change'],
      message: '请输入最大用户数',
    },
    storageLimit: {
      required: true,
      type: 'number',
      trigger: ['blur', 'change'],
      message: '请输入存储限制',
    },
    adminName: {
      required: true,
      trigger: ['blur', 'input'],
      message: '请输入管理员用户名',
      validator: (rule: any, value: string) => {
        if (formParams.value.id === 0 && !value) {
          return new Error('请输入管理员用户名');
        }
        return true;
      },
    },
    adminEmail: {
      required: true,
      trigger: ['blur', 'input'],
      message: '请输入管理员邮箱',
      validator: (rule: any, value: string) => {
        if (formParams.value.id === 0 && !value) {
          return new Error('请输入管理员邮箱');
        }
        if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
          return new Error('请输入正确的邮箱格式');
        }
        return true;
      },
    },
    adminPassword: {
      required: true,
      trigger: ['blur', 'input'],
      message: '请输入管理员密码',
      validator: (rule: any, value: string) => {
        if (formParams.value.id === 0 && !value) {
          return new Error('请输入管理员密码');
        }
        if (value && value.length < 6) {
          return new Error('密码至少6位');
        }
        return true;
      },
    },
  };

  // 加载数据表格
  const loadDataTable = async (params: any) => {
    try {
      const response = await getTenantList({
        page: params.page || 1,
        pageSize: params.pageSize || 20,
        name: params.name,
        code: params.code,
        domain: params.domain,
        status: params.status,
      });

      if (response.code === 0) {
        return response.data;
      }
      return { list: [], total: 0 };
    } catch (error) {
      console.error('获取租户列表失败:', error);
      return { list: [], total: 0 };
    }
  };

  // 刷新表格
  function refreshTable() {
    actionRef.value?.reload();
  }

  // 添加租户
  function addTable() {
    formParams.value = { ...defaultFormParams };
    showModal.value = true;
  }

  // 编辑租户
  function handleEdit(record: TenantInfo) {
    formParams.value = {
      id: record.id,
      name: record.name,
      code: record.code,
      domain: record.domain || '',
      maxUsers: record.maxUsers,
      storageLimit: record.storageLimit,
      expireAt: record.expireAt || '',
      adminName: '',
      adminEmail: '',
      adminPassword: '',
      remark: record.remark || '',
    };
    showModal.value = true;
  }

  // 查看详情
  async function handleDetail(record: TenantInfo) {
    detailLoading.value = true;
    showDetailModal.value = true;
    try {
      const response = await getTenantDetail({ id: record.id });
      if (response.code === 0) {
        tenantDetail.value = response.data;
      }
    } catch (error) {
      console.error('获取租户详情失败:', error);
      message.error('获取租户详情失败');
    } finally {
      detailLoading.value = false;
    }
  }

  // 切换状态
  function handleToggleStatus(record: TenantInfo) {
    const newStatus = record.status === 1 ? 2 : 1;
    const statusText = newStatus === 1 ? '启用' : '禁用';
    
    dialog.warning({
      title: '确认操作',
      content: `确定要${statusText}租户 "${record.name}" 吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await updateTenantStatus({ id: record.id, status: newStatus });
          message.success(`${statusText}成功`);
          refreshTable();
        } catch (error) {
          console.error(`${statusText}租户失败:`, error);
          message.error(`${statusText}失败`);
        }
      },
    });
  }

  // 删除租户
  function handleDelete(record: TenantInfo) {
    dialog.warning({
      title: '确认删除',
      content: `确定要删除租户 "${record.name}" 吗？此操作不可恢复！`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await deleteTenant({ id: record.id });
          message.success('删除成功');
          refreshTable();
        } catch (error) {
          console.error('删除租户失败:', error);
          message.error('删除失败');
        }
      },
    });
  }

  // 确认表单
  function confirmForm() {
    formRef.value.validate(async (errors: any) => {
      if (!errors) {
        formBtnLoading.value = true;
        try {
          if (formParams.value.id && formParams.value.id > 0) {
            // 编辑租户
            const updateParams: UpdateTenantParams = {
              id: formParams.value.id,
              name: formParams.value.name,
              domain: formParams.value.domain,
              maxUsers: formParams.value.maxUsers,
              storageLimit: formParams.value.storageLimit,
              expireAt: formParams.value.expireAt,
              remark: formParams.value.remark,
            };
            await updateTenant(updateParams);
          } else {
            // 创建租户
            const createParams: CreateTenantParams = {
              ...formParams.value,
              adminPassword: btoa(formParams.value.adminPassword), // Base64编码
            };
            await createTenant(createParams);
          }
          message.success('操作成功');
          showModal.value = false;
          refreshTable();
        } catch (error) {
          console.error('操作失败:', error);
          message.error('操作失败');
        } finally {
          formBtnLoading.value = false;
        }
      } else {
        message.error('请填写完整信息');
      }
    });
  }

  // 关闭表单
  function closeForm() {
    showModal.value = false;
    formParams.value = { ...defaultFormParams };
  }

  onMounted(() => {
    // 组件挂载时可以执行一些初始化操作
  });
</script>

<style lang="less" scoped>
  .proCard {
    margin-bottom: 16px;
  }
</style>
