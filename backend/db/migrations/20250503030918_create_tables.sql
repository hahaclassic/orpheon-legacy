-- +goose Up
-- +goose StatementBegin
CREATE TABLE licenses (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    url TEXT NOT NULL
);

CREATE TABLE genres (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL UNIQUE CHECK (length(title) > 0) -- Название не может быть пустым
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (length(name) > 2), -- Имя должно быть хотя бы 3 символа
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата регистрации по умолчанию
    birth_date DATE CHECK (birth_date < registration_date), -- Дата рождения не может быть в будущем
    access_level INT NOT NULL CHECK (access_level IN (1, 2)) -- Ограничиваем возможные роли
);

CREATE TABLE credentials (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    login TEXT NOT NULL UNIQUE CHECK (length(login) > 2), -- Логин должен быть длиннее 3 символов
    password TEXT NOT NULL--- хешированный пароль
);

CREATE TABLE artists (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (length(name) > 0), -- Название артиста должно быть хотя бы 3 символа
    description TEXT,
    country TEXT 
);

CREATE TABLE albums (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL UNIQUE CHECK (length(title) > 0),
    label TEXT,
    license_id UUID REFERENCES licenses(id) ON DELETE SET NULL, -- Если лицензия удалена, оставляем NULL
    release_date DATE NOT NULL CHECK (release_date <= NOW()) -- Альбом не может выйти в будущем
);

CREATE TABLE tracks (
    id UUID PRIMARY KEY,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    duration INT NOT NULL CHECK (duration > 0), -- Длительность трека должна быть положительной
    name TEXT NOT NULL CHECK (length(name) > 0),
    explicit BOOLEAN NOT NULL DEFAULT FALSE, -- Явное указание значения по умолчанию
    license_id UUID REFERENCES licenses(id) ON DELETE SET NULL, -- Лицензия может быть удалена
    total_streams INT NOT NULL CHECK (total_streams >= 0),
    album_id UUID REFERENCES albums(id) ON DELETE CASCADE,
    track_number INT NOT NULL CHECK (track_number >=1), -- позиция в альбоме
    UNIQUE (name, album_id) -- уникальность имени трека в альбоме
);

CREATE TABLE playlists (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(name) > 0),
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    rating INT CHECK (rating >= 0)
);

CREATE TABLE track_segments (
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    index INT NOT NULL CHECK (index >= 0), -- Индексы не могут быть отрицательными
    total_streams INT NOT NULL DEFAULT 0 CHECK (total_streams >= 0), -- Количество прослушиваний не может быть отрицательным
    start_time INT NOT NULL CHECK (start_time >= 0),
    end_time INT NOT NULL CHECK (end_time > start_time), -- Время конца должно быть больше времени начала
    PRIMARY KEY (track_id, index)
);

CREATE TABLE album_genres (
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (album_id, genre_id)
);

CREATE TABLE artist_tracks (
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    PRIMARY KEY (artist_id, track_id)
);

CREATE TABLE artist_albums (
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    PRIMARY KEY (artist_id, album_id)
);

CREATE TABLE playlist_tracks (
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    position INT NOT NULL CHECK (position >= 1), -- Позиция трека в плейлисте начинается с 0
    PRIMARY KEY (playlist_id, track_id),
    UNIQUE (playlist_id, position)
);

CREATE TABLE favorite_playlists (
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (playlist_id, user_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorite_playlists;
DROP TABLE IF EXISTS playlist_tracks;
DROP TABLE IF EXISTS artist_tracks;
DROP TABLE IF EXISTS artist_albums;
DROP TABLE IF EXISTS track_segments;
DROP TABLE IF EXISTS album_genres;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS licenses;

-- +goose StatementEnd
