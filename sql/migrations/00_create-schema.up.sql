CREATE TYPE user_typ AS ENUM ('User', 'Sys');

CREATE TYPE user_role AS ENUM ('Admin', 'Editor', 'User');

-- M = Male
-- F = Female
-- U = Unknown
CREATE TYPE gender_typ AS ENUM ('M', 'F', 'U');

CREATE TABLE "user" (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY (START WITH 1000) PRIMARY KEY,
    -- Personal information
    name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL UNIQUE,
    -- Auth
    pwd VARCHAR(256),
    -- Flags
    typ user_typ NOT NULL DEFAULT 'User',
    role user_role NOT NULL DEFAULT 'User',
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    -- Timestamps
    cid bigint NOT NULL,
    ctime timestamp with time zone NOT NULL,
    mid bigint NOT NULL,
    mtime timestamp with time zone NOT NULL
);

CREATE TABLE "account" (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY (START WITH 1000) PRIMARY KEY,
    -- Personal information
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(64) NOT NULL,
    email VARCHAR(64),
    gender gender_typ NOT NULL DEFAULT 'U',
    address VARCHAR(256),
    phone VARCHAR(16) NOT NULL UNIQUE,
    secondary_phone VARCHAR(16),
    -- Flags
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    -- Foreign keys
    social_media_id BIGINT REFERENCES social_media(id) ON DELETE SET NOT NULL,
    -- Timestamps
    cid bigint NOT NULL,
    ctime timestamp with time zone NOT NULL,
    mid bigint NOT NULL,
    mtime timestamp with time zone NOT NULL
);

CREATE TABLE "social_media" (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY (START WITH 1000) PRIMARY KEY
);

CREATE TABLE "social_media_data" (
    id BIGINT GENERATED BY DEFAULT AS IDENTITY (START WITH 1000) PRIMARY KEY,
    -- Foreign keys
    social_media_id BIGINT NOT NULL REFERENCES social_media(id) ON DELETE CASCADE,
    -- Social media information
    platform VARCHAR(25) NOT NULL,
    -- will store in lowercase
    url VARCHAR(256) NOT NULL,
    -- Timestamps
    cid BIGINT NOT NULL,
    ctime TIMESTAMP WITH TIME ZONE NOT NULL,
    mid BIGINT NOT NULL,
    mtime TIMESTAMP WITH TIME ZONE NOT NULL
)