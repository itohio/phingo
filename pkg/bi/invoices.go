package bi

import (
	"strings"

	"github.com/itohio/phingo/pkg/types"
)

type InvoiceItem struct {
	Name   string
	Amount *types.Price
	Rate   float32
}

type InvoiceSummary struct {
	Subtotal *types.Price
	Discount *types.Price
	Tax      *types.Price
	Total    *types.Price

	Discounts []*InvoiceItem
	Taxes     []*InvoiceItem
}

func NewInvoiceSummary(inv *types.Invoice, denom string, decimals uint32) *InvoiceSummary {
	ret := &InvoiceSummary{
		Subtotal: SubTotal(inv.Items, denom, decimals),
	}

	ret.Discount = ret.Subtotal.Mul(SubtotalPercent(inv.Items, "Discount"))
	intermediate := ret.Subtotal.Sub(ret.Discount)
	ret.Tax = intermediate.Mul(SubtotalPercent(inv.Items, "Tax"))
	ret.Total = intermediate.Add(ret.Tax)

	for _, it := range inv.Items {
		if !it.Extra || it.Rate <= 0 {
			continue
		}

		switch {
		case strings.Contains(it.Name, "Tax"):
			ret.Taxes = append(ret.Taxes, &InvoiceItem{
				Name:   it.Name,
				Amount: intermediate.Mul(it.AdjustedRate()),
				Rate:   it.Rate,
			})
		case strings.Contains(it.Name, "Discount"):
			ret.Discounts = append(ret.Discounts, &InvoiceItem{
				Name:   it.Name,
				Amount: ret.Subtotal.Mul(it.AdjustedRate()),
				Rate:   it.Rate,
			})
		}
	}

	return ret
}

func SubTotal(arr []*types.Invoice_Item, denom string, decimals uint32) *types.Price {
	pr := &types.Price{
		Amount:   0,
		Denom:    denom,
		Decimals: decimals,
	}
	for _, it := range arr {
		if it.Extra {
			continue
		}
		pr.Amount += it.Amount * it.AdjustedRate()
	}
	return pr
}

func SubtotalPercent(arr []*types.Invoice_Item, name string) float32 {
	var rate float32
	for _, it := range arr {
		if !it.Extra {
			continue
		}
		if strings.Contains(it.Name, name) {
			rate += it.AdjustedRate()
		}
	}
	return rate
}

type InvoicePayments []*types.Invoice_Payment
