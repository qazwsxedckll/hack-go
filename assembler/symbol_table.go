package main

type SymbolTable struct{}

func (s *SymbolTable) AddEntry(symbol string, address int) {
}

func (s *SymbolTable) Contains(symbol string) bool {
	return false
}

func (s *SymbolTable) GetAddress(symbol string) int {
	return 0
}
