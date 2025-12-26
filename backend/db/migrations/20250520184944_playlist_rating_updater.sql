-- +goose Up
-- +goose StatementBegin
-- CREATE OR REPLACE FUNCTION update_playlist_rating()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     IF (TG_OP = 'INSERT') THEN
--         UPDATE playlists
--         SET rating = rating + 1
--         WHERE id = NEW.playlist_id;
--     ELSIF (TG_OP = 'DELETE') THEN
--         UPDATE playlists
--         SET rating = rating - 1
--         WHERE id = OLD.playlist_id;
--     END IF;

--     RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_playlist_rating()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE playlists
        SET rating = rating - 1
        WHERE id = OLD.playlist_id;
    ELSE
        UPDATE playlists
        SET rating = rating + 1
        WHERE id = NEW.playlist_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_playlist_rating
AFTER INSERT OR DELETE ON favorite_playlists
FOR EACH ROW
EXECUTE FUNCTION update_playlist_rating();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_playlist_rating ON favorite_playlists;
DROP FUNCTION IF EXISTS update_playlist_rating();
-- +goose StatementEnd
