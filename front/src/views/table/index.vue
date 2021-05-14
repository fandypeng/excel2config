<template>
  <div class="app-container">
    <el-tabs type="border-card">
      <el-tab-pane label="配置列表">
        <div class="operate-box">
          <el-button :type="groupInfo.IsDev ? 'primary': 'danger'" @click="handleSwitchEnv">
            <span v-if="groupInfo.IsDev === true">当前处于测试环境，点击切换到正式环境</span>
            <span v-else>当前处于正式环境，点击切换到测试环境</span>
          </el-button>
          <el-button v-show="groupInfo.IsDev" type="primary" @click="handleAddExcel">添加表格</el-button>
        </div>
        <el-table
          v-loading="listLoading"
          :data="list"
          element-loading-text="Loading"
          border
          fit
          highlight-current-row
          style="width: 100%;margin-top:30px;"
        >
          <!-- <el-table-column align="center" label="ID" width="210">
            <template slot-scope="scope">
              {{ scope.row.id }}
            </template>
          </el-table-column> -->
          <el-table-column label="名称" width="150">
            <template slot-scope="scope">
                <el-link :href="'#/config/excel/' + scope.row.name + '/' + scope.row.id"  type="primary">{{ scope.row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column label="创建者"  align="center" width="150">
            <template slot-scope="scope">
              <span>{{ scope.row.owner }}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注"  align="center">
            <template slot-scope="scope">
              {{ scope.row.remark }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="210" align="center">
            <template slot-scope="scope">
              {{ scope.row.createTime | dateFilter }}
            </template>
          </el-table-column>
          <el-table-column label="编辑时间" width="210" align="center">
            <template slot-scope="scope">
              {{ scope.row.editTime | dateFilter }}
            </template>
          </el-table-column>
          <el-table-column class-name="status-col" label="发布" width="300" align="center">
            <template slot-scope="scope">
              <el-button type="success" size="small" @click="handleExport(scope)">发布</el-button>
              <el-button type="info" size="small" @click="handleExportRecord(scope)">记录</el-button>
            </template>
          </el-table-column>
          <el-table-column class-name="status-col" label="同步正式环境" width="300" align="center">
            <template slot-scope="scope">
              <el-button v-show="groupInfo.IsDev" type="warning" size="small" @click="handleExportProd(scope)">同步正式环境</el-button>
            </template>
          </el-table-column>
          <el-table-column class-name="status-col" label="管理" width="300" align="center">
            <template slot-scope="scope">
              <el-button type="primary" size="small" @click="handleEdit(scope)">编辑</el-button>
              <el-button type="danger" size="small" @click="handleDelete(scope)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-dialog :visible.sync="dialogVisible" :title="dialogType==='edit'?'编辑Excel信息':'新增Excel'">
          <el-form :model="excel" label-width="80px" label-position="left">
            <el-form-item label="名称">
              <el-input v-model="excel.name" placeholder="名称" :disabled="dialogType==='edit'" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input
                v-model="excel.remark"
                :autosize="{ minRows: 2, maxRows: 4}"
                type="textarea"
                placeholder="备注"
              />
            </el-form-item>
            <el-form-item label="创建者" v-show="dialogType==='edit'">
              <el-input v-model="excel.owner" placeholder="创建者" disabled />
            </el-form-item>
          </el-form>
          <div style="text-align:right;">
            <el-button type="danger" @click="dialogVisible=false">取消</el-button>
            <el-button type="primary" @click="confirmExcel">保存</el-button>
          </div>
        </el-dialog>

        <el-dialog :visible.sync="exportDialogVisible" title="发布表格内容" width="80%">
          <div v-loading="dialogLoading">
            <div style="float: left;width: 50%; padding:20px;">
              <el-form label-width="120px" label-position="left">
                <el-form-item label="选择表格">
                  <el-radio-group v-model="sleectedSheet" @change="onSelectedSheetChange(1)">
                    <el-radio v-for="item in sheetList" :key="item" :label="item" >{{item}} </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item label="发布内容">
                  <el-input
                    type="textarea"
                    autosize
                    placeholder="请输入发布内容"
                    v-model="exportRemark">
                  </el-input>
                </el-form-item>
                <el-form-item label="导出">
                  <el-button type="primary" @click="exportJsonFile">导出JSON</el-button>
                </el-form-item>
              </el-form>
            </div>
            <div style="float: left;width: 50%;padding:20px;">
              <el-form label-width="120px" label-position="left">
                <el-form-item label="对比设置">
                </el-form-item>
                <el-form-item label="差异显示行数">
                  <el-input type="number" placeholder="10" v-model.number="contextRowLength"></el-input>
                </el-form-item>
              </el-form>
            </div>
            <div style="width: 100%;text-align: center;clear: both;">
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">修改前版本</p>
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">修改后版本</p>
            </div>
            <div style="margin-top: 80px;">
              <code-diff :old-string="codeDiff.beforeExportStr" :new-string="codeDiff.exportStr" :isShowNoChange="true" :renderNothingWhenEmpty="false" outputFormat="side-by-side"  :context="contextRowLength" />
            </div>
          </div>
          <div style="text-align:right;">
            <el-button type="danger" @click="exportDialogVisible=false">取消</el-button>
            <el-button type="primary" @click="confirmExport">发布</el-button>
          </div>
        </el-dialog>

        <el-dialog :visible.sync="exportProdDialogVisible" title="同步表格内容到正式环境" width="80%">
          <div v-loading="dialogLoading">
            <div style="float: left;width: 50%;padding:20px;">
              <el-form :model="excel" label-width="120px" label-position="left">
                <el-form-item label="选择表格">
                  <el-radio-group v-model="sleectedSheet" @change="onSelectedSheetChange(2)">
                    <el-radio v-for="item in sheetList" :key="item" :label="item" >{{item}} </el-radio>
                  </el-radio-group>
                </el-form-item>
              </el-form>
            </div>
            <div style="float: left;width: 50%;padding:20px;">
              <el-form label-width="120px" label-position="left">
                <el-form-item label="对比设置">
                </el-form-item>
                <el-form-item label="差异显示行数">
                  <el-input type="number" placeholder="10" v-model.number="contextRowLength"></el-input>
                </el-form-item>
              </el-form>
            </div>
            <div style="width: 100%;text-align: center;clear: both;">
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">正式环境版本</p>
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">测试环境版本</p>
            </div>
            <div style="margin-top: 80px;">
              <code-diff :old-string="codeDiff.beforeExportProdStr" :new-string="codeDiff.exportProdStr" :isShowNoChange="true" :renderNothingWhenEmpty="false" outputFormat="side-by-side" :context="contextRowLength" />
            </div>
          </div>
          <div style="text-align:right;">
            <el-button type="danger" @click="exportProdDialogVisible=false">取消</el-button>
            <el-button type="primary" @click="confirmExportProd">同步</el-button>
          </div>
        </el-dialog>

        <el-dialog :visible.sync="exportRecordDialogVisible" title="发布记录" width="80%">
          <div v-loading="dialogLoading">
            <div style="float: left;width: 50%; padding:20px;">
              <el-form :model="excel" label-width="120px" label-position="left">
                <el-form-item label="选择表格">
                  <el-radio-group v-model="sleectedSheet" @change="onSelectedSheetChange(3)">
                    <el-radio v-for="item in sheetList" :key="item" :label="item" >{{item}} </el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item label="发布记录">
                  <el-select v-model="selectRecordId" @change="onSelectRecordChange()" filterable placeholder="选择历史发布记录" style="width:90%;">
                    <el-option
                      v-for="item in recordList"
                      :key="item.id"
                      :label="item.remark"
                      :value="item.id">
                      <span style="float: left;margin-right:100px;">{{ item.remark }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">{{ item.userName }} 发布于 {{ item.time }}</span>
                    </el-option>
                  </el-select>
                  <el-button style="margin-top: 20px;" type="danger" @click="handleRollBack">回退到此版本</el-button>
                </el-form-item>
              </el-form>
            </div>
            <div style="float: left;width: 50%;padding:20px;">
              <el-form label-width="120px" label-position="left">
                <el-form-item label="对比设置">
                </el-form-item>
                <el-form-item label="差异显示行数">
                  <el-input type="number" placeholder="10" v-model.number="contextRowLength"></el-input>
                </el-form-item>
              </el-form>
            </div>
            <div style="width: 100%;text-align: center; clear: both;">
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">发布前版本</p>
              <p style="font-size: 18px;font-weight: bold;width: 50%;float: left;">发布后版本</p>
            </div>
            <div style="margin-top: 80px;">
              <code-diff :old-string="codeDiff.beforeExportRecordStr" :new-string="codeDiff.exportRecordStr" :isShowNoChange="true" :renderNothingWhenEmpty="false" outputFormat="side-by-side" :context="contextRowLength" />
            </div>
          </div>
        </el-dialog>

      </el-tab-pane>


      <el-tab-pane label="项目配置">
        <el-form ref="form" :model="groupInfo" label-width="120px">
          <el-form-item label="名称">
            <el-input v-model="groupInfo.name" value="贪吃蛇大作战" />
          </el-form-item>
          <el-form-item label="头像">
            <el-col :span="11">
              <img id="group-icon" v-bind:src="groupInfo.avatar" alt="">
              <input class="avatar-input" type="file" accept="image/*" @change="uploadAvatar" name="avatar"   />
            </el-col>
          </el-form-item>
          <br>
          <el-form-item label="备注">
            <el-input v-model="groupInfo.remark" type="textarea" />
          </el-form-item>
          <el-form-item label="授权访问Token">
            <div class="el-col-8">
              <el-input v-model="groupInfo.AccessToken" type="text" />
            </div>
            <div class="el-col-4" style="margin-left: 20px;">
              <el-button type="primary" @click="generateAccessToken">刷新AccessToken</el-button>
            </div>
            <br>
            <el-tag type="danger" >
              外部应用可以通过授权访问token调用Api获取Excel内容，具体内容请参见API文档描述
            </el-tag>
          </el-form-item>
          <el-form-item label="数据仓库">
            <el-checkbox-group v-model="groupInfo.store" style="padding: 5px 15px;">
              <el-checkbox v-for="item in storeOptions" :key="item.key" :label="item.key">
                {{ item.name }}
              </el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item v-show="inArray(groupInfo.store, 1)" label="Redis配置" class="store-group">
            <el-form-item label="RedisDSN">
              <el-input v-model="groupInfo.RedisDSN" name="RedisDSN" type="text" placeholder="172.2.1.88:6379" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="groupInfo.RedisPassword" name="RedisPassword" auto-complete="new-password" type="password" placeholder="password" />
            </el-form-item>
            <el-form-item label="KeyPrefix">
              <el-input v-model="groupInfo.RedisKeyPrefix" type="text" placeholder="key prefix 可不填" />
            </el-form-item>
            <el-tag type="danger" style="margin-left: 40px;" >
              应用可以通过订阅redis的channel：{{groupInfo.RedisKeyPrefix}}config_refresh 获取配置更新的通知，message为更新的表名
            </el-tag>
            <el-form-item label=" ">
              <el-button type="success" @click="testRedisConnection">测试连通性</el-button>
            </el-form-item>
          </el-form-item>

          <el-form-item v-show="inArray(groupInfo.store, 2)" label="Mysql配置" class="store-group">
            <el-form-item label="MysqlDSN">
              <el-input v-model="groupInfo.MysqlDSN" type="text" placeholder="username:password@tcp(172.2.1.88:3306)/dbname?charset=utf8mb4" />
            </el-form-item>
            <el-form-item label=" ">
              <el-button type="success" @click="testMysqlConnection">测试连通性</el-button>
            </el-form-item>
          </el-form-item>

          <el-form-item v-show="inArray(groupInfo.store, 3)" label="Mongodb配置" class="store-group">
            <el-form-item label="MongodbDSN">
              <el-input v-model="groupInfo.MongodbDSN" type="text" placeholder="mongodb://username:password@127.0.0.1:27017/?authSource=dbname" />
            </el-form-item>
            <el-form-item label=" ">
              <el-button type="success" @click="testMongodbConnection">测试连通性</el-button>
            </el-form-item>
          </el-form-item>

          <el-form-item v-show="inArray(groupInfo.store, 4)" label="Databus配置" class="store-group">
            <el-form-item label="GrpcDSN">
              <el-input v-model="groupInfo.GrpcDSN" name="GrpcDSN" type="text" placeholder="172.2.1.88:10000" />
            </el-form-item>
            <el-form-item label="GrpcAppKey">
              <el-input v-model="groupInfo.GrpcAppKey" name="GrpcAppKey" type="text" placeholder="appKey" />
            </el-form-item>
            <el-form-item label="GrpcAppSecret">
              <el-input v-model="groupInfo.GrpcAppSecret" name="GrpcAppSecret" type="text" placeholder="appSecret" />
            </el-form-item>
            <el-tag type="danger" style="margin-left: 40px;" >
              请参阅Readme.md文档搭建e2cdatabus服务器
            </el-tag>
            <el-form-item label=" ">
              <el-button type="primary" @click="generateAppKeySecretHandler">获取AppKey和AppSecret</el-button>
              <el-button type="success" @click="testGRpcConnection">测试连通性</el-button>
            </el-form-item>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="saveGroupConfig">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="项目成员">
        <el-button type="primary" @click="handleAddMember">添加成员</el-button>
        <el-table
          v-loading="listLoading"
          :data="groupInfo.members"
          element-loading-text="Loading"
          border
          fit
          highlight-current-row
          style="width: 100%;margin-top:30px;"
        >
        <!-- <el-table-column label="头像"  align="center" >
            <template slot-scope="scope">
              <img :src="scope.row.avatar" style="width:60px;border-radius:50%;" alt="">
            </template>
          </el-table-column> -->
          <el-table-column label="昵称"  align="center" >
            <template slot-scope="scope">
              <span>{{ scope.row.userName }}</span>
            </template>
          </el-table-column>

          <el-table-column label="角色"  align="center" >
            <template slot-scope="scope">
              <el-select v-show="uid !== scope.row.uid" v-model.number="scope.row.role" @change="onChangeUserRole(scope)" placeholder="用户角色" style="width:90%;">
                <el-option
                  v-for="item in roleNames"
                  :key="item.role"
                  :label="item.name"
                  :value="item.role">
                </el-option>
              </el-select>
              <span v-show="uid === scope.row.uid">{{ scope.row.role | roleFilter }}</span>
            </template>
          </el-table-column>

          <el-table-column class-name="status-col" label="操作" width="300" align="center">
            <template slot-scope="scope">
              <el-button v-show="showDeleteBtn(scope)" type="danger" size="small" @click="handleDeleteMember(scope)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-dialog :visible.sync="addMemberdialogVisible" title="添加成员" width="500px">
          <el-form label-width="80px" label-position="left">
            <el-form-item label="名称">
              <el-autocomplete
                v-model="searchUserName"
                :fetch-suggestions="querySearchAsync"
                placeholder="请输入名称:"
                clearable
                :trigger-on-focus="false"
                @select="handleSelect"
              ></el-autocomplete>
            </el-form-item>
          </el-form>
          <br>
          <div style="text-align:right;">
            <el-button type="danger" @click="addMemberdialogVisible=false">取消</el-button>
            <el-button type="primary" @click="addMember()">保存</el-button>
          </div>
        </el-dialog>
      </el-tab-pane>
    </el-tabs>

  </div>
</template>

<script>
import { getList, deleteExcel, updateExcel, addExcel, exportExcel, getSheetList, exportExcelProd } from '@/api/sheet'
import { searchUser } from '@/api/user'
import { groupUpdate, groupTestConnection, getConfigFromDB, exportConfigToDB, generateAppKeySecret,
  groupList, syncToProd, exportRecord, exportRecordContent, exportRollback } from '@/api/group'
import { mapGetters } from 'vuex'
import { deepClone } from '@/utils'
import CodeDiff from 'vue-code-diff'
import { MessageBox } from 'element-ui'

const defaultExcel = {
  id: '',
  name: '',
  owner: '',
  remark: '',
  create_time: '',
  update_time: ''
}

const Role = {
  Admin: 1,
  Developer: 2
}
const roleConfig = [
  { 'role': 1, 'name': '管理员' },
  { 'role': 2, 'name': '开发者' }
]

export default {
  components: { CodeDiff },
  computed: {
    ...mapGetters([
      'uid',
      'name'
    ])
  },
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    },
    dateFilter(datetime) {
      return new Date(parseInt(datetime) * 1000).toLocaleString().replace(/:\d{1,2}$/,' ');
    },
    roleFilter(role) {
      if (role === undefined || role === null) {
        role = Role.Developer
      }
      if (roleConfig.length >= role) {
        return roleConfig[ role - 1 ].name
      }
      return '未定义'
    }
  },
  data() {
    return {
      excel: Object.assign({}, defaultExcel),
      list: [],
      listLoading: true,
      dialogLoading: false,
      dialogVisible: false,
      dialogType: 'new',
      exportDialogVisible: false,
      exportFormat: '',
      form: {},
      groupInfo: {
        name: '贪吃蛇大作战',
        remark: '这是贪吃蛇大作战App项目',
        avatar: require('../../icons/tanchishe-icon.png'),
        store: []
      },
      storeOptions: [
        { key: 1, name: 'Redis' },
        { key: 2, name: 'Mysql' },
        // {key: 'Mongodb', name: 'Mongodb'},
        { key: 4, name: 'Databus' }
      ],
      addMemberdialogVisible: false,
      searchUserName: '',
      searchUid: '',
      compareVisable: false,
      sheetList: [],
      sleectedSheet: '',
      exportProdDialogVisible: false,
      exportRecordDialogVisible: false,
      selectRecordId: '',
      recordList: [],
      exportRemark: '',
      contextRowLength: 10,
      codeDiff: {
        beforeExportStr: '',
        exportStr: '',
        beforeExportProdStr: '',
        exportProdStr: '',
        beforeExportRecordStr: '',
        exportRecordStr: ''
      },
      roleNames: roleConfig
    }
  },
  created() {
    this.excel.owner = this.name
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      var groupInfo = localStorage.getItem('groupInfo')
      groupInfo = JSON.parse(groupInfo)
      console.log(groupInfo)
      if (groupInfo === null) {
        this.$router.push({
          path: '/'
        })
        return
      }
      if (groupInfo.store === undefined || groupInfo.store === null) {
        groupInfo.store = []
      }
      this.groupInfo = groupInfo
      this.groupInfo.members.forEach(m => {
        if (m.uid === this.uid && (m.role === 0 || m.role === undefined)) {
          m.role = Role.Admin
        }
      })
      getList({limit: 10, group_id: groupInfo.gid}).then(response => {
        if (response.data.list !== undefined) {
          this.list = response.data.list
        }
        this.listLoading = false
      })
    },
    handleAddExcel() {
      this.excel = Object.assign({}, defaultExcel),
      this.excel.owner = this.name
      this.dialogType = 'new'
      this.dialogVisible = true
    },
    handleSwitchEnv() {
      var text = '当前处于正式环境，是否确认切换到测试环境'
      if (this.groupInfo.IsDev) {
        text = '当前处于测试环境，是否确认切换到正式环境'
      }
      MessageBox.confirm(text, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        if (res === undefined) {
          return
        }
        var unionGid = this.groupInfo.UnionGroupId
        groupList(unionGid).then(response => {
          if (response.data.list !== undefined && response.data.list.length > 0) {
            this.groupInfo = response.data.list[0]
            localStorage.setItem('groupInfo', JSON.stringify(this.groupInfo))
            // refresh
            var path = this.$route.path
            this.$router.push({
              path,
              query: {
                t: +new Date() //保证每次点击路由的query项都是不一样的，确保会重新刷新view
              }
            })
            this.fetchData()
          } else{
            this.$message({
              type: 'warning',
              message: '暂无权限，请联系管理员添加权限'
            })
          }
        })
      }).catch(e => {
        console.log(e)
      })
    },
    async confirmExcel() {
      if (this.excel.name.length === 0) {
        this.$message({
          type: 'warning',
          message: '请输入名称'
        })
        return
      }
      const isEdit = this.dialogType === 'edit'
      if (isEdit) {
        await updateExcel(this.excel)
        for (let index = 0; index < this.list.length; index++) {
          if (this.list[index].id === this.excel.id) {
            this.list.splice(index, 1, Object.assign({}, this.excel))
            break
          }
        }
      } else {
        var params = this.excel
        params.uid = this.uid
        params.group_id = this.groupInfo.gid
        const { data } = await addExcel(params)
        this.excel.id = data.eid
        var timestamp=new Date().getTime()
        this.excel.createTime = timestamp/1000
        this.excel.editTime = timestamp/1000
        this.list.push(this.excel)
      }

      const { id, name, remark } = this.excel
      this.dialogVisible = false
      this.$notify({
        title: 'Success',
        dangerouslyUseHTMLString: true,
        message: `
            <div>名称: ${name}</div>
            <div>备注: ${remark}</div>
          `,
        type: 'success'
      })
    },
    handleEdit(scope) {
      this.dialogType = 'edit'
      this.dialogVisible = true
      this.checkStrictly = true
      this.excel = deepClone(scope.row)
    },
    handleDelete({$index, row}) {
      this.$confirm('确定删除表格吗?', 'Warning', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        await deleteExcel({id: row.id, name: row.name})
        this.list.splice($index, 1)
        this.$message({
          type: 'success',
          message: '删除成功!'
        })
      }).catch(err => { console.error(err) })
    },
    handleExport(scope) {
      this.excel = deepClone(scope.row)
      this.exportDialogVisible = true
      this.sheetList = []
      this.codeDiff.beforeExportStr = ''
      this.codeDiff.exportStr = ''
      this.exportRemark = ''
      this.dialogLoading = true
      getSheetList({gridKey: this.excel.id}).then(response => {
        if (response.code === 0) {
          this.sheetList = response.data.sheetName
          console.log(this.sheetList)
          if (this.sheetList.length > 0) {
            this.sleectedSheet = this.sheetList[0]
          }
          this.getExcelSheetContent(this.excel.id, this.sleectedSheet)
          this.getExcelSheetFromDB(this.excel.id, this.sleectedSheet)
        }
      })
    },
    handleExportProd(scope) {
      this.excel = deepClone(scope.row)
      this.exportProdDialogVisible = true
      this.sheetList = []
      this.codeDiff.beforeExportProdStr = ''
      this.codeDiff.exportProdStr = ''
      this.dialogLoading = true
      getSheetList({gridKey: this.excel.id}).then(response => {
        if (response.code === 0) {
          this.sheetList = response.data.sheetName
          console.log(this.sheetList)
          if (this.sheetList.length > 0) {
            this.sleectedSheet = this.sheetList[0]
          }
          this.getExcelSheetContent(this.excel.id, this.sleectedSheet)
          this.getExcelSheetFromProd(this.excel.id, this.sleectedSheet)
        }
      })
    },
    handleExportRecord(scope) {
      this.excel = deepClone(scope.row)
      this.exportRecordDialogVisible = true
      this.sheetList = []
      this.exportRecord = []
      this.codeDiff.beforeExportRecordStr = ''
      this.codeDiff.exportRecordStr = ''
      this.dialogLoading = true
      getSheetList({gridKey: this.excel.id}).then(response => {
        if (response.code === 0) {
          this.sheetList = response.data.sheetName
          console.log(this.sheetList)
          if (this.sheetList.length > 0) {
            this.sleectedSheet = this.sheetList[0]
          }
          // 获取record list
          this.getRecordList()
        }
      })
    },
    getRecordList() {
      exportRecord({gridKey: this.excel.id, sheetName: this.sleectedSheet}).then(resp => {
        if (resp.data !== null) {
          this.recordList = resp.data.list
        }
        if (this.recordList.length > 0) {
          this.selectRecordId = this.recordList[0].id
        }
        this.refreshRecordContent()
      })
    },
    refreshRecordContent() {
      var lastRecordId = ''
      // 判断是否有上一个record
      for (var i=0;i<this.recordList.length;i++) {
        if (this.recordList[i].id === this.selectRecordId) {
          if (i+1<this.recordList.length) {
            console.log('record：', this.recordList[i+1])
            lastRecordId = this.recordList[i+1].id
          }
        }
      }
      this.getExcelSheetFromRecord(this.excel.id, this.sleectedSheet, this.selectRecordId)
      this.getExcelSheetFromRecord(this.excel.id, this.sleectedSheet, lastRecordId)
    },
    async getExcelSheetContent(gridKey, sheetName) {
      this.codeDiff.exportStr = ''
      var params = {}
      params.gridKey = gridKey
      params.sheetName = sheetName
      params.format = 'json'
      const { data } = await exportExcel(params)
      if (data.jsonstr !== undefined && data.jsonstr.length > 0) {
        this.codeDiff.exportStr = data.jsonstr
        this.codeDiff.exportProdStr = data.jsonstr
        this.dialogLoading = false
      }
    },
    async getExcelSheetFromDB(gridKey, sheetName) {
      this.codeDiff.beforeExportStr = ''
      var params = {}
      params.gridKey = gridKey
      params.sheetName = sheetName
      const { data } = await getConfigFromDB(params)
      if (data.jsonstr !== undefined && data.jsonstr.length > 0) {
        this.codeDiff.beforeExportStr = data.jsonstr
        this.dialogLoading = false
      }
    },
    async getExcelSheetFromProd(gridKey, sheetName) {
      this.codeDiff.beforeExportProdStr = ''
      var params = {}
      params.gridKey = gridKey
      params.sheetName = sheetName
      params.format = 'json'
      params.gid = this.groupInfo.gid
      const { data } = await exportExcelProd(params)
      if (data.jsonstr !== undefined && data.jsonstr.length > 0) {
        this.codeDiff.beforeExportProdStr = data.jsonstr
        this.dialogLoading = false
      }
    },
    async getExcelSheetFromRecord(gridKey, sheetName, recordId) {
      if (recordId === this.selectRecordId) {
        this.codeDiff.exportRecordStr = ''
      } else {
        this.codeDiff.beforeExportRecordStr = ''
      }
      if (recordId.length === 0) {
        return
      }
      var params = {}
      params.gridKey = gridKey
      params.sheetName = sheetName
      params.recordId = recordId
      params.format = 'json'
      const { data } = await exportRecordContent(params)
      if (data.jsonstr !== undefined && data.jsonstr.length > 0) {
        if (recordId === this.selectRecordId) {
          this.codeDiff.exportRecordStr = data.jsonstr
          this.dialogLoading = false
        } else {
          this.codeDiff.beforeExportRecordStr = data.jsonstr
          this.dialogLoading = false
        }
      }
    },
    handleRollBack() {
      var text = '您当前选择的配置表是 \"' + this.sleectedSheet + '\", 确认要回退版本吗？'
      MessageBox.confirm(text, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        if (res === undefined) {
          return
        }
        var params = {}
        params.gridKey = this.excel.id
        params.sheetName = this.sleectedSheet
        params.recordId = this.selectRecordId
        exportRollback(params).then(response => {
          if (response.code !== 0) {
              this.$notify({
                title: '回退失败',
                dangerouslyUseHTMLString: true,
                message: response.data.msg,
                type: 'error'
              })
            } else {
              this.$notify({
                title: '回退成功',
                dangerouslyUseHTMLString: true,
                message: '',
                type: 'success'
              })
              this.exportRecordDialogVisible = false
            }
          }
        )
      }).catch(e => {
        console.log(e)
      })
    },
    onSelectedSheetChange(e) {
      this.dialogLoading = true
      if (e === 1) { // 发布弹窗
        this.listLoading = true
        this.getExcelSheetContent(this.excel.id, this.sleectedSheet)
        this.getExcelSheetFromDB(this.excel.id, this.sleectedSheet)
        this.listLoading = false
      } else if (e === 2) { // 同步弹窗
        this.getExcelSheetContent(this.excel.id, this.sleectedSheet)
        this.getExcelSheetFromProd(this.excel.id, this.sleectedSheet)
      } else if (e === 3) { //发布记录弹窗
        this.selectRecordId = ''
        this.codeDiff.beforeExportRecordStr = ''
        this.codeDiff.exportRecordStr = ''
        this.recordList = []
        this.getRecordList()
      }
    },
    onSelectRecordChange(e) {
      console.log('select record changed')
      this.refreshRecordContent()
    },
    confirmExport() {
      if (this.exportRemark.length === 0) {
        this.$notify({
          title: '发布失败',
          dangerouslyUseHTMLString: true,
          message: '请填写发布内容',
          type: 'error'
        })
        return
      }
      var text = '您当前选择的配置表是 \"' + this.sleectedSheet + '\", 修改备注为 \"'+this.exportRemark + '\", 确认发布吗？'
      if (!this.groupInfo.IsDev) {
        text = '当前处于正式环境, ' + text
      }
      MessageBox.confirm(text, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(res => {
        var params = {}
        params.gridKey = this.excel.id
        params.sheetName = this.sleectedSheet
        params.remark = this.exportRemark
        exportConfigToDB(params).then(response => {
          if (response.code === 0) {
            if (response.code !== 0) {
              this.$notify({
                title: '发布失败',
                dangerouslyUseHTMLString: true,
                message: response.data.msg,
                type: 'error'
              })
            } else {
              this.$notify({
                title: '发布成功',
                dangerouslyUseHTMLString: true,
                message: '',
                type: 'success'
              })
              this.exportDialogVisible = false
            }
          }
        })
      })
    },
    confirmExportProd() {
      var params = {}
      params.gridKey = this.excel.id
      params.sheetName = this.sleectedSheet
      params.gid = this.groupInfo.gid
      syncToProd(params).then(response => {
        if (response.code === 0) {
          if (response.code !== 0) {
            this.$notify({
              title: '同步失败',
              dangerouslyUseHTMLString: true,
              message: response.data.msg,
              type: 'error'
            })
          } else {
            this.$notify({
              title: '同步成功',
              dangerouslyUseHTMLString: true,
              message: '',
              type: 'success'
            })
            this.exportProdDialogVisible = false
          }
        }
      })
    },
    exportJsonFile() {
      this.downloadJsonFile(this.codeDiff.beforeExportStr, this.sleectedSheet + '.json')
    },
    downloadJsonFile(content, filename) {
      var eleLink = document.createElement('a')
      eleLink.download = filename
      eleLink.style.display = 'none'
      // 字符内容转变成blob地址
      var blob = new Blob([content])
      eleLink.href = URL.createObjectURL(blob)
      // 触发点击
      document.body.appendChild(eleLink)
      eleLink.click()
      // 然后移除
      document.body.removeChild(eleLink)
    },
    saveGroupConfig() {
      groupUpdate({id: this.groupInfo.gid, groupInfo: this.groupInfo}).then(response => {
        if (response.code !== 0) {
          this.$notify({
            title: '保存失败',
            dangerouslyUseHTMLString: true,
            message: response.data.msg,
            type: 'error'
          })
        } else {
          this.$notify({
            title: '保存成功',
            dangerouslyUseHTMLString: true,
            message: '',
            type: 'success'
          })
          localStorage.setItem('groupInfo', JSON.stringify(this.groupInfo))
        }
      })
    },
    uploadAvatar(e) {
      e.preventDefault()
      console.log(e)
      const files = e.target.files || e.dataTransfer.files
      if (this.checkFile(files[0])) {
        this.setSourceImg(files[0])
      }
    },
    /* ---------------------------------------------------------------*/
    // 检测选择的文件是否合适
    checkFile(file) {
      const { lang, maxSize } = this
      // 仅限图片
      if (file.type.indexOf('image') === -1) {
        this.hasError = true
        this.errorMsg = lang.error.onlyImg
        return false
      }
      // 超出大小
      if (file.size / 1024 > maxSize) {
        this.hasError = true
        this.errorMsg = lang.error.outOfSize + maxSize + 'kb'
        return false
      }
      return true
    },
    setSourceImg(file) {
      const fr = new FileReader()
      fr.onload = e => {
        this.groupInfo.avatar = fr.result
      }
      fr.readAsDataURL(file)
    },
    inArray(arr, target) {
      for (var i in arr) {
        if (arr[i] === target) {
          return true
        }
      }
      return false
    },
    testRedisConnection(e) {
      var params = {dsnType:1,dsn:this.groupInfo.RedisDSN,pwd:this.groupInfo.RedisPassword}
      groupTestConnection(params).then(response => {
        if (response.data.connected === 0 || response.data.connected === undefined) {
          this.$notify({
            title: 'Redis连接失败',
            dangerouslyUseHTMLString: true,
            message: `请检查DSN配置或者IP白名单配置`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: 'Redis连接成功',
            dangerouslyUseHTMLString: true,
            message: ``,
            type: 'success'
          })
        }
      })
    },
    testMysqlConnection() {
      var params = {dsnType:2,dsn:this.groupInfo.MysqlDSN}
      groupTestConnection(params).then(response => {
        if (response.data.connected === 0 || response.data.connected === undefined) {
          this.$notify({
            title: 'Mysql连接失败',
            dangerouslyUseHTMLString: true,
            message: `请检查DSN配置或者IP白名单配置`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: 'Mysql连接成功',
            dangerouslyUseHTMLString: true,
            message: ``,
            type: 'success'
          })
        }
      })
    },
    testGRpcConnection() {
      var params = {
        dsnType: 4,
        dsn:    this.groupInfo.GrpcDSN,
        appKey: this.groupInfo.GrpcAppKey,
        appSecret: this.groupInfo.GrpcAppSecret,
      }
      groupTestConnection(params).then(response => {
        if (response.data.connected === 0 || response.data.connected === undefined) {
          this.$notify({
            title: 'Databus连接失败',
            dangerouslyUseHTMLString: true,
            message: `请启动Databus服务，检查DSN和AppKey、AppSecret配置`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: 'Databus连接成功',
            dangerouslyUseHTMLString: true,
            message: ``,
            type: 'success'
          })
        }
      })
    },
    generateAppKeySecretHandler() {
      this.$set(this.groupInfo, 'GrpcDSN', '')
      generateAppKeySecret().then(response => {
        if (response.data === null || response.data.appKey === '' || response.data.appSecret === '') {
          this.$notify({
            title: '获取失败',
            dangerouslyUseHTMLString: true,
            message: `请联系管理员`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: '获取成功',
            dangerouslyUseHTMLString: true,
            message: `请参阅文档搭建e2cdatabus服务器`,
            type: 'success'
          })
          this.$set(this.groupInfo, 'GrpcAppKey', response.data.appKey)
          this.$set(this.groupInfo, 'GrpcAppSecret', response.data.appSecret)
        }
      })
    },
    generateAccessToken() {
      this.$set(this.groupInfo, 'AccessToken', '')
      generateAppKeySecret().then(response => {
        if (response.data === null || response.data.appKey === '' || response.data.appSecret === '') {
          this.$notify({
            title: '获取失败',
            dangerouslyUseHTMLString: true,
            message: `请联系管理员`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: '刷新成功',
            dangerouslyUseHTMLString: true,
            message: `保存后即可生效`,
            type: 'success'
          })
          this.$set(this.groupInfo, 'AccessToken', response.data.appSecret)
        }
      })
    },
    testMongodbConnection() {
      var params = { dsnType: 3, dsn: this.groupInfo.MongodbDSN }
      groupTestConnection(params).then(response => {
        if (response.data.connected === 0 || response.data.connected === undefined) {
          this.$notify({
            title: 'Mongodb连接失败',
            dangerouslyUseHTMLString: true,
            message: `请检查DSN配置或者IP白名单配置`,
            type: 'error'
          })
          return
        } else {
          this.$notify({
            title: 'Mongodb连接成功',
            dangerouslyUseHTMLString: true,
            message: ``,
            type: 'success'
          })
        }
      })
    },
    showDeleteBtn(scope) {
      if (this.uid === scope.row.uid) {
        return false
      }
      let isAdmin = false
      this.groupInfo.members.forEach(m => {
        if (m.uid === this.uid) {
          isAdmin = m.role === Role.Admin
        }
      })
      return isAdmin
    },
    handleAddMember() {
      this.addMemberdialogVisible = true
    },
    addMember() {
      var exist = false
      this.groupInfo.members.forEach(m => {
        if (m.uid === this.searchUid) {
          exist = true
          return
        }
      })
      if (exist) {
        this.$notify({
          title: '添加失败',
          dangerouslyUseHTMLString: true,
          message: `成员已存在`,
          type: 'error'
        })
        return
      }
      this.groupInfo.members.push({
        uid: this.searchUid,
        userName: this.searchUserName,
        role: Role.Developer
      })
      this.updateGroupInfo('保存成功', '添加成员：'+this.searchUserName)
      this.addMemberdialogVisible = false
    },
    onChangeUserRole(scope) {
      var role = scope.row.role
      this.groupInfo.members.forEach(m => {
        if (m.uid === this.searchUid) {
          if (m.role === role) {
            return
          }
          m.role = role
        }
      })
      this.updateGroupInfo('保存成功', '修改' + scope.row.userName + '为：' + this.roleNames[role - 1].name)
    },
    handleSelect(e) {
      this.searchUid = e.id
    },
    handleDeleteMember(scope) {
      var members = []
      if (this.uid === scope.row.uid) {
        this.$notify({
          title: '删除失败',
          message: '无法删除自己',
          type: 'error'
        })
      }
      this.groupInfo.members.forEach(m => {
        if (m.uid !== scope.row.uid) {
          members.push(m)
        }
      })
      this.groupInfo.members = members
      this.updateGroupInfo('保存成功', '删除成员：'+scope.row.userName)
    },
    updateGroupInfo(succTitle, sussMsg) {
      groupUpdate({id: this.groupInfo.gid, groupInfo: this.groupInfo}).then(response => {
        if (response.code !== 0) {
          this.$notify({
            title: '保存失败',
            dangerouslyUseHTMLString: true,
            message: response.data.msg,
            type: 'error'
          })
        } else {
          this.$notify({
            title: succTitle,
            dangerouslyUseHTMLString: true,
            message: sussMsg,
            type: 'success'
          })
          localStorage.setItem('groupInfo', JSON.stringify(this.groupInfo))
        }
      })
    },
    querySearchAsync(queryString, cb) {
      var results = []
      this.searchUid = ''
      if (queryString === '') {
        cb(results);
        return
      }
      searchUser({name: this.searchUserName}).then(response => {
        if (response.data.userInfos === undefined || response.data.userInfos === null) {
          cb(results);
          return
        }
        response.data.userInfos.forEach(u => {
          results.push({ 'id':u.uid, 'value': u.userName })
        })
        console.log(results)
        cb(results);
      })
    },
  }
}
</script>

<style scoped>
  .avatar-input {
    float: left;
    position: absolute;
    height: 60px;
    width: 60px;
    opacity: 0;
  }
  #group-icon {
    float: left;
    width: 60px;
    height: 60px;
    border-radius: 50%;
    position: absolute;
  }
  .store-group {
    border: 1px solid #DCDFE6;
    border-radius: 5px;
    padding: 10px;
  }
  .store-group .el-form-item {
    margin: 10px auto;
  }
  .group-title {
    text-align: center;
    font-size: 24px;
    font-weight: bold;
  }
  .operate-box {
    padding: 10px;
    border: 1px solid #DCDFE6;
    /* box-shadow: 0 2px 4px 0 rgb(0 0 0 / 12%), 0 0 6px 0 rgb(0 0 0 / 4%); */
  }
</style>
