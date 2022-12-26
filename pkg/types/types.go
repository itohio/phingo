package types

//go:generate protoc --proto_path=../../proto --go_out=. --go_opt=paths=source_relative ../../proto/models.proto

type AccountsArr []*Account
type ClientsArr []*Client
type ProjectsArr []*Project
type InvoicesArr []*Invoice

func (arr AccountsArr) ById(id string) *Account {
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

func (arr AccountsArr) ByName(name string) *Account {
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

func (arr AccountsArr) ByDenom(denom string) *Account {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Denom == denom {
			return a
		}
	}
	return nil
}

func (arr ClientsArr) ById(id string) *Client {
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

func (arr ClientsArr) ByName(name string) *Client {
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

func (arr ClientsArr) ByAccount(acc string) *Client {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Account == acc {
			return a
		}
	}
	return nil
}

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

func (arr InvoicesArr) ById(id string) *Invoice {
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

func (arr InvoicesArr) ByDate(date string) *Invoice {
	for _, a := range arr {
		if a == nil {
			continue
		}
		if a.Date == date {
			return a
		}
	}
	return nil
}
