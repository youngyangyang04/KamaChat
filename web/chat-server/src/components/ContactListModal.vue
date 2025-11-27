<template>
  <div class="contactlist-container">
    <div class="contactlist-header">
      <el-input
        v-model="contactSearch"
        class="contact-search-input"
        placeholder="搜索联系人/群聊"
        size="small"
        suffix-icon="Search"
      />
      <div class="contactlist-header-right">
        <el-dropdown placement="bottom" trigger="click">
          <button class="create-group-btn">
            <svg
              t="1733664667695"
              class="create-group-icon"
              viewBox="0 0 1024 1024"
              version="1.1"
              xmlns="http://www.w3.org/2000/svg"
              p-id="2875"
              width="128"
              height="128"
            >
              <path
                d="M488.021333 96a248.021333 248.021333 0 1 1-17.92 495.36l-1.749333 0.341333-4.352 0.298667A304 304 0 0 0 160 896a32 32 0 1 1-64 0 368.170667 368.170667 0 0 1 250.026667-348.672A247.978667 247.978667 0 0 1 488.021333 96z m288 528a32 32 0 0 1 32 32l-0.042666 87.978667H896a32 32 0 0 1 31.701333 27.690666l0.298667 4.352a32 32 0 0 1-32 32l-88.021333-0.042666V896a32 32 0 0 1-27.648 31.701333l-4.352 0.298667a32 32 0 0 1-32-32v-88.021333h-87.978667a32 32 0 0 1-31.701333-27.648l-0.298667-4.352a32 32 0 0 1 32-32h87.978667v-87.978667a32 32 0 0 1 27.690666-31.701333zM488.021333 160a184.021333 184.021333 0 1 0 0 368 184.021333 184.021333 0 0 0 0-368z"
                fill="#2c2c2c"
                p-id="2876"
              ></path>
            </svg>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="showCreateGroupModal">
                创建群聊
              </el-dropdown-item>
              <el-dropdown-item @click="showApplyContactModal">
                添加用户/群聊
              </el-dropdown-item>
              <el-dropdown-item @click="showNewContactModal">
                新的好友
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <SmallModal :isVisible="isNewContactModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitNewContactModal">
                  <el-icon><Close /></el-icon>
                </button>
              </div>
              <div class="modal-header-title">
                <h3>新的朋友</h3>
              </div>
            </div>
          </template>
          <template v-slot:body>
            <div class="newcontact-modal-body">
              <el-scrollbar max-height="400px">
                <ul class="newcontact-list" style="list-style-type: none">
                  <li
                    v-for="newContact in newContactList"
                    :key="newContact.contact_id"
                    class="newcontact-item"
                  >
                    <div
                      style="
                        display: flex;
                        align-items: center;
                        justify-content: center;
                      "
                    >
                      <img
                        :src="newContact.contact_avatar"
                        style="width: 30px; height: 30px; margin-right: 10px"
                      />

                      <el-tooltip
                        effect="customized"
                        :content="newContact.message"
                        placement="top"
                        hide-after="0"
                        enterable="false"
                      >
                        <div>
                          {{ newContact.contact_name }}
                        </div>
                      </el-tooltip>
                    </div>
                    <el-dropdown placement="right" trigger="click">
                      <el-button class="action-btn"> 去处理 </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item
                            @click="handleAgree(newContact.contact_id)"
                            >同意</el-dropdown-item
                          >
                          <el-dropdown-item
                            @click="handleReject(newContact.contact_id)"
                          >
                            拒绝
                          </el-dropdown-item>
                          <el-dropdown-item
                            @click="handleBlack(newContact.contact_id)"
                          >
                            拉黑
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </li>
                </ul>
              </el-scrollbar>
            </div>
          </template>
          <template v-slot:footer>
            <div class="newcontact-modal-footer"></div>
          </template>
        </SmallModal>
        <SmallModal :isVisible="isApplyContactModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitApplyContactModal">
                  <el-icon><Close /></el-icon>
                </button>
              </div>
              <div class="modal-header-title">
                <h3>添加用户/群聊</h3>
              </div>
            </div>
          </template>
          <template v-slot:body>
            <div class="modal-body">
              <el-form
                ref="formRef"
                :model="applyContactReq"
                label-width="100px"
                class="apply-contact-form"
              >
                <el-form-item
                  prop="name"
                  label="用户/群聊id"
                  :rules="[
                    {
                      required: true,
                      message: '此项为必填项',
                      trigger: 'blur',
                    },
                  ]"
                >
                  <el-input
                    v-model="applyContactReq.contact_id"
                    placeholder="请填写申请的用户/群聊id"
                  />
                </el-form-item>
                <el-form-item prop="message" label="申请消息">
                  <el-input
                    v-model="applyContactReq.message"
                    placeholder="选填，填写更容易通过"
                    type="textarea"
                    show-word-limit
                    maxlength="100"
                    :autosize="{ minRows: 3, maxRows: 3 }"
                  />
                </el-form-item>
              </el-form>
            </div>
          </template>
          <template v-slot:footer>
            <div class="modal-footer">
              <el-button
                class="modal-close-btn"
                @click="closeApplyContactModal"
              >
                完成
              </el-button>
            </div>
          </template>
        </SmallModal>
        <Modal :isVisible="isCreateGroupModalVisible">
          <template v-slot:header>
            <div class="modal-header">
              <div class="modal-quit-btn-container">
                <button class="modal-quit-btn" @click="quitCreateGroupModal">
                  <el-icon><Close /></el-icon>
                </button>
              </div>
              <div class="modal-header-title">
                <h3>创建群聊</h3>
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
              <div class="creatgroup-modal-body">
                <el-form
                  ref="formRef"
                  :model="createGroupReq"
                  label-width="80px"
                  class="demo-dynamic"
                >
                  <el-form-item
                    prop="name"
                    label="群名称"
                    :rules="[
                      {
                        required: true,
                        message: '此项为必填项',
                        trigger: 'blur',
                      },
                    ]"
                  >
                    <el-input
                      v-model="createGroupReq.name"
                      placeholder="必填"
                    />
                  </el-form-item>
                  <el-form-item prop="notice" label="群公告">
                    <el-input
                      v-model="createGroupReq.notice"
                      type="textarea"
                      show-word-limit
                      maxlength="500"
                      :autosize="{ minRows: 3, maxRows: 3 }"
                      placeholder="选填"
                    />
                  </el-form-item>
                  <el-form-item
                    prop="add_mode"
                    label="加群方式"
                    :rules="[
                      {
                        required: true,
                        message: 'Please select activity resource',
                        trigger: 'change',
                      },
                    ]"
                  >
                    <el-radio-group v-model="createGroupReq.add_mode">
                      <el-radio :value="0">直接加入</el-radio>
                      <el-radio :value="1">群主审核</el-radio>
                    </el-radio-group>
                  </el-form-item>
                  <el-form-item prop="avatar" label="群头像">
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
            <div class="creategroup-modal-footer">
              <el-button class="modal-close-btn" @click="closeCreateGroupModal">
                完成
              </el-button>
            </div>
          </template>
        </Modal>
      </div>
    </div>
    <div class="contactlist-body">
      <div class="contactlist-user">
        <!-- 通知分类 - 放在联系人上面 -->
        <div class="notification-category-list">
          <div 
            class="notification-category-item"
            :class="{ 'active': store.state.notificationFilterType === 'friend' }"
            @click="handleSelectNotificationCategory('friend')"
          >
            <span class="notification-category-text">好友通知</span>
            <el-icon class="notification-category-arrow"><ArrowRight /></el-icon>
          </div>
          <div class="notification-category-divider"></div>
          <div 
            class="notification-category-item"
            :class="{ 'active': store.state.notificationFilterType === 'group' }"
            @click="handleSelectNotificationCategory('group')"
          >
            <span class="notification-category-text">群通知</span>
            <el-icon class="notification-category-arrow"><ArrowRight /></el-icon>
          </div>
          <div class="notification-category-divider"></div>
          <div 
            class="notification-category-item"
            :class="{ 'active': store.state.notificationFilterType === 'system' }"
            @click="handleSelectNotificationCategory('system')"
          >
            <span class="notification-category-text">系统消息</span>
            <el-icon class="notification-category-arrow"><ArrowRight /></el-icon>
          </div>
        </div>

        <el-menu
          router
          unique-opened
          @open="handleShowUserList"
          @close="handleHideUserList"
        >
          <el-sub-menu index="1">
            <template #title>
              <span class="contactlist-user-title">联系人</span>
            </template>
          </el-sub-menu>
          <el-menu-item
            v-for="user in contactUserList"
            :key="user.user_id"
            @click="handleToChatUser(user)"
            class="contactlist-user-menu-item"
          >
            <el-dropdown
              trigger="contextmenu"
              class="contactlist-dropdown"
              placement="right"
            >
              <div class="contactlist-user-item">
                <img :src="user.avatar" class="contactlist-user-avatar" />
                {{ user.user_name }}
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleCancelBlack(user)"
                    >解除拉黑</el-dropdown-item
                  >
                  <!-- 其他菜单项 -->
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </el-menu-item>
        </el-menu>

        <el-menu
          router
          unique-opened
          @open="handleShowMyGroupList"
          @close="handleHideMyGroupList"
        >
          <el-sub-menu index="1">
            <template #title>
              <span class="contactlist-user-title">我创建的群聊</span>
            </template>
          </el-sub-menu>
          <el-menu-item
            v-for="group in myGroupList"
            :key="group.group_id"
            @click="handleToChatGroup(group)"
          >
            <img :src="group.avatar" class="contactlist-avatar" />
            {{ group.group_name }}
          </el-menu-item>
        </el-menu>
        <el-menu
          router
          unique-opened
          @open="handleShowMyJoinedGroupList"
          @close="handleHideMyJoinedGroupList"
        >
          <el-sub-menu index="1">
            <template #title>
              <span class="contactlist-user-title">我加入的群聊</span>
            </template>
          </el-sub-menu>
          <el-menu-item
            v-for="group in myJoinedGroupList"
            :key="group.group_id"
            @click="handleToChatGroup(group)"
          >
            <img :src="group.avatar" class="contactlist-avatar" />
            {{ group.group_name }}
          </el-menu-item>
        </el-menu>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { ElMessage } from "element-plus";
import { ArrowRight } from "@element-plus/icons-vue";
import Modal from "./Modal.vue";
import SmallModal from "./SmallModal.vue";
export default {
  name: "ContactListModal",
  props: {
    isVisible: true,
  },
  components: {
    Modal,
    SmallModal,
    ArrowRight,
  },
  setup() {
    const router = useRouter();
    const store = useStore();

    // 选择通知分类
    const handleSelectNotificationCategory = (category) => {
      // 如果点击的是当前选中的分类，则取消筛选（显示全部）
      if (store.state.notificationFilterType === category) {
        store.commit('setNotificationFilterType', null);
      } else {
        store.commit('setNotificationFilterType', category);
      }
    };
    const data = reactive({
      chatMessage: "",
      chatName: "",
      userInfo: store.state.userInfo,
      contactSearch: "",
      createGroupReq: {
        owner_id: "",
        name: "",
        notice: "",
        add_mode: null,
        avatar: "",
      },
      isCreateGroupModalVisible: false,
      isApplyContactModalVisible: false,
      isNewContactModalVisible: false,
      contactUserList: [],
      myGroupList: [],
      myJoinedGroupList: [],
      applyContactReq: {
        owner_id: "",
        contact_id: "",
        message: "",
      },
      ownListReq: {
        owner_id: "",
      },
      newContactList: [],
      applyContent: "",
      uploadRef: null,
      uploadPath: store.state.backendUrl + "/message/uploadAvatar",
      fileList: [],
      cnt: 0,
    });
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
    const handleCreateGroup = async () => {
      try {
        data.createGroupReq.owner_id = data.userInfo.uuid;
        if (data.fileList.length > 0) {
          data.createGroupReq.avatar =
            "/static/avatars/" + data.fileList[0].name;
          console.log(data.createGroupReq.avatar);
          data.uploadRef.submit();
        }
        const response = await axios.post(
          store.state.backendUrl + "/group/createGroup",
          data.createGroupReq
        );
      } catch (error) {
        console.error(error);
      }
    };
    const showCreateGroupModal = () => {
      data.isCreateGroupModalVisible = true;
    };
    const quitCreateGroupModal = () => {
      data.isCreateGroupModalVisible = false;
    };
    const closeCreateGroupModal = () => {
      if (data.createGroupReq.name == "") {
        ElMessage("请输入群聊名称");
        return;
      }
      if (data.createGroupReq.add_mode == null) {
        ElMessage("请选择加群方式");
        return;
      }
      data.isCreateGroupModalVisible = false;
      handleCreateGroup();
    };
    const showApplyContactModal = () => {
      data.isApplyContactModalVisible = true;
    };
    const quitApplyContactModal = () => {
      data.isApplyContactModalVisible = false;
    };
    const closeApplyContactModal = () => {
      if (data.applyContactReq.contact_id == "") {
        ElMessage.error("请输入申请用户/群组id");
        return;
      }
      if (data.applyContactReq.contact_id[0] == "G") {
        handleApplyGroup();
      } else {
        handleApplyContact();
      }
    };

    const showNewContactModal = () => {
      handleNewContactList();
    };

    const quitNewContactModal = () => {
      data.isNewContactModalVisible = false;
      data.newContactList = [];
    };
    const handleApplyGroup = async () => {
      try {
        let req = {
          group_id: data.applyContactReq.contact_id,
        };
        let rsp = await axios.post(
          store.state.backendUrl + "/group/checkGroupAddMode",
          req
        );
        if (rsp.data.code == 200) {
          if (rsp.data.data == 0) {
            // 直接加入
            handleEnterDirectly(data.applyContactReq.contact_id);
            return;
          }
        } else {
          ElMessage.error("申请失败");
          return;
        }
        data.applyContactReq.owner_id = data.userInfo.uuid;
        rsp = await axios.post(
          store.state.backendUrl + "/contact/applyContact",
          data.applyContactReq
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          if (rsp.data.message == "申请成功") {
            data.isApplyContactModalVisible = false;
            ElMessage.success("申请成功");
            return;
          } else {
            ElMessage.error(rsp.data.message);
          }
        } else {
          ElMessage.error("申请失败");
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleApplyContact = async () => {
      try {
        data.applyContactReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/applyContact",
          data.applyContactReq
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          if (rsp.data.message == "申请成功") {
            data.isApplyContactModalVisible = false;
            ElMessage.success("申请成功");
            return;
          }
        }
        ElMessage.error(rsp.data.message);
      } catch (error) {
        console.error(error);
      }
    };
    const handleShowUserList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const getUserListRsp = await axios.post(
          store.state.backendUrl + "/contact/getUserList",
          data.ownListReq
        );
        if (getUserListRsp.data.data) {
          for (let i = 0; i < getUserListRsp.data.data.length; i++) {
            if (!getUserListRsp.data.data[i].avatar.startsWith("http")) {
              getUserListRsp.data.data[i].avatar =
                store.state.backendUrl + getUserListRsp.data.data[i].avatar;
            }
          }
        }
        data.contactUserList = getUserListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideUserList = () => {
      data.contactUserList = [];
    };

    const handleShowMyGroupList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const loadMyGroupRsp = await axios.post(
          store.state.backendUrl + "/group/loadMyGroup",
          data.ownListReq
        );
        if (loadMyGroupRsp.data.data) {
          for (let i = 0; i < loadMyGroupRsp.data.data.length; i++) {
            if (!loadMyGroupRsp.data.data[i].avatar.startsWith("http")) {
              loadMyGroupRsp.data.data[i].avatar =
                store.state.backendUrl + loadMyGroupRsp.data.data[i].avatar;
            }
          }
        }
        data.myGroupList = loadMyGroupRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideMyGroupList = () => {
      data.myGroupList = [];
    };
    const handleShowMyJoinedGroupList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const loadMyJoinedGroupRsp = await axios.post(
          store.state.backendUrl + "/contact/loadMyJoinedGroup",
          data.ownListReq
        );
        if (loadMyJoinedGroupRsp.data.data) {
          for (let i = 0; i < loadMyJoinedGroupRsp.data.data.length; i++) {
            if (!loadMyJoinedGroupRsp.data.data[i].avatar.startsWith("http")) {
              loadMyJoinedGroupRsp.data.data[i].avatar =
                store.state.backendUrl +
                loadMyJoinedGroupRsp.data.data[i].avatar;
            }
          }
        }
        data.myJoinedGroupList = loadMyJoinedGroupRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideMyJoinedGroupList = () => {
      data.myJoinedGroupList = [];
    };

    const handleToChatUser = async (user) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: user.user_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/session/checkOpenSessionAllowed",
          req
        );
        if (rsp.data.code == 200) {
          if (rsp.data.data == true) {
            router.push("/chat/" + user.user_id);
          } else {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          }
        } else {
          ElMessage.error(rsp.data.message);
          console.error(rsp.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.error(error);
      }
    };

    const handleToChatGroup = async (group) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: group.group_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/session/checkOpenSessionAllowed",
          req
        );
        if (rsp.data.code == 200) {
          if (rsp.data.data == true) {
            router.push("/chat/" + group.group_id);
          } else {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          }
        } else {
          if (rsp.data.code == 400) {
            ElMessage.warning(rsp.data.message);
            console.error(rsp.data.message);
          } else {
            ElMessage.error(rsp.data.message);
            console.error(rsp.data.message);
          }
        }
      } catch (error) {
        console.error(error);
      }
    };

    const handleNewContactList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/getNewContactList",
          data.ownListReq
        );
        console.log(rsp);
        data.newContactList = rsp.data.data;
        if (data.newContactList == null) {
          ElMessage.warning("没有新的好友申请");
          return;
        }
        for (let i = 0; i < data.newContactList.length; i++) {
          if (!data.newContactList[i].contact_avatar.startsWith("http")) {
            data.newContactList[i].contact_avatar =
              store.state.backendUrl + data.newContactList[i].contact_avatar;
          }
        }
        data.isNewContactModalVisible = true;
        console.log(rsp);
      } catch (error) {
        console.error(error);
      }
    };

    const handleAgree = async (contactId) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/passContactApply",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          
          // 注意：不在这里建立加密会话
          // 正确的流程是：对方发送第一条 PreKeyMessage 时，接收方会自动建立会话
          // 如果双方都主动建立会话，会导致密钥不匹配
          console.log("✅ [ContactListModal] 已同意好友申请，等待对方发送消息时自动建立加密会话");
          
          data.newContactList = data.newContactList.filter(
            (c) => c.contact_id !== contactId
          );
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleEnterDirectly = async (groupId) => {
      try {
        const req = {
          owner_id: groupId,
          contact_id: data.userInfo.uuid,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/group/enterGroupDirectly",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          data.isApplyContactModalVisible = false;
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleReject = async (contactId) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/refuseContactApply",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
          data.newContactList = data.newContactList.filter(
            (c) => c.contact_id !== contactId
          );
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleBlack = async (contactId) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/blackApply",
          req
        );
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
          data.newContactList = data.newContactList.filter(
            (c) => c.contact_id !== contactId
          );
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.error(error);
      }
    };
    const handleCancelBlack = async (user) => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: user.user_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/cancelBlackContact",
          req
        );
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else if (rsp.data.code == 500) {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        ElMessage.error(error);
        console.error(error);
      }
    };

    return {
      ...toRefs(data),
      router,
      store,
      handleSelectNotificationCategory,
      handleCreateGroup,
      showCreateGroupModal,
      closeCreateGroupModal,
      quitCreateGroupModal,
      showApplyContactModal,
      quitApplyContactModal,
      closeApplyContactModal,
      showNewContactModal,
      quitNewContactModal,
      handleShowUserList,
      handleHideUserList,
      handleShowMyGroupList,
      handleHideMyGroupList,
      handleShowMyJoinedGroupList,
      handleHideMyJoinedGroupList,
      handleToChatUser,
      handleToChatGroup,
      handleNewContactList,
      handleAgree,
      handleReject,
      handleCancelBlack,
      handleBlack,
      handleUploadSuccess,
      beforeFileUpload,
    };
  },
};
</script>

<style>
.contactlist-container {
  overflow-y: auto;
  overflow-x: hidden;
  /* 使用 scrollbar-gutter 预留滚动条空间，避免内容位移 */
  scrollbar-gutter: stable;
}

/* 自定义滚动条样式，使其更细且不占用空间 */
.contactlist-container::-webkit-scrollbar {
  width: 6px;
}

.contactlist-container::-webkit-scrollbar-track {
  background: transparent;
}

.contactlist-container::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
  border: 1px solid transparent;
  background-clip: padding-box;
}

.contactlist-container::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
  background-clip: padding-box;
}

/* Firefox 滚动条样式 */
.contactlist-container {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.15) transparent;
}

.contactlist-header {
  display: flex;
  flex-direction: row;
  margin-top: 12px;
  margin-bottom: 12px;
  padding: 0 8px;
}

.contact-search-input {
  width: 185px;
  height: 36px;
  margin-left: 5px;
  margin-right: 5px;
  border-radius: 8px;
}

.contactlist-header-right {
  width: 40px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.create-group-btn {
  background: #4facfe;
  cursor: pointer;
  border: none;
  height: 36px;
  width: 36px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.create-group-btn:hover {
  background: #3d8bfe;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.create-group-icon {
  width: 18px;
  height: 18px;
  color: #ffffff;
}

.contactlist-body,
.contactlist-user {
  padding: 0 !important;
  margin: 0 !important;
  width: 100%;
  overflow: hidden !important;
}

.el-menu {
  background-color: #f8f9fa;
  width: 100% !important;
  border: none;
  padding: 0 !important;
  margin: 0 !important;
  overflow: hidden !important;
  overflow-x: hidden !important;
  overflow-y: hidden !important;
}

.el-sub-menu {
  width: 100% !important;
  overflow: hidden !important;
}

.el-menu-item {
  background-color: #ffffff;
  height: 48px;
  border-radius: 0;
  margin: 0 !important;
  transition: all 0.2s ease;
  padding-left: 12px !important;
  padding-right: 0 !important;
  border-bottom: 1px solid #f0f0f0;
  width: 100% !important;
  box-sizing: border-box;
}

.el-menu-item:last-child {
  border-bottom: none;
}

.el-menu-item:hover {
  background-color: #f3f4f6;
}

.el-menu-item.is-active {
  background-color: #4facfe;
  color: #ffffff;
}

.contactlist-user-title {
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 600;
  color: #1f2937;
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
  /*background-color:aqua;*/
}

.modal-body {
  height: 55%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.creatgroup-modal-body {
  height: 75%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.newcontact-modal-body {
  height: 70%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.newcontact-modal-footer {
  height: 10%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-footer {
  height: 25%;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.creategroup-modal-footer {
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

.contactlist-avatar {
  width: 30px;
  height: 30px;
  margin-right: 20px;
}

.newcontact-list {
  width: 280px;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.newcontact-item {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 40px;
}

.action-btn {
  background: #4facfe;
  border: none;
  color: #ffffff;
  border: none;
  cursor: pointer;
  justify-content: center;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.contactlist-user-menu-item {
  justify-content: center;
  align-items: center;
}

.contactlist-user-item {
  width: 221px;
  height: 45px;
  display: flex;
  align-items: center;
  color: rgba(43, 42, 42, 0.893);
}

.contactlist-user-avatar {
  width: 30px;
  height: 30px;
  margin-left: 20px;
  margin-right: 20px;
}

.notification-category-list {
  margin-top: 20px;
  padding: 0;
}

.notification-category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  color: rgba(43, 42, 42, 0.893);
  transition: background-color 0.2s;
}

.notification-category-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.notification-category-item.active {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.notification-category-text {
  font-size: 14px;
  font-weight: 500;
}

.notification-category-arrow {
  font-size: 14px;
  color: rgba(43, 42, 42, 0.5);
}

.notification-category-item.active .notification-category-arrow {
  color: #409eff;
}

.notification-category-divider {
  height: 1px;
  background-color: rgba(0, 0, 0, 0.1);
  margin: 0 20px;
}

.notification-contact-divider {
  height: 1px;
  background-color: rgba(0, 0, 0, 0.15);
  margin: 15px 20px;
}
</style>