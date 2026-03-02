package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/pubsub/v2"
	"github.com/walterjwhite/go-code/lib/events"
	"github.com/walterjwhite/go-code/lib/events/service"
	"github.com/walterjwhite/go-code/lib/events/transport"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

const (
	commandList    = "list"
	commandRespond = "respond"
	commandListen  = "listen"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case commandList:
		handleList()
	case commandRespond:
		handleRespond()
	case commandListen:
		handleListen()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Event Client - Manage events and send responses

Usage:
  event-client list [--project PROJECT_ID] [--registry FILE]
    List all registered events

  event-client respond --event-id ID --action-id ID [--args ARG1,ARG2,...]
    [--project PROJECT_ID] [--topic TOPIC] [--use-encryption] [--encryption-key-file FILE]
    Send a response to an event

  event-client listen --subscription SUB [--project PROJECT_ID]
    [--use-encryption] [--encryption-key-file FILE]
    Listen for events on a subscription

Examples:
  event-client list --project my-project
  event-client respond --event-id 1 --action-id 1 --project my-project --topic responses
  event-client respond --event-id 2 --action-id 2 --args "git,push" --project my-project --topic responses
  event-client listen --subscription event-sub --project my-project`)
}

func handleList() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	projectID := fs.String("project", os.Getenv("GCP_PROJECT_ID"), "GCP Project ID")
	registryFile := fs.String("registry", "", "Optional registry file with event definitions")

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	if *projectID == "" {
		log.Fatal("Project ID is required. Set --project or GCP_PROJECT_ID environment variable")
	}

	ctx := context.Background()
	eventSvc := initEventService(ctx, *projectID, false, "")

	defer func() {
		if err := eventSvc.Close(); err != nil {
			log.Printf("Failed to close event service: %v", err)
		}
	}()

	if *registryFile != "" {
		loadEventsFromFile(*registryFile, eventSvc)
	}

	events := eventSvc.ListEvents()
	if len(events) == 0 {
		fmt.Println("No events registered")
		return
	}

	fmt.Println("\nRegistered Events:")
	fmt.Println("==================")
	for _, event := range events {
		fmt.Printf("\nEvent ID: %d\n", event.EventID)
		fmt.Printf("Details: %s\n", event.Details)
		fmt.Println("Supported Actions:")
		for _, action := range event.SupportedActions {
			args := "no"
			if action.SupportsArgs {
				args = "yes"
			}
			fmt.Printf("  - Action %d: %s (supports args: %s)\n", action.ActionID, action.Message, args)
		}
	}
}

func handleRespond() {
	fs := flag.NewFlagSet("respond", flag.ExitOnError)
	eventID := fs.Int("event-id", 0, "Event ID")
	actionID := fs.Int("action-id", 0, "Action ID")
	argsStr := fs.String("args", "", "Comma-separated arguments")
	projectID := fs.String("project", os.Getenv("GCP_PROJECT_ID"), "GCP Project ID")
	topic := fs.String("topic", "responses", "Pub/Sub topic to publish response to")
	useEncryption := fs.Bool("use-encryption", false, "Enable encryption")
	keyFile := fs.String("encryption-key-file", "", "Path to encryption key file")

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	if *eventID == 0 || *actionID == 0 {
		log.Fatal("Both --event-id and --action-id are required")
	}

	if *projectID == "" {
		log.Fatal("Project ID is required. Set --project or GCP_PROJECT_ID environment variable")
	}

	ctx := context.Background()
	eventSvc := initEventService(ctx, *projectID, *useEncryption, *keyFile)
	defer func() {
		if err := eventSvc.Close(); err != nil {
			log.Printf("Failed to close event service: %v", err)
		}
	}()

	var args []string
	if *argsStr != "" {
		rawArgs := strings.Split(*argsStr, ",")
		if len(rawArgs) > 100 {
			log.Fatal("Too many arguments: maximum 100 arguments allowed")
		}
		for _, arg := range rawArgs {
			trimmed := strings.TrimSpace(arg)
			if err := validateArgument(trimmed); err != nil {
				log.Fatalf("Invalid argument: %v", err)
			}
			args = append(args, trimmed)
		}
	}

	response := &events.Response{
		EventID:  *eventID,
		ActionID: *actionID,
		Args:     args,
	}

	if err := eventSvc.ValidateResponse(response); err != nil {
		log.Fatalf("Invalid response: %v", err)
	}

	if err := eventSvc.PublishResponse(ctx, *topic, response); err != nil {
		log.Fatalf("Failed to publish response: %v", err)
	}

	fmt.Printf("Response published successfully!\n")
	fmt.Printf("Event ID: %d\n", response.EventID)
	fmt.Printf("Action ID: %d\n", response.ActionID)
	if len(response.Args) > 0 {
		fmt.Printf("Args: %v\n", response.Args)
	}
}

func handleListen() {
	fs := flag.NewFlagSet("listen", flag.ExitOnError)
	subscription := fs.String("subscription", "", "Pub/Sub subscription to listen on")
	projectID := fs.String("project", os.Getenv("GCP_PROJECT_ID"), "GCP Project ID")
	useEncryption := fs.Bool("use-encryption", false, "Enable decryption")
	keyFile := fs.String("encryption-key-file", "", "Path to encryption key file")

	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	if *subscription == "" {
		log.Fatal("--subscription is required")
	}

	if *projectID == "" {
		log.Fatal("Project ID is required. Set --project or GCP_PROJECT_ID environment variable")
	}

	ctx := context.Background()
	eventSvc := initEventService(ctx, *projectID, *useEncryption, *keyFile)
	defer func() {
		if err := eventSvc.Close(); err != nil {
			log.Printf("Failed to close event service: %v", err)
		}
	}()

	fmt.Printf("Listening for events on subscription: %s\n", *subscription)

	if err := eventSvc.Subscribe(ctx, *subscription, func(event *events.Event) error {
		fmt.Printf("\nReceived Event:\n")
		fmt.Printf("Event ID: %d\n", event.EventID)
		fmt.Printf("Details: %s\n", event.Details)
		fmt.Println("Supported Actions:")
		for _, action := range event.SupportedActions {
			args := "no"
			if action.SupportsArgs {
				args = "yes"
			}
			fmt.Printf("  - Action %d: %s (supports args: %s)\n", action.ActionID, action.Message, args)
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to listen for events: %v", err)
	}
}

func initEventService(
	ctx context.Context,
	projectID string,
	useEncryption bool,
	keyFile string,
) *service.EventService {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}

	serializer := serialization.NewJSONSerializer()
	compressor := zstd.NewCompressor()

	var encryptor any
	if useEncryption && keyFile != "" {
		enc, err := aes.NewAESFromFile(keyFile)
		if err != nil {
			log.Fatalf("Failed to load encryption key: %v", err)
		}
		encryptor = enc
	}

	pub := transport.NewGooglePubSubPublisher(
		client,
		serializer,
		compressor,
		encryptor,
		false,         // enableCompression
		useEncryption, // enableEncryption
	)

	sub := transport.NewGooglePubSubSubscriber(
		client,
		serializer,
		compressor,
		encryptor,
		false,         // enableCompression
		useEncryption, // enableEncryption
	)

	return service.NewEventService(pub, sub, serializer)
}

func loadEventsFromFile(filename string, svc *service.EventService) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Failed to resolve file path: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	cleanAbsPath := filepath.Clean(absPath)
	cleanCwd := filepath.Clean(cwd)

	realPath, err := filepath.EvalSymlinks(cleanAbsPath)
	if err != nil {
		log.Fatalf("Failed to resolve symlinks: %v", err)
	}

	if !strings.HasPrefix(realPath, cleanCwd+string(filepath.Separator)) && realPath != cleanCwd {
		log.Fatalf("Invalid file path: file must be within the current working directory")
	}

	fileInfo, err := os.Stat(realPath)
	if err != nil {
		log.Fatalf("Failed to stat file: %v", err)
	}
	const maxFileSize = 10 * 1024 * 1024 // 10MB limit
	if fileInfo.Size() > maxFileSize {
		log.Fatalf("Registry file too large: %d bytes (max %d bytes)", fileInfo.Size(), maxFileSize)
	}

	data, err := os.ReadFile(realPath)
	if err != nil {
		log.Fatalf("Failed to read registry file: %v", err)
	}

	var eventsList []events.Event
	if err := json.Unmarshal(data, &eventsList); err != nil {
		log.Fatalf("Failed to parse registry file: %v", err)
	}

	const maxEvents = 10000
	if len(eventsList) > maxEvents {
		log.Fatalf("Too many events in registry: %d (max %d)", len(eventsList), maxEvents)
	}

	for i := range eventsList {
		if err := svc.RegisterEvent(&eventsList[i]); err != nil {
			log.Printf("Warning: Failed to register event %d: %v", eventsList[i].EventID, err)
		}
	}
}

func validateArgument(arg string) error {
	if len(arg) > 4096 {
		return fmt.Errorf("argument exceeds maximum length (4096 chars): %d", len(arg))
	}

	if strings.ContainsRune(arg, '\x00') {
		return fmt.Errorf("argument contains null bytes")
	}

	suspiciousPatterns := []string{"$(", "`", "&", "|", ";", ">", "<", "\n", "\r"}
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(arg, pattern) {
			return fmt.Errorf("argument contains suspicious pattern: %q", pattern)
		}
	}

	if strings.Contains(arg, "/") || strings.Contains(arg, "\\") {
		if arg != ".." && arg != "." && arg != "~" && !strings.HasPrefix(arg, "..") && !strings.HasPrefix(arg, "./") && !strings.HasPrefix(arg, "../") {
			absPath, err := filepath.Abs(arg)
			if err == nil {
				cwd, _ := os.Getwd()
				cleanAbsPath := filepath.Clean(absPath)
				cleanCwd := filepath.Clean(cwd)
				if !strings.HasPrefix(cleanAbsPath, cleanCwd+string(filepath.Separator)) && cleanAbsPath != cleanCwd {
					return fmt.Errorf("file path argument references outside current directory")
				}
			}
		}
	}

	return nil
}
