package internal

const (
	MspID         = "Org1MSP"
	CryptoPath    = "../organizations/peerOrganizations/org1.example.com"
	CertPath      = CryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	KeyPath       = CryptoPath + "/users/User1@org1.example.com/msp/keystore"
	TlsCertPath   = CryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	PeerEndpoint  = "dns:///localhost:7051"
	GatewayPeer   = "peer0.org1.example.com"
	ChaincodeName = "basic"
	ChannelName   = "mychannel"
)
