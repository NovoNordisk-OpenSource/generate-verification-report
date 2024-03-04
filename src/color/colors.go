package color

func Purple(s string) string {
	return "\033[0;35m" + s + "\033[0m"
}

func Red(s string) string {
	return "\033[0;31m" + s + "\033[0m"
}

func Green(s string) string {
	return "\033[0;32m" + s + "\033[0m"
}

func Blue(s string) string {
	return "\033[0;34m" + s + "\033[0m"
}
