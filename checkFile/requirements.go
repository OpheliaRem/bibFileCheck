package checkFile

var numOfSourcesRequirements = [5]map[string]int{
	{
		"total":        15,
		"foreign":      4,
		"periodic2010": 3,
		"modern":       10,
	},
	{
		"total":        20,
		"foreign":      5,
		"periodic2010": 5,
		"modern":       15,
	},
	{
		"total":        25,
		"foreign":      6,
		"periodic2010": 7,
		"modern":       21,
	},
	{
		"total":        30,
		"foreign":      7,
		"periodic2010": 7,
		"modern":       21,
	},
	{
		"total":        30,
		"foreign":      7,
		"periodic2010": 7,
		"modern":       21,
	},
}

var requiredFieldsOfSources = map[string][]string{
	"book":    {"title", "author", "year", "publisher", "pagetotal"},
	"article": {"title", "author", "journal", "volume", "number", "pages", "year"},
}
