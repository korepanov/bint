package drawModule

import(
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"strconv"
)
const imgWidth = 1600
const imgHeight = 800
var green = color.RGBA{0, 200, 0, 0}


func MakeImg() gocv.Mat{
	img := gocv.NewMatWithSize(imgHeight, imgWidth, gocv.MatTypeCV8UC3)
	gray := color.RGBA{200, 200, 200, 0}

	// красим изображение в серый цвет
	gocv.Rectangle(&img, image.Rect(0, 0, imgWidth - 1, imgHeight - 1), gray, -1)

	return img
}

func makeCircle(img gocv.Mat, x int, y int, thickness int) (gocv.Mat){

	gocv.Circle(&img, image.Pt(x, y), 20, green, thickness)

	return img
}

func makeLeftConnection(img gocv.Mat, x int, y int, isNextFilled bool, delta int) (int, int){
	var th int
	if isNextFilled{
		th = -1
	}else{
		th = 3
	}
	leftX := x - 390 + delta
	leftY := y + 50
	gocv.Line(&img, image.Pt(x - 15, y + 15), image.Pt(leftX - 15, leftY - 5), green, 3)
	makeCircle(img, leftX - 15, leftY + 15, th)

	return leftX - 15, leftY + 15
}

func makeRightConnection(img gocv.Mat, x int, y int, isNextFilled bool, delta int) (int, int) {
	var th int
	if isNextFilled{
		th = -1
	}else{
		th = 3
	}

	leftX := x + 390 - delta
	leftY := y + 50
	gocv.Line(&img, image.Pt(x + 15, y + 15), image.Pt(leftX, leftY), green, 3)
	makeCircle(img, leftX + 15, leftY + 15, th)
	return leftX + 15, leftY + 15
}

func makeLevel(history [][]int) int{
	level := len(history)
	if level > 5{
		level = 5
	}
	return level
}

func makeTree(img gocv.Mat, number string, rootX int, rootY int,
	delta int, history [][]int, infoList []interface{}, nextTopNumber int){

	// history = [[x, y, wasLeft, wasRight]]
	wasNull := false
	x := rootX
	y := rootY
	breakFlag := false
	for i:=0; i < len(number); i++{
		if breakFlag{
			break
		}
		if "1" == string(number[i]) && !wasNull{
			x, y = makeLeftConnection(img, x, y, true, delta * makeLevel(history))
			err := printNumber(img, x, y, nextTopNumber)
			if nil != err{
				fmt.Println(err)
			}
			printInfo(img, x, y, fmt.Sprintf("%v", infoList[nextTopNumber]))
			nextTopNumber += 1
			history[len(history) - 1][2] = 1
			history = append(history, []int{x, y, 0, 0})
		}else if "0" == string(number[i]) && !wasNull{
			wasNull = true
			makeLeftConnection(img, x, y, false, delta * makeLevel(history))
			history[len(history) - 1][2] = 1
		}else if "1" == string(number[i]) && wasNull{
			wasNull = false
			if len(history) > 0{
				x, y = makeRightConnection(img, history[len(history) - 1][0], history[len(history) - 1][1],
					true, delta * makeLevel(history))
				err := printNumber(img, x, y, nextTopNumber)
				if nil != err{
					fmt.Println(err)
				}
				printInfo(img, x, y, fmt.Sprintf("%v", infoList[nextTopNumber]))
				nextTopNumber += 1
				history[len(history)-1][3] = 1
				history = append(history, []int{x, y, 0, 0})
				newNumber := number[i + 1:]
				makeTree(img, newNumber, x, y, delta, history, infoList, nextTopNumber)
				breakFlag = true
			}
		}else if ("0" == string(number[i]) && wasNull) || ("0" == string(number[i]) ||  1 == history[1][2]){
			if len(history) > 0{
				makeRightConnection(img, history[len(history) - 1][0], history[len(history) - 1][1],
					false, delta * makeLevel(history))
				history[len(history)-1][3] = 1
				if len(history) > 1{
					history = history[:len(history) - 1]

					for (1 == history[len(history)-1][2]) && (1 == history[len(history) - 1][3]) && len(history) > 1{
						history = history[:len(history) - 1]
					}
				}

			}
		}

	}
}

func printNumber(img gocv.Mat, x int, y int, number int) error{
	font := gocv.FontHersheySimplex
	fontScale := 0.6
	blue := color.RGBA{0, 0, 255, 0}
	thickness := 1
	strNumber := strconv.Itoa(number)

	if number < 10{
		strNumber = "0" + strNumber
	}

	gocv.PutText(&img, strNumber, image.Pt(x - 40, y - 10), font, fontScale, blue, thickness)
	return nil
}

func printInfo(img gocv.Mat, x int, y int, info string){
	font := gocv.FontHersheySimplex
	fontScale := 0.6
	c := color.RGBA{123, 30, 0, 0}
	thickness := 1

	gocv.PutText(&img, info, image.Pt(x - 6, y + 5), font, fontScale, c, thickness)
}

func DrawTree(number string, infoList []interface{}){
	img := MakeImg()
	rootX := 800
	rootY := 25

	if "1" != string(number[0]){
		panic("drawModule: ERROR: the number must starts with 1")
	}
	if "0" != string(number[len(number)-1]){
		panic("drawModule: ERROR: the number must ends with 0")
	}

	makeCircle(img, rootX, rootY, -1)
	err := printNumber(img, rootX, rootY, 0)
	printInfo(img, rootX, rootY, fmt.Sprintf("%v", infoList[0]))
	if nil != err{
		fmt.Println(err)
	}
	number = number[1:]
	if len(infoList) > 1{
		makeTree(img, number, rootX, rootY, 75, [][]int{{rootX, rootY, 0, 0}}, infoList, 1)
	}else{
		makeLeftConnection(img, rootX, rootY, false, 75)
		makeRightConnection(img, rootX, rootY, false, 75)
	}

	err = showWindow(img)
	if nil != err{
		fmt.Println("drawModule: ERROR: could not show window")
		fmt.Println(err)
	}
}

func showWindow(img gocv.Mat) error{
	window := gocv.NewWindow("tree")

	window.ResizeWindow(imgWidth, imgHeight)


	for {

		window.IMShow(img)

		if window.WaitKey(10) >= 0  {
			break
		}
	}

	err := window.Close()

	if nil != err{
		return err
	}

	err = img.Close()

	return err
}