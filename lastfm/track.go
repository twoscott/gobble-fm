package lastfm

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// https://www.last.fm/api/show/track.addTags
type TrackAddTagsParams struct {
	Artist string   `url:"artist"`
	Track  string   `url:"track"`
	Tags   []string `url:"tags,comma"`
}

// https://www.last.fm/api/show/track.getCorrection
type TrackCorrectionParams struct {
	Artist string `url:"artist"`
	Track  string `url:"track"`
}

type TrackCorrection struct {
	Corrections []struct {
		Index           int     `xml:"index,attr"`
		ArtistCorrected IntBool `xml:"artistcorrected,attr"`
		TrackCorrected  IntBool `xml:"trackcorrected,attr"`
		Track           struct {
			Title  string `xml:"name"`
			URL    string `xml:"url"`
			MBID   string `xml:"mbid"`
			Artist struct {
				Name string `xml:"name"`
				URL  string `xml:"url"`
				MBID string `xml:"mbid"`
			} `xml:"artist"`
		} `xml:"track"`
	} `xml:"correction"`
}

// https://www.last.fm/api/show/track.getInfo
type TrackInfoParams struct {
	Artist      string `url:"artist"`
	Track       string `url:"track"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	User        string `url:"username,omitempty"`
}

// https://www.last.fm/api/show/track.getInfo
type TrackInfoMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	User        string `url:"username,omitempty"`
}

// https://www.last.fm/api/show/track.getInfo#attributes
type TrackInfo struct {
	Title      string        `xml:"name"`
	URL        string        `xml:"url"`
	MBID       string        `xml:"mbid"`
	Duration   DurationMilli `xml:"duration"`
	Streamable struct {
		Preview   IntBool `xml:",chardata"`
		FullTrack IntBool `xml:"fulltrack,attr"`
	} `xml:"streamable"`
	Listeners int `xml:"listeners"`
	Playcount int `xml:"playcount"`
	Artist    struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
		MBID string `xml:"mbid"`
	} `xml:"artist"`
	Album struct {
		Artist   string `xml:"artist"`
		Title    string `xml:"title"`
		URL      string `xml:"url"`
		MBID     string `xml:"mbid"`
		Position int    `xml:"position,attr"`
		Image    Image  `xml:"image"`
	} `xml:"album"`
	UserPlaycount *int     `xml:"userplaycount"`
	UserLoved     *IntBool `xml:"userloved"`
	TopTags       []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"toptags>tag"`
	Wiki struct {
		Summary   string   `xml:"summary"`
		Content   string   `xml:"content"`
		Published DateTime `xml:"published"`
	} `xml:"wiki"`
}

// https://www.last.fm/api/show/track.getSimilar
type TrackSimilarParams struct {
	Artist      string `url:"artist"`
	Track       string `url:"track"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
}

// https://www.last.fm/api/show/track.getSimilar
type TrackSimilarMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
}

type SimilarTracks struct {
	Artist string `xml:"artist,attr"`
	Tracks []struct {
		Title      string   `xml:"name"`
		Playcount  int      `xml:"playcount"`
		Match      float64  `xml:"match"`
		URL        string   `xml:"url"`
		MBID       string   `xml:"mbid"`
		Duration   Duration `xml:"duration"`
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

// https://www.last.fm/api/show/track.getTags
type TrackTagsParams struct {
	Artist      string `url:"artist"`
	Track       string `url:"track"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/track.getTags
type TrackTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/track.getTags
type TrackSelfTagsParams struct {
	Artist      string `url:"artist"`
	Track       string `url:"track"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/track.getTags
type TrackSelfTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

type TrackTags struct {
	Artist string `xml:"artist,attr"`
	Track  string `xml:"track,attr"`
	Tags   []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/track.getTopTags
type TrackTopTagsParams struct {
	Artist      string `url:"artist"`
	Track       string `url:"track"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/track.getTopTags
type TrackTopTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

type TrackTopTags struct {
	Artist string `xml:"artist,attr"`
	Track  string `xml:"track,attr"`
	Tags   []struct {
		Name  string `xml:"name"`
		URL   string `xml:"url"`
		Count int    `xml:"count"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/track.love
type TrackLoveParams struct {
	Artist string `url:"artist"`
	Track  string `url:"track"`
}

// https://www.last.fm/api/show/track.removeTag
type TrackRemoveTagParams struct {
	Artist string `url:"artist"`
	Track  string `url:"track"`
	Tag    string `url:"tag"`
}

// https://www.last.fm/api/show/track.scrobble
type ScrobbleParams struct {
	Artist string    `url:"artist"`
	Track  string    `url:"track"`
	Time   time.Time `url:"timestamp,unix"`

	Album       string   `url:"album,omitempty"`
	AlbumArtist string   `url:"albumArtist,omitempty"`
	TrackNumber int      `url:"trackNumber,omitempty"`
	Duration    Duration `url:"duration,omitempty"`
	MBID        string   `url:"mbid,omitempty"`

	Chosen   *bool  `url:"chosenByUser,int,omitempty"`
	Context  string `url:"context,omitempty"`
	StreamID string `url:"streamId,omitempty"`
}

// EncodeIndexValues sets the indexed "key[index]" values in v from the fields
// of p. It returns an error if p cannot be encoded.
func (p ScrobbleParams) EncodeIndexValues(index int, v *url.Values) error {
	values, err := query.Values(p)
	if err != nil {
		return err
	}

	// Set indexed `key[index]` values
	for key, vals := range values {
		indexedKey := fmt.Sprintf("%s[%d]", key, index)
		for _, val := range vals {
			v.Set(indexedKey, val)
		}
	}

	return nil
}

// https://www.last.fm/api/show/track.scrobble
type ScrobbleMultiParams []ScrobbleParams

// EncodeValues implements the url.ValuesEncoder interface.
func (p ScrobbleMultiParams) EncodeValues(key string, v *url.Values) error {
	for i, params := range p {
		if err := params.EncodeIndexValues(i, v); err != nil {
			return err
		}
	}

	return nil
}

// https://www.last.fm/api/show/track.scrobble#attributes
type ScrobbleIgnoredCode int

const (
	ScrobbleNotIgnored ScrobbleIgnoredCode = iota // 0

	ArtistIgnored               // 1
	TrackIgnored                // 2
	TimestampTooOld             // 3
	TimestampTooNew             // 4
	DailyScrobbledLimitExceeded // 5
)

// Message returns the message for the ignored scrobble code.
func (c ScrobbleIgnoredCode) Message() string {
	switch c {
	case ScrobbleNotIgnored:
		return "Not ignored"
	case ArtistIgnored:
		return "Artist was ignored"
	case TrackIgnored:
		return "Track was ignored"
	case TimestampTooOld:
		return "Timestamp was too old"
	case TimestampTooNew:
		return "Timestamp was too new"
	case DailyScrobbledLimitExceeded:
		return "Daily scrobbled limit exceeded"
	default:
		return "Scrobble ignored"
	}
}

type ScrobbleIgnored struct {
	RawMessage string              `xml:",chardata"`
	Code       ScrobbleIgnoredCode `xml:"code,attr"`
}

// Message returns the message for the ignored scrobble. If RawMessage is set,
// it will be returned, other the message will be determined by the code.
//
// The Last.fm API seems to return code 1 (ArtistIgnored) regardless of the
// reason for ignoring the scrobble.
func (s ScrobbleIgnored) Message() string {
	if s.RawMessage != "" {
		return s.RawMessage
	}

	return s.Code.Message()
}

// https://www.last.fm/api/show/track.scrobble#attributes
type ScrobbleResult struct {
	Accepted IntBool  `xml:"accepted,attr"`
	Ignored  IntBool  `xml:"ignored,attr"`
	Scrobble Scrobble `xml:"scrobble"`
}

// https://www.last.fm/api/show/track.scrobble#attributes
type ScrobbleMultiResult struct {
	Accepted  int        `xml:"accepted,attr"`
	Ignored   int        `xml:"ignored,attr"`
	Scrobbles []Scrobble `xml:"scrobble"`
}

// https://www.last.fm/api/show/track.scrobble#attributes
type Scrobble struct {
	Track struct {
		Title     string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"track"`
	Artist struct {
		Name      string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"artist"`
	Album struct {
		Title     string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"album"`
	AlbumArtist struct {
		Name      string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"albumArtist"`
	Timestamp DateTime        `xml:"timestamp"`
	Ignored   ScrobbleIgnored `xml:"ignoredMessage"`
}

// https://www.last.fm/api/show/track.search
type TrackSearchParams struct {
	Track  string `url:"track"`
	Artist string `url:"artist,omitempty"`
	Limit  uint   `url:"limit,omitempty"`
	Page   uint   `url:"page,omitempty"`
}

type TrackSearchResult struct {
	Query struct {
		Role      string `xml:"role,attr"`
		StartPage int    `xml:"startPage,attr"`
	} `xml:"Query"`
	TotalResults int `xml:"totalResults"`
	StartIndex   int `xml:"startIndex"`
	PerPage      int `xml:"itemsPerPage"`
	Tracks       []struct {
		Title  string `xml:"name"`
		Artist string `xml:"artist"`
		// All values returned from the Last.fm API are "FIXME". API issue?
		Streamable string `xml:"streamable"`
		Listeners  int    `xml:"listeners"`
		URL        string `xml:"url"`
		MBID       string `xml:"mbid"`
		Image      Image  `xml:"image"`
	} `xml:"trackmatches>track"`
}

// https://www.last.fm/api/show/track.unlove
type TrackUnloveParams struct {
	Track  string `url:"track"`
	Artist string `url:"artist"`
}

// https://www.last.fm/api/show/track.updateNowPlaying
type UpdateNowPlayingParams struct {
	Artist string `url:"artist"`
	Track  string `url:"track"`

	Album       string   `url:"album,omitempty"`
	AlbumArtist string   `url:"albumArtist,omitempty"`
	TrackNumber int      `url:"trackNumber,omitempty"`
	Duration    Duration `url:"duration,omitempty"`
	MBID        string   `url:"mbid,omitempty"`

	Context string `url:"context,omitempty"`
}

// https://www.last.fm/api/show/track.updateNowPlaying#attributes
type NowPlayingUpdate struct {
	Track struct {
		Title     string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"track"`
	Artist struct {
		Name      string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"artist"`
	Album struct {
		Title     string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"album"`
	AlbumArtist struct {
		Name      string  `xml:",chardata"`
		Corrected IntBool `xml:"corrected,attr"`
	} `xml:"albumArtist"`
	Ignored struct {
		Message string              `xml:",chardata"`
		Code    ScrobbleIgnoredCode `xml:"code,attr"`
	} `xml:"ignoredMessage"`
}
