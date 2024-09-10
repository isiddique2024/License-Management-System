package models

type LicenseResponse struct {
	Key         string `json:"key"`
	Note        string `json:"note"`
	CreatedOn   string `json:"created_on"`
	Duration    string `json:"duration"`
	GeneratedBy string `json:"generated_by"`
	UsedOn      string `json:"used_on"`
	ExpiresOn   string `json:"expires_on"`
	Status      string `json:"status"`
	IP          string `json:"ip"`
	HWID        string `json:"hwid"`
}

type RedeemLicenseResponse struct {
	Key         string `json:"key"`
	Note        string `json:"note"`
	CreatedOn   string `json:"created_on"`
	Duration    string `json:"duration"`
	GeneratedBy string `json:"generated_by"`
	UsedOn      string `json:"used_on"`
	ExpiresOn   string `json:"expires_on"`
	Status      string `json:"status"`
	IP          string `json:"ip"`
	HWID        string `json:"hwid"`
}
