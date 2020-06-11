package keyman

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Keyman struct {
	Keypre    string
	RedisPool *redis.Pool
}

type HKey struct {
	Key  string `form:"key" json:"key" xml:"key"`
	Name string `form:"name" json:"name" xml:"name"`
}

type Key struct {
	Key    string `form:"key" json:"key" xml:"key" binding:"required"`
	Expday int    `form:"expday" json:"expday" xml:"expday"`
	Number int64  `form:"number" json:"number" xml:"number"`
}

func (keyman *Keyman) keyAddPre(key string) string {
	return keyman.Keypre + key
}

func (keyman *Keyman) keyDelPre(key string) string {
	return strings.Replace(key, keyman.Keypre, "", 1)
}

func (keyman *Keyman) StrToPriv(key string) *ecdsa.PrivateKey {
	key = keyman.keyDelPre(key)
	priv := new(ecdsa.PrivateKey)
	d := big.NewInt(0)
	d.SetString(key, 0)
	priv.D = d
	priv.PublicKey.Curve = crypto.S256()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(priv.D.Bytes())
	return priv
}

func (keyman *Keyman) GetPriv(c *gin.Context) (*ecdsa.PrivateKey, error) {
	key := c.GetHeader("key")
	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()
	isExist, err := redis.Int(redisConn.Do("HEXISTS", "keys", keyman.keyAddPre(key)))
	if err == redis.ErrNil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if isExist == 0 {
		return nil, nil
	}

	priv := keyman.StrToPriv(key)
	return priv, nil
}

func (keyman *Keyman) GetManPriv(c *gin.Context) (*ecdsa.PrivateKey, error) {
	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()
	key := c.GetHeader("key")
	isExist, err := redis.Int(redisConn.Do("HEXISTS", "mkeys", keyman.keyAddPre(key)))
	if err == redis.ErrNil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if isExist == 0 {
		return nil, nil
	}
	priv := keyman.StrToPriv(key)
	return priv, nil
}

func (keyman *Keyman) Enable(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	var key Key

	err = c.BindJSON(&key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()
	isExist, err := redis.Int(redisConn.Do("HEXISTS", "keys", keyman.keyAddPre(key.Key)))
	if err == redis.ErrNil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "key not exist",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if isExist == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "key not exist",
		})
		return
	}

	_, err = redisConn.Do("SET", keyman.keyAddPre(key.Key), key.Number)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	exptime := time.Now()
	exptime = exptime.Add(time.Duration(key.Expday) * time.Hour * 24)
	sec := exptime.Unix()
	_, err = redisConn.Do("EXPIREAT", keyman.keyAddPre(key.Key), sec)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"expdate": exptime.Format("2006-01-02T15:04:05"),
		"number":  key.Number,
	})
}

func (keyman *Keyman) Addkey(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	var key HKey
	err = c.BindJSON(&key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()

	if len(key.Key) < 70 {
		k, _ := crypto.GenerateKey()
		key.Key = k.D.String()
	}

	_, err = redisConn.Do("HSET", "keys", keyman.keyAddPre(key.Key), key.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"key":    key.Key,
		"name":   key.Name,
	})

}

func (keyman *Keyman) Delkey(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	var key HKey
	err = c.BindJSON(&key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if len(key.Key) < 70 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "too less",
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()

	_, err = redisConn.Do("HDEL", "keys", keyman.keyAddPre(key.Key))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"key":    key.Key,
	})
}

func (keyman *Keyman) Getkey(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	var key HKey
	err = c.BindJSON(&key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()

	name, err := redis.String(redisConn.Do("HGET", "keys", keyman.keyAddPre(key.Key)))
	if err == redis.ErrNil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "key not exist",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	sec, err := redis.Int(redisConn.Do("TTL", keyman.keyAddPre(key.Key)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if sec < 0 {
		sec = 0
	}

	number, err := redis.Int(redisConn.Do("GET", keyman.keyAddPre(key.Key)))
	if err == redis.ErrNil {
		number = 0
	} else if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"name":   name,
		"sec":    sec,
		"number": number,
	})

}

func (keyman *Keyman) Listkey(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()

	keys, err := redis.Strings(redisConn.Do("HKEYS", "keys"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	for i := 0; i < len(keys); i++ {
		keys[i] = keyman.keyDelPre(keys[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"keys":   keys,
	})

}

func (keyman *Keyman) Diskey(c *gin.Context) {
	priv, err := keyman.GetManPriv(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if priv == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "access denied",
		})
		return
	}

	var key HKey
	err = c.BindJSON(&key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if len(key.Key) < 70 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "too less",
		})
		return
	}

	redisConn := keyman.RedisPool.Get()
	defer redisConn.Close()

	_, err = redisConn.Do("DEL", keyman.keyAddPre(key.Key))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"key":    key.Key,
	})
}
