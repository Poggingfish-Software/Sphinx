package main

import (
	"io/ioutil"

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
		Url := c.Request.Header.Get("url")
		Desc := c.Request.Header.Get("desc")
		Category := c.Request.Header.Get("category")
		Key := c.Request.Header.Get("key")
		if Url == "" || Desc == "" || Category == "" || Key != string(key) {
			c.JSON(500, gin.H{
				"error": "Incorrect arguments!",
			})
		} else {
			if db.Select("*").Where("URL = ?", Url).Find(&Site{}).RowsAffected != 0 {
				c.JSON(500, gin.H{
					"error": "Site already exists on the index!",
				})
			} else {
				db.Table("sites").Create(&Site{URL: Url, Desc: Desc, Category: Category})
			}
		}
	})
	r.DELETE("/api", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		Url := c.Request.Header.Get("url")
		Key := c.Request.Header.Get("key")
		if Key != string(key) {
			c.JSON(500, gin.H{
				"error": "Incorrect arguments!",
			})
		} else {
			p := db.Table("sites").Where("URL = ?", Url).Delete("*").RowsAffected
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
