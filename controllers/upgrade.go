// @APIVersion 1.0.0
// @Title 升级控制器
// @Description 升级控制器，提供产品升级检查，升级信息等功能
// @Contact xuchuangxin@icanmake.cn
// @TermsOfServiceUrl https://maiyajia.com/
// @License
// @LicenseUrl

package controllers

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	m "upgrade.maiyajia.com/models"
	"upgrade.maiyajia.com/services/mongo"
)

// UpgradeController 升级控制器
type UpgradeController struct {
	BaseController
}

// NestPrepare 初始化函数
// 把控制器的MgoClient赋值到数据库操作客户端
func (upgradeCtl *UpgradeController) NestPrepare() {
	mongo.Client = &upgradeCtl.MgoClient
}

// upgradeCredential 检查版本请求凭证
type upgradeCredential struct {
	OS      string `form:"os"`
	Version string `form:"version"`
}

// upgradeTools 检查版本请求凭证
type upgradeTools struct {
	Tools []*m.Tool `form:"tools"`
}

// upgradeCourses 检查版本请求凭证
type upgradeCourses struct {
	URL     string        `form:"url"`
	Key     string        `form:"key"`
	Serial  string        `form:"serial"`
	Courses []*m.UpCourse `form:"courses"`
}

// Upgrade 返回的升级信息
type Upgrade struct {
	Name      string    `bson:"name" json:"name"`
	Version   string    `bson:"version" json:"version"`
	Asset     *m.Asset  `bson:"asset" json:"asset"`
	Changelog string    `bson:"changelog" json:"changelog"`
	Date      time.Time `bson:"date" json:"date"`
}

// UpgradeCheck 版本检测
func (upgradeCtl *UpgradeController) UpgradeCheck() {

	var credential upgradeCredential
	// 验证升级请求数据
	if err := json.Unmarshal(upgradeCtl.Ctx.Input.RequestBody, &credential); err != nil {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	// 检查请求信息是否符合要求
	if !validateOSArch(credential.OS) {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	if !validateVersion(credential.Version) {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	product, isUpgrade := m.QueryUpgradeProduct(credential.Version)
	if !isUpgrade {
		// 没有最新版本，不需要升级
		logs.Info("product, isUpgrade:", product, isUpgrade)
		out := make(map[string]interface{})
		out["code"] = 0
		out["newver"] = false
		upgradeCtl.jsonResult(out)
	}
	beego.Informational(product)
	if product != nil {
		// 封装升级信息
		upgrade := new(Upgrade)

		for key, value := range product.Assets {
			if value.OS == credential.OS {
				upgrade.Asset = product.Assets[key]
				break
			}
		}
		upgrade.Name = product.Name
		upgrade.Version = product.Version
		upgrade.Changelog = product.Changelog
		upgrade.Date = product.Date

		// 封装返回数据
		out := make(map[string]interface{})
		out["code"] = 0
		out["newver"] = true
		out["upgrade"] = upgrade
		upgradeCtl.jsonResult(out)
	}
}

// UpgradeCheckTools 工具升级新增检测检测
func (upgradeCtl *UpgradeController) UpgradeCheckTools() {
	var upgradeTools *upgradeTools
	var resultTools []*m.Tool
	// 验证升级请求数据
	if err := json.Unmarshal(upgradeCtl.Ctx.Input.RequestBody, &upgradeTools); err != nil {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	//比较工具版本，将有版本升级的放入数组切片中
	for _, tool := range upgradeTools.Tools {
		if !validateVersion(tool.Version) {
			logs.Info("版本号不合法")
			continue
		}
		upgradeTool, isUpgrade, err := m.QueryUpgradeTool(tool.ID, tool.Name, tool.Version)
		if err != nil && "not found" == err.Error() {
			logs.Error("QueryUpgradeTool err:", err)
			continue
		}
		if isUpgrade {
			resultTools = append(resultTools, upgradeTool)
		}
	}
	alltools, err := m.GetToolsInfo()
	if err != nil {
		logs.Error("GetToolsInfo err:", err)
	}
	for _, protool := range alltools {
		flag := true
		for _, tool := range upgradeTools.Tools {
			if protool.ID.Hex() == tool.ID.Hex() {
				//if protool.Name == tool.Name {
				flag = false
				break
			}
		}
		if flag {
			resultTools = append(resultTools, protool)
		}
	}
	if len(resultTools) < 1 {
		out := make(map[string]interface{})
		out["code"] = 0
		out["newver"] = false
		upgradeCtl.jsonResult(out)
	}
	// 封装返回数据
	out := make(map[string]interface{})
	out["code"] = 0
	out["newver"] = true
	out["tools"] = resultTools
	upgradeCtl.jsonResult(out)
}

// UpgradeCheckCourses 课程新增检测
func (upgradeCtl *UpgradeController) UpgradeCheckCourses() {
	var upgradeCourses *upgradeCourses
	var resultCourses []*m.Course
	// 验证升级请求数据
	if err := json.Unmarshal(upgradeCtl.Ctx.Input.RequestBody, &upgradeCourses); err != nil {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	url := upgradeCourses.URL
	productKey := upgradeCourses.Key
	productSerial := upgradeCourses.Serial
	//获取已购买的所有课程
	procourses, err := m.GetCoursesInfo(url, productKey, productSerial)
	if err != nil {
		logs.Error("GetCoursesInfo fail:", err)
		upgradeCtl.abortWithError(m.ERR_ACCOUNT_UPGRADE)
	}
	//比较课程版本，将有版本升级的放入数组切片中

	for _, upgradeCou := range upgradeCourses.Courses {
		if !validateVersion(upgradeCou.Version) {
			logs.Error("版本号不合法")
			continue
		}
		upgradeCourse, isUpgrade, err := m.QueryUpgradeCourse(upgradeCou.ID, upgradeCou.Version)
		if err != nil && "not found" == err.Error() {
			logs.Error("QueryUpgradeTool err:", err)
			continue
		}
		if isUpgrade {
			//与数据库对比时将Purchased赋值为true传到学习平台
			upgradeCourse.Purchased = true
			resultCourses = append(resultCourses, upgradeCourse)
		}
	}
	for _, procourse := range procourses {
		flag := true
		for _, upgradeCou := range upgradeCourses.Courses {
			if procourse.ID.Hex() == upgradeCou.ID.Hex() {
				//if procourse.Name == upgradeCou.Name {
				flag = false
				break
			}
		}

		if flag {
			resultCourses = append(resultCourses, procourse)
		}
	}

	// for _, procourse := range procourses {
	// 	if !validateVersion(procourse.Version) {
	// 		logs.Info("版本号不合法")
	// 		continue
	// 	}
	// 	upgradeCourse, isUpgrade := m.QueryUpgradeCourse(course.Name, course.Version)
	// 	if isUpgrade {
	// 		resultCourses = append(resultCourses, upgradeCourse)
	// 	}
	// 	flag := true
	// 	for _, course := range upgradeCourses.Courses {

	// 		//if procourse.ID.Hex() == course.ID.Hex() {
	// 		if procourse.Name == course.Name {
	// 			flag = false
	// 			break
	// 		}
	// 	}
	// 	if flag {
	// 		resultCourses = append(resultCourses, procourse)
	// 	}
	// }
	if len(resultCourses) < 1 {
		out := make(map[string]interface{})
		out["code"] = 0
		out["newver"] = false
		upgradeCtl.jsonResult(out)
	}
	// 封装返回数据
	out := make(map[string]interface{})
	out["code"] = 0
	out["newver"] = true
	out["courses"] = resultCourses
	upgradeCtl.jsonResult(out)
}

// UpgradeDone 客户端升级成功回调接口
func (upgradeCtl *UpgradeController) UpgradeDone() {

	token := upgradeCtl.checkToken()

	var credential upgradeCredential
	// 验证升级请求数据
	if err := json.Unmarshal(upgradeCtl.Ctx.Input.RequestBody, &credential); err != nil {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}

	// 检查请求信息是否符合要求
	if !validateOSArch(credential.OS) {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}
	if !validateVersion(credential.Version) {
		upgradeCtl.abortWithError(m.ERR_REQUEST_PARAM)
	}

	if err := m.AccountUpgraded(token.ProductKey, token.ProductSerial, credential.OS, credential.Version); err != nil {
		upgradeCtl.abortWithError(m.ERR_ACCOUNT_UPGRADE)
	}
	// 封装返回数据
	out := make(map[string]interface{})
	out["code"] = 0
	upgradeCtl.jsonResult(out)
}

func validateVersion(version string) bool {
	// 产品版本号正则，匹配形如 xxxx.xx.xx形式 比如 1.0.1
	reg := regexp.MustCompile(`^[\d+]+\.[0-9]{1,2}\.[0-9]{1,2}`)
	return reg.MatchString(version)
}

// 检查请求参数中的OS是否符合要求
func validateOSArch(os string) bool {
	ok := false
	for _, value := range m.OS_ARCH {
		if os == value {
			ok = true
			break
		}
	}
	return ok
}
