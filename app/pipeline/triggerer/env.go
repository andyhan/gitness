// Copyright 2023 Harness, Inc.
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

package triggerer

import (
	"github.com/harness/gitness/app/url"
	"github.com/harness/gitness/types"

	"golang.org/x/exp/maps"
)

// combine is a helper function that combines one or more maps of
// environment variables into a single map.
func combine(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		maps.Copy(c, e)
	}
	return c
}

func Envs(
	repo *types.Repository,
	pipeline *types.Pipeline,
	urlProvider url.Provider,
) map[string]string {
	return map[string]string{
		"DRONE_BUILD_LINK": urlProvider.GenerateUIBuildURL(repo.Path, pipeline.UID, pipeline.Seq),
	}
}