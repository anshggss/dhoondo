package input

import (
	"fmt"
	"github.com/anshggss/dhoondo/internal/models"
	"os"
)

// This function takes input from the standard input and
// returns a website
func TakeInput() []models.Website {
	var res []models.Website
	args := os.Args[1:]
	for _, arg := range args {
		if validateUrl(arg) {
			res = append(res, models.Website{
				Link: arg,
			})
		} else {
			fmt.Printf("%s    Not a valid url\n", arg)
		}
	}
	return res
}
