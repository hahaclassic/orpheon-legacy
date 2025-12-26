-- +goose Up
-- +goose StatementBegin
-- CREATE OR REPLACE PROCEDURE change_track_position(
--     p_playlist_id UUID,
--     p_track_id UUID,
--     p_new_position INT
-- )
-- LANGUAGE plpgsql
-- AS $$
-- DECLARE
--     current_position INT;
-- BEGIN
--     SELECT position INTO current_position
--     FROM playlist_tracks
--     WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

--     IF current_position IS NULL THEN
--         RAISE EXCEPTION 'Track not found in playlist';
--     END IF;

--     IF current_position = p_new_position THEN
--         RETURN;
--     END IF;

--     IF current_position > p_new_position THEN
--         UPDATE playlist_tracks
--         SET position = position + 1000001
--         WHERE playlist_id = p_playlist_id
--           AND position >= p_new_position AND position < current_position;

--         UPDATE playlist_tracks
--         SET position = p_new_position
--         WHERE playlist_id = p_playlist_id
--             and track_id = p_track_id;

--         UPDATE playlist_tracks
--         SET position = position - 1000000
--         WHERE playlist_id = p_playlist_id
--           AND position > 1000000;
--     ELSE
--         UPDATE playlist_tracks
--         SET position = position + 1000000
--         WHERE playlist_id = p_playlist_id
--           AND position > current_position AND position <= p_new_position;

--         UPDATE playlist_tracks
--         SET position = p_new_position
--         WHERE playlist_id = p_playlist_id
--             and track_id = p_track_id;

--         UPDATE playlist_tracks
--         SET position = position - 1000001
--         WHERE playlist_id = p_playlist_id
--           AND position > 1000000;
--     END IF;
-- END;
-- $$;


CREATE OR REPLACE PROCEDURE change_track_position(
    p_playlist_id UUID,
    p_track_id UUID,
    p_new_position INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    LARGE_OFFSET CONSTANT INT := 1000000;
    current_position INT;
    left_bound INT; -- граница включается!
    right_bound INT; -- граница включается!
    offset_add INT;
    offset_del INT;
BEGIN
    SELECT position INTO current_position
    FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Track not found in playlist';
    END IF;

    IF current_position = p_new_position THEN
        RETURN;
    END IF;

    IF current_position > p_new_position THEN
        left_bound := p_new_position;
        right_bound := current_position - 1;
        offset_add := LARGE_OFFSET + 1;
        offset_del := LARGE_OFFSET;
    ELSE
        left_bound := current_position + 1;
        right_bound := p_new_position;
        offset_add := LARGE_OFFSET;
        offset_del := LARGE_OFFSET + 1;
    END IF;

    UPDATE playlist_tracks
    SET position = position + offset_add
    WHERE playlist_id = p_playlist_id
        AND position BETWEEN left_bound AND right_bound;
    
    UPDATE playlist_tracks
    SET position = p_new_position
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id;
    
    UPDATE playlist_tracks
    SET position = position - offset_del
    WHERE playlist_id = p_playlist_id
        AND position > LARGE_OFFSET;
END;
$$;


-- CREATE OR REPLACE PROCEDURE delete_track_from_playlist(p_playlist_id UUID, p_track_id UUID)
-- LANGUAGE plpgsql
-- AS $$
-- DECLARE
--     deleted_position INT;
-- BEGIN
--     SELECT position INTO deleted_position
--     FROM playlist_tracks
--     WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

--     IF NOT FOUND THEN
--         RAISE EXCEPTION 'Track with id % not found in playlist %', p_track_id, p_playlist_id;
--     END IF;

--     DELETE FROM playlist_tracks
--     WHERE playlist_id = p_playlist_id AND track_id = p_track_id;

--     UPDATE playlist_tracks
--     SET position = position + 1000000
--     WHERE playlist_id = p_playlist_id AND position > deleted_position;

--     UPDATE playlist_tracks
--     SET position = position - 1000001
--     WHERE playlist_id = p_playlist_id AND position > 1000000;
-- END;
-- $$;

CREATE OR REPLACE PROCEDURE delete_track_from_playlist(
    p_playlist_id UUID,
    p_track_id UUID
)
LANGUAGE plpgsql
AS $$
DECLARE
    deleted_position INT;
    LARGE_OFFSET CONSTANT INT := 1000000;  -- Константа для временного сдвига позиций
BEGIN
     -- Удаляем трек и сразу получаем его позицию
    DELETE FROM playlist_tracks
    WHERE playlist_id = p_playlist_id AND track_id = p_track_id
    RETURNING position INTO deleted_position;
    
    -- Проверяем, был ли удален какой-либо трек
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Track with id % not found in playlist %', p_track_id, p_playlist_id;
    END IF;

    -- Временно сдвигаем позиции всех последующих треков вверх
    UPDATE playlist_tracks
    SET position = position + LARGE_OFFSET
    WHERE playlist_id = p_playlist_id AND position > deleted_position;

    -- Возвращаем сдвинутые треки на правильные позиции (с уменьшением на 1)
    UPDATE playlist_tracks
    SET position = position - (LARGE_OFFSET + 1)
    WHERE playlist_id = p_playlist_id AND position > LARGE_OFFSET;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP PROCEDURE change_track_position(UUID, UUID, INT);
DROP PROCEDURE delete_track_from_playlist(UUID, UUID);
-- +goose StatementEnd
