package types

import "strings"

const (
	ContactFullName                  = "Full Name"
	ContactName                      = "Name"
	ContactPersonalCode              = "Personal Code"
	ContactRegistrationNumber        = "Reg. #"
	ContactCompanyRegistrationNumber = "Company Reg. #"
	ContactVAT                       = "VAT"
	ContactAddress                   = "Address"
	ContactPhone                     = "Phone"
	ContactCellPhone                 = "Cell"
	ContactFax                       = "Cell"
	ContactEmail                     = "Email"
	ContactBankAccount               = "Bank Account"
	ContactWalletAddress             = "Wallet Address"
	ContactDirector                  = "Director"
)

func DefaultContactTypes() []string {
	return []string{
		ContactFullName,
		ContactName,
		ContactPersonalCode,
		ContactRegistrationNumber,
		ContactCompanyRegistrationNumber,
		ContactVAT,
		ContactAddress,
		ContactPhone,
		ContactCellPhone,
		ContactFax,
		ContactEmail,
		ContactBankAccount,
		ContactWalletAddress,
		ContactDirector,
	}
}

func ParseContact(out map[string]string, vals []string) {
	for _, c := range vals {
		kv := strings.SplitN(c, "=", 2)
		if len(kv) == 2 && kv[1] != "" {
			out[kv[0]] = kv[1]
		} else {
			delete(out, kv[0])
		}
	}
}
