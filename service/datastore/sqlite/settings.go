package sqlite

import (
	"database/sql"
	"github.com/rs/xid"
	"github.com/slimjim777/snap-downloader/domain"
	"log"
)

const createSettingsTableSQL string = `
	CREATE TABLE IF NOT EXISTS settings (
		id               varchar(200) primary key not null,
		key              varchar(200) not null,
		name             varchar(200) not null,
		data             text default '',
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp
	)
`

const addSettingSQL = `
	INSERT INTO settings (id, key, name, data) VALUES ($1, $2, $3, $4)
`
const updateSettingSQL = `
	UPDATE settings SET data=$1, modified=current_timestamp WHERE key=$2 and name=$3
`
const getSettingSQL = `
	SELECT id, key, name, data, created, modified
	FROM settings
	WHERE key=$1 and name=$2
`

// SettingsPut stores a new config setting
func (db *DB) SettingsPut(key, name, data string) (string, error) {
	// check if the setting exists
	set, err := db.SettingsGet(key, name)
	if err != nil {
		// does not exist, so create it
		id := xid.New()
		_, err := db.Exec(addSettingSQL, id.String(), key, name, data)
		return id.String(), err
	}

	// update it
	_, err = db.Exec(updateSettingSQL, data, key, name)
	return set.ID, err
}

// SettingsGet fetches an existing config setting
func (db *DB) SettingsGet(key, name string) (domain.ConfigSetting, error) {
	r := domain.ConfigSetting{}
	err := db.QueryRow(getSettingSQL, key, name).Scan(&r.ID, &r.Key, &r.Name, &r.Data, &r.Created, &r.Modified)
	switch {
	case err == sql.ErrNoRows:
		return r, err
	case err != nil:
		log.Printf("Error retrieving database repo: %v\n", err)
		return r, err
	}
	return r, nil
}
