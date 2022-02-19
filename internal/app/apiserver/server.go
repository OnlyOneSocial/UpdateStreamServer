package apiserver

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gobwas/ws"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/katelinlis/AuthBackend/internal/app/store"
	"github.com/katelinlis/AuthBackend/third_party/gopool"
	"github.com/sirupsen/logrus"
)

//Notification ...
type Notification struct {
	Type      string `json:"type"`
	UserID    int    `json:"userid"`
	UserName  string `json:"username"`
	TimeStamp int64  `json:"timestamp"`
}

//NotiffList ...
type NotiffList struct {
	List []Notification
	Lock sync.Mutex
}

//BufferNotiff ...
type BufferNotiff map[int]*NotiffList

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	redis        *redis.Client
	jwtKeys      JwtKeys
	BufferNotiff BufferNotiff
}

const (
	ctxKeyUser ctxKey = iota
)

type ctxKey int8

var (
	errIncorrectEmailOrPassword = errors.New("incorect email or password")
	jwtsignkey                  string
)

func initBufferNotification() BufferNotiff {
	NotificationBufferUser := make(BufferNotiff)

	return NotificationBufferUser
}

func newServer(store store.Store, config *Config) *server {

	s := &server{
		router:  mux.NewRouter(),
		jwtKeys: getJwtKeys("/"),
		logger:  logrus.New(),
		redis: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		store:        store,
		BufferNotiff: initBufferNotification(),
	}
	s.configureRouter()

	jwtsignkey = config.JwtSignKey

	return s
}

func (s *server) GetDataFromToken(token string) (int, error) {

	if token == "" {
		return 0, errors.New("Token is missing")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return 0, msg
		}
		return []byte(jwtsignkey), nil
	})

	if err != nil {
		//s.error(w, r, http.StatusUnauthorized, errors.New("Error parsing token"))
		return 0, errors.New("Error parsing token")
	}
	if parsedToken != nil && parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			userid := claims["userid"].(float64)
			return int(userid), nil
		}
	}
	return 0, nil

}

//Packet2 ...
type Packet2 struct {
	Typeof string                 `json:"typeof"`
	Data   map[string]interface{} `json:"data"`
}

func (s *server) ListenAndServe(Addr string) {

	ln, err := net.Listen("tcp", "0.0.0.0"+Addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		pool := gopool.NewPool(20, 10, 5)
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println(err)
		}

		err = pool.ScheduleTimeout(time.Millisecond, func() {
			u := ws.Upgrader{}
			_, err := u.Upgrade(conn)
			if err != nil {
				fmt.Println(err)
			}

			ch := NewChannel(conn, s)
			go ch.reader()
			ch.writer()

		})
		if err != nil {
			time.Sleep(time.Millisecond)
		}
	}
}
