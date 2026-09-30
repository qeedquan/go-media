package mpd

import (
	"encoding/binary"
	"io"
)

type Header struct {
	NumChunks uint16
	NumActors uint16
	Unk       [6]uint16
}

type ChunkHeader struct {
	MapOff   [3]float32
	Unk1     [8]float32
	NumTiles uint16
	Unk2     uint16
	Index    uint16
	Unk3     [7]uint16
}

type ChunkObject struct {
	Unk [36]byte
}

type ChunkEventTile struct {
	X, Y  int8
	Index uint16
}

type Texture struct {
	U           uint8
	Stretch     uint8
	V           uint8
	Width       uint8
	SizeStretch uint8
	Height      uint8
	Index       uint8
	Mirror      uint8
}

type Tile struct {
	Textures [12]Texture
	Corners  [4]int8
	Corners2 [4]int8
	Corners3 [4]int8
	Unk2     [4]uint8
	Unk3     uint8
	X, Z     int8
	Unk4     uint8
	Unk5     uint8
	Unk6     uint8
	Unk7     uint8
	Mobility uint8
	GeoColor uint8
	GeoMark  uint8
	Pad      [6]uint8
}

type ChunkDef struct {
	Header     ChunkHeader
	Objects    [32]ChunkObject
	EventTiles [16]ChunkEventTile
	BaseTile   Tile
}

type Chunk struct {
	ChunkDef
	Tiles []Tile
}

type Actor struct {
	ID         uint16
	Level      uint16
	Unk2       uint8
	X, Z       int8
	Rotation   int8
	Unk4       uint8
	AI         uint8
	Unk5, Unk6 uint8
	Items      [4]uint16
	Appearance uint8
	Unk7       uint8
	GeoEffect  uint16
	Unk8       uint16
	Magic      [4]uint16
	Unk        [30]uint8
}

type Map struct {
	Header Header
	Chunks []Chunk
	Actors []Actor
}

func Decode(r io.Reader) (*Map, error) {
	var hdr Header
	err := binary.Read(r, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}

	var chunks []Chunk
	for range hdr.NumChunks {
		var (
			chunk Chunk
			tile  Tile
		)

		err = binary.Read(r, binary.LittleEndian, &chunk.ChunkDef)
		if err != nil {
			return nil, err
		}

		for range chunk.Header.NumTiles {
			err = binary.Read(r, binary.LittleEndian, &tile)
			if err != nil {
				return nil, err
			}
			chunk.Tiles = append(chunk.Tiles, tile)
		}
		chunks = append(chunks, chunk)
	}

	var actors []Actor
	for range hdr.NumActors {
		var actor Actor
		err = binary.Read(r, binary.LittleEndian, &actor)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	return &Map{
		Header: hdr,
		Chunks: chunks,
		Actors: actors,
	}, nil
}
