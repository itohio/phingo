package types

import (
	"fmt"
	"log"
	"time"
)

func (p *Project) TotalProgress() float32 {
	var total float32
	for _, l := range p.Log {
		total += l.Progress
	}
	return total
}
func (p *Project) TotalDuration() time.Duration {
	var total time.Duration
	for _, l := range p.Log {
		total += time.Duration(l.Duration)
	}
	return total
}
func (p *Project) TotalPrice() *Price {
	total := &Price{
		Denom: p.Denom(),
	}
	for _, l := range p.Log {
		pr := l.Price(p)
		if pr == nil {
			log.Fatalln("failed to parse price id=", p.Id)
		}
		total.Amount += pr.Amount
	}
	return total
}
func (p *Project) Denom() string {
	switch val := p.Rate.(type) {
	case *Project_Hourly:
		return val.Hourly.Denom
	case *Project_Total:
		return val.Total.Denom
	}
	return ""
}

func (p *Project) RateString() string {
	switch val := p.Rate.(type) {
	case *Project_Hourly:
		return fmt.Sprintf("%s/hr", val.Hourly.Pretty())
	case *Project_Total:
		return val.Total.Pretty()
	}
	return "-"
}

func (p *Project) SetRate(amount float32, denom string, hourly bool) {
	if hourly {
		p.Rate = &Project_Hourly{
			Hourly: &Price{
				Amount: amount,
				Denom:  denom,
			},
		}
	} else {
		p.Rate = &Project_Total{
			Total: &Price{
				Amount: amount,
				Denom:  denom,
			},
		}
	}
}
func (pl *Project_LogEntry) Price(p *Project) *Price {
	switch val := p.Rate.(type) {
	case *Project_Hourly:
		return &Price{
			Denom:  val.Hourly.Denom,
			Amount: float32(time.Duration(pl.Duration).Hours()) * val.Hourly.Amount,
		}
	case *Project_Total:
		return &Price{
			Denom:  val.Total.Denom,
			Amount: pl.Progress * val.Total.Amount,
		}
	}
	return nil
}
