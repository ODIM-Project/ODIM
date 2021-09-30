//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

// Package common ...
package common

func getOutChannel(buff []interface{}, out chan interface{}) chan interface{} {
	if len(buff) == 0 {
		return nil
	}
	return out
}

func getCurrentData(buff []interface{}) interface{} {
	if len(buff) == 0 {
		return nil
	}
	return buff[0]
}

// CreateJobQueue defines the queue which will act as an infinite buffer.
// We can make use of this when we are not sure about the capacity
// required for the creating a buffered channel.
//
// The function will return two channels; the first one will act as
// an entry or input channel and the later will be an exit or output channel.
//
// The in channel must be closed by the user, but the closing of
// the out channel is taken care by the function itself.
func CreateJobQueue(qsize int) (chan<- interface{}, <-chan interface{}) {
	in, out := make(chan interface{}, qsize), make(chan interface{}, qsize)
	go func() {
		var buff []interface{}
		// for loop must continue until buff is empty (all data are read) and the in channel is closed
		for len(buff) > 0 || in != nil {
			select {
			case data, ok := <-in: // new data arrives through in channel
				if !ok { // in channel is closed
					in = nil
				} else {
					if data != nil {
						buff = append(buff, data)
					}
				}
			case getOutChannel(buff, out) <- getCurrentData(buff): // out channel is ready for new data from buff
				buff = buff[1:]
			}
		}
		close(out)
	}()
	return in, out
}
