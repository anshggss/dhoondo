package workerhelper

import (
	"fmt"
	"net/url"

	"github.com/anshggss/dhoondo/internal/models"
)

func ConvertToWebsite(links []string, val models.Website) ([]models.Website, error) {
	websites := []models.Website{}
	for _, v := range links {
		base, err := url.Parse(val.Link)
		if err != nil {
			return websites, err
		}
		href, err := url.Parse(v)
		if err != nil {
			return websites, err
		}
		resolved := base.ResolveReference(href)
		curLink := models.Website{
			Link:  resolved.String(),
			Depth: val.Depth + 1,
		}
		websites = append(websites, curLink)
		fmt.Printf("Added %s to publishing queue\n", curLink.Link)
	}
	return websites, nil

}
