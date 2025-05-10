package lastfm

// https://www.last.fm/api/show/chart.getTopArtists
type ChartTopArtistsParams struct {
	Limit uint `url:"limit,omitempty"`
	Page  uint `url:"page,omitempty"`
}

type ChartTopArtists struct {
	Page       int `xml:"page,attr"`
	PerPage    int `xml:"perPage,attr"`
	TotalPages int `xml:"totalPages,attr"`
	Total      int `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		Playcount  int     `xml:"playcount"`
		Listeners  int     `xml:"listeners"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/chart.getTopTags
type ChartTopTagsParams struct {
	Limit uint `url:"limit,omitempty"`
	Page  uint `url:"page,omitempty"`
}

type ChartTopTags struct {
	Page       int `xml:"page,attr"`
	PerPage    int `xml:"perPage,attr"`
	TotalPages int `xml:"totalPages,attr"`
	Total      int `xml:"total,attr"`
	Tags       []struct {
		Name       string  `xml:"name"`
		URL        string  `xml:"url"`
		Reach      int     `xml:"reach"`
		Count      int     `xml:"taggings"`
		Streamable IntBool `xml:"streamable"`
		Wiki       string  `xml:"wiki"`
	} `xml:"tag"`
}

// https://www.last.fm/api/show/chart.getTopTracks
type ChartTopTracksParams struct {
	Limit uint `url:"limit,omitempty"`
	Page  uint `url:"page,omitempty"`
}

type ChartTopTracks struct {
	Page       int `xml:"page,attr"`
	PerPage    int `xml:"perPage,attr"`
	TotalPages int `xml:"totalPages,attr"`
	Total      int `xml:"total,attr"`
	Tracks     []struct {
		Title      string   `xml:"name"`
		Duration   Duration `xml:"duration"`
		Playcount  int      `xml:"playcount"`
		Listeners  int      `xml:"listeners"`
		URL        string   `xml:"url"`
		MBID       string   `xml:"mbid"`
		Streamable struct {
			Preview   IntBool `xml:",chardata"`
			Fulltrack IntBool `xml:"fulltrack,attr"`
		} `xml:"streamable"`
		Artist struct {
			Name string `xml:"name"`
			URL  string `xml:"url"`
			MBID string `xml:"mbid"`
		} `xml:"artist"`
		Image Image `xml:"image"`
	} `xml:"track"`
}
