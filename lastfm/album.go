package lastfm

// https://www.last.fm/api/show/album.addTags
type AlbumAddTagsParams struct {
	Artist string   `url:"artist"`
	Album  string   `url:"album"`
	Tags   []string `url:"tags,comma"`
}

// https://www.last.fm/api/show/album.getInfo
type AlbumInfoParams struct {
	Artist      string `url:"artist"`
	Album       string `url:"album"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	User        string `url:"username,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/album.getInfo
type AlbumInfoMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	User        string `url:"username,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/album.getInfo#attributes
type AlbumInfo struct {
	Title         string `xml:"name"`
	Artist        string `xml:"artist"`
	URL           string `xml:"url"`
	MBID          string `xml:"mbid"`
	Listeners     int    `xml:"listeners"`
	Playcount     int    `xml:"playcount"`
	UserPlaycount *int   `xml:"userplaycount"`
	Image         Image  `xml:"image"`
	Tracks        []struct {
		Title      string   `xml:"name"`
		Number     int      `xml:"rank,attr"`
		URL        string   `xml:"url"`
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
	} `xml:"tracks>track"`
	Tags []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tags>tag"`
	Wiki struct {
		Summary     string   `xml:"summary"`
		Content     string   `xml:"content"`
		PublishedAt DateTime `xml:"published"`
	} `xml:"wiki"`
}

// https://www.last.fm/api/show/album.getTags
type AlbumTagsParams struct {
	Artist      string `url:"artist"`
	Album       string `url:"album"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/album.getTags
type AlbumTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/album.getTags
type AlbumSelfTagsParams struct {
	Artist      string `url:"artist"`
	Album       string `url:"album"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/album.getTags
type AlbumSelfTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

type AlbumTags struct {
	Artist string `xml:"artist,attr"`
	Album  string `xml:"album,attr"`
	Tags   []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/album.getTopTags
type AlbumTopTagsParams struct {
	Artist      string `url:"artist"`
	Album       string `url:"album"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/album.getTopTags
type AlbumTopTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/album.getTopTags#attributes
type AlbumTopTags struct {
	Artist string `xml:"artist,attr"`
	Album  string `xml:"album,attr"`
	Tags   []struct {
		Name  string `xml:"name"`
		URL   string `xml:"url"`
		Count int    `xml:"count"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/album.removeTag
type AlbumRemoveTagParams struct {
	Artist string `url:"artist"`
	Album  string `url:"album"`
	Tag    string `url:"tag"`
}

// https://www.last.fm/api/show/album.search
type AlbumSearchParams struct {
	Album string `url:"album"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type AlbumSearchResult struct {
	For   string `xml:"for,attr"`
	Query struct {
		Role        string `xml:"role,attr"`
		SearchTerms string `xml:"searchTerms,attr"`
		StartPage   int    `xml:"startPage,attr"`
	} `xml:"Query"`
	TotalResults int `xml:"totalResults"`
	StartIndex   int `xml:"startIndex"`
	PerPage      int `xml:"itemsPerPage"`
	Albums       []struct {
		Title      string  `xml:"name"`
		Artist     string  `xml:"artist"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"albummatches>album"`
}
