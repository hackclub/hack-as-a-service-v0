package biller

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hackclub/hack-as-a-service/pkg/db"
	"github.com/hackclub/hack-as-a-service/pkg/dokku"
	"gorm.io/gorm"
)

func StartBilling(conn *dokku.DokkuConn) error {
	rows, err := db.DB.Model(&db.App{}).Rows()
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var app db.App
		db.DB.ScanRows(rows, &app)

		err = StartBillingApp(conn, app)
		if err != nil {
			return err
		}
	}

	return nil
}

func StartBillingApp(conn *dokku.DokkuConn, app db.App) error {
	// start the client and bill every second, updating the record
	var team db.Team
	result := db.DB.First(&team, "id = ?", app.TeamID)
	if result.Error != nil {
		return result.Error
	}

	// Get the app's container ID
	cid, err := conn.RunCommand(context.Background(), []string{"haas:cid", app.ShortName})
	if err != nil {
		return err
	}

	cid = strings.TrimSpace(cid)

	// Initialize a Docker API client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	stats_stream, err := cli.ContainerStats(context.Background(), cid, true)
	if err != nil {
		return err
	}

	go biller(app, team, stats_stream)

	return nil
}

func biller(app db.App, team db.Team, stream types.ContainerStats) {
	lines := bufio.NewScanner(stream.Body)

	defer stream.Body.Close()

	// Drop the first line since it contains bad data
	// FIXME: why?
	if !lines.Scan() {
		return
	}
	for lines.Scan() {
		line := lines.Text()
		// log.Printf("Got line: %s\n", line)
		var stat Stats
		if err := json.Unmarshal([]byte(line), &stat); err != nil {
			log.Printf("Error decoding json: %+v\n", err)
			break
		}

		output := stat.Process()
		expense := output.price()

		// Accrue this new cost to the team
		result := db.DB.Model(&team).Update("expenses", gorm.Expr("expenses + ?::decimal", expense))
		if result.Error != nil {
			log.Printf("Error: %+v\n", result.Error)
		}
	}
}
