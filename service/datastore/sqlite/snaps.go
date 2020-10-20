package sqlite

import (
	"github.com/rs/xid"
	"github.com/slimjim777/snap-downloader/domain"
)

const createSnapsTableSQL string = `
	CREATE TABLE IF NOT EXISTS snaps (
		id               varchar(200) primary key not null,
		name             varchar(200) not null,
		arch             varchar(200) not null,
		created          timestamp default current_timestamp,
		modified         timestamp default current_timestamp
	)
`
const indexSnapsSQL = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_name_arch on snaps (name,arch)
`
const addSnapSQL = `
	INSERT INTO snaps (id, name, arch) VALUES ($1, $2, $3)
`
const deleteSnapSQL = `
	DELETE FROM snaps WHERE id=$1
`
const listSnapSQL = `
	SELECT id, name, arch, created, modified
	FROM snaps
`

// SnapsCreate adds a new snap to cache
func (db *DB) SnapsCreate(name, arch string) (string, error) {
	id := xid.New()
	_, err := db.Exec(addSnapSQL, id.String(), name, arch)
	return id.String(), err
}

// SnapsDelete removes a snap from the cache
func (db *DB) SnapsDelete(id string) error {
	_, err := db.Exec(deleteSnapSQL, id)
	return err
}

// SnapsList fetches the snaps that are to be cached
func (db *DB) SnapsList() ([]domain.SnapCache, error) {
	records := []domain.SnapCache{}
	rows, err := db.Query(listSnapSQL)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		r := domain.SnapCache{}
		err := rows.Scan(&r.ID, &r.Name, &r.Arch, &r.Created, &r.Modified)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}

	return records, nil
}
