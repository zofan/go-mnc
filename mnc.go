package mnc

const (
	statusActive   = `Operational`
	statusInactive = `Inactive`
)

type MNC struct {
	Country  string
	Network  string
	Operator string

	MCC int
	MNC int
}

func ByCountry(v string) (list []MNC) {
	for _, x := range List {
		if x.Country == v {
			list = append(list, x)
		}
	}

	return nil
}

func ByMCC(v int) (list []MNC) {
	for _, x := range List {
		if x.MCC == v {
			list = append(list, x)
		}
	}

	return nil
}
