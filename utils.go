package sdkcm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//GetMD5Hash excute md5 string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//RandomString create rand string with length given
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// Hash make a password hash
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword verify the hashed password
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//GetPageRequest return page on query param page
func GetPageRequest(c *gin.Context) int {
	r, _ := regexp.Compile("^[0-9]*$")
	if c.Query("page") != "" && r.MatchString(c.Query("page")) {
		page, _ := strconv.Atoi(c.Query("page"))
		if page > 0 {
			return page
		}
	}
	return 1
}

func DelacreChannel(obj1, obj2 string) string {
	if len(obj1) != len(obj2) {
		return ""
	}

	//Convert it to buffer byte about 24 item
	buff1, buff2 := []byte(obj1), []byte(obj2)
	//Cut to 2 partition
	str1 := createPartition(buff1[:len(buff1)/2], buff2[:len(buff2)/2])
	str2 := createPartition(buff1[len(buff1)/2:], buff2[len(buff2)/2:])
	return fmt.Sprintf("%s.%s", str1, str2)
}

func createPartition(buff1, buff2 []byte) string {
	str := ""
	spaces := [][]byte{[]byte{48, 57}, []byte{65, 90}, []byte{97, 122}}
	for i := 0; i < len(buff1); i++ {
		total := buff1[i] + buff2[i]
		//space get
		//0 -> 9 a ->z in ascii is 48->57 97 -> 122
		//if over 122, we will debit into 48->122

		for total > 122 || total < 48 {
			if total > 122 {
				total -= 122
			} else {
				total += 47
			}
		}
		var isBelong = false
		for _, space := range spaces {
			if total >= space[0] && total <= space[1] {
				isBelong = true
				break
			}
		}
		if isBelong {
			str += string(total)
		} else {
			for j := 0; j < len(spaces)-1; j++ {
				if total >= spaces[j][1] && total <= spaces[j+1][0] {
					if 2*total > (spaces[j][1] + spaces[j+1][0]) {
						str += string(spaces[j+1][0])
					} else {
						str += string(spaces[j][1])
					}
					break
				}
			}
		}

	}
	return str
}
