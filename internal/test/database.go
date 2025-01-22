package test

import (
	"context"
	"os"
	"testing"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func ApplyMigrations(ctx context.Context, t *testing.T, databaseURL string) {
	t.Helper()

	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("../../../database/migrations/"),
		),
	)
	require.NoError(t, err, "failed to create atlas workdir")
	defer workdir.Close()

	atlasClient, err := atlasexec.NewClient(workdir.Path(), "atlas")
	require.NoError(t, err, "failed to create atlas client")

	res, err := atlasClient.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL: databaseURL,
	})
	require.NoError(t, err, "failed to apply migrations")
	t.Logf("Applied %d migrations\n", len(res.Applied))
}

func CreateDatabaseContainer(ctx context.Context, t *testing.T) {
	t.Helper()

	postgresContainer, err := postgres.Run(ctx, "postgres:15.1-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase("test_db"),
		postgres.WithSQLDriver("pgx"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "failed to create database test container")

	t.Cleanup(func() {
		err = testcontainers.TerminateContainer(postgresContainer)
		if err != nil {
			t.Logf("failed to terminate database container: %v\n", err)
		}
	})
}
