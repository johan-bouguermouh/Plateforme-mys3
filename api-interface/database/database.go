package database

import (
	"api-interface/models"
	"fmt"
)

// import (
// 	"api-interface/models"
// 	"fmt"
// 	"sync"
// )

// var (
// 	db []*models.User
// 	mu sync.Mutex
// )

// // Connect with database
// func Connect() {
// 	db = make([]*models.User, 0)
// 	fmt.Println("Connected with Database")
// }

//	func Insert(user *models.User) {
//		mu.Lock()
//		db = append(db, user)
//		mu.Unlock()
//	}
func Insert(user *models.User) {
	mu.Lock()
	defer mu.Unlock()
	db = append(db, user)
	fmt.Printf("Utilisateur inséré : %+v\n", user)
}

// func Get() []*models.User {
// 	return db
// }

func GetByName(name string) *models.User {
	for _, user := range db {
		if user.Name == name {
			return user
		}
	}
	return nil
}
