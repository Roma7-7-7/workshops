package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Roma7-7-7/workshops/topcoin/api"
	rest "github.com/Roma7-7-7/workshops/topcoin/internal/repository"
	"github.com/Roma7-7-7/workshops/topcoin/internal/service"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List top crypto currencies",
	Run: func(cmd *cobra.Command, args []string) {
		repository := rest.NewRepository(*CoinMarketApiKey)
		service := service.NewService(&repository)

		if *output != "plain" && *output != "json" {
			fmt.Printf("output format \"%s\" is not supported", *output)
			return
		}

		coins, err := service.GetTopCoin(*limit)
		if err != nil {
			fmt.Println(err)
			return
		}
		render(coins)

	},
}

var limit *int
var output *string
var separator *string

func render(coins []*api.Coin) {
	switch *output {
	case "plain":
		err := writeCsv(os.Stdout, *separator, coins)
		if err != nil {
			log.Fatalf("failed to write csv: %v", err)
		}
	case "json":
		for _, coin := range coins {
			encoded, err := json.Marshal(coin)
			if err != nil {
				log.Fatalf("encoding coin to json: %v", err)
			}
			fmt.Println(string(encoded))
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	limit = listCmd.Flags().IntP("limit", "l", 10, "limit of top crypto currencies")
	output = listCmd.Flags().StringP("format", "f", "plain", "output format. acceptable values: plain, json")
	separator = listCmd.Flags().StringP("separator", "s", "\t", "separator for plain output")
}
