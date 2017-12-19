# Draughts game in Golang

This was the task set at the Bristol Schools' Code Jam 2017 competition where
I was part of the judging panel.

Implement the game of Draughts (English Checkers) as a two-player game. It
should:

1. Create the initial board state
1. Allow players to take alternate turns, and display the new board state
1. Verify that moves are legal
1. Determine if there is a winner
1. Implement capture moves
1. Implement king transitions

A pretty complex task given the 2.5 hours of coding time available. Teams could
pick any programming language they wanted, and part of the task was to produce
an initial plan or design for the solution.

Teams were using Java, C#, JavaScript, Python, VisualBasic and even Unity.
Whilst no team demonstrably got beyond the first three points above, we saw
some good skills and innovative solutions.

This is my solution that I sketched out during the competition and now fleshed
out to serve as an illustration to how it could be done.

## Representation

The game state is represented by the `Board` type, a struct holding an 8x8
array of `int`, as rows (Y) by cols (X). This means that the origin `[0][0]`
is the top-left corner of the board.

```go
type Board struct {
    state [8][8]int
}
```

An empty square is has the value 0, player 1 is 1 and player 2 is 2, with Kings
and Men distinguished by the sign of the value: Kings are negative and
Men are positive.

As viewed from the above, player 1 plays top down (Men moves with increasing Y),
and player 2 bottom up (Men moves with decreasing Y).

A square on the board is referenced via a coordinate pair `Pos{X, Y}`:

```go
type Pos struct {
    X, Y int
}
```

and a `Move` is a struct holding a slice of such squares visited by a piece:

```go
type Move struct {
    Squares []Pos
    Player  int
}
```

## Generating valid moves

By the rules of the game, a valid move is either a single step on either forward
diagonal (Kings can also go backwards on the diagonals) if the square is free,
or a sequence of capture moves (also on the forward diagonals) where the piece
jumps over an opposition piece, landing on a free square, and removing the
opposition piece that was jumped. After each such capture, the moving piece is
allowed to change diagonal (whilst still moving forward), but the sequence must
be followed until no further captures are available.

Generating a list of all valid moves for a player is probably the most difficult
bit of the competition task.

In this implementation, this is achieved by four functions, the entry point of
which is `AllMoves()`, calling `jumpMoves()`, `nonCaptureMoves()`, and
`singleJumps()`.

Here `AllMoves()` just loops over all pieces for the player, collecting every legal
move for each piece. The legal moves consist of any direct, non-capture moves as
returned by `nonCaptureMoves()` (available squares directly diagonal), and any
capture moves which are sequences of single jumps (available squares on the
forward diagonals if jumping over a single opponent piece), returned by
`singleJumps()`.

Stringing together single jump moves into a valid multi-hop capture move forms
a binary tree (actually a directed, acyclic graph, DAG, but `tree` is close
enough for our needs). We can generate this concisely using a recursive
depth-first traversal:

```go
func (b Board) jumpMoves(state *Board, square Pos, move *Move, moves *[]*Move) {
    player := move.Player
    jumps := validDirection(b.singleJumps(player, square), move)

    if len(jumps) == 0 && move.Length() > 1 {
        *moves = append(*moves, move)
        return
    }

    for _, jmp := range jumps {
        newState := state.Apply(NewMove(player, square, jmp))
        move.addSquare(jmp)                     // Record this in current move..
        b.jumpMoves(newState, jmp, move, moves) // Walk the tree depth-first

        // Now make a fresh move for when we follow the next branch
        move = NewMove(player, square)
    }
}
```

## Computer opponent

With the game mechanics in place, implementing a simplistic "AI" opponent is
relatively straight-forward. The function `minimax()` is a recursive descent
optimisation strategy to find the best move that simultaneously maximises the
player's standing and minimising that of the opponent to a certain depth of the
game tree. The key here is to have a good heuristic for evaluating board state.

The current implementation uses a polynomial combination of piece counts and
available moves, as given my the function `HeuristicValue()` in `minimax.go`.

The minimax implementation is naive, and would benefit from taking advantage of
the strategy known as alpha-beta pruning to allow it to search deeper in the
game tree. Currently it goes 6 levels deep.

## Game loop

The main game loop is found in `draughs.go`.

Game play uses the [Portable Draughts Notation](https://en.wikipedia.org/wiki/Portable_Draughts_Notation) square numbers to
enter moves.


    +-----+-----+-----+-----+-----+-----+-----+-----+
    |     |  1  |     |  2  |     |  3  |     |  4  |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    |  5  |     |  6  |     |  7  |     |  8  |     |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    |     |  9  |     | 10  |     | 11  |     | 12  |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    | 13  |     | 14  |     | 15  |     | 16  |     |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    |     | 17  |     | 18  |     | 19  |     | 20  |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    | 21  |     | 22  |     | 23  |     | 24  |     |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    |     | 25  |     | 26  |     | 27  |     | 28  |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    | 29  |     | 30  |     | 31  |     | 32  |     |
    +-----+-----+-----+-----+-----+-----+-----+-----+
    [Red: 12] [Green: 12]
    Red's move: 10 15

(although GitHub markdown won't show colours)
