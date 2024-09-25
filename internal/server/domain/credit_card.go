package domain

// BillingAddress represents the address associated with a billing account,
// including the country, street, city, and postcode.
type BillingAddress struct {
	Country  string
	Street   string
	City     string
	Postcode string
}

// CreditCard represents a credit card with associated metadata, name,
// expiration month and year, number, and billing address.
type CreditCard struct {
	Metadata       string
	Name           string
	Month          string
	Year           string
	Number         string
	BillingAddress BillingAddress
}

// NewBillingAddress creates a new instance of BillingAddress with the provided country, street, city, and postcode.
func NewBillingAddress(country, street, city, postcode string) *BillingAddress {
	return &BillingAddress{
		Country:  country,
		Street:   street,
		City:     city,
		Postcode: postcode,
	}
}

// NewCreditCard creates a new CreditCard instance with given metadata, name, expiration month and year, card number,
// and billing address.
func NewCreditCard(
	metadata string,
	name string,
	month string,
	year string,
	number string,
	billingAddress BillingAddress,
) *CreditCard {
	return &CreditCard{
		Metadata:       metadata,
		Name:           name,
		Month:          month,
		Year:           year,
		Number:         number,
		BillingAddress: billingAddress,
	}
}

// Clone creates and returns a deep copy of the BillingAddress structure.
func (b *BillingAddress) Clone() BillingAddress {
	return BillingAddress{
		Country:  b.Country,
		Street:   b.Street,
		City:     b.City,
		Postcode: b.Postcode,
	}
}

// Clone creates and returns a deep copy of the CreditCard structure.
func (c *CreditCard) Clone() *CreditCard {
	return &CreditCard{
		Metadata:       c.Metadata,
		Name:           c.Name,
		Month:          c.Month,
		Year:           c.Year,
		Number:         c.Number,
		BillingAddress: c.BillingAddress.Clone(),
	}
}
