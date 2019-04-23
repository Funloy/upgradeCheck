// @APIVersion 1.0.0
// @Title 业务逻辑错误模型
// @Description 定义业务逻辑错误时发生的错误代码和错误描述
// @Contact xuchuangxin@icanmake.cn
// @TermsOfServiceUrl https://maiyajia.com/
// @License
// @LicenseUrl

package models

const (

	// 0 代表成功，非0代表错误
	_ = iota

	ERR_REQUEST_PARAM
	ERR_PERMISSION_DENIED

	ERR_LOGIN_FAIL
	ERR_LOGIN_TOKEN_FAIL
	ERR_TOKEN_FMT_FAIL
	ERR_ACCOUNT_EXISTS
	ERR_ACCOUNT_NONE
	ERR_ACCOUNT_LOCKED

	ERR_VERSION_CHECK_FAIL
	ERR_ACCOUNT_UPGRADE

	ERR_QUERY_COURSE_FAIL
)

var errorMsgs map[int]string

// GetErrorMsgs 获取错误表
func GetErrorMsgs() map[int]string {
	if errorMsgs == nil {

		errorMsgs = make(map[int]string)

		errorMsgs[ERR_REQUEST_PARAM] = "请求参数错误"
		errorMsgs[ERR_PERMISSION_DENIED] = "权限不足，无法访问"
		// 登录错误信息
		errorMsgs[ERR_LOGIN_FAIL] = "登录请求错误"
		errorMsgs[ERR_LOGIN_TOKEN_FAIL] = "无法获取登录令牌"
		errorMsgs[ERR_TOKEN_FMT_FAIL] = "token格式错误或token已过期"
		// 注册错误信息
		errorMsgs[ERR_ACCOUNT_EXISTS] = "产品用户已被占用"
		errorMsgs[ERR_ACCOUNT_LOCKED] = "产品用户被冻结，请联系管理员解封"
		errorMsgs[ERR_ACCOUNT_NONE] = "产品用户不存在"

		errorMsgs[ERR_VERSION_CHECK_FAIL] = "版本检测失败，请稍后重试"
		errorMsgs[ERR_ACCOUNT_UPGRADE] = "获取升级信息错误"

		errorMsgs[ERR_QUERY_COURSE_FAIL] = "获取课程升级信息错误"

	}
	return errorMsgs
}
