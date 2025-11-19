<template>
  <div style="height: 100%; width: 100%" v-if="isVisible">
    <el-table
      :data="deleteGroupTableData"
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
        width="200"
        show-overflow-tooltip
      />
      <el-table-column prop="status" label="禁用状态" width="80" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.status == 0"
            >否</el-button
          >
          <el-button type="primary" v-if="scope.row.status == 1"
            >是</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="删除状态" width="100">
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
        style="background-color: rgb(252, 210.9, 210.9);" @click="deleteGroups"
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
      deleteGroupTableData: [],
      uuidList: [],
    });

    onMounted(() => {
      getGroupInfoList();
    });

    const getGroupInfoList = async () => {
      try {
        const rsp = await axios.post(
          store.state.backendUrl + "/group/getGroupInfoList"
        );
        data.deleteGroupTableData = rsp.data.data;
        console.log(rsp);
      } catch (error) {
        console.log(error);
      }
    };

    const selectGroups = (val) => {
      data.uuidList = val.map((item) => item.uuid);
      console.log(data.uuidList);
    };

    const deleteGroups = async () => {
      try {
        const req = {
          uuid_list: data.uuidList,
        };
        const rsp = await axios.post(store.state.backendUrl + "/group/deleteGroups", req);
        console.log(rsp);
        getGroupInfoList();
      } catch (error) {
        console.log(error);
      }
    };

    return {
      ...toRefs(data),
      getGroupInfoList,
      selectGroups,
      deleteGroups,
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