package lastfm

import (
	"encoding/xml"
	"time"
)

// https://www.last.fm/api/show/user.getInfo

type UserInfoParams struct {
	User string `url:"user"`
}

type UserInfo struct {
	XMLName     xml.Name `xml:"user"`
	Name        string   `xml:"name"`
	RealName    string   `xml:"realname"`
	URL         string   `xml:"url"`
	Country     string   `xml:"country"`
	Age         int      `xml:"age"`
	Gender      string   `xml:"gender"`
	Subscriber  IntBool  `xml:"subscriber"`
	Playcount   int      `xml:"playcount"`
	Playlists   int      `xml:"playlists"`
	Bootstrap   int      `xml:"bootstrap"`
	Avatar      Image    `xml:"image"`
	Registered  DateTime `xml:"registered"`
	Type        string   `xml:"type"`
	ArtistCount int      `xml:"artist_count"`
	AlbumCount  int      `xml:"album_count"`
	TrackCount  int      `xml:"track_count"`
}

// https://www.last.fm/api/show/user.getFriends

type FriendsParams struct {
	User  string `url:"user"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type Friends struct {
	XMLName    xml.Name `xml:"friends"`
	User       string   `xml:"user,attr"`
	Page       int      `xml:"page,attr"`
	PerPage    int      `xml:"perPage,attr"`
	TotalPages int      `xml:"totalPages,attr"`
	Total      int      `xml:"total,attr"`
	Users      []struct {
		Name       string   `xml:"name"`
		RealName   string   `xml:"realname"`
		URL        string   `xml:"url"`
		Country    string   `xml:"country"`
		Subscriber IntBool  `xml:"subscriber"`
		Playcount  int      `xml:"playcount"`
		Playlists  int      `xml:"playlists"`
		Bootstrap  int      `xml:"bootstrap"`
		Avatar     Image    `xml:"image"`
		Registered DateTime `xml:"registered"`
		Type       string   `xml:"type"`
	} `xml:"user"`
}

// https://www.last.fm/api/show/user.getRecentTracks

type RecentTracksParams struct {
	User  string    `url:"user"`
	Limit uint      `url:"limit,omitempty"`
	From  time.Time `url:"from,unix,omitempty"`
	To    time.Time `url:"to,unix,omitempty"`
	Page  uint      `url:"page,omitempty"`
}

type RecentTracks struct {
	XMLName    xml.Name `xml:"recenttracks"`
	User       string   `xml:"user,attr"`
	Page       int      `xml:"page,attr"`
	PerPage    int      `xml:"perPage,attr"`
	TotalPages int      `xml:"totalPages,attr"`
	Total      int      `xml:"total,attr"`
	Tracks     []struct {
		Name       string  `xml:"name"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Artist     struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"artist"`
		Album struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"album"`
		Image       Image    `xml:"image"`
		ScrobbledAt DateTime `xml:"date"`
	} `xml:"track"`
}

// RecentTracksExtended is used when extended=1 in the API call.
type RecentTracksExtended struct {
	XMLName    xml.Name `xml:"recenttracks"`
	User       string   `xml:"user,attr"`
	Page       int      `xml:"page,attr"`
	PerPage    int      `xml:"perPage,attr"`
	TotalPages int      `xml:"totalPages,attr"`
	Total      int      `xml:"total,attr"`
	Tracks     []struct {
		Name       string  `xml:"name"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Loved      IntBool `xml:"loved"`
		Streamable IntBool `xml:"streamable"`
		Artist     struct {
			Name  string `xml:"name"`
			URL   string `xml:"url"`
			MBID  string `xml:"mbid"`
			Image Image  `xml:"image"`
		} `xml:"artist"`
		Album struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"album"`
		Image       Image    `xml:"image"`
		ScrobbledAt DateTime `xml:"date"`
	} `xml:"track"`
}
