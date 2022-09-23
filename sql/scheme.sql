CREATE TABLE  IF NOT EXISTS files (
    id INTEGER PRIMARY KEY,
    name TEXT,
    size INTEGER,
    content_type TEXT
)
drop table files;
delete from files where true;