// @APIVersion 1.0.0
// @Title 用户通行控制器
// @Description 本控制器提供的接口为整个平台提供产品用户登录验证以及获取通行令牌TOKEN的功能。
// @Contact xuchuangxin@icanmake.cn
// @TermsOfServiceUrl https://maiyajia.com/
// @License
// @LicenseUrl

package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
	m "upgrade.maiyajia.com/models"
	"upgrade.maiyajia.com/services/mongo"
	"upgrade.maiyajia.com/services/token"
)

// PassportController 登录注册控制器
type PassportController struct {
	BaseController
}

// NestPrepare 初始化函数
// 把控制器的MgoClient赋值到数据库操作客户端
func (passportCtl *PassportController) NestPrepare() {
	mongo.Client = &passportCtl.MgoClient
}

// LoginCredential 登录凭证
type loginCredential struct {
	Key    string `form:"key"`
	Serial string `form:"serial"`
}

// Login 登录
func (passportCtl *PassportController) Login() {
	var credential loginCredential

	// 登录凭证解析错误
	if err := json.Unmarshal(passportCtl.Ctx.Input.RequestBody, &credential); err != nil {
		passportCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}

	// 查找产品用户账号，并检查账号的有效性（包括是否存在，是否被锁...)
	account, err := m.FindAccount(bson.M{"productKey": credential.Key, "productSerial": credential.Serial})
	if account == nil || err != nil {
		beego.Error(err)
		passportCtl.abortWithError(m.ERR_ACCOUNT_NONE)
	}
	if account.Locked {
		passportCtl.abortWithError(m.ERR_ACCOUNT_LOCKED)
	}

	// 登录成功，用AccountID，ProductKey和ProductSerial来获取JWT的TOKEN
	token, err := createClientToken(account.ID.Hex(), account.ProductKey, account.ProductSerial)
	if err != nil {
		passportCtl.abortWithError(m.ERR_LOGIN_TOKEN_FAIL)
	}

	// 用户登录日志

	// 封装返回数据
	out := make(map[string]interface{})
	out["code"] = 0
	out["token"] = token
	// 返回结果
	passportCtl.jsonResult(out)

}

// @Title 创建登录令牌
// @Description 用通行证的信息创建用户登录的令牌TOKEN
// @Param	aid		string	true	"产品用户ID"
// @param	key		string	true	"产品KEY"
// @param	serial	string	true	"产品序列号"
// @Success token {string}
// @Failure err {error}
func createClientToken(aid, key, serial string) (string, error) {
	token := token.Token{
		AccountID:     aid,
		ProductKey:    key,
		ProductSerial: serial,
	}
	return token.CreateToken()
}
