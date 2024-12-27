CREATE TABLE users(
UID SERIAL PRIMARY KEY,
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

INSERT INTO users (token, role) VALUES ('test', 'admin');