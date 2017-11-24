package controllers

import (
	"github.com/satori/go.uuid"
	"github.com/revel/revel"
	"os/exec"
	"santa/game-server-mgr/app"
)

type App struct {
	*revel.Controller
}

type StartGameResponse struct {
	GameID         string	`json:"GameID"`
	ServerHostname string	`json:"ServerHostname"`
	ServerPort     int		`json:"ServerPort"`
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) TestIncr() revel.Result {
	app.RedisCli.Incr("count")
	r, _ := app.RedisCli.Get("count").Result()
	return c.RenderText("count: " + r)
}

func (c App) StartGame() revel.Result {
	gameID := uuid.NewV4().String()

	revel.INFO.Println("[app::StartGame] Starting new game instance id = " + gameID)

	go func() {
		cmd := exec.Command("./test-server", "-batchmode", "-nographics", "-logfile", "log.txt")
		cmd.Dir = "builds/test-server.app/Contents/MacOS/"
		_, err := cmd.Output()
		if err != nil {
			println(err.Error())
		}
		revel.INFO.Println("[app::StartGame] Closed game instance id = " + gameID)
	}()
	
	data := make(map[string]interface{})
	data["error"] = nil
	res := StartGameResponse{
		GameID:         gameID,
		ServerHostname: "localhost",
		ServerPort:     7777,
	}
	data["data"] = res
	
	return c.RenderJSON(data)
}
