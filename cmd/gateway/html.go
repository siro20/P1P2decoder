package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

func runHtml(db *p1p2.DB, sys p1p2.System) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "DAIKIN Heat pump",
			"temperatures": sys.Temperatures,
			"valves":       sys.Valves,
			"pumps":        sys.Pumps,
			"status":       sys.Status,
		})
	})
	rest := router.Group("/api/v1")
	rest.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, sys)
	})
	rest.GET("/temperatures", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, sys.Temperatures)
	})
	rest.GET("/valves", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, sys.Valves)
	})
	rest.GET("/pumps", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, sys.Pumps)
	})
	rest.GET("/status", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, sys.Status)
	})

	rest.GET("/temperatures/:name", func(c *gin.Context) {
		for i := range sys.Temperatures {
			if sys.Temperatures[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, sys.Temperatures[i])

			}
		}
	})
	if db != nil {
		rest.GET("/temperatures/charts/:name", func(c *gin.Context) {
			for i := range sys.Temperatures {
				if sys.Temperatures[i].Name() == c.Params.ByName("name") {
					entries, err := db.GetTemperatures(*sys.Temperatures[i], time.Now().Add(-time.Hour*24))
					if err != nil {
						c.AbortWithError(http.StatusNotFound, fmt.Errorf("Unknown name"))
					}
					c.IndentedJSON(http.StatusOK, entries)
				}
			}
		})

	}

	rest.GET("/valves/:name", func(c *gin.Context) {
		for i := range sys.Valves {
			if sys.Valves[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, sys.Valves[i])

			}
		}
	})
	rest.GET("/pumps/:name", func(c *gin.Context) {
		for i := range sys.Pumps {
			if sys.Pumps[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, sys.Pumps[i])

			}
		}
	})
	rest.GET("/status/:name", func(c *gin.Context) {
		for i := range sys.Status {
			if sys.Status[i].Name() == c.Params.ByName("name") {
				c.IndentedJSON(http.StatusOK, sys.Status[i])

			}
		}
	})
	router.Run(":8080")
}
