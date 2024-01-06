package database

import (
	"context"
	"time"
)

func DeleteFlag(userid int, level int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func DeleteRunning(userid int, level int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}