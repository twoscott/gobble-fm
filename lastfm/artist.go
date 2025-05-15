package lastfm

// https://www.last.fm/api/show/artist.addTags
type ArtistAddTagsParams struct {
	Artist string   `url:"artist"`
	Tags   []string `url:"tags,comma"`
}

// https://www.last.fm/api/show/artist.getCorrection
type ArtistCorrectionParams struct {
	Artist string `url:"artist"`
}

type ArtistCorrection struct {
	Corrections []struct {
		Index  int `xml:"index,attr"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
	} `xml:"correction"`
}

// https://www.last.fm/api/show/artist.getInfo
type ArtistInfoParams struct {
	Artist      string `url:"artist"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/artist.getInfo
type ArtistInfoMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/artist.getInfo
type ArtistUserInfoParams struct {
	Artist      string `url:"artist"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/artist.getInfo
type ArtistUserInfoMBIDParams struct {
	MBID        string `url:"mbid"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

type ArtistInfo struct {
	Name           string  `xml:"name"`
	URL            string  `xml:"url"`
	MBID           string  `xml:"mbid"`
	Image          Image   `xml:"image"`
	Listeners      int     `xml:"stats>listeners"`
	Playcount      int     `xml:"stats>playcount"`
	Streamable     IntBool `xml:"streamable"`
	OnTour         IntBool `xml:"ontour"`
	SimilarArtists []struct {
		Name  string `xml:"name"`
		URL   string `xml:"url"`
		Image Image  `xml:"image"`
	} `xml:"similar>artist"`
	Tags []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tags>tag"`
	Bio struct {
		Links []struct {
			URL      string `xml:"href,attr"`
			Relation string `xml:"rel,attr"`
		} `xml:"links>link"`
		Summary     string   `xml:"summary"`
		Content     string   `xml:"content"`
		PublishedAt DateTime `xml:"published"`
	} `xml:"bio"`
}

type ArtistUserInfo struct {
	ArtistInfo
	UserPlaycount int `xml:"stats>userplaycount"`
}

// https://www.last.fm/api/show/artist.getSimilar
type ArtistSimilarParams struct {
	Artist      string `url:"artist"`
	Limit       uint   `url:"limit,omitempty"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getSimilar
type ArtistSimilarMBIDParams struct {
	MBID        string `url:"mbid"`
	Limit       uint   `url:"limit,omitempty"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getSimilar#attributes
type SimilarArtists struct {
	Artist  string `xml:"artist,attr"`
	Artists []struct {
		Name       string  `xml:"name"`
		Match      float64 `xml:"match"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/artist.getTags
type ArtistTagsParams struct {
	Artist      string `url:"artist"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getTags
type ArtistTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	User        string `url:"username"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getTags
type ArtistSelfTagsParams struct {
	Artist      string `url:"artist"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getTags
type ArtistSelfTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

type ArtistTags struct {
	Artist string `xml:"artist,attr"`
	Tags   []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/artist.getTopAlbums
type ArtistTopAlbumsParams struct {
	Artist      string `url:"artist"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
	Page        uint   `url:"page,omitempty"`
}

// https://www.last.fm/api/show/artist.getTopAlbums
type ArtistTopAlbumsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
	Page        uint   `url:"page,omitempty"`
}

type ArtistTopAlbums struct {
	Artist     string `xml:"artist,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Albums     []struct {
		Title     string `xml:"name"`
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

// https://www.last.fm/api/show/artist.getTopTags
type ArtistTopTagsParams struct {
	Artist      string `url:"artist"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

// https://www.last.fm/api/show/artist.getTopTags
type ArtistTopTagsMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
}

type ArtistTopTags struct {
	Artist string `xml:"artist,attr"`
	Tags   []struct {
		Name  string `xml:"name"`
		Count int    `xml:"count"`
		URL   string `xml:"url"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/artist.getTopTracks
type ArtistTopTracksParams struct {
	Artist      string `url:"artist"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
	Page        uint   `url:"page,omitempty"`
}

// https://www.last.fm/api/show/artist.getTopTracks
type ArtistTopTracksMBIDParams struct {
	MBID        string `url:"mbid"`
	AutoCorrect *bool  `url:"autocorrect,int,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
	Page        uint   `url:"page,omitempty"`
}

type ArtistTopTracks struct {
	Artist     string `xml:"artist,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title      string  `xml:"name"`
		Rank       int     `xml:"rank,attr"`
		Playcount  int     `xml:"playcount"`
		Listeners  int     `xml:"listeners"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Artist     struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Image Image `xml:"image"`
	} `xml:"track"`
}

// https://www.last.fm/api/show/artist.removeTag
type ArtistRemoveTagParams struct {
	Artist string `url:"artist"`
	Tag    string `url:"tag"`
}

// https://www.last.fm/api/show/artist.search
type ArtistSearchParams struct {
	Artist string `url:"artist"`
	Limit  uint   `url:"limit,omitempty"`
	Page   uint   `url:"page,omitempty"`
}

type ArtistSearchResult struct {
	For   string `xml:"for,attr"`
	Query struct {
		Role        string `xml:"role,attr"`
		SearchTerms string `xml:"searchTerms,attr"`
		StartPage   int    `xml:"startPage,attr"`
	} `xml:"Query"`
	TotalResults int `xml:"totalResults"`
	StartIndex   int `xml:"startIndex"`
	PerPage      int `xml:"itemsPerPage"`
	Artists      []struct {
		Name       string  `xml:"name"`
		Listeners  int     `xml:"listeners"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artistmatches>artist"`
}
