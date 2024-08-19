package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Id uuid.UUID

func (id Id) MarshalJSON() ([]byte, error) {
	uuidValue := uuid.UUID(id)
	return json.Marshal(uuidValue.String())
}

func (id *Id) UnmarshalJSON(data []byte) error {
	var idStr string
	if err := json.Unmarshal(data, &idStr); err != nil {
		return err
	}
	uuidValue, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	*id = Id(uuidValue)
	return nil
}

type User struct {
	Id        Id     `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Biography string `json:"bio"`
}

type UsersRepository struct {
	data map[Id]User
}

func NewUsersRepository() UsersRepository {
	return UsersRepository{data: map[Id]User{}}
}

func (db UsersRepository) FindAll() []User {
	users := make([]User, 0, len(db.data))

	for _, user := range db.data {
		users = append(users, user)
	}

	return users
}

func (db UsersRepository) FindById(id Id) User {
	user := db.data[id]

	return user
}

func (db UsersRepository) Insert(user User) User {
	user.Id = Id(uuid.New())

	db.data[user.Id] = user

	return user
}

func (db UsersRepository) Update(user User) User {
	row, ok := db.data[user.Id]
	if !ok {
		return User{}
	}

	row.FirstName = user.FirstName
	row.LastName = user.LastName
	row.Biography = user.Biography

	db.data[user.Id] = row

	return row
}

func (db UsersRepository) Delete(id Id) {
	delete(db.data, id)
}
