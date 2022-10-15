package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/calvinsomething/go-proj/auth"
	"github.com/calvinsomething/go-proj/db"
	"github.com/calvinsomething/go-proj/models"
)

type (
	login struct {
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"min=8,max=20,password"`
	}
)

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

	if err = validate.Struct(&login); err != nil {
		validationErr(w, err)
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

func getPlayersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.QueryContext(r.Context(), `SELECT * FROM players`)
	if err != nil {
		httpErr(w, 500, err)
		return
	}
	defer rows.Close()

	players := make([]models.Player, 0, 5)

	for rows.Next() {
		var p models.Player
		err := rows.Scan(&p.IP, &p.Faction, &p.Race, &p.Class, &p.Profession1, &p.Profession2, &p.WeeklyHours)
		if err != nil {
			httpErr(w, 500, err)
			return
		}
		players = append(players, p)
	}

	returnJSON(w, players)
}

func addPlayerHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		httpErr(w, 400, err)
		return
	}

	var player models.Player
	if err = json.Unmarshal(payload, &player); err != nil {
		httpErr(w, 400, err, "problem reading player data")
		return
	}

	if err = validate.Struct(&player); err != nil {
		validationErr(w, err)
		return
	}

	player.IP = r.RemoteAddr[:strings.LastIndexByte(r.RemoteAddr, ':')]

	ctx := r.Context()

	if err = player.Save(ctx); err != nil {
		httpErr(w, 500, err, "could not save player data")
		return
	}

	savedPlayer, err := models.GetPlayer(ctx, player.IP)
	if err != nil {
		httpErr(w, 500, err, "could not retrieve saved player")
	}

	w.WriteHeader(http.StatusCreated)
	returnJSON(w, &savedPlayer)
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
