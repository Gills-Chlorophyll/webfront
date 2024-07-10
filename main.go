package main

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TypOfIndexPage uint8

const (
	Splash TypOfIndexPage = iota
	AboutAquaponics
	AboutJourney
	AboutJoinus
)

func HandleIndexPage(typOfPage TypOfIndexPage) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := []MoreInfo{}
		if typOfPage == AboutAquaponics {
			data = DataAboutAquaponics
		} else if typOfPage == AboutJourney {
			data = DataAboutJourney
		} else if typOfPage == AboutJoinus {
			data = DataAboutJoinus
		} else if typOfPage == Splash {
			data = DataSplash
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Title":     "Gills & Chlorophyll",
			"TemplData": data,
		})
	}
}

func HandleBlogPage(c *gin.Context) {
	data, ok := DiaryData[c.Param("idx")]
	if !ok {
		c.HTML(http.StatusOK, "blog.html", gin.H{
			"Title":    "Gills & Chlorophyll",
			"BlogData": nil,
		})
		return
	}
	c.HTML(http.StatusOK, "blog.html", gin.H{
		"Title":    "Gills & Chlorophyll",
		"BlogData": data,
		"NavData":  data.Nav,
	})
}

// isNotEmptyString: to be used in templates to see if the value is not empty string
func isNotEmptyString(a string) bool {
	return a != ""
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// r.Static("/images", fmt.Sprintf("%s/images/", dirStatic))
	// r.Static("/js", fmt.Sprintf("%s/js/", dirStatic))
	r.SetFuncMap(template.FuncMap{
		"notEmpty": isNotEmptyString,
	})

	r.LoadHTMLGlob("web/html/**/*")
	r.Static("/templates", "web/templates/")
	r.Static("/js", "web/js/")
	r.Static("/images", "web/images/")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": "aboutme",
		})
	})
	r.GET("", HandleIndexPage(Splash))
	r.GET("/about-aquaponics", HandleIndexPage(AboutAquaponics))
	r.GET("/about-journey", HandleIndexPage(AboutJourney))
	r.GET("/about-joinus", HandleIndexPage(AboutJoinus))
	r.GET("/dear-diary/:idx", HandleBlogPage)
	log.Fatal(r.Run(":8080"))
}
