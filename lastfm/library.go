package lastfm

// https://www.last.fm/api/show/library.getArtists
type LibraryArtistsParams struct {
	User  string `url:"user"`
	Limit uint   `url:"limit,omitempty"`
	Page  uint   `url:"page,omitempty"`
}

type LibraryArtists struct {
	User       string `xml:"user,attr"`
	Page       int    `xml:"page,attr"`
	PerPage    int    `xml:"perPage,attr"`
	TotalPages int    `xml:"totalPages,attr"`
	Total      int    `xml:"total,attr"`
	Artists    []struct {
		Name       string  `xml:"name"`
		Playcount  int     `xml:"playcount"`
		Tagcount   int     `xml:"tagcount"`
		URL        string  `xml:"url"`
		MBID       string  `xml:"mbid"`
		Streamable IntBool `xml:"streamable"`
		Image      Image   `xml:"image"`
	} `xml:"artist"`
}
