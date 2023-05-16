package storage

//type ObjStorage struct {
//	shortUrl string
//}

type Storage map[string]string

func NewStorage() Storage {
	arr := make(Storage)
	return arr
}

func (strg Storage) Save(newUrl string, oldUrl string) string {
	if strg == nil {
		strg[oldUrl] = newUrl
		return newUrl
	}
	val, ok := strg[oldUrl]
	if ok {
		return val
	}
	strg[oldUrl] = newUrl
	return newUrl
}

func (strg Storage) OldUrl(newUrl string) string {
	var ans string
	for key, value := range strg {
		if value == newUrl {
			ans = key
		}
	}
	return ans
}
