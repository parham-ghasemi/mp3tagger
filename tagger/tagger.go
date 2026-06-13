package tagger

import (
	"github.com/bogem/id3v2/v2"
)

func Execute(filepath string, trackNum string, title string, artist string, album string) error {
	tag, err := id3v2.Open(filepath, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	tag.SetTitle(title)
	tag.SetArtist(artist)
	tag.SetAlbum(album)
	tag.AddTextFrame(tag.CommonID("Track number/Position in set"), id3v2.EncodingUTF8, trackNum)

	if err := tag.Save(); err != nil {
		return err
	}

	return nil
}
