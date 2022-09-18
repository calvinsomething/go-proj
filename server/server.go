package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/calvinsomething/go-proj/auth"
	"github.com/calvinsomething/go-proj/db"
)

var (
	hostPort string
)

type (
	handlerT func(http.ResponseWriter, *http.Request)

	muxT struct {
		mux        *http.ServeMux
		middleware []handlerT
	}
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Server
	hostPort = os.Getenv("SERVER_PORT")
	// DB
	db.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	defer db.Pool.Close()

	execArgs()

	log.Printf("Listening on port %s...\n", hostPort)
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}

func execArgs() {
	i := 1
	shouldExit := false
	for len(os.Args) > i {
		switch os.Args[i] {
		case "migrate":
			shouldExit = true
			i++
			if len(os.Args) > i {
				if os.Args[i] == "up" {
					log.Println("Migrating database up...")
					if err := db.Pool.Migrate(); err != nil {
						log.Fatal(err)
					}
				} else if os.Args[i] == "down" {
					log.Println("Migrating database down...")
					if err := db.Pool.Migrate(true); err != nil {
						log.Fatal(err)
					}
				} else if steps, err := strconv.ParseInt(os.Args[i], 10, 32); err != nil {
					log.Fatalf("Invalid migrate option: %s", os.Args[i])
				} else {
					log.Printf("Migrating %d steps...", steps)
					if err = db.Pool.MigrateSteps(int(steps)); err != nil {
						log.Fatal(err)
					}
				}
			} else {
				log.Fatal("Missing migrate option")
			}
		default:
			log.Fatalf("Unknown command line argument: '%s'", os.Args[i])
		}
		i++
	}
	if shouldExit {
		os.Exit(0)
	}
}

// Custom Mux
func newMux() *muxT {
	middleware := []handlerT{
		logger,
	}
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/err", errHandler)
	mux.HandleFunc("/data", dataHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)

	return &muxT{mux, middleware}
}

func (m *muxT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mw := range m.middleware {
		mw(w, r)
	}
	m.mux.ServeHTTP(w, r)
}

// Middleware
func logger(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Header)
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

func httpErr(w http.ResponseWriter, statusCode int, err error, msg ...string) {
	if err != nil {
		log.Output(2, err.Error())
	}

	w.WriteHeader(statusCode)

	if len(msg) != 0 {
		w.Write([]byte(msg[0]))
	}
}

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

func returnJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		httpErr(w, 500, err)
	}
	w.Write(j)
}
