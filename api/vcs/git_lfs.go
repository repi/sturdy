
	"getsturdy.com/api/pkg/codebases"
func (r *repository) LargeFilesClean(codebaseID codebases.ID, paths []string) ([][]byte, error) {
func (r *repository) configureLfs(codebaseID codebases.ID) error {