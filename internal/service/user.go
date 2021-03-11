package service

import (
	"context"
	"crypto/tls"
	pb "excel2config/api"
	"excel2config/internal/def"
	"excel2config/internal/helper"
	"excel2config/internal/model"
	"excel2config/internal/pkg/dingtalk"
	"excel2config/internal/server/sessions"
	"fmt"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/ecode"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-ldap/ldap/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	//ldap "github.com/go-ldap/ldap/v3"
	"github.com/prometheus/common/log"
)

func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (resp *pb.LoginResp, err error) {
	var userInfo *model.UserInfo
	switch req.LoginType {
	case pb.LoginType_Common:
		if !helper.IsValidEmail(req.Email) {
			err = ecode.Int(int(def.ErrEmailFormat))
			return
		}
		userInfo, err = s.dao.GetUser(ctx, req.Email)
		if err != nil {
			err = ecode.Int(int(def.ErrUserNotExist))
			return
		}
		if userInfo.LoginType != int(pb.LoginType_Common) {
			err = ecode.Int(int(def.ErrLoginParam))
			return
		}
		inPwd := helper.Md5Sum(req.Pwd + def.PasswdSalt)
		if userInfo.Passwd != inPwd {
			err = ecode.Int(int(def.ErrUserPasswd))
			return
		}
	case pb.LoginType_Ldap:
		if !helper.IsValidEmail(req.Email) {
			err = ecode.Int(int(def.ErrEmailFormat))
			return
		}
		userInfo, err = s.ldapLogin(ctx, req)
	case pb.LoginType_DingDing:
		if len(req.Code) == 0 {
			err = ecode.Int(int(def.ErrLoginParam))
			return
		}
		userInfo, err = s.dingtalkLogin(ctx, req)
	default:
		err = ecode.Int(int(def.ErrLoginParam))
		return
	}
	if err != nil {
		return
	}

	resp = &pb.LoginResp{
		UserInfo: s.copyUserInfo(userInfo),
		Token:    helper.GetRandomToken(userInfo.Uid),
	}
	// 保存token
	err = s.dao.SaveToken(ctx, userInfo.Uid, resp.Token)
	if err != nil {
		log.With("err", err).With("userInfo", userInfo).With("token", resp.Token).Errorln("save token error")
		return
	}
	sess := sessions.Default(ctx.(*bm.Context))
	sess.Set("uid", userInfo.Uid)
	sess.Set("username", userInfo.UserName)
	sess.Set("email", userInfo.Email)
	sess.Set("avatar", userInfo.Avatar)
	sess.Set("token", resp.Token)
	sess.Save()
	return
}

func (s *Service) ldapLogin(ctx context.Context, req *pb.LoginReq) (userInfo *model.UserInfo, err error) {
	var cfg struct {
		Ldap struct {
			ServerHost     string
			BindBaseDn     string
			BindPassword   string
			SearchBaseDn   string
			SearchBy       string
			SearchUserName string
		}
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	if cfg.Ldap.ServerHost == "" {
		err = ecode.Int(int(def.ErrLdapConfig))
		return
	}
	l, err := ldap.DialURL(cfg.Ldap.ServerHost)
	if err != nil {
		return
	}
	defer l.Close()

	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	_, err = l.SimpleBind(&ldap.SimpleBindRequest{
		Username: cfg.Ldap.BindBaseDn,
		Password: cfg.Ldap.BindPassword,
	})
	if err != nil {
		return
	}
	searchRequest := ldap.NewSearchRequest(
		"dc=wepie,dc=com",
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(&(%s=%s))", cfg.Ldap.SearchBy, req.Email),
		[]string{"dn", cfg.Ldap.SearchUserName, cfg.Ldap.SearchBy},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return
	}
	if len(sr.Entries) == 0 {
		err = ecode.Int(int(def.ErrUserNotExist))
		return
	}
	searchRes := sr.Entries[0]
	err = l.Bind(searchRes.DN, req.Pwd)
	if err != nil {
		err = ecode.Int(int(def.ErrUserPasswd))
		return
	}
	userInfo, err = s.dao.GetUser(ctx, req.Email)
	if err == mongo.ErrNoDocuments {
		//record a random password to avoid common login
		u, _ := uuid.New()
		pwd := helper.Md5Sum(string(u[:]))
		userInfo, err = s.dao.AddUser(ctx, req.Email, pwd, searchRes.GetAttributeValue(cfg.Ldap.SearchUserName), int(pb.LoginType_Ldap))
	}
	return
}

func (s *Service) dingtalkLogin(ctx context.Context, req *pb.LoginReq) (userInfo *model.UserInfo, err error) {
	var cfg struct {
		Dingtalk struct {
			corpHost   string
			corpId     string
			corpSecret string
			chatId     string
		}
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	if cfg.Dingtalk.corpId == "" || cfg.Dingtalk.corpSecret == "" || cfg.Dingtalk.corpHost == "" {
		err = ecode.Int(int(def.ErrDingtalkConfig))
		return
	}
	corp := dingtalk.New(cfg.Dingtalk.corpHost, cfg.Dingtalk.corpId, cfg.Dingtalk.corpSecret)
	token, err := corp.GetAccessToken()
	if err != nil {
		log.With("err", err).Errorln("get token failed")
		err = ecode.Int(int(def.ErrLoginFailed))
		return
	}
	dingUserInfo, err := corp.GetUserInfo(token, req.Code)
	if err != nil || dingUserInfo.ErrCode != 0 {
		log.With("err", err).With("userInfo", userInfo).Errorln("login failed")
		err = ecode.Int(int(def.ErrLoginFailed))
		return
	}
	userDetail, err := corp.GetUserDetail(token, dingUserInfo.UserId)
	if err != nil || userDetail.ErrCode != 0 {
		log.With("err", err).With("userInfo", userInfo).Errorln("login failed")
		err = ecode.Int(int(def.ErrLoginFailed))
		return
	}
	//check whether the user in chat group if chatId is configured
	if len(cfg.Dingtalk.chatId) > 0 {
		uidList, ierr := corp.GetAuthUsers(token, cfg.Dingtalk.chatId)
		if ierr != nil || !helper.Contains(uidList, userDetail.UserId) {
			log.With("err", ierr).Errorln("get chat group users failed")
			err = ecode.Int(int(def.ErrLoginDenied))
			return
		}
	}
	userInfo, err = s.dao.GetUser(ctx, userDetail.Email)
	if err == mongo.ErrNoDocuments {
		//record a random password to avoid common login
		u, _ := uuid.New()
		pwd := helper.Md5Sum(string(u[:]))
		userInfo, err = s.dao.AddUser(ctx, req.Email, pwd, userDetail.Name, int(pb.LoginType_DingDing))
	}
	if err == nil {
		userInfo.Avatar = userDetail.Avatar
	}
	return
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterReq) (resp *pb.RegisterResp, err error) {
	if !helper.IsValidEmail(req.Email) {
		err = ecode.Int(int(def.ErrEmailFormat))
		return
	}
	if req.Pwd != req.ConfirmPwd {
		err = ecode.Int(int(def.ErrPwdNotConfirmed))
		return
	}
	userInfo, err := s.dao.AddUser(ctx, req.Email, req.Pwd, req.Name, int(pb.LoginType_Common))
	if err != nil {
		log.With("err", err).With("userInfo", userInfo).Errorln("add user error")
		return
	}
	resp = &pb.RegisterResp{
		UserInfo: s.copyUserInfo(userInfo),
		Token:    helper.GetRandomToken(userInfo.Uid),
	}
	// 保存token
	err = s.dao.SaveToken(ctx, userInfo.Uid, resp.Token)
	if err != nil {
		log.With("err", err).With("userInfo", userInfo).With("token", resp.Token).Errorln("save token error")
		return
	}
	sess := sessions.Default(ctx.(*bm.Context))
	sess.Set("uid", userInfo.Uid)
	sess.Set("username", userInfo.UserName)
	sess.Set("email", userInfo.Email)
	sess.Set("avatar", userInfo.Avatar)
	sess.Set("token", resp.Token)
	sess.Save()
	return
}

func (s *Service) Logout(ctx context.Context, req *pb.LogoutReq) (resp *pb.LogoutResp, err error) {
	err = s.dao.SaveToken(ctx, req.Uid, "")
	if err != nil {
		sess := sessions.Default(ctx.(*bm.Context))
		sess.Delete("uid")
		sess.Delete("username")
		sess.Delete("email")
		sess.Delete("avatar")
		sess.Delete("token")
	}
	resp = &pb.LogoutResp{}
	return
}

func (s *Service) Info(ctx context.Context, req *pb.UserInfoReq) (resp *pb.UserInfo, err error) {
	var email string
	sess := sessions.Default(ctx.(*bm.Context))
	emailInterface := sess.Get("email")
	if emailInterface != nil {
		email = emailInterface.(string)
	} else {
		err = ecode.Int(int(def.ErrNeedLogin))
		return
	}
	userInfo, err := s.dao.GetUser(ctx, email)
	if err != nil {
		err = ecode.Int(int(def.ErrUserNotExist))
		return
	}
	resp = s.copyUserInfo(userInfo)
	return
}

func (s *Service) Search(ctx context.Context, req *pb.UserSearchReq) (resp *pb.UserSearchResp, err error) {
	resp = &pb.UserSearchResp{}
	if req.Name == "" {
		return
	}
	userInfos, err := s.dao.GetUserInfosByKeyword(ctx, req.Name)
	if err != nil {
		log.With("err", err).With("name", req.Name).Errorln("get user infos error")
		return
	}
	for _, userInfo := range userInfos {
		resp.UserInfos = append(resp.UserInfos, s.copySimpleUserInfo(userInfo))
	}
	return
}

func (s *Service) copyUserInfo(userInfo *model.UserInfo) *pb.UserInfo {
	return &pb.UserInfo{
		Uid:      userInfo.Uid,
		UserName: userInfo.UserName,
		Email:    userInfo.Email,
		RegTime:  userInfo.RegTime,
		Avatar:   userInfo.Avatar,
	}
}

func (s *Service) copySimpleUserInfo(userInfo *model.SimpleUserInfo) *pb.SimpleUserInfo {
	return &pb.SimpleUserInfo{
		Uid:      userInfo.Uid,
		UserName: userInfo.UserName,
		Avatar:   userInfo.Avatar,
	}
}
