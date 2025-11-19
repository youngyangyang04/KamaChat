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
        <el-container class="chat-container">
          <el-header>
            <h2 class="chat-name"></h2>
          </el-header>
          <el-main class="main-container">
            <div class="chat-screen"></div>
            <div class="tool-bar">
              <div class="tool-bar-left">
                <el-tooltip
                  effect="customized"
                  content="表情包"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733502796507"
                      class="sticker-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="1555"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M504.32 31.872a472.256 472.256 0 1 1 0 944.512 472.256 472.256 0 0 1 0-944.512z m0 63.36a408.96 408.96 0 1 0 0 817.856 408.96 408.96 0 0 0 0-817.92z m228.864 487.808v0.192a217.856 217.856 0 1 1-435.712 0V583.04h435.712zM370.496 321.536a73.024 73.024 0 1 1 0 146.048 73.024 73.024 0 0 1 0-146.048z m289.664 0a73.024 73.024 0 1 1 0 146.048 73.024 73.024 0 0 1 0-146.048z"
                        fill="#2c2c2c"
                        p-id="1556"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="文件上传"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503065264"
                      class="upload-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="2430"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M543.7 157v534c0 16.6-13.4 30-30 30s-30-13.4-30-30V157c0-16.6 13.4-30 30-30 16.5 0 30 13.4 30 30z"
                        fill=""
                        p-id="2431"
                      ></path>
                      <path
                        d="M323.1 331c11.8 11.8 30.7 11.8 42.5 0l119.9-119.9c15.6-15.6 40.9-15.6 56.6 0L662 331c11.7 11.7 30.7 11.7 42.4 0s11.7-30.7 0-42.4L541.7 126.1c-15.6-15.6-41-15.6-56.6 0L323 288.6c-11.6 11.8-11.6 30.7 0.1 42.4zM853.7 913h-680c-33.1 0-60-26.9-60-60V583.7c0-16.4 12.8-30.2 29.2-30.7 16.9-0.4 30.8 13.2 30.8 30v240c0 16.6 13.4 30 30 30h620c16.6 0 30-13.4 30-30V583.7c0-16.4 12.8-30.2 29.2-30.7 16.9-0.4 30.8 13.2 30.8 30v270c0 33.1-26.9 60-60 60z"
                        fill=""
                        p-id="2432"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="聊天记录"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733504061769"
                      class="record-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="5492"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M695.04 194.32H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h596.96c18.32 0 33.16 14.85 33.16 33.16 0 18.31-14.84 33.16-33.16 33.16zM298.97 393.3H96.19c-17.27 0-31.27-14-31.27-31.27v-3.79c0-17.27 14-31.27 31.27-31.27h202.78c17.27 0 31.27 14 31.27 31.27v3.79c-0.01 17.28-14.01 31.27-31.27 31.27zM230.74 592.29H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM230.74 791.28H98.08c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16h132.66c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16zM728.2 691.78H595.55c-18.32 0-33.16-14.85-33.16-33.16 0-18.32 14.85-33.16 33.16-33.16H728.2c18.32 0 33.16 14.85 33.16 33.16 0.01 18.31-14.84 33.16-33.16 33.16z"
                        p-id="5493"
                      ></path>
                      <path
                        d="M562.38 658.62V525.96c0-18.32 14.85-33.16 33.16-33.16 18.32 0 33.16 14.85 33.16 33.16v132.66c0 18.32-14.85 33.16-33.16 33.16-18.31 0-33.16-14.85-33.16-33.16z"
                        p-id="5494"
                      ></path>
                      <path
                        d="M960.35 625.45c0 183.16-148.48 331.64-331.64 331.64S297.07 808.62 297.07 625.45s148.48-331.64 331.64-331.64 331.64 148.48 331.64 331.64zM628.71 360.14c-146.53 0-265.31 118.79-265.31 265.31s118.79 265.31 265.31 265.31 265.31-118.79 265.31-265.31-118.78-265.31-265.31-265.31z"
                        p-id="5495"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
                <el-tooltip
                  effect="customized"
                  content="全文复制"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503137487"
                      class="copy-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="3442"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M394.666667 106.666667h448a74.666667 74.666667 0 0 1 74.666666 74.666666v448a74.666667 74.666667 0 0 1-74.666666 74.666667H394.666667a74.666667 74.666667 0 0 1-74.666667-74.666667V181.333333a74.666667 74.666667 0 0 1 74.666667-74.666666z m0 64a10.666667 10.666667 0 0 0-10.666667 10.666666v448a10.666667 10.666667 0 0 0 10.666667 10.666667h448a10.666667 10.666667 0 0 0 10.666666-10.666667V181.333333a10.666667 10.666667 0 0 0-10.666666-10.666666H394.666667z m245.333333 597.333333a32 32 0 0 1 64 0v74.666667a74.666667 74.666667 0 0 1-74.666667 74.666666H181.333333a74.666667 74.666667 0 0 1-74.666666-74.666666V394.666667a74.666667 74.666667 0 0 1 74.666666-74.666667h74.666667a32 32 0 0 1 0 64h-74.666667a10.666667 10.666667 0 0 0-10.666666 10.666667v448a10.666667 10.666667 0 0 0 10.666666 10.666666h448a10.666667 10.666667 0 0 0 10.666667-10.666666v-74.666667z"
                        fill="#000000"
                        p-id="3443"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
              </div>
              <div class="tool-bar-right">
                <el-tooltip
                  effect="customized"
                  content="音视频通话"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button">
                    <svg
                      t="1733503700535"
                      class="av-icon"
                      viewBox="0 0 1024 1024"
                      version="1.1"
                      xmlns="http://www.w3.org/2000/svg"
                      p-id="4492"
                      width="128"
                      height="128"
                    >
                      <path
                        d="M790.207709 1023.317561c-100.48917-0.05687-302.832389-33.89448-528.321671-260.00933C-57.722981 442.903032-9.212929 154.458736 25.02277 119.995557L114.194824 30.709763c19.506387-19.563257 47.372654-30.709763 76.319449-30.709763 28.662446 0 56.073753 10.975897 75.23892 30.141064l3.980896 4.606465 131.881373 176.865489c35.145618 52.377208 33.32578 108.564701-4.720205 146.781295l-39.012773 39.069643c11.942686 71.087415 42.31123 113.398645 87.181606 158.439632l5.686993 5.686993c51.865378 52.092858 96.678885 97.076974 174.021993 103.730756l38.899033-38.955903a99.522381 99.522381 0 0 1 71.883595-30.368544c24.169721 0 49.419971 8.41675 73.020993 24.340331l178.002888 133.303121c21.212485 14.558703 34.918138 38.728424 37.477285 66.253471a113.853604 113.853604 0 0 1-33.26891 89.513274l-89.058314 89.285793c-22.179274 22.236144-85.304898 24.624681-111.465068 24.624681h-0.056869zM190.628013 88.091525a19.278907 19.278907 0 0 0-13.421304 5.402644L94.290348 176.63801c-4.549595 22.861713-44.984116 247.554815 230.607575 523.885815 202.684439 203.196268 377.50261 233.507942 463.774297 233.507942 30.652893 0 50.898589-3.753416 58.121071-5.402643l80.982784-82.006443a26.160169 26.160169 0 0 0 7.67744-18.539598l-178.457847-135.293568c-4.151505-2.786627-12.568255-7.677441-20.302566-7.677441a13.478174 13.478174 0 0 0-10.009108 3.980895l-65.969121 66.196601-18.653338-0.17061c-125.227591-1.080529-193.812729-69.950017-254.322337-130.743974l-5.686993-5.630123c-52.490947-52.661557-102.763968-117.20893-115.445963-232.199934l-2.388537-21.155614L333.826502 295.609908c8.41675-8.41675 1.990448-22.349883-4.833944-32.586471L200.750861 91.105631a17.515939 17.515939 0 0 0-10.122848-3.014106z m350.603132 312.159058c-44.131067 0-79.959125-34.235699-79.959125-76.319449V170.609797c0-42.08375 35.828057-76.376319 79.959125-76.376319h292.311452c37.136066 0 68.812618 77.968677 77.627457 111.863156 8.1324-4.606465 14.103743-8.07553 15.923581-9.269799 8.75797-5.743863 18.937687-62.670665 29.458625-62.670665a53.457736 53.457736 0 0 1 25.36399 6.426303 56.130623 56.130623 0 0 1 29.003666 49.87493v121.303566c0 21.496834-11.373986 40.775741-29.572365 50.443629a52.547817 52.547817 0 0 1-24.681551 6.141953c-10.577807 0-21.041875-56.983672-29.970454-62.955015-2.331667-1.421748-8.814839-5.118294-17.686549-10.179718-11.089637 30.368544-41.515051 105.038765-75.40953 105.038765H541.231145z m283.326003-88.944574V183.178052H550.273464v128.127957h274.283684z"
                        fill="#666666"
                        p-id="4493"
                      ></path>
                    </svg>
                  </button>
                </el-tooltip>
              </div>
            </div>
          </el-main>
          <el-footer>
            <div class="chat-input">
              <el-input
                v-model="chatMessage"
                type="textarea"
                show-word-limit
                maxlength="500"
                :autosize="{ minRows: 7.9, maxRows: 7 }"
                placeholder="请输入内容"
              />
            </div>
            <div class="chat-send">
              <el-button class="send-btn">发送</el-button>
            </div>
          </el-footer>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import Modal from "@/components/Modal.vue";
import SmallModal from "@/components/SmallModal.vue";
import ContactListModal from "@/components/ContactListModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import axios from "@/utils/axios";
import { ElMessage } from "element-plus";
export default {
  name: "ContactList",
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
    });

    const handleCreateGroup = async () => {
      try {
        data.createGroupReq.owner_id = data.userInfo.uuid;
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
      if (data.applyContactReq.contact_id[0] == 'G') {
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
        }
        let rsp = await axios.post(store.state.backendUrl + "/group/checkGroupAddMode", req);
        if (rsp.data.code == 200) {
          if (rsp.data.data == 0) { // 直接加入
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
          owner_id:  groupId,
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
    }
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
    }
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
    };
  },
};
</script>

<style scoped>
.contactlist-header {
  display: flex;
  flex-direction: row;
  margin-top: 10px;
  margin-bottom: 10px;
}

.contact-search-input {
  width: 185px;
  height: 30px;
  margin-left: 5px;
  margin-right: 5px;
}

.contactlist-header-right {
  width: 40px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.create-group-btn {
  background-color: rgb(252, 210.9, 210.9);
  cursor: pointer;
  border: none;
  height: 100%;
  width: 30px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 10px;
}

.create-group-icon {
  width: 15px;
  height: 15px;
}

.el-menu {
  background-color: rgb(252, 210.9, 210.9);
  width: 101%;
}

.el-menu-item {
  background-color: rgb(255, 255, 255);
  height: 45px;
}

.contactlist-user-title {
  font-family: Arial, Helvetica, sans-serif;
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
  background-color: rgb(252, 210.9, 210.9);
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
</style>