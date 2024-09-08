package prisma

import "github.com/CollabTed/CollabTed-Backend/prisma/db"

var (
	Client *db.PrismaClient
)

func Connect() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	Client = client
}
