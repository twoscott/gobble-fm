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
	MBID        string `url:"mbid"`
	AutoCorrect bool   `url:"autocorrect,int,omitempty"`
	User        string `url:"username,omitempty"`
	// The language to return the biography in, as an ISO 639 alpha-2 code.
	Language string `url:"lang,omitempty"`
}

type AlbumInfo struct {
	Name      string `xml:"name"`
	Artist    string `xml:"artist"`
	URL       string `xml:"url"`
	MBID      string `xml:"mbid"`
	Listeners int    `xml:"listeners"`
	Playcount int    `xml:"playcount"`
	Image     Image  `xml:"image"`
	Tracks    []struct {
		Name       string   `xml:"name"`
		Number     int      `xml:"rank,attr"`
		URL        string   `xml:"url"`
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
	} `xml:"tracks>track"`
	Tag []struct {
		Name string `xml:"name"`
		URL  string `xml:"url"`
	} `xml:"tags>tag"`
	Wiki struct {
		Summary   string   `xml:"summary"`
		Content   string   `xml:"content"`
		Published DateTime `xml:"published"`
	} `xml:"wiki"`
}
