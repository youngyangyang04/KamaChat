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
          <div class="sessionlist-container">
            <div class="sessionlist-header">
              <el-input
                v-model="contactSearch"
                class="contact-search-input"
                placeholder="搜索会话"
                size="small"
                suffix-icon="Search"
              />
            </div>
            <div class="contactlist-body">
              <div class="contactlist-user">
                <el-menu router>
                  <el-menu-item
                    v-for="session in allSessionList"
                    :key="session.id"
                    @click="handleToSession(session)"
                  >
                    <img :src="session.avatar" class="sessionlist-avatar" />
                    {{ session.name }}
                  </el-menu-item>
                </el-menu>
              </div>
            </div>
          </div>
        </el-aside>
        <el-container class="chat-container">
          <el-header>
            <div class="chat-title" v-if="contactInfo.contact_avatar">
              <img
                :src="contactInfo.contact_avatar"
                style="width: 40px; height: 40px; margin-right: 10px"
              />
              <h2 class="chat-name">{{ contactInfo.contact_name }}</h2>
            </div>
            <div class="chat-title-right">
              <Modal :isVisible="isUserContactInfoModalVisible">
                <template v-slot:header>
                  <div class="userinfo-modal-quit-btn-container">
                    <button
                      class="userinfo-modal-quit-btn"
                      @click="quitUserContactInfoModal"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </div>
                  <div class="userinfo-modal-header-title">
                    <h3>个人主页</h3>
                  </div>
                </template>
                <template v-slot:body>
                  <el-descriptions
                    direction="vertical"
                    border
                    class="modal-list"
                    size="small"
                  >
                    <el-descriptions-item
                      :rowspan="2"
                      :width="120"
                      label="头像"
                      align="center"
                    >
                      <el-image
                        style="width: 100px; height: 100px"
                        :src="contactInfo.contact_avatar"
                      />
                    </el-descriptions-item>
                    <el-descriptions-item label="Id" :width="140">{{
                      contactInfo.contact_id
                    }}</el-descriptions-item>
                    <el-descriptions-item label="性别">{{
                      contactInfo.contact_gender == 0 ? "男" : "女"
                    }}</el-descriptions-item>
                    <el-descriptions-item label="电话号码">{{
                      contactInfo.contact_phone
                    }}</el-descriptions-item>
                    <el-descriptions-item label="昵称">{{
                      contactInfo.contact_name
                    }}</el-descriptions-item>

                    <el-descriptions-item label="邮箱" :span="2">
                      <div style="height: 30px">
                        {{ contactInfo.contact_email }}
                      </div></el-descriptions-item
                    >

                    <el-descriptions-item label="生日" :span="1" :width="140"
                      >{{ contactInfo.contact_birthday }}
                    </el-descriptions-item>
                    <el-descriptions-item label="个性签名">
                      <div style="height: 70px">
                        {{ contactInfo.contact_signature }}
                      </div>
                    </el-descriptions-item>
                  </el-descriptions>
                </template>
              </Modal>
              <Modal :isVisible="isGroupContactInfoModalVisible">
                <template v-slot:header>
                  <div class="groupcontactinfo-modal-quit-btn-container">
                    <button
                      class="groupcontactinfo-modal-quit-btn"
                      @click="quitGroupContactInfoModal"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </div>
                  <div class="groupcontactinfo-modal-header-title">
                    <h3>群聊主页</h3>
                  </div>
                </template>
                <template v-slot:body>
                  <el-descriptions
                    direction="vertical"
                    border
                    class="modal-list"
                    size="small"
                  >
                    <el-descriptions-item
                      :rowspan="2"
                      :width="120"
                      label="头像"
                      align="center"
                    >
                      <el-image
                        style="width: 100px; height: 100px"
                        :src="contactInfo.contact_avatar"
                      />
                    </el-descriptions-item>
                    <el-descriptions-item label="Id" :width="140">{{
                      contactInfo.contact_id
                    }}</el-descriptions-item>
                    <el-descriptions-item label="群人数">{{
                      contactInfo.contact_member_cnt
                    }}</el-descriptions-item>
                    <el-descriptions-item label="群主id">{{
                      contactInfo.contact_owner_id
                    }}</el-descriptions-item>
                    <el-descriptions-item label="入群方式" :width="140"
                      >{{
                        contactInfo.contact_add_mode == 0
                          ? "直接加入"
                          : "群主审核"
                      }}
                    </el-descriptions-item>
                    <el-descriptions-item label="群名称" :span="3">{{
                      contactInfo.contact_name
                    }}</el-descriptions-item>
                    <el-descriptions-item label="群公告" :span="3">
                      <div style="height: 70px">
                        {{ contactInfo.contact_notice }}
                      </div>
                    </el-descriptions-item>
                  </el-descriptions>
                </template>
              </Modal>
              <el-dropdown placement="bottom" trigger="click">
                <button class="setting-btn">
                  <el-icon><MoreFilled /></el-icon>
                </button>
                <template #dropdown>
                  <el-dropdown-menu v-if="contactInfo.contact_id[0] === 'U'">
                    <el-dropdown-item @click="showUserContactInfoModal">
                      个人信息
                    </el-dropdown-item>

                    <el-dropdown-item @click="preToDeleteSession"
                      >删除该会话</el-dropdown-item
                    >
                    <el-dropdown-item @click="preToDeleteContact"
                      >删除联系人</el-dropdown-item
                    >
                    <el-dropdown-item @click="preToBlackContact"
                      >拉黑联系人</el-dropdown-item
                    >
                  </el-dropdown-menu>
                  <el-dropdown-menu
                    v-else-if="contactInfo.contact_id[0] === 'G'"
                  >
                    <el-dropdown-item @click="showGroupContactInfoModal"
                      >群聊信息</el-dropdown-item
                    >
                    <el-dropdown-item
                      v-if="contactInfo.contact_owner_id == userInfo.uuid"
                      @click="showUpdateGroupInfoModal"
                    >
                      修改群聊信息
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="contactInfo.contact_owner_id == userInfo.uuid"
                      @click="showRemoveGroupMemberModal"
                    >
                      移除群组人员
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="contactInfo.contact_owner_id == userInfo.uuid"
                      @click="showAddGroupModal"
                      >加群申请</el-dropdown-item
                    >
                    <el-dropdown-item @click="preToDeleteSession"
                      >删除该会话</el-dropdown-item
                    >
                    <el-dropdown-item
                      v-if="contactInfo.contact_owner_id == userInfo.uuid"
                      @click="handleDismissGroup"
                      >解散群聊</el-dropdown-item
                    >
                    <el-dropdown-item
                      v-if="contactInfo.contact_owner_id != userInfo.uuid"
                      @click="handleLeaveGroup"
                      >退出群聊</el-dropdown-item
                    >
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <Modal :isVisible="isUpdateGroupInfoModalVisible">
                <template v-slot:header>
                  <div class="updategroupinfo-modal-quit-btn-container">
                    <button
                      class="updategroupinfo-modal-quit-btn"
                      @click="quitUpdateGroupInfoModal"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </div>
                  <div class="updategroupinfo-modal-header-title">
                    <h3>修改群聊信息</h3>
                  </div>
                </template>
                <template v-slot:body>
                  <el-scrollbar
                    max-height="255px"
                    style="
                      width: 400px;
                      height: 255px;
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      margin-top: 20px;
                    "
                  >
                    <div class="modal-body">
                      <el-form
                        ref="formRef"
                        :model="updateGroupInfo"
                        label-width="80px"
                      >
                        <el-form-item
                          prop="name"
                          label="群名称"
                          :rules="[
                            {
                              min: 3,
                              max: 10,
                              message: '群名称长度在 3 到 10 个字符',
                              trigger: 'blur',
                            },
                          ]"
                        >
                          <el-input
                            v-model="updateGroupInfo.name"
                            placeholder="选填"
                          />
                        </el-form-item>
                        <el-form-item prop="add_mode" label="入群方式">
                          <el-radio-group v-model="updateGroupInfo.add_mode">
                            <el-radio :value="0">直接加入</el-radio>
                            <el-radio :value="1">群主审核</el-radio>
                          </el-radio-group>
                        </el-form-item>
                        <el-form-item prop="notice" label="群公告">
                          <el-input
                            v-model="updateGroupInfo.notice"
                            type="textarea"
                            show-word-limit
                            maxlength="500"
                            :autosize="{ minRows: 3, maxRows: 3 }"
                            placeholder="选填"
                          />
                        </el-form-item>
                        <el-form-item prop="avatar" label="群头像">
                          <el-upload
                            v-model:file-list="avatarList"
                            ref="uploadAvatarRef"
                            :auto-upload="false"
                            :action="uploadAvatarPath"
                            :on-success="handleAvatarUploadSuccess"
                            :before-upload="beforeAvatarUpload"
                          >
                            <template #trigger>
                              <el-button
                                style="background: #4facfe; border: none; color: #ffffff;"
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
                  <div class="updategroupinfo-modal-footer">
                    <el-button
                      style="background: #4facfe; border: none; color: #ffffff;"
                      @click="closeUpdateGroupInfoModal"
                    >
                      完成
                    </el-button>
                  </div>
                </template>
              </Modal>
              <Modal :isVisible="isRemoveGroupMemberModalVisible">
                <template v-slot:header>
                  <div class="removegroupmember-modal-quit-btn-container">
                    <button
                      class="removegroupmember-modal-quit-btn"
                      @click="quitRemoveGroupMemberModal"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </div>
                  <div class="removegroupmember-modal-header-title">
                    <h3>移除群组人员</h3>
                  </div>
                </template>
                <template v-slot:body>
                  <span
                    style="
                      font-size: 14px;
                      font-weight: bold;
                      font-family: Arial, Helvetica, sans-serif;
                      color: rgb(57, 57, 57);
                      width: 270px;
                      display: flex;
                      justify-content: left;
                      margin-bottom: 5px;
                    "
                    >群组成员：</span
                  >
                  <el-scrollbar
                    max-height="400px"
                    style="height: 300px; width: 350px"
                  >
                    <div class="modal-body">
                      <ul
                        style="list-style-type: none"
                        class="removegroupmembers-list"
                      >
                        <li
                          v-for="groupMember in groupMemberList"
                          :key="groupMember.user_id"
                          class="removegroupmembers-item"
                        >
                          <div style="display: flex; align-items: center">
                            <el-image
                              :src="groupMember.avatar"
                              class="removegroupmembers-item-avatar"
                            />
                            <span class="removegroupmembers-item-name">{{
                              groupMember.nickname
                            }}</span>
                          </div>
                          <input
                            type="checkbox"
                            :value="groupMember.user_id"
                            v-model="selectedGroupMembers"
                            @change="handleCheckboxChange"
                          />
                        </li>
                      </ul>
                    </div>
                  </el-scrollbar>
                </template>
                <template v-slot:footer>
                  <div
                    style="
                      height: 50px;
                      width: 300px;
                      display: flex;
                      justify-content: right;
                    "
                  >
                    <el-button
                      class="removegroupmembers-button"
                      @click="handleRemoveGroupMembers"
                      >移除所选人员</el-button
                    >
                  </div>
                </template>
              </Modal>
              <SmallModal :isVisible="isAddGroupModalVisible">
                <template v-slot:header>
                  <div class="modal-header">
                    <div class="modal-quit-btn-container">
                      <button class="modal-quit-btn" @click="quitAddGroupModal">
                        <el-icon><Close /></el-icon>
                      </button>
                    </div>
                    <div class="modal-header-title">
                      <h3>加群申请</h3>
                    </div>
                  </div>
                </template>
                <template v-slot:body>
                  <div class="addGroup-modal-body">
                    <el-scrollbar max-height="400px">
                      <ul class="addGroup-list" style="list-style-type: none">
                        <li
                          v-for="addGroup in addGroupList"
                          :key="addGroup.contact_id"
                          class="addGroup-item"
                        >
                          <div
                            style="
                              display: flex;
                              align-items: center;
                              justify-content: center;
                            "
                          >
                            <img
                              :src="addGroup.contact_avatar"
                              style="
                                width: 30px;
                                height: 30px;
                                margin-right: 10px;
                              "
                            />

                            <el-tooltip
                              effect="customized"
                              :content="addGroup.message"
                              placement="top"
                              hide-after="0"
                              enterable="false"
                            >
                              <div style="color: black">
                                {{ addGroup.contact_name }}
                              </div>
                            </el-tooltip>
                          </div>
                          <el-dropdown placement="right" trigger="click">
                            <el-button class="action-btn"> 去处理 </el-button>
                            <template #dropdown>
                              <el-dropdown-menu>
                                <el-dropdown-item
                                  @click="handleAgree(addGroup.contact_id)"
                                  >同意</el-dropdown-item
                                >
                                <el-dropdown-item
                                  @click="handleReject(addGroup.contact_id)"
                                >
                                  拒绝
                                </el-dropdown-item>
                              </el-dropdown-menu>
                            </template>
                          </el-dropdown>
                        </li>
                      </ul>
                    </el-scrollbar>
                  </div>
                </template>
              </SmallModal>
            </div>
          </el-header>
          <el-main class="main-container">
            <el-scrollbar
              max-height="332.5px"
              style="height: 332.5px"
              ref="scrollbarRef"
            >
              <div ref="innerRef">
                <div
                  v-for="(messageItem, index) in messageList"
                  :key="index"
                  class="message-item"
                >
                  <div
                    v-if="
                      messageItem.send_id != userInfo.uuid &&
                      messageItem.type == 0
                    "
                    class="left-message"
                  >
                    <div class="left-message-left">
                      <el-image
                        :src="messageItem.send_avatar"
                        style="
                          width: 40px;
                          height: 40px;
                          margin-left: 10px;
                          margin-right: 10px;
                          margin-top: 10px;
                        "
                      >
                      </el-image>
                    </div>

                    <div class="left-message-right">
                      <div class="left-message-right-top">
                        <div class="left-message-contactname">
                          {{ messageItem.send_name }}
                        </div>
                        <div class="left-message-time">
                          {{ messageItem.created_at }}
                        </div>
                      </div>

                      <div class="left-message-content" :class="{ 'unread-message': messageItem.is_unread }">
                        <span v-if="messageItem.is_encrypted" style="color: #67c23a; font-size: 12px; margin-right: 4px;">🔒</span>
                        <span v-if="messageItem.is_unread" class="unread-indicator"></span>
                        {{ messageItem.content }}
                      </div>
                    </div>
                  </div>
                  <div
                    v-if="
                      messageItem.send_id != userInfo.uuid &&
                      messageItem.type == 2
                    "
                    class="left-message"
                  >
                    <div class="left-message-left">
                      <el-image
                        :src="messageItem.send_avatar"
                        style="
                          width: 40px;
                          height: 40px;
                          margin-left: 10px;
                          margin-right: 10px;
                          margin-top: 10px;
                        "
                      >
                      </el-image>
                    </div>

                    <div class="left-message-right">
                      <div class="left-message-right-top">
                        <div class="left-message-contactname">
                          {{ messageItem.send_name }}
                        </div>
                        <div class="left-message-time">
                          {{ messageItem.created_at }}
                        </div>
                      </div>

                      <div class="left-message-file-container">
                        <div style="display: flex; flex-direction: row">
                          <div class="left-message-file-name">
                            {{ messageItem.file_name }}
                          </div>
                          <div class="left-message-file-size">
                            {{ messageItem.file_size }}
                          </div>
                        </div>

                        <div class="left-message-file-download">
                          <el-button
                            style="
                              background: #4facfe;
                              border: none;
                              color: #ffffff;
                              margin-top: 20px;
                            "
                            size="small"
                            @click="downloadFile(messageItem.file_name)"
                          >
                            下载
                          </el-button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 左侧加密图片消息 (type == 4) -->
                  <div
                    v-if="
                      messageItem.send_id != userInfo.uuid &&
                      messageItem.type == 4
                    "
                    class="left-message"
                  >
                    <div class="left-message-left">
                      <el-image
                        :src="messageItem.send_avatar"
                        style="width: 40px; height: 40px; margin-left: 10px; margin-right: 10px; margin-top: 10px;"
                      />
                    </div>
                    <div class="left-message-right">
                      <div class="left-message-right-top">
                        <div class="left-message-contactname">{{ messageItem.send_name }}</div>
                        <div class="left-message-time">{{ messageItem.created_at }}</div>
                      </div>
                      <FileMessage 
                        :file-info="parseFileInfo(messageItem)" 
                        :is-encrypted="true" 
                      />
                    </div>
                  </div>

                  <!-- 左侧加密文件消息 (type == 5) -->
                  <div
                    v-if="
                      messageItem.send_id != userInfo.uuid &&
                      messageItem.type == 5
                    "
                    class="left-message"
                  >
                    <div class="left-message-left">
                      <el-image
                        :src="messageItem.send_avatar"
                        style="width: 40px; height: 40px; margin-left: 10px; margin-right: 10px; margin-top: 10px;"
                      />
                    </div>
                    <div class="left-message-right">
                      <div class="left-message-right-top">
                        <div class="left-message-contactname">{{ messageItem.send_name }}</div>
                        <div class="left-message-time">{{ messageItem.created_at }}</div>
                      </div>
                      <FileMessage 
                        :file-info="parseFileInfo(messageItem)" 
                        :is-encrypted="true" 
                      />
                    </div>
                  </div>

                  <div
                    style="
                      width: 100%;
                      height: 100%;
                      display: flex;
                      flex-direction: row-reverse;
                    "
                  >
                    <div
                      v-if="
                        messageItem.send_id == userInfo.uuid &&
                        messageItem.type == 0
                      "
                      class="right-message"
                    >
                      <div class="right-message-right">
                        <el-image
                          :src="userInfo.avatar"
                          style="
                            width: 40px;
                            height: 40px;
                            margin-left: 10px;
                            margin-right: 10px;
                            margin-top: 10px;
                          "
                        >
                        </el-image>
                      </div>

                      <div class="right-message-left">
                        <div class="right-message-left-top">
                          <div class="right-message-contactname">
                            {{ userInfo.nickname }}
                          </div>
                          <div class="right-message-time">
                            {{ messageItem.created_at }}
                          </div>
                        </div>
                        <div style="display: flex; flex-direction: row-reverse">
                          <div class="right-message-content">
                            <span v-if="messageItem.is_encrypted" style="color: #67c23a; font-size: 12px; margin-right: 4px;">🔒</span>
                            {{ messageItem.content }}
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      v-if="
                        messageItem.send_id == userInfo.uuid &&
                        messageItem.type == 2
                      "
                      class="right-message"
                    >
                      <div class="right-message-right">
                        <el-image
                          :src="userInfo.avatar"
                          style="
                            width: 40px;
                            height: 40px;
                            margin-left: 10px;
                            margin-right: 10px;
                            margin-top: 10px;
                          "
                        >
                        </el-image>
                      </div>

                      <div class="right-message-left">
                        <div class="right-message-left-top">
                          <div class="right-message-contactname">
                            {{ userInfo.nickname }}
                          </div>
                          <div class="right-message-time">
                            {{ messageItem.created_at }}
                          </div>
                        </div>
                        <div style="display: flex; flex-direction: row-reverse">
                          <div class="right-message-file-container">
                            <div style="display: flex; flex-direction: row">
                              <div class="right-message-file-name">
                                {{ messageItem.file_name }}
                              </div>
                              <div class="right-message-file-size">
                                {{ messageItem.file_size }}
                              </div>
                            </div>

                            <div class="right-message-file-download">
                              已发送
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>

                    <!-- 右侧加密图片消息 (type == 4) -->
                    <div
                      v-if="
                        messageItem.send_id == userInfo.uuid &&
                        messageItem.type == 4
                      "
                      class="right-message"
                    >
                      <div class="right-message-right">
                        <el-image
                          :src="userInfo.avatar"
                          style="width: 40px; height: 40px; margin-left: 10px; margin-right: 10px; margin-top: 10px;"
                        />
                      </div>
                      <div class="right-message-left">
                        <div class="right-message-left-top">
                          <div class="right-message-contactname">{{ userInfo.nickname }}</div>
                          <div class="right-message-time">{{ messageItem.created_at }}</div>
                        </div>
                        <div style="display: flex; flex-direction: row-reverse">
                          <FileMessage 
                            :file-info="parseFileInfo(messageItem)" 
                            :is-encrypted="true" 
                          />
                        </div>
                      </div>
                    </div>

                    <!-- 右侧加密文件消息 (type == 5) -->
                    <div
                      v-if="
                        messageItem.send_id == userInfo.uuid &&
                        messageItem.type == 5
                      "
                      class="right-message"
                    >
                      <div class="right-message-right">
                        <el-image
                          :src="userInfo.avatar"
                          style="width: 40px; height: 40px; margin-left: 10px; margin-right: 10px; margin-top: 10px;"
                        />
                      </div>
                      <div class="right-message-left">
                        <div class="right-message-left-top">
                          <div class="right-message-contactname">{{ userInfo.nickname }}</div>
                          <div class="right-message-time">{{ messageItem.created_at }}</div>
                        </div>
                        <div style="display: flex; flex-direction: row-reverse">
                          <FileMessage 
                            :file-info="parseFileInfo(messageItem)" 
                            :is-encrypted="true" 
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-scrollbar>
            <div class="tool-bar">
              <div class="tool-bar-left">
                <FileUpload 
                  v-if="contactInfo.contact_id && contactInfo.contact_id.startsWith('U')"
                  @upload-complete="handleFileUploadComplete"
                  @upload-error="handleFileUploadError"
                />
              </div>
              <div class="tool-bar-right">
                <el-tooltip
                  effect="customized"
                  content="音视频通话"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button class="image-button" @click="showAVContainerModal">
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
                <div
                  class="video-modal-overlay"
                  v-show="isAVContainerModalVisible"
                >
                  <div class="video-modal-content">
                    <div class="video-modal-header">
                      <h2>聊天室</h2>
                    </div>
                    <div class="video-modal-body">
                      <video autoplay playsinline class="local-video"></video>
                      <video autoplay playsinline class="remote-video"></video>
                    </div>
                    <div class="video-modal-footer">
                      <el-button
                        class="video-modal-footer-btn"
                        @click="startCall(true)"
                        >发起通话</el-button
                      >
                      <el-button
                        class="video-modal-footer-btn"
                        @click="startCall(false)"
                        >接收通话</el-button
                      >
                      <el-button
                        class="video-modal-footer-btn"
                        @click="rejectCall"
                        >拒绝通话</el-button
                      >
                      <el-button
                        class="video-modal-footer-btn"
                        @click="sendEndCall"
                        >挂断通话</el-button
                      >
                      <el-button
                        class="video-modal-footer-btn"
                        @click="closeAVContainerModal"
                        >退出聊天室</el-button
                      >
                    </div>
                  </div>
                </div>
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
                placeholder="请输入内容（Enter发送，Shift+Enter换行）"
                @keydown="handleInputKeydown"
              />
            </div>
            <div class="chat-send">
              <el-button class="send-btn" @click="sendMessage">发送</el-button>
            </div>
          </el-footer>
        </el-container>
      </el-container>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted, onUnmounted, ref, nextTick } from "vue";
import { useRouter, onBeforeRouteUpdate } from "vue-router";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import Modal from "@/components/Modal.vue";
import SmallModal from "@/components/SmallModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import FileUpload from "@/components/FileUpload.vue";
import FileMessage from "@/components/FileMessage.vue";
import { ElMessage, ElMessageBox, ElScrollbar } from "element-plus";
import { ElNotification } from "element-plus";
import {
  hasSession,
  createSession,
  encryptAndSendMessage,
  receiveAndDecryptMessage,
} from "@/crypto";
import { decryptMessageList, decryptMessage } from "@/utils/messageDecryptor";
import eventBus from "@/utils/eventBus";
export default {
  name: "ContactChat",
  components: {
    Modal,
    SmallModal,
    NavigationModal,
    FileUpload,
    FileMessage,
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
      isUserContactInfoModalVisible: false,
      isGroupContactInfoModalVisible: false,
      isAddGroupModalVisible: false,
      isUpdateGroupInfoModalVisible: false,
      isRemoveGroupMemberModalVisible: false,
      getUserListReq: {
        owner_id: "",
      },
      contactUserList: [],
      loadMyGroupReq: {
        owner_id: "",
      },
      myGroupList: [],
      loadMyJoinedGroupReq: {
        owner_id: "",
      },
      myJoinedGroupList: [],
      getContactInfoReq: {
        contact_id: "",
      },
      contactInfo: {
        contact_id: "",
        contact_name: "",
        contact_avatar: "",
        contact_phone: "",
        contact_email: "",
        contact_gender: null,
        contact_signature: "",
        contact_birthday: "",
        contact_notice: "",
        contact_members: [],
        contact_member_cnt: 0,
        contact_owner_id: "",
        contact_add_mode: null,
      },
      ownListReq: {
        owner_id: "",
      },
      allSessionList: [],
      sessionId: "",
      messageList: [],
      innerRef: ref < HTMLDivElement > null,
      scrollbarRef: null,
      addGroupList: [],
      uploadRef: null,
      uploadPath: store.state.backendUrl + "/message/uploadFile",
      fileList: [],
      uploadAvatarRef: null,
      uploadAvatarPath: store.state.backendUrl + "/message/uploadAvatar",
      avatarList: [],
      backendUrl: store.state.backendUrl,
      updateGroupInfo: {
        uuid: "",
        avatar: "",
        add_mode: -1,
        name: "",
        notice: "",
      },
      groupMemberList: [],
      selectedGroupMembers: [],
      removeGroupMembersList: [],
      isAVContainerModalVisible: false,
      videoPlayer: null,
      rtcPeerConn: null,
      ICE_CFG: {},
      localStream: null,
      remoteStream: null,
      remoteVideo: null,
      localVideo: null,
      curContactList: [],
      ableToReceiveOrRejectCall: false,
      ableToStartCall: true,
    });
    
    // 统一的聊天消息处理函数
    const handleChatMessage = async (message) => {
      console.log("💬 [ContactChat] 收到聊天消息:", message);
      
      // 检查消息是否属于当前聊天窗口
      const isGroupChat = data.contactInfo.contact_id && data.contactInfo.contact_id[0] == "G";
      const isCurrentChat = 
        // 群聊：消息的接收方是当前群聊
        (isGroupChat && message.receive_id == data.contactInfo.contact_id) ||
        // 单聊：消息是在当前用户和联系人之间发送的
        (!isGroupChat && 
          ((message.send_id == data.userInfo.uuid && message.receive_id == data.contactInfo.contact_id) ||
           (message.send_id == data.contactInfo.contact_id && message.receive_id == data.userInfo.uuid)));
      
      if (!isCurrentChat) {
        console.log("💬 [ContactChat] 消息不属于当前聊天窗口，忽略");
        return;
      }
      
      if (data.messageList == null) {
        data.messageList = [];
      }
      
      // 如果是加密消息，先解密
      let messageToAdd = message;
      if (message.is_encrypted) {
        try {
          console.log("🔓 开始解密 WebSocket 收到的加密消息");
          messageToAdd = await decryptMessage(message);
          console.log("✅ 消息解密成功");
        } catch (error) {
          console.error("❌ 解密消息失败:", error);
          // 如果解密失败，显示错误提示，但仍然添加到列表
          messageToAdd = {
            ...message,
            content: `[解密失败: ${error.message}]`,
          };
        }
      }
      
      data.messageList.push(messageToAdd);
      scrollToBottom();
      
      // 如果收到的是对方发来的消息（不是自己发的），实时标记为已读
      const isReceivedMessage = message.receive_id == data.userInfo.uuid;
      if (isReceivedMessage && data.sessionId) {
        try {
          await axios.post(store.state.backendUrl + "/session/markAsRead", {
            session_id: data.sessionId
          });
          console.log("✅ [ContactChat] 实时标记会话已读（收到新消息）");
          
          // 清除 Vuex 中该会话的未读数
          store.commit('clearSessionUnreadCount', data.sessionId);
        } catch (error) {
          console.error("❌ [ContactChat] 实时标记会话已读失败:", error);
        }
      }
    };
    
    // 统一的 AV 消息处理函数
    const handleAVMessage = (message) => {
      console.log("📹 [ContactChat] 收到 AV 消息:", message);
      
      var messageAVdata = JSON.parse(message.av_data);
      if (messageAVdata.messageId === "CURRENT_PEERS") {
        console.log(
          "获取CURRENT_PEERS当前在线用户列表，curContactList:",
          messageAVdata.messageData.curContactList
        );
        data.curContactList = messageAVdata.messageData.curContactList;
      } else if (messageAVdata.messageId === "PEER_JOIN") {
        console.log(
          "接受到PEER_JOIN消息，contactId:",
          messageAVdata.messagecontactId
        );
        data.curContactList.push(messageAVdata.messagecontactId);
      } else if (messageAVdata.messageId === "PEER_LEAVE") {
        console.log("接收到PEER_LEAVE消息：", data.userInfo.uuid);
        receiveEndCall();
      } else if (messageAVdata.messageId === "PROXY") {
        console.log("接收到PROXY消息：", message);
        if (messageAVdata.type === "start_call") {
          ElNotification({
            title: "消息提示",
            message: `收到一条来自${message.send_name}的通话请求，请及时前往查看`,
            type: "warning",
          });
          data.ableToReceiveOrRejectCall = true;
          data.ableToStartCall = false;
        } else if (messageAVdata.type === "receive_call") {
          createOffer();
        } else if (messageAVdata.type === "reject_call") {
          endCall();
        } else if (messageAVdata.type === "sdp") {
          if (messageAVdata.messageData.sdp.type === "offer") {
            handleOfferSdp(messageAVdata.messageData.sdp);
          } else if (messageAVdata.messageData.sdp.type === "answer") {
            handleAnswerSdp(messageAVdata.messageData.sdp);
          } else {
            console.log("不支持的sdp类型");
          }
        } else if (messageAVdata.type === "candidate") {
          handleCandidate(messageAVdata.messageData.candidate);
        } else {
          console.log("不支持的proxy类型");
        }
      }
    };
    
    //这是/chat/:id 的id改变时会调用
    onBeforeRouteUpdate(async (to, from, next) => {
      await getChatContactInfo(to.params.id);
      await getSessionId(router.currentRoute.value.params.id);
      if (data.contactInfo.contact_id[0] == "U") {
        await getMessageList();
      } else {
        await getGroupMessageList();
      }
      console.log(data.sessionId);
      scrollToBottom();
      next();
    });
    // 这是刚渲染/chat/:id页面的时候会调用
    onMounted(async () => {
      try {
        console.log(router.currentRoute.value.params.id);
        await loadAllSessions();
        await getChatContactInfo(router.currentRoute.value.params.id);
        await getSessionId(router.currentRoute.value.params.id);
        console.log(data.contactInfo);
        if (data.contactInfo.contact_id[0] == "U") {
          await getMessageList();
        } else {
          await getGroupMessageList();
        }
        console.log(data.sessionId);
        
        // 标记会话为已读
        if (data.sessionId) {
          try {
            await axios.post(store.state.backendUrl + "/session/markAsRead", {
              session_id: data.sessionId
            });
            console.log("✅ [ContactChat] 会话已标记为已读");
            
            // 清除 Vuex 中该会话的未读数
            store.commit('clearSessionUnreadCount', data.sessionId);
          } catch (error) {
            console.error("❌ [ContactChat] 标记会话已读失败:", error);
          }
        }
        
        // 订阅事件总线的消息事件（不再直接设置 onmessage，避免覆盖 App.vue 的处理器）
        console.log("📡 [ContactChat] 订阅事件总线的消息事件");
        eventBus.on('chat:message', handleChatMessage);
        eventBus.on('chat:av_message', handleAVMessage);
        console.log("📡 [ContactChat] ✅ 已订阅事件总线");
        
        scrollToBottom();
      } catch (error) {
        console.error(error);
      }
    });
    
    // 组件卸载时取消订阅
    onUnmounted(() => {
      console.log("📡 [ContactChat] 取消订阅事件总线");
      eventBus.off('chat:message', handleChatMessage);
      eventBus.off('chat:av_message', handleAVMessage);
    });
    const getChatContactInfo = async (id) => {
      try {
        data.getContactInfoReq.contact_id = id;
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/getContactInfo",
          data.getContactInfoReq
        );
        if (!rsp.data.data.contact_avatar.startsWith("http")) {
          rsp.data.data.contact_avatar =
            store.state.backendUrl + rsp.data.data.contact_avatar;
        }
        data.contactInfo = rsp.data.data;
        console.log(data.contactInfo);
      } catch (error) {
        console.log(error);
      }
    };
    const getSessionId = async (contactId) => {
      try {
        const req = {
          send_id: data.userInfo.uuid,
          receive_id: contactId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/session/openSession",
          req
        );
        data.sessionId = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.error(error);
      }
    };

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
    const showUserContactInfoModal = () => {
      data.isUserContactInfoModalVisible = true;
    };
    const quitUserContactInfoModal = () => {
      data.isUserContactInfoModalVisible = false;
    };
    const showGroupContactInfoModal = () => {
      data.isGroupContactInfoModalVisible = true;
    };
    const showUpdateGroupInfoModal = () => {
      data.isUpdateGroupInfoModalVisible = true;
    };
    const quitUpdateGroupInfoModal = () => {
      data.isUpdateGroupInfoModalVisible = false;
    };
    const quitGroupContactInfoModal = () => {
      data.isGroupContactInfoModalVisible = false;
    };
    const showAddGroupModal = () => {
      handleAddGroupList();
    };
    const quitAddGroupModal = () => {
      data.isAddGroupModalVisible = false;
    };
    const handleAddGroupList = async () => {
      try {
        const req = {
          group_id: data.contactInfo.contact_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/getAddGroupList",
          req
        );
        if (rsp.data.code == 200) {
          data.addGroupList = rsp.data.data;
          console.log(rsp.data.data);
          if (data.addGroupList == null) {
            ElMessage.warning("没有新的加群申请");
            return;
          } else {
            for (let i = 0; i < data.addGroupList.length; i++) {
              if (!data.addGroupList[i].contact_avatar.startsWith("http")) {
                data.addGroupList[i].contact_avatar =
                  store.state.backendUrl + data.addGroupList[i].contact_avatar;
              }
            }
            data.isAddGroupModalVisible = true;
            console.log(rsp);
          }
        }
      } catch (error) {
        console.log(error);
      }
    };

    

    const handleToSession = async (session) => {
      router.push("/chat/" + session.id);
    };

    const loadAllSessions = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        
        // 并行加载用户会话和群聊会话
        const [userSessionListRsp, groupSessionListRsp] = await Promise.all([
          axios.post(
            store.state.backendUrl + "/session/getUserSessionList",
            data.ownListReq
          ),
          axios.post(
            store.state.backendUrl + "/session/getGroupSessionList",
            data.ownListReq
          )
        ]);

        const allSessions = [];

        // 处理用户会话
        if (userSessionListRsp.data.data) {
          for (let i = 0; i < userSessionListRsp.data.data.length; i++) {
            const user = userSessionListRsp.data.data[i];
            if (!user.avatar.startsWith("http")) {
              user.avatar = store.state.backendUrl + user.avatar;
            }
            allSessions.push({
              id: user.user_id,
              name: user.user_name,
              avatar: user.avatar,
              type: 'user'
            });
          }
        }

        // 处理群聊会话
        if (groupSessionListRsp.data.data) {
          for (let i = 0; i < groupSessionListRsp.data.data.length; i++) {
            const group = groupSessionListRsp.data.data[i];
            if (!group.avatar.startsWith("http")) {
              group.avatar = store.state.backendUrl + group.avatar;
            }
            allSessions.push({
              id: group.group_id,
              name: group.group_name,
              avatar: group.avatar,
              type: 'group'
            });
          }
        }

        // 按创建时间排序（如果有时间字段，这里先简单按数组顺序）
        data.allSessionList = allSessions;
      } catch (error) {
        console.error(error);
      }
    };
    const preToDeleteSession = () => {
      try {
        ElMessageBox.confirm("确认删除该会话以及其聊天记录？", "Warning", {
          confirmButtonText: "确认",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            deleteSession();
            ElMessage({
              type: "success",
              message: "成功删除",
            });
          })
          .catch(() => {
            ElMessage({
              type: "info",
              message: "取消删除",
            });
          });
      } catch (error) {
        console.error(error);
      }
    };
    const deleteSession = async () => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          session_id: data.sessionId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/session/deleteSession",
          req
        );
        console.log(rsp.data);
      } catch (error) {
        console.error(error);
      }
      router.push("/chat/sessionlist");
    };
    const preToDeleteContact = () => {
      try {
        ElMessageBox.confirm("确认删除该联系人？", "Warning", {
          confirmButtonText: "确认",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            deleteContact();
            ElMessage({
              type: "success",
              message: "成功删除",
            });
          })
          .catch(() => {
            ElMessage({
              type: "info",
              message: "取消删除",
            });
          });
      } catch (error) {
        console.error(error);
      }
    };
    const deleteContact = async () => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: data.contactInfo.contact_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/deleteContact",
          req
        );
        console.log(rsp.data);
      } catch (error) {
        console.error(error);
      }
      router.push("/chat/sessionlist");
    };
    const preToBlackContact = () => {
      try {
        ElMessageBox.confirm("确认拉黑该联系人？", "Warning", {
          confirmButtonText: "确认",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            blackContact();
            ElMessage({
              type: "success",
              message: "成功拉黑",
            });
          })
          .catch(() => {
            ElMessage({
              type: "info",
              message: "取消拉黑",
            });
          });
      } catch (error) {
        console.error(error);
      }
    };
    const blackContact = async () => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          contact_id: data.contactInfo.contact_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/blackContact",
          req
        );
        console.log(rsp.data);
      } catch (error) {
        console.error(error);
      }
      router.push("/chat/sessionlist");
    };

    // 处理输入框键盘事件：Enter 发送，Shift+Enter 换行
    const handleInputKeydown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault(); // 阻止默认换行
        sendMessage();
      }
      // Shift+Enter 保持默认行为（换行）
    };

    const sendMessage = async () => {
      // 检查消息内容
      if (!data.chatMessage || data.chatMessage.trim() === "") {
        console.log("消息内容为空，不发送");
        return;
      }

      // 检查 WebSocket 连接状态
      if (!store.state.socket) {
        console.error("WebSocket 未连接");
        ElMessage.error("WebSocket 未连接，请刷新页面");
        return;
      }

      if (store.state.socket.readyState !== WebSocket.OPEN) {
        console.error("WebSocket 连接状态异常：", store.state.socket.readyState);
        ElMessage.error("WebSocket 连接已断开，请刷新页面");
        return;
      }

      // 检查必要的数据
      if (!data.sessionId || !data.contactInfo.contact_id) {
        console.error("会话信息不完整", {
          sessionId: data.sessionId,
          contactId: data.contactInfo.contact_id
        });
        ElMessage.error("会话信息不完整");
        return;
      }

      console.log("准备发送消息：", data.chatMessage);

      // 检查是否启用了加密
      if (store.state.masterKey && data.contactInfo.contact_id.startsWith('U')) {
        // 启用了加密，且是单聊（不是群聊）
        try {
          await sendEncryptedMessage();
          return;
        } catch (error) {
          console.error("加密消息发送失败：", error);
          ElMessage.error("加密消息发送失败：" + error.message);
          return;
        }
      }

      // 未启用加密，使用原有的明文发送
      const chatMessageRequest = {
        session_id: data.sessionId,
        type: 0,
        content: data.chatMessage,
        url: "",
        send_id: data.userInfo.uuid,
        send_name: data.userInfo.nickname,
        send_avatar: data.userInfo.avatar,
        receive_id: data.contactInfo.contact_id,
        file_size: getFileSize(0),
        file_name: "",
        file_type: "",
      };

      try {
        console.log("发送消息请求（明文）：", chatMessageRequest);
        store.state.socket.send(JSON.stringify(chatMessageRequest));
        console.log("消息已发送");
        data.chatMessage = "";
        scrollToBottom();
      } catch (error) {
        console.error("发送消息失败：", error);
        ElMessage.error("发送消息失败：" + error.message);
      }
    };

    // 发送加密消息
    const sendEncryptedMessage = async () => {
      const contactId = data.contactInfo.contact_id;
      const plaintext = data.chatMessage;

      console.log("🔒 准备发送加密消息...");
      console.log("🔒 [ContactChat] 当前主密钥状态:", {
        has_master_key: !!store.state.masterKey,
        master_key_length: store.state.masterKey?.length,
        master_key_type: typeof store.state.masterKey,
      });

      // 1. 检查会话是否存在
      const sessionExists = await hasSession(contactId);
      let isPreKeyMessage = false;
      let initData = null;

      if (!sessionExists) {
        console.log("会话不存在，正在建立加密会话...");
        ElMessage.info("正在建立安全连接...");

        // 2. 获取对方的公钥束
        try {
          const response = await axios.get("/crypto/getPublicKeyBundle", {
            params: { user_id: contactId },
          });

          if (response.data.code !== 200) {
            throw new Error(response.data.message || "获取公钥束失败");
          }

          const publicKeyBundle = response.data.data;
          console.log("获取到公钥束:", publicKeyBundle);

          // 3. 建立会话
          initData = await createSession(
            store.state.masterKey,
            contactId,
            publicKeyBundle
          );
          isPreKeyMessage = true;
          console.log("✅ 加密会话已建立");
        } catch (error) {
          console.error("建立加密会话失败:", error);
          if (error.message.includes("未启用加密功能")) {
            ElMessage.warning("对方未启用加密功能，将发送明文消息");
            // Fallback 到明文发送（调用原有逻辑）
            // 这里简化处理，抛出错误让用户重试
          }
          throw error;
        }
      }

      // 4. 加密消息
      const encryptedMessage = await encryptAndSendMessage(contactId, plaintext);
      console.log("消息已加密:", encryptedMessage);

      // 5. 构造请求数据
      const requestData = {
        session_id: data.sessionId,
        receiver_id: contactId,
        message_type: isPreKeyMessage ? "PreKeyMessage" : "SignalMessage",
        ...encryptedMessage,
      };

      // 如果是 PreKeyMessage，添加初始化数据
      if (isPreKeyMessage && initData) {
        requestData.sender_identity_key = initData.identity_key;
        requestData.sender_identity_key_curve25519 = initData.identity_key_curve25519; // Curve25519 格式的身份公钥
        requestData.sender_ephemeral_key = initData.ephemeral_key;
        // 确保 used_one_time_pre_key_id 被正确传递（即使是 null 也要传递）
        requestData.used_one_time_pre_key_id = initData.used_one_time_pre_key_id !== undefined 
          ? initData.used_one_time_pre_key_id 
          : initData.usedOneTimePreKeyId; // 兼容两种命名
        console.log('📤 PreKeyMessage 初始化数据:', {
          has_identity_key: !!requestData.sender_identity_key,
          has_identity_key_curve25519: !!requestData.sender_identity_key_curve25519,
          has_ephemeral_key: !!requestData.sender_ephemeral_key,
          used_one_time_pre_key_id: requestData.used_one_time_pre_key_id,
        });
      }

      // 6. 发送到服务器
      console.log("发送加密消息到服务器...");
      console.log("📤 发送的请求数据:", {
        ciphertext_length: requestData.ciphertext?.length,
        ciphertext_preview: requestData.ciphertext?.substring(0, 20),
        iv_length: requestData.iv?.length,
        auth_tag_length: requestData.auth_tag?.length,
        ratchet_key_length: requestData.ratchet_key?.length,
      });
      const response = await axios.post("/message/sendEncryptedMessage", requestData);

      if (response.data.code === 200) {
        console.log("✅ 加密消息发送成功");
        
        // 保存发送方的明文到 IndexedDB（仅发送方可见，用于历史记录）
        // 注意：后端可能返回带尾随空格的 UUID，需要 trim
        let messageId = response.data.data?.message_id;
        if (messageId) {
          messageId = messageId.trim(); // 去除首尾空格
        }
        console.log('📝 尝试保存发送方明文，messageId:', messageId, 'plaintext:', plaintext);
        if (messageId && store.state.masterKey) {
          try {
            const { put, STORES } = await import('@/crypto/cryptoStore');
            // 使用消息 ID 作为 key，存储明文
            const sentMessageData = {
              message_id: messageId,
              plaintext: plaintext,
              contact_id: contactId,
              created_at: Date.now(),
            };
            await put(STORES.SENT_MESSAGES, sentMessageData);
            console.log('✅ 已保存发送方的明文到 IndexedDB:', {
              message_id: messageId,
              plaintext_length: plaintext.length,
              plaintext_preview: plaintext.substring(0, 20),
            });
            
            // 验证保存是否成功
            const { get } = await import('@/crypto/cryptoStore');
            const saved = await get(STORES.SENT_MESSAGES, messageId);
            if (saved) {
              console.log('✅ 验证：IndexedDB 中已存在该消息的明文');
            } else {
              console.warn('⚠️ 验证失败：IndexedDB 中未找到该消息的明文');
            }
          } catch (error) {
            console.error('❌ 保存发送方明文失败:', error);
            // 不影响消息发送，只是历史记录可能无法显示
          }
        } else {
          console.warn('⚠️ 无法保存发送方明文:', {
            has_messageId: !!messageId,
            has_masterKey: !!store.state.masterKey,
          });
        }
        
        // 乐观更新：立即在聊天窗口中显示消息
        // 注意：只有当后端返回 messageId 时才进行乐观更新，否则等待 WebSocket 推送
        if (!messageId) {
          console.warn("⚠️ 后端未返回 messageId，跳过乐观更新，等待 WebSocket 推送");
          data.chatMessage = "";
          scrollToBottom();
          return;
        }
        
        const now = new Date();
        // 格式化时间为 "YYYY-MM-DD HH:mm:ss" 格式（与后端一致）
        const formatDateTime = (date) => {
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, '0');
          const day = String(date.getDate()).padStart(2, '0');
          const hours = String(date.getHours()).padStart(2, '0');
          const minutes = String(date.getMinutes()).padStart(2, '0');
          const seconds = String(date.getSeconds()).padStart(2, '0');
          return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
        };
        
        // 初始化消息列表
        if (data.messageList == null) {
          data.messageList = [];
        }
        
        const messageToDisplay = {
          uuid: messageId, // 使用后端返回的真实 UUID
          send_id: data.userInfo.uuid,
          send_name: data.userInfo.nickname,
          send_avatar: data.userInfo.avatar,
          receive_id: contactId,
          type: 0, // 文本消息
          content: plaintext, // 显示明文（发送方看到的是明文）
          url: "",
          file_type: "",
          file_name: "",
          file_size: "",
          created_at: formatDateTime(now),
          // 加密相关字段
          is_encrypted: true,
          encryption_version: 1,
          message_type: requestData.message_type,
          sender_identity_key: requestData.sender_identity_key,
          sender_identity_key_curve25519: requestData.sender_identity_key_curve25519,
          sender_ephemeral_key: requestData.sender_ephemeral_key,
          used_one_time_pre_key_id: requestData.used_one_time_pre_key_id,
          ratchet_key: requestData.ratchet_key,
          counter: requestData.counter,
          prev_counter: requestData.prev_counter,
          iv: requestData.iv,
          auth_tag: requestData.auth_tag,
        };
        
        // 立即添加到消息列表（乐观更新）
        if (data.messageList == null) {
          data.messageList = [];
        }
        data.messageList.push(messageToDisplay);
        console.log("✅ 已乐观更新消息到聊天窗口，消息ID:", messageId);
        
        // 清空输入框并滚动到底部
        data.chatMessage = "";
        scrollToBottom();
      } else {
        throw new Error(response.data.message || "发送失败");
      }
    };

    const sendFileMessage = async (fileUrl) => {
      const chatFileMessageRequest = {
        session_id: data.sessionId,
        type: 2,
        content: "",
        url: fileUrl,
        send_id: data.userInfo.uuid,
        send_name: data.userInfo.nickname,
        send_avatar: data.userInfo.avatar,
        receive_id: data.contactInfo.contact_id,
        file_size: getFileSize(data.fileList[0].size),
        file_name: data.fileList[0].name,
        file_type: data.fileList[0].type,
      };
      console.log(chatFileMessageRequest);
      store.state.socket.send(JSON.stringify(chatFileMessageRequest));
      scrollToBottom();
    };

    const sendAvatarMessage = (avatarUrl) => {
      const chatAvatarMessageRequest = {
        session_id: data.sessionId,
        type: 2,
        content: "",
        url: avatarUrl,
        send_id: data.userInfo.uuid,
        send_name: data.userInfo.nickname,
        send_avatar: data.userInfo.avatar,
        receive_id: data.contactInfo.contact_id,
        file_size: getFileSize(data.avatarList[0].size),
        file_name: data.avatarList[0].name,
        file_type: data.avatarList[0].type,
      };
      console.log(chatAvatarMessageRequest);
      store.state.socket.send(JSON.stringify(chatAvatarMessageRequest));
      scrollToBottom();
    };

    // 处理加密文件上传完成
    const handleFileUploadComplete = async (uploadResult) => {
      console.log('📁 加密文件上传完成:', uploadResult);
      
      const contactId = data.contactInfo.contact_id;
      if (!contactId || !contactId.startsWith('U')) {
        ElMessage.error('加密文件只能发送给单聊好友');
        return;
      }

      try {
        // 1. 检查是否有加密会话
        const sessionExists = await hasSession(contactId);
        let isPreKeyMessage = false;
        let initData = null;

        if (!sessionExists) {
          console.log("会话不存在，正在建立加密会话...");
          ElMessage.info("正在建立安全连接...");

          // 获取对方的公钥束并建立会话
          const response = await axios.get("/crypto/getPublicKeyBundle", {
            params: { user_id: contactId },
          });

          if (response.data.code !== 200) {
            throw new Error(response.data.message || "获取公钥束失败");
          }

          const publicKeyBundle = response.data.data;
          initData = await createSession(
            store.state.masterKey,
            contactId,
            publicKeyBundle
          );
          isPreKeyMessage = true;
          console.log("✅ 加密会话已建立");
        }

        // 2. 构造文件元数据
        const fileMetadata = JSON.stringify({
          type: uploadResult.width ? 'image' : 'file',
          ossKey: uploadResult.ossKey,
          fileKey: uploadResult.fileKey,
          fileIv: uploadResult.fileIv,
          fileHash: uploadResult.fileHash,
          fileName: uploadResult.fileName,
          fileSize: uploadResult.fileSize,
          mimeType: uploadResult.mimeType,
          width: uploadResult.width || null,
          height: uploadResult.height || null,
          thumbnail: uploadResult.thumbnail || null,
        });

        // 3. 使用 Signal 协议加密文件元数据
        const encryptedMessage = await encryptAndSendMessage(contactId, fileMetadata);
        console.log("文件元数据已加密:", encryptedMessage);

        // 4. 构造请求数据
        const messageType = uploadResult.width ? 4 : 5; // 4=加密图片, 5=加密文件
        const requestData = {
          session_id: data.sessionId,
          receiver_id: contactId,
          message_type: isPreKeyMessage ? "PreKeyMessage" : "SignalMessage",
          file_message_type: messageType,
          ...encryptedMessage,
        };

        // 如果是 PreKeyMessage，添加初始化数据
        if (isPreKeyMessage && initData) {
          requestData.sender_identity_key = initData.identity_key;
          requestData.sender_identity_key_curve25519 = initData.identity_key_curve25519;
          requestData.sender_ephemeral_key = initData.ephemeral_key;
          requestData.used_one_time_pre_key_id = initData.used_one_time_pre_key_id !== undefined 
            ? initData.used_one_time_pre_key_id 
            : initData.usedOneTimePreKeyId;
        }

        // 5. 发送到服务器
        console.log("发送加密文件消息到服务器...");
        const response = await axios.post("/message/sendEncryptedMessage", requestData);

        if (response.data.code === 200) {
          console.log("✅ 加密文件消息发送成功");
          ElMessage.success(uploadResult.width ? '图片发送成功' : '文件发送成功');
          
          // 保存发送方的明文元数据到 IndexedDB
          let messageId = response.data.data?.message_id;
          if (messageId) {
            messageId = messageId.trim();
            try {
              const { put, STORES } = await import('@/crypto/cryptoStore');
              await put(STORES.SENT_MESSAGES, {
                message_id: messageId,
                plaintext: fileMetadata,
                contact_id: contactId,
                created_at: Date.now(),
              });
            } catch (error) {
              console.error('保存文件消息元数据失败:', error);
            }
          }
          
          // 乐观更新：立即显示消息
          if (messageId) {
            const now = new Date();
            const formatDateTime = (date) => {
              const year = date.getFullYear();
              const month = String(date.getMonth() + 1).padStart(2, '0');
              const day = String(date.getDate()).padStart(2, '0');
              const hours = String(date.getHours()).padStart(2, '0');
              const minutes = String(date.getMinutes()).padStart(2, '0');
              const seconds = String(date.getSeconds()).padStart(2, '0');
              return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
            };

            if (data.messageList == null) {
              data.messageList = [];
            }

            data.messageList.push({
              uuid: messageId,
              send_id: data.userInfo.uuid,
              send_name: data.userInfo.nickname,
              send_avatar: data.userInfo.avatar,
              receive_id: contactId,
              type: messageType,
              content: fileMetadata,
              url: "",
              file_type: uploadResult.mimeType,
              file_name: uploadResult.fileName,
              file_size: uploadResult.fileSize,
              created_at: formatDateTime(now),
              is_encrypted: true,
              encryption_version: 1,
            });

            scrollToBottom();
          }
        } else {
          throw new Error(response.data.message || "发送失败");
        }
      } catch (error) {
        console.error("加密文件消息发送失败:", error);
        ElMessage.error("发送失败：" + error.message);
      }
    };

    // 处理文件上传错误
    const handleFileUploadError = (error) => {
      console.error("文件上传错误:", error);
      ElMessage.error("文件上传失败：" + error.message);
    };

    // 解析文件信息（从消息内容中）
    const parseFileInfo = (messageItem) => {
      try {
        // 如果 content 是字符串，尝试解析为 JSON
        if (typeof messageItem.content === 'string') {
          return JSON.parse(messageItem.content);
        }
        return messageItem.content;
      } catch (error) {
        console.error('解析文件信息失败:', error);
        return {
          type: messageItem.type === 4 ? 'image' : 'file',
          fileName: messageItem.file_name || '未知文件',
          fileSize: messageItem.file_size || 0,
        };
      }
    };

    const getMessageList = async () => {
      try {
        console.log(data.contactInfo);
        const req = {
          user_one_id: data.userInfo.uuid,
          user_two_id: data.contactInfo.contact_id,
        };
        console.log(req);
        const rsp = await axios.post(
          "/message/getMessageList",
          req
        );
        if (rsp.data.data) {
          for (let i = 0; i < rsp.data.data.length; i++) {
            if (!rsp.data.data[i].send_avatar.startsWith("http")) {
              rsp.data.data[i].send_avatar =
                store.state.backendUrl + rsp.data.data[i].send_avatar;
            }
          }

          // 解密加密消息
          try {
            rsp.data.data = await decryptMessageList(rsp.data.data);
            console.log('消息列表解密完成');
          } catch (error) {
            console.error('消息解密失败:', error);
            // 不阻断，继续显示消息（可能是明文或解密失败的提示）
          }
        }
        data.messageList = rsp.data.data;
        console.log(data.messageList);
        console.log(rsp);
      } catch (error) {
        console.error(error);
      }
    };

    const getGroupMessageList = async () => {
      try {
        console.log(data.contactInfo);
        const req = {
          group_id: data.contactInfo.contact_id,
        };
        console.log(req);
        const rsp = await axios.post(
          store.state.backendUrl + "/message/getGroupMessageList",
          req
        );
        if (rsp.data.data) {
          for (let i = 0; i < rsp.data.data.length; i++) {
            if (!rsp.data.data[i].send_avatar.startsWith("http")) {
              rsp.data.data[i].send_avatar =
                store.state.backendUrl + rsp.data.data[i].send_avatar;
            }
          }
        }
        data.messageList = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.error(error);
      }
    };

    const scrollToBottom = () => {
      nextTick(() => {
        // 检查 DOM 元素是否存在
        if (!data.innerRef || !data.scrollbarRef) {
          console.warn("scrollToBottom: DOM 元素未准备好");
          return;
        }
        try {
          const scrollHeight = data.innerRef.scrollHeight;
          console.log(scrollHeight);
          data.scrollbarRef.setScrollTop(scrollHeight);
        } catch (error) {
          console.warn("scrollToBottom 失败:", error);
        }
      });
    };

    const handleAgree = async (contactId) => {
      try {
        const req = {
          owner_id: data.contactInfo.contact_id,
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
          console.log("✅ [ContactChat] 已同意好友申请，等待对方发送消息时自动建立加密会话");
          
          data.addGroupList = data.addGroupList.filter(
            (c) => c.contact_id !== contactId
          );
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
          owner_id: data.contactInfo.contact_id,
          contact_id: contactId,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/contact/refuseContactApply",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          data.addGroupList = data.addGroupList.filter(
            (c) => c.contact_id !== contactId
          );
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };

    const handleLeaveGroup = async () => {
      try {
        const req = {
          user_id: data.userInfo.uuid,
          group_id: data.contactInfo.contact_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/group/leaveGroup",
          req
        );
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
          router.push("/chat/sessionlist");
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else {
          ElMessage.error(rsp.data.message);
          console.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };

    const handleDismissGroup = async () => {
      try {
        const req = {
          owner_id: data.userInfo.uuid,
          group_id: data.contactInfo.contact_id,
        };
        const rsp = await axios.post(
          store.state.backendUrl + "/group/dismissGroup",
          req
        );
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          console.log(rsp.data.message);
          router.push("/chat/sessionlist");
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
          console.log(rsp.data.message);
        } else {
          ElMessage.error(rsp.data.message);
          console.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };

    const handleUploadSuccess = () => {
      ElMessage.success("文件上传成功");
      sendFileMessage(
        store.state.backendUrl + "/static/files/" + data.fileList[0].name
      );
      data.fileList = [];
    };

    const handleAvatarUploadSuccess = () => {
      ElMessage.success("头像上传成功");
      data.avatarList = [];
    };

    const beforeAvatarUpload = (avatar) => {
      console.log("上传前avatar====>", avatar);
      console.log(data.avatarList);
      console.log(avatar);
      if (data.avatarList.length > 1) {
        ElMessage.error("只能上传一张头像");
        return false;
      }
      const isLt50M = avatar.size / 1024 / 1024 < 50;
      if (!isLt50M) {
        ElMessage.error("上传头像图片大小不能超过 50MB!");
        return false;
      }
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
    const downloadFile = async (fileName) => {
      try {
        const rsp = await axios.get(
          store.state.backendUrl + "/static/files/" + fileName,
          {
            responseType: "blob",
          }
        );
        console.log(rsp);
        const blob = new Blob([rsp.data], {
          type: rsp.headers["content-type"] || "application/octet-stream",
        });
        const link = document.createElement("a");
        link.href = window.URL.createObjectURL(blob);
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } catch (error) {
        console.error(error);
      }
    };
    const getFileSize = (size) => {
      if (size < 1024) {
        return size + "B";
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + "KB";
      } else if (size < 1024 * 1024 * 1024) {
        return (size / 1024 / 1024).toFixed(2) + "MB";
      } else {
        return (size / 1024 / 1024 / 1024).toFixed(2) + "GB";
      }
    };

    const handleUpdateGroupInfo = async () => {
      try {
        if (
          data.updateGroupInfo.name == "" &&
          data.updateGroupInfo.notice == "" &&
          data.updateGroupInfo.add_mode == -1 &&
          data.avatarList.length == 0
        ) {
          ElMessage.error("请至少修改一项");
          return;
        }
        if (data.avatarList.length > 0) {
          data.updateGroupInfo.avatar =
            "/static/avatars/" + data.avatarList[0].name;
          data.uploadAvatarRef.submit();
        }
        data.updateGroupInfo.uuid = data.contactInfo.contact_id;
        const rsp = await axios.post(
          store.state.backendUrl + "/group/updateGroupInfo",
          data.updateGroupInfo
        );
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          data.isUpdateGroupInfoModalVisible = false;
          await getChatContactInfo(router.currentRoute.value.params.id);
        } else {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };

    const closeUpdateGroupInfoModal = () => {
      handleUpdateGroupInfo();
    };
    const getGroupMemberList = async () => {
      const req = {
        group_id: data.contactInfo.contact_id,
      };
      try {
        const rsp = await axios.post(
          store.state.backendUrl + "/group/getGroupMemberList",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          for (let i = 0; i < rsp.data.data.length; i++) {
            if (!rsp.data.data[i].avatar.startsWith("http")) {
              rsp.data.data[i].avatar =
                store.state.backendUrl + rsp.data.data[i].avatar;
            }
          }
          data.groupMemberList = rsp.data.data;
          console.log(data.groupMemberList);
        } else {
          ElMessage.error(rsp.data.message);
          console.log(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const showRemoveGroupMemberModal = async () => {
      await getGroupMemberList();
      data.isRemoveGroupMemberModalVisible = true;
    };

    const quitRemoveGroupMemberModal = () => {
      data.isRemoveGroupMemberModalVisible = false;
    };

    const closeRemoveGroupMemberModal = () => {};

    const handleCheckboxChange = () => {
      data.removeGroupMembersList = data.selectedGroupMembers;
      console.log(data.removeGroupMembersList);
    };

    const handleRemoveGroupMembers = async () => {
      const req = {
        group_id: data.contactInfo.contact_id,
        owner_id: data.contactInfo.contact_owner_id,
        uuid_list: data.removeGroupMembersList,
      };
      console.log(data.contactInfo);
      try {
        const rsp = await axios.post(
          store.state.backendUrl + "/group/removeGroupMembers",
          req
        );
        console.log(rsp);
        if (rsp.data.code == 200) {
          ElMessage.success(rsp.data.message);
          getGroupMemberList();
        } else if (rsp.data.code == 400) {
          ElMessage.warning(rsp.data.message);
        } else {
          ElMessage.error(rsp.data.message);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const showAVContainerModal = () => {
      data.isAVContainerModalVisible = true;
    };

    const closeAVContainerModal = () => {
      if (data.localVideo || data.remoteVideo) {
        ElMessage.warning("请先结束通话");
        return;
      }
      data.isAVContainerModalVisible = false;
    };

    const createRtcPeerConnection = () => {
      if (data.rtcPeerConn) {
        console.log("peer connection has already been created.");
        return;
      }
      data.rtcPeerConn = new RTCPeerConnection(data.ICE_CFG);
      data.rtcPeerConn.onicecandidate = (event) => {
        if (event.candidate) {
          var proxyCandidateMessage = {
            messageId: "PROXY",
            type: "candidate",
            messageData: {
              candidate: event.candidate,
            },
          };
          const rtcMessageRequest = {
            session_id: data.sessionId,
            type: 3,
            content: "",
            url: "",
            send_id: data.userInfo.uuid,
            send_name: data.userInfo.nickname,
            send_avatar: data.userInfo.avatar,
            receive_id: data.contactInfo.contact_id,
            file_size: "",
            file_name: "",
            file_type: "",
            av_data: JSON.stringify(proxyCandidateMessage),
          };
          console.log(rtcMessageRequest);
          store.state.socket.send(JSON.stringify(rtcMessageRequest));
        }
      };
      data.rtcPeerConn.oniceconnectionstatechange = (event) => {
        console.log(
          "oniceconnectionstatechange",
          data.rtcPeerConn.iceConnectionState
        );
      };
      // 对端传来媒体轨道
      data.rtcPeerConn.ontrack = (event) => {
        if (data.remoteStream === null) {
          data.remoteStream = new MediaStream();
          data.remoteVideo = document.querySelector("video.remote-video");
          data.remoteVideo.srcObject = data.remoteStream;
          data.remoteVideo.style.display = "inline-block";
        }
        data.remoteStream.addTrack(event.track);
      };
    };

    const closeRtcPeerConnection = () => {
      if (data.rtcPeerConn) {
        data.rtcPeerConn.close();
        data.rtcPeerConn = null;
      }
    };

    const getLocalMediaStream = () => {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        console.log("getUserMedia is not supported!");
        return null;
      } else if (data.localStream) {
        console.log("localStream already exist.");
        return data.localStream;
      } else {
        var constraints = {
          video: true,
          audio: true,
        };
        return navigator.mediaDevices.getUserMedia(constraints);
      }
    };

    const closeLocalMediaStream = () => {
      if (data.localStream != null) {
        data.localStream.getTracks().forEach((track) => {
          track.stop();
        });
        data.localStream = null;
      }
    };

    const attachMediaStreamToLocalVideo = () => {
      data.localVideo = document.querySelector("video.local-video");
      data.localVideo.srcObject = data.localStream;
      data.localVideo.muted = true;
      data.localVideo.style.display = "inline-block";
    };

    const attachMediaStreamToPeerConnection = () => {
      data.localStream.getTracks().forEach((track) => {
        data.rtcPeerConn.addTrack(track);
      });
    };

    const createOffer = () => {
      var offerOpts = {
        offerToReceiveAudio: true,
        offerToReceiveVideo: true,
      };
      data.rtcPeerConn
        .createOffer(offerOpts)
        .then((desc) => {
          data.rtcPeerConn.setLocalDescription(desc);
          var proxySdpMessage = {
            messageId: "PROXY",
            type: "sdp",
            messageData: {
              sdp: desc,
            },
          };
          console.log(desc);
          const rtcMessageRequest = {
            session_id: data.sessionId,
            type: 3,
            content: "",
            url: "",
            send_id: data.userInfo.uuid,
            send_name: data.userInfo.nickname,
            send_avatar: data.userInfo.avatar,
            receive_id: data.contactInfo.contact_id,
            file_size: "",
            file_name: "",
            file_type: "",
            av_data: JSON.stringify(proxySdpMessage),
          };
          store.state.socket.send(JSON.stringify(rtcMessageRequest));
        })
        .catch((err) => {
          console.log(
            `createOffer failed, error name: ${err.name}, error message: ${err.message}`
          );
        });
    };

    const createAnswer = () => {
      data.rtcPeerConn
        .createAnswer()
        .then((desc) => {
          data.rtcPeerConn.setLocalDescription(desc);
          var proxySdpMessage = {
            messageId: "PROXY",
            type: "sdp",
            messageData: {
              sdp: desc,
            },
          };
          console.log(desc);
          const rtcMessageRequest = {
            session_id: data.sessionId,
            type: 3,
            content: "",
            url: "",
            send_id: data.userInfo.uuid,
            send_name: data.userInfo.nickname,
            send_avatar: data.userInfo.avatar,
            receive_id: data.contactInfo.contact_id,
            file_size: "",
            file_name: "",
            file_type: "",
            av_data: JSON.stringify(proxySdpMessage),
          };
          store.state.socket.send(JSON.stringify(rtcMessageRequest));
        })
        .catch((err) => {
          console.log(
            `createAnswer failed, error name: ${err.name}, error message: ${err.message}`
          );
        });
    };

    const startCall = async (isInitiator) => {
      console.log(data.localVideo);
      console.log(data.localStream);
      if (data.localVideo) {
        ElMessage.warning("已经在通话中，请勿重复发起");
        return;
      }
      if (isInitiator && !data.ableToStartCall) {
        ElMessage.warning(
          "对方已经发起通话，请先接收通话或拒绝通话，才能发起下一次通话"
        );
        return;
      }
      if (!isInitiator && !data.ableToReceiveOrRejectCall) {
        ElMessage.warning("对方没有发起通话或已在通话中，无法接收通话");
        return;
      }
      createRtcPeerConnection();
      data.localStream = await getLocalMediaStream();
      attachMediaStreamToLocalVideo();
      attachMediaStreamToPeerConnection();
      if (isInitiator) {
        var startCallMessage = {
          messageId: "PROXY",
          type: "start_call",
        };
        const rtcMessageRequest = {
          session_id: data.sessionId,
          type: 3,
          content: "",
          url: "",
          send_id: data.userInfo.uuid,
          send_name: data.userInfo.nickname,
          send_avatar: data.userInfo.avatar,
          receive_id: data.contactInfo.contact_id,
          file_size: "",
          file_name: "",
          file_type: "",
          av_data: JSON.stringify(startCallMessage),
        };
        store.state.socket.send(JSON.stringify(rtcMessageRequest));
      } else {
        var receiveCallMessage = {
          messageId: "PROXY",
          type: "receive_call",
        };
        const rtcMessageRequest = {
          session_id: data.sessionId,
          type: 3,
          content: "",
          url: "",
          send_id: data.userInfo.uuid,
          send_name: data.userInfo.nickname,
          send_avatar: data.userInfo.avatar,
          receive_id: data.contactInfo.contact_id,
          file_size: "",
          file_name: "",
          file_type: "",
          av_data: JSON.stringify(receiveCallMessage),
        };
        store.state.socket.send(JSON.stringify(rtcMessageRequest));
        data.ableToReceiveOrRejectCall = false;
      }
    };

    const sendEndCall = () => {
      if (data.localVideo == null && data.remoteVideo == null) {
        ElMessage.warning("尚未开始通话，无法挂断");
        return;
      }
      if (data.localVideo) data.localVideo.style.display = "none";
      if (data.remoteVideo) data.remoteVideo.style.display = "none";
      closeLocalMediaStream();
      closeRtcPeerConnection();
      data.remoteStream = null;
      data.localStream = null;
      data.localVideo = null;
      data.remoteVideo = null;
      data.ableToReceiveOrRejectCall = false;
      data.ableToStartCall = true;
      var proxyPeerLeaveMessage = {
        messageId: "PEER_LEAVE",
      };
      const rtcMessageRequest = {
        session_id: data.sessionId,
        type: 3,
        content: "",
        url: "",
        send_id: data.userInfo.uuid,
        send_name: data.userInfo.nickname,
        send_avatar: data.userInfo.avatar,
        receive_id: data.contactInfo.contact_id,
        file_size: "",
        file_name: "",
        file_type: "",
        av_data: JSON.stringify(proxyPeerLeaveMessage),
      };
      store.state.socket.send(JSON.stringify(rtcMessageRequest));
    };

    const endCall = () => {
      if (data.localVideo) data.localVideo.style.display = "none";
      if (data.remoteVideo) data.remoteVideo.style.display = "none";
      closeLocalMediaStream();
      closeRtcPeerConnection();
      data.remoteStream = null;
      data.localStream = null;
      data.localVideo = null;
      data.remoteVideo = null;
      data.ableToReceiveOrRejectCall = false;
      data.ableToStartCall = true;
      ElMessage.warning("对方拒绝通话");
    };

    const receiveEndCall = () => {
      if (data.localVideo) data.localVideo.style.display = "none";
      if (data.remoteVideo) data.remoteVideo.style.display = "none";
      closeLocalMediaStream();
      closeRtcPeerConnection();
      data.remoteStream = null;
      data.localStream = null;
      data.localVideo = null;
      data.remoteVideo = null;
      data.ableToReceiveOrRejectCall = false;
      data.ableToStartCall = true;
      ElMessage.warning("对方已挂断");
    };

    const handleOfferSdp = (val) => {
      data.rtcPeerConn
        .setRemoteDescription(new RTCSessionDescription(val))
        .then(() => {
          createAnswer();
        })
        .catch((err) => {
          console.log("rtcPeerConn setRemoteDescription failed", err);
        });
    };

    const handleAnswerSdp = (val) => {
      data.rtcPeerConn
        .setRemoteDescription(new RTCSessionDescription(val))
        .catch((err) => {
          console.log("rtcPeerConn setRemoteDescription failed", err);
        });
    };

    const handleCandidate = (val) => {
      data.rtcPeerConn.addIceCandidate(new RTCIceCandidate(val));
    };

    const rejectCall = () => {
      if (!data.ableToReceiveOrRejectCall) {
        ElMessage.warning("对方没有发起通话或已在通话中，无法拒绝通话");
        return;
      }
      var rejectCallMessage = {
        messageId: "PROXY",
        type: "reject_call",
      };
      const rtcMessageRequest = {
        session_id: data.sessionId,
        type: 3,
        content: "",
        url: "",
        send_id: data.userInfo.uuid,
        send_name: data.userInfo.nickname,
        send_avatar: data.userInfo.avatar,
        receive_id: data.contactInfo.contact_id,
        file_size: "",
        file_name: "",
        file_type: "",
        av_data: JSON.stringify(rejectCallMessage),
      };
      store.state.socket.send(JSON.stringify(rtcMessageRequest));
      data.ableToReceiveOrRejectCall = false;
    };

    return {
      ...toRefs(data),
      router,
      handleCreateGroup,
      showUserContactInfoModal,
      quitUserContactInfoModal,
      showGroupContactInfoModal,
      quitGroupContactInfoModal,
      showAddGroupModal,
      quitAddGroupModal,
      handleToSession,
      loadAllSessions,
      deleteSession,
      preToDeleteSession,
      preToDeleteContact,
      preToBlackContact,
      sendMessage,
      handleInputKeydown,
      getMessageList,
      getGroupMessageList,
      handleAgree,
      handleReject,
      handleAddGroupList,
      handleLeaveGroup,
      handleDismissGroup,
      handleUploadSuccess,
      beforeFileUpload,
      downloadFile,
      getFileSize,
      handleFileUploadComplete,
      handleFileUploadError,
      parseFileInfo,
      showUpdateGroupInfoModal,
      quitUpdateGroupInfoModal,
      beforeAvatarUpload,
      handleAvatarUploadSuccess,
      handleUpdateGroupInfo,
      closeUpdateGroupInfoModal,
      showRemoveGroupMemberModal,
      quitRemoveGroupMemberModal,
      closeRemoveGroupMemberModal,
      getGroupMemberList,
      handleCheckboxChange,
      handleRemoveGroupMembers,
      createRtcPeerConnection,
      closeRtcPeerConnection,
      getLocalMediaStream,
      closeLocalMediaStream,
      attachMediaStreamToLocalVideo,
      attachMediaStreamToPeerConnection,
      createOffer,
      createAnswer,
      startCall,
      sendEndCall,
      receiveEndCall,
      handleOfferSdp,
      handleAnswerSdp,
      handleCandidate,
      showAVContainerModal,
      closeAVContainerModal,
      rejectCall,
      endCall,
    };
  },
};
</script>

<style scoped>
.chat-footer-toolbar {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #ebeef5;
  gap: 8px;
}

.sessionlist-header {
  display: flex;
  flex-direction: row;
  width: 100%;
  margin-top: 10px;
  margin-bottom: 10px;
}

.contact-search-input {
  width: 215px;
  height: 30px;
  margin-left: 5px;
  margin-right: 2px;
}

.sessionlist-container,
.contactlist-body,
.contactlist-user {
  padding: 0 !important;
  margin: 0 !important;
  width: 100%;
}

.el-menu {
  background-color: #f8f9fa;
  width: 100% !important;
  border: none;
  padding: 0 !important;
  margin: 0 !important;
}

.el-menu-item {
  background-color: #ffffff;
  height: 48px;
  border-radius: 0;
  margin: 0 !important;
  padding-left: 8px !important;
  padding-right: 0 !important;
  transition: all 0.2s ease;
  border-bottom: 1px solid #f0f0f0;
  width: 100% !important;
  box-sizing: border-box;
}

.el-menu-item * {
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

.sessionlist-title {
  font-family: Arial, Helvetica, sans-serif;
}

h3 {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(69, 69, 68);
}

.groupcontactinfo-modal-quit-btn-container,
.userinfo-modal-quit-btn-container,
.updategroupinfo-modal-quit-btn-container,
.removegroupmember-modal-quit-btn-container {
  height: 25px;
  width: 100%;
  display: flex;
  flex-direction: row-reverse;
}

.groupcontactinfo-modal-quit-btn,
.userinfo-modal-quit-btn,
.updategroupinfo-modal-quit-btn,
.removegroupmember-modal-quit-btn {
  background-color: rgba(255, 255, 255, 0);
  color: rgb(229, 25, 25);
  padding: 15px;
  border: none;
  cursor: pointer;
  position: fixed;
  justify-content: center;
  align-items: center;
}

.groupcontactinfo-modal-header-title,
.userinfo-modal-header-title,
.removegroupmember-modal-header-title {
  height: 30px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.updategroupinfo-modal-header-title {
  margin-top: 10px;
  height: 30px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.el-header {
  padding: 10px;
  display: flex;
  flex-direction: row;
}

.chat-title {
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 50%;
  height: 100%;
}

.chat-title-right {
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
  width: 50%;
  height: 100%;
  margin-right: 10px;
  color: rgb(201, 139, 139);
}

.sessionlist-avatar {
  width: 30px;
  height: 30px;
  margin-right: 20px;
}

.setting-btn {
  background-color: rgba(255, 255, 255, 0);
  border: none;
  cursor: pointer;
  color: rgb(201, 139, 139);
}

.modal-list {
  height: 270px;
  width: 90%;
  margin-top: 5px;
}

.left-message {
  width: 67%;
  height: 100%;
  display: flex;
  flex-direction: row;
  margin-bottom: 10px;
}

.right-message {
  width: 67%;
  height: 100%;
  display: flex;
  flex-direction: row-reverse;
  margin-bottom: 10px;
}

.left-message-left {
  display: flex;
  justify-content: center;
}

.right-message-right {
  display: flex;
  justify-content: center;
}

.left-message-right-top {
  display: flex;
  flex-direction: row;
  margin-bottom: 5px;
}

.right-message-left-top {
  display: flex;
  flex-direction: row-reverse;
  margin-bottom: 5px;
}

.left-message-contactname {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(77, 61, 192);
  font-weight: bold;
  margin-top: 5px;
  margin-right: 10px;
  font-size: 15px;
}

.right-message-contactname {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(77, 61, 192);
  font-weight: bold;
  margin-top: 5px;
  margin-left: 10px;
  font-size: 15px;
}

.left-message-time {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(237, 161, 161);
  margin-top: 5px;
  font-size: 15px;
}

.right-message-time {
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(237, 161, 161);
  margin-top: 5px;
  font-size: 15px;
}

.left-message-content {
  background-color: #ffffff;
  color: rgb(74, 72, 72);
  display: inline-block;
  max-width: 400px;
  white-space: pre-wrap; /* 保留换行符并允许自动换行 */
  word-break: break-word; /* 长单词换行 */
  font-family: Arial, Helvetica, sans-serif;
  border-radius: 6px;
  padding: 3px;
  padding-right: 5px;
  font-size: 14px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.right-message-content {
  background: #4facfe;
  border: none;
  color: #ffffff;
  color: rgb(74, 72, 72);
  display: inline-block;
  max-width: 400px;
  white-space: pre-wrap; /* 保留换行符并允许自动换行 */
  word-break: break-word; /* 长单词换行 */
  font-family: Arial, Helvetica, sans-serif;
  border-radius: 6px;
  padding: 3px;
  padding-right: 5px;
  font-size: 14px;
  box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.2);
}

.left-message-file-container {
  background-color: #f9f9f9; /* 浅灰色背景 */
  border: 1px solid #ddd; /* 浅灰色边框 */
  border-radius: 8px; /* 圆角边框 */
  padding: 16px; /* 内边距 */
  width: 250px;
  height: 100px;
  margin: 0 auto; /* 水平居中 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 轻微阴影效果 */
}

.right-message-file-container {
  background-color: #f9f9f9; /* 浅灰色背景 */
  border: 1px solid #ddd; /* 浅灰色边框 */
  border-radius: 8px; /* 圆角边框 */
  padding: 16px; /* 内边距 */
  width: 250px;
  height: 100px;
  margin: 0 auto; /* 水平居中 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 轻微阴影效果 */
}

.right-message-file-name {
  font-size: 12px;
  font-weight: bold;
  color: #333;
  margin-right: 5px;
  font-family: Arial, Helvetica, sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
}

.left-message-file-name {
  font-size: 12px;
  font-weight: bold;
  color: #333;
  margin-right: 5px;
  font-family: Arial, Helvetica, sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
}

.right-message-file-size {
  font-size: 11px;
  color: rgb(176, 172, 172);
  font-family: Arial, Helvetica, sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 4px;
}

.left-message-file-size {
  font-size: 11px;
  color: rgb(176, 172, 172);
  font-family: Arial, Helvetica, sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 4px;
}

.right-message-file-download {
  font-size: 12px;
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(176, 172, 172);
  margin-top: 20px;
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

.addGroup-modal-body {
  height: 70%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.newcontact-modal-footer,
.updategroupinfo-modal-footer {
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

.addGroup-list {
  width: 280px;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.addGroup-item {
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

.removegroupmembers-list {
  width: 280px;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: Arial, Helvetica, sans-serif;
}

.removegroupmembers-item {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 40px;
}

.removegroupmembers-item-avatar {
  height: 30px;
  width: 30px;
  margin-right: 20px;
}
.removegroupmembers-item-name {
  font-size: 14px;
  font-weight: bold;
  font-family: Arial, Helvetica, sans-serif;
  color: rgb(57, 57, 57);
}
.removegroupmembers-button {
  background: #4facfe;
  border: none;
  color: #ffffff;
}

.video-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  border-radius: 30px;
}

.video-modal-content {
  background: #fff;
  height: 500px;
  border-radius: 20px;
  width: 800px;
  box-shadow: 0 2px 15px rgb(195, 8, 8);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.local-video,
.remote-video {
  width: 330px;
  height: 320px;
}

.video-modal-header {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-modal-body {
  height: 360px;
  width: 700px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 15px rgb(250, 65, 109);
  border-radius: 20px;
}

.video-modal-footer {
  height: 80px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-modal-footer-btn {
  background: #4facfe;
  border: none;
  color: #ffffff;
}

/* 未读消息样式 */
.unread-message {
  position: relative;
  background-color: #f0f9ff !important; /* 浅蓝色背景 */
  border-left: 3px solid #409EFF !important; /* 左侧蓝色边框 */
}

.unread-indicator {
  display: inline-block;
  width: 6px;
  height: 6px;
  background-color: #f56c6c;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
}
</style>

