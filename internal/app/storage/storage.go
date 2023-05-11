package storage

var mainMap = map[string]string{"hello": "world"}

func SaveToMap(newUrl string, oldUrl string) string {
	count, ok := mainMap[oldUrl]
	if ok {
		return count
	}
	mainMap[oldUrl] = newUrl
	return newUrl
}

func GetFromMap(newUrl string) string {
	var ans string
	for key, value := range mainMap {
		if value == newUrl {
			ans = key
		}
	}
	return ans
}
