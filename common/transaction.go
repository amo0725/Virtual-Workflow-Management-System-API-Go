package common

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionFunc func(context.Context, mongo.Session) error

func WithTransaction(ctx context.Context, client *mongo.Client, fn TransactionFunc) (err error) {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	opts := options.Transaction()

	return mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		err := session.StartTransaction(opts)
		if err != nil {
			return err
		}

		if err = fn(sessionContext, session); err != nil {
			if e := session.AbortTransaction(sessionContext); e != nil {
				logrus.Errorf("Aborting transaction failed: %v", e)
			}
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
}
