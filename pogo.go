// based on https://pkg.go.dev/go.mau.fi/whatsmeow#example-package

package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func menu(client *whatsmeow.Client) {
	target, _ := types.ParseJID("123456789@s.whatsapp.net")
	for {
		var selection string
		fmt.Print("Enter command (write help to list available commands): ")
		fmt.Scanln(&selection)

		switch selection {

		case "t", "target":
			fmt.Print("Enter the target's phone number (E.164 format): ")
			var targetNumber string
			fmt.Scanln(&targetNumber)
			targetNumber = strings.TrimSpace(strings.Trim(targetNumber, "+"))
			matched, _ := regexp.MatchString(`^\d{5,15}$`, targetNumber)
			if !matched {
				fmt.Println("Input must be a valid phone number (5 to 15 digits).")
				break
			}
			target, _ = types.ParseJID(targetNumber + "@s.whatsapp.net")
			fmt.Printf("Target set to %s\n\n", target)

		case "d", "devices":
			fmt.Printf("Query devices for %s...\n", target)
			jids, _ := client.GetUserDevicesContext(context.Background(), []types.JID{target})
			if len(jids) > 0 {

				devices := make([]uint16, len(jids))
				for i, jid := range jids {
					devices[i] = jid.Device
				}
				fmt.Printf("Device list for %s: %d\n\n", target, devices)
			} else {
				fmt.Printf("No devices found for %s\n\n", target)
			}

		case "p", "prekey":
			fmt.Printf("Query prekey bundle for main device: %s...\n", target)
			preKeyResp, err := client.DangerousInternals().FetchPreKeys(context.Background(), []types.JID{target})
			spew.Dump("%#v, error: %v", preKeyResp, err)
			fmt.Println()

		case "c", "combine":
			fmt.Printf("Query devices for %s...\n", target)
			jids, _ := client.GetUserDevicesContext(context.Background(), []types.JID{target})
			if len(jids) > 0 {
				fmt.Printf("Query prekey bundle for all [%d] devices: %s...\n", len(jids), target)
				preKeyResp, err := client.DangerousInternals().FetchPreKeys(context.Background(), jids)
				spew.Dump("%#v, error: %v", preKeyResp, err)
				fmt.Println()
			} else {
				fmt.Printf("No devices found for %s\n\n", target)
			}

		case "h", "help":
			fmt.Println("Available commands:")
			fmt.Println("\t(h)elp     -- Show this help message")
			fmt.Println("\t(t)arget   -- Update the current target number")
			fmt.Println("\t(d)evices  -- Display existing sessions (main and companion devices) for the target number")
			fmt.Println("\t(p)rekey   -- Retrieve a prekey bundle for the target number (main device only)")
			fmt.Println("\t(c)ombine  -- Retrieve prekey bundles for all existing sessions (main and companion devices) of the target number")
			fmt.Println("\t(e)xit     -- Exit the program")

		case "e", "exit":
			return
		}
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:session/session.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(500 * time.Millisecond)

	menu(client)

	client.Disconnect()
}
