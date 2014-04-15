package oauth

import (
	"fmt"
)

var providers = make(map[SocialType]Provider)
var providersByPath = make(map[string]Provider)

// Register the social provider
func RegisterProvider(prov Provider) error {
	typ := prov.GetType()
	if !typ.Available() {
		return fmt.Errorf("Unknown social type `%d`", typ)
	}
	path := prov.GetPath()
	//if providersByPath[path] != nil {
	//	return fmt.Errorf("path `%s` is already in used", path)
	//}
	providers[typ] = prov
	providersByPath[path] = prov
	return nil
}

// Get provider by SocialType
func GetProviderByType(typ SocialType) (Provider, bool) {
	if p, ok := providers[typ]; ok {
		return p, true
	} else {
		return nil, false
	}
}

// Get provider by path name
func GetProviderByPath(path string) (Provider, bool) {
	if p, ok := providersByPath[path]; ok {
		return p, true
	} else {
		return nil, false
	}
}

func GetProviders() (ret []Provider) {
	for _, v := range providers {
		ret = append(ret, v)
	}
	return ret
}
