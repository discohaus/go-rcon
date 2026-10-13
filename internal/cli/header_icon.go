package cli

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Generated from identity/assets/go-rcon/chip.svg. The PNG is embedded so
// installed binaries do not need the submodule or an SVG renderer at runtime.
//
//go:generate make -C ../.. gen
//go:embed go_rcon_chip.png
var goRconChipPNG []byte

const iconImageID = 0x72636f // "rco", encoded in the placeholder's RGB color

type iconProtocol uint8

const (
	noIcon iconProtocol = iota
	kittyIcon
	footSixelIcon
)

func detectIconProtocol() iconProtocol {
	return iconProtocolFor(os.Getenv("TERM"), os.Getenv("TERM_PROGRAM"),
		os.Getenv("TMUX") != "" || os.Getenv("STY") != "", runningInFoot())
}

func iconProtocolFor(term, program string, multiplexer, foot bool) iconProtocol {
	if multiplexer {
		return noIcon
	}
	term = strings.ToLower(term)
	program = strings.ToLower(program)
	if strings.Contains(term, "kitty") || strings.Contains(term, "ghostty") || program == "ghostty" {
		return kittyIcon
	}
	if term == "foot" || term == "foot-direct" || foot {
		return footSixelIcon
	}
	return noIcon
}

// Foot does not set an identifying environment variable, and users can set
// TERM=xterm-256color. On Linux, the shell's parent chain identifies local Foot.
func runningInFoot() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	pid := os.Getppid()
	for range 8 {
		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
		if err != nil {
			return false
		}
		name, parent := "", 0
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "Name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			}
			if strings.HasPrefix(line, "PPid:") {
				parent, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "PPid:")))
			}
		}
		if name == "foot" || name == "footclient" {
			return true
		}
		if parent <= 1 || parent == pid {
			return false
		}
		pid = parent
	}
	return false
}

// Kitty's Unicode placeholders stay attached to cells when Bubble Tea repaints
// the header. Other image placements can be left behind by TUI redraws.
var kittyIconSequence = makeKittyIcon()

func makeKittyIcon() string {
	payload := base64.StdEncoding.EncodeToString(goRconChipPNG)
	return "\x1b_Ga=T,f=100,i=" + strconv.Itoa(iconImageID) + ",U=1,c=2,r=1,q=2;" + payload + "\x1b\\" +
		"\x1b[38;2;114;99;111m" +
		"\U0010eeee\u0305\u0305\U0010eeee\u0305\u030d" +
		"\x1b[39m "
}

// Sixel is painted after the header text, then the cursor is restored before
// Bubble Tea writes the next line. Transparent pixels above the mark align it
// with the text instead of the top edge of the cell.
var footIconSequence = makeFootIcon()

func makeFootIcon() string {
	img, err := png.Decode(bytes.NewReader(goRconChipPNG))
	if err != nil {
		return ""
	}
	const size = 28
	const topPadding = 3
	var pixels [size][size]color.NRGBA
	var palette []color.NRGBA
	indices := make(map[[3]uint8]int)
	for y := range size {
		for x := range size {
			c := color.NRGBAModel.Convert(img.At(x*img.Bounds().Dx()/size, y*img.Bounds().Dy()/size)).(color.NRGBA)
			pixels[y][x] = c
			if c.A < 128 {
				continue
			}
			key := [3]uint8{c.R, c.G, c.B}
			if _, ok := indices[key]; !ok {
				indices[key] = len(palette)
				palette = append(palette, c)
			}
		}
	}

	var out strings.Builder
	fmt.Fprintf(&out, "\x1bP0;1;0q\"1;1;%d;%d", size, size+topPadding) // transparent background
	for i, c := range palette {
		fmt.Fprintf(&out, "#%d;2;%d;%d;%d", i, percent(c.R), percent(c.G), percent(c.B))
	}
	for band := 0; band < size+topPadding; band += 6 {
		if band > 0 {
			out.WriteByte('-')
		}
		for i, c := range palette {
			if i > 0 {
				out.WriteByte('$')
			}
			fmt.Fprintf(&out, "#%d", i)
			for x := range size {
				mask := 0
				for bit := range 6 {
					y := band + bit - topPadding
					if y < 0 || y >= size {
						continue
					}
					p := pixels[y][x]
					if p.A >= 128 && p.R == c.R && p.G == c.G && p.B == c.B {
						mask |= 1 << bit
					}
				}
				out.WriteByte(byte(63 + mask))
			}
		}
	}
	out.WriteString("\x1b\\")
	return out.String()
}

func percent(v uint8) int { return (int(v)*100 + 127) / 255 }
