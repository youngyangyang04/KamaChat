<template>
  <div class="chat-wrap">
    <div
      class="chat-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <el-container class="chat-window-container">
        <el-header class="header-container">
        <div style="display: flex; align-items: center; flex-direction: row;">
          <svg
            t="1736782942164"
            class="header-icon"
            viewBox="0 0 1024 1024"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            p-id="2358"
            width="200"
            height="200"
          >
            <path
              d="M512 597.333333v341.333334H170.666667a341.333333 341.333333 0 0 1 341.333333-341.333334z m0-42.666666c-141.44 0-256-114.56-256-256s114.56-256 256-256 256 114.56 256 256-114.56 256-256 256z m384 170.666666h42.666667v213.333334h-341.333334v-213.333334h42.666667v-42.666666a128 128 0 1 1 256 0v42.666666z m-85.333333 0v-42.666666a42.666667 42.666667 0 1 0-85.333334 0v42.666666h85.333334z"
              p-id="2359"
              fill="#2c2c2c"
            ></path>
          </svg>
          <div class="header-title">admin</div>
          </div>
          <el-button @click="backToChat">返回</el-button>
        </el-header>

        <el-container>
          <el-aside class="manager-aside-container">
            <el-menu unique-opened>
              <el-sub-menu index="1">
                <template #title>
                  <svg
                    t="1736823523206"
                    class="manager-user-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="1477"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M213.333333 896H128v-42.666667a299.008 299.008 0 0 1 298.666667-298.666666h170.666666a299.008 299.008 0 0 1 298.666667 298.666666v42.666667h-85.333333v-42.666667a213.589333 213.589333 0 0 0-213.333334-213.333333h-170.666666a213.589333 213.589333 0 0 0-213.333334 213.333333v42.666667z m298.666667-384a213.333333 213.333333 0 1 1 0-426.666667 213.333333 213.333333 0 0 1 0 426.666667z m0-341.333333a128 128 0 1 0 0 256 128 128 0 0 0 0-256z"
                      p-id="1478"
                    ></path>
                  </svg>
                  <span class="manager-user-title">用户</span>
                </template>

                <el-menu-item index="1-1" @click="showDisableUserTable"
                  ><svg
                    t="1736825139556"
                    class="menu-item-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="4645"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M512 853.333333c-187.733333 0-341.333333-153.6-341.333333-341.333333s153.6-341.333333 341.333333-341.333333 341.333333 153.6 341.333333 341.333333-153.6 341.333333-341.333333 341.333333z m0-85.333333c140.8 0 256-115.2 256-256s-115.2-256-256-256-256 115.2-256 256 115.2 256 256 256zM401.066667 341.333333l273.066666 273.066667-59.733333 59.733333L341.333333 401.066667 401.066667 341.333333z"
                      fill="#444444"
                      p-id="4646"
                    ></path>
                  </svg>
                  启用/禁用
                </el-menu-item>
                <el-menu-item index="1-2" @click="showDeleteUserTable">
                  <svg
                    t="1736825210236"
                    class="menu-item-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="5694"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M736.653061 929.959184H287.346939c-45.97551 0-83.591837-37.616327-83.591837-83.591837V355.265306h616.489796v491.102041c0 45.97551-37.616327 83.591837-83.591837 83.591837zM245.55102 397.061224v449.306123c0 22.987755 18.808163 41.795918 41.795919 41.795918h449.306122c22.987755 0 41.795918-18.808163 41.795919-41.795918V397.061224H245.55102z"
                      fill="#333333"
                      p-id="5695"
                    ></path>
                    <path
                      d="M407.510204 794.122449c-11.493878 0-20.897959-9.404082-20.897959-20.897959V512c0-11.493878 9.404082-20.897959 20.897959-20.897959s20.897959 9.404082 20.897959 20.897959v261.22449c0 11.493878-9.404082 20.897959-20.897959 20.897959zM616.489796 794.122449c-11.493878 0-20.897959-9.404082-20.897959-20.897959V512c0-11.493878 9.404082-20.897959 20.897959-20.897959s20.897959 9.404082 20.897959 20.897959v261.22449c0 11.493878-9.404082 20.897959-20.897959 20.897959z"
                      fill="#333333"
                      p-id="5696"
                    ></path>
                    <path
                      d="M846.367347 397.061224H177.632653c-45.97551 0-83.591837-37.616327-83.591837-83.591836v-31.346939c0-45.97551 37.616327-83.591837 83.591837-83.591837h668.734694c45.97551 0 83.591837 37.616327 83.591837 83.591837v31.346939c0 45.97551-37.616327 83.591837-83.591837 83.591836zM177.632653 240.326531c-22.987755 0-41.795918 18.808163-41.795918 41.795918v31.346939c0 22.987755 18.808163 41.795918 41.795918 41.795918h668.734694c22.987755 0 41.795918-18.808163 41.795918-41.795918v-31.346939c0-22.987755-18.808163-41.795918-41.795918-41.795918H177.632653z"
                      fill="#333333"
                      p-id="5697"
                    ></path>
                    <path
                      d="M650.44898 240.326531h-276.89796c-28.734694 0-52.244898-23.510204-52.244898-52.244898v-41.795919c0-28.734694 23.510204-52.244898 52.244898-52.244898h276.89796c28.734694 0 52.244898 23.510204 52.244898 52.244898v41.795919c0 28.734694-23.510204 52.244898-52.244898 52.244898z m-276.89796-104.489796c-5.746939 0-10.44898 4.702041-10.448979 10.448979v41.795919c0 5.746939 4.702041 10.44898 10.448979 10.448979h276.89796c5.746939 0 10.44898-4.702041 10.448979-10.448979v-41.795919c0-5.746939-4.702041-10.44898-10.448979-10.448979h-276.89796z"
                      fill="#333333"
                      p-id="5698"
                    ></path></svg
                  >删除
                </el-menu-item>
                <el-menu-item index="1-3" @click="showSetAdminTable">
                  <svg
                    t="1736825411779"
                    class="menu-item-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="6703"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M755.2 780.8L896 640l59.733333 59.733333-200.533333 200.533334-119.466667-119.466667 59.733334-59.733333 59.733333 59.733333zM597.333333 768v85.333333H256v-170.666666h4.266667l85.333333-85.333334h375.466667l85.333333 85.333334h-426.666667l-38.4 38.4V768h256z m213.333334-81.066667z m-85.333334-324.266666c0 106.666667-85.333333 192-192 192S341.333333 469.333333 341.333333 362.666667 426.666667 170.666667 533.333333 170.666667 725.333333 256 725.333333 362.666667z m-85.333333 0C640 302.933333 593.066667 256 533.333333 256S426.666667 302.933333 426.666667 362.666667s46.933333 106.666667 106.666666 106.666666S640 422.4 640 362.666667z"
                      fill="#444444"
                      p-id="6704"
                    ></path></svg
                  >设置为管理员
                </el-menu-item>
              </el-sub-menu>
              <el-sub-menu index="2">
                <template #title>
                  <svg
                    t="1736823667080"
                    class="manager-group-icon"
                    viewBox="0 0 1143 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="3627"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M457.238151 0.006532a296.045313 296.045313 0 0 1 285.349798 215.379905 243.465837 243.465837 0 0 1 44.333317-4.082258 237.669031 237.669031 0 0 1 143.777108 428.963617 379.07843 379.07843 0 0 1 212.27739 338.745726v44.904833H38.725113A38.699801 38.699801 0 0 1 0.025312 985.218554a473.541869 473.541869 0 0 1 295.310507-443.578099A294.41241 294.41241 0 0 1 457.238151 0.006532z m-2.939226 589.396336a384.385365 384.385365 0 0 0-374.751237 357.19753h749.502475a384.385365 384.385365 0 0 0-374.751238-357.19753z m358.340563 96.831147h-6.041741a482.359545 482.359545 0 0 1 100.096953 259.713222h156.921978v-0.571516a291.799765 291.799765 0 0 0-251.303771-259.141706z m-25.718222-397.530234a165.902944 165.902944 0 0 0-33.556157 3.429096v2.531a293.922539 293.922539 0 0 1-136.837271 248.446191 450.109711 450.109711 0 0 1 133.32653 82.951472 38.699801 38.699801 0 0 1 33.229576-18.86003q14.369546 0 28.412512 1.143033a160.59601 160.59601 0 0 0-24.493545-319.640762zM457.238151 77.406134a217.257744 217.257744 0 1 0 218.809001 217.257743A218.074195 218.074195 0 0 0 457.238151 77.406134z"
                      p-id="3628"
                    ></path>
                  </svg>
                  <span class="manager-group-title">群聊</span>
                </template>
                <el-menu-item index="2-1" @click="showDisableGroupTable">
                  <svg
                    t="1736825139556"
                    class="menu-item-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="4645"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M512 853.333333c-187.733333 0-341.333333-153.6-341.333333-341.333333s153.6-341.333333 341.333333-341.333333 341.333333 153.6 341.333333 341.333333-153.6 341.333333-341.333333 341.333333z m0-85.333333c140.8 0 256-115.2 256-256s-115.2-256-256-256-256 115.2-256 256 115.2 256 256 256zM401.066667 341.333333l273.066666 273.066667-59.733333 59.733333L341.333333 401.066667 401.066667 341.333333z"
                      fill="#444444"
                      p-id="4646"
                    ></path>
                  </svg>
                  启用/禁用</el-menu-item
                >
                <el-menu-item index="2-2" @click="showDeleteGroupTable">
                  <svg
                    t="1736825642071"
                    class="menu-item-icon"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    p-id="7972"
                    width="200"
                    height="200"
                  >
                    <path
                      d="M960 384h-80V144H144v666.519L231.778 752H448v80H256L64 960V64h896v320z m-704-80h512v80H256v-80z m0 208h192v80H256v-80z"
                      fill="#2c2c2c"
                      p-id="7973"
                    ></path>
                    <path
                      d="M817.497 704l74.247-74.246-49.498-49.498L768 654.503l-74.246-74.247-49.498 49.498L718.503 704l-74.247 74.246 49.498 49.498L768 753.497l74.246 74.247 49.498-49.498L817.497 704zM767.5 960.107C625.798 960.107 511 845.095 511 703.5 511 561.905 625.798 447 767.5 447S1024 562.012 1024 703.607c0 141.595-114.798 256.5-256.5 256.5z"
                      fill="#2c2c2c"
                      p-id="7974"
                    ></path>
                  </svg>
                  删除/解散群聊</el-menu-item
                >
              </el-sub-menu>
            </el-menu>
          </el-aside>
          <el-main class="manager-main-container"> 
            <DeleteUserModal :isVisible="isDeleteUserModalVisible"></DeleteUserModal>
            <DisableUserModal :isVisible="isDisableUserModalVisible"></DisableUserModal>
            <SetAdminModal :isVisible="isSetAdminModalVisible"></SetAdminModal>
            <DeleteGroupModal :isVisible="isDeleteGroupModalVisible"></DeleteGroupModal>
            <DisableGroupModal :isVisible="isDisableGroupModalVisible"></DisableGroupModal>
          </el-main>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted, ref } from "vue";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { useRouter } from "vue-router";
import Modal from "@/components/Modal.vue";
import { checkEmailValid } from "@/assets/js/valid.js";
import { generateString } from "@/assets/js/random.js";
import DisableUserModal from "@/components/DisableUserModal.vue";
import DeleteUserModal from "@/components/DeleteUserModal.vue";
import SetAdminModal from "@/components/SetAdminModal.vue";
import DeleteGroupModal from "@/components/DeleteGroupModal.vue";
import DisableGroupModal from "@/components/DisableGroupModal.vue";
import { ElMessage } from "element-plus";
export default {
  name: "Manager",
  components: {
    DisableUserModal,
    DeleteUserModal,
    SetAdminModal,
    DeleteGroupModal,
    DisableGroupModal,
  },
  setup() {
    const router = useRouter();
    const data = reactive({
      isDisableUserModalVisible: false,
      isDeleteUserModalVisible: false,
      isSetAdminModalVisible: false,
      isDeleteGroupModalVisible: false,
      isDisableGroupModalVisible: false,
    });
    const showDisableUserTable = () => {
      data.isDisableUserModalVisible = true;
      data.isDeleteUserModalVisible = false;
      data.isSetAdminModalVisible = false;
      data.isDeleteGroupModalVisible = false;
      data.isDisableGroupModalVisible = false;
    };

    const showDeleteUserTable = () => {
      data.isDeleteUserModalVisible = true;
      data.isDisableUserModalVisible = false;
      data.isSetAdminModalVisible = false;
      data.isDeleteGroupModalVisible = false;
      data.isDisableGroupModalVisible = false;
    }

    const showSetAdminTable = () => {
      data.isSetAdminModalVisible = true;
      data.isDeleteUserModalVisible = false;
      data.isDisableUserModalVisible = false;
      data.isDeleteGroupModalVisible = false;
      data.isDisableGroupModalVisible = false;
    }

    const showDeleteGroupTable = () => {
      data.isSetAdminModalVisible = false;
      data.isDeleteUserModalVisible = false;
      data.isDisableUserModalVisible = false;
      data.isDeleteGroupModalVisible = true;
      data.isDisableGroupModalVisible = false;
    };

    const showDisableGroupTable = () => {
      data.isSetAdminModalVisible = false;
      data.isDeleteUserModalVisible = false;
      data.isDisableUserModalVisible = false;
      data.isDeleteGroupModalVisible = false;
      data.isDisableGroupModalVisible = true;
    };

    const backToChat = () => {
      router.push("/chat/sessionlist");
    }

    return {
      ...toRefs(data),
      router,
      showDisableUserTable,
      showDeleteUserTable,
      showSetAdminTable,
      showDeleteGroupTable,
      showDisableGroupTable,
      backToChat,
    };
  },
};
</script>

<style scoped>
.header-container {
  height: 70px;
  border-bottom: 3px solid #ccc;
  background-color: rgb(252, 210.9, 210.9);
  border-top-right-radius: 30px;
  border-top-left-radius: 30px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.el-menu {
  width: 196px;
  background-color: rgb(252, 210.9, 210.9);
}

.el-menu-item {
  font-family: Arial, Helvetica, sans-serif;
}

.manager-aside-container {
  width: 200px;
  border-right: 3px solid #ccc;
  display: flex;
  flex-direction: column;
  border-bottom-left-radius: 30px;
}

.header-icon {
  height: 40px;
  width: 40px;
  margin-right: 20px;
}

.header-title {
  font-family: Arial, Helvetica, sans-serif;
  font-weight: bold;
  color: rgb(37, 37, 37);
  font-size: 30px;
}

.user-manager-icon {
  height: 20px;
  width: 20px;
}

.manager-user-title {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(37, 37, 37);
}

.manager-user-icon {
  height: 20px;
  width: 20px;
  margin-right: 10px;
}

.manager-group-title {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(37, 37, 37);
}

.manager-group-icon {
  height: 20px;
  width: 20px;
  margin-right: 10px;
}

.menu-item-icon {
  height: 20px;
  width: 20px;
  margin-right: 10px;
}

</style>