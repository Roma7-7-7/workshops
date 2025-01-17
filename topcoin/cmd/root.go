package cmd

import (
	"fmt"
	"github.com/Roma7-7-7/workshops/topcoin/api"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "topcoin",
	Short: "Top Coin workshop",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var CoinMarketApiKey *string

func init() {
	CoinMarketApiKey = rootCmd.PersistentFlags().String("cm-api-key", "", "CoinMarketCap API key")

	rootCmd.MarkPersistentFlagRequired("cm-api-key")
}

func writeCsv(w io.Writer, separator string, coins []*api.Coin) error {
	_, err := fmt.Fprintf(w, "Rank%sSymbol%sPrice USD%s\n", separator, separator, separator)
	if err != nil {
		return err
	}
	for _, coin := range coins {
		_, err := fmt.Fprintf(w, "%d%s%s%s%f%s\n", coin.Rank, separator, coin.Symbol, separator, *coin.PriceUSD, separator)
		if err != nil {
			return err
		}
	}
	return nil
}
