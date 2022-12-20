package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type TrackSearchResult struct {
	Tracks struct {
		Hits []struct {
			Track struct {
				Key      string
				Title    string
				Subtitle string
				Plays    int
			}
		}
	}
}

func SearchHandler(c *gin.Context) {
	// get the query param for search
	songName := c.Request.URL.Query().Get("name")
	result := &TrackSearchResult{}
	get(fmt.Sprintf("/search?term=%s", songName), result)
	for i, hit := range result.Tracks.Hits {
		result.Tracks.Hits[i].Track.Plays = getPlays(hit.Track.Key)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result.Tracks.Hits,
	})
}

// example: 609790719
func PlaysHandler(c *gin.Context) {
	// get the id from the route
	songID := c.Param("id")
	// return the result
	c.JSON(http.StatusOK, gin.H{
		"number_of_plays": getPlays(songID),
	})
}

func getPlays(songID string) int {
	result := &struct {
		ID    string
		Total int
		Type  string
	}{}
	/// call the shazam api to get number of plays
	get(fmt.Sprintf("/songs/get-count?key=%s", songID), result)
	return result.Total
}
