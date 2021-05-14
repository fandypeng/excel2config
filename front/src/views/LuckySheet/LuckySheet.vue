<template>
  <div id="luckysheet"
      style="margin:0px;padding:0px;position:absolute;width:100%;left: 0px;top: 0px;bottom:0px;height:100%;display:block;"
    ></div>
</template>

<script>
import { title } from '@/settings';
import { getToken } from '@/utils/auth';
import { mapGetters } from 'vuex'

export default {
  name: 'LuckySheet',
  props: {
    msg: String
  },
  data(){
    return {
      selected:"",
      isMaskShow: false,
      gridKey: this.$route.params.name,
    }
  },
  computed: {
    ...mapGetters([
      'uid'
    ])
  },
  mounted() {
    var vm = this
    // In some cases, you need to use $nextTick
    this.$nextTick(() => {
      var baseHost = process.env.VUE_APP_BASE_API
      var wsHost = baseHost.replace("8000", "8001")
      wsHost = wsHost.replace("http://", "ws://")
      var groupInfo = localStorage.getItem("groupInfo")
      groupInfo = JSON.parse(groupInfo)
      var allowEdit = groupInfo.IsDev
      if (groupInfo.IsDev == null) {
        allowEdit = true // TODO:: 临时打开线上环境的编辑权限
      }
      $(function () {
        luckysheet.create({
          container: "luckysheet",
          allowUpdate: true, //是否允许编辑后的后台更新
          gridKey: vm.$route.params.id,
          loadUrl: baseHost + "/excel?uid="+vm.uid+"&token="+getToken(), // 配置loadUrl的地址，luckysheet会通过ajax请求表格数据，默认载入status为1的sheet数据中的所有data，其余的sheet载入除data字段外的所有字段
          loadSheetUrl: baseHost + "/excel/sheet?uid="+vm.uid+"&token="+getToken(), //配置loadSheetUrl的地址，参数为gridKey（表格主键） 和 index（sheet主键合集，格式为[1,2,3]），返回的数据为sheet的data字段数据集合
          updateUrl: wsHost + "/uid/"+vm.uid+"/token/"+getToken(),
          myFolderUrl: "#/config/index",
          title: vm.$route.params.name,
          lang: "zh",
          allowEdit: allowEdit,
        });
      });
    });
  },
  destroyed () {
    // console.log(luckysheet)
    luckysheet.closeWebsocket()
    luckysheet.destroy()
    var retTip = document.getElementById("luckysheet-tooltip-up")
    if (retTip != null) {
      retTip.remove()
    }
  },
  methods:{

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
