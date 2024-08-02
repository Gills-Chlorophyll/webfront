package main

import (
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type TypOfIndexPage uint8

const (
	Splash TypOfIndexPage = iota
	AboutAquaponics
	AboutJourney
	AboutJoinus
)

const (
	ENLISTING_PER_PAGE = 3
	VULTR_S3_BUCKET    = "gillschlorophyll"
)

var (
	AWS_S3 *s3.S3
	// For each of the payload the footer content is populated once
	// For the prototype stage we are leaving it blank, but as we expand the social media footprint in we can populate this
	FOOTER = &Footer{
		FBLink:     "",
		GmailLink:  "",
		LinkedLink: "",
		GitLink:    "",
		Phone:      "",
		Address:    "",
		Email:      "",
	}
)

func init() {
	/* Aws session for vultr object storage */
	// S3 signed urls for images
	var accessKey, secretKey, endpoint, region string
	envs := map[string]*string{
		"S3_ACCESSKEY": &accessKey,
		"S3_SECRETKEY": &secretKey,
		"S3_ENDPOINT":  &endpoint,
		"S3_REGION":    &region,
	}
	for k, v := range envs {
		if os.Getenv(k) == "" {
			log.Panicf("Environment variable value not available %s", k)
			return
		} else {
			*v = os.Getenv(k)
		}
	}
	// Create a new session with the provided credentials
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Errorf("Failed to create new aws session, check params %s", err)
		panic("Could not reach object storage for static files")
	}

	AWS_S3 = s3.New(sess)
}

func main() {

	// key := "tomato_farming.png"

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// r.Static("/images", fmt.Sprintf("%s/images/", dirStatic))
	// r.Static("/js", fmt.Sprintf("%s/js/", dirStatic))
	r.SetFuncMap(template.FuncMap{
		"notEmpty":      isNotEmptyString,
		"countToRange":  countToRange,
		"presignImgURL": presignImageUrl,
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
	r.GET("", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Gills & Chlorophyll is an aquaponics project that attempts to change how we implement urban farming."), IndexPageContent(Splash), PageDispatch("index"))
	r.GET("/about-aquaponics", HndlPagePayload("Gills & Chlorophyll", "tomato_farming.png", "Aquaponics urban farming", "Aquaponics is a modern urban soil-less farming method that can help us grow food sustainably."), IndexPageContent(AboutAquaponics), PageDispatch("index"))
	r.GET("/about-journey", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Gills & Chlorophyll is an aquaponics project that attempts to change how we look at urban farming."), IndexPageContent(AboutJourney), PageDispatch("index"))
	r.GET("/about-joinus", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Gills & Chlorophyll is an aquaponics project that attempts to change how we look at urban farming."), IndexPageContent(AboutJoinus), PageDispatch("index"))

	r.GET("/dear-diary/", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Enlisting of all the monthly events chronologically"), DiaryIndexContent, PageDispatch("enlist-blogs"))
	r.GET("/dear-diary/:idx", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Blogs are a way we share our earned knowledge wit you."), BlogPageContent, PageDispatch("blog")) // blog page

	r.GET("/gallery/", HndlPagePayload("Gills & Chlorophyll", "garden_farm.png", "Gills & Chlorophyll project", "Gallery of all images from our site in Pune."), GalleryPageContent, PageDispatch("gallery"))

	log.Fatal(r.Run(":8080"))
}
