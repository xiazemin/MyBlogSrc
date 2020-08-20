---
title: transaction
layout: post
category: golang
author: 夏泽民
---
https://dev.to/techschoolguru/a-clean-way-to-implement-database-transaction-in-golang-2ba

https://github.com/techschool/simplebank
<!-- more -->

goUI


https://kpfaulkner.wordpress.com/2020/08/17/goui-a-very-simple-ui-framework/


package main
 
import (
    "github.com/hajimehoshi/ebiten"
    "github.com/kpfaulkner/goui/pkg"
    "github.com/kpfaulkner/goui/pkg/widgets"
    log "github.com/sirupsen/logrus"
    "image/color"
)
 
type MyApp struct {
    window pkg.Window
}
 
func NewMyApp() *MyApp {
    a := MyApp{}
    a.window = pkg.NewWindow(800, 600, "test app", false, false)
    return &a
}
 
func (m *MyApp) SetupUI() error {
    vPanel := widgets.NewVPanel("main vpanel", &color.RGBA{0, 0, 0, 0xff})
    m.window.AddPanel(vPanel)
    button1 := widgets.NewTextButton("text button 1", "my button1", true, 0, 0, nil, nil, nil, nil)
    vPanel.AddWidget(button1)
    return nil
}
 
func (m *MyApp) Run() error {
    m.SetupUI()
    ebiten.SetRunnableInBackground(true)
    ebiten.SetWindowResizable(true)
    m.window.MainLoop()
    return nil
}
 
func main() {
    log.SetLevel(log.DebugLevel)
    app := NewMyApp()
    app.Run()
}

https://github.com/kpfaulkner/goui

https://fyne.io/


ent - An Entity Framework For Go

https://github.com/facebook/ent