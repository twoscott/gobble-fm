package lastfm

// https://www.last.fm/api/show/tag.getInfo
type TagInfoParams struct {
	Tag string `url:"tag"`
	// The language to return the info in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

// https://www.last.fm/api/show/tag.getInfo#attributes
type TagInfo struct {
	Name  string `xml:"name"`
	Total int    `xml:"total"`
	Reach int    `xml:"reach"`
	Wiki  struct {
		Summary string `xml:"summary"`
		Content string `xml:"content"`
	} `xml:"wiki"`
}

// https://www.last.fm/api/show/tag.getSimilar
type TagSimilarParams struct {
	Tag string `url:"tag"`
}

type SimilarTags struct {
	Name string `xml:"tag,attr"`
	Tags []struct {
		Name       string  `xml:"name"`
		URL        string  `xml:"url"`
		Streamable IntBool `xml:"streamable"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/tag.getTopAlbums
type TagTopAlbumsParams struct {
	Tag   string `url:"tag"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type TagTopAlbums struct {
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Albums     []struct {
		Title  string `xml:"name"`
		Rank   int    `xml:"rank,attr"`
		URL    string `xml:"url"`
		MBID   string `xml:"mbid"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Cover Image `xml:"image"`
	} `xml:"album"`
}

// https://www.last.fm/api/show/tag.getTopArtists
type TagTopArtistsParams struct {
	Tag   string `url:"tag"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type TagTopArtists struct {
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		Rank       int     `xml:"rank,attr"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/tag.getTopTags
type TagTopTagsParams struct {
	Limit  uint `url:"num_res,omitempty"`
	Offset uint `url:"offset,omitempty"`
}

type TagTopTags struct {
	Offset  int `xml:"offset,attr"`
	Results int `xml:"num_res,attr"`
	Total   int `xml:"total,attr"`
	Tags    []struct {
		Name  string `xml:"name"`
		Count int    `xml:"count"`
		Reach int    `xml:"reach"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/tag.getTopTracks
type TagTopTracksParams struct {
	Tag   string `url:"tag"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type TagTopTracks struct {
	Tag        string `xml:"tag,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Title      string   `xml:"name"`
		Rank       int      `xml:"rank,attr"`
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

// https://www.last.fm/api/show/tag.getWeeklyChartList
type TagWeeklyChartListParams struct {
	Tag string `url:"tag"`
}

type TagWeeklyChartList struct {
	Tag    string `xml:"tag,attr"`
	Charts []struct {
		From DateTime `xml:"from,attr"`
		To   DateTime `xml:"to,attr"`
	} `xml:"chart"`
}
