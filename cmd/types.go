package cmd

type svcMap struct {
	SvcName       string
	SvcPort       int32
	SvcSolacePort int32
	SvcNodePort   int32
}

type HaproxyCRD struct {
	Haproxy Haproxy `json:"haproxy"`
}
type Haproxy struct {
	Namespace string  `json:"namespace"`
	Publish   svcName `json:"publish"`
	Subscribe svcName `json:"subscribe"`
}
type svcName struct {
	ServiceName string `json:"serviceName"`
}
