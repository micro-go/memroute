package memroute

var (
	singleton = NewRouter()
)

func Connect() (Client, error) {
	return singleton.Connect()
}

func Send(route string, data interface{}) error {
	return singleton.Send(route, data)
}
