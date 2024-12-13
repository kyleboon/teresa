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
	Empty BitBoard = 0x0000000000000000
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

func bitBoardToAlgebraic(position BitBoard) string {
	if position == 0 {
		return ""
	}
	index := 0
	for position != 1 {
		position >>= 1
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

func piecesForCurrentPlayer(currentposition Board) []BitBoard {
	if currentposition.whiteToMove {
		return []BitBoard{
			currentposition.WhitePawns, currentposition.WhiteKnights, currentposition.WhiteBishops, currentposition.WhiteRooks, currentposition.WhiteQueens, currentposition.WhiteKing,
		}
	}
	return []BitBoard{
		currentposition.BlackPawns, currentposition.BlackKnights, currentposition.BlackBishops, currentposition.BlackRooks, currentposition.BlackQueens, currentposition.BlackKing,
	}
}

func piecesForOpposingPlayer(currentPosition Board) []BitBoard {
	if currentPosition.whiteToMove {
		return []BitBoard{
			currentPosition.BlackPawns, currentPosition.BlackKnights, currentPosition.BlackBishops, currentPosition.BlackRooks, currentPosition.BlackQueens, currentPosition.BlackKing,
		}
	}
	return []BitBoard{
		currentPosition.WhitePawns, currentPosition.WhiteKnights, currentPosition.WhiteBishops, currentPosition.WhiteRooks, currentPosition.WhiteQueens, currentPosition.WhiteKing,
	}
}

func noPieceIsOnSquare(currentPosition Board, toSquare BitBoard) bool {
	return (currentPosition.WhitePawns|
		currentPosition.WhiteKnights|
		currentPosition.WhiteBishops|
		currentPosition.WhiteRooks|
		currentPosition.WhiteQueens|
		currentPosition.WhiteKing|
		currentPosition.BlackPawns|
		currentPosition.BlackKnights|
		currentPosition.BlackBishops|
		currentPosition.BlackRooks|
		currentPosition.BlackQueens|
		currentPosition.BlackKing)&toSquare == 0
}

func opponentPieceIsOnSquare(currentPosition Board, toSquare BitBoard) bool {
	opposingPieces := piecesForOpposingPlayer(currentPosition)
	for _, piece := range opposingPieces {
		if piece&toSquare != 0 {
			return true
		}
	}

	return false
}

func randomMove(moves []Move) Move {
	return moves[rand.Intn(len(moves))]
}

func generatePawnMoves(currentPosition Board) []Move {
	var moves []Move
	var pawns BitBoard
	direction := 8
	var startRank BitBoard

	if currentPosition.whiteToMove {
		pawns = currentPosition.WhitePawns
		startRank = Rank2
	} else {
		pawns = currentPosition.BlackPawns
		startRank = Rank7
	}

	for pawns != 0 {
		square := pawns & -pawns
		pawns &= pawns - 1

		// Single move forward
		var toSquare BitBoard
		if currentPosition.whiteToMove {
			toSquare = square << direction
		} else {
			toSquare = square >> direction
		}

		if toSquare != 0 && noPieceIsOnSquare(currentPosition, toSquare) {
			moves = append(moves, Move{From: square, To: toSquare})
		}

		// Handle start rank moving 2 squares
		if square&startRank != Empty {
			var doubleMoveSquare BitBoard
			if currentPosition.whiteToMove {
				doubleMoveSquare = square << (2 * direction)
			} else {
				doubleMoveSquare = square >> (2 * direction)
			}

			if doubleMoveSquare != Empty && noPieceIsOnSquare(currentPosition, doubleMoveSquare) && noPieceIsOnSquare(currentPosition, toSquare) {
				moves = append(moves, Move{From: square, To: doubleMoveSquare})
			}
		}

		// Handle captures
		// Capture to the left
		if currentPosition.whiteToMove {
			toSquare = (square << (direction - 1)) & ^FileA
		} else {
			toSquare = (square >> (direction + 1)) & ^FileH
		}
		if toSquare != 0 && opponentPieceIsOnSquare(currentPosition, toSquare) {
			moves = append(moves, Move{From: square, To: toSquare})
		}

		// Capture to the right
		if currentPosition.whiteToMove {
			toSquare = (square << (direction + 1)) & ^FileH
		} else {
			toSquare = (square >> (direction - 1)) & ^FileA
		}
		if toSquare != Empty && opponentPieceIsOnSquare(currentPosition, toSquare) {
			moves = append(moves, Move{From: square, To: toSquare})
		}
		// handle en passant
	}

	return moves
}

func generateKnightMoves(currentPosition Board) []Move {
	var moves []Move
	var knights BitBoard
	var directions = []int{15, 17, 10, 6, -15, -17, -10, -6}

	if currentPosition.whiteToMove {
		knights = currentPosition.WhiteKnights
	} else {
		knights = currentPosition.BlackKnights
	}

	for knights != 0 {
		square := knights & -knights
		knights &= knights - 1

		for _, direction := range directions {
			var toSquare BitBoard

			// Ensure the knight doesn't wrap around the board
			if (direction == 15 || direction == -17) && (square&FileA != 0) {
				continue
			}
			if (direction == 17 || direction == -15) && (square&FileH != 0) {
				continue
			}
			if (direction == 10 || direction == -6) && (square&(FileA|FileB) != 0) {
				continue
			}
			if (direction == 6 || direction == -10) && (square&(FileG|FileH) != 0) {
				continue
			}

			if direction > 0 {
				toSquare = square << direction
			} else {
				toSquare = square >> -direction
			}

			if toSquare != 0 && noPieceIsOnSquare(currentPosition, toSquare) {
				moves = append(moves, Move{From: square, To: toSquare})
			} else if toSquare != 0 && opponentPieceIsOnSquare(currentPosition, toSquare) {
				moves = append(moves, Move{From: square, To: toSquare})
			}
		}
	}

	return moves
}

func generateMoves(currentPosition Board) []Move {
	return append(generatePawnMoves(currentPosition), generateKnightMoves(currentPosition)...)
}

func bitBoardsInterect(position1 BitBoard, position2 BitBoard) bool {
	return position1&position2 != Empty
}

func applyMove(currentPosition Board, currentMove Move) Board {
	var resultingPosition Board = currentPosition

	// Determine which piece is moving
	pieceMoved := Empty
	for _, piece := range piecesForCurrentPlayer(currentPosition) {
		if bitBoardsInterect(currentMove.From, piece) {
			pieceMoved = piece
			break
		}
	}

	// remove captured pice
	pieceRemoved := Empty
	isCapture := false
	for _, piece := range piecesForOpposingPlayer(resultingPosition) {
		if bitBoardsInterect(currentMove.To, piece) {
			fmt.Println("Removing piece!")
			pieceRemoved = piece
			isCapture = true
			break
		}
	}

	// Remove the piece from the source square
	resultingPosition.WhitePawns &= ^currentMove.From
	resultingPosition.WhiteKnights &= ^currentMove.From
	resultingPosition.WhiteBishops &= ^currentMove.From
	resultingPosition.WhiteRooks &= ^currentMove.From
	resultingPosition.WhiteQueens &= ^currentMove.From
	resultingPosition.WhiteKing &= ^currentMove.From
	resultingPosition.BlackPawns &= ^currentMove.From
	resultingPosition.BlackKnights &= ^currentMove.From
	resultingPosition.BlackBishops &= ^currentMove.From
	resultingPosition.BlackRooks &= ^currentMove.From
	resultingPosition.BlackQueens &= ^currentMove.From
	resultingPosition.BlackKing &= ^currentMove.From

	// Place the piece on the destination square
	if pieceRemoved != Empty {
		switch pieceRemoved {
		case resultingPosition.WhitePawns:
			fmt.Println("Removing white pawn")
			resultingPosition.WhitePawns &= ^currentMove.To
		case resultingPosition.WhiteKnights:
			fmt.Println("Removing white knight")
			resultingPosition.WhiteKnights &= ^currentMove.To
		case resultingPosition.WhiteBishops:
			fmt.Println("Removing white bishop")
			resultingPosition.WhiteBishops &= ^currentMove.To
		case resultingPosition.WhiteRooks:
			fmt.Println("Removing white rook")
			resultingPosition.WhiteRooks &= ^currentMove.To
		case resultingPosition.WhiteQueens:
			fmt.Println("Removing white queen")
			resultingPosition.WhiteQueens &= ^currentMove.To
		case resultingPosition.WhiteKing:
			fmt.Println("Removing white king")
			resultingPosition.WhiteKing &= ^currentMove.To
		case resultingPosition.BlackPawns:
			fmt.Println("Removing black pawn")
			resultingPosition.BlackPawns &= ^currentMove.To
		case resultingPosition.BlackKnights:
			fmt.Println("Removing black knight")
			resultingPosition.BlackKnights &= ^currentMove.To
		case resultingPosition.BlackBishops:
			fmt.Println("Removing black bishop")
			resultingPosition.BlackBishops &= ^currentMove.To
		case resultingPosition.BlackRooks:
			fmt.Println("Removing black rook")
			resultingPosition.BlackRooks &= ^currentMove.To
		case resultingPosition.BlackQueens:
			fmt.Println("Removing black queen")
			resultingPosition.BlackQueens &= ^currentMove.To
		case resultingPosition.BlackKing:
			fmt.Println("Removing black king")
			resultingPosition.BlackKing &= ^currentMove.To
		}
	}

	// Place the piece on the destination square
	switch pieceMoved {
	case currentPosition.WhitePawns:
		fmt.Println("Placing white pawn")
		resultingPosition.WhitePawns |= currentMove.To
	case currentPosition.WhiteKnights:
		fmt.Println("Placing white knight")
		resultingPosition.WhiteKnights |= currentMove.To
	case currentPosition.WhiteBishops:
		fmt.Println("Placing white bishop")
		resultingPosition.WhiteBishops |= currentMove.To
	case currentPosition.WhiteRooks:
		fmt.Println("Placing white rook")
		resultingPosition.WhiteRooks |= currentMove.To
	case currentPosition.WhiteQueens:
		fmt.Println("Placing white queen")
		resultingPosition.WhiteQueens |= currentMove.To
	case currentPosition.WhiteKing:
		fmt.Println("Placing white king")
		resultingPosition.WhiteKing |= currentMove.To
	case currentPosition.BlackPawns:
		fmt.Println("Placing black pawn")
		resultingPosition.BlackPawns |= currentMove.To
	case currentPosition.BlackKnights:
		fmt.Println("Placing black knight")
		resultingPosition.BlackKnights |= currentMove.To
	case currentPosition.BlackBishops:
		fmt.Println("Placing black bishop")
		resultingPosition.BlackBishops |= currentMove.To
	case currentPosition.BlackRooks:
		fmt.Println("Placing black rook")
		resultingPosition.BlackRooks |= currentMove.To
	case currentPosition.BlackQueens:
		fmt.Println("Placing black queen")
		resultingPosition.BlackQueens |= currentMove.To
	case currentPosition.BlackKing:
		fmt.Println("Placing black king")
		resultingPosition.BlackKing |= currentMove.To
	default:
		fmt.Println("No piece to place.")
	}

	// Toggle the active color
	resultingPosition.whiteToMove = !currentPosition.whiteToMove

	if isCapture || (pieceMoved == currentPosition.WhitePawns || pieceMoved == currentPosition.BlackPawns) {
		resultingPosition.halfMoveClock = 0
	} else {
		resultingPosition.halfMoveClock = currentPosition.halfMoveClock + 1
	}

	if resultingPosition.whiteToMove {
		resultingPosition.fullMoveNumber = currentPosition.fullMoveNumber + 1
	}

	return resultingPosition
}

func displayBitBoard(bb BitBoard) {
	fmt.Println("  a b c d e f g h")
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("%d ", rank+1)
		for file := 0; file < 8; file++ {
			square := BitBoard(1 << (rank*8 + file))
			if bb&square != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print("- ")
			}
		}
		fmt.Printf("%d\n", rank+1)
	}
	fmt.Println("  a b c d e f g h")
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

	moves := generateMoves(board)

	// fmt.Println(len(moves))

	// for _, move := range moves {
	// 	fmt.Println(moveToAlgebraic(move))
	// }

	for len(moves) > 0 {
		move := randomMove(moves)
		fmt.Println(moveToAlgebraic(move))
		board = applyMove(board, move)
		displayBoard(board)

		moves = generateMoves(board)
	}

	if board.whiteToMove {
		fmt.Println("White's turn but there are no legal moves.")
	} else {
		fmt.Println("Black's turn but there are no legal moves.")
	}
}
