package output

import (
	"fmt"
	"strings"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	tableoutput "github.com/hahaclassic/orpheon/backend/pkg/table"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintAlbum(album *entity.AlbumMeta) {
	fmt.Println("--------------------------------")
	fmt.Printf("Album ID: %s\n", album.ID)
	fmt.Printf("Title: %s\n", album.Title)
	fmt.Printf("Label: %s\n", album.Label)
	fmt.Printf("License ID: %s\n", album.LicenseID)
	fmt.Printf("Release Date: %s\n", album.ReleaseDate.Format("2006-01-02"))
	fmt.Println("--------------------------------")
}

func PrintAlbums(albums []*entity.AlbumMeta) {
	var tableData [][]any
	for _, album := range albums {
		tableData = append(tableData, []any{album.ID, album.Title, album.Label,
			album.LicenseID, album.ReleaseDate.Format("2006-01-02")})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Title", "Label", "License ID", "Release Date"}, tableData)
}

func PrintArtist(artist *entity.ArtistMeta) {
	fmt.Println("--------------------------------")
	fmt.Println("Artist ID:", artist.ID)
	fmt.Println("Name:", artist.Name)
	fmt.Println("Description:", artist.Description)
	fmt.Println("Country:", artist.Country)
	fmt.Println("--------------------------------")
}

func PrintArtists(artists []*entity.ArtistMeta) {
	var tableData [][]any
	for _, artist := range artists {
		tableData = append(tableData, []any{artist.ID, artist.Name, artist.Country,
			artist.Description[:min(len(artist.Description), 70)]})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Name", "Country", "Description"}, tableData)
}

func PrintTrack(track *entity.TrackMeta) {
	fmt.Println("--------------------------------")
	fmt.Println("Track ID:", track.ID)
	fmt.Println("Name:", track.Name)
	fmt.Println("Duration:", track.Duration)
	fmt.Println("Explicit:", track.Explicit)
	fmt.Println("License ID:", track.LicenseID)
	fmt.Println("Genre ID:", track.GenreID)
	fmt.Println("Album ID:", track.AlbumID)
	fmt.Println("Track Number:", track.TrackNumber)
	fmt.Println("Total Streams:", track.TotalStreams)
	fmt.Println("--------------------------------")
}

func PrintTracks(tracks []*entity.TrackMeta) {
	var tableData [][]any
	for _, track := range tracks {
		tableData = append(tableData, []any{track.ID, track.Name, track.Duration,
			track.Explicit, track.LicenseID, track.GenreID, track.TrackNumber, track.TotalStreams})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Name", "Duration", "Explicit", "License ID", "Genre ID", "Track Number", "Total Streams"}, tableData)
}

func PrintPlaylist(playlist *entity.PlaylistMeta) {
	fmt.Println("--------------------------------")
	fmt.Println("Playlist ID:", playlist.ID)
	fmt.Println("Name:", playlist.Name)
	fmt.Println("Description:", playlist.Description)
	fmt.Println("IsPrivate:", playlist.IsPrivate)
	fmt.Println("Owner ID:", playlist.OwnerID)
	fmt.Println("Rating:", playlist.Rating)
	fmt.Println("Created At:", playlist.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("Updated At:", playlist.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("--------------------------------")
}

func PrintPlaylists(playlists []*entity.PlaylistMeta) {
	var tableData [][]any
	for _, playlist := range playlists {
		tableData = append(tableData, []any{playlist.ID, playlist.Name, playlist.IsPrivate, playlist.OwnerID,
			playlist.Rating, playlist.CreatedAt.Format("2006-01-02 15:04:05"),
			playlist.UpdatedAt.Format("2006-01-02 15:04:05"),
			playlist.Description[:min(len(playlist.Description), 20)]})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Name", "Is Private", "Owner ID", "Rating", "Created At", "Updated At", "Description"}, tableData)
}

func PrintLicense(license *entity.License) {
	fmt.Println("--------------------------------")
	fmt.Println("License ID:", license.ID)
	fmt.Println("Title:", license.Title)
	fmt.Println("Description:", license.Description)
	fmt.Println("--------------------------------")
}

func PrintLicenses(licenses []*entity.License) {
	var tableData [][]any
	for _, license := range licenses {
		tableData = append(tableData, []any{license.ID, license.Title,
			license.Description[:min(len(license.Description), 70)]})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Title", "Description"}, tableData)
}

func PrintGenre(genre *entity.Genre) {
	fmt.Println("--------------------------------")
	fmt.Println("Genre ID:", genre.ID)
	fmt.Println("Title:", genre.Title)
	fmt.Println("--------------------------------")
}

func PrintGenres(genres []*entity.Genre) {
	var tableData [][]any
	for _, genre := range genres {
		tableData = append(tableData, []any{genre.ID, genre.Title})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Title"}, tableData)
}

func PrintUser(user *entity.UserInfo) {
	fmt.Println("--------------------------------")
	fmt.Println("User ID:", user.ID)
	fmt.Println("Name:", user.Name)
	fmt.Println("Registration Date:", user.RegistrationDate.Format("2006-01-02 15:04:05"))
	fmt.Println("Birth Date:", user.BirthDate.Format("2006-01-02"))
	fmt.Println("Access Level:", user.AccessLvl)
	fmt.Println("--------------------------------")
}

func PrintUsers(users []*entity.UserInfo) {
	var tableData [][]any
	for _, user := range users {
		tableData = append(tableData, []any{user.ID, user.Name, user.RegistrationDate.Format("2006-01-02 15:04:05"),
			user.BirthDate.Format("2006-01-02 15:04:05"), user.AccessLvl})
	}

	tableoutput.PrintTable(table.StyleColoredDark,
		[]string{"ID", "Name", "Registration Date", "Birth Date", "Access Level"}, tableData)
}

func PrintStatsGraph(segments []*entity.Segment) {
	const (
		graphWidth  = 60 // ширина графика в символах
		graphHeight = 10 // высота графика в символах
	)

	// Находим максимальное количество прослушиваний для масштабирования
	maxStreams := uint64(0)
	for _, seg := range segments {
		if seg.TotalStreams > maxStreams {
			maxStreams = seg.TotalStreams
		}
	}

	// Создаем матрицу для графика
	graph := make([][]rune, graphHeight)
	for i := range graph {
		graph[i] = make([]rune, graphWidth)
		for j := range graph[i] {
			graph[i][j] = ' '
		}
	}

	// Заполняем график
	segmentWidth := graphWidth / len(segments)
	for i, seg := range segments {
		height := int(float64(seg.TotalStreams) / float64(maxStreams) * float64(graphHeight))
		if height > 0 {
			startX := i * segmentWidth
			endX := (i + 1) * segmentWidth
			endX = min(endX, graphWidth)

			for x := startX; x < endX; x++ {
				for y := 0; y < height; y++ {
					graph[graphHeight-1-y][x] = '█'
				}
			}
		}
	}

	// Выводим график
	fmt.Println("\nStreams per segment:")
	fmt.Println("┌" + strings.Repeat("─", graphWidth) + "┐")
	for _, row := range graph {
		fmt.Print("│")
		for _, cell := range row {
			fmt.Print(string(cell))
		}
		fmt.Println("│")
	}
	fmt.Println("└" + strings.Repeat("─", graphWidth) + "┘")
	fmt.Printf("Max streams: %d\n", maxStreams)
}
