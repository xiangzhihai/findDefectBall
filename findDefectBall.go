package findDefectBall

import (
	"fmt"
	"math"
)

// Define a set

var exists = struct{}{}

type set struct {
	m map[int]struct{}
}

func NewSet() *set {
	s := &set{}
	s.m = make(map[int]struct{})
	return s
}

func (s *set) Add(value int) {
	s.m[value] = exists
}

func (s *set) Remove(value int) {
	delete(s.m, value)
}

func (s *set) Contains(value int) bool {
	_, c := s.m[value]
	return c
}

func (s *set) Copy(newS *set) {
	for key := range newS.m {
		s.m[key] = exists
	}
}

func (s *set) Pop() int {
	// Assumes the set has enough elements
	var key int
	for key = range s.m {
		s.Remove(key)
		break
	}
	return key
}

func (s *set) Size() int {
	return len(s.m)
}

// Defect status
type defectStatus int

const (
	light defectStatus = iota
	heavy
)

type defection struct {
	number int
	status defectStatus
}

// Change these parameters for testing
var numOfBalls = 11
var DefectBall = defection{number: 2, status: heavy}

func CreateBalls(numBalls int) *set {
	balls := NewSet()
	for i := 0; i < numBalls; i++ {
		balls.Add(i)
	}
	return balls
}

// Balance status
type status int

const (
	left status = iota
	right
	even
)

// Compares balls, Assumes leftBalls and leftBalls have the same
// number of balls, and only one or less ball is defect
func compare(leftBalls, rightBalls *set) status {
	if leftBalls.Contains(DefectBall.number) {
		if DefectBall.status == heavy {
			return left
		} else {
			return right
		}
	}
	if rightBalls.Contains(DefectBall.number) {
		if DefectBall.status == light {
			return left
		} else {
			return right
		}
	}
	return even
}

// Solve
func Solve(balls *set) defection {
	// Base case
	if balls.Size() <= 2 {
		leftBalls := NewSet()
		leftBalls.Add(balls.Pop())
		rightBalls := NewSet()
		rightBalls.Add(-1)
		res := compare(leftBalls, rightBalls)
		switch res {
		case even:
			return Solve(balls)
		case left:
			return defection{
				number: leftBalls.Pop(),
				status: heavy,
			}
		case right:
			return defection{
				number: leftBalls.Pop(),
				status: light,
			}
		}
	}
	// Devide to three portions
	numBallsOnBal := math.Round(float64(balls.Size()) / 3.0)
	leftBalls := NewSet()
	rightBalls := NewSet()
	for i := 0; i < int(numBallsOnBal); i++ {
		leftBalls.Add(balls.Pop())
		rightBalls.Add(balls.Pop())
	}
	res := compare(leftBalls, rightBalls)
	switch res {
	case even:
		return Solve(balls)
	case left:
		return solveUnequal(leftBalls, rightBalls)
	default:
		return solveUnequal(rightBalls, leftBalls)
	}
}

func solveUnequal(heavyBalls, lightBalls *set) defection {
	totalNumBalls := heavyBalls.Size() + lightBalls.Size()
	if totalNumBalls <= 2 {
		leftBalls := NewSet()
		getBall := func(Balls *set) (int, bool) {
			if Balls.Size() > 0 {
				return Balls.Pop(), true
			}
			return -1, false
		}
		// Add ball to left
		if ball, notEmpty := getBall(heavyBalls); notEmpty {
			leftBalls.Add(ball)
		} else {
			ball, _ := getBall(lightBalls)
			leftBalls.Add(ball)
		}
		// Add a normal ball to right
		rightBalls := NewSet()
		rightBalls.Add(-1)
		res := compare(leftBalls, rightBalls)
		switch res {
		case even:
			if ball, notEmpty := getBall(heavyBalls); notEmpty {
				return defection{
					number: ball,
					status: heavy,
				}
			} else {
				return defection{
					number: lightBalls.Pop(),
					status: light,
				}
			}
		case left:
			return defection{
				number: leftBalls.Pop(),
				status: heavy,
			}
		default:
			return defection{
				number: leftBalls.Pop(),
				status: light,
			}
		}
	}
	numBallsOnBal := int(math.Round(float64(totalNumBalls) / 3.0))
	var largeBalls *set
	var smarllBalls *set
	heavyBallsCopy := NewSet()
	heavyBallsCopy.Copy(heavyBalls)
	lightBallsCopy := NewSet()
	lightBallsCopy.Copy(lightBalls)
	if heavyBalls.Size() > lightBalls.Size() {
		largeBalls = heavyBalls
		smarllBalls = lightBalls
	} else {
		smarllBalls = heavyBalls
		largeBalls = lightBalls
	}
	ballPool := largeBalls
	leftBalls := NewSet()
	rightBalls := NewSet()
	for i := 0; i < numBallsOnBal; i++ {
		leftBalls.Add(ballPool.Pop())
		rightBalls.Add(ballPool.Pop())
		if ballPool == largeBalls {
			ballPool = smarllBalls
		} else {
			ballPool = largeBalls
		}
	}
	res := compare(leftBalls, rightBalls)
	switch res {
	case even:
		return solveUnequal(heavyBalls, lightBalls)
	case left:
		for key := range leftBalls.m {
			if lightBallsCopy.Contains(key) {
				leftBalls.Remove(key)
			}
		}
		for key := range rightBalls.m {
			if heavyBallsCopy.Contains(key) {
				heavyBalls.Remove(key)
			}
		}
		return solveUnequal(leftBalls, rightBalls)
	default:
		for key := range leftBalls.m {
			if heavyBallsCopy.Contains(key) {
				leftBalls.Remove(key)
			}
		}
		for key := range rightBalls.m {
			if lightBallsCopy.Contains(key) {
				heavyBalls.Remove(key)
			}
		}
		return solveUnequal(rightBalls, leftBalls)
	}
}

func main() {
	balls := CreateBalls(numOfBalls)
	def := Solve(balls)
	if def == DefectBall {
		fmt.Println("Pass")
	} else {
		fmt.Println("Error")
	}

	// Todo: Add unit test
}
