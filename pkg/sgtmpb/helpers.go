package sgtmpb

import (
	"encoding/json"
	"fmt"
	"time"

	"ultre.me/calcbiz/pkg/soundcloud"
)

// Post

func (p *Post) ApplyDefaults() {
	if p.Title == "" {
		p.Title = "noname"
	}
	if p.ArtworkURL == "" && p.Provider == Provider_SoundCloud {
		var metadata soundcloud.Track
		err := json.Unmarshal([]byte(p.ProviderMetadata), &metadata)
		if err == nil {
			p.ArtworkURL = metadata.User.AvatarURL
		}
	}
}

func (p *Post) CanonicalURL() string {
	return fmt.Sprintf("/post/%d", p.ID)
}

func (p *Post) GoDuration() time.Duration {
	return time.Millisecond * time.Duration(p.Duration)
}

// User

func (u *User) ApplyDefaults() {

}

func (u *User) CanonicalURL() string {
	return fmt.Sprintf("/@%s", u.Slug)
}

func (u *User) DisplayName() string {
	// FIXME: firstname, lastname
	return fmt.Sprintf("@%s", u.Slug)
}