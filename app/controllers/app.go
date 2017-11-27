package controllers

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/revel/revel"
	"os/exec"
	"santa/game-server-mgr/app"
	"archive/zip"
	"bytes"
	"os"
	"io"
	"path/filepath"
	"strings"
)

type App struct {
	*revel.Controller
	
}

type StartGameResponse struct {
	GameID         string	`json:"GameID"`
	ServerHostname string	`json:"ServerHostname"`
	ServerPort     int		`json:"ServerPort"`
}

const (
	buildPath = "builds"
	executableName = "server"
)

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) TestIncr() revel.Result {
	app.RedisCli.Incr("count")
	r, _ := app.RedisCli.Get("count").Result()
	return c.RenderText("count: " + r)
}

func (c App) UploadBuild(buildZipFile []byte) revel.Result {
	dest := buildPath
	var filenames []string

	r, err := zip.NewReader(bytes.NewReader(buildZipFile), (int64)(len(buildZipFile)))
	if err != nil {
		revel.ERROR.Printf(err.Error())
		return c.RenderError(err)
	}
	for _, f := range r.File {
		
		rc, err := f.Open()
		if err != nil {
			revel.ERROR.Println(err.Error())
			return c.RenderError(err)
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				revel.ERROR.Println(err.Error())
				return c.RenderError(err)
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				revel.ERROR.Println(err.Error())
				return c.RenderError(err)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				revel.ERROR.Println(err.Error())
				return c.RenderError(err)
			}

			revel.INFO.Printf("Extracted file => %v", fpath)
		}
	}

	return c.RenderText("Successfully uploaded!");
}

func (c App) StartGame() revel.Result {
	gameID := uuid.NewV4().String()

	revel.INFO.Println("[app::StartGame] Starting new game instance id = " + gameID)
	app.InstanceCount++
	
	go func() {
		cmd := exec.Command(fmt.Sprintf("./%v", executableName), "-batchmode", "-nographics", "-logfile", "log.txt")
		cmd.Dir = fmt.Sprintf("%v/%v.app/Contents/MacOS/", buildPath, executableName)
		_, err := cmd.Output()
		if err != nil {
			revel.ERROR.Println(err)
		}
		app.InstanceCount--
		revel.INFO.Println("[app::StartGame] Closed game instance id = " + gameID)
	}()
	
	data := make(map[string]interface{})
	data["error"] = nil
	res := StartGameResponse{
		GameID:         gameID,
		ServerHostname: "localhost",
		ServerPort:     app.PortManager.GetNext(),
	}
	data["data"] = res

	app.MatchCount++
	revel.INFO.Printf("[app::StartGame] MatchCount = %v, Running instances = %v", app.MatchCount, app.InstanceCount)

	return c.RenderJSON(data)
}
