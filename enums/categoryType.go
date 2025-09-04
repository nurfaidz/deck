package enums

type CategoryType string

const (
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
		Appetizers,
		MainCourse,
		Desserts,
		Snacks,
		Food,
		Pastry,
		Other,
	}
}
