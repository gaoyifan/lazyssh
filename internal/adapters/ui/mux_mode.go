// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

type muxMode string

const (
	muxModeHerdr   muxMode = "herdr"
	muxModeTmuxCC  muxMode = "tmux-cc"
	muxModeOff     muxMode = "off"
	defaultMuxMode         = muxModeTmuxCC
)

func (m muxMode) valid() bool {
	return m == muxModeHerdr || m == muxModeTmuxCC || m == muxModeOff
}

func (m muxMode) String() string {
	if m == muxModeTmuxCC {
		return "tmux -CC"
	}
	return string(m)
}

func (m muxMode) next() muxMode {
	switch m {
	case muxModeHerdr:
		return muxModeTmuxCC
	case muxModeTmuxCC:
		return muxModeOff
	case muxModeOff:
		return muxModeHerdr
	}
	return defaultMuxMode
}
