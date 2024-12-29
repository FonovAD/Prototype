PRAGMA foreign_keys = ON;
PRAGMA journal_mode=WAL;
PRAGMA cache_size=-1000;
CREATE TABLE IF NOT EXISTS users(
UID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
Token TEXT NOT NULL,
Role varchar(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS links(
UID INTEGER REFERENCES users(UID) ON DELETE CASCADE,
OriginLink TEXT NOT NULL,
ShortLink TEXT UNIQUE NOT NULL,
CreatedAt integer,
ExpirationTime integer NOT NULL,
Status varchar(10) NOT NULL,
ScheduledDeletionTime integer NOT NULL,
PRIMARY KEY (UID, OriginLink)
);

CREATE INDEX ShortLink ON links(ShortLink);

INSERT INTO users(Token, Role) VALUES("%s", "admin");