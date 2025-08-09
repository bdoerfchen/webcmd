package interfaces

type RouteParsingRemark struct {
	Message    string
	IsCritical bool
}

type RouteErrorCollection []RouteParsingRemark

func (c *RouteErrorCollection) Error() string {
	return "one or multiple routes were not loaded successfully"
}
