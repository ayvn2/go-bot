package models

type Config struct {
	SteamcmdPath  string          `json:"steamcmd_path"`
	AdminPassword string          `json:"admin_password"`
	Accounts      map[int]Account `json:"accounts"`
	Games         []Game          `json:"games"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Game struct {
	AppID    string `json:"app_id"`
	Accounts []int  `json:"accounts"`
}
