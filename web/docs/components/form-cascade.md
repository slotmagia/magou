# 表单多级联动实现指南

在 HotGo 2.0 中，表单多级联动是一个常见的业务需求。本指南将详细介绍如何实现各种级联场景，从简单的两级联动到复杂的多级联动。

## 🎯 联动实现原理

### 核心机制
表单联动主要通过以下几种方式实现：

1. **监听表单字段变化** - 使用 `watch` 监听某个字段的值变化
2. **动态更新字段配置** - 通过 `updateSchema` 更新依赖字段的选项
3. **条件显示控制** - 使用 `ifShow` 控制字段的显示/隐藏
4. **异步数据加载** - 根据上级字段值异步加载下级数据

### 技术实现
```typescript
// 基础联动模式
watch(() => formModel.parentField, async (newValue) => {
  if (newValue) {
    // 清空子级字段
    setFieldsValue({ childField: null });
    
    // 加载子级数据
    const childOptions = await loadChildData(newValue);
    
    // 更新子级字段配置
    updateSchema([{
      field: 'childField',
      componentProps: {
        options: childOptions,
      },
    }]);
  }
});
```

## 🏗️ 基础二级联动

### 省市联动示例

```vue
<template>
  <div class="p-4">
    <h3>省市二级联动</h3>
    <BasicForm @register="register" />
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

// 模拟省份数据
const provinces = [
  { label: '广东省', value: 'guangdong' },
  { label: '江苏省', value: 'jiangsu' },
  { label: '浙江省', value: 'zhejiang' },
];

// 模拟城市数据
const cityMap = {
  guangdong: [
    { label: '广州市', value: 'guangzhou' },
    { label: '深圳市', value: 'shenzhen' },
    { label: '佛山市', value: 'foshan' },
  ],
  jiangsu: [
    { label: '南京市', value: 'nanjing' },
    { label: '苏州市', value: 'suzhou' },
    { label: '无锡市', value: 'wuxi' },
  ],
  zhejiang: [
    { label: '杭州市', value: 'hangzhou' },
    { label: '宁波市', value: 'ningbo' },
    { label: '温州市', value: 'wenzhou' },
  ],
};

const schemas: FormSchema[] = [
  {
    field: 'province',
    label: '省份',
    component: 'NSelect',
    componentProps: {
      placeholder: '请选择省份',
      options: provinces,
      clearable: true,
    },
    rules: [{ required: true, message: '请选择省份' }],
  },
  {
    field: 'city',
    label: '城市',
    component: 'NSelect',
    componentProps: {
      placeholder: '请先选择省份',
      options: [],
      clearable: true,
      disabled: true,
    },
    rules: [{ required: true, message: '请选择城市' }],
  },
];

const [register, { updateSchema, setFieldsValue, getFieldsValue }] = useForm({
  schemas,
  labelWidth: 80,
  showActionButtonGroup: false,
});

// 监听省份变化
watch(
  () => getFieldsValue().province,
  (provinceValue) => {
    // 清空城市字段
    setFieldsValue({ city: null });
    
    if (provinceValue) {
      // 获取对应的城市列表
      const cityOptions = cityMap[provinceValue] || [];
      
      // 更新城市字段配置
      updateSchema([
        {
          field: 'city',
          componentProps: {
            placeholder: '请选择城市',
            options: cityOptions,
            disabled: false,
          },
        },
      ]);
    } else {
      // 没有选择省份时禁用城市选择
      updateSchema([
        {
          field: 'city',
          componentProps: {
            placeholder: '请先选择省份',
            options: [],
            disabled: true,
          },
        },
      ]);
    }
  },
  { immediate: true }
);
</script>
```

## 🏢 三级联动示例

### 省市区三级联动

```vue
<template>
  <div class="p-4">
    <h3>省市区三级联动</h3>
    <BasicForm @register="register" />
    
    <div class="mt-4">
      <n-button @click="showCurrentValues">查看当前值</n-button>
      <div v-if="currentValues" class="mt-2">
        <p>当前选择: {{ currentValues }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const currentValues = ref(null);

// 模拟异步加载数据的函数
const loadProvinces = async () => {
  // 模拟网络请求
  return new Promise(resolve => {
    setTimeout(() => {
      resolve([
        { label: '广东省', value: 'guangdong' },
        { label: '江苏省', value: 'jiangsu' },
        { label: '浙江省', value: 'zhejiang' },
      ]);
    }, 500);
  });
};

const loadCities = async (provinceCode: string) => {
  return new Promise(resolve => {
    setTimeout(() => {
      const cityMap = {
        guangdong: [
          { label: '广州市', value: 'guangzhou', code: 'gz' },
          { label: '深圳市', value: 'shenzhen', code: 'sz' },
          { label: '佛山市', value: 'foshan', code: 'fs' },
        ],
        jiangsu: [
          { label: '南京市', value: 'nanjing', code: 'nj' },
          { label: '苏州市', value: 'suzhou', code: 'su' },
          { label: '无锡市', value: 'wuxi', code: 'wx' },
        ],
        zhejiang: [
          { label: '杭州市', value: 'hangzhou', code: 'hz' },
          { label: '宁波市', value: 'ningbo', code: 'nb' },
          { label: '温州市', value: 'wenzhou', code: 'wz' },
        ],
      };
      resolve(cityMap[provinceCode] || []);
    }, 300);
  });
};

const loadDistricts = async (cityCode: string) => {
  return new Promise(resolve => {
    setTimeout(() => {
      const districtMap = {
        gz: [
          { label: '天河区', value: 'tianhe' },
          { label: '海珠区', value: 'haizhu' },
          { label: '越秀区', value: 'yuexiu' },
        ],
        sz: [
          { label: '南山区', value: 'nanshan' },
          { label: '福田区', value: 'futian' },
          { label: '罗湖区', value: 'luohu' },
        ],
        nj: [
          { label: '玄武区', value: 'xuanwu' },
          { label: '秦淮区', value: 'qinhuai' },
          { label: '建邺区', value: 'jianye' },
        ],
        hz: [
          { label: '西湖区', value: 'xihu' },
          { label: '拱墅区', value: 'gongshu' },
          { label: '江干区', value: 'jianggan' },
        ],
      };
      resolve(districtMap[cityCode] || []);
    }, 200);
  });
};

const schemas: FormSchema[] = [
  {
    field: 'province',
    label: '省份',
    component: 'NSelect',
    componentProps: {
      placeholder: '请选择省份',
      options: [],
      clearable: true,
      loading: false,
    },
    rules: [{ required: true, message: '请选择省份' }],
  },
  {
    field: 'city',
    label: '城市',
    component: 'NSelect',
    componentProps: {
      placeholder: '请先选择省份',
      options: [],
      clearable: true,
      disabled: true,
      loading: false,
    },
    rules: [{ required: true, message: '请选择城市' }],
  },
  {
    field: 'district',
    label: '区县',
    component: 'NSelect',
    componentProps: {
      placeholder: '请先选择城市',
      options: [],
      clearable: true,
      disabled: true,
      loading: false,
    },
    rules: [{ required: true, message: '请选择区县' }],
  },
];

const [register, { updateSchema, setFieldsValue, getFieldsValue }] = useForm({
  schemas,
  labelWidth: 80,
  showActionButtonGroup: false,
});

// 初始化省份数据
const initProvinces = async () => {
  updateSchema([
    {
      field: 'province',
      componentProps: { loading: true },
    },
  ]);
  
  const provinces = await loadProvinces();
  
  updateSchema([
    {
      field: 'province',
      componentProps: {
        options: provinces,
        loading: false,
      },
    },
  ]);
};

// 监听省份变化
watch(
  () => getFieldsValue().province,
  async (provinceValue) => {
    // 清空下级字段
    setFieldsValue({ city: null, district: null });
    
    if (provinceValue) {
      // 启用城市选择，显示加载状态
      updateSchema([
        {
          field: 'city',
          componentProps: {
            disabled: false,
            loading: true,
            placeholder: '加载中...',
            options: [],
          },
        },
        {
          field: 'district',
          componentProps: {
            disabled: true,
            placeholder: '请先选择城市',
            options: [],
          },
        },
      ]);
      
      // 加载城市数据
      const cities = await loadCities(provinceValue);
      
      updateSchema([
        {
          field: 'city',
          componentProps: {
            loading: false,
            placeholder: '请选择城市',
            options: cities,
          },
        },
      ]);
    } else {
      // 禁用下级选择
      updateSchema([
        {
          field: 'city',
          componentProps: {
            disabled: true,
            placeholder: '请先选择省份',
            options: [],
          },
        },
        {
          field: 'district',
          componentProps: {
            disabled: true,
            placeholder: '请先选择城市',
            options: [],
          },
        },
      ]);
    }
  }
);

// 监听城市变化
watch(
  () => getFieldsValue().city,
  async (cityValue) => {
    // 清空区县字段
    setFieldsValue({ district: null });
    
    if (cityValue) {
      // 获取选中城市的code
      const formValues = getFieldsValue();
      const provinces = await loadProvinces();
      const cities = await loadCities(formValues.province);
      const selectedCity = cities.find(city => city.value === cityValue);
      
      if (selectedCity) {
        // 启用区县选择，显示加载状态
        updateSchema([
          {
            field: 'district',
            componentProps: {
              disabled: false,
              loading: true,
              placeholder: '加载中...',
              options: [],
            },
          },
        ]);
        
        // 加载区县数据
        const districts = await loadDistricts(selectedCity.code);
        
        updateSchema([
          {
            field: 'district',
            componentProps: {
              loading: false,
              placeholder: '请选择区县',
              options: districts,
            },
          },
        ]);
      }
    } else {
      // 禁用区县选择
      updateSchema([
        {
          field: 'district',
          componentProps: {
            disabled: true,
            placeholder: '请先选择城市',
            options: [],
          },
        },
      ]);
    }
  }
);

// 显示当前值
const showCurrentValues = () => {
  currentValues.value = getFieldsValue();
};

// 组件挂载时初始化省份数据
onMounted(() => {
  initProvinces();
});
</script>
```

## 🏭 复杂业务联动

### 产品类型-品牌-型号联动

```vue
<template>
  <div class="p-4">
    <h3>产品类型-品牌-型号联动</h3>
    <BasicForm @register="register" />
    
    <div class="mt-4">
      <n-card title="选择结果">
        <pre>{{ JSON.stringify(productInfo, null, 2) }}</pre>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const productInfo = ref({});

// 模拟产品数据
const productData = {
  categories: [
    { label: '手机', value: 'mobile', icon: '📱' },
    { label: '电脑', value: 'computer', icon: '💻' },
    { label: '汽车', value: 'car', icon: '🚗' },
  ],
  brands: {
    mobile: [
      { label: 'Apple', value: 'apple' },
      { label: '华为', value: 'huawei' },
      { label: '小米', value: 'xiaomi' },
    ],
    computer: [
      { label: 'Apple', value: 'apple' },
      { label: '联想', value: 'lenovo' },
      { label: '华硕', value: 'asus' },
    ],
    car: [
      { label: '奔驰', value: 'benz' },
      { label: '宝马', value: 'bmw' },
      { label: '奥迪', value: 'audi' },
    ],
  },
  models: {
    'mobile-apple': [
      { label: 'iPhone 15 Pro', value: 'iphone15pro', price: 9999 },
      { label: 'iPhone 15', value: 'iphone15', price: 5999 },
      { label: 'iPhone 14', value: 'iphone14', price: 4999 },
    ],
    'mobile-huawei': [
      { label: 'Mate 60 Pro', value: 'mate60pro', price: 6999 },
      { label: 'P60 Pro', value: 'p60pro', price: 5988 },
    ],
    'mobile-xiaomi': [
      { label: '小米14 Pro', value: 'mi14pro', price: 4999 },
      { label: '小米14', value: 'mi14', price: 3999 },
    ],
    'computer-apple': [
      { label: 'MacBook Pro M3', value: 'mbpm3', price: 14999 },
      { label: 'MacBook Air M2', value: 'mbam2', price: 8999 },
    ],
    'computer-lenovo': [
      { label: 'ThinkPad X1', value: 'tpx1', price: 12999 },
      { label: 'ThinkPad T14', value: 'tpt14', price: 8999 },
    ],
    'car-benz': [
      { label: '奔驰C级', value: 'benzc', price: 320000 },
      { label: '奔驰E级', value: 'benze', price: 450000 },
    ],
  },
};

const schemas: FormSchema[] = [
  {
    field: 'category',
    label: '产品类型',
    component: 'NSelect',
    componentProps: {
      placeholder: '请选择产品类型',
      options: productData.categories,
      clearable: true,
    },
    rules: [{ required: true, message: '请选择产品类型' }],
  },
  {
    field: 'brand',
    label: '品牌',
    component: 'NSelect',
    componentProps: {
      placeholder: '请先选择产品类型',
      options: [],
      clearable: true,
      disabled: true,
    },
    rules: [{ required: true, message: '请选择品牌' }],
  },
  {
    field: 'model',
    label: '型号',
    component: 'NSelect',
    componentProps: {
      placeholder: '请先选择品牌',
      options: [],
      clearable: true,
      disabled: true,
    },
    rules: [{ required: true, message: '请选择型号' }],
  },
  {
    field: 'price',
    label: '价格',
    component: 'NInputNumber',
    componentProps: {
      placeholder: '自动填充',
      disabled: true,
      prefix: '¥',
    },
  },
  {
    field: 'quantity',
    label: '数量',
    component: 'NInputNumber',
    componentProps: {
      placeholder: '请输入数量',
      min: 1,
      max: 999,
      defaultValue: 1,
    },
    rules: [{ required: true, message: '请输入数量' }],
  },
  {
    field: 'totalAmount',
    label: '总金额',
    component: 'NInputNumber',
    componentProps: {
      placeholder: '自动计算',
      disabled: true,
      prefix: '¥',
    },
  },
];

const [register, { updateSchema, setFieldsValue, getFieldsValue }] = useForm({
  schemas,
  labelWidth: 100,
  showActionButtonGroup: false,
});

// 监听产品类型变化
watch(
  () => getFieldsValue().category,
  (categoryValue) => {
    // 清空下级字段
    setFieldsValue({ 
      brand: null, 
      model: null, 
      price: null,
      totalAmount: null,
    });
    
    if (categoryValue) {
      const brands = productData.brands[categoryValue] || [];
      
      updateSchema([
        {
          field: 'brand',
          componentProps: {
            placeholder: '请选择品牌',
            options: brands,
            disabled: false,
          },
        },
        {
          field: 'model',
          componentProps: {
            placeholder: '请先选择品牌',
            options: [],
            disabled: true,
          },
        },
      ]);
    } else {
      updateSchema([
        {
          field: 'brand',
          componentProps: {
            placeholder: '请先选择产品类型',
            options: [],
            disabled: true,
          },
        },
        {
          field: 'model',
          componentProps: {
            placeholder: '请先选择品牌',
            options: [],
            disabled: true,
          },
        },
      ]);
    }
  }
);

// 监听品牌变化
watch(
  () => getFieldsValue().brand,
  (brandValue) => {
    const formValues = getFieldsValue();
    
    // 清空下级字段
    setFieldsValue({ 
      model: null, 
      price: null,
      totalAmount: null,
    });
    
    if (brandValue && formValues.category) {
      const modelKey = `${formValues.category}-${brandValue}`;
      const models = productData.models[modelKey] || [];
      
      updateSchema([
        {
          field: 'model',
          componentProps: {
            placeholder: '请选择型号',
            options: models,
            disabled: false,
          },
        },
      ]);
    } else {
      updateSchema([
        {
          field: 'model',
          componentProps: {
            placeholder: '请先选择品牌',
            options: [],
            disabled: true,
          },
        },
      ]);
    }
  }
);

// 监听型号变化
watch(
  () => getFieldsValue().model,
  (modelValue) => {
    const formValues = getFieldsValue();
    
    if (modelValue && formValues.category && formValues.brand) {
      const modelKey = `${formValues.category}-${formValues.brand}`;
      const models = productData.models[modelKey] || [];
      const selectedModel = models.find(m => m.value === modelValue);
      
      if (selectedModel) {
        setFieldsValue({ 
          price: selectedModel.price,
        });
        
        // 触发总金额计算
        calculateTotal();
      }
    } else {
      setFieldsValue({ 
        price: null,
        totalAmount: null,
      });
    }
  }
);

// 监听数量变化
watch(
  () => getFieldsValue().quantity,
  () => {
    calculateTotal();
  }
);

// 计算总金额
const calculateTotal = () => {
  const formValues = getFieldsValue();
  const { price, quantity } = formValues;
  
  if (price && quantity) {
    const total = price * quantity;
    setFieldsValue({ totalAmount: total });
  } else {
    setFieldsValue({ totalAmount: null });
  }
};

// 监听所有字段变化，更新产品信息展示
watch(
  () => getFieldsValue(),
  (values) => {
    productInfo.value = values;
  },
  { deep: true }
);
</script>
```

## 🔄 动态字段联动

### 条件字段显示/隐藏

```vue
<template>
  <div class="p-4">
    <h3>条件字段联动</h3>
    <BasicForm @register="register" />
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const schemas: FormSchema[] = [
  {
    field: 'userType',
    label: '用户类型',
    component: 'NRadioGroup',
    componentProps: {
      options: [
        { label: '个人用户', value: 'personal' },
        { label: '企业用户', value: 'enterprise' },
      ],
    },
    rules: [{ required: true, message: '请选择用户类型' }],
  },
  {
    field: 'name',
    label: '姓名',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入姓名',
    },
    rules: [{ required: true, message: '请输入姓名' }],
    // 只有个人用户才显示
    ifShow: ({ model }) => model.userType === 'personal',
  },
  {
    field: 'idCard',
    label: '身份证号',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入身份证号',
    },
    rules: [
      { required: true, message: '请输入身份证号' },
      { pattern: /^\d{17}[\dxX]$/, message: '身份证号格式不正确' },
    ],
    // 只有个人用户才显示
    ifShow: ({ model }) => model.userType === 'personal',
  },
  {
    field: 'companyName',
    label: '公司名称',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入公司名称',
    },
    rules: [{ required: true, message: '请输入公司名称' }],
    // 只有企业用户才显示
    ifShow: ({ model }) => model.userType === 'enterprise',
  },
  {
    field: 'businessLicense',
    label: '营业执照号',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入营业执照号',
    },
    rules: [{ required: true, message: '请输入营业执照号' }],
    // 只有企业用户才显示
    ifShow: ({ model }) => model.userType === 'enterprise',
  },
  {
    field: 'taxNumber',
    label: '税号',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入税号',
    },
    rules: [{ required: true, message: '请输入税号' }],
    // 只有企业用户才显示
    ifShow: ({ model }) => model.userType === 'enterprise',
  },
  {
    field: 'hasInvoice',
    label: '是否需要发票',
    component: 'NSwitch',
    componentProps: {
      checkedValue: true,
      uncheckedValue: false,
    },
  },
  {
    field: 'invoiceType',
    label: '发票类型',
    component: 'NSelect',
    componentProps: {
      placeholder: '请选择发票类型',
      options: [
        { label: '普通发票', value: 'normal' },
        { label: '专用发票', value: 'special' },
      ],
    },
    rules: [{ required: true, message: '请选择发票类型' }],
    // 只有需要发票时才显示
    ifShow: ({ model }) => model.hasInvoice === true,
  },
  {
    field: 'invoiceTitle',
    label: '发票抬头',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入发票抬头',
    },
    rules: [{ required: true, message: '请输入发票抬头' }],
    // 只有需要发票时才显示
    ifShow: ({ model }) => model.hasInvoice === true,
  },
];

const [register, { getFieldsValue, setFieldsValue }] = useForm({
  schemas,
  labelWidth: 120,
  showActionButtonGroup: true,
  submitButtonText: '提交注册',
});

// 监听用户类型变化，清空相关字段
watch(
  () => getFieldsValue().userType,
  (userType) => {
    if (userType === 'personal') {
      // 切换到个人用户，清空企业相关字段
      setFieldsValue({
        companyName: null,
        businessLicense: null,
        taxNumber: null,
      });
    } else if (userType === 'enterprise') {
      // 切换到企业用户，清空个人相关字段
      setFieldsValue({
        name: null,
        idCard: null,
      });
    }
  }
);

// 监听发票选择变化
watch(
  () => getFieldsValue().hasInvoice,
  (hasInvoice) => {
    if (!hasInvoice) {
      // 不需要发票时清空发票相关字段
      setFieldsValue({
        invoiceType: null,
        invoiceTitle: null,
      });
    }
  }
);
</script>
```

## 🔧 高级联动技巧

### 1. 使用 Composable 封装联动逻辑

```typescript
// composables/useCascadeForm.ts
import { watch, ref } from 'vue';

export function useCascadeForm(formMethods: any) {
  const { updateSchema, setFieldsValue, getFieldsValue } = formMethods;
  
  // 创建联动关系
  const createCascade = (config: {
    parentField: string;
    childField: string;
    loadOptions: (parentValue: any) => Promise<any[]>;
    clearOnChange?: boolean;
  }) => {
    const { parentField, childField, loadOptions, clearOnChange = true } = config;
    
    watch(
      () => getFieldsValue()[parentField],
      async (parentValue) => {
        if (clearOnChange) {
          setFieldsValue({ [childField]: null });
        }
        
        if (parentValue) {
          // 显示加载状态
          updateSchema([{
            field: childField,
            componentProps: {
              loading: true,
              disabled: false,
            },
          }]);
          
          try {
            const options = await loadOptions(parentValue);
            
            updateSchema([{
              field: childField,
              componentProps: {
                loading: false,
                options,
                disabled: false,
              },
            }]);
          } catch (error) {
            console.error('联动数据加载失败:', error);
            updateSchema([{
              field: childField,
              componentProps: {
                loading: false,
                disabled: true,
              },
            }]);
          }
        } else {
          updateSchema([{
            field: childField,
            componentProps: {
              loading: false,
              options: [],
              disabled: true,
            },
          }]);
        }
      }
    );
  };
  
  return {
    createCascade,
  };
}
```

### 2. 批量联动配置

```typescript
// utils/cascadeConfig.ts
export interface CascadeConfig {
  field: string;
  dependsOn: string[];
  loadOptions: (...values: any[]) => Promise<any[]>;
  clearOnChange?: boolean;
  condition?: (...values: any[]) => boolean;
}

export function setupCascades(formMethods: any, configs: CascadeConfig[]) {
  const { updateSchema, setFieldsValue, getFieldsValue } = formMethods;
  
  configs.forEach(config => {
    const { field, dependsOn, loadOptions, clearOnChange = true, condition } = config;
    
    // 创建依赖字段的监听器
    const dependencies = dependsOn.map(dep => () => getFieldsValue()[dep]);
    
    watch(
      dependencies,
      async (values) => {
        if (clearOnChange) {
          setFieldsValue({ [field]: null });
        }
        
        // 检查条件
        if (condition && !condition(...values)) {
          updateSchema([{
            field,
            componentProps: {
              disabled: true,
              options: [],
            },
          }]);
          return;
        }
        
        // 检查是否所有依赖都有值
        const hasAllDependencies = values.every(v => v != null && v !== '');
        
        if (hasAllDependencies) {
          updateSchema([{
            field,
            componentProps: {
              loading: true,
              disabled: false,
            },
          }]);
          
          try {
            const options = await loadOptions(...values);
            
            updateSchema([{
              field,
              componentProps: {
                loading: false,
                options,
                disabled: false,
              },
            }]);
          } catch (error) {
            updateSchema([{
              field,
              componentProps: {
                loading: false,
                disabled: true,
              },
            }]);
          }
        } else {
          updateSchema([{
            field,
            componentProps: {
              loading: false,
              options: [],
              disabled: true,
            },
          }]);
        }
      },
      { immediate: true }
    );
  });
}
```

### 3. 使用批量联动配置

```vue
<script setup lang="ts">
import { setupCascades, type CascadeConfig } from '@/utils/cascadeConfig';

// 定义联动配置
const cascadeConfigs: CascadeConfig[] = [
  {
    field: 'city',
    dependsOn: ['province'],
    loadOptions: async (province) => await loadCities(province),
  },
  {
    field: 'district',
    dependsOn: ['province', 'city'],
    loadOptions: async (province, city) => await loadDistricts(province, city),
  },
  {
    field: 'model',
    dependsOn: ['category', 'brand'],
    loadOptions: async (category, brand) => await loadModels(category, brand),
    condition: (category, brand) => category && brand, // 只有两个都选择了才加载
  },
];

const [register, formMethods] = useForm({
  schemas,
});

// 设置联动关系
onMounted(() => {
  setupCascades(formMethods, cascadeConfigs);
});
</script>
```

## 📝 最佳实践

### 1. 性能优化
- **防抖处理**: 对频繁变化的字段使用防抖
- **缓存机制**: 缓存已加载的数据避免重复请求
- **懒加载**: 只在需要时加载数据

### 2. 用户体验
- **加载状态**: 显示数据加载中的状态
- **错误处理**: 优雅处理加载失败的情况
- **清空提示**: 在联动变化时给用户适当提示

### 3. 代码组织
- **逻辑分离**: 将联动逻辑提取到单独的 composable
- **配置驱动**: 使用配置化的方式管理复杂联动
- **类型安全**: 使用 TypeScript 确保类型安全

### 4. 测试策略
- **单元测试**: 测试联动逻辑的正确性
- **集成测试**: 测试表单整体的联动效果
- **用户测试**: 验证用户操作流程的合理性

通过以上方法，您可以在 HotGo 2.0 中实现各种复杂的表单多级联动功能，提供流畅的用户体验。 🚀






