package pano2cube

import (
	"image"
	"math"
	"github.com/gonum/floats"
	"image/color"
)


// get x,y,z coords from out image pixels coords
// i,j are pixel coords
// faceIdx is face number
// faceSize is edge length
func outImgToXYZ(i int, j int, faceIdx int, faceSize int)(float64, float64, float64) {
	a := 2.0 * float64(i) / float64(faceSize)
	b := 2.0 * float64(j) / float64(faceSize)

	var x,y,z float64

	if faceIdx == 0 { // back
		x,y,z = -1.0, 1.0 - a, 1.0 - b
	}else if faceIdx == 1{ // left
		x,y,z = a - 1.0, -1.0, 1.0 - b
	}else if faceIdx == 2 { // front
		x,y,z = 1.0, a - 1.0, 1.0 - b
	}else if faceIdx == 3 { // right
		x,y,z = 1.0 - a, 1.0, 1.0 - b
	}else if faceIdx == 4 { // top
		x,y,z = b - 1.0, a - 1.0, 1.0
	}else if faceIdx == 5 { //bottom
		x,y,z = 1.0 - b, a - 1.0, -1.0
	}

	return x, y, z
}

func  clip(vi float64, min int, max int) int {
	if int(vi) < min {
		return min
	} else if int(vi) > max {
		return max
	}
	return int(vi)
}

func round(src float64) int{
	return int(floats.Round(src, 0))
}

// convert using an inverse transformation
func convertFace(imgIn *image.RGBA, imgOut *image.RGBA, faceIdx int) {
	inSize := imgIn.Rect.Size()
	outSize := imgOut.Rect.Size()
	faceSize := outSize.X

	for xOut := 0; xOut < faceSize; xOut++ {
		for yOut := 0; yOut < faceSize; yOut++ {
			x,y,z := outImgToXYZ(xOut, yOut, faceIdx, faceSize)

			theta := math.Atan2(y,x) //# range -pi to pi
			rad := math.Hypot(x,y)
			phi := math.Atan2(z,rad) //# range -pi/2 to pi/2

			// source img coords
			uf := 0.5 * float64(inSize.X) * (theta + math.Pi) / math.Pi
			vf := 0.5 * float64(inSize.X) * (math.Pi/2 - phi) / math.Pi

			// Use bilinear interpolation between the four surrounding pixels
			ui := math.Floor(uf)  //# coord of pixel to bottom left
			vi := math.Floor(vf)
			u2 := ui+1       //# coords of pixel to top right
			v2 := vi+1
			mu := uf-ui      //# fraction of way across pixel
			nu := vf-vi

			// Pixel values of four corners
			A := imgIn.RGBAAt(int(ui) % inSize.X, clip(vi, 0, inSize.Y-1))
			B := imgIn.RGBAAt(int(u2) % inSize.X, clip(vi, 0, inSize.Y-1))
			C := imgIn.RGBAAt(int(ui) % inSize.X, clip(v2, 0, inSize.Y-1))
			D := imgIn.RGBAAt(int(u2) % inSize.X, clip(v2, 0, inSize.Y-1))

			// interpolate
			r,g,b := float64(A.R)*(1.0-mu)*(1.0-nu)+ float64(B.R)*((mu)*(1.0-nu)) + float64(C.R)*((1.0-mu)*nu) + float64(D.R)*(mu*nu),
						float64(A.G)*((1-mu)*(1-nu)) + float64(B.G)*((mu)*(1-nu)) + float64(C.G)*((1-mu)*nu) + float64(D.G)*(mu*nu),
						float64(A.B)*((1-mu)*(1-nu)) + float64(B.B)*((mu)*(1-nu)) + float64(C.B)*((1-mu)*nu) + float64(D.B)*(mu*nu)

			imgOut.SetRGBA(xOut, yOut, color.RGBA{R:uint8(round(r)), G:uint8(round(g)), B:uint8(round(b)), A:uint8(255)})
		}
	}
}
