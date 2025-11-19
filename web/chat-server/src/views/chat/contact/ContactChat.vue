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
                <el-menu
                  router
                  unique-opened
                  @open="handleShowUserSessionList"
                  @close="handleHideUserSessionList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="sessionlist-title">用户</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="user in userSessionList"
                    :key="user.user_id"
                    @click="handleToChatUser(user)"
                  >
                    <img :src="user.avatar" class="sessionlist-avatar" />
                    {{ user.user_name }}
                  </el-menu-item>
                </el-menu>
                <el-menu
                  router
                  unique-opened
                  @open="handleShowGroupSessionList"
                  @close="handleHideGroupSessionList"
                >
                  <el-sub-menu index="1">
                    <template #title>
                      <span class="sessionlist-title">群聊</span>
                    </template>
                  </el-sub-menu>
                  <el-menu-item
                    v-for="group in groupSessionList"
                    :key="group.group_id"
                    @click="handleToChatGroup(group)"
                  >
                    <img :src="group.avatar" class="sessionlist-avatar" />
                    {{ group.group_name }}
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
                  <div class="updategroupinfo-modal-footer">
                    <el-button
                      style="background-color: rgb(252, 210.9, 210.9)"
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

                      <div class="left-message-content">
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
                              background-color: rgb(252, 210.9, 210.9);
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
                  </div>
                </div>
              </div>
            </el-scrollbar>
            <div class="tool-bar">
              <div class="tool-bar-left">
                <el-tooltip
                  effect="customized"
                  content="表情包"
                  placement="top"
                  hide-after="0"
                  enterable="false"
                >
                  <button
                    class="image-button"
                    @click="
                      downloadFile(backendUrl + '/static/avatars', '头像.jpg')
                    "
                  >
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
                    <el-upload
                      v-model:file-list="fileList"
                      ref="uploadRef"
                      :auto-upload="true"
                      :show-file-list="false"
                      :action="uploadPath"
                      :on-success="handleUploadSuccess"
                      :before-upload="beforeFileUpload"
                      style="
                        display: flex;
                        align-items: center;
                        justify-content: center;
                      "
                    >
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
                    </el-upload>
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
                placeholder="请输入内容"
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
import { reactive, toRefs, onMounted, ref, nextTick } from "vue";
import { useRouter, onBeforeRouteUpdate } from "vue-router";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import Modal from "@/components/Modal.vue";
import SmallModal from "@/components/SmallModal.vue";
import NavigationModal from "@/components/NavigationModal.vue";
import { ElMessage, ElMessageBox, ElScrollbar } from "element-plus";
import { ElNotification } from "element-plus";
export default {
  name: "ContactChat",
  components: {
    Modal,
    SmallModal,
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
      userSessionList: [],
      groupSessionList: [],
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
      store.state.socket.onmessage = (jsonMessage) => {
        const message = JSON.parse(jsonMessage.data);
        if (message.type != 3) {
          if (
            // 群聊过来的消息，且当前会话是该群聊
            (message.receive_id[0] == "G" &&
              message.receive_id == data.contactInfo.contact_id) ||
            // 其他用户过来的消息，且当前会话是该用户
            (message.receive_id[0] == "U" &&
              message.receive_id == data.userInfo.uuid) ||
            // 自己发送的消息
            message.send_id == data.userInfo.uuid
          ) {
            console.log("收到消息：", message);
            if (data.messageList == null) {
              data.messageList = [];
            }
            data.messageList.push(message);
            scrollToBottom();
          }
          // 其他接受的消息都不显示在messageList中，而是通过切换页面或刷新页面getMessageList来获取
        } else {
          var messageAVdata = JSON.parse(message.av_data); // 后端message的该字段命名为av_data
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
          console.log("收到消息：", message);
          if (data.messageList == null) {
            data.messageList = [];
          }
          data.messageList.push(message);
          scrollToBottom();
        }
      };
      scrollToBottom();
      next();
    });
    // 这是刚渲染/chat/:id页面的时候会调用
    onMounted(async () => {
      try {
        /*  */
        console.log(router.currentRoute.value.params.id);
        await getChatContactInfo(router.currentRoute.value.params.id);
        await getSessionId(router.currentRoute.value.params.id);
        console.log(data.contactInfo);
        if (data.contactInfo.contact_id[0] == "U") {
          await getMessageList();
        } else {
          await getGroupMessageList();
        }
        console.log(data.sessionId);
        store.state.socket.onmessage = (jsonMessage) => {
          const message = JSON.parse(jsonMessage.data);
          if (message.type != 3) {
            if (
              // 群聊过来的消息，且当前会话是该群聊
              (message.receive_id[0] == "G" &&
                message.receive_id == data.contactInfo.contact_id) ||
              // 其他用户过来的消息，且当前会话是该用户
              (message.receive_id[0] == "U" &&
                message.receive_id == data.userInfo.uuid) ||
              // 自己发送的消息
              message.send_id == data.userInfo.uuid
            ) {
              console.log("收到消息：", message);
              if (data.messageList == null) {
                data.messageList = [];
              }
              data.messageList.push(message);
              scrollToBottom();
            }
            // 其他接受的消息都不显示在messageList中，而是通过切换页面或刷新页面getMessageList来获取
          } else {
            var messageAVdata = JSON.parse(message.av_data); // 后端message的该字段命名为av_data
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
              } else if (messageAVdata.type === "reject_call") {
                endCall();
              } else if (messageAVdata.type === "receive_call") {
                console.log("接收到receive_call消息", data.userInfo.nickname);
                createOffer();
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
            console.log("收到消息：", message);
            if (data.messageList == null) {
              data.messageList = [];
            }
            data.messageList.push(message);
            scrollToBottom();
          }
        };
        scrollToBottom();
      } catch (error) {
        console.error(error);
      }
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

    

    const handleToChatUser = async (user) => {
      router.push("/chat/" + user.user_id);
    };

    const handleToChatGroup = async (group) => {
      router.push("/chat/" + group.group_id);
    };

    const handleShowUserSessionList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const userSessionListRsp = await axios.post(
          store.state.backendUrl + "/session/getUserSessionList",
          data.ownListReq
        );
        if (userSessionListRsp.data.data) {
          for (let i = 0; i < userSessionListRsp.data.data.length; i++) {
            if (!userSessionListRsp.data.data[i].avatar.startsWith("http")) {
              userSessionListRsp.data.data[i].avatar =
                store.state.backendUrl + userSessionListRsp.data.data[i].avatar;
            }
          }
        }
        data.userSessionList = userSessionListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideUserSessionList = () => {
      data.userSessionList = [];
    };
    const handleShowGroupSessionList = async () => {
      try {
        data.ownListReq.owner_id = data.userInfo.uuid;
        const groupSessionListRsp = await axios.post(
          store.state.backendUrl + "/session/getGroupSessionList",
          data.ownListReq
        );
        if (groupSessionListRsp.data.data) {
          for (let i = 0; i < groupSessionListRsp.data.data.length; i++) {
            if (!groupSessionListRsp.data.data[i].avatar.startsWith("http")) {
              groupSessionListRsp.data.data[i].avatar =
                store.state.backendUrl +
                groupSessionListRsp.data.data[i].avatar;
            }
          }
        }
        data.groupSessionList = groupSessionListRsp.data.data;
      } catch (error) {
        console.error(error);
      }
    };
    const handleHideGroupSessionList = () => {
      data.groupSessionList = [];
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
    const sendMessage = () => {
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
      store.state.socket.send(JSON.stringify(chatMessageRequest));
      data.chatMessage = "";
      scrollToBottom();
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

    const getMessageList = async () => {
      try {
        console.log(data.contactInfo);
        const req = {
          user_one_id: data.userInfo.uuid,
          user_two_id: data.contactInfo.contact_id,
        };
        console.log(req);
        const rsp = await axios.post(
          store.state.backendUrl + "/message/getMessageList",
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
        const scrollHeight = data.innerRef.scrollHeight;
        console.log(scrollHeight);
        data.scrollbarRef.setScrollTop(scrollHeight);
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
      handleToChatUser,
      handleToChatGroup,
      handleShowUserSessionList,
      handleShowGroupSessionList,
      handleHideUserSessionList,
      handleHideGroupSessionList,
      deleteSession,
      preToDeleteSession,
      preToDeleteContact,
      preToBlackContact,
      sendMessage,
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

.el-menu {
  background-color: rgb(252, 210.9, 210.9);
  width: 100%;
}

.el-menu-item {
  background-color: rgb(255, 255, 255);
  height: 45px;
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
  background-color: rgb(239, 255, 174);
  color: rgb(74, 72, 72);
  display: inline-block;
  max-width: 400px;
  white-space: normal; /* 允许文本换行 */
  font-family: Arial, Helvetica, sans-serif;
  border-radius: 6px;
  padding: 3px;
  padding-right: 5px;
  font-size: 14px;
  box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.2);
}

.right-message-content {
  background-color: rgb(252, 210.9, 210.9);
  color: rgb(74, 72, 72);
  display: inline-block;
  max-width: 400px;
  white-space: normal; /* 允许文本换行 */
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
  background-color: rgb(252, 210.9, 210.9);
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
  background-color: rgb(252, 210.9, 210.9);
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
  background-color: rgb(252, 210.9, 210.9);
}
</style>
