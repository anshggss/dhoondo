package workerhelper

import (
	"context"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func AddWebsiteToDatabase(val models.Website, q *database.Queries) (uuid.UUID, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return id, err
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
		return id, err
	}
	return id, nil
}

func AddWordsToDatabase(words []string, q *database.Queries) ([]int, error) {
	var ids []int
	for _, word := range words {
		arg := pgtype.Text{
			String: word,
			Valid:  true,
		}
		var finid int
		id, err := q.FindWord(context.Background(), arg)
		if err != nil {
			nid, err := q.InsertWord(context.Background(), arg)
			if err != nil {
				return nil, err
			}
			finid = int(nid)
		} else {
			finid = int(id)
		}
		ids = append(ids, finid)
	}
	return ids, nil
}

func IndexWords(ids []int, webId uuid.UUID, q *database.Queries) error {
	for _, id := range ids {
		arg := database.IndexWordParams{
			Wordid: pgtype.Int4{
				Int32: int32(id),
				Valid: true,
			},
			Websiteid: pgtype.UUID{
				Bytes: [16]byte(webId),
				Valid: true,
			},
		}
		err := q.IndexWord(context.Background(), arg)
		if err != nil {
			return err
		}
	}

	return nil
}
