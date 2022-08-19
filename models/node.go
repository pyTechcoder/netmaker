package models

import (
	"bytes"
	"math/rand"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const TEN_YEARS_IN_SECONDS = 300000000
const MAX_NAME_LENGTH = 62

// == ACTIONS == (can only be set by GRPC)
const NODE_UPDATE_KEY = "updatekey"
const NODE_SERVER_NAME = "netmaker"
const NODE_DELETE = "delete"
const NODE_IS_PENDING = "pending"
const NODE_NOOP = "noop"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Node - struct for node model
type Node struct {
	ID                  string   `json:"id,omitempty" bson:"id,omitempty" yaml:"id,omitempty" validate:"required,min=5"`
	Address             string   `json:"address" bson:"address" yaml:"address" validate:"omitempty,ipv4"`
	Address6            string   `json:"address6" bson:"address6" yaml:"address6" validate:"omitempty,ipv6"`
	LocalAddress        string   `json:"localaddress" bson:"localaddress" yaml:"localaddress" validate:"omitempty,ip"`
	Name                string   `json:"name" bson:"name" yaml:"name" validate:"omitempty,max=62,in_charset"`
	NetworkSettings     Network  `json:"networksettings" bson:"networksettings" yaml:"networksettings" validate:"-"`
	ListenPort          int32    `json:"listenport" bson:"listenport" yaml:"listenport" validate:"omitempty,numeric,min=1024,max=65535"`
	PublicKey           string   `json:"publickey" bson:"publickey" yaml:"publickey" validate:"required,base64"`
	Endpoint            string   `json:"endpoint" bson:"endpoint" yaml:"endpoint" validate:"required,ip"`
	PostUp              string   `json:"postup" bson:"postup" yaml:"postup"`
	PostDown            string   `json:"postdown" bson:"postdown" yaml:"postdown"`
	AllowedIPs          []string `json:"allowedips" bson:"allowedips" yaml:"allowedips"`
	PersistentKeepalive int32    `json:"persistentkeepalive" bson:"persistentkeepalive" yaml:"persistentkeepalive" validate:"omitempty,numeric,max=1000"`
	SaveConfig          string   `json:"saveconfig" bson:"saveconfig" yaml:"saveconfig" validate:"checkyesorno"`
	AccessKey           string   `json:"accesskey" bson:"accesskey" yaml:"accesskey"`
	Interface           string   `json:"interface" bson:"interface" yaml:"interface"`
	LastModified        int64    `json:"lastmodified" bson:"lastmodified" yaml:"lastmodified"`
	KeyUpdateTimeStamp  int64    `json:"keyupdatetimestamp" bson:"keyupdatetimestamp" yaml:"keyupdatetimestamp"`
	ExpirationDateTime  int64    `json:"expdatetime" bson:"expdatetime" yaml:"expdatetime"`
	LastPeerUpdate      int64    `json:"lastpeerupdate" bson:"lastpeerupdate" yaml:"lastpeerupdate"`
	LastCheckIn         int64    `json:"lastcheckin" bson:"lastcheckin" yaml:"lastcheckin"`
	MacAddress          string   `json:"macaddress" bson:"macaddress" yaml:"macaddress" validate:"macaddress_unique"`
	// checkin interval is depreciated at the network level. Set on server with CHECKIN_INTERVAL
	CheckInInterval     int32       `json:"checkininterval" bson:"checkininterval" yaml:"checkininterval"`
	Password            string      `json:"password" bson:"password" yaml:"password" validate:"required,min=6"`
	Network             string      `json:"network" bson:"network" yaml:"network" validate:"network_exists"`
	IsRelayed           string      `json:"isrelayed" bson:"isrelayed" yaml:"isrelayed"`
	IsPending           string      `json:"ispending" bson:"ispending" yaml:"ispending"`
	IsRelay             string      `json:"isrelay" bson:"isrelay" yaml:"isrelay" validate:"checkyesorno"`
	IsDocker            string      `json:"isdocker" bson:"isdocker" yaml:"isdocker" validate:"checkyesorno"`
	IsK8S               string      `json:"isk8s" bson:"isk8s" yaml:"isk8s" validate:"checkyesorno"`
	IsEgressGateway     string      `json:"isegressgateway" bson:"isegressgateway" yaml:"isegressgateway"`
	IsIngressGateway    string      `json:"isingressgateway" bson:"isingressgateway" yaml:"isingressgateway"`
	EgressGatewayRanges []string    `json:"egressgatewayranges" bson:"egressgatewayranges" yaml:"egressgatewayranges"`
	RelayAddrs          []string    `json:"relayaddrs" bson:"relayaddrs" yaml:"relayaddrs"`
	IngressGatewayRange string      `json:"ingressgatewayrange" bson:"ingressgatewayrange" yaml:"ingressgatewayrange"`
	IsStatic            string      `json:"isstatic" bson:"isstatic" yaml:"isstatic" validate:"checkyesorno"`
	UDPHolePunch        string      `json:"udpholepunch" bson:"udpholepunch" yaml:"udpholepunch" validate:"checkyesorno"`
	PullChanges         string      `json:"pullchanges" bson:"pullchanges" yaml:"pullchanges" validate:"checkyesorno"`
	DNSOn               string      `json:"dnson" bson:"dnson" yaml:"dnson" validate:"checkyesorno"`
	IsDualStack         string      `json:"isdualstack" bson:"isdualstack" yaml:"isdualstack" validate:"checkyesorno"`
	IsServer            string      `json:"isserver" bson:"isserver" yaml:"isserver" validate:"checkyesorno"`
	Action              string      `json:"action" bson:"action" yaml:"action"`
	IsLocal             string      `json:"islocal" bson:"islocal" yaml:"islocal" validate:"checkyesorno"`
	LocalRange          string      `json:"localrange" bson:"localrange" yaml:"localrange"`
	Roaming             string      `json:"roaming" bson:"roaming" yaml:"roaming" validate:"checkyesorno"`
	IPForwarding        string      `json:"ipforwarding" bson:"ipforwarding" yaml:"ipforwarding" validate:"checkyesorno"`
	OS                  string      `json:"os" bson:"os" yaml:"os"`
	MTU                 int32       `json:"mtu" bson:"mtu" yaml:"mtu"`
	Version             string      `json:"version" bson:"version" yaml:"version"`
	ExcludedAddrs       []string    `json:"excludedaddrs" bson:"excludedaddrs" yaml:"excludedaddrs"`
	TrafficKeys         TrafficKeys `json:"traffickeys" bson:"traffickeys" yaml:"traffickeys"`
}

// NodesArray - used for node sorting
type NodesArray []Node

// NodesArray.Len - gets length of node array
func (a NodesArray) Len() int { return len(a) }

// NodesArray.Less - gets returns lower rank of two node addresses
func (a NodesArray) Less(i, j int) bool { return isLess(a[i].Address, a[j].Address) }

// NodesArray.Swap - swaps two nodes in array
func (a NodesArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func isLess(ipA string, ipB string) bool {
	ipNetA := net.ParseIP(ipA)
	ipNetB := net.ParseIP(ipB)
	return bytes.Compare(ipNetA, ipNetB) < 0
}

// Node.SetDefaultMTU - sets default MTU of a node
func (node *Node) SetDefaultMTU() {
	if node.MTU == 0 {
		node.MTU = 1280
	}
}

// Node.SetDefaulIsPending - sets ispending default
func (node *Node) SetDefaulIsPending() {
	if node.IsPending == "" {
		node.IsPending = "no"
	}
}

// Node.SetDefaultIsRelayed - set default is relayed
func (node *Node) SetDefaultIsRelayed() {
	if node.IsRelayed == "" {
		node.IsRelayed = "no"
	}
}

// Node.SetDefaultIsRelay - set default isrelay
func (node *Node) SetDefaultIsRelay() {
	if node.IsRelay == "" {
		node.IsRelay = "no"
	}
}

// Node.SetDefaultIsDocker - set default isdocker
func (node *Node) SetDefaultIsDocker() {
	if node.IsDocker == "" {
		node.IsDocker = "no"
	}
}

// Node.SetDefaultIsK8S - set default isk8s
func (node *Node) SetDefaultIsK8S() {
	if node.IsK8S == "" {
		node.IsK8S = "no"
	}
}

// Node.SetDefaultEgressGateway - sets default egress gateway status
func (node *Node) SetDefaultEgressGateway() {
	if node.IsEgressGateway == "" {
		node.IsEgressGateway = "no"
	}
}

// Node.SetDefaultIngressGateway - sets default ingress gateway status
func (node *Node) SetDefaultIngressGateway() {
	if node.IsIngressGateway == "" {
		node.IsIngressGateway = "no"
	}
}

// Node.SetDefaultAction - sets default action status
func (node *Node) SetDefaultAction() {
	if node.Action == "" {
		node.Action = NODE_NOOP
	}
}

// Node.SetRoamingDefault - sets default roaming status
func (node *Node) SetRoamingDefault() {
	if node.Roaming == "" {
		node.Roaming = "yes"
	}
}

// Node.SetPullChangesDefault - sets default pull changes status
func (node *Node) SetPullChangesDefault() {
	if node.PullChanges == "" {
		node.PullChanges = "no"
	}
}

// Node.SetIPForwardingDefault - set ip forwarding default
func (node *Node) SetIPForwardingDefault() {
	if node.IPForwarding == "" {
		node.IPForwarding = "yes"
	}
}

// Node.SetIsLocalDefault - set is local default
func (node *Node) SetIsLocalDefault() {
	if node.IsLocal == "" {
		node.IsLocal = "no"
	}
}

// Node.SetDNSOnDefault - sets dns on default
func (node *Node) SetDNSOnDefault() {
	if node.DNSOn == "" {
		node.DNSOn = "yes"
	}
}

// Node.SetIsDualStackDefault - set is dual stack default status
func (node *Node) SetIsDualStackDefault() {
	if node.IsDualStack == "" {
		node.IsDualStack = "no"
	}
}

// Node.SetIsServerDefault - sets node isserver default
func (node *Node) SetIsServerDefault() {
	if node.IsServer != "yes" {
		node.IsServer = "no"
	}
}

// Node.SetIsStaticDefault - set is static default
func (node *Node) SetIsStaticDefault() {
	if node.IsServer == "yes" {
		node.IsStatic = "yes"
	} else if node.IsStatic != "yes" {
		node.IsStatic = "no"
	}
}

// Node.SetLastModified - set last modified initial time
func (node *Node) SetLastModified() {
	node.LastModified = time.Now().Unix()
}

// Node.SetLastCheckIn - time.Now().Unix()
func (node *Node) SetLastCheckIn() {
	node.LastCheckIn = time.Now().Unix()
}

// Node.SetLastPeerUpdate - sets last peer update time
func (node *Node) SetLastPeerUpdate() {
	node.LastPeerUpdate = time.Now().Unix()
}

// Node.SetExpirationDateTime - sets node expiry time
func (node *Node) SetExpirationDateTime() {
	node.ExpirationDateTime = time.Now().Unix() + TEN_YEARS_IN_SECONDS
}

// Node.SetDefaultName - sets a random name to node
func (node *Node) SetDefaultName() {
	if node.Name == "" {
		node.Name = GenerateNodeName()
	}
}

// Node.SetDefaultExcludedAddrs - sets ExcludedAddrs to empty array if nil
func (node *Node) SetDefaultExcludedAddrs() {
	if node.ExcludedAddrs == nil {
		node.ExcludedAddrs = make([]string, 0)
	}
}

// Node.Fill - fills other node data into calling node data if not set on calling node
func (newNode *Node) Fill(currentNode *Node) {
	newNode.ID = currentNode.ID

	if newNode.Address == "" && newNode.IsStatic != "yes" {
		newNode.Address = currentNode.Address
	}
	if newNode.Address6 == "" && newNode.IsStatic != "yes" {
		newNode.Address6 = currentNode.Address6
	}
	if newNode.LocalAddress == "" {
		newNode.LocalAddress = currentNode.LocalAddress
	}
	if newNode.Name == "" {
		newNode.Name = currentNode.Name
	}
	if newNode.ListenPort == 0 && newNode.IsStatic != "yes" {
		newNode.ListenPort = currentNode.ListenPort
	}
	if newNode.PublicKey == "" && newNode.IsStatic != "yes" {
		newNode.PublicKey = currentNode.PublicKey
	} else {
		newNode.KeyUpdateTimeStamp = time.Now().Unix()
	}
	if newNode.Endpoint == "" && newNode.IsStatic != "yes" {
		newNode.Endpoint = currentNode.Endpoint
	}
	if newNode.PostUp == "" {
		newNode.PostUp = currentNode.PostUp
	}
	if newNode.PostDown == "" {
		newNode.PostDown = currentNode.PostDown
	}
	if newNode.AllowedIPs == nil {
		newNode.AllowedIPs = currentNode.AllowedIPs
	}
	if newNode.PersistentKeepalive == 0 {
		newNode.PersistentKeepalive = currentNode.PersistentKeepalive
	}
	if newNode.SaveConfig == "" {
		newNode.SaveConfig = currentNode.SaveConfig
	}
	if newNode.AccessKey == "" {
		newNode.AccessKey = currentNode.AccessKey
	}
	if newNode.Interface == "" {
		newNode.Interface = currentNode.Interface
	}
	if newNode.LastModified == 0 {
		newNode.LastModified = currentNode.LastModified
	}
	if newNode.KeyUpdateTimeStamp == 0 {
		newNode.LastModified = currentNode.LastModified
	}
	if newNode.ExpirationDateTime == 0 {
		newNode.ExpirationDateTime = currentNode.ExpirationDateTime
	}
	if newNode.LastPeerUpdate == 0 {
		newNode.LastPeerUpdate = currentNode.LastPeerUpdate
	}
	if newNode.LastCheckIn == 0 {
		newNode.LastCheckIn = currentNode.LastCheckIn
	}
	if newNode.MacAddress == "" {
		newNode.MacAddress = currentNode.MacAddress
	}
	if newNode.CheckInInterval == 0 {
		newNode.CheckInInterval = currentNode.CheckInInterval
	}
	if newNode.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(newNode.Password), []byte(currentNode.Password))
		if err != nil && currentNode.Password != newNode.Password {
			hash, err := bcrypt.GenerateFromPassword([]byte(newNode.Password), 5)
			if err == nil {
				newNode.Password = string(hash)
			}
		}
	} else {
		newNode.Password = currentNode.Password
	}
	if newNode.Network == "" {
		newNode.Network = currentNode.Network
	}
	if newNode.IsPending == "" {
		newNode.IsPending = currentNode.IsPending
	}
	if newNode.IsEgressGateway == "" {
		newNode.IsEgressGateway = currentNode.IsEgressGateway
	}
	if newNode.IsIngressGateway == "" {
		newNode.IsIngressGateway = currentNode.IsIngressGateway
	}
	if newNode.EgressGatewayRanges == nil {
		newNode.EgressGatewayRanges = currentNode.EgressGatewayRanges
	}
	if newNode.IngressGatewayRange == "" {
		newNode.IngressGatewayRange = currentNode.IngressGatewayRange
	}
	if newNode.IsStatic == "" {
		newNode.IsStatic = currentNode.IsStatic
	}
	if newNode.UDPHolePunch == "" {
		newNode.UDPHolePunch = currentNode.SaveConfig
	}
	if newNode.DNSOn == "" {
		newNode.DNSOn = currentNode.DNSOn
	}
	if newNode.IsDualStack == "" {
		newNode.IsDualStack = currentNode.IsDualStack
	}
	if newNode.IsLocal == "" {
		newNode.IsLocal = currentNode.IsLocal
	}
	if newNode.IPForwarding == "" {
		newNode.IPForwarding = currentNode.IPForwarding
	}
	if newNode.PullChanges == "" {
		newNode.PullChanges = currentNode.PullChanges
	}
	if newNode.Roaming == "" {
		newNode.Roaming = currentNode.Roaming
	}
	if newNode.Action == "" {
		newNode.Action = currentNode.Action
	}
	if newNode.IsServer == "" {
		newNode.IsServer = currentNode.IsServer
	}
	if newNode.IsServer == "yes" {
		newNode.IsStatic = "yes"
	}
	if newNode.MTU == 0 {
		newNode.MTU = currentNode.MTU
	}
	if newNode.OS == "" {
		newNode.OS = currentNode.OS
	}
	if newNode.RelayAddrs == nil {
		newNode.RelayAddrs = currentNode.RelayAddrs
	}
	if newNode.IsRelay == "" {
		newNode.IsRelay = currentNode.IsRelay
	}
	if newNode.IsRelayed == "" {
		newNode.IsRelayed = currentNode.IsRelayed
	}
	if newNode.IsDocker == "" {
		newNode.IsDocker = currentNode.IsDocker
	}
	if newNode.IsK8S == "" {
		newNode.IsK8S = currentNode.IsK8S
	}
	if newNode.Version == "" {
		newNode.Version = currentNode.Version
	}
	if newNode.ExcludedAddrs == nil || len(newNode.ExcludedAddrs) != len(currentNode.ExcludedAddrs) {
		newNode.ExcludedAddrs = currentNode.ExcludedAddrs
	}
}

// StringWithCharset - returns random string inside defined charset
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// IsIpv4Net - check for valid IPv4 address
// Note: We dont handle IPv6 AT ALL!!!!! This definitely is needed at some point
// But for iteration 1, lets just stick to IPv4. Keep it simple stupid.
func IsIpv4Net(host string) bool {
	return net.ParseIP(host) != nil
}

// Node.NameInNodeCharset - returns if name is in charset below or not
func (node *Node) NameInNodeCharSet() bool {

	charset := "abcdefghijklmnopqrstuvwxyz1234567890-"

	for _, char := range node.Name {
		if !strings.Contains(charset, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}
