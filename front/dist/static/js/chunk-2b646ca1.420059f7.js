(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2b646ca1"],{"0568":function(e,t,c){},"3e41":function(e,t,c){"use strict";c("68dc")},"68dc":function(e,t,c){},7382:function(e,t,c){"use strict";c("0568")},7896:function(e,t,c){"use strict";c.r(t);var n=function(){var e=this,t=e.$createElement,c=e._self._c||t;return c("div",{attrs:{id:"excel"}},[c("LuckySheet",{attrs:{msg:"Welcome to Your Vue.js App"}})],1)},o=[],a=function(){var e=this,t=e.$createElement,c=e._self._c||t;return c("div",{staticStyle:{margin:"0px",padding:"0px",position:"absolute",width:"100%",left:"0px",top:"0px",bottom:"0px",height:"100%",display:"block"},attrs:{id:"luckysheet"}})},l=[],u=(c("b0c0"),c("ac1f"),c("5319"),c("5530")),i=(c("83d6"),c("5f87")),r=c("2f62"),s={name:"LuckySheet",props:{msg:String},data:function(){return{selected:"",isMaskShow:!1,gridKey:this.$route.params.name}},computed:Object(u["a"])({},Object(r["b"])(["uid"])),mounted:function(){var e=this;this.$nextTick((function(){var t="http://e2c.17zjh.com:8000",c=t.replace("8000","8001");c=c.replace("http://","ws://");var n=localStorage.getItem("groupInfo");n=JSON.parse(n);var o=n.IsDev;null==n.IsDev&&(o=!1),$((function(){luckysheet.create({container:"luckysheet",allowUpdate:!0,gridKey:e.$route.params.id,loadUrl:t+"/excel?uid="+e.uid+"&token="+Object(i["a"])(),loadSheetUrl:t+"/excel/sheet?uid="+e.uid+"&token="+Object(i["a"])(),updateUrl:c+"/uid/"+e.uid+"/token/"+Object(i["a"])(),myFolderUrl:"#/config/index",title:e.$route.params.name,lang:"zh",allowEdit:o})}))}))},destroyed:function(){luckysheet.closeWebsocket(),luckysheet.destroy();var e=document.getElementById("luckysheet-tooltip-up");null!=e&&e.remove()},methods:{}},d=s,p=(c("3e41"),c("2877")),h=Object(p["a"])(d,a,l,!1,null,"6eda6a32",null),m=h.exports,f={name:"excel",components:{LuckySheet:m}},k=f,y=(c("7382"),Object(p["a"])(k,n,o,!1,null,null,null));t["default"]=y.exports}}]);