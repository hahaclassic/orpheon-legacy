-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_playlist_updated_at() 
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE playlists
        SET updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.playlist_id;
    ELSE
        UPDATE playlists
        SET updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.playlist_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_playlist_updated_at
AFTER INSERT OR UPDATE OR DELETE ON playlist_tracks
FOR EACH ROW
EXECUTE FUNCTION update_playlist_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_playlist_updated_at ON playlist_tracks;
DROP FUNCTION IF EXISTS update_playlist_updated_at();
-- +goose StatementEnd
