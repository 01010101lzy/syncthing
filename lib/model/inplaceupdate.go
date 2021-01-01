package model

import "github.com/syncthing/syncthing/lib/protocol"

// Calculate whether the file should be updated in-place from current file status.
func shouldDoInPlaceUpdate(state *sharedPullerState) bool {
	/*
		The following aspects should be considered:

		- Should this file be created or deleted?
			(no if true)
		- How many blocks need to be updated?
			(no if it needs more space to store the delta than creating a new file)
	*/

	state.mut.Lock()
	defer state.mut.Unlock()

	if !state.hasCurFile {
		return false
	}

	// Delta file size is estimated as:
	//   delta_file_size_blocks = num_appending_or_deleting_blocks + num_replacing_blocks * 2
	// because for replacing blocks we need twice the space to store old data
	// to support rollback.
	//
	// If one day variable-sized blocks are supported, this function should
	// be updated to using true block sizes instead of block counts.

	curLen := len(state.curFile.Blocks)
	newLen := len(state.file.Blocks)
	appendingBlocks := newLen - curLen
	// Is there really no Math.Abs(int) in this language, or am I missing something?
	if appendingBlocks < 0 {
		appendingBlocks = -appendingBlocks
	}
	replacingBlocks := curLen - state.copyTotal

	if (replacingBlocks*2 + appendingBlocks) > newLen {
		return false
	}
	return true
}

// Perform in-place update for the given file.
func doInPlaceUpdate() {
	// it performs the following steps:
	// 1. calculate update commands, write into a command file
	// 2. lock the destination file from writes, if possible
	//    lock other files we need to read from from writes, if possible
	// 3. pull blocks from other devices, write into a delta file
	// 4. perform the in-place update using the given commands
	// 5. release locks

}

type inPlaceUpdateDeltaFile struct {
	copies         []inPlaceCopyCommand            // topologically sorted in-file copy commands
	externalCopies map[string][]inPlaceCopyCommand // blocks to copy from other files
	pulls          []inPlacePullCommand            // blocks to be pulled from other devices
}

type inPlaceCopyCommand struct {
	sourceBlock int
	destBlock   int
}

type inPlacePullCommand struct {
	sourceBlock protocol.BlockInfo
	destBlock   int
}

func calculateDeltaFile(state *sharedPullerState) {

}
