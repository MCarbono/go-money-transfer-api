package dockertest

//dockertest package is a custom common package that wraps the use of the library called "testcontainer".
//This lib has the purpose to create a docker container so we can do end-to-end or integration tests.
//doc: https://golang.testcontainers.org/
import (
	"context"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	DbName     = "money-api-test"
	DbUser     = "money-api-test"
	DbPassword = "money-api-test"
)

func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase(DbName),
		postgres.WithUsername(DbUser),
		postgres.WithPassword(DbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}
	endpoint, err := postgresContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	hostPost := strings.Split(endpoint, ":")
	return &PostgresContainer{
		PostgresContainer: postgresContainer,
		Host:              hostPost[0],
		Port:              hostPost[1],
	}, nil
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	Host string
	Port string
}

// postgresContainer, err := postgres.RunContainer(ctx,
//     testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
//     postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
//     postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
//     postgres.WithDatabase(dbName),
//     postgres.WithUsername(dbUser),
//     postgres.WithPassword(dbPassword),
//     testcontainers.WithWaitStrategy(
//         wait.ForLog("database system is ready to accept connections").
//             WithOccurrence(2).
//             WithStartupTimeout(5*time.Second)),
// )
// if err != nil {
//     panic(err)
// }
