package helpers

import (
	"math/rand"
	"time"
)

type CoverPhoto struct {
	Source       string
	Photographer string
	CreditURL    string
}

func GetCover() CoverPhoto {
	rand.Seed(time.Now().UnixNano())
	return images[rand.Intn(len(images))]
}

var images = []CoverPhoto{
	CoverPhoto{
		Source:       "LibraryAmsterdam.jpg",
		Photographer: "Antonio Molinari",
		CreditURL:    "https://unsplash.com/@amolinari?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "LibraryKiev.jpg",
		Photographer: "Mariia Zakatiura",
		CreditURL:    "https://unsplash.com/@mzakatiura?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "LibraryLondon.jpg",
		Photographer: "Francesca Grima",
		CreditURL:    "https://unsplash.com/@francescagrima?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "LibrarySeattle.jpg",
		Photographer: "Sylvia Yang",
		CreditURL:    "https://unsplash.com/@sylviasyang?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "LibraryStuttgart.jpg",
		Photographer: "Gabriel Sollmann",
		CreditURL:    "https://unsplash.com/@gabons?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "LibraryWhiteGold.jpg",
		Photographer: "Valdemaras D.",
		CreditURL:    "https://unsplash.com/@deko_lt?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "BooksLampPitcher.jpg",
		Photographer: "Jez Timms",
		CreditURL:    "https://unsplash.com/@jeztimms?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
	CoverPhoto{
		Source:       "Archways.jpg",
		Photographer: "J Zamora",
		CreditURL:    "https://unsplash.com/@jzamora?utm_source=unsplash&amp;utm_medium=referral&amp;utm_content=creditCopyText",
	},
}
