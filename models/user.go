package models

import (
	"time"
	"log"
	"fmt"
	
	"github.com/credli/hcsg/base"
)

const AdminRoleID = "4F57EDFA-F6ED-4C1A-B29A-18A27BEDE63A"

type User struct {
	UserID         string
	UserName       string
	Password       string
	PasswordFormat string
	PasswordSalt   string
	PartnerID      string
	LoggedInAt     time.Time
	Email          string
	IsLockedOut    bool
	IsAdmin        bool
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	var newPasswd string
	switch u.PasswordFormat {
	case "0":
		newPasswd = base.EncodeMD5(u.Password)
		
	case "1":
		newPasswd = base.EncodeMD5(u.PasswordSalt + u.Password)
	default:
		newPasswd = u.Password
	}
	u.Password = fmt.Sprintf("%x", newPasswd)
}

func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, PasswordSalt: u.PasswordSalt}
	newUser.EncodePasswd()
	return u.Password == newUser.Password
}

func GetUser(userId string) (*User, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
		SELECT a.UserId, a.UserName, b.Password, b.PasswordFormat, b.PasswordSalt, d.ID AS PartnerID, b.LoweredEmail as Email, b.IsLockedOut,
		(SELECT 1 FROM aspnet_UsersInRoles WHERE RoleId = '`+AdminRoleID+`' AND UserId = a.UserId) AS IsAdmin
		FROM aspnet_Users AS a
		INNER JOIN aspnet_Membership AS b ON a.UserId = b.UserId
		INNER JOIN PartnerUsers AS c ON c.UserID = a.UserId
		INNER JOIN Partners AS d ON d.ID = c.PartnerID
		WHERE a.UserId = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)
		err = rows.Scan(&u.UserID, &u.UserName, &u.Password, &u.PasswordFormat, &u.PasswordSalt, &u.PartnerID, &u.Email, &u.IsLockedOut, &u.IsAdmin)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserNotExist{userId, ""}
}

func GetUserByName(uname string) (*User, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
		SELECT a.UserId, a.UserName, b.Password, b.PasswordFormat, b.PasswordSalt, d.ID AS PartnerID, b.LoweredEmail as Email, b.IsLockedOut,
		(SELECT 1 FROM aspnet_UsersInRoles WHERE RoleId = '`+AdminRoleID+`' AND UserId = a.UserId) AS IsAdmin
		FROM aspnet_Users AS a
		INNER JOIN aspnet_Membership AS b ON a.UserId = b.UserId
		INNER JOIN PartnerUsers AS c ON c.UserID = a.UserId
		INNER JOIN Partners AS d ON d.ID = c.PartnerID
		WHERE a.UserName = ?`, uname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)
		err = rows.Scan(&u.UserID, &u.UserName, &u.Password, &u.PasswordFormat, &u.PasswordSalt, &u.PartnerID, &u.Email, &u.IsLockedOut, &u.IsAdmin)
		if err != nil {
			return nil, err
		}
		return u, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserNotExist{"", uname}
}

func UserSignIn(uname, passwd string) (user *User, err error) {
	user, err = GetUserByName(uname)
	if err != nil {
		return nil, err
	}
	
	if user.ValidatePassword(passwd) {
		log.Printf("Failed to login user '%s', password incorrect.\n", uname)
		return user, nil
	}
	
	return nil, ErrUserNotExist{user.UserID, ""} 
}