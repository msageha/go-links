package short_link

import "time"

type ShortLink struct {
	ID                     uint64
	UserID                 uint64
	Shortpath              string
	ShortpathPrefix        string
	DestinationUrl         string
	VisitsCount            uint64
	VisitsCountLastUpdated time.Time
	Key                    string
	Namespace              string
	Type                   string
	DisplayShortpath       string
}

const (
	SHORT_LINK_ROW = "id, user_id, shortpath, shortpath_prefix, destination_url, visits_count, visits_count_last_updated, key, namespace, type, display_shortpath"
)
