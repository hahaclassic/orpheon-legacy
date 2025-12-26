# HTTP Routes

## /v1/api

### /auth
    * POST /auth/register
    * POST /auth/login
    * POST/auth/logout
    * POST /auth/refresh
    * POST /auth/password/update

### /albums

    * GET /albums/:id (meta)
    * POST /albums/:id 
    * DELETE /albums/:id 
    * GET /albums/:id/tracks

    /albums/:id/cover
        * GET /albums/:id/cover
        * POST /albums/:id/cover
        * DELETE /albums/:id/cover


### /artists

    * GET /artists/:id
    * POST /artists/:id
    * DELETE /artists/:id
    * GET /artists/:id/albums
    * GET /artists/:id/tracks
    * PUT /artists/:id/albums/:album_id
    * PUT /artists/:id/tracks/:track_id

    /artists/:id/avatar
        * GET /artists/:id/avatar
        * POST /artists/:id/avatar
        * DELETE /artists/:id/avatar

### /tracks

    * GET /tracks/:id
    * POST /tracks/:id
    * DELETE /tracks/:id
    * GET /tracks/:id/

    /tracks/:id/audio
        * GET /tracks/:id/audio (HTTP Range request)
        * POST /tracks/:id/audio
        * DELETE /tracks/:id/audio

    * GET /tracks/segments - получение статистики по сегментам
    
    * POST /tracks/stat/ - отправка статистики

### /genres

    * GET /genres/:id
    * POST /genres/:id
    * DELETE /genres/:id
    * GET /genres/:id/tracks

### /licenses

    * GET /licenses
    * GET /licenses/:id
    * POST /licenses/:id
    * DELETE /licenses/:id

### /search

    * GET /search?query=:query&limit=:limit&offset=:offset&type=:type&genre=:genre&country=:country

### /playlists

    * GET /playlists/:id
    * POST /playlists/:id
    * PUT /playlists/:id
    * DELETE /playlists/:id
    * PATCH /playlists/:id/privacy

    * GET /playlists/:id/tracks
    * POST /playlists/:id/tracks
    * DELETE /playlists/:id/tracks/:track_id
    * PATCH /playlists/:id/tracks/:track_id/position
    
    /playlists/:id/cover
        * GET /playlists/:id/cover
        * POST /playlists/:id/cover
        * DELETE /playlists/:id/cover

### /user (+ /me)

    * GET /me -> мета инфа
    * GET /me/playlists -> плейлисты
    * GET /me/favorites
    * POST /me/favorites/:playlist_id
    * DELETE /me/favorites/:playlist_id

    * GET /user/:id
    * POST /user/:id
    * DELETE /user/:id
    * GET /user/:id/playlists
