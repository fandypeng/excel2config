<template>
  <el-row :gutter="40" class="panel-group" v-loading="listLoading" element-loading-text="Loading">
    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col" v-for="item in list" :key="item.gid">
      <div class="card-panel" @click="handleClickGroupPanel(item)">
        <div class="card-panel-icon-wrapper icon-box">
          <img :src="item.avatar" alt="" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            {{item.name}}
          </div>
        </div>
      </div>
    </el-col>

    <el-col :xs="12" :sm="12" :lg="6" class="card-panel-col">
      <div class="card-panel" @click=handleClickAddGroup()>
        <div class="card-panel-icon-wrapper icon-box">
          <i class="el-icon-plus" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
            添加项目
          </div>
        </div>
      </div>
    </el-col>

    <el-dialog :visible.sync="dialogVisible" title="添加项目">
      <el-form :model="groupInfo" label-width="80px" label-position="left">
        <el-form-item label="项目名称">
          <el-input v-model="groupInfo.name" placeholder="项目名称" />
        </el-form-item>
        <el-form-item label="头像">
          <el-col :span="11">
            <img id="group-icon" v-bind:src="groupInfo.avatar" alt="">
            <input class="avatar-input" type="file" accept="image/*" @change="uploadAvatar" name="avatar"   />
          </el-col>
        </el-form-item>
        <br>
        <el-form-item label="备注">
          <el-input v-model="groupInfo.remark" type="textarea" placeholder="项目备注" />
        </el-form-item>
      </el-form>

      <div style="text-align:right;">
        <el-button type="danger" @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="addGroupInfo">保存</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>

<script>
import { groupList, groupAdd } from '@/api/group'

export default {
  data() {
    return {
      dialogVisible: false,
      groupInfo: {
        name: "",
        avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
      },
      listLoading: false,
      list: [],
    }
  },
  components: {
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      groupList({}).then(response => {
        if (response.data.list != undefined) {
          this.list = response.data.list
        }
        this.listLoading = false
      })
    },
    handleClickGroupPanel(groupInfo) {
      localStorage.setItem("groupInfo", JSON.stringify(groupInfo))
      this.$router.push({
        path: "/config/index",
      })
    },
    handleClickAddGroup() {
      this.dialogVisible = true
    },
    addGroupInfo() {
      console.log(this.groupInfo)
      groupAdd(this.groupInfo).then(response => {
        if (response.data.groupInfo != undefined) {
          this.list.push(response.data.groupInfo)
          this.dialogVisible = false
        }
        console.log(this.list)
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
        if (arr[i] == target) {
          return true
        }
      }
      return false
    },
  }
}
</script>

<style lang="scss" scoped>
.panel-group {
  margin-top: 18px;

  .card-panel-col {
    margin-bottom: 32px;
  }

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

  .card-panel {
    height: 108px;
    cursor: pointer;
    font-size: 12px;
    position: relative;
    overflow: hidden;
    color: #666;
    background: #fff;
    box-shadow: 4px 4px 40px rgba(0, 0, 0, .05);
    border-color: rgba(0, 0, 0, .05);

    // &:hover {
      // .card-panel-icon-wrapper {
      //   color: #fff;
      // }
    // }

    .icon-box img {
      width: 64px;
      border-radius: 50%;
    }

    .card-panel-icon-wrapper {
      float: left;
      margin: 22px 0 22px 10px;
      // padding: 16px;
      transition: all 0.38s ease-out;
      border-radius: 6px;
    }

    .card-panel-icon {
      float: left;
      font-size: 48px;
    }

    .card-panel-description {
      float: right;
      font-weight: bold;
      margin-right: 20px;
      margin-left: 0px;
      max-width: 300px;

      .card-panel-text {
        line-height: 108px;
        color: rgba(0, 0, 0, 0.45);
        font-size: 24px;
        // margin-bottom: 12px;
      }
    }
    .icon-box i {
      font-size: 64px;
      font-weight: bold;
    }
  }
}

@media (max-width:550px) {
  // .card-panel-description {
  //   display: none;
  // }
  .card-panel-text {
    display: block!important;
    margin: 10 auto;
    text-align: center;
  }
  .icon-box {
    text-align: center;
  }
  .icon-box img {
    margin: 14px auto!important;
    width: 48px;
  }
  .card-panel-icon-wrapper {
    float: none !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;

    .svg-icon {
      display: block;
      margin: 14px auto !important;
      float: none !important;
    }
  }
}
</style>
