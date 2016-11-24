package controllers

import (

)

type MainController struct {
	CommonController
}



func (c *MainController) Get() {
	c.TplName = "index.html"
}
