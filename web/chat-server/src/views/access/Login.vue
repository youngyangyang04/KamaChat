<template>
  <div class="login-wrap">
    <div
      class="login-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <h2 class="login-item">登录</h2>
      <el-form
        :model="loginData"
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
          ]"
        >
          <el-input v-model="loginData.account" placeholder="请输入账号" />
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
          ]"
        >
          <el-input type="password" v-model="loginData.password" placeholder="请输入密码" />
        </el-form-item>
      </el-form>
      <div class="login-button-container">
        <el-button type="primary" class="login-btn" @click="handleLogin"
          >登录</el-button
        >
      </div>

      <div class="go-register-button-container">
        <button class="go-register-btn" @click="handleRegister">注册</button>
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
import { loginAndDeriveMasterKey } from "@/crypto";
import { setCurrentUserId } from "@/crypto/cryptoStore";

export default {
  name: "Login",
  setup() {
    const data = reactive({
      loginData: {
        account: "",
        password: "",
      },
    });
    const router = useRouter();
    const store = useStore();
    const handleLogin = async () => {
      try {
        if (!data.loginData.account || !data.loginData.password) {
          ElMessage.error("请填写完整登录信息。");
          return;
        }
	console.log(store.state.backendUrl, store.state.wsUrl);
        const response = await axios.post("/login", data.loginData);
        console.log(response);
        if (response.data.code == 200) {
          if (response.data.data.status == 1) {
            ElMessage.error("该账号已被封禁，请联系管理员。");
            return;
          }
          try {
            if (!response.data.data.avatar.startsWith("http")) {
              response.data.data.avatar =
                store.state.backendUrl + response.data.data.avatar;
            }
            store.commit("setUserInfo", response.data.data);
            
            // 🔥 关键：必须先设置当前用户 ID，确保 IndexedDB 数据隔离
            // 这样 loginAndDeriveMasterKey 才能读取到正确的 salt
            setCurrentUserId(response.data.data.uuid);
            console.log(`🔐 [Login.vue] 已设置当前用户 ID: ${response.data.data.uuid}`);
            
            // 尝试重新派生主密钥（如果用户启用了加密）
            try {
              const masterKey = await loginAndDeriveMasterKey(data.loginData.password);
              if (masterKey) {
                store.commit("setMasterKey", masterKey);
                console.log("主密钥已重新派生并保存到内存（加密功能已启用）");
                
                // 检查 sessionStorage 中是否有主密钥（说明用户开启了"保存主密钥"开关）
                const savedMasterKey = sessionStorage.getItem('masterKey');
                if (savedMasterKey) {
                  // 如果开关是打开的，更新保存的主密钥
                  store.commit("saveMasterKeyToStorage", masterKey);
                  console.log("✅ 已更新 sessionStorage 中的主密钥");
                }
                
                ElMessage.success(response.data.message + " (端到端加密已启用)");
              } else {
                console.log("主密钥验证失败，可能未启用加密或密码错误");
                ElMessage.success(response.data.message);
              }
            } catch (error) {
              console.log("未找到加密密钥，使用普通模式:", error.message);
              ElMessage.success(response.data.message);
            }
            
            // 准备创建websocket连接
            const wsUrl =
              store.state.wsUrl + "/wss?client_id=" + response.data.data.uuid + "&token=" + encodeURIComponent(response.data.data.token);
            console.log(wsUrl);
            store.state.socket = new WebSocket(wsUrl);
            store.state.socket.onopen = () => {
              console.log("🌐 [Login.vue] WebSocket连接已打开");
              console.log("🌐 [Login.vue] 连接建立后，App.vue 将自动设置全局消息处理器");
            };
            // 不在这里设置 onmessage，让 App.vue 统一管理
            // App.vue 中的 watch 会监听到 socket 的变化并设置全局处理器
            store.state.socket.onclose = () => {
              console.log("🌐 [Login.vue] WebSocket连接已关闭");
            };
            store.state.socket.onerror = (error) => {
              console.log("🌐 [Login.vue] WebSocket连接发生错误", error);
            };
            router.push("/chat/sessionlist");
          } catch (error) {
            console.log(error);
          }
        } else {
          ElMessage.error(response.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
      }
    };
    const handleRegister = () => {
      router.push("/register");
    };

    return {
      ...toRefs(data),
      router,
      handleLogin,
      handleRegister,
    };
  },
};
</script>

<style>
.login-wrap {
  height: 100vh;
  /* 简约风格：淡灰色渐变背景 */
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.login-window {
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

.login-item {
  text-align: center;
  margin-bottom: 20px;
  color: #374151;
  font-size: 14px;
}

.login-button-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
  width: 100%;
}

.login-btn,
.login-btn:hover {
  background: #4facfe;
  border: none;
  color: #ffffff;
  font-weight: 500;
  padding: 12px 32px;
  border-radius: 8px;
  width: 100%;
  transition: all 0.3s ease;
}

.login-btn:hover {
  background: #3d8bfe;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.go-register-button-container {
  display: flex;
  flex-direction: row-reverse;
  margin-top: 16px;
}

.go-register-btn,
.go-sms-btn {
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
