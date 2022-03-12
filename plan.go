package tentez

func getTargetNames(targets map[string]Targets) []string {
	targetNames := []string{}

	for _, targetResources := range targets {
		for _, target := range targetResources.targetsSlice() {
			targetNames = append(targetNames, target.getName())
		}
	}

	return targetNames
}
