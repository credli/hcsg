package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	FirstName      string
	LastName       string
	Email          string
	IsLockedOut    bool
	IsAdmin        bool
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	if u.PasswordFormat == "0" {
		newPasswd := base.EncodeMD5(u.PasswordSalt + u.Password)
		u.Password = fmt.Sprintf("%x", newPasswd)
	}
}

func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, PasswordSalt: u.PasswordSalt}
	newUser.EncodePasswd()
	return u.Password == newUser.Password
}

func GetUser(userId string) (*User, error) {
	log.Println("Getting user from database...")

	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`
		SELECT a.UserId, a.UserName, b.Password, b.PasswordFormat, b.PasswordSalt, d.ID AS PartnerID, b.LoweredEmail as Email, b.IsLockedOut,
		(SELECT 1 FROM aspnet_UsersInRoles WHERE RoleId = '`+AdminRoleID+`' AND UserId = a.UserId) AS IsAdmin,
		e.PropertyNames, e.PropertyValuesString
		FROM aspnet_Users AS a
		INNER JOIN aspnet_Membership AS b ON a.UserId = b.UserId
		INNER JOIN PartnerUsers AS c ON c.UserID = a.UserId
		INNER JOIN Partners AS d ON d.ID = c.PartnerID
		LEFT JOIN aspnet_Profile AS e ON e.UserId = a.UserId
		WHERE a.UserId = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)
		var (
			propNames   string
			propVals    string
			isLockedOut *sql.NullBool
			isAdmin     *sql.NullBool
		)
		err = rows.Scan(&u.UserID, &u.UserName, &u.Password, &u.PasswordFormat, &u.PasswordSalt, &u.PartnerID, &u.Email, &isLockedOut, &isAdmin, &propNames, &propVals)
		if err != nil {
			return nil, err
		}
		if isLockedOut != nil && isLockedOut.Valid {
			u.IsLockedOut = isLockedOut.Bool
		}
		if isAdmin != nil && isAdmin.Valid {
			u.IsAdmin = isAdmin.Bool
		}
		u.parseProfileProperties(propNames, propVals)
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
		(SELECT 1 FROM aspnet_UsersInRoles WHERE RoleId = '`+AdminRoleID+`' AND UserId = a.UserId) AS IsAdmin,
		e.PropertyNames, e.PropertyValuesString
		FROM aspnet_Users AS a
		INNER JOIN aspnet_Membership AS b ON a.UserId = b.UserId
		INNER JOIN PartnerUsers AS c ON c.UserID = a.UserId
		INNER JOIN Partners AS d ON d.ID = c.PartnerID
		LEFT JOIN aspnet_Profile AS e ON e.UserId = a.UserId
		WHERE a.UserName = ?`, uname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)
		var (
			propNames   string
			propVals    string
			isLockedOut *sql.NullBool
			isAdmin     *sql.NullBool
		)
		err = rows.Scan(&u.UserID, &u.UserName, &u.Password, &u.PasswordFormat, &u.PasswordSalt, &u.PartnerID, &u.Email, &isLockedOut, &isAdmin, &propNames, &propVals)
		if err != nil {
			return nil, err
		}
		if isLockedOut != nil && isLockedOut.Valid {
			u.IsLockedOut = isLockedOut.Bool
		}
		if isAdmin != nil && isAdmin.Valid {
			u.IsAdmin = isAdmin.Bool
		}
		u.parseProfileProperties(propNames, propVals)
		return u, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrUserNotExist{"", uname}
}

func (u *User) DisplayName() string {
	if len(u.FirstName) == 0 && len(u.LastName) == 0 {
		return u.UserName
	} else if len(u.LastName) == 0 {
		return u.FirstName
	} else if len(u.FirstName) == 0 {
		return u.LastName
	} else {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	}
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

func (u *User) parseProfileProperties(names, vals string) {
	_ = "breakpoint"
	u.FirstName, _ = getProfileProperty("FirstName", names, vals)
	u.LastName, _ = getProfileProperty("LastName", names, vals)
}

func getProfileProperty(key, propertyNames, propertyValues string) (string, error) {
	keyIndex := strings.Index(propertyNames, key)
	if keyIndex == -1 {
		return "", fmt.Errorf("Property '%s' was not found", key)
	}

	//Takes care of finding S
	dataType := propertyNames[keyIndex+len(key)+1:]
	dataType = dataType[:strings.Index(dataType, ":")]
	if dataType != "S" {
		return "", fmt.Errorf("Property '%s' is not of type string", key)
	}

	pos1 := propertyNames[keyIndex+len(key)+1+len(dataType)+1:]
	pos1 = pos1[:strings.Index(pos1, ":")]

	pos2 := propertyNames[keyIndex+len(key)+1+len(dataType)+1+len(pos1)+1:]
	pos2 = pos2[:strings.Index(pos2, ":")]

	startPos, _ := strconv.ParseInt(pos1, 0, 64)
	endPos, _ := strconv.ParseInt(pos2, 0, 64)

	return propertyValues[startPos : startPos+endPos], nil
}
