package spanner

import (
	"os"
	"fmt"
	"log"
	"context"

	"cloud.google.com/go/spanner"
	//"github.com/low-latency-game-telemetry/utils"
)

func SpannerWriteDML(ctx context.Context, keyString, valueString string) error {

	gcpProjectId    := os.Getenv("GCP_PROJECT_ID")
	spannerInstance := os.Getenv("SPANNER_INSTANCE")
	spannerDatabase := os.Getenv("SPANNER_DATABASE")
	spannerTable    := os.Getenv("SPANNER_TABLE_GAME_TELEMETRY")

	connectionStr := fmt.Sprintf("projects/%v/instances/%v/databases/%v", gcpProjectId, spannerInstance, spannerDatabase)

	spannerClient, err := spanner.NewClient(ctx, connectionStr)
	if err != nil {
		return err
	}
	defer spannerClient.Close()

	//keyString, valueString := utils.FormatInterface(data)

	// Generate DML
	dml := fmt.Sprintf("INSERT %v (%v) VALUES (%v)", spannerTable, keyString, valueString)
	fmt.Printf("dml: %v\n", dml)

	_, err = spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: dml,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		log.Printf("%d record(s) inserted.\n", rowCount)
		return err
	})
	return err

}