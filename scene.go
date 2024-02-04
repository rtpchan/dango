// adopted from Hajime Hoshi Ebiten example Blocks
// See https://github.com/hajimehoshi/ebiten/blob/main/examples/blocks/blocks/scenemanager.go

package dango

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type SceneManager struct {
	current            Scene
	next               Scene
	transitionCount    int
	transitionMaxCount int
	transitionFrom     *ebiten.Image
	transitionTo       *ebiten.Image
}

// Create new scene manager, where w - screen width, h - screen height,
// transitionFrames - number of frames to transit between scenes
func NewSceneManager(w, h, transitionFrames int) *SceneManager {
	sm := &SceneManager{transitionMaxCount: transitionFrames}
	sm.transitionFrom = ebiten.NewImage(w, h)
	sm.transitionTo = ebiten.NewImage(w, h)
	return sm
}

func (s *SceneManager) Update() error {
	if s.transitionCount == 0 {
		return s.current.Update()
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	s.transitionFrom.Clear()
	s.current.Draw(s.transitionFrom)

	s.transitionTo.Clear()
	s.next.Draw(s.transitionTo)

	r.DrawImage(s.transitionFrom, nil)

	alpha := 1 - float32(s.transitionCount)/float32(s.transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.Scale(1, 1, 1, alpha)
	r.DrawImage(s.transitionTo, op)
}

func (s *SceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = s.transitionMaxCount
	}
}
