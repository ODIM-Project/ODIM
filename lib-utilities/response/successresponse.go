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

package response

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// CreateGenericResponse will fill the response with respective data
func (r *Response) CreateGenericResponse(statusMessage string) {

	r.MessageID = statusMessage
	r.Severity = "OK"

	switch statusMessage {
	case Success:
		r.Message = "Successfully Completed Request"
	case Created:
		r.Message = "The resource has been created successfully"
	case AccountRemoved:
		r.Message = "The account was successfully removed."
	case AccountModified:
		r.Message = "The account was successfully modified."
	case ResourceRemoved:
		r.Message = "The resource has been removed successfully."
	case ResourceCreated:
		r.Message = "The resource has been created successfully."
	case TaskStarted:
		r.NumberOfArgs = len(r.MessageArgs)
		if r.NumberOfArgs < 1 {
			log.Warn("MessageArgs in Response is missing")
		}
		r.Message = fmt.Sprintf("The task with id %v has started.", r.MessageArgs[0])
	}

}
