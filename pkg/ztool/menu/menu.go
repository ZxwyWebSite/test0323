package menu

import (
	"fmt"
	"os"

	"github.com/ZxwyWebSite/ztool"
)

type (
	// 显示接口
	View interface {
		Show() string
	}
	// 配置
	Conf struct {
		List List
	}
	// 页面
	Page struct {
		Name string // 页面名称
		Next string // 下一页
		Show func(app *App, this *Page) string
	}
	PageFunc func(this *App) string // app -> next
	// 缓存
	Cache struct {
		Last string // 上一页
		Now  string // 当前页
	}
	// 数据
	Data map[string]PageFunc // [名称]页面
	// 主菜单
	App struct {
		Name  string // 程序名称
		Conf  Conf
		cache Cache
		Data  Data
	}
)

const (
	Line  = "\n"                   // 下行
	Clear = "\033[H\033[2J" + Line // 清屏
)

// 新建菜单对象
func NewApp(name string) *App {
	return &App{
		Name: name,
		Conf: Conf{
			List: List{
				Prefix: `[`,
				Suffix: `]`,
				Priter: true,
			},
		},
		Data: Data{
			`Main`: func(this *App) string {
				this.print(this.ShowList([]string{
					"Item1", "Item2", "Item3", "Item4", "Item5", "Item6", "Item7", "Item8", "Item9", "Item10",
				}))
				fmt.Scanln()
				return `Exit`
			},
			`Exit`: func(this *App) string {
				os.Exit(0)
				return ``
			},
			`Page`: func(this *App) string {
				this.print(`Page`)
				return ``
			},
		},
	}
}

func (app *App) print(str string) {
	ztool.Cmd_FastPrintln(ztool.Str_FastConcat(
		`# `, app.cache.Now, ` | `, app.Name, ` #`, Line, Line,
		str,
	))
}

// 递归展示页面
func (app *App) Show(name string) {
	// page := data[name]
	// page.Show(page)
	// data.Show(page.Next)
	this, ok := app.Data[name]
	if !ok {
		return
	}
	app.cache.Now = name
	next := this(app)
	app.cache.Last = name
	app.Show(next)
}

// 运行程序
func (app *App) Run() {
	ztool.Cmd_FastPrintln(Clear)
	app.Show(`Main`)
	os.Exit(1)
	ztool.Cmd_FastPrintln(ztool.Str_FastConcat(
		`# `, app.cache.Now, ` | `, app.Name, ` #`, Line, Line,
		app.ShowMaps(map[string]any{
			`Ah`:      114514,
			`Hahah`:   `红红火火恍恍惚惚哈哈`,
			`1314520`: ``,
			`Notice`:  `逸一时误一世`,
			``:        `逸久逸久罢已零`,
			`nil`:     3.1415926,
		}),
		// app.showList([]string{
		// 	"Item1", "Item2", "Item3", "Item4", "Item5", "Item6", "Item7", "Item8", "Item9", "Item10",
		// 	"Item11", "Item12", "Item13", "Item14", "Item15", "Item16", "Item17", "Item18", "Item19", "Item20",
		// 	"Item1", "Item2", "Item3", "Item4", "Item5", "Item6", "Item7", "Item8", "Item9", "Item10",
		// 	"Item11", "Item12", "Item13", "Item14", "Item15", "Item16", "Item17", "Item18", "Item19", "Item20",
		// 	"Item1", "Item2", "Item3", "Item4", "Item5", "Item6", "Item7", "Item8", "Item9", "Item10",
		// 	"Item11", "Item12", "Item13", "Item14", "Item15", "Item16", "Item17", "Item18", "Item19", "Item20",
		// 	"Item1", "Item2", "Item3", "Item4", "Item5", "Item6", "Item7", "Item8", "Item9", "Item10",
		// 	"Item11", "Item12", "Item13", "Item14", "Item15", "Item16", "Item17", "Item18", "Item19", "Item20",
		// }),
	))
}
