package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func InitMysql() (db *gorm.DB, err error) {
	dsn := "root:Wxwklyxjtwdy666@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}

func main() {
	// 链接数据库
	db, err := InitMysql()
	if err := db.AutoMigrate(&Todo{}); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		panic(errors.New("sql connect failure "))
	}

	// 默认路由
	r := gin.Default()

	// 加载资源
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	toDoGroup := r.Group("v1")
	{
		// add
		toDoGroup.POST("/todo", func(c *gin.Context) {
			var td Todo
			if err := c.BindJSON(&td); err != nil {
				fmt.Println(err)
			}
			if err := db.Create(&td).Error; err != nil {
				fmt.Println(err)
			}
			c.JSON(http.StatusOK, td)
		})

		// update
		toDoGroup.PUT("/todo/:id", func(c *gin.Context) {
			var td Todo
			if err := c.BindJSON(&td); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}

			res := db.Find(&Todo{}, &td)

			if res.Error != nil || res.RowsAffected != 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}
			db.Save(&td)
			c.JSON(http.StatusOK, td)
		})

		// view
		toDoGroup.GET("/todo", func(c *gin.Context) {
			tds := []Todo{}

			if err := db.Find(&tds).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
			}

			c.JSON(http.StatusOK, &tds)
		})

		// TODO:
		// need to Get by id
		toDoGroup.GET("/todo/:id", func(c *gin.Context) {
		})

		// remove
		toDoGroup.DELETE("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")

			if err := db.Where("id=?", id).Delete(&Todo{}).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}
			c.JSON(http.StatusOK, gin.H{"id": id})
		})
	}

	r.Run(":8080")
}
