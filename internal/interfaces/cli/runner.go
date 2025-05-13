package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"trading-bot/internal/application/usecase"
	"trading-bot/internal/domain/model"
	"trading-bot/internal/domain/service"
	"trading-bot/internal/infrastructure/exchange/foxbit"
)

// GOOS represents the operating system on which the program is running.
var GOOS string

// GOARCH represents the architecture of the operating system on which the program is running.
var GOARCH string

// CODEVERSION represents the version of the code.
var CODEVERSION string

// CODEBUILDDATE represents the date when the program was built.
var CODEBUILDDATE string

// CODEBUILDREVISION represents the revision of the code build.
var CODEBUILDREVISION string

// ExecuteCLI is the entry point for the command-line interface.
// It first parses any global flags (-info / -license), then
// dispatches on the sub-command.
func ExecuteCLI() {
	// --- GLOBAL FLAGS ---
	infoFlag := flag.Bool("info", false, "Display program compilation and version information")
	licenseFlag := flag.Bool("license", false, "Display program license information")
	flag.Parse()

	if *infoFlag {
		fmt.Printf("Version: %s\n", CODEVERSION)
		fmt.Printf("Operating System: %s\n", runtime.GOOS)
		fmt.Printf("System Architecture: %s\n", runtime.GOARCH)
		fmt.Printf("Build Date: %s\n", CODEBUILDDATE)
		fmt.Printf("Build Revision: %s\n", CODEBUILDREVISION)
		return
	}

	if *licenseFlag {
		fmt.Println("Copyright 2025 Isak Ruas")
		fmt.Println("Licensed under the Apache License, Version 2.0 (the 'License');")
		fmt.Println("you may not use this file except in compliance with the License.")
		fmt.Println("You may obtain a copy of the License at")
		fmt.Println("    http://www.apache.org/licenses/LICENSE-2.0")
		fmt.Println("Unless required by applicable law or agreed to in writing, software")
		fmt.Println("distributed under the License is distributed on an 'AS IS' BASIS,")
		fmt.Println("WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.")
		fmt.Println("See the License for the specific language governing permissions and")
		fmt.Println("limitations under the License.")
		return
	}

	// now handle subcommands
	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}
	cmd := args[0]

	switch cmd {
	case "help", "-h", "--help":
		usage()
		os.Exit(0)

	case "fetch-markets":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter: foxbit|binance|coinbase")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s fetch-markets [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		ex := mustInitExchange(*exch)
		mkts, err := (&usecase.FetchMarkets{Ex: ex}).Execute()
		if err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayMarkets(mkts)

	case "fetch-order-book":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter: foxbit|binance|coinbase")
		market := fs.String("market", "", "Market symbol (required), e.g. BTCBRL")
		depth := fs.Int("depth", 10, "Order book depth")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s fetch-order-book [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		if *market == "" {
			fmt.Fprintln(os.Stderr, "error: -market is required")
			fs.Usage()
			os.Exit(1)
		}
		ex := mustInitExchange(*exch)
		ob, err := (&usecase.FetchOrderBook{Ex: ex}).Execute(strings.ToUpper(*market), *depth)
		if err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayOrderBook(ob)

	case "place-order":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter: foxbit|binance|coinbase")
		market := fs.String("market", "", "Market symbol (required), e.g. BTCBRL")
		qty := fs.String("quantity", "", "Order quantity (required)")
		prc := fs.String("price", "", "Order price (required)")
		sideF := fs.String("side", "buy", "Order side: buy|sell")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s place-order [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		if *market == "" || *qty == "" || *prc == "" {
			fmt.Fprintln(os.Stderr, "error: -market, -quantity and -price are required")
			fs.Usage()
			os.Exit(1)
		}
		ex := mustInitExchange(*exch)
		var side model.OrderSide
		switch strings.ToLower(*sideF) {
		case "buy":
			side = model.Buy
		case "sell":
			side = model.Sell
		default:
			fmt.Fprintln(os.Stderr, "error: invalid side, use buy or sell")
			os.Exit(1)
		}
		order := model.Order{
			MarketSymbol: *market,
			Side:         side,
			Type:         model.Limit,
			Quantity:     *qty,
			Price:        *prc,
		}
		res, err := (&usecase.PlaceOrder{Ex: ex}).Execute(order)
		if err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayOrders([]model.Order{*res})

	case "cancel-order":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter")
		orderID := fs.String("order-id", "", "Order ID to cancel (required)")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s cancel-order [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		if *orderID == "" {
			fmt.Fprintln(os.Stderr, "error: -order-id is required")
			fs.Usage()
			os.Exit(1)
		}
		ex := mustInitExchange(*exch)
		if err := (&usecase.CancelOrder{Ex: ex}).Execute(*orderID); err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayCancel(*orderID)

	case "list-active-orders":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter")
		market := fs.String("market", "", "Market symbol (required)")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s list-active-orders [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		if *market == "" {
			fmt.Fprintln(os.Stderr, "error: -market is required")
			fs.Usage()
			os.Exit(1)
		}
		ex := mustInitExchange(*exch)
		act, err := (&usecase.ListActiveOrders{Ex: ex}).Execute(*market)
		if err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayOrders(act)

	case "get-order":
		fs := flag.NewFlagSet(cmd, flag.ExitOnError)
		exch := fs.String("exchange", "foxbit", "Exchange adapter")
		orderID := fs.String("order-id", "", "Order ID (required)")
		fs.Usage = func() {
			fmt.Fprintf(fs.Output(), "Usage: %s get-order [options]\n\n", os.Args[0])
			fmt.Fprintln(fs.Output(), "Options:")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])

		if *orderID == "" {
			fmt.Fprintln(os.Stderr, "error: -order-id is required")
			fs.Usage()
			os.Exit(1)
		}
		ex := mustInitExchange(*exch)
		o, err := (&usecase.GetOrder{Ex: ex}).Execute(*orderID)
		if err != nil {
			DisplayError(err)
			os.Exit(1)
		}
		DisplayOrders([]model.Order{*o})

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [global flags] <command> [options]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Global Flags:")
	fmt.Fprintln(os.Stderr, "  -info       Display build/version information")
	fmt.Fprintln(os.Stderr, "  -license    Display license information")
	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "  fetch-markets           List all markets")
	fmt.Fprintln(os.Stderr, "  fetch-order-book        Fetch order book for a market")
	fmt.Fprintln(os.Stderr, "  place-order             Place a new limit order")
	fmt.Fprintln(os.Stderr, "  cancel-order            Cancel an existing order")
	fmt.Fprintln(os.Stderr, "  list-active-orders      List active orders for a market")
	fmt.Fprintln(os.Stderr, "  get-order               Get details of a single order")
	fmt.Fprintln(os.Stderr, "\nUse “<command> --help” for more information about a command.")
}

func mustInitExchange(name string) service.Exchange {
	switch strings.ToLower(name) {
	case "foxbit":
		return foxbit.New(os.Getenv("FOXBIT_API_KEY"), os.Getenv("FOXBIT_API_SECRET"))
	default:
		log.Fatalf("Unknown exchange: %s", name)
		return nil
	}
}
