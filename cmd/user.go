package main

type UserAgreement struct {
	InstitionId        string   `json:"instition_id"`
	MaxHistoricalDays  int      `json:"max_historical_days"`
	AccessValidForDays int      `json:"access_valid_for_days"`
	AccessScope        []string `json:"access_scope"`
}

func GetEndUserAgreement(institutionid string) (UserAgreement, error) {
	url := "https://bankaccountdata.gocardless.com/api/v2/agreements/enduser/"

	payload := UserAgreement{
		InstitionId:        institutionid,
		MaxHistoricalDays:  180,
		AccessValidForDays: 180,
		AccessScope:        []string{"balances", "details", "transactions"},
	}
}
