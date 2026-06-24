package types

type TelegramConfig struct {
	ChatID int64  `json:"chatID"`
	Token  string `json:"token"`
}

type BrowserConfig struct {
	Subscription string `json:"subscription"`
}

type UserBody struct {
	ID          uint64 `json:"user_id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	OldPassword string `json:"oldpassword"`
	Email       string `json:"email"`
}

type TaskBody struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Link        string `json:"link"`
	DueDate     string `json:"duedate"`
	IsBackLog   bool   `json:"isbacklog"`
	IsCompleted bool   `json:"iscompleted"`
	IsStarted   bool   `json:"isstarted"`
	Priority    bool   `json:"priority"`
	Urgency     bool   `json:"urgency"`
	CategoryID  uint   `json:"categoryid,omitempty"`
}

type UserPreferencesBody struct {
	ActiveCategoryID uint `json:"active_category_id"`
}

type CategoryBody struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Active  bool   `json:"active"`
	GroupID uint   `json:"group_id,omitempty"`
}

type DeviceBody struct {
	DeviceID string `json:"deviceid"`
	UserID   uint   `json:"userid"`
}

type DeviceRegisterBody struct {
	DeviceName string `json:"device_name" binding:"required"`
}

type GroupBody struct {
	Name string `json:"name" binding:"required"`
}

type GroupMemberBody struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type CategoryShareBody struct {
	GroupID uint `json:"group_id" binding:"required"`
}
