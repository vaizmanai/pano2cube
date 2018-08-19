# pano2cube

Simple converter equirectangular to cubemap.

Example for using:

	f, _ := os.Open("pano.bmp")
	s, _ := f.Stat()
	b := make([]byte, s.Size())
	f.Read(b)
	imgInSrc, _ := bmp.Decode(bytes.NewReader(b))

	imgIn := image.NewRGBA(imgInSrc.Bounds())
	draw.Draw(imgIn, imgInSrc.Bounds(), imgInSrc, imgInSrc.Bounds().Min, draw.Over)

	inSize := imgInSrc.Bounds().Size()
	faceSize := inSize.X / 4

	for face := 0; face < 6; face++ {
		imgOut := image.NewRGBA(image.Rect(0,0,faceSize,faceSize))
		convertFace(imgIn, imgOut, face)
		dstFile, _ := os.Create("cube_" + fmt.Sprint(face) + ".bmp")
		bmp.Encode(dstFile, imgOut)
		dstFile.Close()
	}
