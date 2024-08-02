package main

import "html/template"

type GalleryImage struct {
	Src      string
	Header   string
	Caption  template.HTML
	IsActive bool
}
type Gallery []GalleryImage

var (
	imageGallery = Gallery{
		{Header: "Tomatoes", IsActive: true, Src: "camera/tomato_1.jpg", Caption: `Tomato plants in their full glory.`},
		{Header: "Tomatoes", IsActive: false, Src: "camera/tomato_2.jpg", Caption: `Dense web of tomato foliage.`},
		{Header: "Tomatoes", IsActive: false, Src: "camera/tomato_3.jpg", Caption: `I guess the system is doing quite well.`},
		{Header: "Tomatoes", IsActive: false, Src: "camera/tomato_4.jpg", Caption: `Sun kissed unripe tomatoes.`},
		{Header: "Siphon", IsActive: false, Src: "camera/siphon_drain.jpg", Caption: `Siphon works well keeps the flood-drain cycle working`},
		{Header: "Random", IsActive: false, Src: "camera/rainy_afternoon.jpg", Caption: `One rainy afternoon in May of 2024`},
		{Header: "Chillies", IsActive: false, Src: "camera/chilly_ripe.jpg", Caption: `Ripening chilles boosted our confidence much.`},
		{Header: "Automation", IsActive: false, Src: "camera/automation.jpg", Caption: `Cloud connected motor control automation.`},
		{Header: "Water test", IsActive: false, Src: "camera/water_params.jpg", Caption: `Water composition measurements, API test kit`},
		{Header: "Aquaponics", IsActive: false, Src: "camera/tomato_fish.jpg", Caption: `System in full bloom`},
	}
)
