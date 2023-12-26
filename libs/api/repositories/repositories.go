package api_repositories

import "github.com/jmoiron/sqlx"

type Repositories struct {
	Activity ActivityRepository
	Group    GroupRepository
}

func SetupRepositories(db *sqlx.DB) Repositories {
	return Repositories{Activity: ActivityRepository{DB: db}, Group: GroupRepository{DB: db}}
}
