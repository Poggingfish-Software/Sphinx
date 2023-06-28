package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	URL      string
	Desc     string
	Category string
}
type Body struct {
	Url      string `json:"url"`
	Desc     string `json:"desc"`
	Category string `json:"category"`
	Key      string `json:"key"`
}
type RemBody struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func main() {
	key, err := ioutil.ReadFile(".env")
	if err != nil {
		panic("failed to open env")
	}
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Site{})
	r := gin.Default()
	r.LoadHTMLGlob("pages/*")
	r.GET("/", func(c *gin.Context) {
		q := new(int64)
		db.Select("*").Find(&Site{}).Count(q)
		c.HTML(200, "index.tmpl", gin.H{
			"links": q,
		})
	})
	r.GET("/index", func(c *gin.Context) {
		list := []Site{}
		out := ""
		db.Table("sites").Select("*").Find(&list)
		for _, i := range list {
			out += i.URL + " -- " + i.Desc + " -- " + i.Category + "\n"
		}
		c.String(200, out)
	})
	r.POST("/api", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		body := Body{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, "Incorrect arguments!")
			return
		}
		if body.Key != string(key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Incorrect arguments!",
			})
		} else {
			if db.Select("*").Where("URL = ?", body.Url).Find(&Site{}).RowsAffected != 0 {
				c.JSON(http.StatusNotModified, gin.H{
					"error": "Site already exists on the index!",
				})
			} else {
				db.Table("sites").Create(&Site{URL: (body.Url), Desc: body.Desc, Category: body.Category})
			}
		}
	})
	r.DELETE("/api", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		body := RemBody{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, "Incorrect arguments!")
			return
		}
		if body.Key != string(key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Incorrect arguments!",
			})
		} else {
			p := db.Table("sites").Where("URL = ?", body.Url).Delete("*").RowsAffected
			c.JSON(200, gin.H{
				"count": p,
			})
		}
	})
	r.OPTIONS("/api", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET")
		c.JSON(200, gin.H{})
	})
	r.Run("0.0.0.0:9332") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
