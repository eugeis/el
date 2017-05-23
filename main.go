package main

import (
	"encoding/json"
	"gee/lg"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/urfave/cli"
	"el/core"
)

var l = lg.NewLogger("EL")

func main() {
	app := cli.NewApp()
	app.Name = "el"
	app.Usage = "elasticsearch exporter"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration files, separated by ','",
			Value: "el.yml",
		}, cli.StringFlag{
			Name: "environments, e",
			Usage: "Environments (file suffixes) for configurations and properties files." +
				"This files - <file>_<suffix>.<ext> -will be loaded from same directory as a configuration, properties file.",
			Value: "",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"sc"},
			Usage:   "Start El server",
			Action: func(c *cli.Context) error {
				config, err := loadConfig(c)
				if err != nil {
					return err
				}
				if err != nil {
					return err
				}
				return start(config)
			},

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		l.Err("%s", err)
	}
}

func loadConfig(c *cli.Context) (*core.Config, error) {
	configFiles := strings.Split(c.GlobalString("c"), ",")
	environments := strings.Split(c.GlobalString("e"), ",")
	return core.LoadConfig(configFiles, environments)
}

func prepareDebug(config *core.Config) {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
		lg.Debug = false
	}
	return
}

func start(config *core.Config) (err error) {
	prepareDebug(config)
	controller := core.NewEl(config)
	defer controller.Close()
	engine := gin.Default()
	defineRoutes(engine, controller, config)
	return engine.Run(fmt.Sprintf(":%d", config.Port))
}

func defineRoutes(engine *gin.Engine, controller *core.El, config *core.Config) {
	engine.StaticFile("/help", "html/doc.html")
	engine.StaticFile("/", "html/doc.html")

	engine.GET("/export", func(c *gin.Context) {
		response(nil, c)
	})

	adminGroup := engine.Group("/admin")
	{
		adminGroup.GET("/reload", func(c *gin.Context) {
			config, err := config.Reload()
			if err == nil {
				controller.UpdateConfig(config)
			}
			response(err, c)
		})

		adminGroup.GET("/config", func(c *gin.Context) {
			c.Header("Content-Type", "application/json; charset=UTF-8")
			c.IndentedJSON(http.StatusOK, config)
		})

		adminGroup.GET("/shutdown", func(c *gin.Context) {
			panic(nil)
		})
	}
}

func queryInt(key string, c *gin.Context) (ret int) {
	if value := c.Query(key); value != "" {
		ret, _ = strconv.Atoi(value)
	}
	return ret
}

func queryFlag(key string, c *gin.Context) (ret bool) {
	if value, ok := c.GetQuery(key); ok {
		if value == "" {
			ret = true
		} else {
			ret, _ = strconv.ParseBool(value)
		}
	}
	return ret
}

func response(err error, c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	if err == nil {
		c.String(http.StatusOK, "{ \"ok\": true }")
	} else {
		jsonDesc, _ := json.Marshal(err.Error())
		c.String(http.StatusConflict, fmt.Sprintf("{ \"ok\": false, \"desc:\": %s }", jsonDesc))
	}
}
