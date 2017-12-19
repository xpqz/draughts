package main

// Score ...
type Score struct {
	value int
	move  *Move
}

type boardEval struct {
	Player         int
	PieceDiff      int
	KingDiff       int
	MovesCountDiff int
	CapturableDiff int
	KingMakerDiff  int
}

const (
	pieceWeight      = 15
	kingWeight       = 20
	movesCountWeight = 2
	capturableWeight = 10
	kingMakerWeight  = 6
)

// minimax is a naive implementation of the classic optimisation algorithm
func minimax(b *Board, player, depth, maxDepth int) *Score {
	allMoves := b.AllMoves(player)
	if len(allMoves) == 0 || depth == maxDepth {
		return &Score{value: HeuristicValue(b, player)}
	}

	scores := []*Score{}
	for _, move := range allMoves {
		newState := b.Apply(move)
		score := minimax(newState, Opposition(player), depth+1, maxDepth)
		score.move = move
		scores = append(scores, score)
	}

	if player == 1 {
		max := scores[0]
		for _, sc := range scores {
			if sc.value > max.value {
				max = sc
			}
		}
		return max
	}

	min := scores[0]
	for _, sc := range scores {
		if sc.value < min.value {
			min = sc
		}
	}

	return min
}

// KingsCaptures ...
func KingsCaptures(movesList []*Move) (int, int) {
	kingMakers, capturable := 0, 0
	for _, m := range movesList {
		squares := len(m.Squares)
		if m.Squares[squares-1].Y == 7 {
			kingMakers++
		}
		if jmp, _ := IsJump(m.Squares[0], m.Squares[1]); jmp {
			capturable += squares - 1
		}
	}

	return kingMakers, capturable
}

func eval(b *Board, player int) *boardEval {
	allMoves := [2][]*Move{}
	players := [2]int{1, 2}

	pieceCount := b.CountPieces()

	kingMakers := [2]int{}
	capturable := [2]int{}

	for i := 0; i < 2; i++ {
		allMoves[i] = b.AllMoves(players[i])
		kingMakers[i], capturable[i] = KingsCaptures(allMoves[i])
	}

	pieceDiff := (pieceCount.men[player-1] + pieceCount.kings[player-1]) -
		(pieceCount.men[Opposition(player)-1] + pieceCount.kings[Opposition(player)-1])

	kingDiff := pieceCount.kings[player-1] - pieceCount.kings[Opposition(player)-1]

	movesCountDiff := len(allMoves[player-1]) - len(allMoves[Opposition(player)-1])

	kingMakerDiff := kingMakers[player-1] - kingMakers[Opposition(player)-1]
	capturableDiff := capturable[player-1] - capturable[Opposition(player)-1]

	return &boardEval{
		player,
		pieceDiff,
		kingDiff,
		movesCountDiff,
		capturableDiff,
		kingMakerDiff,
	}
}

// HeuristicValue puts a number on `player`'s standing
func HeuristicValue(b *Board, player int) int {
	be := eval(b, player)

	return be.PieceDiff*pieceWeight +
		be.KingDiff*kingWeight +
		be.KingMakerDiff*kingMakerWeight +
		be.MovesCountDiff*movesCountWeight +
		be.CapturableDiff*capturableWeight
}
