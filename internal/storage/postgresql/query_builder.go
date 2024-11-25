package postgresql

import (
	"net/url"
	"strings"
	"time"
)

type Query struct {
	builder strings.Builder
}

func NewQuery(fields []string, s, limit, offset string, vals url.Values) (*Query, error) {
	var builder strings.Builder
	params := make(map[string]string)

	builder.WriteString(s)

	for _, key := range fields {
		// fmt.Println("key: ", key)

		var v string
		var val []string
		var ok bool
		switch key {
		case "group_name", "song_name", "song_text", "link":
			val, ok = vals[key]
			if ok {
				v = val[0]
			}

		case "release_date":
			val, ok = vals[key]
			if ok {
				t, err := time.Parse("02.01.2006", val[0])
				if err != nil {
					return nil, err
				}
				v = t.Format("2006-01-02")
			}
		}

		if v != "" {
			params[key] = v
		}
	}

	filter := ""
	for key, val := range params {
		if filter == "" {
			filter += " WHERE"
		} else {
			filter += " AND"
		}

		filter += " " + key + "='" + val + "'"
	}

	builder.WriteString(filter)
	builder.WriteString(" ORDER BY group_name ")
	builder.WriteString(" LIMIT ")
	builder.WriteString(limit)
	builder.WriteString(" OFFSET ")
	builder.WriteString(offset)
	builder.WriteString(";")
	return &Query{
		builder: builder,
	}, nil

}

func (q *Query) GetQuery() string {
	return q.builder.String()
}
