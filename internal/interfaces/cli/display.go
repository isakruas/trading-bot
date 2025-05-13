package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"trading-bot/internal/domain/model"
)

// DisplayMarkets prints a table of Market entries.
func DisplayMarkets(markets []model.Market) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SYMBOL\tPRICE_MIN\tPRICE_INCREMENT\tPRICE_PRECISION\tQUANTITY_MIN\tQUANTITY_INCREMENT\tQUANTITY_PRECISION")
	for _, m := range markets {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\t%d\n",
			m.Symbol, m.PriceMin, m.PriceIncrement, m.PricePrecision,
			m.QuantityMin, m.QuantityIncrement, m.QuantityPrecision,
		)
	}
	w.Flush()
}

// DisplayOrderBook prints bids and asks in tabular form.
func DisplayOrderBook(ob *model.OrderBook) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SIDE\tPRICE\tQUANTITY")
	for _, bid := range ob.Bids {
		p, _ := strconv.ParseFloat(bid[0], 64)
		q, _ := strconv.ParseFloat(bid[1], 64)
		fmt.Fprintf(w, "BID\t%.2f\t%.8f\n", p, q)
	}
	for _, ask := range ob.Asks {
		p, _ := strconv.ParseFloat(ask[0], 64)
		q, _ := strconv.ParseFloat(ask[1], 64)
		fmt.Fprintf(w, "ASK\t%.2f\t%.8f\n", p, q)
	}
	w.Flush()
}

// DisplayOrders prints a list of orders in tabular form.
func DisplayOrders(orders []model.Order) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tMARKET\tSIDE\tTYPE\tPRICE\tQUANTITY\tSTATE")
	for _, o := range orders {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			o.ID, o.MarketSymbol, string(o.Side), string(o.Type),
			o.Price, o.Quantity, o.State,
		)
	}
	w.Flush()
}

// DisplayCancel prints the result of a cancel operation in two columns.
func DisplayCancel(orderID string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ORDER_ID\tSTATUS")
	fmt.Fprintf(w, "%s\t%s\n", orderID, "CANCELLED")
	w.Flush()
}
