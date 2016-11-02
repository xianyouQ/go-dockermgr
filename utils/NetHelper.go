package utils
import (
    "net"
    "fmt"
    "errors"
)

type IpHelper struct {
    Ip net.IP
    IpInt uint32
}

type CidrHelper struct  {
   Net net.IPNet
   StartIp IpHelper
   EndIp IpHelper
}


func NewIpfromString(ip string) (IpHelper,error) {
    var resultIp IpHelper
    Ip := net.ParseIP(ip).To4()
    if Ip == nil {
        return resultIp,errors.New("valid ip string")
    }
    var sum uint32
    sum += uint32(Ip[0]) << 24
    sum += uint32(Ip[1]) << 16
    sum += uint32(Ip[2]) << 8
    sum += uint32(Ip[3])
    resultIp.Ip = Ip
    resultIp.IpInt = sum
    return resultIp,nil
}


func NewIPv4fromUint(ip uint32) net.IP {
    var bytes [4]byte
    bytes[0] = byte(ip & 0xFF)
    bytes[1] = byte((ip >> 8) & 0xFF)
    bytes[2] = byte((ip >> 16) & 0xFF)
    bytes[3] = byte((ip >> 24) & 0xFF)
    Ip := net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0])
    return Ip
}
func NewIpfromUint(ip uint32) IpHelper {
    var resultIp IpHelper
    resultIp = IpHelper{}
    Ip := NewIPv4fromUint(ip)
    resultIp.Ip = Ip
    resultIp.IpInt = ip
    return resultIp
}

func NewIpfromIp(ip net.IP) IpHelper {
    ip = ip.To4()
    var resultIp IpHelper
    var sum uint32
    sum += uint32(ip[0]) << 24
    sum += uint32(ip[1]) << 16
    sum += uint32(ip[2]) << 8
    sum += uint32(ip[3])
    resultIp.Ip = ip
    resultIp.IpInt = sum
    return resultIp
}

func (self IpHelper) String() string {
    return self.Ip.String()
}

//generate no ieee mac for container
func (self IpHelper)GetMacAddr(base string) string {
    result := fmt.Sprintf("%s:%x:%x:%x",base,self.Ip[1],self.Ip[2],self.Ip[3])
    return result
}
//
func NewCidrfromString(netstr string) (CidrHelper,error) {
    var resultCidr CidrHelper
    resultCidr = CidrHelper{}
     startIp,ipNet,err := net.ParseCIDR(netstr)
     if err != nil {
         return resultCidr,err
     }
     resultCidr.Net = *ipNet
     resultCidr.StartIp = NewIpfromIp(startIp)
     lastIpInt := NewIpfromIp(resultCidr.Net.IP).IpInt + uint32(resultCidr.Size()) - 1
     lastIp := NewIpfromUint(uint32(lastIpInt))
     if err != nil {

     }
     resultCidr.EndIp = lastIp
     return resultCidr,nil
}

func NewCidrwithStartEnd(netstr string,start string,end string) (CidrHelper,error) {
     resultCidr,err := NewCidrfromString(netstr)
     if err != nil {
         return resultCidr,err
     }
     startPar,err := NewIpfromString(start)
     endPar,err := NewIpfromString(end)
     if endPar.IpInt < startPar.IpInt {
         return resultCidr,errors.New("start Ip cant not greater than end Ip")
     }
     if startPar.IpInt >= resultCidr.StartIp.IpInt && startPar.IpInt <= resultCidr.EndIp.IpInt {
         resultCidr.StartIp = startPar
     }
     if endPar.IpInt >= resultCidr.StartIp.IpInt && startPar.IpInt <= resultCidr.EndIp.IpInt {
         resultCidr.EndIp = endPar
     }
     return resultCidr,nil
}

func (self CidrHelper) Size() int {
    ones,bits := self.Net.Mask.Size()
    sum := 1 << (uint(bits-ones))
    return sum
}
func (self CidrHelper) Overlaps(oths CidrHelper) bool {
    if oths.StartIp.IpInt > self.EndIp.IpInt || self.StartIp.IpInt > oths.EndIp.IpInt {
        return false
    }
    return true
}
func (self CidrHelper) Subnet(oths CidrHelper) bool {
    if self.StartIp.IpInt <= oths.StartIp.IpInt && self.EndIp.IpInt >= oths.EndIp.IpInt {
        return true
    }
    return false
}

func (self CidrHelper) String() string {
    return self.Net.String()
}
func (self CidrHelper) IpList() []net.IP {
    Ips := make([]net.IP,0,125)
    for ipInt := self.StartIp.IpInt ; ipInt <= self.EndIp.IpInt; ipInt = ipInt + 1{
        Ip := NewIPv4fromUint(uint32(ipInt))
        Ips = append(Ips,Ip)
    }
    return Ips
}

func (self CidrHelper) IpHelperList() []IpHelper {
    Ips := make([]IpHelper,0,125)
    for ipInt := self.StartIp.IpInt ; ipInt <= self.EndIp.IpInt; ipInt = ipInt + 1{
        Ip := NewIpfromUint(uint32(ipInt))
        Ips = append(Ips,Ip)
    }
    return Ips
}
