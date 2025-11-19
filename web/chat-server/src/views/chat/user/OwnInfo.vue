<template>
  <div class="chat-wrap">
    <div
      class="chat-window"
      :style="{
        boxShadow: `var(${'--el-box-shadow-dark'})`,
      }"
    >
      <el-container class="chat-window-container">
        <el-aside class="aside-container">
          <NavigationModal></NavigationModal>
          <ContactListModal></ContactListModal>
        </el-aside>
        <div class="owner-info-window">
          <div class="my-homepage-title"><h2>我的主页</h2></div>

          <p class="owner-prefix">用户id：{{ userInfo.uuid }}</p>
          <p class="owner-prefix">昵称：{{ userInfo.nickname }}</p>
          <p class="owner-prefix">电话：{{ userInfo.telephone }}</p>
          <p class="owner-prefix">邮箱：{{ userInfo.email }}</p>
          <p class="owner-prefix">
            性别：{{ userInfo.gender === 0 ? "男" : "女" }}
          </p>
          <p class="owner-prefix">生日：{{ userInfo.birthday }}</p>
          <p class="owner-prefix">个性签名：{{ userInfo.signature }}</p>
          <p class="owner-prefix">
            加入kama chat server的时间：{{ userInfo.created_at }}
          </p>
          <div class="owner-opt">
            <p class="owner-prefix">头像：</p>
            <img style="width: 40px; height: 40px" :src="userInfo.avatar" />
          </div>
        </div>
        <div class="edit-window">
          <el-button class="edit-btn" @click="showMyInfoModal">编辑</el-button>
        </div>
        <Modal :isVisible="isMyInfoModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitMyInfoModal">
                  <el-icon><Close /></el-icon>
                </button>
              </div>
              <div class="modal-header-title">
                <h3>修改主页</h3>
              </div>
            </div>
          </template>
          <template v-slot:body>
            <el-scrollbar
              max-height="300px"
              style="
                width: 400px;
                display: flex;
                align-items: center;
                justify-content: center;
                margin-top: 20px;
              "
            >
              <div class="modal-body">
                <el-form ref="formRef" :model="updateInfo" label-width="70px">
                  <el-form-item
                    prop="nickname"
                    label="昵称"
                    :rules="[
                      {
                        min: 3,
                        max: 10,
                        message: '昵称长度在 3 到 10 个字符',
                        trigger: 'blur',
                      },
                    ]"
                  >
                    <el-input
                      v-model="updateInfo.nickname"
                      placeholder="选填"
                    />
                  </el-form-item>
                  <el-form-item prop="email" label="邮箱">
                    <el-input v-model="updateInfo.email" placeholder="选填" />
                  </el-form-item>
                  <el-form-item prop="birthday" label="生日">
                    <el-input
                      v-model="updateInfo.birthday"
                      placeholder="选填，格式为2024.1.1"
                    />
                  </el-form-item>
                  <el-form-item prop="signature" label="个性签名">
                    <el-input
                      v-model="updateInfo.signature"
                      placeholder="选填"
                    />
                  </el-form-item>
                  <el-form-item prop="avatar" label="头像">
                    <el-upload
                      v-model:file-list="fileList"
                      ref="uploadRef"
                      :auto-upload="false"
                      :action="uploadPath"
                      :on-success="handleUploadSuccess"
                      :before-upload="beforeFileUpload"
                    >
                      <template #trigger>
                        <el-button
                          style="background-color: rgb(252, 210.9, 210.9)"
                          >上传图片</el-button
                        >
                      </template>
                    </el-upload>
                  </el-form-item>
                </el-form>
              </div>
            </el-scrollbar>
          </template>
          <template v-slot:footer>
            <div class="modal-footer">
              <el-button class="modal-close-btn" @click="closeMyInfoModal">
                完成
              </el-button>
            </div>
          </template>
        </Modal>
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
import SmallModal from "@/components/SmallModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import ContactListModal from "@/components/ContactListModal.vue";
import { ElMessage } from "element-plus";
export default {
  name: "OwnInfo",
  components: {
    Modal,
    SmallModal,
    ContactListModal,
    NavigationModal,
  },
  setup() {
    const router = useRouter();
    const store = useStore();
    const data = reactive({
      userInfo: store.state.userInfo,
      updateInfo: {
        uuid: "",
        nickname: "",
        email: "",
        birthday: "",
        signature: "",
        avatar: "",
      },
      isMyInfoModalVisible: false,
      ownListReq: {
        owner_id: "",
      },
      uploadRef: null,
      uploadPath: store.state.backendUrl + "/message/uploadAvatar",
      fileList: [],
      cnt: 0,
    });
    const showMyInfoModal = () => {
      data.isMyInfoModalVisible = true;
    };
    const closeMyInfoModal = async () => {
      console.log(data.fileList);
      if (
        data.updateInfo.nickname == "" &&
        data.fileList.length == 0 &&
        data.updateInfo.email == "" &&
        data.updateInfo.birthday == "" &&
        data.updateInfo.signature == ""
      ) {
        ElMessage("请至少修改一项");
        return;
      }
      if (data.updateInfo.nickname != "") {
        if (
          data.updateInfo.nickname.length < 3 ||
          data.updateInfo.nickname.length > 10
        ) {
          return;
        }
      }
      if (data.updateInfo.email != "") {
        if (!checkEmailValid(data.updateInfo.email)) {
          ElMessage("请输入有效的邮箱。");
          return;
        }
      }
      if (data.updateInfo.nickname != "") {
        data.userInfo.nickname = data.updateInfo.nickname;
      }
      if (data.updateInfo.email != "") {
        data.userInfo.email = data.updateInfo.email;
      }
      if (data.fileList.length != 0) {
        try {
          data.updateInfo.avatar = "/static/avatars/" + data.fileList[0].name;
          console.log(data.updateInfo.avatar);
          data.userInfo.avatar = store.state.backendUrl + data.updateInfo.avatar;
          store.commit("setUserInfo", data.userInfo);
          data.uploadRef.submit();
        } catch (error) {
          console.log(error);
        }
      }

      if (data.updateInfo.birthday != "") {
        data.userInfo.birthday = data.updateInfo.birthday;
      }
      if (data.updateInfo.signature != "") {
        data.userInfo.signature = data.updateInfo.signature;
      }
      data.isMyInfoModalVisible = false;
      data.fileList = [];
      data.cnt = 0;
      data.updateInfo.uuid = data.userInfo.uuid;
      store.commit("setUserInfo", data.userInfo);
      try {
        const rsp = await axios.post(
          store.state.backendUrl + "/user/updateUserInfo",
          data.updateInfo
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
        } else if (rsp.data.code == 400) {
          ElMessage.error(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.log(error);
      }
      router.go(0);
    };
    const quitMyInfoModal = () => {
      data.isMyInfoModalVisible = false;
      data.fileList = [];
      data.cnt = 0;
    };
    const handleUploadSuccess = () => {
      ElMessage.success("头像上传成功");
      data.fileList = [];
    };
    const beforeFileUpload = (file) => {
      console.log("上传前file====>", file);
      console.log(data.fileList);
      console.log(file);
      if (data.fileList.length > 1) {
        ElMessage.error("只能上传一张头像");
        return false;
      }
      const isLt50M = file.size / 1024 / 1024 < 50;
      if (!isLt50M) {
        ElMessage.error("上传头像图片大小不能超过 50MB!");
        return false;
      }
    };
    const getFileExtension = (filename) => {
      const parts = filename.split(".");
      return parts.length > 1 ? parts.pop() : "";
    };

    
    return {
      ...toRefs(data),
      router,
      showMyInfoModal,
      closeMyInfoModal,
      quitMyInfoModal,
      handleUploadSuccess,
      beforeFileUpload,
    };
  },
};
</script>

<style scoped>
.owner-info-window {
  width: 84%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.owner-prefix {
  font-family: Arial, Helvetica, sans-serif;
  margin: 6px;
}

.owner-opt {
  margin: 6px;
  display: flex;
  flex-direction: row;
}

.edit-window {
  width: 10%;
  display: flex;
  flex-direction: column-reverse;
}

h3 {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(69, 69, 68);
}

.modal-quit-btn-container {
  height: 30%;
  width: 100%;
  display: flex;
  flex-direction: row-reverse;
}

.modal-quit-btn {
  background-color: rgba(255, 255, 255, 0);
  color: rgb(229, 25, 25);
  padding: 15px;
  border: none;
  cursor: pointer;
  position: fixed;
  justify-content: center;
  align-items: center;
}

.modal-header {
  height: 20%;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.modal-body {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.modal-footer {
  height: 20%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-header-title {
  height: 70%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

h2 {
  margin-bottom: 20px;
  font-family: Arial, Helvetica, sans-serif;
}

.el-menu {
  background-color: rgb(252, 210.9, 210.9);
  width: 101%;
}

.el-menu-item {
  background-color: rgb(255, 255, 255);
  height: 45px;
}
</style>