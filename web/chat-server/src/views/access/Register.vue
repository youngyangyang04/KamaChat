<template>
  <div class="register-wrap">
    <div
      class="register-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <h2 class="register-item">注册</h2>
      <el-form
        ref="formRef"
        :model="registerData"
        label-width="70px"
        class="demo-dynamic"
      >
        <el-form-item
          prop="account"
          label="账号"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
            {
              min: 3,
              max: 20,
              message: '账号长度在 3 到 20 个字符',
              trigger: 'blur',
            },
          ]"
        >
          <el-input v-model="registerData.account" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item
          prop="nickname"
          label="昵称"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
            {
              min: 1,
              max: 20,
              message: '昵称长度在 1 到 20 个字符',
              trigger: 'blur',
            },
          ]"
        >
          <el-input v-model="registerData.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item
          prop="password"
          label="密码"
          :rules="[
            {
              required: true,
              message: '此项为必填项',
              trigger: 'blur',
            },
            {
              min: 6,
              max: 50,
              message: '密码长度在 6 到 50 个字符',
              trigger: 'blur',
            },
          ]"
        >
          <el-input type="password" v-model="registerData.password" placeholder="请输入密码" />
        </el-form-item>
      </el-form>
      <div class="register-button-container">
        <el-button type="primary" class="register-btn" @click="handleRegister"
          >注册</el-button
        >
      </div>
      <div class="go-login-button-container">
        <button class="go-password-login-btn" @click="handleLogin">
          返回登录
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import axios from "@/utils/axios";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useStore } from "vuex";
import { initializeUserKeys, loginAndDeriveMasterKey } from "@/crypto";
import { setCurrentUserId } from "@/crypto/cryptoStore";

export default {
  name: "Register",
  setup() {
    const data = reactive({
      registerData: {
        account: "",
        nickname: "",
        password: "",
      },
    });
    const router = useRouter();
    const store = useStore();
    const handleRegister = async () => {
      try {
        if (
          !data.registerData.account ||
          !data.registerData.nickname ||
          !data.registerData.password
        ) {
          ElMessage.error("请填写完整注册信息。");
          return;
        }
        if (
          data.registerData.account.length < 3 ||
          data.registerData.account.length > 20
        ) {
          ElMessage.error("账号长度在 3 到 20 个字符。");
          return;
        }
        if (
          data.registerData.nickname.length < 1 ||
          data.registerData.nickname.length > 20
        ) {
          ElMessage.error("昵称长度在 1 到 20 个字符。");
          return;
        }
        if (
          data.registerData.password.length < 6 ||
          data.registerData.password.length > 50
        ) {
          ElMessage.error("密码长度在 6 到 50 个字符。");
          return;
        }

        // 第一步：调用普通注册接口获取用户 ID
        const registerResponse = await axios.post("/register", data.registerData);
        if (registerResponse.data.code !== 200) {
          ElMessage.error(registerResponse.data.message);
          console.log(registerResponse.data.message);
          return;
        }

        // 获取用户信息
        if (!registerResponse.data.data.avatar.startsWith("http")) {
          registerResponse.data.data.avatar =
            store.state.backendUrl + registerResponse.data.data.avatar;
        }
        store.commit("setUserInfo", registerResponse.data.data);
        const userId = registerResponse.data.data.uuid;

        // 第二步：立即设置当前用户 ID（确保 IndexedDB 使用正确的数据库）
        setCurrentUserId(userId);
        console.log(`🔐 [Register.vue] 已设置当前用户 ID: ${userId}`);

        // 第三步：生成加密密钥（此时已设置用户 ID，密钥会保存到正确的用户专属数据库）
        ElMessage.info("正在生成加密密钥，请稍候...");
        let cryptoKeys;
        try {
          cryptoKeys = await initializeUserKeys(data.registerData.password);
          console.log("✅ 加密密钥生成成功");
        } catch (error) {
          console.error("❌ 生成加密密钥失败:", error);
          ElMessage.error("生成加密密钥失败: " + error.message);
          return;
        }

        // 第四步：上传公钥束到服务器
        try {
          const uploadResponse = await axios.post("/crypto/uploadPublicKeyBundle", {
            user_id: userId,
            ...cryptoKeys,
          });
          if (uploadResponse.data.code === 200) {
            console.log("✅ 公钥束上传成功");
            ElMessage.success("注册成功！(端到端加密已启用)");
          } else {
            console.warn("⚠️ 上传公钥束失败:", uploadResponse.data.message);
            ElMessage.warning("公钥上传失败，但账号已创建成功");
          }
        } catch (error) {
          console.error("❌ 上传公钥束出错:", error);
          ElMessage.warning("公钥上传失败，但账号已创建成功");
        }

        // 第五步：派生主密钥并保存到内存
        const masterKey = await loginAndDeriveMasterKey(data.registerData.password);
        if (masterKey) {
          store.commit("setMasterKey", masterKey);
          console.log("✅ 主密钥已保存到内存");
        }

        // 第六步：创建 WebSocket 连接
        const wsUrl =
          store.state.wsUrl + "/wss?client_id=" + userId + "&token=" + encodeURIComponent(registerResponse.data.data.token);
        console.log(wsUrl);
        store.state.socket = new WebSocket(wsUrl);
        store.state.socket.onopen = () => {
          console.log("🌐 [Register.vue] WebSocket连接已打开");
        };
        // 不设置 onmessage，让 App.vue 统一管理
        store.state.socket.onclose = () => {
          console.log("🌐 [Register.vue] WebSocket连接已关闭");
        };
        store.state.socket.onerror = (error) => {
          console.log("🌐 [Register.vue] WebSocket连接发生错误", error);
        };
        router.push("/chat/sessionlist");
      } catch (error) {
        ElMessage.error(error);
        console.log(error);
      }
    };
    const handleLogin = () => {
      router.push("/login");
    };

    return {
      ...toRefs(data),
      router,
      handleRegister,
      handleLogin,
    };
  },
};
</script>

<style>
.register-wrap {
  height: 100vh;
  /* 简约风格：淡灰色渐变背景 */
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.register-window {
  background-color: #ffffff;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 40px 50px;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  width: 400px;
}

.register-item {
  text-align: center;
  margin-bottom: 20px;
  color: #374151;
  font-size: 14px;
}

.register-button-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
  width: 100%;
}

.register-btn,
.register-btn:hover {
  background: #4facfe;
  border: none;
  color: #ffffff;
  font-weight: 500;
  padding: 12px 32px;
  border-radius: 8px;
  width: 100%;
  transition: all 0.3s ease;
}

.register-btn:hover {
  background: #3d8bfe;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.go-login-button-container {
  display: flex;
  flex-direction: row-reverse;
  margin-top: 16px;
}

.go-sms-login-btn,
.go-password-login-btn {
  background-color: transparent;
  border: none;
  cursor: pointer;
  color: #4facfe;
  font-weight: 500;
  font-size: 14px;
  transition: color 0.2s ease;
  text-decoration: underline;
  text-underline-offset: 0.2em;
  margin-left: 10px;
}

.el-alert {
  margin-top: 20px;
}
</style>