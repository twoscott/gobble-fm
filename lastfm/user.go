package lastfm

import (
	"time"
)

// https://www.last.fm/api/show/user.getFriends
type FriendsParams struct {
	User  string `url:"user"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type Friends struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Users      []struct {
		Name         string   `xml:"name"`
		RealName     string   `xml:"realname"`
		URL          string   `xml:"url"`
		Country      string   `xml:"country"`
		Subscriber   IntBool  `xml:"subscriber"`
		Playcount    int      `xml:"playcount"`
		Playlists    int      `xml:"playlists"`
		Bootstrap    int      `xml:"bootstrap"`
		Avatar       Image    `xml:"image"`
		RegisteredAt DateTime `xml:"registered"`
		Type         string   `xml:"type"`
	} `xml:"user"`
}

// https://www.last.fm/api/show/user.getInfo
type UserInfoParams struct {
	User string `url:"user"`
}

type UserInfo struct {
	Name         string   `xml:"name"`
	RealName     string   `xml:"realname"`
	URL          string   `xml:"url"`
	Country      string   `xml:"country"`
	Age          int      `xml:"age"`
	Gender       string   `xml:"gender"`
	Subscriber   IntBool  `xml:"subscriber"`
	Playcount    int      `xml:"playcount"`
	Playlists    int      `xml:"playlists"`
	Bootstrap    int      `xml:"bootstrap"`
	Avatar       Image    `xml:"image"`
	RegisteredAt DateTime `xml:"registered"`
	Type         string   `xml:"type"`
	ArtistCount  int      `xml:"artist_count"`
	AlbumCount   int      `xml:"album_count"`
	TrackCount   int      `xml:"track_count"`
}

// https://www.last.fm/api/show/user.getLovedTracks
type LovedTracksParams struct {
	User  string `url:"user"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type LovedTracks struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title  string `xml:"name"`
		URL    string `xml:"url"`
		MBID   string `xml:"mbid"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Image      Image `xml:"image"`
		Streamable struct {
			Preview   IntBool `xml:",chardata"`
			FullTrack IntBool `xml:"fulltrack,attr"`
		} `xml:"streamable"`
		LovedAt DateTime `xml:"date"`
	} `xml:"track"`
}

// https://www.last.fm/api/show/user.getPersonalTags
type UserTagsParams struct {
	User  string `url:"user"`
	Tag   string `url:"tag"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type UserAlbumTags struct {
	User       string `xml:"user,attr"`
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Albums     []struct {
		Title  string `xml:"name"`
		URL    string `xml:"url"`
		MBID   string `xml:"mbid"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Cover Image `xml:"image"`
	} `xml:"albums>album"`
}

type UserArtistTags struct {
	User       string `xml:"user,attr"`
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artists>artist"`
}

type UserTrackTags struct {
	User       string `xml:"user,attr"`
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title string `xml:"name"`
		// All values returned from the Last.fm API are "FIXME". API issue?
		Duration   string `xml:"duration"`
		URL        string `xml:"url"`
		MBID       string `xml:"mbid"`
		Streamable struct {
			Preview   IntBool `xml:",chardata"`
			FullTrack IntBool `xml:"fulltrack,attr"`
		} `xml:"streamable"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Image Image `xml:"image"`
	} `xml:"tracks>track"`
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
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title      string  `xml:"name"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		NowPlaying bool    `xml:"nowplaying,attr"`
		Streamable IntBool `xml:"streamable"`
		Artist     struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"artist"`
		Album struct {
			Title string `xml:",chardata"`
			MBID  string `xml:"mbid,attr"`
		} `xml:"album"`
		Image       Image    `xml:"image"`
		ScrobbledAt DateTime `xml:"date"`
	} `xml:"track"`
}

// RecentTracksExtended is used when extended=1 in the API call.
type RecentTracksExtended struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title      string  `xml:"name"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		NowPlaying bool    `xml:"nowplaying,attr"`
		Loved      IntBool `xml:"loved"`
		Streamable IntBool `xml:"streamable"`
		Artist     struct {
			Name  string `xml:"name"`
			URL   string `xml:"url"`
			MBID  string `xml:"mbid"`
			Image Image  `xml:"image"`
		} `xml:"artist"`
		Album struct {
			Title string `xml:",chardata"`
			MBID  string `xml:"mbid,attr"`
		} `xml:"album"`
		Image       Image    `xml:"image"`
		ScrobbledAt DateTime `xml:"date"`
	} `xml:"track"`
}

// https://www.last.fm/api/show/user.getTopAlbums
type UserTopAlbumsParams struct {
	User   string `url:"user"`
	Period Period `url:"period,omitempty"`
	Limit  uint   `url:"limit,omitempty"`
	Page   uint   `url:"page,omitempty"`
}

type UserTopAlbums struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Albums     []struct {
		Title     string `xml:"name"`
		Rank      int    `xml:"rank,attr"`
		Playcount int    `xml:"playcount"`
		URL       string `xml:"url"`
		MBID      string `xml:"mbid"`
		Artist    struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Cover Image `xml:"image"`
	} `xml:"album"`
}

// https://www.last.fm/api/show/user.getTopArtists
type UserTopArtistsParams struct {
	User   string `url:"user"`
	Period Period `url:"period,omitempty"`
	Limit  uint   `url:"limit,omitempty"`
	Page   uint   `url:"page,omitempty"`
}

type UserTopArtists struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		Rank       int     `xml:"rank,attr"`
		Playcount  int     `xml:"playcount"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/user.getTopTags
type UserTopTagsParams struct {
	User  string `url:"user"`
	Limit uint   `url:"limit,omitempty"`
}

type UserTopTags struct {
	User string `xml:"user,attr"`
	Tags []struct {
		Name  string `xml:"name"`
		Count int    `xml:"count"`
		URL   string `xml:"url"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/user.getTopTracks
type UserTopTracksParams struct {
	User   string `url:"user"`
	Period Period `url:"period,omitempty"`
	Limit  uint   `url:"limit,omitempty"`
	Page   uint   `url:"page,omitempty"`
}

type UserTopTracks struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title      string   `xml:"name"`
		Rank       int      `xml:"rank,attr"`
		Playcount  int      `xml:"playcount"`
		Duration   Duration `xml:"duration"`
		URL        string   `xml:"url"`
		MBID       string   `xml:"mbid"`
		Streamable struct {
			Preview   IntBool `xml:",chardata"`
			FullTrack IntBool `xml:"fulltrack,attr"`
		} `xml:"streamable"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Image Image `xml:"image"`
	} `xml:"track"`
}

// https://www.last.fm/api/show/user.getWeeklyAlbumChart
type WeeklyAlbumChartParams struct {
	User  string    `url:"user"`
	Limit uint      `url:"limit,omitempty"`
	From  time.Time `url:"from,unix,omitempty"`
	To    time.Time `url:"to,unix,omitempty"`
}

type WeeklyAlbumChart struct {
	User   string   `xml:"user,attr"`
	From   DateTime `xml:"from,attr"`
	To     DateTime `xml:"to,attr"`
	Albums []struct {
		Title     string `xml:"name"`
		Rank      int    `xml:"rank,attr"`
		Playcount int    `xml:"playcount"`
		URL       string `xml:"url"`
		MBID      string `xml:"mbid"`
		Artist    struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"artist"`
	} `xml:"album"`
}

// https://www.last.fm/api/show/user.getWeeklyArtistChart
type WeeklyArtistChartParams struct {
	User  string    `url:"user"`
	Limit uint      `url:"limit,omitempty"`
	From  time.Time `url:"from,unix,omitempty"`
	To    time.Time `url:"to,unix,omitempty"`
}

type WeeklyArtistChart struct {
	User    string   `xml:"user,attr"`
	From    DateTime `xml:"from,attr"`
	To      DateTime `xml:"to,attr"`
	Artists []struct {
		Name      string `xml:"name"`
		Rank      int    `xml:"rank,attr"`
		Playcount int    `xml:"playcount"`
		URL       string `xml:"url"`
		MBID      string `xml:"mbid"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/user.getWeeklyChartList
type WeeklyChartListParams struct {
	User string `url:"user"`
}

type WeeklyChartList struct {
	User   string `xml:"user,attr"`
	Charts []struct {
		From string `xml:"from,attr"`
		To   string `xml:"to,attr"`
	} `xml:"chart"`
}

// https://www.last.fm/api/show/user.getWeeklyTrackChart
type WeeklyTrackChartParams struct {
	User  string    `url:"user"`
	Limit uint      `url:"limit,omitempty"`
	From  time.Time `url:"from,unix,omitempty"`
	To    time.Time `url:"to,unix,omitempty"`
}

type WeeklyTrackChart struct {
	User   string   `xml:"user,attr"`
	From   DateTime `xml:"from,attr"`
	To     DateTime `xml:"to,attr"`
	Tracks []struct {
		Title     string `xml:"name"`
		Rank      int    `xml:"rank,attr"`
		Playcount int    `xml:"playcount"`
		URL       string `xml:"url"`
		MBID      string `xml:"mbid"`
		Artist    struct {
			Name string `xml:",chardata"`
			MBID string `xml:"mbid,attr"`
		} `xml:"artist"`
		Image Image `xml:"image"`
	} `xml:"track"`
}
