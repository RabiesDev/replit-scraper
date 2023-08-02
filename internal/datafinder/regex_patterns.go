package datafinder

import "regexp"

var (
	DiscordTokenPattern   = regexp.MustCompile("[\\w-]{24,}\\.[\\w-]{6}\\.[\\w-]{27,38}")
	EmailPattern          = regexp.MustCompile("(?:[a-z0-9!#$%&'*+=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])")
	PasswordPattern       = regexp.MustCompile("(?i)(?U)password[\\s\\W]{1,5}\"(\\S{4,24})\"")
	PhonePattern          = regexp.MustCompile("((0[7-9]0)-?([1-9]{4})-?([1-9]{4}))")
	PaidProxyPattern      = regexp.MustCompile("(?:((http|https|socks5?)://)|\")(\\S+):(\\w+|\\d{5})@(\\S+):(\\w+|\\d{5})\"?")
	CaptchaServicePattern = regexp.MustCompile("(?i)(?U)((twocaptcha|capmonster)\\W[\"|'](.*?)[\"|']\\W)")
	OpenAiKeyPattern      = regexp.MustCompile("(sk-[a-zA-Z0-9]{48,})")
	GoogleApiKeyPattern   = regexp.MustCompile("AIzaSyC\\w{32}")
	TelegramTokenPattern  = regexp.MustCompile("(\\d{10}:[a-zA-Z0-9]{35})")
)
