package main
import ( 
	"errors" 
	"fmt"
)

var ErrNoAvatarURL = errors.New("chat: あばたーのURLを取得できません.")
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			fmt.Print("this ok\n")
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
