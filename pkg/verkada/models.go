package verkada

type User struct {
	Email      string `json:"email"`
	EmployeeID string `json:"employee_id"`
	FullName   string `json:"full_name"`
	UserID     string `json:"user_id"`
}

type UserAccess struct {
	AccessGroups []Group     `json:"access_groups"`
	BleUnlock    bool        `json:"ble_unlock"`
	EndDate      string      `json:"end_date"`
	EntryCode    string      `json:"entry_code"`
	ExternalID   interface{} `json:"external_id"`
	RemoteUnlock bool        `json:"remote_unlock"`
	StartDate    interface{} `json:"start_date"`
	UserID       string      `json:"user_id"`
}

type Group struct {
	GroupID string   `json:"group_id"`
	Name    string   `json:"name"`
	UserIDs []string `json:"user_ids,omitempty"`
}
