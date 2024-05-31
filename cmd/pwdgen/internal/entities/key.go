package entities

type XyKey struct {
	Version int   `json:"version"`
	Key     []Key `json:"key"`
}

type Key struct {
	Name      string  `json:"name"`
	Account   string  `json:"account"`
	Password  string  `json:"password"`
	Password2 string  `json:"password2"`
	Url       string  `json:"url"`
	Note      string  `json:"note"`
	Extra     []Extra `json:"extra"`
}

type Extra struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ChromeCSV chrome, edge csv password
type ChromeCSV struct {
	Name     string `json:"name" csv:"name"`
	URL      string `json:"url" csv:"url"`
	Username string `json:"username" csv:"username"`
	Password string `json:"password" csv:"password"`
	Note     string `json:"note" csv:"note"`
}
