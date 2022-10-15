package models

import (
	"context"

	"github.com/calvinsomething/go-proj/db"
)

type (
	// Player reflects a row in the players table.
	Player struct {
		IP          string  `json:"ip"`
		Faction     string  `json:"faction"`
		Race        string  `json:"race" validate:"oneof=dwarf gnome human night elf orc tauren troll undead"`
		Class       string  `json:"class" validate:"oneof=druid hunter mage paladin priest rogue shaman warlock warrior"`
		Profession1 *string `json:"profession1" validate:"oneof=alchemy blacksmithing enchanting engineering herbalism mining tailoring"`
		Profession2 *string `json:"profession2" validate:"oneof=alchemy blacksmithing enchanting engineering herbalism mining tailoring,nefield=Profession1"`
		WeeklyHours *int    `json:"weeklyHours" validate:"gt=0,lt=51"`
	}
)

// GetPlayer gets the Player associated with the ip address.
func GetPlayer(ctx context.Context, ip string) (p *Player, err error) {
	p.IP = ip
	err = db.Pool.QueryRowContext(ctx, `
		SELECT faction, race, class, profession1, profession2, weekly_hours
		FROM players
		WHERE ip = ?;
	`, ip).Scan(&p.Faction, &p.Race, &p.Class, &p.Profession1, &p.Profession2, p.WeeklyHours)
	return
}

// GetPlayers returns all players in the db.
func GetPlayers(ctx context.Context) ([]*Player, error) {
	rows, err := db.Pool.QueryContext(ctx, `
		SELECT ip, faction, race, class, profession1, profession2, weekly_hours
		FROM players;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*Player
	for rows.Next() {
		p := &Player{}
		err = rows.Scan(&p.IP, &p.Faction, &p.Race, &p.Class, &p.Profession1, &p.Profession2, &p.WeeklyHours)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

// Save upserts the Player into the db.
func (p *Player) Save(ctx context.Context) error {
	return db.Pool.MustAffect(ctx, `
		INSERT INTO players (ip, race, class, profession1, profession2, weekly_hours)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE players ip = ?, race = ?, class = ?, profession1 = ?,
			profession2 = ?, weekly_hours = ?
	`, p.IP, p.Race, p.Class, p.Profession1, p.Profession2, p.WeeklyHours, p.IP, p.Race,
		p.Class, p.Profession1, p.Profession2, p.WeeklyHours)
}
