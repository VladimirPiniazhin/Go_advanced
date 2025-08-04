package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const usersFilePath = "internals/storage/users.json"

// Сохранение пользователя
func SaveUser(email string, password string, name string) error {
	user, err := FindUserByEmail(email)
	if err == nil && user != nil {
		return errors.New("USER ALREADY EXIST")
	}
	data, err := loadUsersData()
	if err != nil {
		return err
	}
	newUser := User{Email: email, Password: password, Name: name}
	data.Users = append(data.Users, newUser)

	return saveUsersData(data)
}

// Сохранение хэша пользователя
func SaveHash(user *User, hash string) error {
	data, err := loadUsersData()
	if err != nil {
		return err
	}
	for i := range data.Users {
		if data.Users[i].Email == user.Email {
			data.Users[i].Hash = hash
			break
		}
	}
	return saveUsersData(data)
}

// Записываем данные в JSON
func saveUsersData(data *UsersData) error {
	file, _ := json.Marshal(data)
	return os.WriteFile(usersFilePath, file, 0644)
}

func FindUserByEmail(email string) (*User, error) {
	// Загружаем данные пользователей
	data, err := loadUsersData()
	if err != nil {
		return nil, err
	}
	// Проходим по ним и проверяем на совпадения
	for _, user := range data.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

// ИСпользуем для загрузки данных их JSON
func loadUsersData() (*UsersData, error) {
	file, err := os.ReadFile(usersFilePath)
	if err != nil {
		// Файл не существует - возвращаем пустую структуру
		return &UsersData{Users: []User{}}, nil
	}

	// Файл пустой - возвращаем пустую структуру
	if len(file) == 0 {
		return &UsersData{Users: []User{}}, nil
	}

	var data UsersData
	err = json.Unmarshal(file, &data)
	return &data, err
}

// Проверяем пользователей с таким хэшем и возвращаем ссылку на него
func GetUserHash(hash string) (*User, error) {
	data, err := loadUsersData()
	if err != nil {
		return nil, err
	}
	// Проходим по пользователям в базе и возвращаем ссылку на найленного
	// У найденого пользователя удаляем хэш
	// Если пользователь не найден то возвращаем нулевой указатель

	for i := range data.Users {
		if data.Users[i].Hash == hash {
			foundUser := data.Users[i]
			data.Users[i].Hash = ""

			if err := saveUsersData(data); err != nil {
				return nil, fmt.Errorf("Failed to save updated data: %w", err)
			}

			return &foundUser, nil
		}
	}
	return nil, nil
}
