package sandpay

// ClientOption 客户端配置项
type SandV5Option func(c *sandV5Client)

// merNo, prvKeyFile, pubKeyFile, notifyUrl, returnUrl string
func WithMerNo(merNo string) SandV5Option {
	return func(c *sandV5Client) {
		c.merNo = merNo
	}
}

func WithPrvKeyFile(file string) SandV5Option {
	return func(c *sandV5Client) {
		c.prvKeyFile = file
	}
}

func WithPubKeyFile(file string) SandV5Option {
	return func(c *sandV5Client) {
		c.pubKeyFile = file
	}
}

func WithNotifyUrl(file string) SandV5Option {
	return func(c *sandV5Client) {
		c.notifyUrl = file
	}
}

func WithReturnUrl(file string) SandV5Option {
	return func(c *sandV5Client) {
		c.returnUrl = file
	}
}

func WithUserFlag(flag string) SandV5Option {
	return func(c *sandV5Client) {
		c.userFlag = flag
	}
}
