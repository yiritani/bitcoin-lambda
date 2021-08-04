package bitflyer

type ProductCode int
const (
	BtcJpy ProductCode = iota
	Ethjpy
	Fxbtcjpy
	Ethbtc
	bchbtc
)

func (code ProductCode) String() string{
	switch code {
	case BtcJpy:
		return "BTC_JPY"
	case Ethjpy:
		return "ETH_JPY"
	case Fxbtcjpy:
		return "FX_BTC_JPY"
	case Ethbtc:
		return "ETH_BTC"
	case bchbtc:
		return "BCH_BTC"
	default:
		return "BTC_JPY"
	}
}