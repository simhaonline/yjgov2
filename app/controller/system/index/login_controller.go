package index

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/mssola/user_agent"
	"net/http"
	"strings"
	"time"
	"yj-app/app/model"
	logininforModel "yj-app/app/model/monitor/logininfor"
	"yj-app/app/model/monitor/online"
	logininforService "yj-app/app/service/monitor/logininfor"
	userService "yj-app/app/service/system/user"
	"yj-app/app/yjgframe/response"
	"yj-app/app/yjgframe/utils/gconv"
	"yj-app/app/yjgframe/utils/ip"
)

type RegisterReq struct {
	UserName     string `form:"username"  binding:"required,min=5,max=30"`
	Password     string `form:"password" binding:"required,min=5,max=30"`
	ValidateCode string `form:"validateCode" binding:"required,min=4,max=10"`
	IdKey        string `form:"idkey" binding:"required,min=5,max=30"`
}

// 登陆页面
func Login(c *gin.Context) {

	if strings.EqualFold(c.Request.Header.Get("X-Requested-With"), "XMLHttpRequest") {
		response.ErrorResp(c).SetMsg("未登录或登录超时。请重新登录").WriteJsonExit()
		return
	}

	response.BuildTpl(c, "login").WriteTpl()
}

// 图形验证码
func CaptchaImage(c *gin.Context) {
	//config struct for digits
	//数字验证码配置
	//var configD = base64Captcha.ConfigDigit{
	//	Height:     80,
	//	Width:      240,
	//	MaxSkew:    0.7,
	//	DotCount:   80,
	//	CaptchaLen: 5,
	//}
	//config struct for audio
	//声音验证码配置
	//var configA = base64Captcha.ConfigAudio{
	//	CaptchaLen: 6,
	//	Language:   "zh",
	//}
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height: 60,
		Width:  240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	//创建声音验证码
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	//idKeyA, capA := base64Captcha.GenerateCaptcha("", configA)
	//以base64编码
	//base64stringA := base64Captcha.CaptchaWriteToBase64Encoding(capA)
	//创建字符公式验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	idKeyC, capC := base64Captcha.GenerateCaptcha("", configC)
	//以base64编码
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)
	//创建数字验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	//idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	//以base64编码
	//base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	c.JSON(http.StatusOK, model.CaptchaRes{
		Code:  0,
		IdKey: idKeyC,
		Data:  base64stringC,
		Msg:   "操作成功",
	})
}

//验证登陆
func CheckLogin(c *gin.Context) {
	var req *RegisterReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}

	if req == nil {
		response.ErrorResp(c).SetMsg("用户名或密码错误").WriteJsonExit()
		return
	}

	//比对验证码
	verifyResult := base64Captcha.VerifyCaptcha(req.IdKey, req.ValidateCode)

	if !verifyResult {
		response.ErrorResp(c).SetMsg("验证码不正确").WriteJsonExit()
		return
	}

	isLock := logininforService.CheckLock(req.UserName)

	if isLock {
		response.ErrorResp(c).SetMsg("账号已锁定，请30分钟后再试").WriteJsonExit()
		return
	}

	//验证账号密码
	if sessionId, err := userService.SignIn(req.UserName, req.Password, c); err != nil {

		errTimes := logininforService.SetPasswordCounts(req.UserName)

		having := 5 - errTimes

		//记录日志
		var logininfor logininforModel.Entity
		logininfor.LoginName = req.UserName
		logininfor.Ipaddr = c.ClientIP()

		userAgent := c.Request.Header.Get("User-Agent")
		ua := user_agent.New(userAgent)
		logininfor.Os = ua.OS()
		logininfor.Browser, _ = ua.Browser()
		logininfor.LoginTime = time.Now()
		logininfor.LoginLocation = ip.GetCityByIp(logininfor.Ipaddr)
		logininfor.Msg = "账号或密码不正确"
		logininfor.Status = "0"

		logininfor.Insert()

		response.ErrorResp(c).SetMsg("账号或密码不正确,还有" + gconv.String(having) + "次之后账号将锁定").WriteJsonExit()
	} else {
		//保存在线状态

		userAgent := c.Request.Header.Get("User-Agent")
		ua := user_agent.New(userAgent)
		os := ua.OS()
		browser, _ := ua.Browser()
		loginIp := c.ClientIP()
		loginLocation := ip.GetCityByIp(loginIp)

		var userOnline online.Entity
		userOnline.Sessionid = sessionId
		userOnline.LoginName = req.UserName
		userOnline.Browser = browser
		userOnline.Os = os
		userOnline.DeptName = ""
		userOnline.Ipaddr = loginIp
		userOnline.ExpireTime = 1440
		userOnline.StartTimestamp = time.Now()
		userOnline.LastAccessTime = time.Now()
		userOnline.Status = "on_line"
		userOnline.LoginLocation = loginLocation
		userOnline.Delete()
		userOnline.Insert()

		//移除登陆次数记录
		logininforService.RemovePasswordCounts(req.UserName)
		//记录日志
		var logininfor logininforModel.Entity
		logininfor.LoginName = req.UserName
		logininfor.Ipaddr = loginIp

		logininfor.Os = os
		logininfor.Browser = browser
		logininfor.LoginTime = time.Now()
		logininfor.LoginLocation = loginLocation
		logininfor.Msg = "登陆成功"
		logininfor.Status = "0"

		logininfor.Insert()
		response.SucessResp(c).SetMsg("登陆成功").WriteJsonExit()
	}
}
