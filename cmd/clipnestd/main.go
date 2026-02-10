package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clipnest/internal/clipboard"
	"clipnest/internal/config"
	"clipnest/internal/socket"
	"clipnest/internal/storage"
)

func main() {
	cfg := config.DefaultConfig()

	store, err := storage.NewStorage(cfg.MaxMemoryClips)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create storage: %v\n", err)
		os.Exit(1)
	}

	var server *socket.Server

	// Command handler: dispatches incoming commands from CLI clients
	handler := func(conn net.Conn, msg socket.SocketMessage) {
		switch msg.Type {
		case "list":
			limit := 20
			if msg.Data != nil {
				if m, ok := msg.Data.(map[string]interface{}); ok {
					if l, ok := m["limit"].(float64); ok && l > 0 {
						limit = int(l)
					}
				}
			}
			clips, _ := store.List(limit)
			sendClipList(conn, clips)

		case "search":
			query := ""
			limit := 20
			if m, ok := msg.Data.(map[string]interface{}); ok {
				if q, ok := m["query"].(string); ok {
					query = q
				}
				if l, ok := m["limit"].(float64); ok && l > 0 {
					limit = int(l)
				}
			}
			clips, _ := store.Search(query, limit)
			sendClipList(conn, clips)

		case "pins":
			clips, _ := store.GetPinned()
			sendClipList(conn, clips)

		case "copy_clip":
			id := extractID(msg)
			if id == 0 {
				sendError(conn, "missing clip id")
				return
			}
			clip, err := store.Get(id)
			if err != nil {
				sendError(conn, err.Error())
				return
			}
			if err := clipboard.Copy(clip.Content); err != nil {
				sendError(conn, fmt.Sprintf("failed to copy: %v", err))
				return
			}
			sendOK(conn)

		case "pin":
			id := extractID(msg)
			if id == 0 {
				sendError(conn, "missing clip id")
				return
			}
			if err := store.Pin(id); err != nil {
				sendError(conn, err.Error())
				return
			}
			sendOK(conn)

		case "unpin":
			id := extractID(msg)
			if id == 0 {
				sendError(conn, "missing clip id")
				return
			}
			if err := store.Unpin(id); err != nil {
				sendError(conn, err.Error())
				return
			}
			sendOK(conn)

		case "clear":
			_ = store.Clear()
			sendOK(conn)

		default:
			sendError(conn, fmt.Sprintf("unknown command: %s", msg.Type))
		}
	}

	server, err = socket.NewServer(cfg.SocketPath, handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start socket server: %v\n", err)
		os.Exit(1)
	}

	// Clipboard monitor: detects changes and stores them
	monitor := clipboard.NewMonitor(500*time.Millisecond, func(content, clipType string) {
		clip := storage.Clip{
			Content:   content,
			Type:      clipType,
			Timestamp: time.Now(),
		}
		id, err := store.Add(clip)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to store clip: %v\n", err)
			return
		}

		// Broadcast to connected clients
		stored, _ := store.Get(id)
		_ = server.Broadcast(socket.SocketMessage{
			Type: "new_clip",
			Data: clipToData(stored),
		})
	})
	monitor.Start()

	fmt.Printf("clipnestd running (socket: %s, max clips: %d)\n", cfg.SocketPath, cfg.MaxMemoryClips)

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nShutting down...")
	monitor.Stop()
	server.Close()
	store.Close()
}

func extractID(msg socket.SocketMessage) int64 {
	if msg.Data == nil {
		return 0
	}
	m, ok := msg.Data.(map[string]interface{})
	if !ok {
		return 0
	}
	id, ok := m["id"].(float64)
	if !ok {
		return 0
	}
	return int64(id)
}

func clipToData(c storage.Clip) socket.ClipData {
	return socket.ClipData{
		ID:        c.ID,
		Content:   c.Content,
		Type:      c.Type,
		Timestamp: c.Timestamp.Unix(),
		Pinned:    c.Pinned,
	}
}

func sendClipList(conn net.Conn, clips []storage.Clip) {
	clipDatas := make([]socket.ClipData, len(clips))
	for i, c := range clips {
		clipDatas[i] = clipToData(c)
	}
	resp := socket.ResponseMessage{
		Success: true,
		Data: socket.ClipListData{
			Clips: clipDatas,
			Count: len(clipDatas),
		},
	}
	data, _ := json.Marshal(resp)
	_ = socket.SendMessage(conn, socket.SocketMessage{Type: "response", Data: json.RawMessage(data)})
}

func sendOK(conn net.Conn) {
	resp := socket.ResponseMessage{Success: true}
	data, _ := json.Marshal(resp)
	_ = socket.SendMessage(conn, socket.SocketMessage{Type: "response", Data: json.RawMessage(data)})
}

func sendError(conn net.Conn, errMsg string) {
	resp := socket.ResponseMessage{Success: false, Error: errMsg}
	data, _ := json.Marshal(resp)
	_ = socket.SendMessage(conn, socket.SocketMessage{Type: "response", Data: json.RawMessage(data)})
}
