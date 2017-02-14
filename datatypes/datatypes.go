package datatypes

type Server struct {
	Domain string
	Host string
	Port string
}

type Database struct {
	Type string
	Name string
	Host string
	Port string
	User string
	Password string
}

type Page struct {
    Url string
    Title string
    Content  string
}

type Event struct {
	Class string
	From string
	Message string
	Data map[string]interface{}
}

/*
type GormPage struct {
	Domain string
	Url string `gorm:"not null;unique"`
    Title string
    Content  string
	
}

type Command struct {
	Name string
	From string
	Reason string
	Values interface{}
	Status string
}

type Hit struct {
	Url string
	Method string
	Ip string
	User_agent string
	Referer string
	Date time.Time
}

type DataChanges struct {
	Msg string
	Type string
	Values map[string]interface{}
}*/
