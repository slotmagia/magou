<template>
  <n-form
    ref="formRef"
    label-placement="left"
    size="large"
    :model="mode === 'account' ? formInline : formMobile"
    :rules="mode === 'account' ? rules : mobileRules"
  >
    <template v-if="mode === 'account'">
      <n-form-item path="tenantCode">
        <n-input
          @keyup.enter="debounceHandleSubmit"
          v-model:value="formInline.tenantCode"
          placeholder="请输入租户编码"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <SafetyCertificateOutlined />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="username">
        <n-input
          @keyup.enter="debounceHandleSubmit"
          v-model:value="formInline.username"
          placeholder="请输入用户名"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <PersonOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="pass">
        <n-input
          @keyup.enter="debounceHandleSubmit"
          v-model:value="formInline.pass"
          type="password"
          show-password-on="click"
          placeholder="请输入密码"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <LockClosedOutline />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="code" v-show="showCaptcha">
        <n-input-group>
          <n-input
            :style="{ width: '70%' }"
            placeholder="请输入验证码"
            @keyup.enter="debounceHandleSubmit"
            v-model:value="formInline.code"
          >
            <template #prefix>
              <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
            </template>
          </n-input>
          
          <div style="width: 30%; position: relative;">
            <n-spin :show="captchaLoading" size="small">
              <div 
                style="width: 100%; height: 40px; cursor: pointer; border: 1px solid #d9d9d9; border-left: none; display: flex; align-items: center; justify-content: center; background: #fafafa;"
                @click="refreshCode"
                title="点击刷新验证码"
              >
                <img
                  v-if="codeBase64"
                  style="width: 100%; height: 100%; object-fit: contain;"
                  :src="codeBase64"
                  alt="验证码"
                />
                <span v-else style="color: #999; font-size: 12px;">点击获取</span>
              </div>
            </n-spin>
          </div>
        </n-input-group>
      </n-form-item>
    </template>

    <template v-if="mode === 'mobile'">
      <n-form-item path="mobile">
        <n-input
          @keyup.enter="handleMobileSubmit"
          v-model:value="formMobile.mobile"
          placeholder="请输入手机号码"
        >
          <template #prefix>
            <n-icon size="18" color="#808695">
              <MobileOutlined />
            </n-icon>
          </template>
        </n-input>
      </n-form-item>

      <n-form-item path="code">
        <n-input-group>
          <n-input
            @keyup.enter="handleMobileSubmit"
            v-model:value="formMobile.code"
            placeholder="请输入验证码"
          >
            <template #prefix>
              <n-icon size="18" color="#808695" :component="SafetyCertificateOutlined" />
            </template>
          </n-input>
          <n-button
            type="primary"
            ghost
            @click="sendMobileCode"
            :disabled="isCounting"
            :loading="sendLoading"
          >
            {{ sendLabel }}
          </n-button>
        </n-input-group>
      </n-form-item>
    </template>

    <n-space :vertical="true" :size="24">
      <div class="flex-y-center justify-between">
        <n-checkbox v-model:checked="autoLogin">自动登录</n-checkbox>
        <n-button :text="true" @click="handleResetPassword">忘记密码？</n-button>
      </div>
      <n-button type="primary" size="large" :block="true" :loading="loading" @click="handleLogin">
        登录
      </n-button>

      <FormOther moduleKey="register" tag="注册账号" @updateActiveModule="updateActiveModule" />
    </n-space>

    <DemoAccount @login="handleDemoAccountLogin" />
  </n-form>
</template>

<script lang="ts" setup>
  import '../components/style.less';
  import { ref, onMounted, computed } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useUserStore } from '@/store/modules/user';
  import { useMessage } from 'naive-ui';
  import { ResultEnum } from '@/enums/httpEnum';
  import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5';
  import { PageEnum } from '@/enums/pageEnum';
  import { SafetyCertificateOutlined, MobileOutlined } from '@vicons/antd';
  import { getCaptcha } from '@/api/auth';
  import { aesEcb } from '@/utils/encrypt';
  import DemoAccount from './demo-account.vue';
  import FormOther from '../components/form-other.vue';
  import { useSendCode } from '@/hooks/common';
  import { SendSms } from '@/api/system/user';
  import { validate } from '@/utils/validateUtil';
  import { useDebounceFn } from '@vueuse/core';

  interface Props {
    mode: string;
  }

  const props = withDefaults(defineProps<Props>(), {
    mode: 'account',
  });

  interface FormState {
    tenantCode: string;
    username: string;
    pass: string;
    cid: string;
    code: string;
    password: string;
  }

  interface FormMobileState {
    mobile: string;
    code: string;
  }

  const formRef = ref();
  const message = useMessage();
  const loading = ref(false);
  const autoLogin = ref(true);
  const codeBase64 = ref('');
  const captchaLoading = ref(false);
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();
  const { sendLabel, isCounting, loading: sendLoading, activateSend } = useSendCode();
  const emit = defineEmits(['updateActiveModule']);
  const LOGIN_NAME = PageEnum.BASE_LOGIN_NAME;
  const debounceHandleSubmit = useDebounceFn((e) => {
    handleSubmit(e);
  }, 500);
  const formInline = ref<FormState>({
    tenantCode: '',
    username: '',
    pass: '',
    cid: '',
    code: '',
    password: '',
  });

  const formMobile = ref<FormMobileState>({
    mobile: '',
    code: '',
  });

  const rules = computed(() => {
    const baseRules: any = {
      tenantCode: { required: true, message: '请输入租户编码', trigger: 'blur' },
      username: { required: true, message: '请输入用户名', trigger: 'blur' },
      pass: { required: true, message: '请输入密码', trigger: 'blur' },
    };

    // 账号密码登录强制要求验证码
    if (props.mode === 'account') {
      baseRules.code = { required: true, message: '请输入验证码', trigger: 'blur' };
    } else {
      // 其他登录方式根据配置决定
      const captchaEnabled = userStore.loginConfig?.loginCaptchaSwitch === 1;
      if (captchaEnabled) {
        baseRules.code = { required: true, message: '请输入验证码', trigger: 'blur' };
      }
    }

    return baseRules;
  });

  const mobileRules = {
    mobile: { required: true, message: '请输入手机号码', trigger: 'blur' },
    code: { required: true, message: '请输入验证码', trigger: 'blur' },
  };

  // 计算是否显示验证码
  const showCaptcha = computed(() => {
    // 账号密码登录强制显示验证码
    if (props.mode === 'account') {
      return true; // 强制显示验证码
    }
    // 其他登录方式根据配置决定
    const captchaEnabled = userStore.loginConfig?.loginCaptchaSwitch === 1;
    return captchaEnabled && codeBase64.value !== '';
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    formRef.value.validate(async (errors) => {
      if (!errors) {
        // 账号密码登录强制验证验证码
        if (props.mode === 'account') {
          if (!formInline.value.code) {
            message.error('请输入验证码');
            return;
          }
          if (!formInline.value.cid) {
            message.error('验证码ID缺失，请刷新验证码');
            await refreshCode();
            return;
          }
        }

        const params = {
          tenantCode: formInline.value.tenantCode,
          username: formInline.value.username,
          password: aesEcb.encrypt(formInline.value.pass),
          captchaId: formInline.value.cid,
          captcha: formInline.value.code,
        };
        
        console.log('登录参数:', {
          ...params,
          password: '[加密后的密码]' // 不打印实际密码
        });
        
        await handleLoginResp(userStore.tenantLogin(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };

  async function refreshCode() {
    // 账号密码登录强制加载验证码
    if (props.mode === 'account') {
      console.log('账号密码登录，强制加载验证码');
    } else {
      // 其他登录方式检查配置
      const captchaEnabled = userStore.loginConfig?.loginCaptchaSwitch === 1;
      if (!captchaEnabled) {
        console.log('验证码功能未启用');
        return;
      }
    }
    
    captchaLoading.value = true;
    
    try {
      console.log('正在调用验证码API: /captcha');
      const response: any = await getCaptcha();
      console.log('验证码API响应:', response);
      console.log('响应类型:', typeof response);
      console.log('响应详情:', JSON.stringify(response, null, 2));
      
      if (response) {
        console.log('解析响应数据...');
        
        // axios配置会自动提取data字段，所以直接从response中获取
        const { captchaImage, captchaId } = response;
        console.log('captchaId:', captchaId, '类型:', typeof captchaId);
        console.log('captchaImage存在:', !!captchaImage, '长度:', captchaImage?.length);
        
        if (captchaImage && captchaId) {
          codeBase64.value = captchaImage;
          formInline.value.cid = captchaId;
          formInline.value.code = '';
          console.log('✅ 验证码加载成功，ID:', captchaId);
          return;
        } else {
          console.error('验证码数据缺失 - captchaImage:', !!captchaImage, 'captchaId:', !!captchaId);
          throw new Error(`验证码数据缺失: captchaImage=${!!captchaImage}, captchaId=${!!captchaId}`);
        }
      } else {
        console.error('响应为空');
        throw new Error('API响应为空');
      }
    } catch (error) {
      console.error('❌ 验证码API调用失败:', error instanceof Error ? error.message : String(error));
      console.error('错误详情:', error);
      
      message.error('验证码加载失败，请检查网络连接后重试');
    } finally {
      captchaLoading.value = false;
    }
  }

  async function handleDemoAccountLogin(user: { username: string; password: string }) {
    const params = {
      username: user.username,
      password: aesEcb.encrypt(user.password),
      isLock: true,
    };
    await handleLoginResp(userStore.login(params));
  }

  const handleMobileSubmit = (e) => {
    e.preventDefault();
    formRef.value.validate(async (errors) => {
      if (!errors) {
        const params = {
          mobile: formMobile.value.mobile,
          code: formMobile.value.code,
        };
        await handleLoginResp(userStore.mobileLogin(params));
      } else {
        message.error('请填写完整信息，并且进行验证码校验');
      }
    });
  };

  function updateActiveModule(key: string) {
    emit('updateActiveModule', key);
  }

  function sendMobileCode() {
    validate.phone(mobileRules.mobile, formMobile.value.mobile, function (error?: Error) {
      if (error === undefined) {
        activateSend(SendSms({ mobile: formMobile.value.mobile, event: 'login' }));
        return;
      }
      message.error(error.message);
    });
  }

  function handleResetPassword() {
    message.info('如果你忘记了密码，请联系管理员找回');
  }

  function handleLogin(e) {
    if (props.mode === 'account') {
      debounceHandleSubmit(e);
      return;
    }

    handleMobileSubmit(e);
  }

  async function handleLoginResp(request: Promise<any>) {
    message.loading('登录中...');
    loading.value = true;
    try {
      const { code, message: msg } = await request;
      message.destroyAll();
      if (code == ResultEnum.SUCCESS) {
        const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
        message.success('登录成功，即将进入系统');
        if (route.name === LOGIN_NAME) {
          await router.replace('/');
        } else {
          await router.replace(toPath);
        }
      } else {
        message.destroyAll();
        message.info(msg || '登录失败');
        await refreshCode();
      }
    } finally {
      loading.value = false;
    }
  }

  // 初始化验证码
  async function initCaptcha() {
    console.log('🔄 开始初始化验证码，当前模式:', props.mode);
    
    // 账号密码登录模式强制加载验证码
    if (props.mode === 'account') {
      console.log('📝 账号密码登录模式，强制加载验证码');
      await refreshCode();
      return;
    }
    
    // 其他模式根据配置决定
    console.log('📱 其他登录模式，检查配置');
    try {
      if (!userStore.loginConfig) {
        console.log('⚙️ 登录配置未加载，正在加载...');
        await userStore.LoadLoginConfig();
      }
      console.log('✅ 登录配置已加载，开始加载验证码');
      await refreshCode();
    } catch (error) {
      console.warn('❌ 初始化验证码失败:', error);
    }
  }

  onMounted(async () => {
    console.log('🚀 登录组件已挂载，准备初始化验证码');
    // 延迟一点时间确保组件完全挂载
    setTimeout(async () => {
      console.log('⏰ 延迟100ms后开始初始化验证码');
      await initCaptcha();
    }, 100);
  });
</script>
