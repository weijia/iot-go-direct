package shared


type SftpInfo struct {
	IP                    string `json:"ip"`
	Port                  int    `json:"port"`
	User                  string `json:"user"`
	Pwd                   string `json:"pwd"`
	Path                  string `json:"path"`
}