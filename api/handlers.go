package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type Track struct {
	Key      string
	Title    string
	Subtitle string
	Plays    int
}

type TrackSearchResult struct {
	Tracks struct {
		Hits []struct {
			Track Track
		}
	}
}

func SearchHandler(c *gin.Context) {
	// get the query param for search
	songName := c.Request.URL.Query().Get("name")
	result := &TrackSearchResult{}
	err := get(fmt.Sprintf("/search?term=%s", songName), result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	for i, hit := range result.Tracks.Hits {
		plays, err := getPlays(hit.Track.Key)
		if err != nil {
			log.Printf("error getting plays for %s: %v", hit.Track.Key, err)
		}
		result.Tracks.Hits[i].Track.Plays = plays
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
	numOfPlays, err := getPlays(songID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"number_of_plays": numOfPlays,
	})
}

func getPlays(songID string) (int, error) {
	result := &struct {
		ID    string
		Total int
		Type  string
	}{}
	// call the shazam api to get number of plays
	err := get(fmt.Sprintf("/songs/get-count?key=%s", songID), result)
	return result.Total, err
}
