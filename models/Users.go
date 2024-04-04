package models

import (
	"Onepiece/middleware"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Users struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u Users) Register(db *sql.DB) []byte {
	Response := make(map[string]interface{})

	sha := sha256.New()
	sha.Write([]byte(u.Password))
	pass := sha.Sum(nil)
	u.Password = fmt.Sprintf("%x", pass)

	// registered := findUser(db, u.Email)

	// if !registered {
	id := uuid.New()
	u.Id = fmt.Sprintf("%v", id)

	_, err := db.Exec("INSERT INTO Users(id, firstname, lastname, email, password) VALUES(?,?,?,?,?)", u.Id, u.Firstname, u.Lastname, u.Email, u.Password)

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "حدث خطأ يرجى إعادة المحاولة"
		responseBytes, err := json.Marshal(Response)

		if err != nil {
			fmt.Println(err.Error())
		}
		return responseBytes
	}

	Response["Done"] = true

	responseBytes, err := json.Marshal(Response)

	if err != nil {
		fmt.Println(err.Error())
	}

	return responseBytes
	// }
	// Response["Error"] = true
	// Response["Message"] = "البريد الإلكتروني موجود بالفعل"
	// responseBytes, err := json.Marshal(Response)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// return responseBytes
}

func (u Users) Login(db *sql.DB) []byte {
	Response := make(map[string]interface{})

	sha := sha256.New()
	sha.Write([]byte(u.Password))
	pass := sha.Sum(nil)
	u.Password = fmt.Sprintf("%x", pass)

	userQuery := db.QueryRow("SELECT email FROM Users WHERE (email, password) = (?,?)", u.Email, u.Password)

	var User Users

	if err := userQuery.Scan(&User.Email); err != nil {
		Response["Error"] = true
		Response["Message"] = "حدث خطأ يرجى إعادة المحاولة"
		responseBytes, err := json.Marshal(Response)

		if err != nil {
			fmt.Println(err.Error())
		}

		return responseBytes
	}

	token, err := middleware.GenerateJWT(User.Email)
	if err != nil {
		Response["Error"] = true
		Response["Message"] = "حدث خطأ يرجى إعادة المحاولة"
		responseBytes, err := json.Marshal(Response)

		if err != nil {
			fmt.Println(err.Error())
		}

		return responseBytes
	}
	Response["Token"] = token
	responseBytes, err := json.Marshal(Response)

	if err != nil {
		fmt.Println(err.Error())
	}

	return responseBytes

}
