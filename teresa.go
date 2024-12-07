package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type BitBoard uint64
type Square string
type PieceType string

const (
	Rank1 BitBoard = 0x00000000000000FF
	Rank2 BitBoard = 0x000000000000FF00
	Rank3 BitBoard = 0x0000000000FF0000
	Rank4 BitBoard = 0x00000000FF000000
	Rank5 BitBoard = 0x000000FF00000000
	Rank6 BitBoard = 0x0000FF0000000000
	Rank7 BitBoard = 0x00FF000000000000
	Rank8 BitBoard = 0xFF00000000000000

	FileA BitBoard = 0x0101010101010101
	FileB BitBoard = 0x0202020202020202
	FileC BitBoard = 0x0404040404040404
	FileD BitBoard = 0x0808080808080808
	FileE BitBoard = 0x1010101010101010
	FileF BitBoard = 0x2020202020202020
	FileG BitBoard = 0x4040404040404040
	FileH BitBoard = 0x8080808080808080
)

type Move struct {
	From BitBoard
	To   BitBoard
}

func moveToAlgebraic(m Move) string {
	fromSquare := bitBoardToAlgebraic(m.From)
	toSquare := bitBoardToAlgebraic(m.To)
	return fromSquare + toSquare
}

func bitBoardToAlgebraic(bb BitBoard) string {
	if bb == 0 {
		return ""
	}
	index := 0
	for bb != 1 {
		bb >>= 1
		index++
	}
	file := index % 8
	rank := index / 8
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

func algebraicToBitBoard(square string) BitBoard {
	if len(square) != 2 {
		panic("Invalid square format")
	}

	file := square[0] - 'a'
	rank := square[1] - '1'

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		panic("Invalid square. Must be in the format a1-h8")
	}

	return BitBoard(1 << (rank*8 + file))
}

type Board struct {
	whiteToMove             bool
	WhitePawns              BitBoard
	WhiteKnights            BitBoard
	WhiteBishops            BitBoard
	WhiteRooks              BitBoard
	WhiteQueens             BitBoard
	WhiteKing               BitBoard
	BlackPawns              BitBoard
	BlackKnights            BitBoard
	BlackBishops            BitBoard
	BlackRooks              BitBoard
	BlackQueens             BitBoard
	BlackKing               BitBoard
	enpassantSquare         BitBoard
	halfMoveClock           int
	fullMoveNumber          int
	whiteCanCastleKingSide  bool
	whiteCanCastleQueenSide bool
	blackCanCastleKingSide  bool
	blackCanCastleQueenSide bool
}

func noPieceIsOnSquare(b Board, toSquare BitBoard) bool {
	return (b.WhitePawns|
		b.WhiteKnights|
		b.WhiteBishops|
		b.WhiteRooks|
		b.WhiteQueens|
		b.WhiteKing|
		b.BlackPawns|
		b.BlackKnights|
		b.BlackBishops|
		b.BlackRooks|
		b.BlackQueens|
		b.BlackKing)&toSquare == 0
}

func randomMove(moves []Move) Move {
	return moves[rand.Intn(len(moves))]
}

func generatePawnMoves(b Board, isWhite bool) []Move {
	var moves []Move
	var pawns BitBoard
	direction := 8
	var startRank BitBoard
	//var promotionRank BitBoard

	if isWhite {
		pawns = b.WhitePawns
		startRank = Rank2
		//promotionRank = 0xFF00000000000000
	} else {
		pawns = b.BlackPawns
		startRank = Rank7
		//promotionRank = 0x00000000000000FF
	}

	for pawns != 0 {
		square := pawns & -pawns
		pawns &= pawns - 1

		// Single move forward
		var toSquare BitBoard
		if isWhite {
			toSquare = square << direction
		} else {
			toSquare = square >> direction
		}

		if toSquare != 0 && noPieceIsOnSquare(b, toSquare) {
			moves = append(moves, Move{From: square, To: toSquare})
		}

		// Handle start rank moving 2 squares
		if square&startRank != 0 {
			var doubleMoveSquare BitBoard
			if isWhite {
				doubleMoveSquare = square << (2 * direction)
			} else {
				doubleMoveSquare = square >> (2 * direction)
			}

			if doubleMoveSquare != 0 && noPieceIsOnSquare(b, doubleMoveSquare) && noPieceIsOnSquare(b, toSquare) {
				moves = append(moves, Move{From: square, To: doubleMoveSquare})
			}
		}
	}

	return moves
}

func applyMove(b Board, m Move) Board {
	var nextBoard Board = b

	// Determine which piece is moving
	pieceMoved := BitBoard(0)
	for _, piece := range []BitBoard{
		b.WhitePawns, b.WhiteKnights, b.WhiteBishops, b.WhiteRooks, b.WhiteQueens, b.WhiteKing,
		b.BlackPawns, b.BlackKnights, b.BlackBishops, b.BlackRooks, b.BlackQueens, b.BlackKing,
	} {
		if m.From&piece != 0 {
			pieceMoved = piece
			break
		}
	}

	// Remove the piece from the source square
	nextBoard.WhitePawns &= ^m.From
	nextBoard.WhiteKnights &= ^m.From
	nextBoard.WhiteBishops &= ^m.From
	nextBoard.WhiteRooks &= ^m.From
	nextBoard.WhiteQueens &= ^m.From
	nextBoard.WhiteKing &= ^m.From
	nextBoard.BlackPawns &= ^m.From
	nextBoard.BlackKnights &= ^m.From
	nextBoard.BlackBishops &= ^m.From
	nextBoard.BlackRooks &= ^m.From
	nextBoard.BlackQueens &= ^m.From
	nextBoard.BlackKing &= ^m.From

	// Place the piece on the destination square
	switch pieceMoved {
	case b.WhitePawns:
		nextBoard.WhitePawns |= m.To
	case b.WhiteKnights:
		nextBoard.WhiteKnights |= m.To
	case b.WhiteBishops:
		nextBoard.WhiteBishops |= m.To
	case b.WhiteRooks:
		nextBoard.WhiteRooks |= m.To
	case b.WhiteQueens:
		nextBoard.WhiteQueens |= m.To
	case b.WhiteKing:
		nextBoard.WhiteKing |= m.To
	case b.BlackPawns:
		nextBoard.BlackPawns |= m.To
	case b.BlackKnights:
		nextBoard.BlackKnights |= m.To
	case b.BlackBishops:
		nextBoard.BlackBishops |= m.To
	case b.BlackRooks:
		nextBoard.BlackRooks |= m.To
	case b.BlackQueens:
		nextBoard.BlackQueens |= m.To
	case b.BlackKing:
		nextBoard.BlackKing |= m.To
	}

	// Toggle the active color
	nextBoard.whiteToMove = !b.whiteToMove

	if pieceMoved == b.WhitePawns || pieceMoved == b.BlackPawns {
		nextBoard.halfMoveClock = b.halfMoveClock + 1
	}

	if nextBoard.whiteToMove {
		nextBoard.fullMoveNumber = b.fullMoveNumber + 1
	}

	return nextBoard
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
	parts := strings.Split(fen, " ")
	if len(parts) < 4 {
		panic("Invalid FEN string")
	}
	piecePlacement := parts[0]
	activeColor := parts[1]
	castlingAvailability := parts[2]
	enPassantTargetSquare := parts[3]

	board.whiteToMove = activeColor == "w"

	if len(parts) > 4 {
		board.halfMoveClock = int(parts[4][0] - '0')
	}

	if len(parts) > 5 {
		board.fullMoveNumber = int(parts[5][0] - '0')
	}

	if castlingAvailability != "-" {
		for _, char := range castlingAvailability {
			switch char {
			case 'K':
				board.whiteCanCastleKingSide = true
			case 'Q':
				board.whiteCanCastleQueenSide = true
			case 'k':
				board.blackCanCastleKingSide = true
			case 'q':
				board.blackCanCastleQueenSide = true
			}
		}
	}

	if enPassantTargetSquare != "-" {
		file := enPassantTargetSquare[0] - 'a'
		rank := enPassantTargetSquare[1] - '1'
		board.enpassantSquare = BitBoard(1 << (rank*8 + file))
	}

	for _, char := range piecePlacement {
		if char == '/' {
			rankIndex++
			fileIndex = 0
		} else if char >= '1' && char <= '8' {
			fileIndex += int(char - '0')
		} else {
			square := 1 << (rankIndex*8 + fileIndex)
			switch char {
			case 'p':
				board.WhitePawns |= BitBoard(square)
			case 'n':
				board.WhiteKnights |= BitBoard(square)
			case 'b':
				board.WhiteBishops |= BitBoard(square)
			case 'r':
				board.WhiteRooks |= BitBoard(square)
			case 'q':
				board.WhiteQueens |= BitBoard(square)
			case 'k':
				board.WhiteKing |= BitBoard(square)
			case 'P':
				board.BlackPawns |= BitBoard(square)
			case 'N':
				board.BlackKnights |= BitBoard(square)
			case 'B':
				board.BlackBishops |= BitBoard(square)
			case 'R':
				board.BlackRooks |= BitBoard(square)
			case 'Q':
				board.BlackQueens |= BitBoard(square)
			case 'K':
				board.BlackKing |= BitBoard(square)
			}
			fileIndex++
		}
	}

	return board
}

func boardToFen(b Board) string {
	var fen strings.Builder

	for rank := 0; rank < 8; rank++ {
		emptyCount := 0
		for file := 0; file < 8; file++ {
			square := 1 << (rank*8 + file)
			piece := ""

			switch {
			case b.WhitePawns&BitBoard(square) != 0:
				piece = "p"
			case b.WhiteKnights&BitBoard(square) != 0:
				piece = "n"
			case b.WhiteBishops&BitBoard(square) != 0:
				piece = "b"
			case b.WhiteRooks&BitBoard(square) != 0:
				piece = "r"
			case b.WhiteQueens&BitBoard(square) != 0:
				piece = "q"
			case b.WhiteKing&BitBoard(square) != 0:
				piece = "k"
			case b.BlackPawns&BitBoard(square) != 0:
				piece = "P"
			case b.BlackKnights&BitBoard(square) != 0:
				piece = "N"
			case b.BlackBishops&BitBoard(square) != 0:
				piece = "B"
			case b.BlackRooks&BitBoard(square) != 0:
				piece = "R"
			case b.BlackQueens&BitBoard(square) != 0:
				piece = "Q"
			case b.BlackKing&BitBoard(square) != 0:
				piece = "K"
			default:
				emptyCount++
			}

			if piece != "" {
				if emptyCount > 0 {
					fen.WriteString(fmt.Sprintf("%d", emptyCount))
					emptyCount = 0
				}
				fen.WriteString(piece)
			}
		}
		if emptyCount > 0 {
			fen.WriteString(fmt.Sprintf("%d", emptyCount))
		}
		if rank < 7 {
			fen.WriteString("/")
		}
	}

	fen.WriteString(" ")
	if b.whiteToMove {
		fen.WriteString("w")
	} else {
		fen.WriteString("b")
	}

	fen.WriteString(" ")
	castling := ""
	if b.whiteCanCastleKingSide {
		castling += "K"
	}
	if b.whiteCanCastleQueenSide {
		castling += "Q"
	}
	if b.blackCanCastleKingSide {
		castling += "k"
	}
	if b.blackCanCastleQueenSide {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	fen.WriteString(castling)

	fen.WriteString(" ")
	if b.enpassantSquare != 0 {
		file := (b.enpassantSquare & 0xFF) % 8
		rank := (b.enpassantSquare & 0xFF) / 8
		fen.WriteString(fmt.Sprintf("%c%d", 'a'+file, rank+1))
	} else {
		fen.WriteString("-")
	}

	fen.WriteString(fmt.Sprintf(" %d %d", b.halfMoveClock, b.fullMoveNumber))

	return fen.String()
}

func main() {
	board := fenToBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	displayBoard(board)

	moves := generatePawnMoves(board, board.whiteToMove)

	for len(moves) > 0 {
		for _, num := range moves {
			fmt.Println(moveToAlgebraic(num))
		}

		board = applyMove(board, randomMove(moves))
		displayBoard(board)

		moves = generatePawnMoves(board, board.whiteToMove)
	}

	fmt.Println("No more legal moves")
}
