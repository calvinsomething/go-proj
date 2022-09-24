package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/calvinsomething/go-proj/auth"
	"github.com/calvinsomething/go-proj/db"
)

type player struct {
	IP          string         `json:"ip"`
	Faction     string         `json:"faction"`
	Race        string         `json:"race"`
	Class       string         `json:"class"`
	Profession1 sql.NullString `json:"profession1"`
	Profession2 sql.NullString `json:"profession2"`
	WeeklyHours int            `json:"weeklyHours"`
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpErr(w, 400, err)
		return
	}
	defer r.Body.Close()

	var login login

	if err = json.Unmarshal(body, &login); err != nil {
		httpErr(w, 400, err)
		return
	}

	sid, err := auth.LogIn(ctx, login.Email, login.Password)
	if err == auth.ErrBadLogin {
		httpErr(w, 400, err)
		return
	} else if err == auth.ErrBadMAC {
		httpErr(w, 401, err)
		return
	} else if err != nil {
		httpErr(w, 500, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sid,
		MaxAge:   int(auth.SessionMaxAge),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.QueryContext(r.Context(), `SELECT * FROM players`)
	if err != nil {
		httpErr(w, 500, err)
		return
	}
	defer rows.Close()

	players := make([]player, 0, 5)

	for rows.Next() {
		var p player
		err := rows.Scan(&p.IP, &p.Faction, &p.Race, &p.Class, &p.Profession1, &p.Profession2, &p.WeeklyHours)
		if err != nil {
			httpErr(w, 500, err)
			return
		}
		players = append(players, p)
	}

	returnJSON(w, players)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpErr(w, 400, err)
		return
	}
	defer r.Body.Close()

	var login login
	if err = json.Unmarshal(body, &login); err != nil {
		httpErr(w, 400, err)
		return
	}

	err = auth.CreateUser(r.Context(), login.Email, login.Password)
	if err != nil {
		httpErr(w, 500, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// temp...
func errHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(struct {
		Err string
	}{"test err"})
	if err != nil {
		log.Println(err)
	}
}
