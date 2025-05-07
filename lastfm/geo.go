package lastfm

// https://www.last.fm/api/show/geo.getTopArtists

type GeoTopArtistsParams struct {
	// A country name, as defined by the ISO 3166-1 country names standard
	Country string `url:"country"`
	Limit   uint   `url:"limit,omitempty"`
	Page    uint   `url:"page,omitempty"`
}

type GeoTopArtists struct {
	Country    string `xml:"country,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		Listeners  int     `xml:"listeners"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}

// https://www.last.fm/api/show/geo.getTopTracks

type GeoTopTracksParams struct {
	// A country name, as defined by the ISO 3166-1 country names standard
	Country string `url:"country"`
	// A metro name to fetch the charts for, (must be within the country
	// specified)
	Location string `url:"location,omitempty"`
	Limit    uint   `url:"limit,omitempty"`
	Page     uint   `url:"page,omitempty"`
}

type GeoTopTracks struct {
	Country    string `xml:"country,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Tracks     []struct {
		Name       string   `xml:"name"`
		Rank       int      `xml:"rank,attr"`
		Listeners  int      `xml:"listeners"`
		URL        string   `xml:"url"`
		MBID       string   `xml:"mbid"`
		Duration   Duration `xml:"duration"`
		Streamable struct {
			Value     IntBool `xml:",chardata"`
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
