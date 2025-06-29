-- +goose Up
-- +goose StatementBegin
create table if not exists links (
link text unique,
filename text,
path text,
work_status text,
stage_config text,
retry int,
message text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists links;
-- +goose StatementEnd
