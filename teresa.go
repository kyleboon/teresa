package main

import "fmt"

type BitBoard uint64

type Board struct {
	WhitePawns   BitBoard
	WhiteKnights BitBoard
	WhiteBishops BitBoard
	WhiteRooks   BitBoard
	WhiteQueens  BitBoard
	WhiteKing    BitBoard
	BlackPawns   BitBoard
	BlackKnights BitBoard
	BlackBishops BitBoard
	BlackRooks   BitBoard
	BlackQueens  BitBoard
	BlackKing    BitBoard
}

func displayBoard(b Board) {
	pieces := map[BitBoard]string{
		b.WhitePawns:   "♙",
		b.WhiteKnights: "♘",
		b.WhiteBishops: "♗",
		b.WhiteRooks:   "♖",
		b.WhiteQueens:  "♕",
		b.WhiteKing:    "♔",
		b.BlackPawns:   "♟",
		b.BlackKnights: "♞",
		b.BlackBishops: "♝",
		b.BlackRooks:   "♜",
		b.BlackQueens:  "♛",
		b.BlackKing:    "♚",
	}

	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			square := 1 << (rank*8 + file)
			piece := "."

			for bitboard, symbol := range pieces {
				if bitboard&BitBoard(square) != 0 {
					piece = symbol
					break
				}
			}

			fmt.Print(piece, " ")
		}
		fmt.Println()
	}
}

func fenToBoard(fen string) Board {
	var board Board
	rankIndex := 0
	fileIndex := 0

	for _, char := range fen {
		if char == '/' {
			rankIndex++
			fileIndex = 0
		} else if char >= '1' && char <= '8' {
			fileIndex += int(char - '0')
		} else {
			square := 1 << (rankIndex*8 + fileIndex)
			switch char {
			case 'P':
				board.WhitePawns |= BitBoard(square)
			case 'N':
				board.WhiteKnights |= BitBoard(square)
			case 'B':
				board.WhiteBishops |= BitBoard(square)
			case 'R':
				board.WhiteRooks |= BitBoard(square)
			case 'Q':
				board.WhiteQueens |= BitBoard(square)
			case 'K':
				board.WhiteKing |= BitBoard(square)
			case 'p':
				board.BlackPawns |= BitBoard(square)
			case 'n':
				board.BlackKnights |= BitBoard(square)
			case 'b':
				board.BlackBishops |= BitBoard(square)
			case 'r':
				board.BlackRooks |= BitBoard(square)
			case 'q':
				board.BlackQueens |= BitBoard(square)
			case 'k':
				board.BlackKing |= BitBoard(square)
			}
			fileIndex++
		}
	}

	return board
}

func main() {
	// Create a new board from a FEN string
	// https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation
	// The starting position is:
	// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR

	board := fenToBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	displayBoard(board)
}
