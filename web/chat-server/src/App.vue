<template>
  <router-view />
</template>

<script>
import { onMounted, onUnmounted, watch } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import axios from "@/utils/axios";
import { ElMessage } from "element-plus";
import eventBus from "@/utils/eventBus";
import { setCurrentUserId } from "@/crypto/cryptoStore";

export default {
  name: "App",
  setup() {
    const store = useStore();
    const router = useRouter();
    
    const getUserInfo = async () => {
      try {
        const req = {
          uuid: store.state.userInfo.uuid,
        };
        const rsp = await axios.post("/user/getUserInfo", req);
        if (rsp.data.code == 200) {
          if (!rsp.data.data.avatar.startsWith("http")) {
            rsp.data.data.avatar = store.state.backendUrl + rsp.data.data.avatar;
          }
          store.commit("setUserInfo", rsp.data.data);
        } else {
          console.error(rsp.data.message);
        }
        console.log(rsp);
      } catch (error) {
        console.log(error);
      }
    };
    
    // 获取未读通知数量
    const getUnreadNotificationCount = async () => {
      try {
        console.log("🔔 [App.vue] 开始获取未读通知数量...");
        const req = {
          user_id: store.state.userInfo.uuid,
          type: null, // null 表示获取所有类型的未读数量
        };
        const rsp = await axios.post("/notification/getUnreadCount", req);
        if (rsp.data.code === 200) {
          // 后端返回的数据格式可能是 {count: 3} 或者直接是 3
          let unreadCount = 0;
          if (typeof rsp.data.data === 'object' && rsp.data.data !== null) {
            unreadCount = rsp.data.data.count || 0;
          } else {
            unreadCount = rsp.data.data || 0;
          }
          console.log("🔔 [App.vue] 获取到未读通知数量:", unreadCount);
          store.commit('setUnreadNotificationCount', unreadCount);
          console.log("🔔 [App.vue] 已更新 store 中的未读数量");
        } else {
          console.error("🔔 [App.vue] 获取未读数量失败:", rsp.data.message);
        }
      } catch (error) {
        console.error("🔔 [App.vue] 获取未读通知数量出错:", error);
      }
    };
    
    const logout = async () => {
      store.commit("cleanUserInfo");
      const req = {
        owner_id: store.state.userInfo.uuid,
      };
      const rsp = await axios.post(
        store.state.backendUrl + "/user/wsLogout",
        req
      );
      if (rsp.data.code == 200) {
        router.push("/login");
        ElMessage.success("账号被封禁，退出登录");
      } else {
        ElMessage.error(rsp.data.message);
      }
    };
    
    // 全局 WebSocket 消息处理器 - 作为唯一的消息入口
    const handleWebSocketMessage = async (jsonMessage) => {
      try {
        const message = JSON.parse(jsonMessage.data);
        console.log("🌐 [App.vue] 全局收到 WebSocket 消息：", message);
        
        // 1. 优先处理通知推送消息
        if (message.type === 'notification') {
          console.log("🔔 [App.vue] 收到通知推送，未读数量:", message.unread_count);
          // 只更新未读通知数量，不处理通知对象
          // 前端打开通知界面时会自动从后端获取完整的通知列表
          if (message.unread_count !== undefined && message.unread_count !== null) {
            const count = Number(message.unread_count);
            console.log("🔔 [App.vue] 更新前 store 中的未读数量:", store.state.unreadNotificationCount);
            store.commit('setUnreadNotificationCount', count);
            console.log("🔔 [App.vue] 更新后 store 中的未读数量:", store.state.unreadNotificationCount);
          } else {
            console.warn("🔔 [App.vue] 未读数量无效:", message.unread_count);
          }
          return; // 通知消息处理完毕，不再传递给其他组件
        }
        
        // 2. 处理聊天消息（文本、文件等）
        if (typeof message.type === 'number' && message.type !== 3) {
          console.log("💬 [App.vue] 收到聊天消息，通过事件总线分发");
          // 通过事件总线分发给 ContactChat.vue
          eventBus.emit('chat:message', message);
          return;
        }
        
        // 3. 处理 AV 消息（type = 3）
        if (message.type === 3) {
          console.log("📹 [App.vue] 收到 AV 消息，通过事件总线分发");
          // 通过事件总线分发给 ContactChat.vue
          eventBus.emit('chat:av_message', message);
          return;
        }
        
        // 4. 其他未知类型的消息
        console.warn("⚠️ [App.vue] 收到未知类型的消息:", message);
      } catch (error) {
        console.error("🌐 [App.vue] 处理 WebSocket 消息失败：", error);
        console.error("🌐 [App.vue] 错误堆栈:", error.stack);
      }
    };
    
    // 设置 WebSocket 消息处理器的函数（可复用）
    const setupWebSocketHandler = () => {
      if (store.state.socket) {
        console.log("🌐 [App.vue] 设置唯一的全局消息处理器");
        // 清除所有旧的监听器
        const oldOnMessage = store.state.socket.onmessage;
        if (oldOnMessage) {
          store.state.socket.removeEventListener('message', oldOnMessage);
        }
        // 设置唯一的全局消息处理器
        store.state.socket.onmessage = handleWebSocketMessage;
        console.log("🌐 [App.vue] ✅ 已设置唯一的全局消息处理器，所有消息将在此统一处理");
      }
    };
    
    onMounted(() => {
      if (store.state.userInfo.uuid) {
        // 设置当前用户 ID，确保 IndexedDB 数据隔离
        setCurrentUserId(store.state.userInfo.uuid);
        console.log(`🔐 [App.vue] 已设置当前用户 ID: ${store.state.userInfo.uuid}`);
        
        getUserInfo();
        
        // 初始化时获取未读通知数量
        getUnreadNotificationCount();
        
        if (store.state.userInfo.status == 1) {
          logout();
        }
        
        // 如果 WebSocket 已存在（可能是 Login.vue 创建的），重新设置消息处理器
        if (store.state.socket && store.state.socket.readyState === WebSocket.OPEN) {
          console.log("🌐 [App.vue] WebSocket 已存在且已连接");
          setupWebSocketHandler();
          return;
        }
        
        // 如果 WebSocket 不存在，创建新的连接
        const wsUrl =
          store.state.wsUrl + "/wss?client_id=" + store.state.userInfo.uuid + "&token=" + encodeURIComponent(store.state.token);
        console.log("🌐 [App.vue] 创建新的 WebSocket 连接:", wsUrl);
        
        store.state.socket = new WebSocket(wsUrl);
        store.state.socket.onopen = () => {
          console.log("🌐 [App.vue] WebSocket连接已打开");
          console.log("🌐 [App.vue] 连接信令服务器成功");
        };
        // 设置唯一的全局消息处理器
        setupWebSocketHandler();
        
        store.state.socket.onclose = () => {
          console.log("🌐 [App.vue] WebSocket连接已关闭");
          console.log("🌐 [App.vue] 连接信令服务器断开");
        };
        store.state.socket.onerror = (error) => {
          console.log("🌐 [App.vue] WebSocket连接发生错误");
          console.log("🌐 [App.vue] 连接信令服务器失败，错误信息：", error);
        };
        console.log("🌐 [App.vue] WebSocket 对象:", store.state.socket);
      }
    });
    
    // 监听 WebSocket 对象的变化（登录后会创建新的 WebSocket）
    watch(
      () => store.state.socket,
      (newSocket, oldSocket) => {
        if (newSocket && newSocket !== oldSocket) {
          console.log("🌐 [App.vue] 检测到 WebSocket 对象变化，重新设置消息处理器");
          // 等待连接建立后设置处理器
          if (newSocket.readyState === WebSocket.OPEN) {
            setupWebSocketHandler();
          } else {
            // 监听 onopen 事件
            const originalOnOpen = newSocket.onopen;
            newSocket.onopen = (event) => {
              console.log("🌐 [App.vue] WebSocket 连接已建立（通过 watch 监听）");
              if (originalOnOpen) {
                originalOnOpen.call(newSocket, event);
              }
              setupWebSocketHandler();
            };
          }
        }
      },
      { immediate: true }
    );
    
    // 监听用户信息的变化（登录后会设置用户信息）
    watch(
      () => store.state.userInfo.uuid,
      (newUuid, oldUuid) => {
        if (newUuid && newUuid !== oldUuid) {
          console.log("🔔 [App.vue] 检测到用户登录");
          
          // 设置当前用户 ID，确保 IndexedDB 数据隔离
          setCurrentUserId(newUuid);
          console.log(`🔐 [App.vue] 已设置当前用户 ID: ${newUuid}`);
          
          // 获取未读通知数量
          getUnreadNotificationCount();
        }
      }
    );
    
    onUnmounted(() => {
      // 清理事件总线
      eventBus.clear();
      console.log("🌐 [App.vue] onUnmounted: 已清理事件总线");
    });
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box; /* 推荐使用，以确保布局计算的一致性 */
}
</style>