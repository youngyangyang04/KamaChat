<template>
  <router-view />
</template>

<script>
import { onMounted } from "vue";
import { useStore } from "vuex";
import axios from "@/utils/axios";
export default {
  name: "App",
  setup() {
    const store = useStore();
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
    const logout = async () => {
      store.commit("cleanUserInfo");
      const req = {
        owner_id: data.userInfo.uuid,
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
    onMounted(() => {
      if (store.state.userInfo.uuid) {
        getUserInfo();
        if (store.state.userInfo.status == 1) {
          logout();
        }
        const wsUrl =
          store.state.wsUrl + "/wss?client_id=" + store.state.userInfo.uuid + "&token=" + encodeURIComponent(store.state.token);
          console.log(wsUrl);
        store.state.socket = new WebSocket(wsUrl);
        store.state.socket.onopen = () => {
          console.log("WebSocket连接已打开");console.log("连接信令服务器成功");
        };
        store.state.socket.onmessage = (message) => {
          console.log("收到消息：", message.data);
        };
        store.state.socket.onclose = () => {
          console.log("WebSocket连接已关闭");
        console.log("连接信令服务器断开");
        };
        store.state.socket.onerror = (error) => {
          console.log("WebSocket连接发生错误");
          console.log("连接信令服务器失败，错误信息：", error);
        };
        console.log(store.state.socket);
      }
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