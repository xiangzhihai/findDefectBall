package findDefectBall

import (
	"testing"
)

func TestFindBall(t *testing.T) {
	balls := CreateBalls(12)
	def := Solve(balls)
	if def != DefectBall {
		t.Error("Wrong ball was found!")
	}
}
