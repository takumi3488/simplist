package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	// DBに接続する
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	// テーブルを削除する
	if _, err := db.Exec(`DROP TABLE IF EXISTS lists`); err != nil {
		panic(err)
	}

	// テーブルを作成する
	createTable := `CREATE TABLE lists (
		key TEXT PRIMARY KEY, items TEXT[] DEFAULT '{}' NOT NULL
	)`
	if _, err := db.Exec(createTable); err != nil {
		panic(err)
	}

	// データを挿入する
	for _, k := range []string{"foo", "bar", "baz"} {
		if _, err := db.Exec(`INSERT INTO lists (key) VALUES ($1)`, k); err != nil {
			panic(err)
		}
	}

	return db
}

func resetDB(*sql.DB) {
	// DBに接続する
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	// テーブルを削除する
	if _, err := db.Exec(`DROP TABLE IF EXISTS lists`); err != nil {
		panic(err)
	}

	// データベースを閉じる
	db.Close()
}

func TestGetLists(t *testing.T) {
	db := initDB()
	defer resetDB(db)
	h := &Handler{db}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/lists", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, h.GetLists(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[{"key":"foo","items":[]},{"key":"bar","items":[]},{"key":"baz","items":[]}]`, rec.Body.String())
	}
}

func TestGetList(t *testing.T) {
	db := initDB()
	defer resetDB(db)
	h := &Handler{db}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/lists/foo", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("key")
	c.SetParamValues("foo")
	if assert.NoError(t, h.GetList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"key":"foo","items":[]}`, rec.Body.String())
	}
}

func TestUpdateList(t *testing.T) {
	db := initDB()
	defer resetDB(db)
	h := &Handler{db}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/lists/foo", strings.NewReader(`{"items":["a","b","c"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("key")
	c.SetParamValues("foo")
	if assert.NoError(t, h.UpdateList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Body.String())

	}
}
