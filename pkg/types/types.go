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
	DueDate     string `json:"duedate"`
	IsBackLog   bool   `json:"isbacklog,omitempty"`
	IsCompleted bool   `json:"iscompleted,omitempty"`
	Priority    bool   `json:"priority,omitempty"`
	Urgency     bool   `json:"urgency,omitempty"`
	CategoryID  string `json:"categoryid,omitempty"`
}
