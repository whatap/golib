package ansi

/*
 *  Copyright 2015 Scouter Project.
 *  @https://github.com/scouter-project/scouter
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

import (
	"fmt"
)

var (
	enable bool
)

const (
	ANSI_RESET  = "\u001B[0m"
	ANSI_BOLD   = "\u001B[1m"
	ANSI_DIM    = "\u001B[2m"
	ANSI_BLACK  = "\u001B[30m"
	ANSI_RED    = "\u001B[31m"
	ANSI_GREEN  = "\u001B[32m"
	ANSI_YELLOW = "\u001B[33m"
	ANSI_BLUE   = "\u001B[34m"
	ANSI_PURPLE = "\u001B[35m"
	ANSI_CYAN   = "\u001B[36m"
	ANSI_WHITE  = "\u001B[37m"

	// Bright colors
	ANSI_BRIGHT_RED     = "\u001B[91m"
	ANSI_BRIGHT_GREEN   = "\u001B[92m"
	ANSI_BRIGHT_YELLOW  = "\u001B[93m"
	ANSI_BRIGHT_BLUE    = "\u001B[94m"
	ANSI_BRIGHT_MAGENTA = "\u001B[95m"
	ANSI_BRIGHT_CYAN    = "\u001B[96m"
)

func init() {
	//SystemUtil.IS_WINDOWS == false;
	enable = true
}

func Red(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_RED, s, ANSI_RESET)
}

func Yellow(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_YELLOW, s, ANSI_RESET)
}

func Green(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_GREEN, s, ANSI_RESET)
}

func Cyan(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_CYAN, s, ANSI_RESET)
}

func Blue(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BLUE, s, ANSI_RESET)
}

func Magenta(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_PURPLE, s, ANSI_RESET)
}

func White(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_WHITE, s, ANSI_RESET)
}

func Bold(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BOLD, s, ANSI_RESET)
}

func BrightRed(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_RED, s, ANSI_RESET)
}

func BrightGreen(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_GREEN, s, ANSI_RESET)
}

func BrightYellow(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_YELLOW, s, ANSI_RESET)
}

func BrightBlue(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_BLUE, s, ANSI_RESET)
}

func BrightMagenta(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_MAGENTA, s, ANSI_RESET)
}

func BrightCyan(s string) string {
	if enable == false {
		return s
	}
	return fmt.Sprintf("%s%s%s", ANSI_BRIGHT_CYAN, s, ANSI_RESET)
}

func RedOut(s string) {
	fmt.Printf("%s\n", Red(s))
}
