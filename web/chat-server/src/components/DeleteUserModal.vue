<template>
  <div style="height: 100%; width: 100%" v-if="isVisible">
    <el-table
      :data="deleteUserTableData"
      style="width: 100%; height: 90%"
      @selection-change="selectUsers"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="uuid" label="Uuid" width="200" />
      <el-table-column
        prop="nickname"
        label="昵称"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column
        prop="telephone"
        label="电话"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column prop="is_admin" label="管理员" width="120" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_admin == 0"
            >否</el-button
          >
          <el-button type="primary" v-if="scope.row.is_admin == 1"
            >是</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="删除状态" width="120">
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_deleted == false"
            >正常</el-button
          >
          <el-button type="danger" v-if="scope.row.is_deleted == true"
            >已删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <div class="main-able-button-container">
      <el-button
        style="background-color: rgb(252, 210.9, 210.9);" @click="deleteUsers"
        >删除/全部删除</el-button
      >
    </div>
  </div>
</template>

<script>
import { onMounted, reactive, toRefs } from 'vue';
import { useStore } from 'vuex';
import axios from '@/utils/axios';
export default {
  name: "DeleteUserModal",
  props: {
    isVisible: false,
  },
  setup() {
    const store = useStore();
    const data = reactive({
      deleteUserTableData: [],
    });

    onMounted(() => {
      getUserList();
    });

    const getUserList = async () => {
      try {
        const req = {
          owner_id: store.state.userInfo.uuid,
        }
        const rsp = await axios.post(
          store.state.backendUrl + "/user/getUserInfoList", req
        );
        data.deleteUserTableData = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.log(error);
      }
    };

    const selectUsers = (val) => {
      data.uuidList = val.map((item) => item.uuid);
      console.log(data.uuidList);
    };

    const deleteUsers = async () => {
      try {
        const req = {
          uuid_list: data.uuidList,
        };
        const rsp = await axios.post(store.state.backendUrl + "/user/deleteUsers", req);
        console.log(rsp);
        getUserList();
      } catch (error) {
        console.log(error);
      }
    };

    return {
      ...toRefs(data),
      getUserList,
      selectUsers,
      deleteUsers,
    }
  }
};
</script>

<style scoped>
.main-able-button-container {
  height: 10%;
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
}
</style>