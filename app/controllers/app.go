package controllers

import (
	"github.com/revel/revel"
	"santa/game-server-mgr/app"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) TestIncr() revel.Result {
	app.RedisCli.Incr("count")
	r, _ := app.RedisCli.Get("count").Result()
	return c.RenderText("count: " + r)
}
