package dingtalk

import (
	"excel2config/internal/pkg/httplib"
)

type Dingtalk struct {
	corpHost   string
	corpId     string
	corpSecret string
}

func New(corpHost, corpId, corpSecret string) (dt *Dingtalk) {
	dt = &Dingtalk{
		corpHost:   corpHost,
		corpId:     corpId,
		corpSecret: corpSecret,
	}
	return
}

type (
	userInfo struct {
		ErrCode  int    `json:"errcode"`
		ErrMsg   string `json:"errmsg"`
		UserId   string `json:"userid"`
		DeviceId string `json:"deviceid"`
		IsSys    bool   `json:"is_sys"`
		SysLevel int    `json:"sys_level"`
	}
	userDetail struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		UserId  string `json:"userid"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Avatar  string `json:"avatar"`
		Unionid string `json:"unionid"`
	}
)

func (dt *Dingtalk) GetAccessToken() (accessToken string, err error) {
	var s struct {
		AccessToken string `json:"access_token"`
	}
	get := httplib.Get(dt.corpHost + "/gettoken?corpid=" + dt.corpId + "&corpsecret=" + dt.corpSecret)
	err = get.ToJSON(&s)
	if err == nil {
		accessToken = s.AccessToken
	}
	return
}

func (dt *Dingtalk) GetUserInfo(token, code string) (uinfo *userInfo, err error) {
	uinfo = new(userInfo)
	get := httplib.Get(dt.corpHost + "/user/getuserinfo?access_token=" + token + "&code=" + code)
	err = get.ToJSON(uinfo)
	return
}

func (dt *Dingtalk) GetUserDetail(token, userid string) (detail *userDetail, err error) {
	detail = new(userDetail)
	get := httplib.Get(dt.corpHost + "/user/get?access_token=" + token + "&userid=" + userid)
	err = get.ToJSON(detail)
	return
}

func (dt *Dingtalk) GetAuthUsers(token, chatId string) (uidList []string, err error) {
	var s struct {
		ChatInfo struct {
			UseridList []string `json:"useridlist"`
		} `json:"chat_info"`
	}
	get := httplib.Get(dt.corpHost + "/chat/get?access_token=" + token + "&chatid=" + chatId)
	err = get.ToJSON(&s)
	if err == nil {
		uidList = s.ChatInfo.UseridList
	}
	return
}
