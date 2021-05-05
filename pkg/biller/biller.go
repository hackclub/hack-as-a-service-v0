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
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func StartBilling(conn *dokku.DokkuConn) error {
	// Setup the event channel
	go func() {
		ch := dokku.CreateEventChannel()
		for {
			ev := <-ch
			if ev.Event == "post-deploy" {
				// Switch container ID
				var app db.App
				result := db.DB.First(&app, "short_name = ?", ev.AppName)
				if result.Error != nil {
					log.Printf(
						"Error while trying to rebill after post-deploy (%s): %+v\n",
						ev.AppName,
						result.Error,
					)
				} else {
					err := StartBillingApp(conn, app)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

	rows, err := db.DB.Model(&db.App{}).Rows()
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var app db.App
		err = db.DB.ScanRows(rows, &app)
		if err != nil {
			return err
		}

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

	// No container to bill
	if cid == "" {
		return nil
	}

	return StartBillerFromCID(app, team, cid)
}

func startBillerFromCID(app db.App, team db.Team, cid string) error {
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

func StartBillerFromCID(app db.App, team db.Team, cid string) error {
	if ch, ok := billerEventsChannels[app.ID]; ok {
		ch <- BillerEvent{NewContainerId: cid}
		return nil
	}

	return startBillerFromCID(app, team, cid)
}

type BillerEvent struct {
	NewContainerId string
	Delete         bool
}

var billerOutputs map[uint]map[chan decimal.Decimal]struct{} = make(map[uint]map[chan decimal.Decimal]struct{})

var statsOutputs map[uint]map[chan ProcessedOutput]struct{} = make(map[uint]map[chan ProcessedOutput]struct{})
var billerEventsChannels map[uint]chan BillerEvent = make(map[uint]chan BillerEvent)

func CreateStatsOutput(appId uint) chan ProcessedOutput {
	ch := make(chan ProcessedOutput)
	if _, ok := statsOutputs[appId]; !ok {
		statsOutputs[appId] = make(map[chan ProcessedOutput]struct{})
	}
	statsOutputs[appId][ch] = struct{}{}
	return ch
}

func RemoveStatsOutput(appId uint, ch chan ProcessedOutput) {
	if outputs, ok := statsOutputs[appId]; ok {
		delete(outputs, ch)
	}
}

func CreateBillerOutput(teamId uint) chan decimal.Decimal {
	ch := make(chan decimal.Decimal)
	if _, ok := billerOutputs[teamId]; !ok {
		billerOutputs[teamId] = make(map[chan decimal.Decimal]struct{})
	}
	billerOutputs[teamId][ch] = struct{}{}
	return ch
}

func RemoveBillerOutput(teamId uint, ch chan decimal.Decimal) {
	if outputs, ok := billerOutputs[teamId]; ok {
		delete(outputs, ch)
	}
}

func StopBiller(appId uint) {
	if ch, ok := billerEventsChannels[appId]; ok {
		ch <- BillerEvent{Delete: true}
	}
}

func biller(app db.App, team db.Team, stream types.ContainerStats) {
	lines := bufio.NewScanner(stream.Body)
	eventsCh := make(chan BillerEvent)
	billerEventsChannels[app.ID] = eventsCh

	defer stream.Body.Close()

	// Drop the first line since it contains bad data
	// FIXME: why?
	if !lines.Scan() {
		return
	}
loop:
	for lines.Scan() {
		line := lines.Text()
		// log.Printf("Got line: %s\n", line)
		var stat Stats
		if err := json.Unmarshal([]byte(line), &stat); err != nil {
			log.Printf("Error decoding json: %+v\n", err)
			break
		}

		output := stat.Process()
		// Notify listeners
		if outputs, ok := statsOutputs[app.ID]; ok {
			for ch := range outputs {
				ch <- output
			}
		}
		expense := output.price()
		if outputs, ok := billerOutputs[team.ID]; ok {
			for ch := range outputs {
				ch <- expense
			}
		}

		// Accrue this new cost to the team
		result := db.DB.Model(&team).Update("expenses", gorm.Expr("expenses + ?::decimal", expense))
		if result.Error != nil {
			log.Printf("Error: %+v\n", result.Error)
		}

		select {
		case ev := <-eventsCh:
			if ev.Delete {
				log.Printf("[biller][%s] Requested delete, bye\n", app.ShortName)
				break loop
			}
			if ev.NewContainerId != "" {
				err := startBillerFromCID(app, team, ev.NewContainerId)
				if err != nil {
					log.Printf("[biller][%s] Error while trying to migrate to new container ID %s, continuing with old one\n", app.ShortName, ev.NewContainerId)
					break
				}
				log.Printf("[biller][%s] Migrated to new container ID %s\n", app.ShortName, ev.NewContainerId)
				return
			}
		default:
		}
	}

	delete(billerEventsChannels, app.ID)
}
