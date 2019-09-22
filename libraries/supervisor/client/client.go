package irpc

import (
	"context"
	"log"
	"net/rpc"
	"time"

	"github.com/walterjwhite/go-application/libraries/io/disk"
	"github.com/walterjwhite/go-application/libraries/os/unix/interfaces"
	"github.com/walterjwhite/go-application/libraries/os/unix/service"
	"github.com/walterjwhite/go-application/libraries/supervisor/data"
	"github.com/walterjwhite/go-application/libraries/supervisor/server"
)

const REFRESH_DURATION = 5 * time.Second

type Client struct {
	Rpc     *rpc.Client
	Context *context.Context
	Cancel  *context.CancelFunc
}

func New() *Client {
	client, err := rpc.DialHTTP(server.PROTOCOL, server.LISTEN_HOST+":"+server.PORT)
	if err != nil {
		log.Fatalf("Error dialing: %v\n", err)
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	return &Client{client, &context, &cancel}
}

func (client Client) Services() []data.Row {
	args := &server.Args{}
	var servicesResponse []service.Service
	_ = client.Rpc.Call("ServiceServer.Services", args, &servicesResponse)

	return service.Convert(servicesResponse)
}

func (client Client) BuildDateTime() string {
	args := &server.Args{}
	var buildDateTimeResponse string
	_ = client.Rpc.Call("Server.BuildDateTime", args, &buildDateTimeResponse)

	return buildDateTimeResponse
}

func (client Client) Uptime() string {
	args := &server.Args{}
	var uptimeResponse string
	_ = client.Rpc.Call("Server.Uptime", args, &uptimeResponse)

	return uptimeResponse
}

func (client Client) Logs() string {
	args := &server.Args{}
	var logResponse string
	_ = client.Rpc.Call("Server.Logs", args, &logResponse)

	return logResponse
}

func (client Client) Interfaces() []data.Row {
	args := &server.Args{}
	var interfacesResponse []interfaces.Interface
	_ = client.Rpc.Call("InterfaceServer.Interfaces", args, &interfacesResponse)

	return interfaces.Convert(interfacesResponse)
}

func (client Client) Disks() []data.Row {
	args := &server.Args{}
	var disksResponse []disk.Disk
	_ = client.Rpc.Call("DiskServer.Disks", args, &disksResponse)

	return disk.Convert(disksResponse)
}
