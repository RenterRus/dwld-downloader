package persistent

import (
	"dwld-downloader/internal/entity"
	"dwld-downloader/pkg/sqldb"
	"fmt"
)

type persistentRepo struct {
	db      *sqldb.DB
	workDir string
}

func NewSQLRepo(db *sqldb.DB, workDir string) SQLRepo {
	return &persistentRepo{
		db:      db,
		workDir: workDir,
	}
}

func (p *persistentRepo) SelectHistory() ([]LinkModel, error) {
	rows, err := p.db.Select("select link, filename, path, work_status, stage_config, retry, message, target_quantity from links")
	defer func() {
		rows.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("SelectHistory: %w", err)
	}

	resp := make([]LinkModel, 0)
	var row LinkModel
	for rows.Next() {
		err := rows.Scan(&row.Link, &row.Filename, &row.Path, &row.WorkStatus, &row.StageConfig, &row.Retry, &row.Message, &row.TargetQuantity)
		if err != nil {
			fmt.Println(err)
		}

		resp = append(resp, LinkModel{
			Link:           row.Link,
			Filename:       row.Filename,
			Path:           row.Path,
			WorkStatus:     row.WorkStatus,
			StageConfig:    row.StageConfig,
			Message:        row.Message,
			Retry:          row.Retry,
			TargetQuantity: row.TargetQuantity,
		})
	}

	return resp, nil
}

func (p *persistentRepo) Insert(link string, maxQuality int) ([]LinkModel, error) {
	_, err := p.db.Exec("insert into links (link, target_quantity, work_status, path) values($1, $2, $3, $4);", link, maxQuality, entity.StatusMapping[entity.NEW], p.workDir)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory()
}

func (p *persistentRepo) UpdateStatus(link string, status entity.Status) ([]LinkModel, error) {
	_, err := p.db.Exec("update links set work_status = $1 where link = $2;", entity.StatusMapping[status], link)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory()
}

func (p *persistentRepo) Delete(link string) ([]LinkModel, error) {
	_, err := p.db.Exec("delete from links where link = $1;", link)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory()
}

func (p *persistentRepo) DeleteHistory() ([]LinkModel, error) {
	_, err := p.db.Exec("delete from links where work_status = $1;", entity.StatusMapping[entity.DONE])
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}

	return p.SelectHistory()
}

func (p *persistentRepo) SelectOne() (*LinkModel, error) {
	rows, err := p.db.Select(`select link, filename, path, work_status, stage_config, retry, message, target_quantity from links
	 where work_status = $1 order by RANDOM() limit 1;`, entity.StatusMapping[entity.NEW])
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

	err = rows.Scan(&row.Link, &row.Filename, &row.Path, &row.WorkStatus, &row.StageConfig, &row.Retry, &row.Message, &row.TargetQuantity)
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
    path = $3,
    stage_config = $4,
    message = $5,
    retry = $6,
	target_quantity = $7
 	where link = $8;`, entity.StatusMapping[l.WorkStatus], *l.Filename, l.Path, *l.StageConfig, *l.Message, *l.Retry, l.TargetQuantity, l.Link)
	if err != nil {
		fmt.Println("|||||||||||||")
		fmt.Println(err)
		fmt.Println("|||||||||||||")
		return fmt.Errorf("update link: %w", err)
	}

	return nil
}
