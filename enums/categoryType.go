package enums

type CategoryType string

const (
	Classic    CategoryType = "classic"
	Sparkling  CategoryType = "sparkling"
	Smoothies  CategoryType = "smoothies"
	Tea        CategoryType = "tea"
	Powders    CategoryType = "powders"
	IceCream   CategoryType = "ice_cream"
	Appetizers CategoryType = "appetizers"
	MainCourse CategoryType = "main_course"
	Desserts   CategoryType = "desserts"
	Snacks     CategoryType = "snacks"
	Food       CategoryType = "food"
	Pastry     CategoryType = "pastry"
	Other      CategoryType = "other"
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
	case Appetizers:
		return "Appetizers"
	case MainCourse:
		return "Main Course"
	case Desserts:
		return "Desserts"
	case Snacks:
		return "Snacks"
	case Food:
		return "Food"
	case Pastry:
		return "Pastry"
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
		Appetizers,
		MainCourse,
		Desserts,
		Snacks,
		Food,
		Pastry,
		Other,
	}
}
