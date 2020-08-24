package types

import "github.com/nicpottier/decent/parser"

// init initializes all our types with our parser library
func init() {
	parser.RegisterReader("M", ReadShotSample)
	parser.RegisterReader("Q", ReadWaterLevels)
	parser.RegisterReader("N", ReadStateInfo)
	parser.RegisterReader("P", ReadShotFrame)
	parser.RegisterReader("K", ReadShotSettings)
	parser.RegisterReader("O", ReadShotDescHeader)
}
