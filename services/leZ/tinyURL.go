package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"web-back-go/pkg/structs"
	tinyURL "web-back-go/pkg/tinyURL"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func (app *App) getShortCode(c *gin.Context) {
	baseURL := c.DefaultQuery("url", "NULL")
	if baseURL == "NULL" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Missing origin url",
		})
		return
	}

	shortCode, err := getShortCodeFromDB(app.db, baseURL)
	if err != nil {
		log.Errorf("Could not get short code from db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code": shortCode,
	})
}

func getShortCodeFromDB(db *sqlx.DB, baseURL string) (string, error) {
	query := `SELECT md5 FROM tiny_url WHERE origin_url = ?`

	row := db.QueryRow(query, baseURL)

	shortCode := ""
	if err := row.Scan(&shortCode); err != nil {
		if err == sql.ErrNoRows {
			shortCode, err := setShortCodeInDB(db, baseURL)
			if err == nil {
				return shortCode, nil
			}
		}
		return "", err
	}

	return shortCode, nil
}

func (app *App) getBaseURL(c *gin.Context) {
	shortCode := c.DefaultQuery("shortCode", "NULL")
	if shortCode == "NULL" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Missing short code",
		})
		return
	}

	baseURL, err := getBaseURLFromDB(app.db, shortCode)
	if err != nil {
		log.Errorf("Could not get base URL from db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base_url": baseURL,
	})
}

func getBaseURLFromDB(db *sqlx.DB, shortCode string) (string, error) {
	query := `SELECT origin_url FROM tiny_url WHERE md5 = ?`

	row := db.QueryRow(query, shortCode)

	baseURL := ""
	if err := row.Scan(&baseURL); err != nil {
		return "", err
	}

	return baseURL, nil
}

func setShortCodeInDB(db *sqlx.DB, baseURL string) (string, error) {
	shortCode := tinyURL.GeneragteMD5Value(baseURL)

	query := `INSERT INTO tiny_url (origin_url, md5) VALUES (?, ?)`

	result := db.MustExec(query, baseURL, shortCode)

	lastInsertId, _ := result.LastInsertId()

	log.Infof("Successfullt inserted a record in DB : %d, %s, %s", lastInsertId, baseURL, shortCode)

	return shortCode, nil
}

func (app *App) listAllTinyURL(c *gin.Context) {
	rows, err := app.db.Query(`SELECT id, origin_url, md5 FROM tiny_url`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
	}

	defer func() { _ = rows.Close() }()

	var (
		tinyURLs  []structs.TinyURL
		id        int
		baseURL   string
		shortCode string
	)

	for rows.Next() {
		if err := rows.Scan(&id, &baseURL, &shortCode); err != nil {
			log.Errorf("Error occurred when scanning tinyURLs : %v", err)
			continue
		}

		tinyURLs = append(tinyURLs, structs.TinyURL{
			ID:        id,
			BaseURL:   baseURL,
			ShortCode: shortCode,
		})
	}

	urlsStr, _ := json.Marshal(tinyURLs)

	c.JSON(http.StatusOK, gin.H{
		"urls": string(urlsStr),
	})
}
