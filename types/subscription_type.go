package types

type SubscriptionType string

const (
	NoSubscription SubscriptionType = "nosubs"
	Monthly        SubscriptionType = "month"
	Yearly         SubscriptionType = "year"
)

func (s SubscriptionType) IsValid() bool {
	return s == Monthly || s == Yearly || s == NoSubscription
}
