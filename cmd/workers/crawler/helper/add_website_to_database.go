package workerhelper

import (
	"context"
	"fmt"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func AddWebsiteToDatabase(val models.Website, q *database.Queries) error {

	id, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	arg := database.InsertWebsiteParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		Url: pgtype.Text{
			String: val.Link,
			Valid:  true,
		},
	}
	err = q.InsertWebsite(context.Background(), arg)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
