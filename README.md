# rygelXVI
![Dominar Rygel XVI] (https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcQCymqxYzLXe0em0f9Jvss35DL1VxVdgC_Y8LUKlSwHE0nKHB8H6w)

rygelXVI is a GO implementation of the Mashape Domainr API version 2. 

## Installation
```golang
go get github.com/Dystopi/rygelXVI
```

## Overview

This package is a wrapper around the Mashape implementation of the Domainr API. We currently only implement the /v2/search and /v2/status endpoints.

## Usage

```golang
import (
	"fmt"
	"github.com/Dystopi/rygelXVI"
)

client, _ := rygelXVI.NewClient({Super_Awesome_API_Key}, nil) // For testing you can pass your own http.Client, or leave nil for a standard http.Client
domains, _ := client.SearchDomainr("foo.com")
for _, d := range domains.Results {
	fmt.Println(d.Domain)
}

activeDomains, _ := client.SearchActive("bar.com")
for _, d := range activeDomains.Results {
	fmt.Println(d.Domain)
}

status, _ := client.DomainrStatus([]string{"biz.co"})
isAvailable := status["biz.co"] // bool
```
