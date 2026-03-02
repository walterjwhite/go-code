package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"github.com/walterjwhite/go-code/lib/events"
	"github.com/walterjwhite/go-code/lib/events/service"
	"github.com/walterjwhite/go-code/lib/events/transport"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/serialization"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable not set")
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close Pub/Sub client: %v", err)
		}
	}()

	serializer := serialization.NewJSONSerializer()
	compressor := zstd.NewCompressor()

	publisher := transport.NewGooglePubSubPublisher(client, serializer, compressor, nil, false, false)
	subscriber := transport.NewGooglePubSubSubscriber(client, serializer, compressor, nil, false, false)

	eventSvc := service.NewEventService(publisher, subscriber, serializer)
	defer func() {
		if err := eventSvc.Close(); err != nil {
			log.Printf("Failed to close event service: %v", err)
		}
	}()

	defineAndRegisterEvents(eventSvc)

	publishExampleEvents(ctx, eventSvc)

	sendResponses(ctx, eventSvc)

}

func defineAndRegisterEvents(svc *service.EventService) {
	log.Println("Registering events...")

	event1 := &events.Event{
		EventID: 1,
		Details: "Server ran out of memory",
		SupportedActions: []events.Action{
			{
				ActionID:     1,
				Message:      "Reboot Server",
				SupportsArgs: false,
			},
			{
				ActionID:     2,
				Message:      "Restart Process",
				SupportsArgs: true,
			},
			{
				ActionID:     3,
				Message:      "Ignore / Acknowledge",
				SupportsArgs: false,
			},
		},
	}

	if err := svc.RegisterEvent(event1); err != nil {
		log.Printf("Failed to register event 1: %v", err)
	} else {
		log.Println("✓ Event 1 registered")
	}

	event2 := &events.Event{
		EventID: 2,
		Details: "Build completed",
		SupportedActions: []events.Action{
			{
				ActionID:     1,
				Message:      "Start App",
				SupportsArgs: false,
			},
			{
				ActionID:     2,
				Message:      "Run Command",
				SupportsArgs: true,
			},
			{
				ActionID:     3,
				Message:      "Ignore / Acknowledge",
				SupportsArgs: false,
			},
		},
	}

	if err := svc.RegisterEvent(event2); err != nil {
		log.Printf("Failed to register event 2: %v", err)
	} else {
		log.Println("✓ Event 2 registered")
	}
}

func publishExampleEvents(ctx context.Context, svc *service.EventService) {
	log.Println("\nPublishing example events...")

	event1 := &events.Event{
		EventID: 1,
		Details: "Server CPU usage exceeded 95%",
		SupportedActions: []events.Action{
			{ActionID: 1, Message: "Reboot Server", SupportsArgs: false},
			{ActionID: 2, Message: "Restart Process", SupportsArgs: true},
			{ActionID: 3, Message: "Ignore", SupportsArgs: false},
		},
	}

	if err := svc.PublishEvent(ctx, "events", event1); err != nil {
		log.Printf("Failed to publish event 1: %v", err)
	} else {
		log.Println("✓ Event 1 published")
	}

	event2 := &events.Event{
		EventID: 2,
		Details: "Build completed successfully",
		SupportedActions: []events.Action{
			{ActionID: 1, Message: "Start App", SupportsArgs: false},
			{ActionID: 2, Message: "Run Command", SupportsArgs: true},
			{ActionID: 3, Message: "Ignore", SupportsArgs: false},
		},
	}

	if err := svc.PublishEvent(ctx, "events", event2); err != nil {
		log.Printf("Failed to publish event 2: %v", err)
	} else {
		log.Println("✓ Event 2 published")
	}
}

func sendResponses(ctx context.Context, svc *service.EventService) {
	log.Println("\nSending example responses...")

	response1 := &events.Response{
		EventID:  1,
		ActionID: 1,
	}

	if err := svc.PublishResponse(ctx, "responses", response1); err != nil {
		log.Printf("Failed to publish response 1: %v", err)
	} else {
		log.Printf("✓ Response 1 sent: Event %d -> Action %d\n", response1.EventID, response1.ActionID)
	}

	response2 := &events.Response{
		EventID:  2,
		ActionID: 2,
		Args:     []string{"git", "push"},
	}

	if err := svc.PublishResponse(ctx, "responses", response2); err != nil {
		log.Printf("Failed to publish response 2: %v", err)
	} else {
		log.Printf("✓ Response 2 sent: Event %d -> Action %d with args: %v\n",
			response2.EventID, response2.ActionID, response2.Args)
	}

	invalidResponse := &events.Response{
		EventID:  1,
		ActionID: 1,
		Args:     []string{"should-fail"}, // Action 1 doesn't support args
	}

	if err := svc.PublishResponse(ctx, "responses", invalidResponse); err != nil {
		log.Printf("✓ Invalid response correctly rejected: %v\n", err)
	}
}
