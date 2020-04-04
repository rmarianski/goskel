// Package magic defines readers and writer to read and write payloads that are
// magic packed. The intended usage for this is to append some bytes to a
// binary that can later be read.
package magic

// Bytes is the magic bytes used in the command trailer.
var Bytes = "goskel"
