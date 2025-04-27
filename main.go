package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Colors struct {
		Background string `toml:"background"`
		Text       string `toml:"text"`
		Highlight  string `toml:"highlight"`
	} `toml:"colors"`

	Browsers []BrowserEntry `toml:"browsers"`
}

type BrowserEntry struct {
	Name string `toml:"name"`
	Exec string `toml:"exec"`
	Icon string `toml:"icon"`
}

func loadDefaultConfig() *Config {
	return &Config{
		Colors: struct {
			Background string "toml:\"background\""
			Text       string "toml:\"text\""
			Highlight  string "toml:\"highlight\""
		}{
			Background: "#2E3440CC",
			Text:       "#ECEFF4FF",
			Highlight:  "#5E81ACAA",
		},
		Browsers: []BrowserEntry{
			{
				Name: "Firefox",
				Exec: "firefox",
				Icon: "assets/firefox.png",
			},
			{
				Name: "Chrome",
				Exec: "google-chrome-stable",
				Icon: "assets/chrome.png",
			},
			{
				Name: "Brave",
				Exec: "brave",
				Icon: "assets/brave.png",
			},
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func parseHexColor(s string) color.NRGBA {
	var r, g, b, a uint8 = 0, 0, 0, 255
	if len(s) == 7 { // "#RRGGBB"
		fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
		a = 255
	} else if len(s) == 9 { // "#RRGGBBAA"
		fmt.Sscanf(s, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	} else {
		log.Printf("Invalid color format: %s", s)
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// ClickableBox is a transparent clickable container
type ClickableBox struct {
	widget.BaseWidget
	OnTapped func()
	OnHover  func()
	content  fyne.CanvasObject
}

func NewClickableBox(content fyne.CanvasObject, tapped func(), hover func()) *ClickableBox {
	b := &ClickableBox{OnTapped: tapped, OnHover: hover, content: content}
	b.ExtendBaseWidget(b)
	return b
}

func (b *ClickableBox) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.content)
}

func (b *ClickableBox) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *ClickableBox) MouseIn(*desktop.MouseEvent) {
	if b.OnHover != nil {
		b.OnHover()
	}
}

func (b *ClickableBox) MouseOut() {}

func (b *ClickableBox) MouseMoved(*desktop.MouseEvent) {}

func main() {
	// Load config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "bifrost", "bifrost.toml")

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Printf("Config not found, using defaults: %v", err)
		config = loadDefaultConfig()
	}

	// Parse colors
	canvasBackgroundColor := parseHexColor(config.Colors.Background)
	textColor := parseHexColor(config.Colors.Text)
	highlightFillColor := parseHexColor(config.Colors.Highlight)

	myApp := app.New()
	w := myApp.NewWindow("Bifrost")

	w.SetFixedSize(true)
	w.SetPadded(false)
	w.SetMaster()

	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		url = "https://example.com"
	}

	selectedIndex := 0
	var browserBoxes []fyne.CanvasObject
	var browserFuncs []func()
	var highlightRects []*canvas.Rectangle

	updateHighlight := func() {
		for i, rect := range highlightRects {
			if i == selectedIndex {
				rect.FillColor = highlightFillColor
			} else {
				rect.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			}
			rect.Refresh()
		}
	}

	for i, browser := range config.Browsers {
		browser := browser // capture

		iconPath := filepath.Join(browser.Icon)
		fileBytes, err := os.ReadFile(iconPath)
		if err != nil {
			log.Printf("Failed to read image: %v", err)
		}

		img := canvas.NewImageFromResource(
			fyne.NewStaticResource(filepath.Base(iconPath), fileBytes),
		)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(80, 80))
		img.Refresh()

		label := canvas.NewText(fmt.Sprintf("%d. %s", i+1, browser.Name), textColor)
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.TextSize = 16
		label.Alignment = fyne.TextAlignCenter

		spacer := canvas.NewRectangle(color.Transparent)
		spacer.SetMinSize(fyne.NewSize(0, 10))

		stack := container.NewVBox(
			container.NewCenter(img),
			spacer,
			container.NewCenter(label),
		)

		centeredStack := container.NewCenter(stack)

		bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		bg.CornerRadius = 10
		highlightRects = append(highlightRects, bg)

		content := container.NewStack(bg, centeredStack)

		idx := i

		clickable := NewClickableBox(content, func() {
			fmt.Printf("Launching %s for URL %s\n", browser.Name, url)
			cmd := exec.Command(browser.Exec, url)
			err := cmd.Start()
			if err != nil {
				fmt.Printf("Failed to launch %s: %v\n", browser.Exec, err)
			}
			myApp.Quit()
		}, func() {
			selectedIndex = idx
			updateHighlight()
		})

		entry := container.NewGridWrap(fyne.NewSize(150, 150), clickable)

		browserBoxes = append(browserBoxes, entry)
		browserFuncs = append(browserFuncs, clickable.OnTapped)
	}

	hbox := container.NewHBox()
	for i, entry := range browserBoxes {
		hbox.Add(entry)
		if i != len(browserBoxes)-1 {
			hSpacer := canvas.NewRectangle(color.Transparent)
			hSpacer.SetMinSize(fyne.NewSize(20, 0))
			hbox.Add(hSpacer)
		}
	}

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyRight, fyne.KeyL:
			selectedIndex = (selectedIndex + 1) % len(browserBoxes)
			updateHighlight()
		case fyne.KeyLeft, fyne.KeyH:
			selectedIndex = (selectedIndex - 1 + len(browserBoxes)) % len(browserBoxes)
			updateHighlight()
		case fyne.KeyReturn, fyne.KeyEnter:
			browserFuncs[selectedIndex]()
		case fyne.Key1, fyne.Key2, fyne.Key3:
			num := int(k.Name[0] - '1')
			if num >= 0 && num < len(browserFuncs) {
				browserFuncs[num]()
			}
		case fyne.KeyEscape, fyne.KeyQ:
			myApp.Quit()
		}
	})

	updateHighlight()

	width := 160 * len(config.Browsers)
	if width < 480 {
		width = 480
	}
	height := 150 + 40

	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.SetFixedSize(true)

	leftSpacer := canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(20, 0))
	rightSpacer := canvas.NewRectangle(color.Transparent)
	rightSpacer.SetMinSize(fyne.NewSize(20, 0))

	background := canvas.NewRectangle(canvasBackgroundColor)

	w.SetContent(
		container.NewStack(
			background,
			container.NewHBox(leftSpacer, container.NewCenter(hbox), rightSpacer),
		),
	)

	w.ShowAndRun()
}
