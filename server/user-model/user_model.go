/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package userModel

import (
    "math/rand"
    "log"
    "strings"
    "errors"

    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"

    "github.com/GehirnInc/crypt"
    _ "github.com/GehirnInc/crypt/sha256_crypt"

)

const schema = `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        id          INTEGER PRIMARY KEY,
        username    VARCHAR(255) NOT NULL UNIQUE,
        password    VARCHAR(255) NOT NULL,
        isadmin     BOOLEAN DEFAULT FALSE
    );`

type Model struct {
    db *sqlx.DB
}

type User struct {
    Id          int     `db:"id"        json:"id"`
    Username    string  `db:"username"  json:"username"`
    Password    string  `db:"password"  json:"password"`
    IsAdmin     bool    `db:"isadmin"   json:"isadmin"`
    Limit       int     `db:"-"         json:"limit,omitempty"`
    Offset      int     `db:"-"         json:"offset,omitempty"`
}

type Page struct {
    Total       int     `json:"total"`
    Offset      int     `json:"offset"`
    Limit       int     `json:"limit"`
    UserPattern string  `json:"user_pattern"`
    Users       *[]User `json:"users,omitempty"`
}

func randString(n int) string {
    const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    arr := make([]byte, n)
    for i := range arr {
        arr[i] = letters[rand.Intn(len(letters))]
    }
    return string(arr)
}

func createHash(key string) (string, error) {
    crypt := crypt.SHA256.New()
    return crypt.Generate([]byte(key), []byte("$5$" + randString(12)))
}

func CheckHash(hash, password string) error {
    arr := strings.Split(hash, "$")
    if len(arr) < 3 {
        return errors.New("incorrect hash structure")
    }
    hashType := arr[1]
    hashSalt := arr[2]

    if hashType != "5" {
        return errors.New("incorrect hash type")
    }
    crypt := crypt.SHA256.New()
    newHash, _ := crypt.Generate([]byte(password), []byte("$" + hashType + "$" + hashSalt))

    if hash != newHash {
        return errors.New("password is incorrect")
    }
    return nil
}


func (this *Model) Migrate() error {
    _, err := this.db.Exec(schema)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) List(page *Page) (error) {
    var request string
    var err error
    var total int

    userPattern := "%" + page.UserPattern + "%"
    request = `SELECT COUNT(id) as total FROM users WHERE username LIKE $1`
    err = this.db.QueryRow(request, userPattern).Scan(&total)
    if err != nil {
        log.Println(err)
        return err
    }
    page.Total = total

    var users []User
    request = `SELECT id, username, '' as password, isadmin FROM users
                WHERE username LIKE $1
                ORDER BY username LIMIT $2 OFFSET $3`

    err = this.db.Select(&users, request, userPattern, page.Limit, page.Offset)
    if err != nil {
        log.Println(err)
        return err
    }
    page.Users = &users
    return nil
}

func (this *Model) Create(user User) error {
    user.Password, _ = createHash(user.Password)
    request := `INSERT INTO users(username, password, isadmin) VALUES ($1, $2, $3)`
    _, err := this.db.Exec(request, user.Username, user.Password, user.IsAdmin)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Delete(user User) error {
    request := `DELETE FROM users WHERE id = $1`
    _, err := this.db.Exec(request, user.Id)
    if err != nil {
        log.Println(err)
        return err

    }
    return nil
}

func (this *Model) Find(user User) (User, error) {
    request := `SELECT id, username, '' as password, isadmin FROM users WHERE username = $1 LIMIT 1`
    var out User
    err := this.db.Get(&out, request, user.Username)
    if err != nil {
        log.Println(err)
        return out, err
    }
    return out, nil
}

func (this *Model) Update(user User) error {
    var err error
    if len(user.Password) > 0 {
        user.Password, _ = createHash(user.Password)
        request := `UPDATE users SET username = $1, password = $2, isadmin = $3 WHERE id = $4`
        _, err = this.db.Exec(request, user.Username, user.Password, user.IsAdmin, user.Id)
    } else {
        request := `UPDATE users SET username = $1, isadmin = $2 WHERE id = $3`
        _, err = this.db.Exec(request, user.Username, user.IsAdmin, user.Id)
    }
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Check(user *User) error {
    username := user.Username
    password := user.Password

    request := `SELECT * FROM users WHERE username = $1 LIMIT 1`
    err := this.db.Get(user, request, username)
    if err != nil {
        log.Println(err)
        return err
    }

    err = CheckHash(user.Password, password)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func New(db *sqlx.DB) *Model {
    model := Model{
        db: db,
    }
    return &model
}
