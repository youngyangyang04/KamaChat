<template>
  <div style="height: 100%; width: 100%" v-if="isVisible">
    <el-table
      :data="setAdminTableData"
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
      <el-button style="background-color: rgb(252, 210.9, 210.9); margin-left: 20px;" @click="setAdmin(1)"
        >设置/全部设置为管理员</el-button
      >
      <el-button style="background-color: rgb(252, 210.9, 210.9)" @click="setAdmin(0)"
        >取消/全部取消管理员</el-button
      >
    </div>
  </div>
</template>

<script>
import { ElTable } from "element-plus";
import { onMounted, reactive, toRefs } from "vue";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { useRouter } from 'vue-router';
export default {
  name: "SetAdminModal",
  props: {
    isVisible: false,
  },
  setup() {
    const store = useStore();
    const router = useRouter();
    const data = reactive({
      setAdminTableData: [],
      uuidList: [],
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
        data.setAdminTableData = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.log(error);
      }
    };

    const selectUsers = (val) => {
      data.uuidList = val.map((item) => item.uuid);
      console.log(data.uuidList);
    };

    const setAdmin = async (isAdmin) => {
      try {
        const req = {
          uuid_list: data.uuidList,
          is_admin: isAdmin,
        }
        const rsp = await axios.post(
          store.state.backendUrl + "/user/setAdmin", req);
        console.log(rsp);
        // router.go(0);
        getUserList();
      } catch (error) {
        console.log(error);
      }
    };


    return {
      ...toRefs(data),
      router,
      getUserList,
      selectUsers,
      setAdmin,
    };
  },
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