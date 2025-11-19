<template>
  <div style="height: 100%; width: 100%" v-if="isVisible">
    <el-table
      :data="disableGroupTableData"
      style="width: 100%; height: 90%"
      @selection-change="selectGroups"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="uuid" label="Uuid" width="200" />
      <el-table-column
        prop="name"
        label="群名称"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column
        prop="owner_id"
        label="群主id"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column prop="status" label="禁用状态" width="120" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.status == 0"
            >否</el-button
          >
          <el-button type="primary" v-if="scope.row.status == 1"
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
        style="background-color: rgb(252, 210.9, 210.9); margin-left: 20px;" @click="setGroupsStatus(1)"
        >禁用/全部禁用</el-button
      >
    <el-button
        style="background-color: rgb(252, 210.9, 210.9);" @click="setGroupsStatus(0)"
        >启用/全部启用</el-button
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
      disableGroupTableData: [],
    });

    onMounted(() => {
      getGroupInfoList();
    });

    const getGroupInfoList = async () => {
      try {
        const rsp = await axios.post(
          store.state.backendUrl + "/group/getGroupInfoList"
        );
        data.disableGroupTableData = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.log(error);
      }
    };

    const selectGroups = (val) => {
      data.uuidList = val.map((item) => item.uuid);
      console.log(data.uuidList);
    };

    const setGroupsStatus = async (status) => {
      try {
        const req = {
          uuid_list: data.uuidList,
          status: status,
        };
        const rsp = await axios.post(store.state.backendUrl + "/group/setGroupsStatus", req);
        console.log(rsp);
        getGroupInfoList();
      } catch (error) {
        console.log(error);
      }
    };

    return {
      ...toRefs(data),
      getGroupInfoList,
      setGroupsStatus,
      selectGroups,
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