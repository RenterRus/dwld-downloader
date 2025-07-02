package persistent

import (
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
	_, err := p.db.Exec("insert into links (link, target_quantity, work_status, path) values($1, $2, $3, $4);", link, maxQuality, "NEW", p.workDir)
	if err != nil {
		return nil, fmt.Errorf("insert new link: %w", err)
	}
	return nil, nil
}

func (p *persistentRepo) Update(data LinkModel) ([]LinkModel, error) {
	return nil, nil
}

func (p *persistentRepo) Delete(link string) ([]LinkModel, error) {
	return nil, nil
}

func (p *persistentRepo) DeleteHistory() ([]LinkModel, error) {
	return nil, nil
}
