package main
import "testing"

func TestAuthAvatar(t *testing.T){
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、 AuthAvatar.GetAvatarURL は" + "ErrNoAvatarURL を返すべきです。" )
			
	}
}
