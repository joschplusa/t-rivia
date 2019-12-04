/**
* Dear dev who is reading this source code,
* I'm really really sorry! I'm not a developer I'm just a marketing dude trying some new stuf
*/

package main

import (
	"log"
	"strconv"
	"os"
	//"math/rand"
	//"time"
	"image/color"

	_ "image/gif"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/gocarina/gocsv"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"github.com/golang/freetype/truetype"
	//"golang.org/x/image/font/gofont/gomonobold"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"

)

var isGameRunning bool;
var currentQuestion Question;
var currentQuestionIndex int = 0;
var width int = 640;
var height int = 360;
var displayedText string;
var QuestionText string;
//var correctAnswer string;
var AnswerAText string;
var AnswerBText string;
var AnswerCText string;
var AnswerDText string;
var answerPlayer1 string;
var answerPlayer2 string;
var player1Answered bool;
var player2Answered bool;
var no *ebiten.Image;
var scorePlayer1 int = 0;
var scorePlayer2 int = 0;
var questions []Question;
var splashScreen *ebiten.Image;
var backgroundImage *ebiten.Image;
var uiFont font.Face;
var uiFontMHeight int32;

var (
	gamepadIDs = map[int]struct{}{}
)
func init() {
	var err error
	// read my questions CSV file
	questionsFile, err := os.OpenFile("ressources/questions.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	splashScreen , _, err = ebitenutil.NewImageFromFile("ressources/splashscreen.png", ebiten.FilterDefault)
	backgroundImage , _, err = ebitenutil.NewImageFromFile("ressources/ui.png", ebiten.FilterDefault)
    if err != nil {
        panic(err)
    }
    defer questionsFile.Close()

    if err := gocsv.UnmarshalFile(questionsFile, &questions); err != nil {
        panic(err)
	}	
	// end reading my questions
	

    tt, err := truetype.Parse(fonts.ArcadeN_ttf)
    if err != nil {
            log.Fatal(err)
    }
    uiFont = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
        DPI:     72,
        Hinting: font.HintingFull,
    })
    b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = int32((b.Max.Y - b.Min.Y).Ceil())
		
	ebiten.SetFullscreen(true)
	
	if err != nil {
		log.Fatal(err)
	}
}

func update(screen *ebiten.Image) error {	
    
	if (isGameRunning){
		//fmt.Printf("P1 answered %v; P2 answered %v \n", player1Answered, player2Answered )
		screen.DrawImage(backgroundImage, nil)


		// get the current question
		currentQuestion = questions[currentQuestionIndex]
		// and set the right text
		displayedText = "\n\n\n\n\n                         " + currentQuestion.Question
		// shuffle our answers
		//currentQuestion.randomizeAnswers()
		// display our answers
		ebitenutil.DebugPrint(screen, "\n\n\n\n\n\n\n" + "                         A: " + currentQuestion.AnswerA)
		ebitenutil.DebugPrint(screen, "\n\n\n\n\n\n\n\n" + "                         B: " + currentQuestion.AnswerB)
		ebitenutil.DebugPrint(screen, "\n\n\n\n\n\n\n\n\n" + "                         C: " + currentQuestion.AnswerC)
		ebitenutil.DebugPrint(screen, "\n\n\n\n\n\n\n\n\n\n" + "                         D: " + currentQuestion.AnswerD)

		printLeft(screen, strconv.Itoa(scorePlayer1), 21, color.RGBA{0xff, 0, 0, 0xff})
		printRight(screen, strconv.Itoa(scorePlayer2), 21, color.RGBA{0xff, 0, 0, 0xff})
		// display the question
		//printCenter(screen, currentQuestion.Question , 12, color.RGBA{255, 255, 255, 0xff})

		// Player 1 answer A
		if (inpututil.IsKeyJustPressed(ebiten.KeyA) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0)) && !player1Answered {
			answerPlayer1 = currentQuestion.AnswerA
			player1Answered = true
		}
		// Player 1 answer B
		if (inpututil.IsKeyJustPressed(ebiten.KeyB) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1)) && !player1Answered {
			answerPlayer1 = currentQuestion.AnswerB
			player1Answered = true
		}
		// Player 1 answer C
		if (inpututil.IsKeyJustPressed(ebiten.KeyC) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton2)) && !player1Answered {
			answerPlayer1 = currentQuestion.AnswerC
			player1Answered = true
		}
		// Player 1 answer D
		if (inpututil.IsKeyJustPressed(ebiten.KeyD) || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton3)) && !player1Answered {
			answerPlayer1 = currentQuestion.AnswerD
			player1Answered = true
		}

		// Player 2 answer A
		if (inpututil.IsKeyJustPressed(ebiten.KeyUp) || ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton0)) && !player2Answered {
			answerPlayer2 = currentQuestion.AnswerA
			player2Answered = true
		}
		// Player 2 answer B
		if (inpututil.IsKeyJustPressed(ebiten.KeyRight)|| ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton1)) && !player2Answered {
			answerPlayer2 = currentQuestion.AnswerB
			player2Answered = true
		}
		// Player 2 answer C
		if (inpututil.IsKeyJustPressed(ebiten.KeyLeft) || ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton2)) && !player2Answered {
			answerPlayer2 = currentQuestion.AnswerC
			player2Answered = true
		}
		// Player 2 answer D
		if (inpututil.IsKeyJustPressed(ebiten.KeyDown) || ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton3)) && !player2Answered {
			answerPlayer2 = currentQuestion.AnswerD
			player2Answered = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			resetEverything()
			screen.Clear()
			isGameRunning = false
		}

		if player1Answered && player2Answered {
			if answerPlayer1 == currentQuestion.CorrectAnswer {
				//that's right!
				scorePlayer1++
			} else {
				//wrong you fool!
			}

			if answerPlayer2 == currentQuestion.CorrectAnswer {
				//that's right!
				scorePlayer2++
			} else {
				//wrong you fool!
			}

			if currentQuestionIndex < len(questions)-1 {
				currentQuestionIndex++
				screen.Clear()
				//displayedText = "Off to question number " + strconv.Itoa(currentQuestionIndex + 1) + "!"
				player1Answered = false
				player2Answered = false
			} else {
				// Our game is finished! 
				
				if scorePlayer1 > scorePlayer2 {
					//winnerText = "Player 1 wins!\nPlayer 2 is a loser"
					printCenter(screen, "Player 1 wins!" , 12, color.RGBA{255, 255, 255, 0xff})
					printCenter(screen, "Player 2 is a loser!" , 13, color.RGBA{255, 255, 255, 0xff})

				} else if scorePlayer2 > scorePlayer1 {
					//winnerText = "Player 2 wins!\nPlayer 1 is a loser"
					printCenter(screen, "Player 2 wins!" , 12, color.RGBA{255, 255, 255, 0xff})
					printCenter(screen, "Player 1 is a loser!" , 13, color.RGBA{255, 255, 255, 0xff})

				} else {
					//winnerText = "Tie! You both win!\n(or loose both?)"
					printCenter(screen, "Draw! You both win!" , 12, color.RGBA{255, 255, 255, 0xff})
					printCenter(screen, "(or both loose?)" , 13, color.RGBA{255, 255, 255, 0xff})

				}
				
			}
			
		}

	} else if (inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyLeft)){
		resetEverything()
		screen.Clear()
		isGameRunning = true
		ebitenutil.DebugPrint(screen, "\n\n\n" + "  LET'S GO! " )

	} else {
		screen.DrawImage(splashScreen, nil)
	}

	if inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton1){
		ebitenutil.DebugPrint(screen, "A: 1 pressed!")
	}

	if inpututil.IsGamepadButtonJustPressed(1, ebiten.GamepadButton1){
		ebitenutil.DebugPrint(screen, "B: 1 pressed!")
	}

	if ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton0) && ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton1) && ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton2) && ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton3) && ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton0) && ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton1) && ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton2) && ebiten.IsGamepadButtonPressed(1, ebiten.GamepadButton3) {
		os.Exit(0)
	}

	ebitenutil.DebugPrint(screen, displayedText)	

    return nil
}

func main() {
	if err := ebiten.Run(update, width, height, 2, "t-rivia - Das tarent Quiz"); err != nil {
		log.Fatal(err)
	}
}

/*func setQuestion(){
	var q = Question{"Was ist der korrekte Farbcode für das tarent rot?", "#CC000", "#FF2222", "#002355", "#4257f5", ""}

	QuestionText = q.Question
	AnswerAText = q.AnswerA
	AnswerBText = q.AnswerB
	AnswerCText = q.AnswerC
	AnswerDText = q.AnswerD
	//correctAnswer = AnswerAText
}*/

func reset(){
	QuestionText = ""
	//correctAnswer = ""
	AnswerAText  = ""
	AnswerBText = ""
	AnswerCText = ""
	AnswerDText = ""
	scorePlayer1 = 0
	scorePlayer2 = 0
	player1Answered = false
	player2Answered = false
}

type Question struct {
    Question string `csv:"question"`
	AnswerA string `csv:"answer_a"`
	AnswerB string `csv:"answer_b"`
	AnswerC string `csv:"answer_c"`
	AnswerD string `csv:"answer_d"`
	CorrectAnswer string `csv:"answer_correct"`

}

func (q *Question) reset() {
	q.Question = ""
	q.AnswerA = ""
	q.AnswerB = ""
	q.AnswerC = ""
	q.AnswerD = ""
	q.CorrectAnswer = ""
}

func (q *Question) randomizeAnswers() {
	/*randomIndex := []int{1, 2, 3, 4}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

	for i := 0; i < 4; i++ {
		
	}
	q.AnswerA = ""
	q.AnswerB = ""
	q.AnswerC = ""
	q.AnswerD = ""*/
}

func printRight(screen *ebiten.Image, s string, line int32, c color.Color) {
	bounds, _ := font.BoundString(uiFont, s)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	x := int(width - 68 - w)
	y := int(35 + (uiFontMHeight + 4) * line)
	text.Draw(screen, s, uiFont, x, y, c)
}

func printLeft(screen *ebiten.Image, s string, line int32, c color.Color) {
	x := 68
	y := int(35 + (uiFontMHeight + 4) * line)
	text.Draw(screen, s, uiFont, x, y, c)
}

func printCenter(screen *ebiten.Image, s string, line int32, c color.Color) {
	bounds, _ := font.BoundString(uiFont, s)
	w := (bounds.Max.X - bounds.Min.X).Ceil()
	x := int(width/2 - w)
	y := int(16 + (uiFontMHeight + 4) * line)
	text.Draw(screen, s, uiFont, x, y, c)
}

func resetEverything(){
	isGameRunning = false
	currentQuestionIndex = 0
	answerPlayer2 = ""
	answerPlayer1 = ""
	player1Answered = false
	player2Answered = false
	displayedText = ""
	scorePlayer1 = 0
	scorePlayer2 = 0
}