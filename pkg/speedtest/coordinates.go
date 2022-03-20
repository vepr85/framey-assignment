package speedtest

import (
	"framey/assignment/internal/geo"
	"sort"
)

func SortServersByDistance(servers []Server, org geo.Coordinates) map[ServerID]geo.Kilometers {
	m := make(map[ServerID]geo.Kilometers)
	for _, s := range servers {
		if _, ok := m[s.ID]; ok {
			continue
		}
		m[s.ID] = org.DistanceTo(s.Coordinates)
	}

	sort.Slice(servers, func(i, j int) bool {
		return m[servers[i].ID].Less(m[servers[j].ID])
	})

	return m
}
