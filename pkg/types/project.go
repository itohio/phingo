package types

import (
	"fmt"
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
		total.Amount += pr.Amount
	}
	return total
}

func (p *Project) Denom() string {
	if p.Account == nil {
		return "-"
	}
	return p.Account.Denom
}

func (p *Project) RateString() string {
	switch val := p.Rate.(type) {
	case *Project_Hourly:
		return fmt.Sprintf("%f/hr", val.Hourly)
	case *Project_Total:
		return fmt.Sprint(val.Total)
	}
	return "-"
}

func (p *Project) SetRate(amount float32, hourly bool) {
	if hourly {
		p.Rate = &Project_Hourly{
			Hourly: amount,
		}
	} else {
		p.Rate = &Project_Total{
			Total: amount,
		}
	}
}
func (pl *Project_LogEntry) Price(p *Project) *Price {
	denom := p.Denom()
	switch val := p.Rate.(type) {
	case *Project_Hourly:
		return &Price{
			Denom:  denom,
			Amount: float32(time.Duration(pl.Duration).Hours()) * val.Hourly,
		}
	case *Project_Total:
		return &Price{
			Denom:  denom,
			Amount: pl.Progress * val.Total,
		}
	}
	return nil
}

func (a *Project) MakeId() string {
	return SanitizePath(a.Name)
}

type ProjectsArr []*Project

func (arr ProjectsArr) ById(id string) *Project {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (arr ProjectsArr) ByName(name string) *Project {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Name == name {
			return a
		}
	}
	return nil
}
