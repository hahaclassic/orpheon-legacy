package entity

import "github.com/google/uuid"

type ArtistMeta struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Country     string    `json:"country"`
}

type ArtistMetaBuilder struct {
	artistMeta ArtistMeta
}

func NewArtistMetaBuilder() *ArtistMetaBuilder {
	return &ArtistMetaBuilder{
		artistMeta: ArtistMeta{
			ID:          uuid.New(),
			Name:        "Default Artist",
			Description: "Default description",
			Country:     "RUS",
		},
	}
}

func (b *ArtistMetaBuilder) WithID(id uuid.UUID) *ArtistMetaBuilder {
	b.artistMeta.ID = id
	return b
}

func (b *ArtistMetaBuilder) WithName(name string) *ArtistMetaBuilder {
	b.artistMeta.Name = name
	return b
}

func (b *ArtistMetaBuilder) WithDescription(description string) *ArtistMetaBuilder {
	b.artistMeta.Description = description
	return b
}

func (b *ArtistMetaBuilder) WithCountry(country string) *ArtistMetaBuilder {
	b.artistMeta.Country = country
	return b
}

func (b *ArtistMetaBuilder) Build() *ArtistMeta {
	return &b.artistMeta
}
