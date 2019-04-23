// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"upgrade.maiyajia.com/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/upgrade",
			beego.NSRouter("/login", &controllers.PassportController{}, "post:Login"),
			beego.NSRouter("/check", &controllers.UpgradeController{}, "post:UpgradeCheck"),
			beego.NSRouter("/done", &controllers.UpgradeController{}, "post:UpgradeDone"),
			beego.NSRouter("/tools", &controllers.UpgradeController{}, "post:UpgradeCheckTools"),
			beego.NSRouter("/courses", &controllers.UpgradeController{}, "post:UpgradeCheckCourses"),
		),
	)
	beego.AddNamespace(ns)
}
