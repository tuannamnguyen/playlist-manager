package spotifyconverter

func chunkBy[T any](items []T, chunkSize int) [][]T {
	var chunks = make([][]T, 0, (len(items)/chunkSize)+1)
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
