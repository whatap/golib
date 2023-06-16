package logo

import (
	"fmt"

	"github.com/whatap/golib/util/ansi"
)

func Print(app string, version string) {
	fmt.Println(ansi.Green(" _       ____         ______"))
	fmt.Println(ansi.Green("| |     / / /_  ____ /_  __/___ _____   "))
	fmt.Println(ansi.Green("| | /| / / __ \\/ __ `// / / __ `/ __ \\ "))
	fmt.Println(ansi.Green("| |/ |/ / / / / /_/ // / / /_/ / /_/ / "))
	fmt.Println(ansi.Green("|__/|__/_/ /_/\\__,_//_/  \\__,_/ .___/  "))
	fmt.Println(ansi.Green("                             /_/      "))
	fmt.Println(ansi.Green("                                                "))
	fmt.Printf(ansi.Green("WhaTap %s ver %s                   \n"), app, version)
	fmt.Println(ansi.Green("Copyright ⓒ 2019 WhaTap Labs Inc. All rights reserved.\n"))
}

func Print2(app string, version string) {
	fmt.Println(ansi.Green(" _      ____       ______        		"))
	fmt.Println(ansi.Green("| | /| / / /  ___ /_  __/__ ____ "))
	fmt.Println(ansi.Green("| |/ |/ / _ \\/ _ `// / / _ `/ _ \\"))
	fmt.Println(ansi.Green("|__/|__/_//_/\\_,_//_/  \\_,_/ .__/"))
	fmt.Println(ansi.Green("                          /_/    "))
	fmt.Println(ansi.Green("                                           "))
	fmt.Printf(ansi.Green("WhaTap %s ver %s                   \n"), app, version)
	fmt.Println(ansi.Green("Copyright ⓒ 2019 WhaTap Labs Inc. All rights reserved.\n"))
}
