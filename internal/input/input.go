package input

import (
	"context"
	"fmt"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

// This function takes input from the standard input and
// returns a website
func WebsiteHandler(args []string) []models.Website {
	var res []models.Website
	for _, arg := range args {
		if validateUrl(arg) {
			res = append(res, models.Website{
				Link:  arg,
				Depth: 1,
			})
		} else {
			fmt.Printf("%s    Not a valid url\n", arg)
		}
	}
	return res
}

func QueryHandler(query []string, dbq *database.Queries) []string {
	var res []string
	for _, q := range query {
		arg := pgtype.Text{
			String: q,
			Valid:  true,
		}
		matchedWeb, err := dbq.QueryIndex(context.Background(), arg)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		for _, word := range matchedWeb {
			if word.Valid {
				res = append(res, word.String)
			}
		}

	}
	return res
}

func PrintHelp() {
	fmt.Printf("No command specified\n")
	fmt.Println("Use index to index a website")
	fmt.Println("Use query to query")
}
