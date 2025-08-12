package types

type SubscriptionType string

const (
	SubscriptionMonthly SubscriptionType = "month"
	SubscriptionYearly  SubscriptionType = "year"
)

func (s SubscriptionType) IsValid() bool {
	return s == SubscriptionMonthly || s == SubscriptionYearly
}
