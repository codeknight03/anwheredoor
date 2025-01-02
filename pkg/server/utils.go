package server

import "sort"

func addToPathHandlerMap(currSlice []pathHandlerPair, path string, handler *BackendHandler) []pathHandlerPair {

	newPair := pathHandlerPair{
		path:    path,
		handler: handler,
	}

	if currSlice == nil || len(currSlice) == 0 {
		return []pathHandlerPair{newPair}
	}

	currSlice = append(currSlice, newPair)

	sort.Slice(currSlice, func(i, j int) bool {
		return len(currSlice[i].path) > len(currSlice[j].path)
	})

	return currSlice

}
