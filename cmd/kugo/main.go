package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/SundaeSwap-finance/kugo"
	"github.com/SundaeSwap-finance/ogmigo/ouroboros/chainsync"
	"github.com/urfave/cli/v2"
)

var opts struct {
	Endpoint     string
	Spent        bool
	Unspent      bool
	Pattern      string
	PolicyID     string
	AssetName    string
	AssetNameHex string

	CreatedBefore uint64
	CreatedAfter  uint64
	SpentBefore   uint64
	SpentAfter    uint64
	Overlapping   uint64
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "endpoint",
			Usage:       "The Kupo Endpoint",
			Value:       "http://localhost:1442",
			EnvVars:     []string{"ENDPOINT"},
			Destination: &opts.Endpoint,
		},
		&cli.BoolFlag{
			Name:        "spent",
			Usage:       "Retrieve spent UTXOs only",
			Value:       false,
			Destination: &opts.Spent,
		},
		&cli.BoolFlag{
			Name:        "unspent",
			Usage:       "Retrieve unspent UTXOs only",
			Value:       false,
			EnvVars:     []string{"UNSPENT"},
			Destination: &opts.Unspent,
		},
		&cli.StringFlag{
			Name:        "pattern",
			Usage:       "Pattern to filter the address by",
			EnvVars:     []string{"PATTERN"},
			Destination: &opts.Pattern,
		},
		&cli.StringFlag{
			Name:        "policy-id",
			Usage:       "The policy ID to filter to",
			EnvVars:     []string{"POLICY_ID"},
			Destination: &opts.PolicyID,
		},
		&cli.StringFlag{
			Name:        "asset-name",
			Usage:       "The asset name to filter to",
			EnvVars:     []string{"ASSET_NAME"},
			Destination: &opts.AssetName,
		},
		&cli.StringFlag{
			Name:        "asset-name-hex",
			Usage:       "The hex encoded asset name to filter to",
			EnvVars:     []string{"ASSET_NAME_HEX"},
			Destination: &opts.AssetNameHex,
		},

		&cli.Uint64Flag{
			Name:        "created-before",
			Usage:       "Only print UTXOs that were created before a specific slot",
			EnvVars:     []string{"CREATED_BEFORE"},
			Destination: &opts.CreatedBefore,
		},
		&cli.Uint64Flag{
			Name:        "created-after",
			Usage:       "Only print UTXOs that were created after a specific slot",
			EnvVars:     []string{"CREATED_AFTER"},
			Destination: &opts.CreatedAfter,
		},
		&cli.Uint64Flag{
			Name:        "spent-before",
			Usage:       "Only print UTXOs that were spent, and spent before a specific slot",
			EnvVars:     []string{"SPENT_BEFORE"},
			Destination: &opts.SpentBefore,
		},
		&cli.Uint64Flag{
			Name:        "spent-after",
			Usage:       "Only print UTXOs that were spent, and spent after a specific slot",
			EnvVars:     []string{"SPENT_AFTER"},
			Destination: &opts.SpentAfter,
		},
		&cli.Uint64Flag{
			Name:        "overlapping",
			Usage:       "Only print UTXOs that 'overlap' a specific slot, i.e. were created before, or spent after",
			EnvVars:     []string{"OVERLAPPING"},
			Destination: &opts.Overlapping,
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func action(_ *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := kugo.New(
		kugo.WithEndpoint(opts.Endpoint),
	)

	var filters []kugo.MatchesFilter
	if opts.Unspent {
		filters = append(filters, kugo.OnlyUnspent())
	} else if opts.Spent {
		filters = append(filters, kugo.OnlySpent())
	}
	if opts.Pattern != "" {
		filters = append(filters, kugo.Pattern(opts.Pattern))
	}
	if opts.PolicyID != "" && opts.AssetName == "" && opts.AssetNameHex == "" {
		filters = append(filters, kugo.PolicyID(opts.PolicyID))
	} else if opts.PolicyID != "" {
		assetName := opts.AssetName
		if assetName == "" {
			assetNameBytes, err := hex.DecodeString(opts.AssetNameHex)
			if err != nil {
				return fmt.Errorf("invalid asset-name-hex %v: %w", opts.AssetNameHex, err)
			}
			assetName = string(assetNameBytes)
		}
		filters = append(filters, kugo.AssetID(chainsync.AssetID(fmt.Sprintf("%v.%v", opts.PolicyID, assetName))))
	}
	if opts.Overlapping > 0 {
		filters = append(filters, kugo.Overlapping(opts.Overlapping))
	}

	if opts.CreatedBefore > 0 {
		filters = append(filters, kugo.CreatedBefore(opts.CreatedBefore))
	} else if opts.CreatedAfter > 0 {
		filters = append(filters, kugo.CreatedAfter(opts.CreatedAfter))
	}

	if opts.SpentBefore > 0 {
		filters = append(filters, kugo.CreatedBefore(opts.SpentBefore))
	} else if opts.SpentAfter > 0 {
		filters = append(filters, kugo.CreatedAfter(opts.SpentAfter))
	}

	matches, err := client.Matches(ctx, filters...)
	if err != nil {
		return fmt.Errorf("failed to find matches: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(&matches)

	return nil
}
