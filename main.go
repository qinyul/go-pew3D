package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Seconds float32

type Proj struct {
	lifetime Seconds
	pos      rl.Vector3
	dir      rl.Vector3
}

const (
	SCREEN_HEIGHT = 800
	SCREEN_WIDTH  = 800

	PROJS_CAP     = 1000
	PROJ_WIDTH    = 0.1
	PROJ_HEIGHT   = 0.1
	PROJ_LEN      = 0.1
	PROJ_VEL      = 5
	PROJ_LIFETIME = 5

	PLAYER_GUN_LEN = 0.5
	PLAYER_HEIGHT  = 1
)

var (
	camera = rl.Camera{
		Position: rl.Vector3{
			X: 1,
			Y: PLAYER_HEIGHT,
			Z: 1,
		},
		Target: rl.Vector3Zero(),
		Up: rl.Vector3{
			X: 0,
			Y: 1,
			Z: 0,
		},
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}
	position  = rl.Vector3Zero()
	projs     = [PROJS_CAP]Proj{}
	projColor = rl.Red
)

func drawProj() {
	for i := 0; i < PROJS_CAP; i++ {
		if projs[i].lifetime > 0 {
			rl.DrawCube(projs[i].pos, PROJ_WIDTH, PROJ_HEIGHT, PROJ_LEN, projColor)
		}
	}

}

func updateProj() {
	dt := rl.GetFrameTime()
	for i := 0; i < PROJS_CAP; i++ {
		if projs[i].lifetime > 0 {
			projs[i].lifetime -= Seconds(dt)
			projs[i].pos = rl.Vector3Add(
				projs[i].pos,
				rl.Vector3Scale(projs[i].dir, rl.GetFrameTime()),
			)
		}
	}
}

func spawnProj(pos rl.Vector3, dir rl.Vector3) {
	for i := 0; i < PROJS_CAP; i++ {
		if projs[i].lifetime <= 0 {
			projs[i].lifetime = PROJ_LIFETIME
			projs[i].pos = pos
			projs[i].dir = dir
			break
		}
	}
}

func cameraDirection(camera rl.Camera) rl.Vector3 {
	return rl.Vector3Normalize(rl.Vector3Subtract(camera.Target, camera.Position))
}

func main() {

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "3D-Probe")
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetTargetFPS(60)
	model := rl.LoadModel("./guy/guy.iqm")
	texture := rl.LoadTexture("./guy/guytex.png")
	rl.SetMaterialTexture(model.Materials, rl.MapDiffuse, texture)

	defer rl.CloseWindow()
	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			dir := cameraDirection(camera)
			spawnProj(
				rl.Vector3Add(camera.Position, rl.Vector3Scale(dir, PLAYER_GUN_LEN)),
				rl.Vector3Scale(dir, PROJ_VEL),
			)
		}

		updateProj()
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.Black)
			rl.BeginMode3D(camera)
			{
				drawProj()
				rl.DrawModelEx(
					model,
					position,
					rl.Vector3{
						X: 1,
						Y: 0,
						Z: 0,
					},
					float32(-90),
					rl.Vector3{
						X: 1,
						Y: 1,
						Z: 1,
					},
					rl.White,
				)
				rl.DrawGrid(10, 1)
			}
			rl.EndMode3D()

			hudBuffer := fmt.Sprintf("Target: %f,%f,%f", camera.Target.X, camera.Target.Y, camera.Target.Z)
			rl.DrawText(hudBuffer, 10, 10, 20, rl.Maroon)
		}
		rl.EndDrawing()
	}
}
