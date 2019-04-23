// @APIVersion 1.0.0
// @Title 控制器基类
// @Description 控制器基础类，为后面继承该类的控制器提供一些共用的方法或操作
// @Contact xuchuangxin@icanmake.cn
// @TermsOfServiceUrl https://maiyajia.com/
// @License
// @LicenseUrl

package controllers

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	m "upgrade.maiyajia.com/models"
	"upgrade.maiyajia.com/services/mongo"
	"upgrade.maiyajia.com/services/token"
)

//BaseController 基础控制器
type (
	BaseController struct {
		beego.Controller
		mongo.MgoClient
	}
)

//NestPreparer 作为子类自定义prepare使用
type NestPreparer interface {
	NestPrepare()
}

//Prepare 作初始化处理
func (base *BaseController) Prepare() {
	base.MgoClient.StartSession()
	if app, ok := base.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

//Finish 后处理
func (base *BaseController) Finish() {
	defer func() {
		base.MgoClient.CloseSession()
	}()
}

// jsonResult 服务端返回json
func (base *BaseController) jsonResult(out interface{}) {
	base.Data["json"] = out
	base.ServeJSON()
}

// abortWithError 根据错误码获取错误描述信息，然后发送到请求客户端
func (base *BaseController) abortWithError(code int) {
	var msg string
	if v, ok := m.GetErrorMsgs()[code]; ok {
		msg = v
	} else {
		msg = "未知错误"
	}

	result := ErrorResult{
		Code:    code,
		Message: msg,
	}
	base.Data["json"] = result
	base.ServeJSON()

	/* 调用 StopRun 之后，如果你还定义了 Finish 函数就不会再执行。
	* 如果需要释放资源，那么请自己在调用 StopRun 之前手工调用 Finish 函数。
	* https://beego.me/docs/mvc/controller/controller.md
	 */

	// 调用Finish()函数，释放数据库资源
	base.Finish()
	// 终止执行
	base.StopRun()
}

// checkToken 读取请求的token
func (base *BaseController) checkToken() *token.Token {
	// JWT验证
	token, err := parseClientToken(base.Ctx.Input.Header("Authorization"))
	if err != nil {
		base.abortWithError(m.ERR_TOKEN_FMT_FAIL)
	}
	return token
}

// parseClientToken 根据客户端请求的令牌TOKEN字符串，解析出TOKEN信息
func parseClientToken(authHeader string) (*token.Token, error) {

	auths := strings.Split(authHeader, " ")
	if len(auths) != 2 || auths[0] != "Bearer" {
		return nil, errors.New("authorization invalid")
	}

	token, err := token.ValidateToken(auths[1])
	if err != nil {
		return nil, errors.New("authorization invalid")
	}
	return token, nil
}
