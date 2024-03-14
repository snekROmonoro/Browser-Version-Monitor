package global

import (
	"fmt"

	"github.com/snekROmonoro/Browser-Version-Monitor/db"
)

var DatabaseClient *db.PrismaClient

func InitDatabase() error {
	if DatabaseClient != nil {
		return fmt.Errorf("database already initialized")
	}

	DatabaseClient = db.NewClient()
	return DatabaseClient.Connect()
}

func CloseDatabase() error {
	if DatabaseClient == nil {
		return fmt.Errorf("database not initialized")
	}

	var ret = DatabaseClient.Disconnect()
	DatabaseClient = nil
	return ret
}
