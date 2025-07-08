package persistent

import (
	"database/sql"
	"dwld-downloader/internal/entity"
	"dwld-downloader/pkg/sqldb"
	"fmt"
)

type persistentRepo struct {
	db      *sqldb.DB
	workDir string
}

func NewSQLRepo(db *sqldb.DB, workDir string) SQLRepo {
	resp := &persistentRepo{
		db:      db,
		workDir: workDir,
	}

	resp.workToNew()

	return resp
}

func (p *persistentRepo) workToNew() {
	_, err := p.db.Exec("update links set work_status = $1 where work_status = $2", entity.StatusMapping[entity.NEW], entity.StatusMapping[entity.WORK])
	if err != nil {
		fmt.Println(err)
	}
}

func (p *persistentRepo) SelectHistory(withoutStatus *entity.Status) ([]LinkModel, error) {
	var rows *sql.Rows
	var err error

	if withoutStatus != nil {
		rows, err = p.db.Select("select link, filename, work_status, message, target_quantity from links where work_status != $1", entity.StatusMapping[*withoutStatus])
	} else {
		rows, err = p.db.Select("select link, filename, work_status, message, target_quantity from links")
	}

	defer func() {
		rows.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("SelectHistory: %w", err)
	}

	resp := make([]LinkModel, 0)
	var row LinkModel
	for rows.Next() {
		err := rows.Scan(&row.Link, &row.Filename, &row.WorkStatus, &row.Message, &row.TargetQuantity)
		if err != nil {
			fmt.Println(err)
		}

		resp = append(resp, LinkModel{
			Link:           row.Link,
			Filename:       row.Filename,
			WorkStatus:     row.WorkStatus,
			Message:        row.Message,
			TargetQuantity: row.TargetQuantity,
		})
	}

	return resp, nil
}

func (p *persistentRepo) Insert(link string, maxQuality int) ([]LinkModel, error) {
	_, err := p.db.Exec("insert into links (link, filename, target_quantity, work_status) values($1, $2, $3, $4);", link, "COMING SOON", maxQuality, entity.StatusMapping[entity.NEW])
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory(nil)
}

func (p *persistentRepo) UpdateStatus(link string, status entity.Status) ([]LinkModel, error) {
	_, err := p.db.Exec("update links set work_status = $1 where link = $2;", entity.StatusMapping[status], link)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory(nil)
}

func (p *persistentRepo) Delete(link string) ([]LinkModel, error) {
	_, err := p.db.Exec("delete from links where link = $1;", link)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory(nil)
}

func (p *persistentRepo) DeleteHistory() ([]LinkModel, error) {
	_, err := p.db.Exec("delete from links where work_status = $1;", entity.StatusMapping[entity.DONE])
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory(nil)
}

func (p *persistentRepo) SelectOne(status entity.Status) (*LinkModel, error) {
	rows, err := p.db.Select(`select link, filename, work_status, message, target_quantity from links
	 where work_status = $1 order by RANDOM() limit 1;`, entity.StatusMapping[status])
	defer func() {
		rows.Close()
	}()

	if err != nil {
		return nil, fmt.Errorf("db.SelectOne(query): %w", err)
	}

	isNext := rows.Next()
	if !isNext {
		return nil, nil
	}

	row := &LinkModel{}

	err = rows.Scan(&row.Link, &row.Filename, &row.WorkStatus, &row.Message, &row.TargetQuantity)
	if err != nil {
		return nil, fmt.Errorf("db.SelectOne(Scan): %w", err)
	}

	return row, nil
}

func (p *persistentRepo) Update(l *LinkModelRequest) error {
	_, err := p.db.Exec(`update links 
	set 
	work_status = $1,
	filename = $2,
    message = $3,
    target_quantity = $4
 	where link = $5;`, entity.StatusMapping[l.WorkStatus], *l.Filename, *l.Message, l.TargetQuantity, l.Link)
	if err != nil {
		return fmt.Errorf("update link: %w", err)
	}

	return nil
}
