package cmd

import (
	"encoding/json"
	rest "github.com/Roma7-7-7/workshops/topcoin/internal/repository"
	"github.com/Roma7-7-7/workshops/topcoin/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Provide top coin http server",
	Run: func(cmd *cobra.Command, args []string) {
		repo := rest.NewRepository(*CoinMarketApiKey)
		server := &Server{
			service: service.NewService(&repo),
		}

		http.HandleFunc("/topcoin", server.topCoinHandler)

		http.ListenAndServe(":8080", nil)
	},
}

type Server struct {
	service *service.Service
}

func (s *Server) topCoinHandler(writer http.ResponseWriter, request *http.Request) {
	accept := parseAccept(request.Header.Get("Accept"))
	if accept == "" {
		writer.WriteHeader(http.StatusNotAcceptable)
		_, _ = writer.Write([]byte("Provides only application/json and text/csv"))
		return
	}
	l := request.URL.Query().Get("limit")
	if l == "" {
		l = "10"
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("limit must be a number"))
		return
	}

	coins, err := s.service.GetTopCoin(limit)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if accept == "application/json" {
		err := json.NewEncoder(writer).Encode(coins)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		return
	}

	err = writeCsv(writer, "\t", coins)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/csv")
}

func parseAccept(val string) string {
	if val == "" {
		return "application/json"
	}

	parts := strings.Split(val, ",")
	for _, part := range parts {
		if part == "text/csv" {
			return "text/csv"
		}
		if part == "application/json" {
			return "application/json"
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
