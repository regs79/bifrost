package main

import (
	"embed"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pelletier/go-toml"
)

//go:embed assets/*
var embeddedAssets embed.FS

type Config struct {
	Colors struct {
		Background string `toml:"background"`
		Text       string `toml:"text"`
		Highlight  string `toml:"highlight"`
	} `toml:"colors"`

	Browsers []BrowserEntry `toml:"browsers"`
	ShowURL  bool           `toml:"show_url"` // Ensure this tag matches the key in the config file
}

type BrowserEntry struct {
	Name string `toml:"name"`
	Exec string `toml:"exec"`
	Icon string `toml:"icon"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	log.Printf("Config file contents:\n%s", string(data)) // Log the file contents

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	log.Printf("Unmarshaled config: %+v", config) // Log the unmarshaled config
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
	box := &ClickableBox{OnTapped: tapped, OnHover: hover, content: content}
	box.ExtendBaseWidget(box)
	return box
}

func (box *ClickableBox) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(box.content)
}

func (box *ClickableBox) Tapped(*fyne.PointEvent) {
	if box.OnTapped != nil {
		box.OnTapped()
	}
}

func (box *ClickableBox) MouseIn(*desktop.MouseEvent) {
	if box.OnHover != nil {
		box.OnHover()
	}
}

func (box *ClickableBox) MouseOut() {}

func (box *ClickableBox) MouseMoved(*desktop.MouseEvent) {}

func loadAsset(path string) ([]byte, error) {
	// Try loading from embedded assets
	data, err := embeddedAssets.ReadFile(path)
	if err == nil {
		// log.Printf("Loaded embedded asset: %s", path)
		return data, nil
	}

	// Fallback to loading from the working directory
	data, err = os.ReadFile(path)
	if err == nil {
		// log.Printf("Loaded asset from working directory: %s", path)
		return data, nil
	}

	log.Printf("Failed to load asset: %s, error: %v", path, err)
	return nil, err
}

func main() {
	// Load config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "bifrost", "bifrost.toml")

	myApp := app.New()

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	} else {
		log.Printf("Config loaded: %+v", config)
	}

	if err != nil {
		// Config not found or invalid: create default config file with example content
		var defaultConfig Config
		defaultConfig.Colors.Background = "#2E3440"
		defaultConfig.Colors.Text = "#ECEFF4"
		defaultConfig.Colors.Highlight = "#5E81AC"
		configDir := filepath.Dir(configPath)
		err := os.MkdirAll(configDir, 0o755)
		if err != nil {
			log.Fatalf("Could not create config directory: %v", err)
		}
		defaultConfigContent := `# Bifrost browser picker configuration
# Add your preferred browsers below. Example:

[[browsers]]
name = "Firefox"
exec = "firefox"
icon = "assets/firefox.png"

[colors]
background = "#2E3440"
text = "#ECEFF4"
highlight = "#5E81AC"

show_url = false
`
		err = os.WriteFile(configPath, []byte(defaultConfigContent), 0o644)
		if err != nil {
			log.Fatalf("Could not write default config file: %v", err)
		}
		// Show a Fyne window with instructions and exit
		w := myApp.NewWindow("Bifrost")
		w.SetFixedSize(true)
		w.SetPadded(false)
		w.SetMaster()

		msg1 := canvas.NewText("No configuration found.", parseHexColor(defaultConfig.Colors.Text))
		msg2 := canvas.NewText("A new config file has been created in:", parseHexColor(defaultConfig.Colors.Text))
		msg3 := canvas.NewText(configPath, parseHexColor(defaultConfig.Colors.Highlight))
		msg4 := canvas.NewText("Please edit this file and restart Bifrost.", parseHexColor(defaultConfig.Colors.Text))

		for _, msg := range []*canvas.Text{msg1, msg2, msg3, msg4} {
			msg.Alignment = fyne.TextAlignCenter
			msg.TextSize = 16
		}

		w.SetContent(container.NewStack(
			canvas.NewRectangle(parseHexColor(defaultConfig.Colors.Background)),
			container.NewVBox(
				layout.NewSpacer(),
				container.NewCenter(msg1),
				container.NewCenter(msg2),
				container.NewCenter(msg3),
				container.NewCenter(msg4),
				layout.NewSpacer(),
			),
		))

		w.Resize(fyne.NewSize(480, 220))
		w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			switch k.Name {
			case fyne.KeyEscape, fyne.KeyQ:
				myApp.Quit()
			}
		})
		w.ShowAndRun()
		return
	}

	if config.ShowURL == false {
		log.Println("show_url not set in config, defaulting to false")
	}

	bgColorHex := config.Colors.Background
	textColorHex := config.Colors.Text
	highlightHex := config.Colors.Highlight

	canvasBackgroundColor := parseHexColor(bgColorHex)
	textColor := parseHexColor(textColorHex)
	highlightFillColor := parseHexColor(highlightHex)

	w := myApp.NewWindow("Bifrost")

	w.SetFixedSize(true)
	w.SetPadded(false)
	w.SetMaster()

	if len(config.Browsers) == 0 {
		msg1 := canvas.NewText("No browsers found in config.", textColor)
		msg2 := canvas.NewText("You can add your browsers in:", textColor)
		msg3 := canvas.NewText(configPath, highlightFillColor)

		for _, msg := range []*canvas.Text{msg1, msg2, msg3} {
			msg.Alignment = fyne.TextAlignCenter
			msg.TextSize = 16
		}

		w.SetContent(container.NewStack(
			canvas.NewRectangle(canvasBackgroundColor),
			container.NewVBox(
				layout.NewSpacer(),
				container.NewCenter(msg1),
				container.NewCenter(msg2),
				container.NewCenter(msg3),
				layout.NewSpacer(),
			),
		))
		w.Resize(fyne.NewSize(480, 220))
		w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			switch k.Name {
			case fyne.KeyEscape, fyne.KeyQ:
				myApp.Quit()
			}
		})
		w.ShowAndRun()
		return
	}

	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		url = "https://example.com?lorem=ipsum&dolor=sit&amet=consectetur&adipiscing=elit"
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

		iconPath := browser.Icon // Use the Icon field directly
		// log.Printf("Attempting to load asset: %s", iconPath)
		fileBytes, err := loadAsset(iconPath)
		if err != nil {
			log.Printf("Failed to read image: %v", err)
			continue
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
			var cmd *exec.Cmd
			if runtime.GOOS == "darwin" {
				cmd = exec.Command("open", "-a", browser.Exec, url)
			} else {
				cmd = exec.Command(browser.Exec, url)
			}
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
		case fyne.KeyRight, fyne.KeyL, fyne.KeyUp, fyne.KeyK, fyne.KeyTab:
			selectedIndex = (selectedIndex + 1) % len(browserBoxes)
			updateHighlight()
		case fyne.KeyLeft, fyne.KeyH, fyne.KeyDown, fyne.KeyJ:
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

	width := max(160*len(config.Browsers), 400)
	height := 150 + 40

	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.SetFixedSize(true)

	leftSpacer := canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(40, 0))
	rightSpacer := canvas.NewRectangle(color.Transparent)
	rightSpacer.SetMinSize(fyne.NewSize(40, 0))

	background := canvas.NewRectangle(canvasBackgroundColor)

	// Add a text display for the URL
	urlLabel := canvas.NewText("Transporting to:", textColor)
	urlLabel.Alignment = fyne.TextAlignCenter
	urlLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Use a widget.Label for the URL to enable wrapping
	urlDisplay := widget.NewLabel(url)
	urlDisplay.Wrapping = fyne.TextWrapWord // Enable word wrapping

	// Create a container for the URL display
	urlWrapper := container.NewVBox(
		urlLabel,   // Add the "Transporting to:" label
		urlDisplay, // Add the URL display directly without centering each line
	)

	verticalPadding := func(height ...float32) fyne.CanvasObject {
		padHeight := float32(10) // Default height
		if len(height) > 0 {
			padHeight = height[0]
		}
		rect := canvas.NewRectangle(color.Transparent)
		rect.SetMinSize(fyne.NewSize(0, padHeight))
		return rect
	}

	var testText fyne.CanvasObject

	if config.ShowURL {
		testText = canvas.NewText("config on", textColor)
	} else {
		testText = canvas.NewText("config off", textColor)
	}

	// Create a container for the URL display and browser selection
	content := container.NewVBox(
		// Fixed top padding
		verticalPadding(),
		testText,
		urlWrapper,                             // URL display with padding
		canvas.NewRectangle(color.Transparent), // Add some spacing
		container.NewCenter(hbox),              // Add the browser selection
		verticalPadding(),                      // Fixed bottom padding
	)

	// Set the window content
	w.SetContent(
		container.NewStack(
			background,
			content,
		),
	)

	iconPath := filepath.Join("assets", "bifrost.png")
	// log.Printf("Attempting to load asset: %s", iconPath)
	iconBytes, err := loadAsset(iconPath)
	if err == nil {
		w.SetIcon(fyne.NewStaticResource("bifrost.png", iconBytes))
	} else {
		log.Printf("Failed to load window icon: %v", err)
	}

	w.ShowAndRun()
}
