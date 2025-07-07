-- +goose Up
-- +goose StatementBegin
create table if not exists links (
link text unique,
filename text,
target_quantity integer,
work_status text,
message text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists links;
-- +goose StatementEnd
