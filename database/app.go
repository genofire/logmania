package database

type Application struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
	//Entries []*Entry `json:"entries" gorm:"-"`
}

func IsTokenValid(token string) bool {
	result := db.Where("token = ?", token).First(&Application{})
	return !result.RecordNotFound()
}
