package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

var mv *MainView

func main() {
	pv := NewPriceView()
	vecty.SetTitle("WASM Benzin")
	mv = &MainView{PriceView: pv, InputView: InputView{}}
	vecty.RenderBody(mv)
}

// MainView is our main page component.
type MainView struct {
	vecty.Core
	PriceView
	InputView
}

type PriceView struct {
	vecty.Core
	Gasoline *Gasoline
	selected int
}

type InputView struct {
	vecty.Core
	Liters      int
	TotalPrice  float64
	TotalLiters int
}

func (iv *InputView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Paragraph(vecty.Text(fmt.Sprintf("Liters: %d", iv.TotalLiters))),
		elem.Paragraph(vecty.Text(fmt.Sprintf("Price: %.2f BYR", iv.TotalPrice))),
		elem.Input(
			vecty.Markup(
				vecty.Property("value", "0"),
				event.Input(
					func(e *vecty.Event) {
						liters, err := strconv.Atoi(e.Target.Get("value").String())
						if err != nil {
							return
						}
						iv.Liters = liters
					},
				),
			),
		),
		elem.Button(
			vecty.Text("Add"),
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					// document.querySelector('input[name="gasoline"]:checked')
					gas := js.Global().Get("document").Call("querySelector", "input[name='gasoline']:checked").Get("value")
					n, err := strconv.Atoi(gas.String())
					if err != nil {
						return
					}
					iv.TotalPrice += mv.PriceView.CalculatePrice(n, iv.Liters)
					iv.TotalLiters += iv.Liters
					vecty.Rerender(mv)
				}),
			),
		),
	)
}

func NewPriceView() PriceView {
	gp := GetPrices()
	return PriceView{Gasoline: gp}
}

func (pv *PriceView) GasToString(n int) string {
	var price float64
	switch n {
	case 92:
		price = pv.Gasoline.Ai92
	case 95:
		price = pv.Gasoline.Ai95
	case 98:
		price = pv.Gasoline.Ai98
	case 1:
		price = pv.Gasoline.Diesel
		return fmt.Sprintf("Diesel: %.2f BYR", price)
	default:
		return ""
	}
	return fmt.Sprintf("AI-%d: %.2f BYR\n", n, price)
}

func (pv *PriceView) CalculatePrice(n, liters int) float64 {
	var price float64
	switch n {
	case 92:
		price = pv.Gasoline.Ai92
	case 95:
		price = pv.Gasoline.Ai95
	case 98:
		price = pv.Gasoline.Ai98
	case 1:
		price = pv.Gasoline.Diesel
	}
	return price * float64(liters)
}

func (pv *PriceView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Div(
			elem.Input(
				vecty.Markup(
					vecty.Property("type", "radio"),
					vecty.Property("id", "ai-92"),
					vecty.Property("name", "gasoline"),
					vecty.Property("value", "92"),
				),
			),
			elem.Label(
				vecty.Text(pv.GasToString(92)),
				vecty.Markup(
					vecty.Attribute("for", "ai-92"),
				),
			),
		),
		elem.Div(
			elem.Input(
				vecty.Markup(
					vecty.Property("type", "radio"),
					vecty.Property("id", "ai-95"),
					vecty.Property("name", "gasoline"),
					vecty.Property("value", "95"),
				),
			),
			elem.Label(
				vecty.Text(pv.GasToString(95)),
				vecty.Markup(
					vecty.Attribute("for", "ai-95"),
				),
			),
		),
		elem.Div(
			elem.Input(
				vecty.Markup(
					vecty.Property("type", "radio"),
					vecty.Property("id", "ai-98"),
					vecty.Property("name", "gasoline"),
					vecty.Property("value", "98"),
				),
			),
			elem.Label(
				vecty.Text(pv.GasToString(98)),
				vecty.Markup(
					vecty.Attribute("for", "ai-98"),
				),
			),
		),
		elem.Div(
			elem.Input(
				vecty.Markup(
					vecty.Property("type", "radio"),
					vecty.Property("id", "diesel"),
					vecty.Property("name", "gasoline"),
					vecty.Property("value", "1"),
				),
			),
			elem.Label(
				vecty.Text(pv.GasToString(1)),
				vecty.Markup(
					vecty.Attribute("for", "diesel"),
				),
			),
		),
	)
}

type Gasoline struct {
	Ai92   float64 `json:"ai-92"`
	Ai95   float64 `json:"ai-95"`
	Ai98   float64 `json:"ai-98"`
	Diesel float64 `json:"diesel"`
}

func GetPrices() *Gasoline {
	req, err := http.NewRequest("GET", "http://localhost:3000/gasoline", nil)
	req.Header.Add("js.fetch:mode", "cors")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var gp Gasoline
	err = json.Unmarshal(body, &gp)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &gp
}

// Render implements the vecty.Component interface.
func (p *MainView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		p.PriceView.Render(),
		p.InputView.Render(),
	)
}
