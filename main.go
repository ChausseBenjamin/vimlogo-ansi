/*
 * ----------------------------------------------------------------------------
 * "THE BEER-WARE LICENSE" (Revision 69):
 * <benjamin@chausse.xyz> wrote this file. As long as you retain this notice you
 * can do whatever you want with this stuff. If we meet some day, and you think
 * this stuff is worth it, you can buy me a beer in return.   Benjamin Chausse
 * ----------------------------------------------------------------------------
 */

package main

import (
	"fmt"
	"io"
	"os"
)

const (
	topFmt  = "\033[38;2;%d;%d;%dm"
	botFmt  = "\033[48;2;%d;%d;%dm"
	postFmt = "▀\033[0m"
)

type pixelArt struct {
	height uint32
	width  uint32
	pixels []uint32
}

type pixel struct {
	r uint8
	g uint8
	b uint8
}

func int2Pixel(src uint32) pixel {
	return pixel{
		b: uint8((src) & 0xFF),
		g: uint8((src >> 8) & 0xFF),
		r: uint8((src >> 16) & 0xFF),
	}
}

type drawing struct {
	height int
	width  int
	pixels []pixel
}

func mkPixelStack(top, bot pixel) string {
	noTop := (int(top.r+top.g+top.b) == 0)
	noBot := (int(bot.r+bot.g+bot.b) == 0)
	switch {
	case noTop && noBot:
		return " "
	case noTop:
		return fmt.Sprintf(botFmt+postFmt,
			bot.r, bot.g, bot.b)
	case noBot:
		return fmt.Sprintf(topFmt+postFmt,
			top.r, top.g, top.b)
	default:
		return fmt.Sprintf(topFmt+botFmt+postFmt,
			top.r, top.g, top.b,
			bot.r, bot.g, bot.b,
		)
	}
}

func (p pixelArt) draw(stream io.Writer) {
	step := uint32(p.width * 2)

	topIdx, botIdx := uint32(0), p.width
	row := uint32(0)

	top, bot := make([]pixel, p.width), make([]pixel, p.width)

	if p.height%2 != 0 {
		for i := range top {
			top[i] = pixel{0, 0, 0}
			bot[i] = int2Pixel(p.pixels[0])
		}
		drawRowStack(top, bot)
		topIdx += step
		botIdx += step
		row++
	}

	for row < p.height/2 {
		for i := range p.width {
			top[i] = int2Pixel(p.pixels[topIdx+i])
			bot[i] = int2Pixel(p.pixels[botIdx+i])
		}
		io.WriteString(stream, drawRowStack(top, bot))
		topIdx += step
		botIdx += step
		row++
	}
}

func drawRowStack(top, bot []pixel) string {
	var str string
	for i := 0; i < len(top); i++ {
		str += mkPixelStack(top[i], bot[i])
	}
	return str + "\n"
}

func main() {
	vim_logo := pixelArt{height: 72, width: 72, pixels: vim_data}
	vim_logo.draw(os.Stdout)
}
