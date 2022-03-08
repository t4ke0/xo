package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	RedColor  string = "\x1b[31m"
	BlueColor string = "\x1b[34m"
	Reset     string = "\x1b[39m\x1b[0m"
)

const (
	width  int = 3
	height int = 3
)

// cell
type cell string

const (
	X       cell = "X"
	O            = "O"
	Default      = "."
)

// board
type board [][]cell

// drawBoard
func (b board) drawBoard() {

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var color string
			switch b[i][j] {
			case X:
				color = RedColor
			case O:
				color = BlueColor
			}
			if j == width-1 {
				if color != "" {
					fmt.Printf(color)
				}
				fmt.Printf("%+2v", b[i][j])
				fmt.Printf(Reset)
				continue
			}
			fmt.Printf(color)
			fmt.Printf("%+2v ", b[i][j])
			fmt.Printf(Reset)
			fmt.Printf(" %-2s", "|")
		}
		fmt.Println()
		if i == height-1 {
			continue
		}
		fmt.Printf("----------------\n")
	}

}

// updateBoard
func (b *board) updateBoard(x, y int, c cell) bool {

	if (*b)[x][y] == X || (*b)[x][y] == O {
		return false
	}

	(*b)[x][y] = c
	return true
}

// result
type result struct {
	winner cell
	won    bool
	isDraw bool
}

func (b board) checkVertPoints() result {

	for rowIndex, _ := range b {
		count := 0
		var currentCell cell
		for colIndex, _ := range b[rowIndex] {
			if currentCell != Default && currentCell == b[colIndex][rowIndex] {
				count++
			}
			currentCell = b[colIndex][rowIndex]
		}
		if count == height-1 {
			return result{
				winner: currentCell,
				won:    true,
			}
		}
	}

	return result{won: false}
}

func (b board) checkHorzPoints() result {

	for rowIndex, _ := range b {
		count := 0
		var currentCell cell
		for colIndex, _ := range b[rowIndex] {
			if currentCell != Default && currentCell == b[rowIndex][colIndex] {
				count++
			}
			currentCell = b[rowIndex][colIndex]
		}
		if count == width-1 {
			return result{
				winner: currentCell,
				won:    true,
			}
		}
	}

	return result{won: false}
}

func (b board) checkDrawSituation() bool {
	countEmptyCells := 0
	for _, r := range b {
		for _, c := range r {
			if c != X && c != O {
				countEmptyCells++
			}
		}
	}

	if countEmptyCells == 0 {
		return true
	}
	return false
}

func (b board) checkDiagonals() result {

	var (
		currentLeftCell  cell
		currentRightCell cell
	)

	xIndex := 0
	zIndex := len(b[xIndex]) - 1

	xCount := 0
	yCount := 0
	for r, _ := range b {
		if currentRightCell != Default && b[xIndex][r] == currentRightCell {
			xCount++
		}

		if currentLeftCell != Default && b[r][zIndex] == currentLeftCell {
			yCount++
		}

		currentRightCell = b[xIndex][r]
		currentLeftCell = b[r][zIndex]

		xIndex++
		zIndex--

		if xCount == width-1 {
			return result{won: true, winner: currentRightCell}
		}

		if yCount == width-1 {
			return result{won: true, winner: currentLeftCell}
		}

	}

	return result{won: false}
}

// checkBoard
func (b board) checkBoard() result {

	if r := b.checkVertPoints(); r.won {
		fmt.Println("verticall")
		return r
	}

	if r := b.checkHorzPoints(); r.won {
		fmt.Println("horizontal")
		return r
	}

	if r := b.checkDiagonals(); r.won {
		fmt.Println("diagonal")
		return r
	}

	isDraw := b.checkDrawSituation()

	if isDraw {
		return result{isDraw: isDraw}
	}

	return result{}
}

type position struct {
	x, y int
}

func (b board) numberToPosition(num string) (position, error) {

	if _, err := strconv.Atoi(num); err != nil {
		return position{}, err
	}

	for i, n := range b {
		for j, m := range n {
			if m != X && m != O && cell(num) == m {
				return position{
					x: i,
					y: j,
				}, nil
			}
		}
	}

	return position{}, nil

}

func generateBoard() *board {
	b := make(board, height)

	n := 1

	for rowIndex, _ := range b {
		for i := 0; i < width; i++ {
			b[rowIndex] = append(b[rowIndex], cell(strconv.Itoa(n)))
			n++
		}
	}

	return &b
}

func main() {

	// brd := &board{
	// 	[]cell{O, O, X},
	// 	[]cell{X, X, O},
	// 	[]cell{O, O, X},
	// }

	// brd.drawBoard()
	// r := brd.checkBoard()

	// fmt.Printf("%#v", r)

	brd := generateBoard()

	// return

	scanner := bufio.NewScanner(os.Stdin)

	var cellTurn cell = X

	var (
		stop      bool
		situation string
	)

	for {

		fmt.Printf("\033[H\033[J")
		brd.drawBoard()

		fmt.Printf("Enter a position number (%v) >  ", cellTurn)

		if stop {
			fmt.Println(situation)
			break
		}

		if !scanner.Scan() {
			break
		}
		cmd := scanner.Text()

		p, err := brd.numberToPosition(strings.TrimSpace(cmd))
		if err != nil {
			_ = err
			continue
		}

		brd.updateBoard(p.x, p.y, cellTurn)
		rlt := brd.checkBoard()
		if rlt.won {
			situation = fmt.Sprintf("%v won the game\n", rlt.winner)
			stop = true
			continue
		}

		if rlt.isDraw {
			situation = fmt.Sprintf("DRAW game\n")
			stop = true
			continue
		}

		if cellTurn == X {
			cellTurn = O
		} else {
			cellTurn = X
		}

	}

}
