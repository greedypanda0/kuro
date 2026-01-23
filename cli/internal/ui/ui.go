package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

/*
Color palette (edit once, theme everywhere)
*/
var (
	ColorPrimary = lipgloss.Color("#7D56F4")
	ColorInfo    = lipgloss.Color("#3FA9F5")
	ColorSuccess = lipgloss.Color("#2ECC71")
	ColorWarn    = lipgloss.Color("#F1C40F")
	ColorError   = lipgloss.Color("#E74C3C")
	ColorMuted   = lipgloss.Color("#666666")
)

/*
Base styles
*/
var (
	Base = lipgloss.NewStyle().
		Padding(0, 1)

	Bold = lipgloss.NewStyle().
		Bold(true)

	Muted = lipgloss.NewStyle().
		Foreground(ColorMuted)
)

/*
Headers
*/
func Header(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Padding(1, 2).
		Render(strings.ToUpper(text))
}

/*
Alerts
*/
func Info(text string) string {
	return alertBox("INFO", text, ColorInfo)
}

func Success(text string) string {
	return alertBox("OK", text, ColorSuccess)
}

func Warn(text string) string {
	return alertBox("WARN", text, ColorWarn)
}

func Error(text string) string {
	return alertBox("ERR", text, ColorError)
}

func alertBox(label, text string, color lipgloss.Color) string {
	tag := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Render(" " + label + " ")

	body := lipgloss.NewStyle().
		Foreground(color).
		Render(text)

	return lipgloss.JoinHorizontal(lipgloss.Left, tag, body)
}

/*
Info block (boxed text)
*/
func InfoBox(title, content string) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorInfo).
		Padding(1, 2).
		Width(50)

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorInfo).
		Render(title)

	return box.Render(header + "\n\n" + content)
}

/*
Bullet / indicator primitives
*/
func Bullet(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("• " + text)
}

func Step(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("› " + text)
}

func Check(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Render("✔ " + text)
}

func Cross(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorError).
		Render("✖ " + text)
}

func Simple(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render(text)
}

/*
Arrows / indicators
*/
func ArrowRight(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("→ " + text)
}

func ArrowDown(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("↓ " + text)
}

func ArrowUp(text string) string {
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Render("↑ " + text)
}

/*
Key-value row (nice for config / status)
*/
func KV(key, value string) string {
	k := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render(key)

	v := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render(value)

	return fmt.Sprintf("%s: %s", k, v)
}

/*
Progress bar (0.0 - 1.0)
*/
func Progress(width int, ratio float64) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(float64(width) * ratio)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	return lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Render(bar)
}

/*
Print helpers (UI-aware fmt)
*/
func Print(v ...any) string {
	return lipgloss.NewStyle().
		Render(fmt.Sprint(v...))
}

func Println(v ...any) {
	fmt.Println(Print(v...))
}
