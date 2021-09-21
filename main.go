package main

import (
	"context"
	"net/http"
	"tehranifar/fflow/follower"
	"tehranifar/fflow/storage"

	flowClient "github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO LIST
// ---------
//
//
// - fix client: rpc error: code = ResourceExhausted desc = grpc: received message larger than max (7902438 vs. 4194304)
// - event listener (from tx list) + metrics

func main() {
	ctx := context.Background()

	storage, err := storage.NewSqliteStorage()
	if err != nil {
		panic(err)
	}

	client, err := flowClient.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	startBlock, err := client.GetLatestBlock(ctx, true)
	if err != nil {
		panic(err)
	}

	f := follower.New(ctx, client, storage)
	go f.Follow(startBlock)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
