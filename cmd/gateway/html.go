package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

type HtmlConfig struct {
	Enable        bool   `yaml:"enable"`
	AssetPath     string `yaml:"asset_path"`
	DefaultPort   int    `yaml:"port"`
	ListenAddress string `yaml:"listen_address"`
}

func runHtml(cfg HtmlConfig) {
	if !cfg.Enable {
		return
	}
	if cfg.AssetPath == "" {
		cfg.AssetPath = "."
	}
	router := gin.Default()
	router.LoadHTMLGlob(cfg.AssetPath + "/templates/*.tmpl")
	router.Static("/assets", cfg.AssetPath+"/assets")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "DAIKIN Heat pump",
			"temperatures": p1p2.Sys.Temperatures,
			"valves":       p1p2.Sys.Valves,
			"pumps":        p1p2.Sys.Pumps,
			"status":       p1p2.Sys.Status,
		})
	})
	rest := router.Group("/api/v1")
	rest.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, p1p2.Sys)
	})
	rest.GET("/temperatures", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, p1p2.Sys.Temperatures)
	})
	rest.GET("/valves", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, p1p2.Sys.Valves)
	})
	rest.GET("/pumps", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, p1p2.Sys.Pumps)
	})
	rest.GET("/status", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, p1p2.Sys.Status)
	})

	rest.GET("/temperatures/:name", func(c *gin.Context) {
		for i := range p1p2.Sys.Temperatures {
			if p1p2.Sys.Temperatures[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, p1p2.Sys.Temperatures[i])

			}
		}
	})

	rest.GET("/valves/:name", func(c *gin.Context) {
		for i := range p1p2.Sys.Valves {
			if p1p2.Sys.Valves[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, p1p2.Sys.Valves[i])

			}
		}
	})
	rest.GET("/pumps/:name", func(c *gin.Context) {
		for i := range p1p2.Sys.Pumps {
			if p1p2.Sys.Pumps[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, p1p2.Sys.Pumps[i])

			}
		}
	})
	rest.GET("/status/:name", func(c *gin.Context) {
		for i := range p1p2.Sys.Status {
			if p1p2.Sys.Status[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, p1p2.Sys.Status[i])

			}
		}
	})
	if cfg.DefaultPort == 0 {
		cfg.DefaultPort = 8080
	}
	router.Run(fmt.Sprintf("%s:%d", cfg.ListenAddress, cfg.DefaultPort))
}
