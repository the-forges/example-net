package model

import (
	"fmt"
	"sync"
)

// User Error codes
const (
	ErrUserUnexplainable = iota + 400
	ErrUserEmpty
	ErrUserChangeFailure
	ErrUserDeletionFailure
	ErrUserNotFound
)

// UsersMap is used to control recieving and kicking users
type UsersMap interface {
	RecieveUser(*User) UsersMap
	KickUser(int64) (UsersMap, error)
	FindUser(int64) (*User, error)
	FindAllUsers() ([]*User, error)
	Len() int
	Iter() []*User
}

type usersMap struct {
	users map[int64]*User
	lock  sync.Mutex
}

// NewUsersMap builds a usersMap to hold a list of unique users that are
// currently online. If you pass one or more users as an argument, they will get
// added to the list automatically.
func NewUsersMap(users ...*User) UsersMap {
	var um UsersMap = &usersMap{users: make(map[int64]*User)}
	for _, user := range users {
		um = um.RecieveUser(user)
	}
	return um
}

// RecieveUser is how we add new users to the usersMap. It will always return a
// pointer to the usersMap in case you need it for future operations.
func (u *usersMap) RecieveUser(user *User) UsersMap {
	u.lock.Lock()
	if u.users == nil {
		u.users = make(map[int64]*User, 0)
	}
	u.users[int64(user.ID)] = user
	u.lock.Unlock()
	return u
}

// KickUser handles the removal of a user from the usersMap.
func (u *usersMap) KickUser(id int64) (UsersMap, error) {
	_, err := u.FindUser(id)
	if err != nil {
		return u, err
	}
	u.lock.Lock()
	delete(u.users, id)
	u.lock.Unlock()
	return u, nil
}

// FindUser takes an int64 and searches the userMap for it. If not found a 404
// will be raised.
func (u *usersMap) FindUser(id int64) (*User, error) {
	u.lock.Lock()
	user, ok := u.users[id]
	u.lock.Unlock()
	if !ok {
		return nil, fmt.Errorf(
			"user with ID %d is no longer logged on: error code %d",
			id,
			ErrUserNotFound,
		)
	}
	return user, nil
}

// FindAllUsers return a slice of User pointers if there are any, otherwise an
// empty slice.
func (u *usersMap) FindAllUsers() ([]*User, error) {
	u.lock.Lock()
	tempUsers := u.users
	u.lock.Unlock()
	users := make([]*User, len(tempUsers))
	for _, user := range tempUsers {
		users = append(users, user)
	}
	return users, nil
}

// Len returns the amount of users connected at the moment
func (u *usersMap) Len() int {
	return len(u.users)
}

// Iter is for use in range.
// Example:
//
// 	import "fmt"
// 	func main() {
// 		um := NewUsersMap()
// 		for _, user := range um.Iter() {
// 			fmt.Printf("%#v\n", user)
// 		}
// 	}
func (u *usersMap) Iter() []*User {
	u.lock.Lock()
	tempUsers := u.users
	u.lock.Unlock()
	users := make([]*User, len(tempUsers))
	for _, user := range tempUsers {
		users = append(users, user)
	}
	return users
}
