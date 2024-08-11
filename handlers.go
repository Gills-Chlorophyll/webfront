package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/* All the route handlers with specialized logic for each of the specialized routes */

// HndlPagePayload : middleware to fill in the header footer template content
// Will not change any part of content - that is for downstream handlers to fill in
// Any data sent to the template has 3 parts
//
// - open graph, title data: filled in by the HndlPagePayload
// - Content 	: to be specific to the data and the template used
// - Footer 	: external / social navigation links to be filled in by the HndlPagePayload
func HndlPagePayload(title, ogimg, ogtitle, ogdesc string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("payload", &PagePayload{
			OpnGraph: &OpenGraph{
				Title:   title,                  // page title
				OgImage: presignImageUrl(ogimg), // getting the temp signed url for image resource on object storage
				OgTitle: ogtitle,                // Open graph title
				OgDesc:  ogdesc,                 // Open graph description
			},
			Content: nil, // to be set by downstream handlers
			Foot:    FOOTER,
		})
	}
}

// PageDispatch : final function in the chain of handlers where it takes the payload, loads the content on the pagepayload object, then calls the page to be dispatched
func PageDispatch(pageName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pylod, exists := ctx.Get("payload")
		if !exists {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("something unusual went wrong on the server"))
			return
		}
		content, _ := ctx.Get("content")
		if content != nil {
			// if its nil it would carry over the nil value from the previous handler
			pylod.(*PagePayload).Content = content // get content from Content dispatch
		}
		ctx.HTML(http.StatusOK, fmt.Sprintf("%s.html", pageName), pylod) // its the pagepayload object that is eventually dispatched
	}
}

// IndexPageContent : Serves content to media-text pages of the type []MoreInfo
// But it depends on the type of the page
// Should be called between HndlPagePayload and PageDispatch
// picks up the correct content on the page to be loaded
func IndexPageContent(typOfPage TypOfIndexPage) gin.HandlerFunc {
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
		} // extensible - if there are more pages
		ctx.Set("content", data)
	}
}

// DiaryIndexContent : content insertable as paginated list of blogs/diary
// From the query params this the current page number requested translates to the exact items tobe displayed
// Should be called between HndlPagePayload and PageDispatch
func DiaryIndexContent(c *gin.Context) {
	/* ==============
	Getting the current page requested
	==============*/
	page := 1
	currPage := c.Query("page")
	if currPage != "" {
		val, _ := strconv.ParseInt(currPage, 10, 64)
		page = int(val)
	}
	result := DiaryData.Paginate(ENLISTING_PER_PAGE, page)
	if yes, herr := result.HasError(); yes {
		herr.ToHttpCtx(c)
		return
	}
	c.Set("content", result.Result)
}

func BlogPageContent(c *gin.Context) {
	data, err := DiaryData.SearchWith(c.Param("idx"))
	if err != nil {
		c.HTML(http.StatusOK, "blog.html", gin.H{
			"Title":    "Gills & Chlorophyll",
			"BlogData": nil,
		})
		return
	}
	c.Set("content", data)
}

func GalleryPageContent(c *gin.Context) {
	c.Set("content", imageGallery)
}

// FootprintCalcContent : Handles form for Co2 emissions and footprint calculations.
func FootprintCalcContent(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.Set("content", &Co2FtPrintParams{Vegeterian: "off", ElectricKwh: 0, FishFeedKgs: 0, PlantYeildKgs: 0, FishYeildKgs: 0, Emissions: 0, Footprint: 0})
	} else if c.Request.Method == "POST" {
		// do calculations and send across
		// but for now we just send some dummy values
		c.Request.ParseForm()
		fmt.Println(c.Request.PostForm)
		result := Co2FtPrintParams{}
		if err := c.ShouldBind(&result); err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("failed to bind form input tomodel")
			// will send back the relevant page when error
			c.HTML(http.StatusBadRequest, "400.html", gin.H{
				"err_msg": "failed reading the inputs, this happens when one or more inputs you privided is invalidated",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"kwh":   result.ElectricKwh,
			"feed":  result.FishFeedKgs,
			"fish":  result.FishYeildKgs,
			"plant": result.PlantYeildKgs,
		}).Debug("form data..")
		// First we get the emissions
		// footprint is emissions per kg of yeild
		// we have gotten this formula from chatgpt discussions
		result.Emissions = 0.5*result.ElectricKwh + 1.5*result.FishFeedKgs
		// this can be calculated for vegeterian as well, in that case only the plant yeild is considered
		if result.Vegeterian != "on" {
			if result.PlantYeildKgs != 0 {
				result.Footprint = result.Emissions / result.PlantYeildKgs
			} else {
				// No yield would mean not footprint calcuations.
				logrus.WithFields(logrus.Fields{
					"kwh":   result.ElectricKwh,
					"feed":  result.FishFeedKgs,
					"fish":  result.FishYeildKgs,
					"plant": result.PlantYeildKgs,
				}).Error("Zero yield condition")
				c.HTML(http.StatusBadRequest, "400.html", gin.H{
					"err_msg": fmt.Sprintf("Without any yield cannot calculate footprint, Emissions stand to be %f Kgs of Co2", result.Emissions),
				})
				return
			}
		} else {
			result.Footprint = result.Emissions / (result.FishYeildKgs + result.PlantYeildKgs)
		}
		c.Set("content", &result)
	}

}
