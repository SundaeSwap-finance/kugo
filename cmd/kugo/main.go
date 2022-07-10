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

	var filters []kugo.Filter
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

	matches, err := client.Matches(ctx, filters...)
	if err != nil {
		return fmt.Errorf("failed to find matches: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(&matches)

	return nil
}
