package sandpay

// ClientOption 客户端配置项
type SandV4Option func(c *sandV4Client)

// merNo, prvKeyFile, pubKeyFile, notifyUrl, returnUrl string
func WithV4MerNo(merNo string) SandV4Option {
	return func(c *sandV4Client) {
		c.merNo = merNo
	}
}

func WithV4PrvKeyFile(file string) SandV4Option {
	return func(c *sandV4Client) {
		c.prvKeyFile = file
	}
}

func WithV4PubKeyFile(file string) SandV4Option {
	return func(c *sandV4Client) {
		c.pubKeyFile = file
	}
}

func WithV4NotifyUrl(file string) SandV4Option {
	return func(c *sandV4Client) {
		c.notifyUrl = file
	}
}

func WithV4ReturnUrl(file string) SandV4Option {
	return func(c *sandV4Client) {
		c.returnUrl = file
	}
}

func WithV4UserFlag(flag string) SandV4Option {
	return func(c *sandV4Client) {
		c.userFlag = flag
	}
}
