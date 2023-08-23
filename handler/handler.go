package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type List struct {
	Key   string   `json:"key"`
	Items []string `json:"items"`
}
type ListCreateRequest struct {
	Items []string `json:"items"`
}
type Handler struct {
	DB *sql.DB
}

func (h *Handler) GetLists(c echo.Context) error {
	log.Println("GetLists")

	var lists []List
	rows, err := h.DB.Query(`SELECT * FROM lists`)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var items []string
		if err := rows.Scan(&key, pq.Array(&items)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		lists = append(lists, List{
			Key:   key,
			Items: items,
		})
	}
	return c.JSON(http.StatusOK, lists)
}

func (h *Handler) GetList(c echo.Context) error {
	log.Panicln("GetList")

	key := c.Param("key")
	var items []string
	if err := h.DB.QueryRow(`SELECT items FROM lists WHERE key = $1`, key).Scan(pq.Array(&items)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	list := List{
		Key:   key,
		Items: items,
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateList(c echo.Context) error {
	log.Println("UpdateList")

	key := c.Param("key")
	var req ListCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	items := req.Items
	if _, err := h.DB.Exec(`UPDATE lists SET items = $2 WHERE key = $1`, key, pq.Array(&items)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
