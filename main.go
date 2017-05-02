// main
package main

import (
	"littleform/controllers"

	"github.com/astaxie/beego"
)

func main() {
	//	fmt.Println("Hello World!")
	//	createFile()
	beego.Router("/weixin", &controllers.WeixinController{})
	//	beego.Router("/:func", &mainController{})
	beego.Run()
}
