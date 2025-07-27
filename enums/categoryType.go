package enums

type CategoryType string

const (
	Classic   CategoryType = "classic"
	Sparkling CategoryType = "sparkling"
	Smoothies CategoryType = "smoothies"
	Tea       CategoryType = "tea"
	Powders   CategoryType = "powders"
	IceCream  CategoryType = "ice_cream"
	Other     CategoryType = "other"
)

func (c CategoryType) GetDisplayName() string {
	switch c {
	case Classic:
		return "Classic"
	case Sparkling:
		return "Sparkling"
	case Smoothies:
		return "Smoothies"
	case Tea:
		return "Tea"
	case Powders:
		return "Powders"
	case IceCream:
		return "IceCream"
	case Other:
		return "Other"
	default:
		return string(c)
	}
}

func GetAllCategories() []CategoryType {
	return []CategoryType{
		Classic,
		Sparkling,
		Smoothies,
		Tea,
		Powders,
		IceCream,
		Other,
	}
}
