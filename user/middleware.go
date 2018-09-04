package user

import (
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"strings"
	"time"

	"sync"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/btcsuite/btcutil/base58"
	gin "gopkg.in/gin-gonic/gin.v1"
)

var once sync.Once

var session *Session

type Session struct {
	kvs map[string]interface{}
	kes map[string]time.Time
	sync.RWMutex
}

func (s *Session) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.kvs[key]
}

func (s *Session) Set(key string, val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.kvs[key] = val
	s.kes[key] = time.Now().Add(time.Minute * 5)
}

func (s *Session) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.kes, key)
	delete(s.kvs, key)
}

func (s *Session) Expired() {
	s.Lock()
	defer s.Unlock()
	for k, e := range s.kes {
		if e.Before(time.Now()) {
			delete(s.kes, k)
			delete(s.kvs, k)
		}
	}
}

func NewSession() *Session {
	once.Do(func() {
		session = &Session{
			kvs: make(map[string]interface{}),
			kes: make(map[string]time.Time),
		}
		go func() {
			ticker := time.NewTicker(time.Hour)
			for {
				select {
				case <-ticker.C:
					session.Expired()
				}
			}
		}()
	})
	return session
}

func getsessions(c *gin.Context) *Session {
	return NewSession()
}

// Token 用户
type Token struct {
	UserName   string `json:"username"`
	SignInTime int64  `json:"signintime"`
	Expire     int64  `json:"expire"`
}

func (t *Token) token() string {
	h := sha256.New()
	h.Write([]byte(t.UserName))
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t.SignInTime))
	h.Write(buf)
	return base58.Encode(h.Sum(nil))
}

func Checklogged() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := getsessions(c)
		if strings.EqualFold(c.Request.RequestURI, "/signin") {
			c.Next()
			return
		}

		if strings.EqualFold(c.Request.RequestURI, "/signup") {
			c.Next()
			return
		}

		if strings.EqualFold(c.Request.RequestURI, "/sendsms") {
			c.Next()
			return
		}

		//auth
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.JSON(http.StatusUnauthorized, &common.APIRespone{
				ErrCode: common.UnauthorizedCode,
				ErrMsg:  "未登录",
			})
			c.Abort()
			return
		}

		t := session.Get(strings.TrimPrefix(h, "Bearer "))
		if t == nil {
			c.JSON(http.StatusUnauthorized, &common.APIRespone{
				ErrCode: common.UnauthorizedCode,
				ErrMsg:  "无效登陆",
			})
			c.Abort()
			return
		}
		token := t.(*Token)
		defer func() {
			session.Set(token.token(), token)
		}()

		if token.Expire == 0 {
			session.Delete(token.token())
			c.JSON(http.StatusUnauthorized, &common.APIRespone{
				ErrCode: common.UnauthorizedCode,
				ErrMsg:  "未登陆成功",
			})
			c.Abort()
			return
		} else if time.Now().Sub(time.Unix(token.Expire, 0)) > 0 {
			session.Delete(token.token())
			c.JSON(http.StatusUnauthorized, &common.APIRespone{
				ErrCode: common.UnauthorizedCode,
				ErrMsg:  "会话超时，请重新登录",
			})
			c.Abort()
			return
		}
		// 更新过期时间
		token.Expire = time.Now().Add(600 * time.Second).Unix()
		c.Next()
	}
}
