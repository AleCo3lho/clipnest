package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"clipnest/internal/config"
	"clipnest/internal/socket"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "version" {
		fmt.Printf("clipnest %s\n", version)
		return
	}

	client, err := socket.NewClient(config.GetSocketPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\nIs clipnestd running?\n", err)
		os.Exit(1)
	}
	defer client.Close()

	switch cmd {
	case "list":
		limit := 20
		if len(os.Args) > 2 {
			if l, err := strconv.Atoi(os.Args[2]); err == nil && l > 0 {
				limit = l
			}
		}
		sendAndPrintList(client, socket.SocketMessage{
			Type: "list",
			Data: map[string]interface{}{"limit": limit},
		})

	case "search":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: clipnest search <query>")
			os.Exit(1)
		}
		sendAndPrintList(client, socket.SocketMessage{
			Type: "search",
			Data: map[string]interface{}{"query": os.Args[2], "limit": 20},
		})

	case "pins":
		sendAndPrintList(client, socket.SocketMessage{Type: "pins"})

	case "copy":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: clipnest copy <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid id")
			os.Exit(1)
		}
		sendAndPrintStatus(client, socket.SocketMessage{
			Type: "copy_clip",
			Data: map[string]interface{}{"id": id},
		})

	case "pin":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: clipnest pin <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid id")
			os.Exit(1)
		}
		sendAndPrintStatus(client, socket.SocketMessage{
			Type: "pin",
			Data: map[string]interface{}{"id": id},
		})

	case "unpin":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: clipnest unpin <id>")
			os.Exit(1)
		}
		id, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: invalid id")
			os.Exit(1)
		}
		sendAndPrintStatus(client, socket.SocketMessage{
			Type: "unpin",
			Data: map[string]interface{}{"id": id},
		})

	case "clear":
		sendAndPrintStatus(client, socket.SocketMessage{Type: "clear"})

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func sendAndPrintList(client *socket.Client, msg socket.SocketMessage) {
	if err := client.Send(msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Receive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse response
	rawData, err := json.Marshal(resp.Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var respMsg socket.ResponseMessage
	if err := json.Unmarshal(rawData, &respMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if !respMsg.Success {
		fmt.Fprintf(os.Stderr, "Error: %s\n", respMsg.Error)
		os.Exit(1)
	}

	// Parse clip list data
	listData, err := json.Marshal(respMsg.Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var clipList socket.ClipListData
	if err := json.Unmarshal(listData, &clipList); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing clips: %v\n", err)
		os.Exit(1)
	}

	if len(clipList.Clips) == 0 {
		fmt.Println("No clips.")
		return
	}

	for _, clip := range clipList.Clips {
		pin := " "
		if clip.Pinned {
			pin = "*"
		}
		ts := time.Unix(clip.Timestamp, 0).Format("15:04:05")
		content := clip.Content
		if len(content) > 80 {
			content = content[:77] + "..."
		}
		// Replace newlines for single-line display
		for i, c := range content {
			if c == '\n' || c == '\r' {
				content = content[:i] + "\\n" + content[i+1:]
			}
		}
		fmt.Printf("[%s] %-4d %s  %s\n", pin, clip.ID, ts, content)
	}
}

func sendAndPrintStatus(client *socket.Client, msg socket.SocketMessage) {
	if err := client.Send(msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Receive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	rawData, err := json.Marshal(resp.Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var respMsg socket.ResponseMessage
	if err := json.Unmarshal(rawData, &respMsg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if !respMsg.Success {
		fmt.Fprintf(os.Stderr, "Error: %s\n", respMsg.Error)
		os.Exit(1)
	}

	fmt.Println("OK")
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: clipnest <command> [args]

Commands:
  list [limit]     List recent clips (default: 20)
  search <query>   Search clips by content
  copy <id>        Copy clip back to system clipboard
  pin <id>         Pin a clip (protect from eviction)
  unpin <id>       Unpin a clip
  pins             List pinned clips only
  clear            Clear all clips
  version          Show version
`)
}
