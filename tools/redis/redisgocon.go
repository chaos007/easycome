package redis

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
    "net/url"
    "strings"
)

func Connect(serverAddr string, auth string, db string) (con redis.Conn, err error) {
    if serverAddr == "" {
        return
    }

    con, err = redis.Dial("tcp", serverAddr)

    if err != nil {
        fmt.Println(err)
        return
    }

    if auth != "" && len(auth) > 0 {
        _, err = con.Do("AUTH", auth)

        if err != nil {
            fmt.Println(err)
            return
        }
    }

    if db != "" && len(db) > 0 {
        _, err = con.Do("SELECT", db)

        if err != nil {
            fmt.Println(err)
            return
        }
    }

    return
}

func ConnectUrl(conStr string) (con redis.Conn, err error) {
    parsed, err := url.Parse(conStr)

    if err != nil {
        fmt.Print(err)
        return
    }

    authPw := ""

    if parsed.User != nil {
        password, ok := parsed.User.Password()

        if ok {
            authPw = password
        }
    }

    con, err = redis.Dial("tcp", parsed.Host)

    if err != nil {
        fmt.Println(err)
        return
    }

    if authPw != "" && len(authPw) > 0 {
        _, err = con.Do("AUTH", authPw)

        if err != nil {
            fmt.Println(err)
            return
        }
    }

    if parsed.Path != "" && len(parsed.Path) > 1 {
        db := strings.TrimPrefix(parsed.Path, "/")
        con.Do("SELECT", db)
    }

    return
}